package testutils

import (
	"github.com/gin-gonic/gin"
	authhandlers "github.com/nas-ai/api/src/handlers/auth"
	auth_repo "github.com/nas-ai/api/src/repository/auth"
	"github.com/nas-ai/api/src/services/operations"
	"github.com/nas-ai/api/src/services/security"
)

// ============================================================
// HELPERS
// ============================================================

// MakeRealUserRepository creates a real UserRepository connected to TestEnv DB
func MakeRealUserRepository(env *TestEnv) *auth_repo.UserRepository {
	return auth_repo.NewUserRepository(env.DB, env.Logger)
}

// MakeTestEmailService creates a real EmailService configured for testing (won't send)
func MakeTestEmailService(env *TestEnv) *operations.EmailService {
	return operations.NewEmailService(env.Config, env.Logger)
}

// ============================================================
// REAL HANDLER WRAPPERS (Deprecated usage - prefer NewTestRouter)
// These import and call the REAL production handlers
// ============================================================

func TestableRegisterHandler(env *TestEnv) gin.HandlerFunc {
	// Create REAL repository (uses in-memory SQLite for tests)
	userRepo := MakeRealUserRepository(env)

	// Create REAL services
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	passwordService := security.NewPasswordService()
	tokenService := security.NewTokenService(env.RedisClient, env.Logger)
	emailService := MakeTestEmailService(env)

	// Call the REAL handler constructor (now accepting interfaces)
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

func TestableLoginHandler(env *TestEnv) gin.HandlerFunc {
	userRepo := MakeRealUserRepository(env)
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	passwordService := security.NewPasswordService()

	return authhandlers.LoginHandler(
		userRepo,
		jwtService,
		passwordService,
		env.RedisClient,
		env.Logger,
	)
}

func TestableRefreshHandler(env *TestEnv) gin.HandlerFunc {
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	return authhandlers.RefreshHandler(
		jwtService,
		env.RedisClient,
		env.Logger,
	)
}

func TestableLogoutHandler(env *TestEnv) gin.HandlerFunc {
	jwtService, _ := security.NewJWTService(env.Config, env.Logger)
	return authhandlers.LogoutHandler(
		jwtService,
		env.RedisClient,
		env.Logger,
	)
}
