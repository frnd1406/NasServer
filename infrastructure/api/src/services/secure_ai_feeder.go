package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// SecureAIFeeder handles secure content feeding to the AI agent.
// It decrypts encrypted files in-memory and pushes the plaintext
// to the AI agent for indexing, ensuring no plaintext touches the disk.
type SecureAIFeeder struct {
	encryption *EncryptionService
	aiAgentURL string
	httpClient *http.Client
	logger     *logrus.Logger
}

// IngestDirectPayload is the JSON payload for the AI agent's /ingest_direct endpoint
type IngestDirectPayload struct {
	Content  string `json:"content"`
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
	MimeType string `json:"mime_type"`
}

// IngestDirectResponse is the response from the AI agent
type IngestDirectResponse struct {
	Status          string `json:"status"`
	FileID          string `json:"file_id"`
	FilePath        string `json:"file_path"`
	ContentLength   int    `json:"content_length"`
	EmbeddingDim    int    `json:"embedding_dim"`
	EncryptedSource bool   `json:"encrypted_source"`
	Error           string `json:"error,omitempty"`
}

// NewSecureAIFeeder creates a new secure AI feeder service
func NewSecureAIFeeder(
	encryption *EncryptionService,
	aiAgentURL string,
	logger *logrus.Logger,
) *SecureAIFeeder {
	return &SecureAIFeeder{
		encryption: encryption,
		aiAgentURL: aiAgentURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Embedding generation can take time
		},
		logger: logger,
	}
}

// FeedEncryptedFile decrypts a file and pushes the content to the AI agent.
// The plaintext only exists in RAM and is securely wiped after the request.
//
// Parameters:
//   - encryptedPath: Full path to the encrypted file on disk (e.g., /media/frnd14/DEMO/geheim.pdf.enc)
//   - originalPath: Original path for source citations (stored in metadata)
//   - fileID: Unique identifier for the file
//   - mimeType: MIME type of the original file
//
// Returns error if decryption or AI agent communication fails.
func (f *SecureAIFeeder) FeedEncryptedFile(
	encryptedPath string,
	originalPath string,
	fileID string,
	mimeType string,
) error {
	// Check if encryption is available
	if f.encryption == nil || !f.encryption.IsUnlocked() {
		return ErrVaultLocked
	}

	f.logger.WithFields(logrus.Fields{
		"encryptedPath": encryptedPath,
		"originalPath":  originalPath,
		"fileID":        fileID,
	}).Info("SecureAIFeeder: Starting encrypted file ingestion")

	// Step 1: Read encrypted file from disk
	encryptedData, err := os.ReadFile(encryptedPath)
	if err != nil {
		return fmt.Errorf("read encrypted file: %w", err)
	}

	// Step 2: Decrypt in-memory
	plaintext, err := f.encryption.DecryptData(encryptedData)
	if err != nil {
		return fmt.Errorf("decrypt file: %w", err)
	}

	// Ensure plaintext is wiped after we're done (even on error)
	defer func() {
		for i := range plaintext {
			plaintext[i] = 0
		}
		f.logger.Debug("SecureAIFeeder: Plaintext securely wiped from memory")
	}()

	// Step 3: Convert to string (for text files)
	content := string(plaintext)

	if len(content) == 0 {
		f.logger.WithField("fileID", fileID).Warn("SecureAIFeeder: Empty file, skipping")
		return nil
	}

	// Step 4: Build payload for AI agent
	payload := IngestDirectPayload{
		Content:  content,
		FileID:   fileID,
		FilePath: originalPath, // Store original path for citations
		MimeType: mimeType,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	// Step 5: POST to AI agent
	url := f.aiAgentURL + "/ingest_direct"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		f.logger.WithFields(logrus.Fields{
			"url":    url,
			"fileID": fileID,
			"error":  err.Error(),
		}).Warn("SecureAIFeeder: Failed to contact AI agent")
		return fmt.Errorf("AI agent request failed: %w", err)
	}
	defer resp.Body.Close()

	// Step 6: Handle response
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var result IngestDirectResponse
		if err := json.Unmarshal(body, &result); err == nil {
			f.logger.WithFields(logrus.Fields{
				"fileID":        fileID,
				"contentLength": result.ContentLength,
				"embeddingDim":  result.EmbeddingDim,
			}).Info("SecureAIFeeder: Encrypted file indexed successfully")
		}
		return nil
	}

	// Error response
	f.logger.WithFields(logrus.Fields{
		"fileID":     fileID,
		"statusCode": resp.StatusCode,
		"response":   string(body),
	}).Warn("SecureAIFeeder: AI agent returned error")

	return fmt.Errorf("AI agent error (status %d): %s", resp.StatusCode, string(body))
}

// FeedContent pushes raw content (already decrypted) to the AI agent.
// Used when the caller handles decryption themselves.
func (f *SecureAIFeeder) FeedContent(
	content string,
	fileID string,
	filePath string,
	mimeType string,
) error {
	if len(content) == 0 {
		return nil
	}

	payload := IngestDirectPayload{
		Content:  content,
		FileID:   fileID,
		FilePath: filePath,
		MimeType: mimeType,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := f.aiAgentURL + "/ingest_direct"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("AI agent request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		f.logger.WithField("fileID", fileID).Info("SecureAIFeeder: Content indexed successfully")
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("AI agent error (status %d): %s", resp.StatusCode, string(body))
}
