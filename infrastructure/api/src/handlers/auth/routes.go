package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/middleware/logic"
	auth_repo "github.com/nas-ai/api/src/repository/auth"
	"github.com/nas-ai/api/src/services/operations"
	"github.com/nas-ai/api/src/services/security"
	"github.com/sirupsen/logrus"
)

// Handler holds dependencies for auth handlers
type Handler struct {
	cfg             *config.Config
	userRepo        *auth_repo.UserRepository
	jwtService      *security.JWTService
	passwordService *security.PasswordService
	tokenService    *security.TokenService
	emailService    *operations.EmailService
	redis           *database.RedisClient
	logger          *logrus.Logger
}

// NewHandler creates a new Auth Handler
func NewHandler(
	cfg *config.Config,
	userRepo *auth_repo.UserRepository,
	jwtService *security.JWTService,
	passwordService *security.PasswordService,
	tokenService *security.TokenService,
	emailService *operations.EmailService,
	redis *database.RedisClient,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		cfg:             cfg,
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		tokenService:    tokenService,
		emailService:    emailService,
		redis:           redis,
		logger:          logger,
	}
}

// RegisterGlobalRoutes registers public auth routes (usually under /auth)
func (h *Handler) RegisterGlobalRoutes(rg *gin.RouterGroup) {
	authLimiter := logic.NewRateLimiter(&config.Config{RateLimitPerMin: 5})

	rg.POST("/register",
		authLimiter.Middleware(),
		RegisterHandler(h.cfg, h.userRepo, h.jwtService, h.passwordService, h.tokenService, h.emailService, h.redis, h.logger),
	)
	rg.POST("/login",
		authLimiter.Middleware(),
		LoginHandler(h.userRepo, h.jwtService, h.passwordService, h.redis, h.logger),
	)
	rg.POST("/refresh", RefreshHandler(h.jwtService, h.redis, h.logger))
	rg.POST("/logout",
		logic.AuthMiddleware(h.jwtService, h.tokenService, h.redis, h.logger),
		LogoutHandler(h.jwtService, h.redis, h.logger),
	)

	// Email verification
	rg.POST("/verify-email", VerifyEmailHandler(h.userRepo, h.tokenService, h.emailService, h.logger))
	rg.POST("/resend-verification",
		logic.AuthMiddleware(h.jwtService, h.tokenService, h.redis, h.logger),
		ResendVerificationHandler(h.userRepo, h.tokenService, h.emailService, h.logger),
	)

	// Password reset
	rg.POST("/forgot-password", ForgotPasswordHandler(h.userRepo, h.tokenService, h.emailService, h.logger))
	rg.POST("/reset-password", ResetPasswordHandler(h.userRepo, h.tokenService, h.passwordService, h.jwtService, h.redis, h.logger))
}

// RegisterV1Routes registers API v1 auth routes (usually under /api/v1/auth)
func (h *Handler) RegisterV1Routes(rg *gin.RouterGroup) {
	rg.GET("/csrf", GetCSRFToken(h.redis, h.logger))
}
