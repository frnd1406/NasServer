package ai

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/handlers/settings"
	auth_repo "github.com/nas-ai/api/src/repository/auth"
	"github.com/nas-ai/api/src/services/common"
	"github.com/nas-ai/api/src/services/intelligence"
	"github.com/nas-ai/api/src/services/operations"
	"github.com/sirupsen/logrus"
)

// Handler holds dependencies for AI handlers
type Handler struct {
	db             *database.DB
	cfg            *config.Config
	aiHTTPClient   *http.Client
	jobService     *operations.JobService
	secureAIFeeder *intelligence.SecureAIFeeder
	userRepo       auth_repo.UserRepositoryInterface
	logger         *logrus.Logger
}

// NewHandler creates a new AI Handler
func NewHandler(
	db *database.DB,
	cfg *config.Config,
	aiHTTPClient *http.Client,
	jobService *operations.JobService,
	secureAIFeeder *intelligence.SecureAIFeeder,
	userRepo auth_repo.UserRepositoryInterface,
	logger *logrus.Logger,
) *Handler {
	return &Handler{
		db:             db,
		cfg:            cfg,
		aiHTTPClient:   aiHTTPClient,
		jobService:     jobService,
		secureAIFeeder: secureAIFeeder,
		userRepo:       userRepo,
		logger:         logger,
	}
}

// RegisterV1Routes registers API v1 AI routes
func (h *Handler) RegisterV1Routes(rg *gin.RouterGroup) {
	// Search & Query
	rg.GET("/search", SearchHandler(h.db, h.cfg.AIServiceURL, h.aiHTTPClient, h.logger))
	rg.POST("/query", UnifiedQueryHandler(h.cfg.AIServiceURL, common.NewSecureHTTPClient(h.cfg.InternalAPISecret, 90*time.Second), h.jobService, h.logger))
	rg.GET("/ask", AskHandler(h.db, h.cfg.AIServiceURL, h.cfg.OllamaURL, h.cfg.LLMModel, nil, h.logger))

	// AI Settings
	secureClient := common.NewSecureHTTPClient(h.cfg.InternalAPISecret, 10*time.Second)
	rg.GET("/ai/status", settings.AIStatusHandler(h.cfg.AIServiceURL, secureClient, h.logger))
	rg.GET("/ai/settings", settings.AISettingsGetHandler(h.logger))
	rg.POST("/ai/settings", settings.AISettingsSaveHandler(h.logger))
	rg.POST("/ai/reindex", settings.AIReindexHandler(h.cfg.AIServiceURL, secureClient, h.logger))
	rg.POST("/ai/warmup", settings.AIWarmupHandler(h.logger))
}

// RegisterAdminRoutes registers admin AI routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	// Reconcile Knowledge
	rg.POST("/system/reconcile-knowledge", ReconcileKnowledgeHandler(h.secureAIFeeder, "/mnt/data", h.logger))
}
