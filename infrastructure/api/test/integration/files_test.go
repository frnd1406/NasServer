package integration

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/domain/files"
	"github.com/nas-ai/api/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// File Upload Integration Tests
// Uses mock StorageManager via testutils
// ============================================================

func setupFileRouter(env *testutils.TestEnv) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Request ID middleware
	router.Use(func(c *gin.Context) {
		c.Set("request_id", "test-file-upload")
		c.Next()
	})

	// Mock upload endpoint that uses env mocks
	router.POST("/api/v1/storage/upload", func(c *gin.Context) {
		path := c.PostForm("path")
		encryptionOverride := c.PostForm("encryption_override")

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		// Call mock policy service
		encryptionMode := env.PolicyService.DetermineMode(
			fileHeader.Filename,
			fileHeader.Size,
			encryptionOverride,
		)

		// Open file
		src, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read upload"})
			return
		}
		defer src.Close()

		// Call mock storage manager
		result, err := env.StorageManager.Save(path, src, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Call mock AI service for unencrypted files
		if encryptionMode == string(files.EncryptionNone) {
			env.AIService.NotifyUpload(result.Path, result.FileID, result.MimeType, "")
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":            "ok",
			"encryption_status": encryptionMode,
			"size_bytes":        result.SizeBytes,
			"checksum":          result.Checksum,
		})
	})

	return router
}

func createMultipartRequest(path, filename string, content []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	io.Copy(part, bytes.NewReader(content))

	// Add path
	writer.WriteField("path", path)
	writer.WriteField("encryption_override", "")

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestUpload_Success(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := setupFileRouter(env)

	// Mock expectations
	env.PolicyService.On("DetermineMode", "test.txt", int64(13), "").
		Return(string(files.EncryptionNone))

	mockResult := &testutils.MockSaveResult{
		Path:      "/uploads/test.txt",
		FileID:    "file-123",
		MimeType:  "text/plain",
		SizeBytes: 13,
		Checksum:  "abc123",
	}
	env.StorageManager.On("Save", "/uploads", mock.Anything, mock.Anything).
		Return(mockResult, nil)

	env.AIService.On("NotifyUpload", "/uploads/test.txt", "file-123", "text/plain", "").
		Return()

	// Create request
	req, err := createMultipartRequest("/uploads", "test.txt", []byte("Hello, World!"))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
	assert.Contains(t, w.Body.String(), `"checksum":"abc123"`)

	// Verify mocks
	env.PolicyService.AssertExpectations(t)
	env.StorageManager.AssertExpectations(t)
	env.AIService.AssertExpectations(t)
}

func TestUpload_MissingFile(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := setupFileRouter(env)

	// Request without file
	req := httptest.NewRequest(http.MethodPost, "/api/v1/storage/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "file is required")
}

func TestUpload_StorageError(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := setupFileRouter(env)

	// Mock expectations - storage fails
	env.PolicyService.On("DetermineMode", "fail.txt", int64(4), "").
		Return(string(files.EncryptionNone))

	env.StorageManager.On("Save", "/uploads", mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	// Create request
	req, err := createMultipartRequest("/uploads", "fail.txt", []byte("fail"))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify mocks
	env.PolicyService.AssertExpectations(t)
	env.StorageManager.AssertExpectations(t)
}
