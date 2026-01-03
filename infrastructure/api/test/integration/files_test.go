package integration

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nas-ai/api/src/domain/files"
	"github.com/nas-ai/api/src/services/content"
	"github.com/nas-ai/api/src/services/security"
	"github.com/nas-ai/api/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestFileUpload_Success(t *testing.T) {
	// 1. Setup Mocks
	mockStorage := new(testutils.MockStorageService)
	mockPolicy := new(testutils.MockEncryptionPolicyService)
	mockAI := new(testutils.MockAIAgentService)
	mockHoney := new(testutils.MockHoneyfileService)
	mockJWT := new(testutils.MockJWTService)
	mockToken := new(testutils.MockTokenService)
	mockEncryption := new(testutils.MockEncryptionService)
	mockEncryptedStorage := new(testutils.MockEncryptedStorageService)
	env := testutils.NewTestEnv(t)

	// 2. Setup Expectations
	// Auth Mocks
	mockJWT.On("ValidateToken", "valid-token").Return(&security.TokenClaims{
		UserID: "user-1",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)
	mockToken.On("IsTokenRevoked", mock.Anything, "user-1", mock.Anything).Return(false)

	// Policy determines no encryption
	mockPolicy.On("DetermineMode", "test.txt", int64(12), "").Return(files.EncryptionNone)

	expectedResult := &content.SaveResult{
		Path:      "uploads/test.txt",
		FileID:    "file-123",
		MimeType:  "text/plain",
		SizeBytes: 12,
		Checksum:  "checksum",
	}
	// Storage Save called
	mockStorage.On("Save", "/uploads", mock.Anything, mock.Anything).Return(expectedResult, nil)

	// AI Notified (since unencrypted)
	mockAI.On("NotifyUpload", "uploads/test.txt", "file-123", "text/plain", "Hello World!").Return()

	// 3. Setup Router
	router := testutils.SetupTestRouter(
		env,
		new(testutils.MockUserRepository),
		mockJWT,
		new(testutils.MockPasswordService),
		mockToken,
		mockStorage,
		mockHoney,
		mockAI,
		mockPolicy,
		mockEncryption,
		mockEncryptedStorage,
	)

	// 4. Create Request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("Hello World!"))
	writer.WriteField("path", "/uploads")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/storage/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	mockStorage.AssertExpectations(t)
	mockAI.AssertExpectations(t)
	mockPolicy.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestFileDownload_Success(t *testing.T) {
	// 1. Setup Mocks
	mockStorage := new(testutils.MockStorageService)
	mockHoney := new(testutils.MockHoneyfileService)
	mockJWT := new(testutils.MockJWTService)
	mockToken := new(testutils.MockTokenService)
	mockEncryption := new(testutils.MockEncryptionService)
	mockEncryptedStorage := new(testutils.MockEncryptedStorageService)
	env := testutils.NewTestEnv(t)

	// 2. Prepare Temp File for "Open" to return valid *os.File
	tmpFile, err := os.CreateTemp("", "test-download")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("File Content")
	tmpFile.Seek(0, 0)

	fileInfo, _ := tmpFile.Stat()

	// 3. Expectations
	// Auth Mocks
	mockJWT.On("ValidateToken", "valid-token").Return(&security.TokenClaims{
		UserID: "user-1",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)
	mockToken.On("IsTokenRevoked", mock.Anything, "user-1", mock.Anything).Return(false)

	// Honeyfile check returns false (not a honeyfile)
	mockHoney.On("CheckAndTrigger", mock.Anything, "/mnt/data/test/path.txt", mock.Anything).Return(false)

	// Storage Open returns the temp file
	mockStorage.On("Open", "test/path.txt").Return(tmpFile, fileInfo, "text/plain", nil)

	// 4. Router
	router := testutils.SetupTestRouter(
		env,
		new(testutils.MockUserRepository),
		mockJWT,
		new(testutils.MockPasswordService),
		mockToken,
		mockStorage,
		mockHoney,
		new(testutils.MockAIAgentService),
		new(testutils.MockEncryptionPolicyService),
		mockEncryption,
		mockEncryptedStorage,
	)

	// 5. Request
	req := httptest.NewRequest("GET", "/api/v1/files/content?path=test/path.txt", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 6. Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "File Content", w.Body.String())
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")

	mockStorage.AssertExpectations(t)
	mockHoney.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockToken.AssertExpectations(t)

	// Close temp file handle if handler didn't (handler defers Close, but it's good practice)
	tmpFile.Close()
}

func TestEncryptedFileUpload_Success(t *testing.T) {
	// 1. Mocks
	mockEncryptedStorage := new(testutils.MockEncryptedStorageService)
	mockEncryption := new(testutils.MockEncryptionService)
	mockJWT := new(testutils.MockJWTService)
	mockToken := new(testutils.MockTokenService)
	mockStorage := new(testutils.MockStorageService)
	mockHoney := new(testutils.MockHoneyfileService)
	env := testutils.NewTestEnv(t)

	// 2. Expectations
	// Auth Mocks
	mockJWT.On("ValidateToken", "valid-token").Return(&security.TokenClaims{
		UserID: "user-1",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, nil)
	mockToken.On("IsTokenRevoked", mock.Anything, "user-1", mock.Anything).Return(false)

	// Encryption Enabled Check
	mockEncryptedStorage.On("IsEncryptionEnabled").Return(true)

	// Save Encrypted
	expectedResult := &content.SaveResult{
		Path:      "secret/file.enc",
		FileID:    "enc-file-1",
		MimeType:  "text/plain",
		SizeBytes: 20,
	}
	mockEncryptedStorage.On("SaveEncrypted", "/secret", mock.Anything, mock.Anything).Return(expectedResult, nil)

	// 3. Router
	router := testutils.SetupTestRouter(
		env,
		new(testutils.MockUserRepository),
		mockJWT,
		new(testutils.MockPasswordService),
		mockToken,
		mockStorage,
		mockHoney,
		new(testutils.MockAIAgentService),
		new(testutils.MockEncryptionPolicyService),
		mockEncryption,
		mockEncryptedStorage,
	)

	// 4. Request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "secret.txt")
	part.Write([]byte("Secret Info"))
	writer.WriteField("path", "/secret")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/encrypted/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	mockEncryptedStorage.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}
