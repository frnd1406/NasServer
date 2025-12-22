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
	encryption     *EncryptionService
	aiAgentURL     string
	internalSecret string
	httpClient     *http.Client
	logger         *logrus.Logger
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
	internalSecret string,
	logger *logrus.Logger,
) *SecureAIFeeder {
	return &SecureAIFeeder{
		encryption:     encryption,
		aiAgentURL:     aiAgentURL,
		internalSecret: internalSecret,
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
	url := f.aiAgentURL + "/process"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if f.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", f.internalSecret)
	}

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

	url := f.aiAgentURL + "/process"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if f.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", f.internalSecret)
	}

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

// RemoveDocument removes a document from the AI agent's vector index.
// Used for synchronous deletion when files are removed from storage.
func (f *SecureAIFeeder) RemoveDocument(fileID string, filePath string) error {
	payload := map[string]string{
		"file_id":   fileID,
		"file_path": filePath,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	url := f.aiAgentURL + "/delete"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if f.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", f.internalSecret)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("AI agent unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		f.logger.WithFields(logrus.Fields{
			"fileID":   fileID,
			"filePath": filePath,
		}).Info("SecureAIFeeder: Document removed from index")
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("AI agent error (status %d): %s", resp.StatusCode, string(body))
}

// ListVectorsResponse is the response from the AI agent's /list_vectors endpoint
type ListVectorsResponse struct {
	FileIDs []string `json:"file_ids"`
	Count   int      `json:"count"`
}

// ReconcileIndex identifies and removes orphaned embeddings (ghost knowledge).
// It fetches all file IDs from the vector store and compares them against
// actual files in storage, deleting any that no longer exist.
//
// Parameters:
//   - existingFileIDs: Set of file IDs that currently exist in storage
//
// Returns:
//   - deleted: Number of orphaned embeddings removed
//   - err: Error if reconciliation failed
func (f *SecureAIFeeder) ReconcileIndex(existingFileIDs map[string]bool) (int, error) {
	f.logger.Info("SecureAIFeeder: Starting index reconciliation (garbage collection)")

	// Step 1: Get all file IDs from vector store
	url := f.aiAgentURL + "/list_vectors"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	if f.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", f.internalSecret)
	}
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to contact AI agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("AI agent error (status %d): %s", resp.StatusCode, string(body))
	}

	var vectorResponse ListVectorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&vectorResponse); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	f.logger.WithField("vectorCount", vectorResponse.Count).Info("SecureAIFeeder: Found vectors in index")

	// Step 2: Find zombies (vectors without corresponding files)
	var zombies []string
	for _, fileID := range vectorResponse.FileIDs {
		if !existingFileIDs[fileID] {
			zombies = append(zombies, fileID)
		}
	}

	if len(zombies) == 0 {
		f.logger.Info("SecureAIFeeder: No orphaned embeddings found, index is clean")
		return 0, nil
	}

	f.logger.WithField("zombieCount", len(zombies)).Warn("SecureAIFeeder: Found orphaned embeddings")

	// Step 3: Delete zombies
	deleted := 0
	for _, fileID := range zombies {
		if err := f.RemoveDocument(fileID, ""); err != nil {
			f.logger.WithFields(logrus.Fields{
				"fileID": fileID,
				"error":  err.Error(),
			}).Warn("SecureAIFeeder: Failed to delete orphaned embedding")
			continue
		}
		deleted++
	}

	f.logger.WithFields(logrus.Fields{
		"deleted": deleted,
		"total":   len(zombies),
	}).Info("SecureAIFeeder: Index reconciliation complete")

	return deleted, nil
}
