package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmartDownloadHandler_UnencryptedFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create temp storage directory
	tmpDir, err := os.MkdirTemp("", "smart_download_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create storage service
	storage, err := services.NewStorageService(tmpDir, logger)
	require.NoError(t, err)

	// Create test file
	testContent := []byte("Hello, this is test content for download!")
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)

	// Setup router
	router := gin.New()
	router.GET("/download", SmartDownloadHandler(storage, nil, logger))

	// Test download
	req := httptest.NewRequest("GET", "/download?path=test.txt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "test.txt")
	assert.Equal(t, testContent, w.Body.Bytes())
}

func TestSmartDownloadHandler_EncryptedFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create temp storage directory
	tmpDir, err := os.MkdirTemp("", "smart_download_enc_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create storage service
	storage, err := services.NewStorageService(tmpDir, logger)
	require.NoError(t, err)

	// Create and encrypt test file
	password := "testpassword123"
	testContent := []byte("Secret encrypted content that must be decrypted!")

	// Create encrypted file
	encryptedPath := filepath.Join(tmpDir, "secret.txt.enc")
	encFile, err := os.Create(encryptedPath)
	require.NoError(t, err)

	err = services.EncryptStream(password, bytes.NewReader(testContent), encFile)
	encFile.Close()
	require.NoError(t, err)

	// Setup router
	router := gin.New()
	router.GET("/download", SmartDownloadHandler(storage, nil, logger))

	// Test download without password (should fail)
	req := httptest.NewRequest("GET", "/download?path=secret.txt.enc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "password required")

	// Test download with wrong password
	req = httptest.NewRequest("GET", "/download?path=secret.txt.enc&password=wrongpassword", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Will start streaming but has incomplete/corrupted data

	// Test download with correct password
	req = httptest.NewRequest("GET", "/download?path=secret.txt.enc&password="+password, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "secret.txt")
	assert.Equal(t, testContent, w.Body.Bytes())
}

func TestSmartDownloadHandler_RangeRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create temp storage directory
	tmpDir, err := os.MkdirTemp("", "smart_download_range_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create storage service
	storage, err := services.NewStorageService(tmpDir, logger)
	require.NoError(t, err)

	// Create test file with known content
	testContent := make([]byte, 1000)
	for i := range testContent {
		testContent[i] = byte(i % 256)
	}
	testFile := filepath.Join(tmpDir, "range_test.bin")
	err = os.WriteFile(testFile, testContent, 0644)
	require.NoError(t, err)

	// Setup router
	router := gin.New()
	router.GET("/download", SmartDownloadHandler(storage, nil, logger))

	// Test Range request
	req := httptest.NewRequest("GET", "/download?path=range_test.bin", nil)
	req.Header.Set("Range", "bytes=0-99")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPartialContent, w.Code)
	assert.Equal(t, 100, len(w.Body.Bytes()))
	assert.Equal(t, testContent[:100], w.Body.Bytes())
}

func TestParseRangeHeader(t *testing.T) {
	tests := []struct {
		name        string
		rangeHeader string
		fileSize    int64
		wantStart   int64
		wantEnd     int64
		wantErr     bool
	}{
		{
			name:        "Simple range",
			rangeHeader: "bytes=0-99",
			fileSize:    1000,
			wantStart:   0,
			wantEnd:     99,
			wantErr:     false,
		},
		{
			name:        "Open-ended range",
			rangeHeader: "bytes=500-",
			fileSize:    1000,
			wantStart:   500,
			wantEnd:     999,
			wantErr:     false,
		},
		{
			name:        "Suffix range",
			rangeHeader: "bytes=-100",
			fileSize:    1000,
			wantStart:   900,
			wantEnd:     999,
			wantErr:     false,
		},
		{
			name:        "Invalid format",
			rangeHeader: "invalid",
			fileSize:    1000,
			wantErr:     true,
		},
		{
			name:        "Range exceeds file",
			rangeHeader: "bytes=0-1500",
			fileSize:    1000,
			wantStart:   0,
			wantEnd:     999, // Clamped to file size
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := parseRangeHeader(tt.rangeHeader, tt.fileSize)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStart, start)
				assert.Equal(t, tt.wantEnd, end)
			}
		})
	}
}

func TestDetectContentType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"video.mp4", "video/mp4"},
		{"audio.mp3", "audio/mpeg"},
		{"document.pdf", "application/pdf"},
		{"image.jpg", "image/jpeg"},
		{"image.png", "image/png"},
		{"text.txt", "text/plain; charset=utf-8"},
		{"data.json", "application/json; charset=utf-8"},
		{"page.html", "text/html; charset=utf-8"},
		{"unknown.xyz", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := detectContentType("/path/"+tt.filename, tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSkipLimitWriter(t *testing.T) {
	// Test skipping and limiting
	buf := new(bytes.Buffer)
	slw := &skipLimitWriter{
		writer:     buf,
		skipBytes:  5,
		limitBytes: 10,
	}

	// Write data in chunks to simulate streaming
	data := []byte("0123456789ABCDEFGHIJ") // 20 bytes

	// First write: "01234" - all skipped
	slw.Write(data[:5])
	assert.Equal(t, "", buf.String(), "First 5 bytes should be skipped")

	// Second write: "56789" - all written
	slw.Write(data[5:10])
	assert.Equal(t, "56789", buf.String(), "Next 5 bytes should be written")

	// Third write: "ABCDE" - all written (now at limit)
	slw.Write(data[10:15])
	assert.Equal(t, "56789ABCDE", buf.String(), "Should have written up to limit")

	// Fourth write should be blocked
	n, err := slw.Write(data[15:])
	assert.Equal(t, 0, n) // Nothing written
	assert.Equal(t, errLimitReached, err)
}

func TestDetectEncryptionStatus(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "detect_enc_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test 1: Regular file
	regularFile := filepath.Join(tmpDir, "regular.txt")
	os.WriteFile(regularFile, []byte("hello"), 0644)

	status := detectEncryptionStatus(regularFile, logger)
	assert.Equal(t, "NONE", string(status))

	// Test 2: File with .enc extension but no magic bytes
	fakeEncFile := filepath.Join(tmpDir, "fake.enc")
	os.WriteFile(fakeEncFile, []byte("not encrypted"), 0644)

	status = detectEncryptionStatus(fakeEncFile, logger)
	assert.Equal(t, "NONE", string(status)) // Should detect it's not actually encrypted

	// Test 3: Actually encrypted file
	realEncFile := filepath.Join(tmpDir, "real.txt.enc")
	encFile, _ := os.Create(realEncFile)
	services.EncryptStream("password", bytes.NewReader([]byte("secret")), encFile)
	encFile.Close()

	status = detectEncryptionStatus(realEncFile, logger)
	assert.Equal(t, "USER", string(status))
}
