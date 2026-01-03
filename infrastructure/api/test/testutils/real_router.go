package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/middleware/logic"
	"github.com/nas-ai/api/src/services/security"
)

// NewRealRouter creates a Gin router with REAL AuthMiddleware.
// This is for integration tests that need to test the actual middleware code paths.
// It uses real JWT validation with test secrets.
func NewRealRouter(env *TestEnv) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Real services with test config
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	tokenService := security.NewTokenService(env.RedisClient, env.Logger)

	// Request ID middleware (required by AuthMiddleware)
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-id")
		c.Next()
	})

	// Protected routes with REAL AuthMiddleware
	protectedGroup := router.Group("/api/v1")
	protectedGroup.Use(logic.AuthMiddleware(jwtService, tokenService, env.RedisClient, env.Logger))
	{
		protectedGroup.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetString("user_id"),
				"email":   c.GetString("user_email"),
			})
		})
	}

	return router
}
