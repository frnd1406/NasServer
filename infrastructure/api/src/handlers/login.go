package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/middleware"
	"github.com/nas-ai/api/src/repository"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	CSRFToken    string      `json:"csrf_token"`
}

// LoginHandler godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful with tokens"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func LoginHandler(
	userRepo *repository.UserRepository,
	jwtService *services.JWTService,
	passwordService *services.PasswordService,
	redis *database.RedisClient,
	logger *logrus.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("Invalid login request")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "invalid_request",
					"message":    "Invalid request body",
					"request_id": requestID,
				},
			})
			return
		}

		// Find user by email
		ctx := c.Request.Context()
		user, err := userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			logger.WithError(err).Error("Failed to find user")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":       "internal_error",
					"message":    "Login failed",
					"request_id": requestID,
				},
			})
			return
		}

		// User not found or invalid password - same error message (security)
		if user == nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"email":      req.Email,
				"ip":         c.ClientIP(),
			}).Warn("Login failed: user not found")
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
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"user_id":    user.ID,
				"email":      req.Email,
				"ip":         c.ClientIP(),
			}).Warn("Login failed: invalid password")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "invalid_credentials",
					"message":    "Invalid email or password",
					"request_id": requestID,
				},
			})
			return
		}

		// Generate new session-scoped CSRF token (rotate session to avoid fixation)
		sessionID, err := middleware.EnsureCSRFSession(c)
		if err != nil {
			logger.WithError(err).Error("Failed to create CSRF session")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		// Generate JWT tokens
		accessToken, err := jwtService.GenerateAccessToken(user.ID, user.Email)
		if err != nil {
			logger.WithError(err).Error("Failed to generate access token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		refreshToken, err := jwtService.GenerateRefreshToken(user.ID, user.Email)
		if err != nil {
			logger.WithError(err).Error("Failed to generate refresh token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		// Generate CSRF token
		csrfToken, err := middleware.GenerateCSRFToken(redis, sessionID)
		if err != nil {
			logger.WithError(err).Error("Failed to generate CSRF token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}
		middleware.SetCSRFCookie(c, sessionID)

		// Audit log - SUCCESS
		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    user.ID,
			"email":      user.Email,
			"ip":         c.ClientIP(),
		}).Info("User logged in successfully")

		c.JSON(http.StatusOK, LoginResponse{
			User:         user.ToResponse(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			CSRFToken:    csrfToken,
		})
	}
}
