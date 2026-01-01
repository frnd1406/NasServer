package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// AI-indexable MIME types (text-based files only)
var aiIndexableMimeTypes = map[string]bool{
	"text/plain":      true,
	"application/pdf": true,
	"text/markdown":   true,
	"text/csv":        true,
}

// AIAgentPayload represents the data sent to the AI knowledge agent
type AIAgentPayload struct {
	FilePath string `json:"file_path"`
	FileID   string `json:"file_id"`
	MimeType string `json:"mime_type"`
	Content  string `json:"content,omitempty"` // Optional: if set, agent uses this instead of reading from disk
}

type AIAgentService struct {
	logger         *logrus.Logger
	honeySvc       *HoneyfileService
	internalSecret string
	client         *http.Client
	baseURL        string
}

func NewAIAgentService(logger *logrus.Logger, honeySvc *HoneyfileService, internalSecret string) *AIAgentService {
	return &AIAgentService{
		logger:         logger,
		honeySvc:       honeySvc,
		internalSecret: internalSecret,
		client:         &http.Client{Timeout: 10 * time.Second},
		baseURL:        "http://ai-knowledge-agent:5000",
	}
}

// NotifyUpload sends a fire-and-forget notification to the AI knowledge agent
func (s *AIAgentService) NotifyUpload(filePath, fileID, mimeType string, content string) {
	// Run asynchronously to avoid blocking the upload response
	go func() {
		// SECURITY: Never index monitored resources
		if s.honeySvc != nil && s.honeySvc.CheckAndTrigger(context.Background(), filePath, RequestMetadata{Action: "index_scan", IPAddress: "internal"}) {
			s.logger.WithField("file_path", filePath).Info("Skipping AI indexing for monitored resource")
			return
		}

		// Check if file is eligible for AI indexing
		if !aiIndexableMimeTypes[mimeType] {
			s.logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"mime_type": mimeType,
			}).Info("Skipping AI indexing for non-text file")
			return
		}

		// If content is empty but mime type is text, log a debug message (Legacy Mode / Disk Read)
		if content == "" {
			s.logger.WithField("file_id", fileID).Debug("Indexing triggered without inline content (Legacy Mode)")
		}

		payload := AIAgentPayload{
			FilePath: filePath,
			FileID:   fileID,
			MimeType: mimeType,
			Content:  content,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			s.logger.WithError(err).Warn("Failed to marshal AI agent payload")
			return
		}

		// AI agent endpoint
		url := s.baseURL + "/process"

		// Create HTTP request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			s.logger.WithError(err).Warn("Failed to create AI agent request")
			return
		}

		req.Header.Set("Content-Type", "application/json")
		// SECURITY: Internal Auth
		if s.internalSecret != "" {
			req.Header.Set("X-Internal-Secret", s.internalSecret)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"url":       url,
				"file_path": filePath,
				"error":     err.Error(),
			}).Warn("Failed to trigger AI agent")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			s.logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"file_id":   fileID,
				"mime_type": mimeType,
			}).Info("Triggered AI agent successfully")
		} else {
			s.logger.WithFields(logrus.Fields{
				"status_code": resp.StatusCode,
				"file_path":   filePath,
			}).Warn("AI agent returned non-200 status")
		}
	}()
}

// NotifyDelete sends a SYNCHRONOUS deletion notification to the AI knowledge agent.
// This prevents "ghost knowledge" by removing embeddings when files are deleted.
// Returns error if AI agent is unreachable, but caller should soft-fail (log, don't block user).
func (s *AIAgentService) NotifyDelete(filePath, fileID string) error {
	payload := map[string]string{
		"file_path": filePath,
		"file_id":   fileID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to marshal AI agent delete payload")
		return fmt.Errorf("marshal payload: %w", err)
	}

	// AI agent delete endpoint
	url := s.baseURL + "/delete"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.logger.WithError(err).Warn("Failed to create AI agent delete request")
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", s.internalSecret)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"url":       url,
			"file_path": filePath,
			"file_id":   fileID,
			"error":     err.Error(),
		}).Warn("Failed to notify AI agent of deletion")
		return fmt.Errorf("AI agent unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"file_id":   fileID,
		}).Info("AI agent deletion completed successfully")
		return nil
	}

	// Non-2xx status - log but treat as soft error
	s.logger.WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
		"file_path":   filePath,
		"file_id":     fileID,
	}).Warn("AI agent delete returned non-2xx status")
	return fmt.Errorf("AI agent returned status %d", resp.StatusCode)
}
