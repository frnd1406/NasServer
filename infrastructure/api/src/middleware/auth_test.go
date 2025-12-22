package middleware

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthMiddleware(t *testing.T) (*services.JWTService, *database.RedisClient, *miniredis.Miniredis, *logrus.Logger) {
	ensureTCPAllowed(t)

	// Setup JWT Service
	cfg := &config.Config{
		JWTSecret: "test-secret-key-minimum-32-characters-long-for-security",
	}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	jwtService, err := services.NewJWTService(cfg, logger)
	require.NoError(t, err)

	// Setup miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	redisClient := &database.RedisClient{
		Client: client,
	}

	return jwtService, redisClient, mr, logger
}

func ensureTCPAllowed(t *testing.T) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("TCP sockets nicht erlaubt in dieser Umgebung: %v", err)
		return
	}
	ln.Close()
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService, redisClient, mr, logger := setupAuthMiddleware(t)
	defer mr.Close()

	// Generate valid token
	userID := "test-user-123"
	email := "test@example.com"
	token, err := jwtService.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Setup router with auth middleware
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-123")
		c.Next()
	})
	router.Use(AuthMiddleware(jwtService, redisClient, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with valid token
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should succeed
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService, redisClient, mr, logger := setupAuthMiddleware(t)
	defer mr.Close()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-123")
		c.Next()
	})
	router.Use(AuthMiddleware(jwtService, redisClient, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing authorization token")
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService, redisClient, mr, logger := setupAuthMiddleware(t)
	defer mr.Close()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-123")
		c.Next()
	})
	router.Use(AuthMiddleware(jwtService, redisClient, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "missing Bearer prefix",
			header: "just-a-token",
		},
		{
			name:   "wrong prefix",
			header: "Basic dGVzdDp0ZXN0",
		},
		{
			name:   "no space",
			header: "Bearertoken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid authorization header format")
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService, redisClient, mr, logger := setupAuthMiddleware(t)
	defer mr.Close()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-123")
		c.Next()
	})
	router.Use(AuthMiddleware(jwtService, redisClient, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestAuthMiddleware_BlacklistedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService, redisClient, mr, logger := setupAuthMiddleware(t)
	defer mr.Close()

	// Generate valid token
	userID := "test-user-123"
	email := "test@example.com"
	token, err := jwtService.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	// Blacklist the token
	ctx := context.Background()
	err = redisClient.Set(ctx, "blacklist:"+token, "1", 1*time.Hour).Err()
	require.NoError(t, err)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-123")
		c.Next()
	})
	router.Use(AuthMiddleware(jwtService, redisClient, logger))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 401 for blacklisted token
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token has been revoked")
}

func TestAuthMiddleware_UserContextPropagation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService, redisClient, mr, logger := setupAuthMiddleware(t)
	defer mr.Close()

	userID := "test-user-456"
	email := "user@example.com"
	token, err := jwtService.GenerateAccessToken(userID, email)
	require.NoError(t, err)

	var capturedUserID string
	var capturedEmail string

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-123")
		c.Next()
	})
	router.Use(AuthMiddleware(jwtService, redisClient, logger))
	router.GET("/protected", func(c *gin.Context) {
		capturedUserID = c.GetString("user_id")
		capturedEmail = c.GetString("user_email")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify user context was set correctly
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, userID, capturedUserID)
	assert.Equal(t, email, capturedEmail)
}
