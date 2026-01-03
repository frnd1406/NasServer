package testutils

import (
	"github.com/gin-gonic/gin"
	authhandlers "github.com/nas-ai/api/src/handlers/auth"
	auth_repo "github.com/nas-ai/api/src/repository/auth"
	"github.com/nas-ai/api/src/services/operations"
	"github.com/nas-ai/api/src/services/security"
)

// ============================================================
// REAL HANDLER WRAPPERS
// These import and call the REAL production handlers
// ============================================================

// TestableRegisterHandler wraps the REAL RegisterHandler from src/handlers/auth
// It creates real service instances configured for testing
func TestableRegisterHandler(env *TestEnv) gin.HandlerFunc {
	// Create REAL repository (uses in-memory SQLite for tests)
	userRepo := auth_repo.NewUserRepository(env.DB, env.Logger)

	// Create REAL services
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	passwordService := security.NewPasswordService()
	tokenService := security.NewTokenService(env.RedisClient, env.Logger)

	// Create REAL email service (configured for test mode - won't actually send emails)
	emailService := operations.NewEmailService(env.Config, env.Logger)

	// Call the REAL handler constructor
	return authhandlers.RegisterHandler(
		env.Config,
		userRepo,
		jwtService,
		passwordService,
		tokenService,
		emailService,
		env.RedisClient,
		env.Logger,
	)
}

// TestableLoginHandler wraps the REAL LoginHandler from src/handlers/auth
func TestableLoginHandler(env *TestEnv) gin.HandlerFunc {
	// Create REAL repository
	userRepo := auth_repo.NewUserRepository(env.DB, env.Logger)

	// Create REAL services
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	passwordService := security.NewPasswordService()

	// Call the REAL handler constructor
	return authhandlers.LoginHandler(
		userRepo,
		jwtService,
		passwordService,
		env.RedisClient,
		env.Logger,
	)
}
