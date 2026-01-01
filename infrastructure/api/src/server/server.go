package server

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
	"github.com/nas-ai/api/src/handlers/files"
	"github.com/nas-ai/api/src/middleware"
	"github.com/nas-ai/api/src/repository"
	"github.com/nas-ai/api/src/scheduler"
	"github.com/nas-ai/api/src/services"
	"github.com/nas-ai/api/src/services/storage"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Server holds all dependencies for the API server (Clean Architecture)
type Server struct {
	cfg    *config.Config
	logger *logrus.Logger
	router *gin.Engine
	db     *database.DB
	redis  *database.RedisClient
	dbx    *sqlx.DB

	// Repositories
	userRepo          *repository.UserRepository
	settingsRepo      *repository.SystemSettingsRepository
	systemMetricsRepo *repository.SystemMetricsRepository
	systemAlertsRepo  *repository.SystemAlertsRepository
	monitoringRepo    *repository.MonitoringRepository
	embeddingsRepo    *repository.FileEmbeddingsRepository
	honeyfileRepo     *repository.HoneyfileRepository
	fileRepo          *repository.FileRepository

	// Services
	jwtService              *services.JWTService
	passwordService         *services.PasswordService
	tokenService            *services.TokenService
	emailService            *services.EmailService
	backupService           *services.BackupService
	settingsService         *services.SettingsService
	encryptionService       *services.EncryptionService
	storageService          *services.StorageManager
	encryptedStorageService *services.EncryptedStorageService
	secureAIFeeder          *services.SecureAIFeeder
	aiHTTPClient            *http.Client
	jobService              *services.JobService
	honeyfileService        *services.HoneyfileService
	encryptionPolicyService *services.EncryptionPolicyService
	archiveService          *services.ArchiveService
	contentDeliveryService  *services.ContentDeliveryService
	aiAgentService          *services.AIAgentService
	alertService            *services.AlertService
	benchmarkService        *services.BenchmarkService
	consistencyService      *services.ConsistencyService

	// Handlers
	blobStorageHandler *files.BlobStorageHandler
}

// NewServer creates and initializes all server dependencies
func NewServer(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	s := &Server{
		cfg:    cfg,
		logger: logger,
	}

	if err := s.initDatabase(); err != nil {
		return nil, fmt.Errorf("database init failed: %w", err)
	}

	if err := s.initRepositories(); err != nil {
		return nil, fmt.Errorf("repository init failed: %w", err)
	}

	if err := s.initServices(); err != nil {
		return nil, fmt.Errorf("service init failed: %w", err)
	}

	s.initHandlers()
	s.initRouter()
	s.SetupRoutes()
	s.startBackgroundWorkers()

	return s, nil
}

// initDatabase establishes database connections
func (s *Server) initDatabase() error {
	var err error

	s.db, err = database.NewPostgresConnection(s.cfg, s.logger)
	if err != nil {
		return fmt.Errorf("postgres connection failed: %w", err)
	}

	s.redis, err = database.NewRedisConnection(s.cfg, s.logger)
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	// Health check
	if err := s.checkDependencies(); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	s.dbx = sqlx.NewDb(s.db.DB, "postgres")
	return nil
}

// initRepositories initializes all data access layers
func (s *Server) initRepositories() error {
	var err error

	s.settingsRepo = repository.NewSystemSettingsRepository(s.dbx, s.logger)
	if err = s.settingsRepo.EnsureTable(context.Background()); err != nil {
		return fmt.Errorf("settings table init failed: %w", err)
	}
	s.applyPersistedBackupSettings()

	s.userRepo = repository.NewUserRepository(s.db, s.logger)
	s.systemMetricsRepo = repository.NewSystemMetricsRepository(s.dbx, s.logger)
	s.systemAlertsRepo = repository.NewSystemAlertsRepository(s.dbx, s.logger)
	s.monitoringRepo = repository.NewMonitoringRepository(s.db, s.logger)
	s.embeddingsRepo = repository.NewFileEmbeddingsRepository(s.dbx, s.logger)
	s.honeyfileRepo = repository.NewHoneyfileRepository(s.dbx, s.logger)

	if err = s.honeyfileRepo.EnsureTable(context.Background()); err != nil {
		return fmt.Errorf("honeyfiles table init failed: %w", err)
	}

	s.fileRepo = repository.NewFileRepository(s.dbx, s.logger)

	return nil
}

// initServices initializes all business logic services
func (s *Server) initServices() error {
	var err error

	// Auth Services
	s.jwtService, err = services.NewJWTService(s.cfg, s.logger)
	if err != nil {
		return fmt.Errorf("JWT service init failed: %w", err)
	}
	s.passwordService = services.NewPasswordService()
	s.tokenService = services.NewTokenService(s.redis, s.logger)
	s.emailService = services.NewEmailService(s.cfg, s.logger)

	// Backup Service
	s.backupService, err = services.NewBackupService("/mnt/data", s.cfg.BackupStoragePath, s.logger)
	if err != nil {
		return fmt.Errorf("backup service init failed: %w", err)
	}

	// Settings Service
	onRestartScheduler := func() error {
		return scheduler.RestartScheduler()
	}
	s.settingsService = services.NewSettingsService(s.cfg, s.settingsRepo, s.backupService, onRestartScheduler, s.logger)
	s.logger.Info("SettingsService initialized")

	// Encryption Service (Zero-Knowledge)
	vaultPath := "/tmp/nas-vault-demo"
	if s.cfg.Environment == "production" {
		if _, err := os.Stat("/var/lib/nas/vault"); err == nil {
			vaultPath = "/var/lib/nas/vault"
			s.logger.Warn("⚠️  Vault persistence enabled: Keys survive restarts (security trade-off)")
		}
	}
	s.encryptionService = services.NewEncryptionService(vaultPath, s.logger)
	s.logger.WithField("vaultPath", vaultPath).Info("Encryption service initialized")

	// Storage Services
	localStore, err := storage.NewLocalStore("/mnt/data")
	if err != nil {
		return fmt.Errorf("local store init failed: %w", err)
	}
	s.storageService = services.NewStorageManager(localStore, s.encryptionService, s.fileRepo, s.logger)

	// Encrypted Storage (optional)
	encryptedStoragePath := "/media/frnd14/DEMO"
	s.encryptedStorageService, err = services.NewEncryptedStorageService(
		s.storageService,
		s.encryptionService,
		encryptedStoragePath,
		s.logger,
	)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to initialize encrypted storage service (non-fatal)")
	} else {
		s.logger.WithField("path", encryptedStoragePath).Info("Encrypted storage service initialized")
	}

	// AI Services
	if s.encryptionService != nil {
		s.secureAIFeeder = services.NewSecureAIFeeder(
			s.encryptionService,
			s.cfg.AIServiceURL,
			s.cfg.InternalAPISecret,
			s.logger,
		)
		s.logger.Info("SecureAIFeeder initialized")
	}

	s.aiHTTPClient = services.NewSecureHTTPClient(s.cfg.InternalAPISecret, 15*time.Second)
	s.jobService = services.NewJobService(s.redis, s.logger)
	if err := s.jobService.EnsureConsumerGroup(context.Background()); err != nil {
		s.logger.WithError(err).Warn("Failed to ensure AI job consumer group (non-fatal)")
	}
	s.logger.Info("JobService initialized")

	// Security Services
	s.honeyfileService = services.NewHoneyfileService(s.honeyfileRepo, s.encryptionService, s.logger)
	s.logger.Info("HoneyfileService initialized")

	s.encryptionPolicyService = services.NewEncryptionPolicyService()
	s.logger.Info("EncryptionPolicyService initialized")

	// Utility Services
	s.archiveService = services.NewArchiveService(s.logger)
	s.logger.Info("ArchiveService initialized")

	s.contentDeliveryService = services.NewContentDeliveryService(s.storageService, s.encryptionService, s.logger)
	s.logger.Info("ContentDeliveryService initialized")

	s.aiAgentService = services.NewAIAgentService(s.logger, s.honeyfileService, s.cfg.InternalAPISecret)
	s.logger.Info("AIAgentService initialized")

	s.alertService = services.NewAlertService(s.emailService, s.cfg, s.logger)
	s.logger.Info("AlertService initialized")

	s.benchmarkService = services.NewBenchmarkService(s.logger)

	// Consistency Service
	s.consistencyService = services.NewConsistencyService(
		s.dbx,
		s.embeddingsRepo,
		"/mnt/data",
		time.Duration(s.cfg.ConsistencyCheckIntervalMin)*time.Minute,
		s.logger,
	)

	// Initial reconciliation
	s.logger.Info("Running initial consistency reconciliation...")
	if err := s.consistencyService.RunReconciliation(context.Background()); err != nil {
		s.logger.WithError(err).Warn("Initial reconciliation failed (non-fatal)")
	}

	return nil
}

// initHandlers initializes HTTP handlers
func (s *Server) initHandlers() {
	s.blobStorageHandler = files.NewBlobStorageHandler(s.storageService, s.logger)
	s.logger.Info("BlobStorageHandler initialized")
}

// initRouter creates and configures the Gin router
func (s *Server) initRouter() {
	if s.cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.New()

	// Middleware chain (Onion Principle)
	rateLimiter := middleware.NewRateLimiter(s.cfg)
	s.router.Use(
		middleware.PanicRecovery(s.logger),
		middleware.RequestID(),
		middleware.GinSecureHeaders(),
		middleware.CORS(s.cfg, s.logger),
		rateLimiter.Middleware(),
		middleware.AuditLogger(s.logger),
	)

	// OPTIONS preflight
	s.router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Environment context
	s.router.Use(func(c *gin.Context) {
		c.Set("environment", s.cfg.Environment)
		c.Next()
	})
}

// startBackgroundWorkers starts all background goroutines
func (s *Server) startBackgroundWorkers() {
	// Alert ticker
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.alertService.RunSystemChecks()
		}
	}()

	// Startup benchmark
	go func() {
		if err := s.benchmarkService.RunStartupBenchmark(); err != nil {
			s.logger.WithError(err).Warn("Startup benchmark failed (non-fatal)")
		}
	}()

	// Backup scheduler
	go func() {
		if err := scheduler.StartBackupScheduler(s.backupService, s.cfg); err != nil {
			s.logger.WithError(err).Error("Failed to start backup scheduler")
		}
	}()

	// Consistency worker
	go s.consistencyService.Start(context.Background())
}

// Run starts the HTTP server and waits for shutdown signal
func (s *Server) Run() error {
	secureHandler := middleware.SecureHeaders(s.router)

	srv := &http.Server{
		Addr:           "0.0.0.0:" + s.cfg.Port,
		Handler:        secureHandler,
		ReadTimeout:    600 * time.Second,
		WriteTimeout:   600 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	go func() {
		s.logger.WithField("port", s.cfg.Port).Info("Server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")
	s.consistencyService.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.WithError(err).Error("Server forced to shutdown")
		return err
	}

	s.logger.Info("Server exited")
	return nil
}

// Close cleans up all resources
func (s *Server) Close() {
	if s.db != nil {
		s.db.Close()
	}
	if s.redis != nil {
		s.redis.Close()
	}
}

// checkDependencies verifies database connectivity
func (s *Server) checkDependencies() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres unreachable: %w", err)
	}

	if err := s.redis.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis unreachable: %w", err)
	}

	return nil
}

// applyPersistedBackupSettings loads settings from database
func (s *Server) applyPersistedBackupSettings() {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	settings, err := s.settingsRepo.GetAll(context.Background())
	if err != nil {
		s.logger.WithError(err).Warn("Failed to load persisted settings; continuing with defaults")
		return
	}

	if schedule, ok := settings[repository.SystemSettingBackupSchedule]; ok {
		sc := strings.TrimSpace(schedule)
		if sc != "" {
			if _, err := parser.Parse(sc); err != nil {
				s.logger.WithError(err).Warn("Ignoring invalid persisted backup schedule")
			} else {
				s.cfg.BackupSchedule = sc
			}
		}
	}

	if retentionStr, ok := settings[repository.SystemSettingBackupRetention]; ok {
		if n, err := strconv.Atoi(retentionStr); err == nil && n > 0 {
			s.cfg.BackupRetentionCount = n
		} else if err != nil {
			s.logger.WithError(err).Warn("Ignoring invalid persisted backup retention")
		}
	}

	if path, ok := settings[repository.SystemSettingBackupPath]; ok {
		p := filepath.Clean(strings.TrimSpace(path))
		if p != "" && p != "." && p != string(os.PathSeparator) {
			s.cfg.BackupStoragePath = p
		} else {
			s.logger.Warn("Ignoring invalid persisted backup path")
		}
	}
}
