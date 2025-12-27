package main

import (
	"context"
	"fmt"
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

	// Startup Resilience: Fail fast if dependencies are not healthy
	if err := CheckDependencies(db, redis); err != nil {
		logger.WithError(err).Fatal("Startup dependency check failed")
	}

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
	embeddingsRepo := repository.NewFileEmbeddingsRepository(dbx, logger)
	honeyfileRepo := repository.NewHoneyfileRepository(dbx, logger)
	if err := honeyfileRepo.EnsureTable(context.Background()); err != nil {
		logger.WithError(err).Fatal("Failed to ensure honeyfiles table")
	}

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

	// Initialize encryption service with configurable vault path
	// SECURITY: Zero-Knowledge Encryption!
	// Default: /tmp/nas-vault-demo (NICHT persistent)
	//   → Container restart = Vault weg = User muss neu einrichten
	//   → Maximale Sicherheit: Kein physischer Zugriff auf verschlüsselte Keys
	// Production: /var/lib/nas/vault (persistent, nur wenn Volume gemountet)
	//   → Convenience, aber Sicherheitsrisiko bei physischem Zugriff
	vaultPath := "/tmp/nas-vault-demo" // Default: Non-persistent (Zero-Knowledge)
	if cfg.Environment == "production" {
		// In Production: Verwende nur persistent wenn explizit gewünscht
		// User muss bewusst Volume mounten in docker-compose
		if _, err := os.Stat("/var/lib/nas/vault"); err == nil {
			vaultPath = "/var/lib/nas/vault"
			logger.Warn("⚠️  Vault persistence enabled: Keys survive restarts (security trade-off)")
		}
	}
	encryptionService := services.NewEncryptionService(vaultPath)
	logger.WithField("vaultPath", vaultPath).Info("Encryption service initialized")

	// Initialize encrypted storage service for /media/frnd14/DEMO
	// Files stored here are encrypted - only visible via web UI when vault is unlocked
	encryptedStoragePath := "/media/frnd14/DEMO"
	encryptedStorageService, err := services.NewEncryptedStorageService(
		storageService,
		encryptionService,
		encryptedStoragePath,
		logger,
	)
	if err != nil {
		logger.WithError(err).Warn("Failed to initialize encrypted storage service (non-fatal)")
		// Don't fatal - encrypted storage is optional
	} else {
		logger.WithField("path", encryptedStoragePath).Info("Encrypted storage service initialized")
	}

	// Initialize SecureAIFeeder for encrypted file indexing
	// This pushes decrypted content to the AI agent without writing plaintext to disk
	var secureAIFeeder *services.SecureAIFeeder
	if encryptionService != nil {
		secureAIFeeder = services.NewSecureAIFeeder(
			encryptionService,
			cfg.AIServiceURL,
			cfg.InternalAPISecret,
			logger,
		)
		logger.Info("SecureAIFeeder initialized for encrypted content indexing")
	}

	aiHTTPClient := services.NewSecureHTTPClient(cfg.InternalAPISecret, 15*time.Second)

	// Initialize JobService for async AI queries
	jobService := services.NewJobService(redis, logger)
	if err := jobService.EnsureConsumerGroup(context.Background()); err != nil {
		logger.WithError(err).Warn("Failed to ensure AI job consumer group (non-fatal)")
	}
	logger.Info("JobService initialized for async AI queries")

	// Initialize HoneyfileService for intrusion detection
	honeyfileService := services.NewHoneyfileService(honeyfileRepo, encryptionService, logger)
	logger.Info("HoneyfileService initialized for intrusion detection")

	// Phase 8: Dead Man's Switch (Alerting)
	alertService := services.NewAlertService(emailService, cfg, logger)
	// Start background alert ticker (every 5 minutes)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			alertService.RunSystemChecks()
		}
	}()
	logger.Info("AlertService initialized (Dead Man's Switch active)")

	// Initialize BlobStorageHandler for chunked uploads (Zero-Knowledge Phase 4)
	blobStorageHandler := handlers.NewBlobStorageHandler(storageService, logger)
	logger.Info("BlobStorageHandler initialized for chunked encrypted uploads")

	// ==== PHASE 2: Performance Guard - Startup Benchmark ====
	// Run crypto benchmark to measure system encryption speed
	// This enables intelligent warnings for large file encryption
	benchmarkService := services.NewBenchmarkService(logger)
	go func() {
		if err := benchmarkService.RunStartupBenchmark(); err != nil {
			logger.WithError(err).Warn("Startup benchmark failed (non-fatal)")
		}
	}()

	go func() {
		if err := scheduler.StartBackupScheduler(backupService, cfg); err != nil {
			logger.WithError(err).Error("Failed to start backup scheduler")
		}
	}()

	// Initialize Consistency Service (orphan cleanup)
	consistencyService := services.NewConsistencyService(
		dbx,
		embeddingsRepo,
		"/mnt/data",
		time.Duration(cfg.ConsistencyCheckIntervalMin)*time.Minute,
		logger,
	)

	// Run initial reconciliation BEFORE HTTP server starts (ensures clean state)
	logger.Info("Running initial consistency reconciliation...")
	if err := consistencyService.RunReconciliation(context.Background()); err != nil {
		logger.WithError(err).Warn("Initial reconciliation failed (non-fatal)")
	}

	// Start background consistency worker
	go consistencyService.Start(context.Background())

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
			handlers.RegisterHandler(cfg, userRepo, jwtService, passwordService, tokenService, emailService, redis, logger),
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
		v1.GET("/system/metrics/live", handlers.SystemMetricsLiveHandler(logger)) // New Phase 8 Endpoint
		v1.GET("/system/alerts", handlers.SystemAlertsListHandler(systemAlertsRepo, logger))
		v1.POST("/system/alerts", handlers.SystemAlertCreateHandler(systemAlertsRepo, logger))
		v1.POST("/system/alerts/:id/resolve", handlers.SystemAlertResolveHandler(systemAlertsRepo, logger))
		v1.GET("/search", handlers.SearchHandler(db, cfg.AIServiceURL, aiHTTPClient, logger))
		v1.POST("/query", handlers.UnifiedQueryHandler(cfg.AIServiceURL, services.NewSecureHTTPClient(cfg.InternalAPISecret, 90*time.Second), jobService, logger))
		v1.GET("/jobs/:id", handlers.GetJobStatusHandler(jobService, logger))
		v1.GET("/ask", handlers.AskHandler(db, cfg.AIServiceURL, cfg.OllamaURL, cfg.LLMModel, nil, logger))
		v1.GET("/files/content", handlers.FileContentHandler(logger))

		// ==== PHASE 2: System Capabilities (Performance Guard) ====
		v1.GET("/system/capabilities", handlers.Capabilities(benchmarkService))

		// AI Settings endpoints
		v1.GET("/ai/status", handlers.AIStatusHandler(cfg.AIServiceURL, aiHTTPClient, logger))
		v1.GET("/ai/settings", handlers.AISettingsGetHandler(logger))
		v1.POST("/ai/settings", handlers.AISettingsSaveHandler(logger))
		v1.POST("/ai/reindex", handlers.AIReindexHandler(cfg.AIServiceURL, aiHTTPClient, logger))
		v1.POST("/ai/warmup", handlers.AIWarmupHandler(logger))

		// Network Settings endpoints
		v1.GET("/network/settings", handlers.NetworkSettingsGetHandler(logger))
		v1.PUT("/network/settings", handlers.NetworkSettingsSaveHandler(logger))

		// Backup Settings endpoints
		v1.GET("/backup/settings", handlers.BackupSettingsGetHandler(logger))
		v1.PUT("/backup/settings", handlers.BackupSettingsSaveHandler(logger))

		// Security Settings endpoints
		v1.GET("/security/settings", handlers.SecuritySettingsGetHandler(logger))
		v1.PUT("/security/settings", handlers.SecuritySettingsSaveHandler(logger))

		// Storage Monitor Settings endpoints
		v1.GET("/storage/settings", handlers.StorageSettingsGetHandler(logger))
		v1.PUT("/storage/settings", handlers.StorageSettingsSaveHandler(logger))

		// Setup endpoints (first-time wizard)
		v1.GET("/system/setup-status", handlers.SetupStatusHandler(logger))

		// Vault endpoints (encryption management)
		v1.GET("/vault/status", handlers.VaultStatusHandler(encryptionService))
		v1.GET("/vault/config", handlers.VaultConfigGetHandler(encryptionService))
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

		// Protected vault endpoints (require auth)
		settingsV1.POST("/vault/setup", handlers.VaultSetupHandler(encryptionService, logger))
		settingsV1.POST("/vault/unlock", handlers.VaultUnlockHandler(encryptionService, logger))
		settingsV1.POST("/vault/lock", handlers.VaultLockHandler(encryptionService, logger))
		settingsV1.POST("/vault/panic", handlers.VaultPanicHandler(encryptionService, logger)) // EMERGENCY: Destroys all keys
		settingsV1.PUT("/vault/config", handlers.VaultConfigUpdateHandler(encryptionService, logger))
		settingsV1.GET("/vault/export-config", handlers.VaultExportConfigHandler(encryptionService, logger)) // Backup download

		// Setup wizard (first-time configuration)
		settingsV1.POST("/setup", handlers.SetupHandler(logger))
	}

	// Create a dedicated group for Vault Uploads (Chunked)
	vaultUploadV1 := v1.Group("/vault/upload")
	vaultUploadV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
	)
	{
		vaultUploadV1.POST("/init", blobStorageHandler.InitUpload)
		vaultUploadV1.POST("/chunk/:id", blobStorageHandler.UploadChunk)
		vaultUploadV1.POST("/finalize/:id", blobStorageHandler.FinalizeUpload)
	}

	storageV1 := r.Group("/api/v1/storage")
	storageV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
	)
	{
		storageV1.GET("/files", handlers.StorageListHandler(storageService, logger))
		storageV1.POST("/upload", handlers.StorageUploadHandler(storageService, honeyfileService, cfg, logger))
		storageV1.GET("/download", handlers.StorageDownloadHandler(storageService, honeyfileService, logger))

		// ==== PHASE 4: Smart Download (Hybrid Streaming) ====
		// X-Accel-Redirect for unencrypted, streaming decrypt for encrypted
		storageV1.GET("/smart-download", handlers.SmartDownloadHandler(storageService, honeyfileService, logger))
		storageV1.GET("/download-zip", handlers.StorageDownloadZipHandler(storageService, logger))
		storageV1.POST("/batch-download", handlers.StorageBatchDownloadHandler(storageService, logger))
		storageV1.DELETE("/delete", handlers.StorageDeleteHandler(storageService, cfg, logger))
		storageV1.GET("/trash", handlers.StorageTrashListHandler(storageService, logger))
		storageV1.POST("/trash/restore/:id", handlers.StorageTrashRestoreHandler(storageService, logger))
		storageV1.DELETE("/trash/:id", handlers.StorageTrashDeleteHandler(storageService, logger))
		storageV1.POST("/trash/empty", handlers.StorageTrashEmptyHandler(storageService, logger))
		storageV1.POST("/rename", handlers.StorageRenameHandler(storageService, logger))
		storageV1.POST("/move", handlers.StorageMoveHandler(storageService, logger))
		storageV1.POST("/mkdir", handlers.StorageMkdirHandler(storageService, logger))
		storageV1.POST("/upload-zip", handlers.StorageUploadZipHandler(storageService, logger))
	}

	// Encrypted Storage API (separate from regular storage)
	// Files in /media/frnd14/DEMO are encrypted and only visible via web UI
	if encryptedStorageService != nil {
		encV1 := r.Group("/api/v1/encrypted")
		encV1.Use(
			middleware.AuthMiddleware(jwtService, redis, logger),
			middleware.CSRFMiddleware(redis, logger),
		)
		{
			encV1.GET("/status", handlers.EncryptedStorageStatusHandler(encryptedStorageService))
			encV1.GET("/files", handlers.EncryptedStorageListHandler(encryptedStorageService, logger))
			encV1.POST("/upload", handlers.EncryptedStorageUploadHandler(encryptedStorageService, secureAIFeeder, logger))
			encV1.GET("/download", handlers.EncryptedStorageDownloadHandler(encryptedStorageService, logger))
			encV1.GET("/preview", handlers.EncryptedStoragePreviewHandler(encryptedStorageService, logger))
			encV1.DELETE("/delete", handlers.EncryptedStorageDeleteHandler(encryptedStorageService, logger))
		}
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

	// === ADMIN ROUTES (requires JWT + CSRF + Admin role) ===
	adminV1 := r.Group("/api/v1/admin")
	adminV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
		middleware.AdminOnly(userRepo, logger),
	)
	{
		// Admin settings
		adminV1.GET("/settings", handlers.GetAdminSettingsHandler(cfg, logger))
		adminV1.PUT("/settings", handlers.UpdateAdminSettingsHandler(cfg, settingsRepo, logger))

		// System status
		adminV1.GET("/status", handlers.SystemStatusHandler(db, logger))

		// User management
		adminV1.GET("/users", handlers.UserListHandler(userRepo, logger))
		adminV1.PUT("/users/:id/role", handlers.UpdateUserRoleHandler(userRepo, logger))

		// Maintenance mode
		adminV1.POST("/maintenance", handlers.ToggleMaintenanceModeHandler(logger))

		// Audit logs
		adminV1.GET("/audit-logs", handlers.AuditLogHandler(db, logger))

		// Knowledge Index Reconciliation (garbage collection for ghost knowledge)
		adminV1.POST("/system/reconcile-knowledge", handlers.ReconcileKnowledgeHandler(secureAIFeeder, "/mnt/data", logger))
	}

	// === SYSTEM INTEGRITY MONITORING (Stealth - No Swagger exposure) ===
	sysV1 := r.Group("/api/v1/sys")
	sysV1.Use(
		middleware.AuthMiddleware(jwtService, redis, logger),
		middleware.CSRFMiddleware(redis, logger),
		middleware.AdminOnly(userRepo, logger),
	)
	{
		// Stealth endpoint - not documented in Swagger
		sysV1.POST("/integrity/checkpoints", handlers.CreateCheckpointHandler(honeyfileService, logger))
	}

	// Create HTTP server
	secureHandler := middleware.SecureHeaders(r)

	srv := &http.Server{
		Addr:           "0.0.0.0:" + cfg.Port,
		Handler:        secureHandler,
		ReadTimeout:    600 * time.Second,
		WriteTimeout:   600 * time.Second,
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

	// Stop background workers
	consistencyService.Stop()

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

// CheckDependencies verifies that critical infrastructure is reachable
func CheckDependencies(db *database.DB, redis *database.RedisClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check Postgres
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres unreachable: %w", err)
	}

	// Check Redis
	if err := redis.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis unreachable: %w", err)
	}

	return nil
}
