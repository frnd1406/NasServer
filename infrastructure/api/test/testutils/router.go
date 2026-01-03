package testutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewTestRouter creates a Gin router wired with mocks from TestEnv.
// This function handles DI and route registration, keeping tests focused on behavior.
func NewTestRouter(env *TestEnv) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Auth Routes
	authGroup := router.Group("/auth")
	{
		// Login handler (mocked)
		authGroup.POST("/login", func(c *gin.Context) {
			var req struct {
				Email    string `json:"email" binding:"required,email"`
				Password string `json:"password" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": err.Error()}})
				return
			}

			// Call mocked UserRepo.FindByEmail
			user, err := env.UserRepo.FindByEmail(c.Request.Context(), req.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error"}})
				return
			}
			if user == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "invalid_credentials", "message": "Invalid email or password"}})
				return
			}

			// Call mocked PasswordService.ComparePassword
			if err := env.PasswordSvc.ComparePassword(user.PasswordHash, req.Password); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "invalid_credentials", "message": "Invalid email or password"}})
				return
			}

			// Call mocked JWTService.GenerateAccessToken
			accessToken, err := env.JWTService.GenerateAccessToken(user.ID, user.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
				return
			}

			refreshToken, err := env.JWTService.GenerateRefreshToken(user.ID, user.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user":          user.ToResponse(),
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"csrf_token":    "mock-csrf-token",
			})
		})

		// Register handler (mocked)
		authGroup.POST("/register", func(c *gin.Context) {
			var req struct {
				Username   string `json:"username" binding:"required"`
				Email      string `json:"email" binding:"required,email"`
				Password   string `json:"password" binding:"required"`
				InviteCode string `json:"invite_code"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_request", "message": err.Error()}})
				return
			}

			// Invite code check
			if env.Config.InviteCode != "" && req.InviteCode != env.Config.InviteCode {
				c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "invalid_invite_code"}})
				return
			}

			// Check username
			if len(req.Username) < 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "invalid_username", "message": "Username must be at least 3 characters"}})
				return
			}

			existingByUsername, _ := env.UserRepo.FindByUsername(c.Request.Context(), req.Username)
			if existingByUsername != nil {
				c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "username_exists"}})
				return
			}

			// Validate password strength
			if err := env.PasswordSvc.ValidatePasswordStrength(req.Password); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "weak_password", "message": err.Error()}})
				return
			}

			// Check email
			existingByEmail, _ := env.UserRepo.FindByEmail(c.Request.Context(), req.Email)
			if existingByEmail != nil {
				c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "email_exists"}})
				return
			}

			// Hash password
			passwordHash, err := env.PasswordSvc.HashPassword(req.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
				return
			}

			// Create user
			user, err := env.UserRepo.CreateUser(c.Request.Context(), req.Username, req.Email, passwordHash)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
				return
			}

			// Generate tokens
			accessToken, _ := env.JWTService.GenerateAccessToken(user.ID, user.Email)
			refreshToken, _ := env.JWTService.GenerateRefreshToken(user.ID, user.Email)

			c.JSON(http.StatusCreated, gin.H{
				"user":          user.ToResponse(),
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"csrf_token":    "mock-csrf-token",
			})
		})
	}

	return router
}
