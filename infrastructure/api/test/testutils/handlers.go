package testutils

import (
	"net/http"
	"regexp"

	"crypto/subtle"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/sirupsen/logrus"
)

// ============================================================
// REAL HANDLERS with Interface-based DI
// These mirror the real handlers but accept interfaces for testability
// ============================================================

// RegisterRequest matches the real handler's request
type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	InviteCode string `json:"invite_code"`
}

// TestableRegisterHandler creates a register handler that accepts interfaces.
// This allows testing the REAL business logic with mocked dependencies.
func TestableRegisterHandler(
	cfg *config.Config,
	userRepo UserRepositoryInterface,
	jwtService JWTServiceInterface,
	passwordService PasswordServiceInterface,
	logger *logrus.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "invalid_request",
					"message":    "Invalid request body",
					"request_id": requestID,
				},
			})
			return
		}

		// Security Check: Invite Code
		if cfg.InviteCode != "" {
			if subtle.ConstantTimeCompare([]byte(req.InviteCode), []byte(cfg.InviteCode)) != 1 {
				c.JSON(http.StatusForbidden, gin.H{
					"error": gin.H{
						"code":       "invalid_invite_code",
						"message":    "Invalid or missing invite code",
						"request_id": requestID,
					},
				})
				return
			}
		}

		// Validate username
		if len(req.Username) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "invalid_username",
					"message":    "Username must be at least 3 characters",
					"request_id": requestID,
				},
			})
			return
		}

		// Validate email format
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "invalid_email",
					"message":    "Invalid email format",
					"request_id": requestID,
				},
			})
			return
		}

		ctx := c.Request.Context()

		// Check if username already exists
		existingByUsername, err := userRepo.FindByUsername(ctx, req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_error", "message": "Failed to create user"},
			})
			return
		}
		if existingByUsername != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":       "username_exists",
					"message":    "Username already registered",
					"request_id": requestID,
				},
			})
			return
		}

		// Validate password strength
		if err := passwordService.ValidatePasswordStrength(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "weak_password",
					"message":    err.Error(),
					"request_id": requestID,
				},
			})
			return
		}

		// Check if email already exists
		existingByEmail, err := userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_error", "message": "Failed to create user"},
			})
			return
		}
		if existingByEmail != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":       "email_exists",
					"message":    "Email already registered",
					"request_id": requestID,
				},
			})
			return
		}

		// Hash password
		passwordHash, err := passwordService.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_error", "message": "Failed to create user"},
			})
			return
		}

		// Create user
		user, err := userRepo.CreateUser(ctx, req.Username, req.Email, passwordHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_error", "message": "Failed to create user"},
			})
			return
		}

		// Generate JWT tokens
		accessToken, err := jwtService.GenerateAccessToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		refreshToken, err := jwtService.GenerateRefreshToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		logger.WithFields(logrus.Fields{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		}).Info("User registered successfully")

		c.JSON(http.StatusCreated, gin.H{
			"user":          user.ToResponse(),
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

// LoginRequest matches the real handler's request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TestableLoginHandler creates a login handler that accepts interfaces.
func TestableLoginHandler(
	userRepo UserRepositoryInterface,
	jwtService JWTServiceInterface,
	passwordService PasswordServiceInterface,
	redis *database.RedisClient,
	logger *logrus.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "invalid_request",
					"message":    "Invalid request body",
					"request_id": requestID,
				},
			})
			return
		}

		ctx := c.Request.Context()

		// Find user by email
		user, err := userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{"code": "internal_error", "message": "Login failed"},
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "invalid_credentials",
					"message":    "Invalid email or password",
					"request_id": requestID,
				},
			})
			return
		}

		// Verify password
		if err := passwordService.ComparePassword(user.PasswordHash, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "invalid_credentials",
					"message":    "Invalid email or password",
					"request_id": requestID,
				},
			})
			return
		}

		// Generate JWT tokens
		accessToken, err := jwtService.GenerateAccessToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		refreshToken, err := jwtService.GenerateRefreshToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Info("User logged in successfully")

		c.JSON(http.StatusOK, gin.H{
			"user":          user.ToResponse(),
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
