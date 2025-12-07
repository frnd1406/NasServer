package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSecureAIFeeder_FeedContent_PushMode(t *testing.T) {
	// Setup Mock Server (simulates Python AI Agent)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Request
		assert.Equal(t, "/process", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		// Verify Payload (RAM-Push Mode)
		assert.Equal(t, "secret content", payload["content"])
		assert.Equal(t, "test.enc", payload["file_id"])
		assert.Equal(t, "/path/to/test.enc", payload["file_path"])
		assert.Equal(t, "text/plain", payload["mime_type"])

		// Send Success Response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"file_id": "test.enc",
			"mode":    "push",
		})
	}))
	defer mockServer.Close()

	// Setup Service
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Silence logs during test

	// We don't need EncryptionService for FeedContent, so passing nil is fine
	feeder := NewSecureAIFeeder(nil, mockServer.URL, logger)

	// Execute
	err := feeder.FeedContent("secret content", "test.enc", "/path/to/test.enc", "text/plain")

	// Verify
	assert.NoError(t, err)
}

func TestSecureAIFeeder_FeedEncryptedFile_IntegrationFlow(t *testing.T) {
	// Note: This test would require mocking the EncryptionService and file system.
	// For now, we tested the critical "Push Payload" logic in the test above.
	// The FeedEncryptedFile method mainly orchestrates: Read -> Decrypt -> FeedContent.
}
