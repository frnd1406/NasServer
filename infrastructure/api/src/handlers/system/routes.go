package system

import (
	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/handlers/files"
	"github.com/nas-ai/api/src/handlers/settings"
	"github.com/nas-ai/api/src/middleware/logic"
	auth_repo "github.com/nas-ai/api/src/repository/auth"
	system_repo "github.com/nas-ai/api/src/repository/system"
	servicesConfig "github.com/nas-ai/api/src/services/config"
	"github.com/nas-ai/api/src/services/content"
	"github.com/nas-ai/api/src/services/operations"
	"github.com/nas-ai/api/src/services/security"
	"github.com/sirupsen/logrus"
)

// Handler holds dependencies for system handlers
type Handler struct {
	db                 *database.DB
	redis              *database.RedisClient
	cfg                *config.Config
	userRepo           *auth_repo.UserRepository
	systemMetricsRepo  *system_repo.SystemMetricsRepository
	systemAlertsRepo   *system_repo.SystemAlertsRepository
	monitoringRepo     *system_repo.MonitoringRepository
	jobService         *operations.JobService
	benchmarkService   *operations.BenchmarkService
	settingsService    *servicesConfig.SettingsService
	encryptionService  *security.EncryptionService
	honeyfileService   *content.HoneyfileService
	diagnosticsHandler *DiagnosticsHandler
	jwtService         *security.JWTService
	tokenService       *security.TokenService
	logger             *logrus.Logger
}

// NewHandler creates a new System Handler
func NewHandler(
	db *database.DB,
	redis *database.RedisClient,
	cfg *config.Config,
	userRepo *auth_repo.UserRepository,
	systemMetricsRepo *system_repo.SystemMetricsRepository,
	systemAlertsRepo *system_repo.SystemAlertsRepository,
	monitoringRepo *system_repo.MonitoringRepository,
	jobService *operations.JobService,
	benchmarkService *operations.BenchmarkService,
	settingsService *servicesConfig.SettingsService,
	encryptionService *security.EncryptionService,
	honeyfileService *content.HoneyfileService,
	diagnosticsHandler *DiagnosticsHandler,
	jwtService *security.JWTService,
	tokenService *security.TokenService,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		db:                 db,
		redis:              redis,
		cfg:                cfg,
		userRepo:           userRepo,
		systemMetricsRepo:  systemMetricsRepo,
		systemAlertsRepo:   systemAlertsRepo,
		monitoringRepo:     monitoringRepo,
		jobService:         jobService,
		benchmarkService:   benchmarkService,
		settingsService:    settingsService,
		encryptionService:  encryptionService,
		honeyfileService:   honeyfileService,
		diagnosticsHandler: diagnosticsHandler,
		jwtService:         jwtService,
		tokenService:       tokenService,
		logger:             logger,
	}
}

// RegisterPublicRoutes registers public system routes
func (h *Handler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", Health(h.db, h.redis, h.logger))
	rg.POST("/monitoring/ingest", MonitoringIngestHandler(h.monitoringRepo, h.cfg.MonitoringToken, h.logger))
}

// RegisterV1Routes registers API v1 system routes
func (h *Handler) RegisterV1Routes(rg *gin.RouterGroup) {
	// Public (Token based)
	rg.POST("/system/metrics", SystemMetricsHandler(h.systemMetricsRepo, h.cfg.MonitoringToken, h.logger))

	// Public (Frontend Logging)
	rg.POST("/system/logs/frontend", FrontendLogHandler(h.logger))

	// Protected Read-Only
	rg.GET("/system/metrics", SystemMetricsListHandler(h.systemMetricsRepo, h.logger))
	rg.GET("/system/metrics/live", SystemMetricsLiveHandler(h.logger))
	rg.GET("/system/alerts", SystemAlertsListHandler(h.systemAlertsRepo, h.logger))
	rg.GET("/system/capabilities", Capabilities(h.benchmarkService))
	rg.GET("/system/setup-status", settings.SetupStatusHandler(h.logger))
	rg.GET("/jobs/:id", GetJobStatusHandler(h.jobService, h.logger)) // Generic job status

	// Protected Write
	rg.POST("/system/alerts", SystemAlertCreateHandler(h.systemAlertsRepo, h.logger))
	rg.POST("/system/alerts/:id/resolve", SystemAlertResolveHandler(h.systemAlertsRepo, h.logger))

	// Protected System Settings & Vault Management
	settingsV1 := rg.Group("/system")
	settingsV1.Use(
		logic.AuthMiddleware(h.jwtService, h.tokenService, h.redis, h.logger),
		logic.CSRFMiddleware(h.redis, h.logger),
	)
	{
		settingsV1.GET("/settings", settings.SystemSettingsHandler(h.settingsService))
		settingsV1.PUT("/settings/backup", settings.UpdateBackupSettingsHandler(h.settingsService))
		settingsV1.POST("/validate-path", settings.ValidatePathHandler(h.settingsService))

		// Vault management
		settingsV1.POST("/vault/setup", files.VaultSetupHandler(h.encryptionService, h.logger))
		settingsV1.POST("/vault/unlock", files.VaultUnlockHandler(h.encryptionService, h.logger))
		settingsV1.POST("/vault/lock", files.VaultLockHandler(h.encryptionService, h.logger))
		settingsV1.POST("/vault/panic", files.VaultPanicHandler(h.encryptionService, h.logger))
		settingsV1.PUT("/vault/config", files.VaultConfigUpdateHandler(h.encryptionService, h.logger))
		settingsV1.GET("/vault/export-config", files.VaultExportConfigHandler(h.encryptionService, h.logger))

		// Setup wizard
		settingsV1.POST("/setup", settings.SetupHandler(h.logger))
	}

	// Integrity Checkpoints (Admin Only)
	sysV1 := rg.Group("/sys")
	sysV1.Use(
		logic.AuthMiddleware(h.jwtService, h.tokenService, h.redis, h.logger),
		logic.CSRFMiddleware(h.redis, h.logger),
		logic.AdminOnly(h.userRepo, h.logger),
	)
	{
		sysV1.POST("/integrity/checkpoints", CreateCheckpointHandler(h.honeyfileService, h.logger))
	}
}

// RegisterSettingsRoutes registers additional settings routes (Network, Backup, Security, Storage Settings)
func (h *Handler) RegisterSettingsRoutes(rg *gin.RouterGroup) {
	// Network Settings
	rg.GET("/network/settings", settings.NetworkSettingsGetHandler(h.logger))
	rg.PUT("/network/settings", settings.NetworkSettingsSaveHandler(h.logger))

	// Backup Settings
	rg.GET("/backup/settings", settings.BackupSettingsGetHandler(h.logger))
	rg.PUT("/backup/settings", settings.BackupSettingsSaveHandler(h.logger))

	// Security Settings
	rg.GET("/security/settings", settings.SecuritySettingsGetHandler(h.logger))
	rg.PUT("/security/settings", settings.SecuritySettingsSaveHandler(h.logger))

	// Storage Settings
	rg.GET("/storage/settings", settings.StorageSettingsGetHandler(h.logger))
	rg.PUT("/storage/settings", settings.StorageSettingsSaveHandler(h.logger))
}

// RegisterDiagnosticsRoutes registers diagnostics routes
func (h *Handler) RegisterDiagnosticsRoutes(rg *gin.RouterGroup) {
	rg.GET("/system/diagnostics",
		logic.AuthMiddleware(h.jwtService, h.tokenService, h.redis, h.logger),
		logic.AdminOnly(h.userRepo, h.logger),
		h.diagnosticsHandler.RunSelfTest,
	)
}
