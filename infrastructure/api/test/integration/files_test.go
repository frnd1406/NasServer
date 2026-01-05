package integration

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/domain/files"
	"github.com/nas-ai/api/src/services/content"
	"github.com/nas-ai/api/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestFileUpload_Success(t *testing.T) {
	// 1. Setup TestEnv with REAL security services
	env := testutils.NewTestEnv(t)

	// 2. Setup Mock Expectations for Storage (only data service that's mocked)
	// Policy determines no encryption (REAL service will be called)
	env.StorageService.On("Save", "/uploads", mock.Anything, mock.Anything).Return(&content.SaveResult{
		Path:      "uploads/test.txt",
		FileID:    "file-123",
		MimeType:  "text/plain",
		SizeBytes: 12,
		Checksum:  "checksum",
	}, nil)

	// AI Notified (mock - external service)
	env.AIService.On("NotifyUpload", "uploads/test.txt", "file-123", "text/plain", "Hello World!").Return()

	// 3. Setup Router with real security + mock storage/AI
	router := testutils.SetupTestRouter(env)

	// 4. Generate REAL token for authentication
	token, err := env.GenerateTestToken("user-1", "test@example.com")
	assert.NoError(t, err)

	// 5. Create Request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("Hello World!"))
	writer.WriteField("path", "/uploads")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 6. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	env.StorageService.AssertExpectations(t)
	env.AIService.AssertExpectations(t)
}

func TestFileContent_Success(t *testing.T) {
	// 1. Setup TestEnv with REAL security services
	env := testutils.NewTestEnv(t)

	// 2. Prepare Temp File for "Open" to return valid *os.File
	tmpFile, err := os.CreateTemp("", "test-content")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("File Content")
	tmpFile.Seek(0, 0)

	fileInfo, _ := tmpFile.Stat()

	// 3. Setup Mock Expectations - Storage returns temp file
	env.StorageService.On("Open", "test/path.txt").Return(tmpFile, fileInfo, "text/plain", nil)

	// 4. Setup Router
	router := testutils.SetupTestRouter(env)

	// 5. Generate REAL token
	token, err := env.GenerateTestToken("user-1", "test@example.com")
	assert.NoError(t, err)

	// 6. Request
	req := httptest.NewRequest("GET", "/api/v1/files/content?path=test/path.txt", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 7. Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "File Content", w.Body.String())
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")

	env.StorageService.AssertExpectations(t)

	tmpFile.Close()
}

func TestEncryptedFileUpload_Success(t *testing.T) {
	t.Skip("Encrypted upload requires EncryptedStorageService implementation - skipping for now")

	// This test would require:
	// 1. Setting up encryption vault with master password
	// 2. Creating an EncryptedStorageService instance
	// 3. Registering the /encrypted routes in the test router
	//
	// For now, we skip until EncryptedStorageService is properly wired up
}

func TestFileDownload_WithHoneyfileCheck(t *testing.T) {
	// 1. Setup TestEnv - Honeyfile service is REAL, storage is MOCK
	env := testutils.NewTestEnv(t)

	// 2. Prepare temp file
	tmpFile, err := os.CreateTemp("", "test-download")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("Download Content")
	tmpFile.Seek(0, 0)
	fileInfo, _ := tmpFile.Stat()

	// 3. Mock storage
	env.StorageService.On("Open", "documents/file.txt").Return(tmpFile, fileInfo, "text/plain", nil)

	// Note: HoneyfileSvc is REAL - it will check the DB for honeyfiles
	// Since we haven't created any honeyfiles, this file won't trigger

	// 4. Setup Router
	router := testutils.SetupTestRouter(env)

	// 5. Generate REAL token
	token, err := env.GenerateTestToken("user-1", "test@example.com")
	assert.NoError(t, err)

	// 6. Request download
	req := httptest.NewRequest("GET", "/api/v1/storage/download?path=documents/file.txt", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 7. Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Download Content", w.Body.String())

	env.StorageService.AssertExpectations(t)
	tmpFile.Close()
}

func TestUpload_WithRealEncryptionPolicy(t *testing.T) {
	// 1. Setup TestEnv - EncryptionPolicyService is REAL
	env := testutils.NewTestEnv(t)

	// 2. Mock storage
	env.StorageService.On("Save", "/documents", mock.Anything, mock.Anything).Return(&content.SaveResult{
		Path:      "documents/large.bin",
		FileID:    "file-456",
		MimeType:  "application/octet-stream",
		SizeBytes: 1024,
	}, nil)

	// AI notification for unencrypted file
	env.AIService.On("NotifyUpload", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	// 3. Setup Router (uses REAL PolicyService)
	router := testutils.SetupTestRouter(env)

	// 4. Generate REAL token
	token, err := env.GenerateTestToken("user-1", "test@example.com")
	assert.NoError(t, err)

	// 5. Create upload request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "large.bin")
	part.Write([]byte("Binary content"))
	writer.WriteField("path", "/documents")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 6. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	env.StorageService.AssertExpectations(t)
}

// Helper to set mock expectations for encryption mode
func setupEncryptionPolicyMock(env *testutils.TestEnv, mode files.EncryptionMode) {
	// Note: PolicyService is now REAL, so this helper is deprecated
	// The real service will determine encryption mode based on file type/size
	_ = mode
}
