package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/handlers"
	"github.com/nas-ai/api/src/middleware"
	"github.com/nas-ai/api/src/repository"
	"github.com/nas-ai/api/src/scheduler"
	"github.com/nas-ai/api/src/services"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	_ "github.com/nas-ai/api/docs" // swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title NAS.AI API
// @version 1.0
// @description Secure file storage and management API with authentication, email verification, and password reset.
// @termsOfService https://your-domain.com/terms

// @contact.name API Support
// @contact.url https://your-domain.com/support
// @contact.email support@your-domain.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host api.your-domain.com
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.apikey CSRFToken
// @in header
// @name X-CSRF-Token
// @description CSRF token for state-changing operations

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Load configuration (FAIL-FAST if secrets missing!)
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	// Set log level
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Log startup
	logger.WithFields(logrus.Fields{
		"port":         cfg.Port,
		"environment":  cfg.Environment,
		"log_level":    cfg.LogLevel,
		"cors_origins": cfg.CORSOrigins,
		"rate_limit":   cfg.RateLimitPerMin,
	}).Info("Starting NAS.AI API server")

	// Initialize database connections (FAIL-FAST if can't connect!)
	db, err := database.NewPostgresConnection(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to PostgreSQL")
	}
	defer db.Close()

	redis, err := database.NewRedisConnection(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to Redis")
	}
	defer redis.Close()

	// Initialize repositories
	dbx := sqlx.NewDb(db.DB, "postgres")
	settingsRepo := repository.NewSystemSettingsRepository(dbx, logger)
	if err := settingsRepo.EnsureTable(context.Background()); err != nil {
		logger.WithError(err).Fatal("Failed to ensure system settings table")
	}
	applyPersistedBackupSettings(context.Background(), settingsRepo, cfg, logger)

	userRepo := repository.NewUserRepository(db, logger)
	systemMetricsRepo := repository.NewSystemMetricsRepository(dbx, logger)
	systemAlertsRepo := repository.NewSystemAlertsRepository(dbx, logger)
	monitoringRepo := repository.NewMonitoringRepository(db, logger)

	// Initialize services
	jwtService, err := services.NewJWTService(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize JWT service")
	}
	passwordService := services.NewPasswordService()
	tokenService := services.NewTokenService(redis, logger)
	emailService := services.NewEmailService(cfg, logger)
	storageService, err := services.NewStorageService("/mnt/data", logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize storage service")
	}
	backupService, err := services.NewBackupService("/mnt/data", cfg.BackupStoragePath, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize backup service")
	}
	aiHTTPClient := &http.Client{Timeout: 8 * time.Second}

	// Start backup scheduler in background
	go func() {
		if err := scheduler.StartBackupScheduler(backupService, cfg); err != nil {
			logger.WithError(err).Error("Failed to start backup scheduler")
		}
	}()

	// Create Gin engine (without default middleware)
	r := gin.New()

	// Build middleware chain (ZWIEBEL-PRINZIP / ONION PRINCIPLE)
	// Order matters! Outer layers execute first.
	//
	// Request Flow:
	//   1. Panic Recovery (catch crashes)
	//   2. Request ID (generate UUID)
	//   3. Security Headers (set security headers)
	//   4. CORS (check origin whitelist)
	//   5. Rate Limit (check request limits)
	//   6. Audit Logger (log request/response)
	//   7. Handler (business logic)
	//
	// Response flows back through the same layers

	rateLimiter := middleware.NewRateLimiter(cfg)

	r.Use(
		middleware.PanicRecovery(logger), // 1. Catch panics
		middleware.RequestID(),           // 2. Generate request ID
		middleware.GinSecureHeaders(),    // 3. Security headers
		middleware.CORS(cfg, logger),     // 4. CORS whitelist
		rateLimiter.Middleware(),         // 5. Rate limiting
		middleware.AuditLogger(logger),   // 6. Audit logging
	)

	// Allow OPTIONS preflight globally without auth/CSRF
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Store environment in context for middleware
	r.Use(func(c *gin.Context) {
		c.Set("environment", cfg.Environment)
		c.Next()
	})

	// === PUBLIC ROUTES (no auth, but rate-limited) ===
	r.GET("/health", handlers.Health(db, redis, logger))
	r.POST("/monitoring/ingest", handlers.MonitoringIngestHandler(monitoringRepo, cfg.MonitoringToken, logger))

	// Swagger documentation (only in development)
	if cfg.Environment != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// === AUTH ROUTES (public, but rate-limited) ===
	authGroup := r.Group("/auth")
	{
		authLimiter := middleware.NewRateLimiter(&config.Config{RateLimitPerMin: 5})
		authGroup.POST("/register",
			authLimiter.Middleware(),
			handlers.RegisterHandler(userRepo, jwtService, passwordService, tokenService, emailService, redis, logger),
		)
		authGroup.POST("/login",
			authLimiter.Middleware(),
			handlers.LoginHandler(userRepo, jwtService, passwordService, redis, logger),
		)
		authGroup.POST("/refresh", handlers.RefreshHandler(jwtService, redis, logger))
		authGroup.POST("/logout",
			middleware.AuthMiddleware(jwtService, redis, logger), // Require auth for logout
			handlers.LogoutHandler(jwtService, redis, logger),
		)

		// Email verification endpoints
		authGroup.POST("/verify-email", handlers.VerifyEmailHandler(userRepo, tokenService, emailService, logger))
		authGroup.POST("/resend-verification",
			middleware.AuthMiddleware(jwtService, redis, logger), // Require auth
			handlers.ResendVerificationHandler(userRepo, tokenService, emailService, logger),
		)

		// Password reset endpoints
		authGroup.POST("/forgot-password", handlers.ForgotPasswordHandler(userRepo, tokenService, emailService, logger))
		authGroup.POST("/reset-password", handlers.ResetPasswordHandler(userRepo, tokenService, passwordService, jwtService, redis, logger))
	}

	// === PROTECTED API ROUTES (requires JWT + CSRF) ===
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(jwtService, redis, logger)) // JWT validation
	apiGroup.Use(middleware.CSRFMiddleware(redis, logger))             // CSRF validation
	{
		apiGroup.GET("/profile", handlers.ProfileHandler(userRepo, logger))
		apiGroup.GET("/monitoring", handlers.MonitoringListHandler(monitoringRepo, logger))
	}

	// === SYSTEM METRICS (API-Key geschützt, ohne JWT) ===
	v1 := r.Group("/api/v1")
	{
		v1.GET("/auth/csrf", handlers.GetCSRFToken(redis, logger))
		v1.POST("/system/metrics", handlers.SystemMetricsHandler(systemMetricsRepo, cfg.MonitoringToken, logger))
		v1.GET("/system/metrics", handlers.SystemMetricsListHandler(systemMetricsRepo, logger))
		v1.GET("/system/alerts", handlers.SystemAlertsListHandler(systemAlertsRepo, logger))
		v1.POST("/system/alerts", handlers.SystemAlertCreateHandler(systemAlertsRepo, logger))
		v1.POST("/system/alerts/:id/resolve", handlers.SystemAlertResolveHandler(systemAlertsRepo, logger))
		v1.GET("/search", handlers.SearchHandler(db, cfg.AIServiceURL, aiHTTPClient, logger))
		v1.GET("/ask", handlers.AskHandler(db, cfg.AIServiceURL, cfg.OllamaURL, cfg.LLMModel, nil, logger))
	}

	settingsV1 := v1.Group("/system")
	settingsV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
	)
	{
		settingsV1.GET("/settings", handlers.SystemSettingsHandler(cfg))
		settingsV1.PUT("/settings/backup", handlers.UpdateBackupSettingsHandler(cfg, backupService, settingsRepo, logger))
		settingsV1.POST("/validate-path", handlers.ValidatePathHandler(logger))
	}

	storageV1 := r.Group("/api/v1/storage")
	storageV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
	)
	{
		storageV1.GET("/files", handlers.StorageListHandler(storageService, logger))
		storageV1.POST("/upload", handlers.StorageUploadHandler(storageService, logger))
		storageV1.GET("/download", handlers.StorageDownloadHandler(storageService, logger))
		storageV1.DELETE("/delete", handlers.StorageDeleteHandler(storageService, logger))
		storageV1.GET("/trash", handlers.StorageTrashListHandler(storageService, logger))
		storageV1.POST("/trash/restore/:id", handlers.StorageTrashRestoreHandler(storageService, logger))
		storageV1.DELETE("/trash/:id", handlers.StorageTrashDeleteHandler(storageService, logger))
		storageV1.POST("/rename", handlers.StorageRenameHandler(storageService, logger))
	}

	backupV1 := r.Group("/api/v1/backups")
	backupV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
	)
	{
		// List and create backups - all authenticated users
		backupV1.GET("", handlers.BackupListHandler(backupService, logger))
		backupV1.POST("", handlers.BackupCreateHandler(backupService, cfg, logger))

		// SECURITY: Restore and delete require admin role (destructive operations)
		backupV1.POST("/:id/restore",
			middleware.AdminOnly(userRepo, logger),
			handlers.BackupRestoreHandler(backupService, cfg, logger),
		)
		backupV1.DELETE("/:id",
			middleware.AdminOnly(userRepo, logger),
			handlers.BackupDeleteHandler(backupService, logger),
		)
	}

	// Create HTTP server
	secureHandler := middleware.SecureHeaders(r)

	srv := &http.Server{
		Addr:           "0.0.0.0:" + cfg.Port,
		Handler:        secureHandler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		logger.WithField("port", cfg.Port).Info("Server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	}

	logger.Info("Server exited")
}

func applyPersistedBackupSettings(ctx context.Context, settingsRepo *repository.SystemSettingsRepository, cfg *config.Config, logger *logrus.Logger) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	settings, err := settingsRepo.GetAll(ctx)
	if err != nil {
		logger.WithError(err).Warn("Failed to load persisted settings; continuing with defaults")
		return
	}

	if schedule, ok := settings[repository.SystemSettingBackupSchedule]; ok {
		s := strings.TrimSpace(schedule)
		if s != "" {
			if _, err := parser.Parse(s); err != nil {
				logger.WithError(err).Warn("Ignoring invalid persisted backup schedule")
			} else {
				cfg.BackupSchedule = s
			}
		}
	}

	if retentionStr, ok := settings[repository.SystemSettingBackupRetention]; ok {
		if n, err := strconv.Atoi(retentionStr); err == nil && n > 0 {
			cfg.BackupRetentionCount = n
		} else if err != nil {
			logger.WithError(err).Warn("Ignoring invalid persisted backup retention")
		}
	}

	if path, ok := settings[repository.SystemSettingBackupPath]; ok {
		p := filepath.Clean(strings.TrimSpace(path))
		if p != "" && p != "." && p != string(os.PathSeparator) {
			cfg.BackupStoragePath = p
		} else {
			logger.Warn("Ignoring invalid persisted backup path")
		}
	}
}
