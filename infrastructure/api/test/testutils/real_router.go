package testutils

import (
	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/middleware/logic"
	"github.com/nas-ai/api/src/services/security"
)

// NewRealRouter creates a Gin router with REAL business logic handlers.
// The handlers use interfaces for dependencies, allowing mock injection.
// This tests the ACTUAL code paths, not inline fake handlers.
func NewRealRouter(env *TestEnv) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Request ID middleware (required by handlers)
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-id")
		c.Next()
	})

	// ============================================================
	// PUBLIC AUTH ROUTES (no auth middleware)
	// Uses REAL handlers with real dependencies
	// ============================================================
	authGroup := router.Group("/auth")
	{
		// Register uses REAL RegisterHandler
		authGroup.POST("/register", TestableRegisterHandler(env))

		// Login uses REAL LoginHandler
		authGroup.POST("/login", TestableLoginHandler(env))
	}

	// ============================================================
	// PROTECTED ROUTES (uses REAL AuthMiddleware)
	// ============================================================
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	tokenService := security.NewTokenService(env.RedisClient, env.Logger)

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
