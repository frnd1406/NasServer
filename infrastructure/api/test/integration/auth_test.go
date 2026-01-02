package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ============================================================
// Mock Services for Integration Tests
// ============================================================

type MockUserRepository struct {
	users map[string]map[string]interface{}
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]map[string]interface{}),
	}
}

func (m *MockUserRepository) CreateUser(email, username, passwordHash string) (map[string]interface{}, error) {
	user := map[string]interface{}{
		"id":            "user-" + email,
		"email":         email,
		"username":      username,
		"password_hash": passwordHash,
		"created_at":    time.Now(),
		"verified":      false,
	}
	m.users[email] = user
	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (map[string]interface{}, bool) {
	user, ok := m.users[email]
	return user, ok
}

// ============================================================
// Integration Tests - Auth Flow
// ============================================================

func TestIntegration_RegisterLoginFlow(t *testing.T) {
	router := gin.New()
	userRepo := NewMockUserRepository()

	// Setup mock register endpoint
	router.POST("/auth/register", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		// Check if user exists
		if _, exists := userRepo.FindByEmail(req.Email); exists {
			c.JSON(http.StatusConflict, gin.H{"error": "email_exists"})
			return
		}

		// Password validation
		if len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "weak_password"})
			return
		}

		user, _ := userRepo.CreateUser(req.Email, req.Username, "hashed_"+req.Password)
		c.JSON(http.StatusCreated, gin.H{
			"user":          user,
			"access_token":  "mock-access-token",
			"refresh_token": "mock-refresh-token",
			"csrf_token":    "mock-csrf-token",
		})
	})

	// Setup mock login endpoint
	router.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		user, exists := userRepo.FindByEmail(req.Email)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
			return
		}

		// Check password (mock)
		expectedHash := "hashed_" + req.Password
		if user["password_hash"] != expectedHash {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":          user,
			"access_token":  "mock-access-token-login",
			"refresh_token": "mock-refresh-token-login",
			"csrf_token":    "mock-csrf-token-login",
		})
	})

	// Test 1: Register new user
	t.Run("Register_Success", func(t *testing.T) {
		body := `{"username": "testuser", "email": "test@example.com", "password": "SecurePassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.NotEmpty(t, resp["access_token"])
		assert.NotEmpty(t, resp["csrf_token"])
	})

	// Test 2: Duplicate registration
	t.Run("Register_DuplicateEmail", func(t *testing.T) {
		body := `{"username": "testuser2", "email": "test@example.com", "password": "SecurePassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// Test 3: Login with registered user
	t.Run("Login_Success", func(t *testing.T) {
		body := `{"email": "test@example.com", "password": "SecurePassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.NotEmpty(t, resp["access_token"])
	})

	// Test 4: Login with wrong password
	t.Run("Login_WrongPassword", func(t *testing.T) {
		body := `{"email": "test@example.com", "password": "WrongPassword!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test 5: Login with non-existent user
	t.Run("Login_UserNotFound", func(t *testing.T) {
		body := `{"email": "nonexistent@example.com", "password": "SomePassword123!"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestIntegration_PasswordValidation(t *testing.T) {
	router := gin.New()

	router.POST("/auth/register", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
			return
		}

		// Password strength validation
		if len(req.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "weak_password",
					"message": "Password must be at least 8 characters",
				},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "ok"})
	})

	tests := []struct {
		name       string
		password   string
		wantStatus int
	}{
		{"TooShort", "short", http.StatusBadRequest},
		{"ExactlyMinLength", "12345678", http.StatusCreated},
		{"StrongPassword", "SecurePassword123!", http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]string{
				"username": "testuser",
				"email":    tt.name + "@example.com",
				"password": tt.password,
			})

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestIntegration_CSRFProtection(t *testing.T) {
	router := gin.New()

	// Mock CSRF middleware
	csrfMiddleware := func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing_csrf_token"})
			return
		}

		if csrfToken != "valid-csrf-token" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid_csrf_token"})
			return
		}

		c.Next()
	}

	router.GET("/api/v1/auth/csrf", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"csrf_token": "valid-csrf-token"})
	})

	protected := router.Group("/api/v1")
	protected.Use(csrfMiddleware)
	protected.POST("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	t.Run("GetCSRFToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/csrf", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, "valid-csrf-token", resp["csrf_token"])
	})

	t.Run("ProtectedEndpoint_WithoutCSRF", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("ProtectedEndpoint_WithValidCSRF", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/protected", nil)
		req.Header.Set("X-CSRF-Token", "valid-csrf-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ProtectedEndpoint_WithInvalidCSRF", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/protected", nil)
		req.Header.Set("X-CSRF-Token", "invalid-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
