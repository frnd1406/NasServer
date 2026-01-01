package server

import (
	"time"

	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/handlers"
	"github.com/nas-ai/api/src/middleware"
	"github.com/nas-ai/api/src/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all HTTP routes (SRP: Routing Only)
func (s *Server) SetupRoutes() {
	// === PUBLIC ROUTES (no auth, but rate-limited) ===
	s.router.GET("/health", handlers.Health(s.db, s.redis, s.logger))
	s.router.POST("/monitoring/ingest", handlers.MonitoringIngestHandler(s.monitoringRepo, s.cfg.MonitoringToken, s.logger))

	// Swagger documentation (only in development)
	if s.cfg.Environment != "production" {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.setupAuthRoutes()
	s.setupAPIRoutes()
	s.setupV1Routes()
	s.setupStorageRoutes()
	s.setupEncryptedStorageRoutes()
	s.setupBackupRoutes()
	s.setupAdminRoutes()
	s.setupSystemRoutes()
}

// setupAuthRoutes configures authentication endpoints
func (s *Server) setupAuthRoutes() {
	authGroup := s.router.Group("/auth")
	{
		authLimiter := middleware.NewRateLimiter(&config.Config{RateLimitPerMin: 5})

		authGroup.POST("/register",
			authLimiter.Middleware(),
			handlers.RegisterHandler(s.cfg, s.userRepo, s.jwtService, s.passwordService, s.tokenService, s.emailService, s.redis, s.logger),
		)
		authGroup.POST("/login",
			authLimiter.Middleware(),
			handlers.LoginHandler(s.userRepo, s.jwtService, s.passwordService, s.redis, s.logger),
		)
		authGroup.POST("/refresh", handlers.RefreshHandler(s.jwtService, s.redis, s.logger))
		authGroup.POST("/logout",
			middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
			handlers.LogoutHandler(s.jwtService, s.redis, s.logger),
		)

		// Email verification
		authGroup.POST("/verify-email", handlers.VerifyEmailHandler(s.userRepo, s.tokenService, s.emailService, s.logger))
		authGroup.POST("/resend-verification",
			middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
			handlers.ResendVerificationHandler(s.userRepo, s.tokenService, s.emailService, s.logger),
		)

		// Password reset
		authGroup.POST("/forgot-password", handlers.ForgotPasswordHandler(s.userRepo, s.tokenService, s.emailService, s.logger))
		authGroup.POST("/reset-password", handlers.ResetPasswordHandler(s.userRepo, s.tokenService, s.passwordService, s.jwtService, s.redis, s.logger))
	}
}

// setupAPIRoutes configures protected API endpoints
func (s *Server) setupAPIRoutes() {
	apiGroup := s.router.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(s.jwtService, s.redis, s.logger))
	apiGroup.Use(middleware.CSRFMiddleware(s.redis, s.logger))
	{
		apiGroup.GET("/profile", handlers.ProfileHandler(s.userRepo, s.logger))
		apiGroup.GET("/monitoring", handlers.MonitoringListHandler(s.monitoringRepo, s.logger))
	}
}

// setupV1Routes configures API v1 endpoints
func (s *Server) setupV1Routes() {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/auth/csrf", handlers.GetCSRFToken(s.redis, s.logger))
		v1.POST("/system/metrics", handlers.SystemMetricsHandler(s.systemMetricsRepo, s.cfg.MonitoringToken, s.logger))
		v1.GET("/system/metrics", handlers.SystemMetricsListHandler(s.systemMetricsRepo, s.logger))
		v1.GET("/system/metrics/live", handlers.SystemMetricsLiveHandler(s.logger))
		v1.GET("/system/alerts", handlers.SystemAlertsListHandler(s.systemAlertsRepo, s.logger))
		v1.POST("/system/alerts", handlers.SystemAlertCreateHandler(s.systemAlertsRepo, s.logger))
		v1.POST("/system/alerts/:id/resolve", handlers.SystemAlertResolveHandler(s.systemAlertsRepo, s.logger))

		v1.GET("/search", handlers.SearchHandler(s.db, s.cfg.AIServiceURL, s.aiHTTPClient, s.logger))
		v1.POST("/query", handlers.UnifiedQueryHandler(s.cfg.AIServiceURL, services.NewSecureHTTPClient(s.cfg.InternalAPISecret, 90*time.Second), s.jobService, s.logger))
		v1.GET("/jobs/:id", handlers.GetJobStatusHandler(s.jobService, s.logger))
		v1.GET("/ask", handlers.AskHandler(s.db, s.cfg.AIServiceURL, s.cfg.OllamaURL, s.cfg.LLMModel, nil, s.logger))
		v1.GET("/files/content", handlers.FileContentHandler(s.logger))

		// System capabilities
		v1.GET("/system/capabilities", handlers.Capabilities(s.benchmarkService))

		// AI Settings
		v1.GET("/ai/status", handlers.AIStatusHandler(s.cfg.AIServiceURL, s.aiHTTPClient, s.logger))
		v1.GET("/ai/settings", handlers.AISettingsGetHandler(s.logger))
		v1.POST("/ai/settings", handlers.AISettingsSaveHandler(s.logger))
		v1.POST("/ai/reindex", handlers.AIReindexHandler(s.cfg.AIServiceURL, s.aiHTTPClient, s.logger))
		v1.POST("/ai/warmup", handlers.AIWarmupHandler(s.logger))

		// Network Settings
		v1.GET("/network/settings", handlers.NetworkSettingsGetHandler(s.logger))
		v1.PUT("/network/settings", handlers.NetworkSettingsSaveHandler(s.logger))

		// Backup Settings
		v1.GET("/backup/settings", handlers.BackupSettingsGetHandler(s.logger))
		v1.PUT("/backup/settings", handlers.BackupSettingsSaveHandler(s.logger))

		// Security Settings
		v1.GET("/security/settings", handlers.SecuritySettingsGetHandler(s.logger))
		v1.PUT("/security/settings", handlers.SecuritySettingsSaveHandler(s.logger))

		// Storage Settings
		v1.GET("/storage/settings", handlers.StorageSettingsGetHandler(s.logger))
		v1.PUT("/storage/settings", handlers.StorageSettingsSaveHandler(s.logger))

		// Setup endpoints
		v1.GET("/system/setup-status", handlers.SetupStatusHandler(s.logger))

		// Vault endpoints (public)
		v1.GET("/vault/status", handlers.VaultStatusHandler(s.encryptionService))
		v1.GET("/vault/config", handlers.VaultConfigGetHandler(s.encryptionService))
	}

	// Protected system settings
	settingsV1 := v1.Group("/system")
	settingsV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
	)
	{
		settingsV1.GET("/settings", handlers.SystemSettingsHandler(s.settingsService))
		settingsV1.PUT("/settings/backup", handlers.UpdateBackupSettingsHandler(s.settingsService))
		settingsV1.POST("/validate-path", handlers.ValidatePathHandler(s.settingsService))

		// Vault management
		settingsV1.POST("/vault/setup", handlers.VaultSetupHandler(s.encryptionService, s.logger))
		settingsV1.POST("/vault/unlock", handlers.VaultUnlockHandler(s.encryptionService, s.logger))
		settingsV1.POST("/vault/lock", handlers.VaultLockHandler(s.encryptionService, s.logger))
		settingsV1.POST("/vault/panic", handlers.VaultPanicHandler(s.encryptionService, s.logger))
		settingsV1.PUT("/vault/config", handlers.VaultConfigUpdateHandler(s.encryptionService, s.logger))
		settingsV1.GET("/vault/export-config", handlers.VaultExportConfigHandler(s.encryptionService, s.logger))

		// Setup wizard
		settingsV1.POST("/setup", handlers.SetupHandler(s.logger))
	}

	// Vault uploads (chunked)
	vaultUploadV1 := v1.Group("/vault/upload")
	vaultUploadV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
	)
	{
		vaultUploadV1.POST("/init", s.blobStorageHandler.InitUpload)
		vaultUploadV1.POST("/chunk/:id", s.blobStorageHandler.UploadChunk)
		vaultUploadV1.POST("/finalize/:id", s.blobStorageHandler.FinalizeUpload)
	}
}

// setupStorageRoutes configures storage endpoints
func (s *Server) setupStorageRoutes() {
	storageV1 := s.router.Group("/api/v1/storage")
	storageV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
	)
	{
		storageV1.GET("/files", handlers.StorageListHandler(s.storageService, s.logger))
		storageV1.POST("/upload", handlers.StorageUploadHandler(s.storageService, s.encryptionPolicyService, s.honeyfileService, s.aiAgentService, s.logger))
		storageV1.GET("/download", handlers.StorageDownloadHandler(s.storageService, s.honeyfileService, s.logger))
		storageV1.GET("/smart-download", handlers.SmartDownloadHandler(s.storageService, s.honeyfileService, s.contentDeliveryService, s.logger))
		storageV1.GET("/download-zip", handlers.StorageDownloadZipHandler(s.storageService, s.logger))
		storageV1.POST("/batch-download", handlers.StorageBatchDownloadHandler(s.storageService, s.logger))
		storageV1.DELETE("/delete", handlers.StorageDeleteHandler(s.storageService, s.aiAgentService, s.logger))
		storageV1.GET("/trash", handlers.StorageTrashListHandler(s.storageService, s.logger))
		storageV1.POST("/trash/restore/:id", handlers.StorageTrashRestoreHandler(s.storageService, s.logger))
		storageV1.DELETE("/trash/:id", handlers.StorageTrashDeleteHandler(s.storageService, s.logger))
		storageV1.POST("/trash/empty", handlers.StorageTrashEmptyHandler(s.storageService, s.logger))
		storageV1.POST("/rename", handlers.StorageRenameHandler(s.storageService, s.logger))
		storageV1.POST("/move", handlers.StorageMoveHandler(s.storageService, s.logger))
		storageV1.POST("/mkdir", handlers.StorageMkdirHandler(s.storageService, s.logger))
		storageV1.POST("/upload-zip", handlers.StorageUploadZipHandler(s.storageService, s.archiveService, s.logger))
	}
}

// setupEncryptedStorageRoutes configures encrypted storage endpoints
func (s *Server) setupEncryptedStorageRoutes() {
	if s.encryptedStorageService == nil {
		return
	}

	encV1 := s.router.Group("/api/v1/encrypted")
	encV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
	)
	{
		encV1.GET("/status", handlers.EncryptedStorageStatusHandler(s.encryptedStorageService))
		encV1.GET("/files", handlers.EncryptedStorageListHandler(s.encryptedStorageService, s.logger))
		encV1.POST("/upload", handlers.EncryptedStorageUploadHandler(s.encryptedStorageService, s.secureAIFeeder, s.logger))
		encV1.GET("/download", handlers.EncryptedStorageDownloadHandler(s.encryptedStorageService, s.logger))
		encV1.GET("/preview", handlers.EncryptedStoragePreviewHandler(s.encryptedStorageService, s.logger))
		encV1.DELETE("/delete", handlers.EncryptedStorageDeleteHandler(s.encryptedStorageService, s.logger))
	}
}

// setupBackupRoutes configures backup endpoints
func (s *Server) setupBackupRoutes() {
	backupV1 := s.router.Group("/api/v1/backups")
	backupV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
	)
	{
		backupV1.GET("", handlers.BackupListHandler(s.backupService, s.logger))
		backupV1.POST("", handlers.BackupCreateHandler(s.backupService, s.cfg, s.logger))

		backupV1.POST("/:id/restore",
			middleware.AdminOnly(s.userRepo, s.logger),
			handlers.BackupRestoreHandler(s.backupService, s.cfg, s.logger),
		)
		backupV1.DELETE("/:id",
			middleware.AdminOnly(s.userRepo, s.logger),
			handlers.BackupDeleteHandler(s.backupService, s.logger),
		)
	}
}

// setupAdminRoutes configures admin-only endpoints
func (s *Server) setupAdminRoutes() {
	adminV1 := s.router.Group("/api/v1/admin")
	adminV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
		middleware.AdminOnly(s.userRepo, s.logger),
	)
	{
		adminV1.GET("/settings", handlers.GetAdminSettingsHandler(s.cfg, s.logger))
		adminV1.PUT("/settings", handlers.UpdateAdminSettingsHandler(s.cfg, s.settingsRepo, s.logger))
		adminV1.GET("/status", handlers.SystemStatusHandler(s.db, s.logger))
		adminV1.GET("/users", handlers.UserListHandler(s.userRepo, s.logger))
		adminV1.PUT("/users/:id/role", handlers.UpdateUserRoleHandler(s.userRepo, s.logger))
		adminV1.POST("/maintenance", handlers.ToggleMaintenanceModeHandler(s.logger))
		adminV1.GET("/audit-logs", handlers.AuditLogHandler(s.db, s.logger))
		adminV1.POST("/system/reconcile-knowledge", handlers.ReconcileKnowledgeHandler(s.secureAIFeeder, "/mnt/data", s.logger))
	}
}

// setupSystemRoutes configures stealth system routes
func (s *Server) setupSystemRoutes() {
	sysV1 := s.router.Group("/api/v1/sys")
	sysV1.Use(
		middleware.AuthMiddleware(s.jwtService, s.redis, s.logger),
		middleware.CSRFMiddleware(s.redis, s.logger),
		middleware.AdminOnly(s.userRepo, s.logger),
	)
	{
		sysV1.POST("/integrity/checkpoints", handlers.CreateCheckpointHandler(s.honeyfileService, s.logger))
	}
}
