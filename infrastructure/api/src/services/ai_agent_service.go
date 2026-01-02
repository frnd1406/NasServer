package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// =============================================================================
// AI-SPECIFIC ERRORS (wrapping generic resilient client errors)
// =============================================================================

var (
	// ErrAIUnavailable indicates the AI service is not reachable
	ErrAIUnavailable = errors.New("AI service unavailable")

	// ErrAITimeout indicates the AI service did not respond in time
	ErrAITimeout = errors.New("AI service timeout")

	// ErrAIBadResponse indicates the AI service returned an error response
	ErrAIBadResponse = errors.New("AI service returned error")
)

// =============================================================================
// CONFIGURATION
// =============================================================================

const (
	// DefaultAITimeout is generous for slow CPU-only NAS devices
	DefaultAITimeout = 120 * time.Second

	// DefaultAIBaseURL is the internal Docker network address
	DefaultAIBaseURL = "http://ai-knowledge-agent:5000"
)

// AIAgentConfig holds configuration for the AI Agent Service
type AIAgentConfig struct {
	BaseURL        string
	Timeout        time.Duration
	InternalSecret string
	RetryConfig    RetryConfig
}

// DefaultAIAgentConfig returns sensible defaults
func DefaultAIAgentConfig() AIAgentConfig {
	return AIAgentConfig{
		BaseURL:     DefaultAIBaseURL,
		Timeout:     DefaultAITimeout,
		RetryConfig: DefaultRetryConfig(),
	}
}

// =============================================================================
// PAYLOAD TYPES
// =============================================================================

// AIAgentPayload represents the data sent to the AI knowledge agent
type AIAgentPayload struct {
	FilePath string `json:"file_path"`
	FileID   string `json:"file_id"`
	MimeType string `json:"mime_type"`
	Content  string `json:"content,omitempty"`
}

// =============================================================================
// AI AGENT SERVICE
// =============================================================================

// AIAgentService handles communication with the Python AI knowledge agent
// Single Responsibility: HTTP communication with AI service
type AIAgentService struct {
	logger     *logrus.Logger
	httpClient *ResilientHTTPClient
	config     AIAgentConfig
	mimePolicy *MimePolicy
}

// NewAIAgentService creates a new AI Agent Service with default configuration
func NewAIAgentService(logger *logrus.Logger, honeySvc *HoneyfileService, internalSecret string) *AIAgentService {
	config := DefaultAIAgentConfig()
	config.InternalSecret = internalSecret
	return NewAIAgentServiceWithConfig(logger, config)
}

// NewAIAgentServiceWithConfig creates a new AI Agent Service with custom configuration
func NewAIAgentServiceWithConfig(logger *logrus.Logger, config AIAgentConfig) *AIAgentService {
	httpClient := NewResilientHTTPClient(config.Timeout, config.RetryConfig, logger)

	return &AIAgentService{
		logger:     logger,
		httpClient: httpClient,
		config:     config,
		mimePolicy: NewMimePolicy(),
	}
}

// =============================================================================
// CORE API METHODS
// =============================================================================

// Ask sends a query to the AI service and returns the response
func (s *AIAgentService) Ask(ctx context.Context, query string, options map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"query": query,
	}
	for k, v := range options {
		payload[k] = v
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	url := s.config.BaseURL + "/ask"

	resp, err := s.httpClient.DoWithRetry(ctx, "Ask", func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		s.setHeaders(req)
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		return nil, s.wrapError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrAIBadResponse, resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

// NotifyUpload sends a fire-and-forget notification to the AI knowledge agent
// NOTE: Caller is responsible for honeyfile checks before calling this
func (s *AIAgentService) NotifyUpload(filePath, fileID, mimeType string, content string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
		defer cancel()

		if err := s.NotifyUploadSync(ctx, filePath, fileID, mimeType, content); err != nil {
			s.logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"file_id":   fileID,
				"error":     err.Error(),
			}).Warn("Failed to notify AI agent of upload")
		}
	}()
}

// NotifyUploadSync sends a synchronous notification to the AI knowledge agent
// NOTE: Caller is responsible for honeyfile checks and MIME type filtering
func (s *AIAgentService) NotifyUploadSync(ctx context.Context, filePath, fileID, mimeType string, content string) error {
	// MIME type filtering is still here for safety, but caller should pre-filter
	if !s.mimePolicy.IsIndexable(mimeType) {
		s.logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"mime_type": mimeType,
		}).Debug("Skipping AI indexing for non-indexable file type")
		return nil
	}

	payload := AIAgentPayload{
		FilePath: filePath,
		FileID:   fileID,
		MimeType: mimeType,
		Content:  content,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := s.config.BaseURL + "/process"

	resp, err := s.httpClient.DoWithRetry(ctx, "NotifyUpload", func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		s.setHeaders(req)
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		return s.wrapError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"file_id":   fileID,
			"mime_type": mimeType,
		}).Info("Triggered AI agent successfully")
		return nil
	}

	return fmt.Errorf("%w: status %d", ErrAIBadResponse, resp.StatusCode)
}

// NotifyDelete sends a deletion notification to the AI knowledge agent
func (s *AIAgentService) NotifyDelete(ctx context.Context, filePath, fileID string) error {
	payload := map[string]string{
		"file_path": filePath,
		"file_id":   fileID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := s.config.BaseURL + "/delete"

	resp, err := s.httpClient.DoWithRetry(ctx, "NotifyDelete", func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		s.setHeaders(req)
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		return s.wrapError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"file_id":   fileID,
		}).Info("AI agent deletion completed successfully")
		return nil
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	return fmt.Errorf("%w: status %d, body: %s", ErrAIBadResponse, resp.StatusCode, string(body))
}

// =============================================================================
// HELPER METHODS
// =============================================================================

// setHeaders applies common headers to requests
func (s *AIAgentService) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if s.config.InternalSecret != "" {
		req.Header.Set("X-Internal-Secret", s.config.InternalSecret)
	}
}

// wrapError converts generic resilient client errors to AI-specific errors
func (s *AIAgentService) wrapError(err error) error {
	if errors.Is(err, ErrServiceUnavailable) || errors.Is(err, ErrMaxRetriesExceeded) {
		return fmt.Errorf("%w: %v", ErrAIUnavailable, err)
	}
	if errors.Is(err, ErrServiceTimeout) {
		return fmt.Errorf("%w: %v", ErrAITimeout, err)
	}
	return err
}

// IsFileIndexable checks if a file type can be indexed by AI
func (s *AIAgentService) IsFileIndexable(mimeType string) bool {
	return s.mimePolicy.IsIndexable(mimeType)
}
