package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UnifiedQueryRequest is the request body for the /query endpoint
type UnifiedQueryRequest struct {
	Query string `json:"query" binding:"required"`
}

// UnifiedQueryHandler proxies requests to the AI Knowledge Agent's /query endpoint.
// This endpoint uses AI to classify the query intent and dynamically determines:
// - Whether to return search results or an AI-generated answer
// - How many results to return based on query type
func UnifiedQueryHandler(aiServiceURL string, httpClient *http.Client, logger *logrus.Logger) gin.HandlerFunc {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second} // Longer timeout for LLM operations
	}
	baseURL := strings.TrimRight(aiServiceURL, "/")

	return func(c *gin.Context) {
		var req UnifiedQueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid 'query' field"})
			return
		}

		if strings.TrimSpace(req.Query) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query cannot be empty"})
			return
		}

		logger.WithField("query", req.Query).Info("Proxying unified query to AI agent")

		// Build request to AI agent
		payload := map[string]string{"query": req.Query}
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			logger.WithError(err).Error("Failed to encode query payload")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		aiReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, baseURL+"/query", &buf)
		if err != nil {
			logger.WithError(err).Error("Failed to create AI agent request")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		aiReq.Header.Set("Content-Type", "application/json")

		// Call AI agent
		resp, err := client.Do(aiReq)
		if err != nil {
			logger.WithError(err).Error("Failed to call AI agent")
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB limit
		if err != nil {
			logger.WithError(err).Error("Failed to read AI agent response")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read AI response"})
			return
		}

		// Forward status code and response
		c.Data(resp.StatusCode, "application/json", body)
	}
}
