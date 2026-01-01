package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// =============================================================================
// CUSTOM ERRORS (Sentinel Errors for Clean Error Handling)
// =============================================================================

var (
	// ErrAIUnavailable indicates the AI service is not reachable (connection refused, DNS failure)
	ErrAIUnavailable = errors.New("AI service unavailable")

	// ErrAITimeout indicates the AI service did not respond within the configured timeout
	ErrAITimeout = errors.New("AI service timeout")

	// ErrAIOverloaded indicates the AI service is temporarily overloaded (503/504)
	ErrAIOverloaded = errors.New("AI service overloaded")

	// ErrAIBadResponse indicates the AI service returned an unexpected error
	ErrAIBadResponse = errors.New("AI service returned error")

	// ErrAIMaxRetriesExceeded indicates all retry attempts failed
	ErrAIMaxRetriesExceeded = errors.New("AI service: max retries exceeded")
)

// =============================================================================
// CONFIGURATION DEFAULTS (Overridable via AIAgentConfig)
// =============================================================================

const (
	// DefaultAITimeout is generous for slow CPU-only NAS devices (PDF processing, embeddings)
	DefaultAITimeout = 120 * time.Second

	// DefaultAIBaseURL is the internal Docker network address
	DefaultAIBaseURL = "http://ai-knowledge-agent:5000"

	// Retry Configuration
	DefaultMaxRetries     = 3
	DefaultInitialBackoff = 1 * time.Second
	DefaultMaxBackoff     = 30 * time.Second
	DefaultBackoffFactor  = 2.0
)

// =============================================================================
// AI AGENT SERVICE CONFIGURATION
// =============================================================================

// AIAgentConfig holds all configurable parameters for the AI Agent Service
type AIAgentConfig struct {
	// BaseURL of the Python AI service
	BaseURL string

	// Timeout for individual HTTP requests (should be HIGH for slow hardware)
	Timeout time.Duration

	// InternalSecret for service-to-service authentication
	InternalSecret string

	// Retry settings
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

// DefaultAIAgentConfig returns sensible defaults for variable hardware
func DefaultAIAgentConfig() AIAgentConfig {
	return AIAgentConfig{
		BaseURL:        DefaultAIBaseURL,
		Timeout:        DefaultAITimeout,
		MaxRetries:     DefaultMaxRetries,
		InitialBackoff: DefaultInitialBackoff,
		MaxBackoff:     DefaultMaxBackoff,
		BackoffFactor:  DefaultBackoffFactor,
	}
}

// =============================================================================
// AI-INDEXABLE MIME TYPES
// =============================================================================

// aiIndexableMimeTypes defines which file types can be processed by the AI
var aiIndexableMimeTypes = map[string]bool{
	"text/plain":      true,
	"application/pdf": true,
	"text/markdown":   true,
	"text/csv":        true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/msword": true,
}

// =============================================================================
// PAYLOAD TYPES
// =============================================================================

// AIAgentPayload represents the data sent to the AI knowledge agent
type AIAgentPayload struct {
	FilePath string `json:"file_path"`
	FileID   string `json:"file_id"`
	MimeType string `json:"mime_type"`
	Content  string `json:"content,omitempty"` // Optional: if set, agent uses this instead of reading from disk
}

// =============================================================================
// AI AGENT SERVICE
// =============================================================================

// AIAgentService handles communication with the Python AI knowledge agent
// with built-in resilience for variable hardware performance
type AIAgentService struct {
	logger   *logrus.Logger
	honeySvc *HoneyfileService
	client   *http.Client
	config   AIAgentConfig
}

// NewAIAgentService creates a new AI Agent Service with default configuration
func NewAIAgentService(logger *logrus.Logger, honeySvc *HoneyfileService, internalSecret string) *AIAgentService {
	config := DefaultAIAgentConfig()
	config.InternalSecret = internalSecret
	return NewAIAgentServiceWithConfig(logger, honeySvc, config)
}

// NewAIAgentServiceWithConfig creates a new AI Agent Service with custom configuration
func NewAIAgentServiceWithConfig(logger *logrus.Logger, honeySvc *HoneyfileService, config AIAgentConfig) *AIAgentService {
	// Create HTTP client with configured timeout
	// Note: Per-request timeout is handled via context, this is the absolute max
	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			// Connection pooling for performance
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     90 * time.Second,
			// Timeouts for connection establishment
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second, // Connection timeout (not request timeout!)
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}

	return &AIAgentService{
		logger:   logger,
		honeySvc: honeySvc,
		client:   client,
		config:   config,
	}
}

// =============================================================================
// EXPONENTIAL BACKOFF RETRY HELPER
// =============================================================================

// retryWithBackoff executes the given function with exponential backoff
// Returns the result or ErrAIMaxRetriesExceeded if all attempts fail
func (s *AIAgentService) retryWithBackoff(ctx context.Context, operation string, fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error
	backoff := s.config.InitialBackoff

	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		// Check context before each attempt
		if ctx.Err() != nil {
			return nil, fmt.Errorf("%w: %v", ErrAITimeout, ctx.Err())
		}

		if attempt > 0 {
			s.logger.WithFields(logrus.Fields{
				"operation": operation,
				"attempt":   attempt,
				"backoff":   backoff.String(),
			}).Debug("Retrying AI agent request after backoff")

			// Wait with context awareness (can be cancelled)
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("%w: %v", ErrAITimeout, ctx.Err())
			case <-time.After(backoff):
			}

			// Exponential backoff with cap
			backoff = time.Duration(float64(backoff) * s.config.BackoffFactor)
			if backoff > s.config.MaxBackoff {
				backoff = s.config.MaxBackoff
			}
		}

		resp, err := fn()
		if err != nil {
			lastErr = s.classifyError(err)

			// Only retry on transient errors
			if !s.isRetryable(lastErr) {
				return nil, lastErr
			}

			s.logger.WithFields(logrus.Fields{
				"operation": operation,
				"attempt":   attempt + 1,
				"error":     err.Error(),
			}).Warn("AI agent request failed, will retry")
			continue
		}

		// Check for retryable HTTP status codes
		if s.isRetryableStatusCode(resp.StatusCode) {
			lastErr = fmt.Errorf("%w: status %d", ErrAIOverloaded, resp.StatusCode)
			resp.Body.Close()

			s.logger.WithFields(logrus.Fields{
				"operation":   operation,
				"attempt":     attempt + 1,
				"status_code": resp.StatusCode,
			}).Warn("AI agent returned retryable status, will retry")
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("%w after %d attempts: %v", ErrAIMaxRetriesExceeded, s.config.MaxRetries+1, lastErr)
}

// classifyError converts a network error to a sentinel error
func (s *AIAgentService) classifyError(err error) error {
	if err == nil {
		return nil
	}

	// Check for context cancellation/timeout
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", ErrAITimeout, err)
	}

	// Check for connection refused (service not running)
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		if netErr.Op == "dial" {
			return fmt.Errorf("%w: %v", ErrAIUnavailable, err)
		}
	}

	// Check for DNS errors
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return fmt.Errorf("%w: DNS lookup failed: %v", ErrAIUnavailable, err)
	}

	// Generic network error
	return fmt.Errorf("%w: %v", ErrAIUnavailable, err)
}

// isRetryable determines if an error warrants a retry
func (s *AIAgentService) isRetryable(err error) bool {
	// Retry on unavailable/overloaded, but NOT on timeout (user cancelled)
	return errors.Is(err, ErrAIUnavailable) || errors.Is(err, ErrAIOverloaded)
}

// isRetryableStatusCode returns true for HTTP status codes that warrant retry
func (s *AIAgentService) isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusServiceUnavailable, // 503
		http.StatusGatewayTimeout,  // 504
		http.StatusTooManyRequests: // 429
		return true
	default:
		return false
	}
}

// =============================================================================
// CORE API METHODS
// =============================================================================

// Ask sends a query to the AI service and returns the response
// This is the main RAG/Chat endpoint with full context support
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

	resp, err := s.retryWithBackoff(ctx, "Ask", func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		s.setHeaders(req)
		return s.client.Do(req)
	})
	if err != nil {
		return nil, err
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
// Runs asynchronously to avoid blocking the upload response
func (s *AIAgentService) NotifyUpload(filePath, fileID, mimeType string, content string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
		defer cancel()

		if err := s.notifyUploadInternal(ctx, filePath, fileID, mimeType, content); err != nil {
			s.logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"file_id":   fileID,
				"error":     err.Error(),
			}).Warn("Failed to notify AI agent of upload")
		}
	}()
}

// notifyUploadInternal is the synchronous implementation for testing
func (s *AIAgentService) notifyUploadInternal(ctx context.Context, filePath, fileID, mimeType string, content string) error {
	// SECURITY: Never index monitored resources (honeyfiles)
	if s.honeySvc != nil && s.honeySvc.CheckAndTrigger(ctx, filePath, RequestMetadata{Action: "index_scan", IPAddress: "internal"}) {
		s.logger.WithField("file_path", filePath).Info("Skipping AI indexing for monitored resource")
		return nil
	}

	// Check if file is eligible for AI indexing
	if !aiIndexableMimeTypes[mimeType] {
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

	resp, err := s.retryWithBackoff(ctx, "NotifyUpload", func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		s.setHeaders(req)
		return s.client.Do(req)
	})
	if err != nil {
		return err
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

// NotifyDelete sends a SYNCHRONOUS deletion notification to the AI knowledge agent
// This prevents "ghost knowledge" by removing embeddings when files are deleted
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

	resp, err := s.retryWithBackoff(ctx, "NotifyDelete", func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		s.setHeaders(req)
		return s.client.Do(req)
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"file_id":   fileID,
		}).Info("AI agent deletion completed successfully")
		return nil
	}

	// Non-2xx status - return error but caller should soft-fail
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

// IsIndexable checks if a MIME type is eligible for AI indexing
func IsIndexable(mimeType string) bool {
	return aiIndexableMimeTypes[mimeType]
}

// CalculateBackoff computes the backoff duration for a given attempt
// Exported for testing purposes
func CalculateBackoff(attempt int, initial, max time.Duration, factor float64) time.Duration {
	if attempt <= 0 {
		return 0
	}
	backoff := float64(initial) * math.Pow(factor, float64(attempt-1))
	if time.Duration(backoff) > max {
		return max
	}
	return time.Duration(backoff)
}
