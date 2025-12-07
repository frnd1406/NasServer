package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AIStatusHandler proxies requests to the AI Knowledge Agent's /status endpoint.
// Returns Ollama connection status, available models, and index stats.
func AIStatusHandler(aiServiceURL string, httpClient *http.Client, logger *logrus.Logger) gin.HandlerFunc {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	baseURL := strings.TrimRight(aiServiceURL, "/")

	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		// Build request to AI agent
		aiReq, err := http.NewRequestWithContext(c.Request.Context(), "GET", baseURL+"/status", nil)
		if err != nil {
			logger.WithError(err).Error("Failed to build AI status request")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		resp, err := client.Do(aiReq)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("Failed to call AI agent /status")
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI service unavailable"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			logger.WithError(err).Error("Failed to read AI agent status response")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	}
}

// AIReindexHandler proxies requests to the AI Knowledge Agent's /reindex endpoint.
// Triggers re-indexing of all files.
func AIReindexHandler(aiServiceURL string, httpClient *http.Client, logger *logrus.Logger) gin.HandlerFunc {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	baseURL := strings.TrimRight(aiServiceURL, "/")

	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		// Build request to AI agent
		aiReq, err := http.NewRequestWithContext(c.Request.Context(), "POST", baseURL+"/reindex", nil)
		if err != nil {
			logger.WithError(err).Error("Failed to build AI reindex request")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		resp, err := client.Do(aiReq)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("Failed to call AI agent /reindex")
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI service unavailable"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			logger.WithError(err).Error("Failed to read AI agent reindex response")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response"})
			return
		}

		logger.WithField("request_id", requestID).Info("AI reindex triggered")
		c.Data(resp.StatusCode, "application/json", body)
	}
}

// AISettingsResponse represents the AI settings
type AISettingsResponse struct {
	LLMModel         string  `json:"llm_model"`
	ClassifierModel  string  `json:"classifier_model"`
	EmbeddingModel   string  `json:"embedding_model"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	ContextDocuments int     `json:"context_documents"`
	AutoIndex        bool    `json:"auto_index"`
	OllamaURL        string  `json:"ollama_url"`
}

// AISettingsGetHandler returns current AI settings
// For now returns defaults - can be extended to read from config/database
func AISettingsGetHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Default settings
		settings := AISettingsResponse{
			LLMModel:         "llama3.2",
			ClassifierModel:  "llama3.2:1b",
			EmbeddingModel:   "mxbai-embed-large",
			Temperature:      0.7,
			MaxTokens:        500,
			ContextDocuments: 10,
			AutoIndex:        true,
			OllamaURL:        "http://host.docker.internal:11434",
		}

		c.JSON(http.StatusOK, settings)
	}
}

// AISettingsSaveHandler saves AI settings
// For now just acknowledges - can be extended to persist to config/database
func AISettingsSaveHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		var settings AISettingsResponse
		if err := json.NewDecoder(c.Request.Body).Decode(&settings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid settings format"})
			return
		}

		logger.WithFields(logrus.Fields{
			"request_id":    requestID,
			"llm_model":     settings.LLMModel,
			"classifier":    settings.ClassifierModel,
			"temperature":   settings.Temperature,
		}).Info("AI settings updated")

		// TODO: Persist settings to database/config file

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Settings saved",
		})
	}
}
