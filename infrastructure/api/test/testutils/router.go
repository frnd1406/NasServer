package testutils

import (
	"github.com/gin-gonic/gin"
	authhandlers "github.com/nas-ai/api/src/handlers/auth"
	filehandlers "github.com/nas-ai/api/src/handlers/files"
	"github.com/nas-ai/api/src/middleware/logic"
	auth_repo "github.com/nas-ai/api/src/repository/auth"
	"github.com/nas-ai/api/src/services/content"
	"github.com/nas-ai/api/src/services/intelligence"
	"github.com/nas-ai/api/src/services/security"
)

// SetupTestRouter wraps the REAL handlers using injected interfaces.
// This allows testing the real handlers with either:
// 1. Real Services (via TestEnv) - For Integration Tests
// 2. Mock Services (via mocks) - For Unit Tests
//
// Arguments are now Interfaces (from src/services/security and src/repository/auth)
func SetupTestRouter(
	env *TestEnv,
	userRepo auth_repo.UserRepositoryInterface,
	jwtService security.JWTServiceInterface,
	passwordService security.PasswordServiceInterface,
	tokenService security.TokenServiceInterface,
	storageService content.StorageService,
	honeyService content.HoneyfileServiceInterface,
	aiService intelligence.AIAgentServiceInterface,
	policyService security.EncryptionPolicyServiceInterface,
	encryptionService security.EncryptionServiceInterface,
	encryptedStorageService content.EncryptedStorageServiceInterface,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Request ID middleware (required by handlers)
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-id")
		c.Next()
	})

	// Real Email Service (real config) - kept as struct for now
	realEmailService := MakeTestEmailService(env)

	// Auth Group
	authGroup := router.Group("/auth")
	{
		// Inject dependencies into REAL Handlers
		authGroup.POST("/register", authhandlers.RegisterHandler(
			env.Config,
			userRepo,
			jwtService,
			passwordService,
			tokenService,
			realEmailService, // struct
			env.RedisClient,  // struct
			env.Logger,
		))

		// Logout handler (wrapper)

		authGroup.POST("/login", authhandlers.LoginHandler(
			userRepo,
			jwtService,
			passwordService,
			env.RedisClient,
			env.Logger,
		))

		authGroup.POST("/refresh", authhandlers.RefreshHandler(
			jwtService,
			env.RedisClient,
			env.Logger,
		))

		authGroup.POST("/logout", authhandlers.LogoutHandler(
			jwtService,
			env.RedisClient,
			env.Logger,
		))
	}

	// Helper for valid protected route testing
	// We use the REAL AuthMiddleware here for testing it.
	// Since we are creating a "TestRouter" that mimics production but with specific injected services.
	// Note: AuthMiddleware also takes interfaces now.
	protectedGroup := router.Group("/api/v1")
	protectedGroup.Use(logic.AuthMiddleware(jwtService, tokenService, env.RedisClient, env.Logger))
	{
		// Sample protected route for testing middleware
		protectedGroup.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			userEmail, _ := c.Get("user_email")
			c.JSON(200, gin.H{
				"message":    "protected_resource",
				"user_id":    userID,
				"user_email": userEmail,
			})
		})

		// File Routes
		// Only register if services are provided (allows partial setup)
		if storageService != nil {
			storageGroup := protectedGroup.Group("/storage")
			{
				storageGroup.POST("/upload", filehandlers.StorageUploadHandler(storageService, policyService, honeyService, aiService, env.Logger))
				storageGroup.GET("/list", filehandlers.StorageListHandler(storageService, env.Logger))
				// Add others if needed
			}

			filesGroup := protectedGroup.Group("/files")
			{
				filesGroup.GET("/content", filehandlers.StorageDownloadHandler(storageService, honeyService, env.Logger))
			}
		}

		// Encrypted Storage Routes
		if encryptedStorageService != nil {
			encV1 := protectedGroup.Group("/encrypted")
			{
				encV1.GET("/status", filehandlers.EncryptedStorageStatusHandler(encryptedStorageService, encryptionService))
				encV1.GET("/files", filehandlers.EncryptedStorageListHandler(encryptedStorageService, env.Logger))
				encV1.POST("/upload", filehandlers.EncryptedStorageUploadHandler(encryptedStorageService, nil, env.Logger)) // nil aiFeeder
				encV1.GET("/download", filehandlers.EncryptedStorageDownloadHandler(encryptedStorageService, env.Logger))
				encV1.GET("/preview", filehandlers.EncryptedStoragePreviewHandler(encryptedStorageService, env.Logger))
				encV1.DELETE("/delete", filehandlers.EncryptedStorageDeleteHandler(encryptedStorageService, env.Logger))
			}
		}
	}

	return router
}

// NewTestRouter is the helper for Integration Tests (backward compatibility / convenience)
// It injects the REAL services from TestUtils
func NewTestRouter(env *TestEnv) *gin.Engine {
	// Real Repos & Services
	userRepo := MakeRealUserRepository(env)
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	passwordService := security.NewPasswordService()
	tokenService := security.NewTokenService(env.RedisClient, env.Logger)

	// Pass them to SetupTestRouter (with nil for file services to maintain back-compat)
	return SetupTestRouter(env, userRepo, jwtService, passwordService, tokenService, nil, nil, nil, nil, nil, nil)
}
