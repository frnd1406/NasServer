package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

// ==============================================================================
// Smart Download Handler - Phase 4: Hybrid Streaming
// ==============================================================================
//
// This handler implements intelligent file delivery based on encryption status:
//
// NONE (Performance Path):
//   - Uses Nginx X-Accel-Redirect for zero-copy file serving
//   - Go doesn't read the file - Nginx handles it directly
//   - Supports Range requests natively via Nginx
//
// USER (Secure Path):
//   - Streams file through DecryptStream to HTTP response
//   - Supports partial Range requests for video seeking
//   - Uses chunked AEAD decryption (64KB chunks)
//
// ==============================================================================

// DownloadRequest represents parameters for file download
type DownloadRequest struct {
	Path               string // Relative path from storage root
	EncryptionPassword string // Required for USER-encrypted files
	Inline             bool   // If true, display inline instead of download
}

// SmartDownloadHandler handles file downloads with encryption-aware streaming.
// Automatically detects encrypted files and routes to appropriate handler.
//
// Query Parameters:
//   - path: Required. Relative path to file
//   - password: Required for encrypted files. Decryption password
//   - inline: Optional. If "true", use inline Content-Disposition
//
// Headers:
//   - Range: Optional. Byte range for partial content (video seeking)
func SmartDownloadHandler(
	storage *services.StorageService,
	honeySvc *services.HoneyfileService,
	logger *logrus.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.Query("path")

		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}

		// Parse optional parameters
		password := c.Query("password")
		inline := c.Query("inline") == "true"

		// Get full filesystem path
		fullPath, err := storage.GetFullPath(path)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		// ==== SECURITY: Honeyfile Check ====
		if honeySvc != nil {
			meta := services.RequestMetadata{
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Action:    "download",
			}
			if honeySvc.CheckAndTrigger(c.Request.Context(), fullPath, meta) {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"path":       path,
					"ip":         meta.IPAddress,
				}).Error("🔒 HONEYFILE TRIGGERED - ACCESS DENIED")
				c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
				return
			}
		}

		// Check if file exists
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
				return
			}
			handleStorageError(c, err, logger, requestID)
			return
		}

		if fileInfo.IsDir() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot download a directory"})
			return
		}

		// ==== DETECT ENCRYPTION STATUS ====
		encryptionStatus := detectEncryptionStatus(fullPath, logger)

		logger.WithFields(logrus.Fields{
			"request_id":        requestID,
			"path":              path,
			"encryption_status": encryptionStatus,
			"size":              fileInfo.Size(),
		}).Debug("Download request")

		// ==== ROUTE BASED ON ENCRYPTION STATUS ====
		switch encryptionStatus {
		case models.EncryptionNone:
			// Performance Path: X-Accel-Redirect or direct serve
			serveUnencryptedFile(c, fullPath, fileInfo, inline, logger, requestID)

		case models.EncryptionUser:
			// Secure Path: Streaming decryption
			if password == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":             "password required for encrypted file",
					"encryption_status": "USER",
				})
				return
			}
			serveEncryptedFile(c, fullPath, fileInfo, password, inline, logger, requestID)

		case models.EncryptionSystem:
			// Future: System encryption not yet implemented
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "SYSTEM encryption not yet supported for download",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown encryption status"})
		}
	}
}

// detectEncryptionStatus determines how a file is encrypted based on:
// 1. File extension (.enc suffix)
// 2. Magic bytes (NASC header)
func detectEncryptionStatus(fullPath string, logger *logrus.Logger) models.EncryptionMode {
	// Check file extension first (fast path)
	if strings.HasSuffix(strings.ToLower(fullPath), ".enc") {
		// Verify with magic bytes
		isEnc, err := services.IsEncryptedFile(fullPath)
		if err != nil {
			logger.WithError(err).Warn("Failed to check encryption magic bytes")
			return models.EncryptionNone // Fail open for availability
		}
		if isEnc {
			return models.EncryptionUser
		}
	}

	// No .enc extension = unencrypted
	return models.EncryptionNone
}

// serveUnencryptedFile serves a file without encryption using optimal method
func serveUnencryptedFile(
	c *gin.Context,
	fullPath string,
	fileInfo os.FileInfo,
	inline bool,
	logger *logrus.Logger,
	requestID string,
) {
	filename := fileInfo.Name()
	contentType := detectContentType(fullPath, filename)

	// Set Content-Disposition
	disposition := "attachment"
	if inline {
		disposition = "inline"
	}
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, filename))
	c.Header("Content-Type", contentType)

	// ==== X-ACCEL-REDIRECT (Nginx Zero-Copy) ====
	// If running behind Nginx with X-Accel-Redirect configured, use it
	// This is the fastest possible serving method
	//
	// Nginx config required:
	// location /protected-files/ {
	//     internal;
	//     alias /mnt/data/;
	// }

	// Check if we should use X-Accel-Redirect (configurable via env)
	useXAccel := os.Getenv("USE_NGINX_XACCEL") == "true"

	if useXAccel {
		// Convert full path to X-Accel-Redirect path
		// /mnt/data/folder/file.txt -> /protected-files/folder/file.txt
		xAccelPath := strings.Replace(fullPath, "/mnt/data", "/protected-files", 1)

		c.Header("X-Accel-Redirect", xAccelPath)
		c.Header("X-Accel-Buffering", "no") // Disable buffering for streaming

		logger.WithFields(logrus.Fields{
			"request_id":   requestID,
			"x_accel_path": xAccelPath,
		}).Debug("Using X-Accel-Redirect for file serving")

		c.Status(http.StatusOK)
		return
	}

	// ==== FALLBACK: Direct File Serve ====
	// If not using Nginx, serve directly with Range support
	file, err := os.Open(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	// Use http.ServeContent for automatic Range request handling
	http.ServeContent(c.Writer, c.Request, filename, fileInfo.ModTime(), file)
}

// serveEncryptedFile streams a decrypted file to the client
func serveEncryptedFile(
	c *gin.Context,
	fullPath string,
	fileInfo os.FileInfo,
	password string,
	inline bool,
	logger *logrus.Logger,
	requestID string,
) {
	// Open encrypted file
	file, err := os.Open(fullPath)
	if err != nil {
		logger.WithError(err).Error("Failed to open encrypted file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	// Get original filename (strip .enc suffix)
	filename := fileInfo.Name()
	if strings.HasSuffix(strings.ToLower(filename), ".enc") {
		filename = filename[:len(filename)-4]
	}

	contentType := detectContentType(fullPath, filename)

	// Check for Range header (video seeking)
	rangeHeader := c.Request.Header.Get("Range")
	if rangeHeader != "" {
		// Handle partial content request for encrypted files
		serveEncryptedRange(c, file, fileInfo, password, filename, contentType, rangeHeader, inline, logger, requestID)
		return
	}

	// ==== FULL FILE STREAMING DECRYPTION ====

	// Set Content-Disposition
	disposition := "attachment"
	if inline {
		disposition = "inline"
	}
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, filename))
	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes") // Signal range support

	// Note: We can't set Content-Length for encrypted files because
	// the decrypted size is different from encrypted size
	// The chunked transfer encoding will be used instead

	c.Status(http.StatusOK)

	// Stream decrypted content directly to response
	err = services.DecryptStream(password, file, c.Writer)
	if err != nil {
		// Check for authentication failure (wrong password)
		if errors.Is(err, services.ErrCorruptedData) || errors.Is(err, services.ErrInvalidHeader) {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("Decryption failed - possible wrong password")
			// Can't change status after writing started, just log
			return
		}

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"error":      err.Error(),
		}).Error("Decryption stream failed")
		return
	}

	logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"filename":   filename,
	}).Info("Encrypted file served successfully")
}

// serveEncryptedRange handles Range requests for encrypted files
// This enables video seeking in encrypted files using chunk-level seeking
func serveEncryptedRange(
	c *gin.Context,
	file *os.File,
	fileInfo os.FileInfo,
	password string,
	filename string,
	contentType string,
	rangeHeader string,
	inline bool,
	logger *logrus.Logger,
	requestID string,
) {
	// Get encrypted file info to calculate plaintext size
	encInfo, err := services.GetEncryptedFileInfo(file)
	if err != nil {
		logger.WithError(err).Error("Failed to get encrypted file info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file metadata"})
		return
	}

	if !encInfo.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is not properly encrypted"})
		return
	}

	// Use estimated plaintext size for range validation
	plaintextSize := encInfo.EstimatedPlainSize

	// Parse Range header using plaintext size
	rangeStart, rangeEnd, err := parseRangeHeader(rangeHeader, plaintextSize)
	if err != nil {
		c.Header("Content-Range", fmt.Sprintf("bytes */%d", plaintextSize))
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"error": "invalid range"})
		return
	}

	// Calculate content length for this range
	contentLength := rangeEnd - rangeStart + 1

	logger.WithFields(logrus.Fields{
		"request_id":     requestID,
		"range_start":    rangeStart,
		"range_end":      rangeEnd,
		"content_length": contentLength,
		"plaintext_size": plaintextSize,
	}).Debug("Processing encrypted range request with chunk-seeking")

	// Set headers for partial content
	disposition := "attachment"
	if inline {
		disposition = "inline"
	}
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, filename))
	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, plaintextSize))
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))

	c.Status(http.StatusPartialContent)

	// ==== USE OPTIMIZED CHUNK-SEEKING DECRYPTION ====
	// DecryptStreamWithSeek:
	// 1. Reads header from position 0 (always required)
	// 2. Calculates target chunk from rangeStart
	// 3. Seeks directly to that chunk's encrypted position
	// 4. Decrypts and outputs only the requested bytes

	bytesWritten, err := services.DecryptStreamWithSeek(password, file, c.Writer, rangeStart, contentLength)
	if err != nil {
		// Check for authentication failure (wrong password)
		if errors.Is(err, services.ErrCorruptedData) || errors.Is(err, services.ErrInvalidHeader) {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("Range decryption failed - possible wrong password")
			return
		}

		logger.WithFields(logrus.Fields{
			"request_id":    requestID,
			"error":         err.Error(),
			"bytes_written": bytesWritten,
		}).Error("Range decryption stream failed")
		return
	}

	logger.WithFields(logrus.Fields{
		"request_id":    requestID,
		"range_start":   rangeStart,
		"range_end":     rangeEnd,
		"bytes_written": bytesWritten,
	}).Debug("Encrypted range served successfully")
}

// skipLimitWriter skips initial bytes and limits output
type skipLimitWriter struct {
	writer     io.Writer
	skipBytes  int64
	limitBytes int64
	written    int64
	skipped    int64
}

var errLimitReached = errors.New("limit reached")

func (w *skipLimitWriter) Write(p []byte) (n int, err error) {
	if w.written >= w.limitBytes {
		return 0, errLimitReached
	}

	// Skip initial bytes
	if w.skipped < w.skipBytes {
		toSkip := w.skipBytes - w.skipped
		if int64(len(p)) <= toSkip {
			w.skipped += int64(len(p))
			return len(p), nil
		}
		// Skip portion, write rest
		p = p[toSkip:]
		w.skipped = w.skipBytes
	}

	// Limit output
	remaining := w.limitBytes - w.written
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err = w.writer.Write(p)
	w.written += int64(n)

	if w.written >= w.limitBytes {
		return n, errLimitReached
	}
	return n, err
}

// parseRangeHeader parses the Range header and returns start/end bytes
func parseRangeHeader(rangeHeader string, fileSize int64) (start, end int64, err error) {
	// Format: "bytes=0-1023" or "bytes=1024-" or "bytes=-500"
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, errors.New("invalid range format")
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid range format")
	}

	if parts[0] == "" {
		// Suffix range: "-500" means last 500 bytes
		suffixLen, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		start = fileSize - suffixLen
		end = fileSize - 1
	} else if parts[1] == "" {
		// Open-ended range: "1024-" means from 1024 to end
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		end = fileSize - 1
	} else {
		// Explicit range: "0-1023"
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	// Validate range
	if start < 0 || start > end || start >= fileSize {
		return 0, 0, errors.New("range not satisfiable")
	}

	// Clamp end to file size
	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}

// detectContentType determines the MIME type for a file
func detectContentType(fullPath, filename string) string {
	// Try extension-based detection first
	ext := strings.ToLower(filepath.Ext(filename))

	// Common types
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".txt":
		return "text/plain; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".md":
		return "text/markdown; charset=utf-8"
	case ".zip":
		return "application/zip"
	case ".tar":
		return "application/x-tar"
	case ".gz":
		return "application/gzip"
	}

	// Fallback to mime package
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType
	}

	return "application/octet-stream"
}
