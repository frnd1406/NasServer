package server

import (
	"time"

	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/handlers"
	"github.com/nas-ai/api/src/handlers/ai"
	"github.com/nas-ai/api/src/handlers/auth"
	"github.com/nas-ai/api/src/handlers/files"
	"github.com/nas-ai/api/src/handlers/settings"
	"github.com/nas-ai/api/src/handlers/system"
	"github.com/nas-ai/api/src/middleware/logic"

	"github.com/nas-ai/api/src/services/common"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all HTTP routes (SRP: Routing Only)
func (s *Server) SetupRoutes() {
	// === PUBLIC ROUTES (no auth, but rate-limited) ===
	s.router.GET("/health", system.Health(s.db, s.redis, s.logger))
	s.router.POST("/monitoring/ingest", system.MonitoringIngestHandler(s.monitoringRepo, s.cfg.MonitoringToken, s.logger))

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
		authLimiter := logic.NewRateLimiter(&config.Config{RateLimitPerMin: 5})

		authGroup.POST("/register",
			authLimiter.Middleware(),
			auth.RegisterHandler(s.cfg, s.userRepo, s.jwtService, s.passwordService, s.tokenService, s.emailService, s.redis, s.logger),
		)
		authGroup.POST("/login",
			authLimiter.Middleware(),
			auth.LoginHandler(s.userRepo, s.jwtService, s.passwordService, s.redis, s.logger),
		)
		authGroup.POST("/refresh", auth.RefreshHandler(s.jwtService, s.redis, s.logger))
		authGroup.POST("/logout",
			logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
			auth.LogoutHandler(s.jwtService, s.redis, s.logger),
		)

		// Email verification
		authGroup.POST("/verify-email", auth.VerifyEmailHandler(s.userRepo, s.tokenService, s.emailService, s.logger))
		authGroup.POST("/resend-verification",
			logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
			auth.ResendVerificationHandler(s.userRepo, s.tokenService, s.emailService, s.logger),
		)

		// Password reset
		authGroup.POST("/forgot-password", auth.ForgotPasswordHandler(s.userRepo, s.tokenService, s.emailService, s.logger))
		authGroup.POST("/reset-password", auth.ResetPasswordHandler(s.userRepo, s.tokenService, s.passwordService, s.jwtService, s.redis, s.logger))
	}
}

// setupAPIRoutes configures protected API endpoints
func (s *Server) setupAPIRoutes() {
	apiGroup := s.router.Group("/api")
	apiGroup.Use(logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger))
	apiGroup.Use(logic.CSRFMiddleware(s.redis, s.logger))
	{
		apiGroup.GET("/profile", handlers.ProfileHandler(s.userRepo, s.logger))
		apiGroup.GET("/monitoring", system.MonitoringListHandler(s.monitoringRepo, s.logger))
	}
}

// setupV1Routes configures API v1 endpoints
func (s *Server) setupV1Routes() {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/auth/csrf", auth.GetCSRFToken(s.redis, s.logger))
		v1.POST("/system/metrics", system.SystemMetricsHandler(s.systemMetricsRepo, s.cfg.MonitoringToken, s.logger))
		v1.GET("/system/metrics", system.SystemMetricsListHandler(s.systemMetricsRepo, s.logger))
		v1.GET("/system/metrics/live", system.SystemMetricsLiveHandler(s.logger))
		v1.GET("/system/alerts", system.SystemAlertsListHandler(s.systemAlertsRepo, s.logger))
		v1.POST("/system/alerts", system.SystemAlertCreateHandler(s.systemAlertsRepo, s.logger))
		v1.POST("/system/alerts/:id/resolve", system.SystemAlertResolveHandler(s.systemAlertsRepo, s.logger))

		v1.GET("/search", ai.SearchHandler(s.db, s.cfg.AIServiceURL, s.aiHTTPClient, s.logger))
		v1.POST("/query", ai.UnifiedQueryHandler(s.cfg.AIServiceURL, common.NewSecureHTTPClient(s.cfg.InternalAPISecret, 90*time.Second), s.jobService, s.logger))
		v1.GET("/jobs/:id", system.GetJobStatusHandler(s.jobService, s.logger))
		v1.GET("/ask", ai.AskHandler(s.db, s.cfg.AIServiceURL, s.cfg.OllamaURL, s.cfg.LLMModel, nil, s.logger))
		v1.GET("/files/content", files.FileContentHandler(s.logger))

		// System capabilities
		v1.GET("/system/capabilities", system.Capabilities(s.benchmarkService))

		// AI Settings
		v1.GET("/ai/status", settings.AIStatusHandler(s.cfg.AIServiceURL, s.aiHTTPClient, s.logger))
		v1.GET("/ai/settings", settings.AISettingsGetHandler(s.logger))
		v1.POST("/ai/settings", settings.AISettingsSaveHandler(s.logger))
		v1.POST("/ai/reindex", settings.AIReindexHandler(s.cfg.AIServiceURL, s.aiHTTPClient, s.logger))
		v1.POST("/ai/warmup", settings.AIWarmupHandler(s.logger))

		// Network Settings
		v1.GET("/network/settings", settings.NetworkSettingsGetHandler(s.logger))
		v1.PUT("/network/settings", settings.NetworkSettingsSaveHandler(s.logger))

		// Backup Settings
		v1.GET("/backup/settings", settings.BackupSettingsGetHandler(s.logger))
		v1.PUT("/backup/settings", settings.BackupSettingsSaveHandler(s.logger))

		// Security Settings
		v1.GET("/security/settings", settings.SecuritySettingsGetHandler(s.logger))
		v1.PUT("/security/settings", settings.SecuritySettingsSaveHandler(s.logger))

		// Storage Settings
		v1.GET("/storage/settings", settings.StorageSettingsGetHandler(s.logger))
		v1.PUT("/storage/settings", settings.StorageSettingsSaveHandler(s.logger))

		// Setup endpoints
		v1.GET("/system/setup-status", settings.SetupStatusHandler(s.logger))

		// Vault endpoints (public)
		v1.GET("/vault/status", files.VaultStatusHandler(s.encryptionService))
		v1.GET("/vault/config", files.VaultConfigGetHandler(s.encryptionService))
	}

	// Protected system settings
	settingsV1 := v1.Group("/system")
	settingsV1.Use(
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
	)
	{
		settingsV1.GET("/settings", settings.SystemSettingsHandler(s.settingsService))
		settingsV1.PUT("/settings/backup", settings.UpdateBackupSettingsHandler(s.settingsService))
		settingsV1.POST("/validate-path", settings.ValidatePathHandler(s.settingsService))

		// Vault management
		settingsV1.POST("/vault/setup", files.VaultSetupHandler(s.encryptionService, s.logger))
		settingsV1.POST("/vault/unlock", files.VaultUnlockHandler(s.encryptionService, s.logger))
		settingsV1.POST("/vault/lock", files.VaultLockHandler(s.encryptionService, s.logger))
		settingsV1.POST("/vault/panic", files.VaultPanicHandler(s.encryptionService, s.logger))
		settingsV1.PUT("/vault/config", files.VaultConfigUpdateHandler(s.encryptionService, s.logger))
		settingsV1.GET("/vault/export-config", files.VaultExportConfigHandler(s.encryptionService, s.logger))

		// Setup wizard
		settingsV1.POST("/setup", settings.SetupHandler(s.logger))
	}

	// Vault uploads (chunked)
	vaultUploadV1 := v1.Group("/vault/upload")
	vaultUploadV1.Use(
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
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
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
	)
	{
		storageV1.GET("/files", files.StorageListHandler(s.storageService, s.logger))
		storageV1.POST("/upload", files.StorageUploadHandler(s.storageService, s.encryptionPolicyService, s.honeyfileService, s.aiAgentService, s.logger))
		storageV1.GET("/download", files.StorageDownloadHandler(s.storageService, s.honeyfileService, s.logger))
		storageV1.GET("/smart-download", files.SmartDownloadHandler(s.storageService, s.honeyfileService, s.contentDeliveryService, s.logger))
		storageV1.GET("/download-zip", files.StorageDownloadZipHandler(s.storageService, s.logger))
		storageV1.POST("/batch-download", files.StorageBatchDownloadHandler(s.storageService, s.logger))
		storageV1.DELETE("/delete", files.StorageDeleteHandler(s.storageService, s.aiAgentService, s.logger))
		storageV1.GET("/trash", files.StorageTrashListHandler(s.storageService, s.logger))
		storageV1.POST("/trash/restore/:id", files.StorageTrashRestoreHandler(s.storageService, s.logger))
		storageV1.DELETE("/trash/:id", files.StorageTrashDeleteHandler(s.storageService, s.logger))
		storageV1.POST("/trash/empty", files.StorageTrashEmptyHandler(s.storageService, s.logger))
		storageV1.POST("/rename", files.StorageRenameHandler(s.storageService, s.logger))
		storageV1.POST("/move", files.StorageMoveHandler(s.storageService, s.logger))
		storageV1.POST("/mkdir", files.StorageMkdirHandler(s.storageService, s.logger))
		storageV1.POST("/upload-zip", files.StorageUploadZipHandler(s.storageService, s.archiveService, s.logger))
	}
}

// setupEncryptedStorageRoutes configures encrypted storage endpoints
func (s *Server) setupEncryptedStorageRoutes() {
	if s.encryptedStorageService == nil {
		return
	}

	encV1 := s.router.Group("/api/v1/encrypted")
	encV1.Use(
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
	)
	{
		encV1.GET("/status", files.EncryptedStorageStatusHandler(s.encryptedStorageService))
		encV1.GET("/files", files.EncryptedStorageListHandler(s.encryptedStorageService, s.logger))
		encV1.POST("/upload", files.EncryptedStorageUploadHandler(s.encryptedStorageService, s.secureAIFeeder, s.logger))
		encV1.GET("/download", files.EncryptedStorageDownloadHandler(s.encryptedStorageService, s.logger))
		encV1.GET("/preview", files.EncryptedStoragePreviewHandler(s.encryptedStorageService, s.logger))
		encV1.DELETE("/delete", files.EncryptedStorageDeleteHandler(s.encryptedStorageService, s.logger))
	}
}

// setupBackupRoutes configures backup endpoints
func (s *Server) setupBackupRoutes() {
	backupV1 := s.router.Group("/api/v1/backups")
	backupV1.Use(
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
	)
	{
		backupV1.GET("", handlers.BackupListHandler(s.backupService, s.logger))
		backupV1.POST("", handlers.BackupCreateHandler(s.backupService, s.cfg, s.logger))

		backupV1.POST("/:id/restore",
			logic.AdminOnly(s.userRepo, s.logger),
			handlers.BackupRestoreHandler(s.backupService, s.cfg, s.logger),
		)
		backupV1.DELETE("/:id",
			logic.AdminOnly(s.userRepo, s.logger),
			handlers.BackupDeleteHandler(s.backupService, s.logger),
		)
	}
}

// setupAdminRoutes configures admin-only endpoints
func (s *Server) setupAdminRoutes() {
	adminV1 := s.router.Group("/api/v1/admin")
	adminV1.Use(
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
		logic.AdminOnly(s.userRepo, s.logger),
	)
	{
		adminV1.GET("/settings", settings.GetAdminSettingsHandler(s.cfg, s.logger))
		adminV1.PUT("/settings", settings.UpdateAdminSettingsHandler(s.cfg, s.settingsRepo, s.logger))
		adminV1.GET("/status", settings.SystemStatusHandler(s.db, s.logger))
		adminV1.GET("/users", settings.UserListHandler(s.userRepo, s.logger))
		adminV1.PUT("/users/:id/role", settings.UpdateUserRoleHandler(s.userRepo, s.logger))
		adminV1.POST("/maintenance", settings.ToggleMaintenanceModeHandler(s.logger))
		adminV1.GET("/audit-logs", settings.AuditLogHandler(s.db, s.logger))
		adminV1.POST("/system/reconcile-knowledge", ai.ReconcileKnowledgeHandler(s.secureAIFeeder, "/mnt/data", s.logger))
	}
}

// setupSystemRoutes configures stealth system routes
func (s *Server) setupSystemRoutes() {
	sysV1 := s.router.Group("/api/v1/sys")
	sysV1.Use(
		logic.AuthMiddleware(s.jwtService, s.tokenService, s.redis, s.logger),
		logic.CSRFMiddleware(s.redis, s.logger),
		logic.AdminOnly(s.userRepo, s.logger),
	)
	{
		sysV1.POST("/integrity/checkpoints", system.CreateCheckpointHandler(s.honeyfileService, s.logger))
	}
}
