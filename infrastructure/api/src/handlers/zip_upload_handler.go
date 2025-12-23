package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

// Security constants for ZIP extraction
const (
	// ZipMagicBytes - First 4 bytes of a valid ZIP file
	ZipMagicBytes = "PK\x03\x04"

	// MaxDecompressedSize - Maximum total decompressed size (1 GB)
	// Prevents zip bombs from exhausting disk space
	MaxDecompressedSize int64 = 1024 * 1024 * 1024

	// MaxFileCount - Maximum number of files allowed in a ZIP
	// Prevents DoS attacks with millions of small files
	MaxFileCount = 10000

	// MaxSingleFileSize - Maximum size for a single extracted file (500MB)
	MaxSingleFileSize int64 = 500 * 1024 * 1024

	// MaxCompressionRatio - Maximum allowed ratio of uncompressed/compressed size
	// Additional protection against sophisticated zip bombs
	MaxCompressionRatio = 100
)

// UnzipResult contains information about the extraction result
type UnzipResult struct {
	ExtractedFiles []string
	TotalBytes     int64
	FileCount      int
}

// UnzipSecure extracts a ZIP archive securely with multiple protection layers:
// 1. Path Traversal Prevention (Zip-Slip protection)
// 2. Zip Bomb Protection (size limit with real-time byte counting)
// 3. File Count Limit (prevents DoS with many small files)
//
// Returns an error immediately if any security check fails.
func UnzipSecure(zipData []byte, destDir string, logger *logrus.Logger) (*UnzipResult, error) {
	// Validate destination directory exists and is absolute
	destDir = filepath.Clean(destDir)
	if !filepath.IsAbs(destDir) {
		return nil, fmt.Errorf("destination must be an absolute path")
	}

	// Create ZIP reader
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("invalid or corrupted archive: %w", err)
	}

	// SECURITY CHECK 1: File Count Limit
	if len(zipReader.File) > MaxFileCount {
		return nil, fmt.Errorf("archive contains too many files (%d > %d)", len(zipReader.File), MaxFileCount)
	}

	// Pre-check: Validate compression ratio before extraction
	compressedSize := int64(len(zipData))
	var declaredUncompressed uint64
	for _, f := range zipReader.File {
		declaredUncompressed += f.UncompressedSize64
	}
	if compressedSize > 0 {
		ratio := float64(declaredUncompressed) / float64(compressedSize)
		if ratio > MaxCompressionRatio {
			return nil, fmt.Errorf("compression ratio too high (%.1f > %d) - possible zip bomb", ratio, MaxCompressionRatio)
		}
	}

	// Pre-check: Declared size check
	if int64(declaredUncompressed) > MaxDecompressedSize {
		return nil, fmt.Errorf("declared size exceeds limit (%d > %d bytes)", declaredUncompressed, MaxDecompressedSize)
	}

	result := &UnzipResult{
		ExtractedFiles: make([]string, 0, len(zipReader.File)),
	}

	// Track total extracted bytes in real-time
	var totalExtractedBytes int64

	for _, f := range zipReader.File {
		// SECURITY CHECK 2: Path Traversal Prevention (Zip-Slip)
		// This is the most critical security check
		fpath := filepath.Join(destDir, f.Name)

		// Normalize both paths for comparison
		cleanDest := filepath.Clean(destDir) + string(os.PathSeparator)
		cleanFpath := filepath.Clean(fpath)

		// The extracted file path MUST start with the destination directory
		if !strings.HasPrefix(cleanFpath, cleanDest) && cleanFpath != filepath.Clean(destDir) {
			return nil, fmt.Errorf("🚨 SECURITY: illegal file path detected: %s (attempted path traversal)", f.Name)
		}

		// Handle directories
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0755); err != nil {
				if logger != nil {
					logger.WithError(err).Warnf("Failed to create directory: %s", fpath)
				}
			}
			continue
		}

		// Skip files that are too large individually
		if int64(f.UncompressedSize64) > MaxSingleFileSize {
			if logger != nil {
				logger.Warnf("Skipping oversized file: %s (%d bytes > %d limit)", f.Name, f.UncompressedSize64, MaxSingleFileSize)
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for %s: %w", f.Name, err)
		}

		// Extract file with byte counting
		bytesWritten, err := extractFileSecure(f, fpath, MaxDecompressedSize-totalExtractedBytes)
		if err != nil {
			// Check if it's a size limit error
			if strings.Contains(err.Error(), "size limit exceeded") {
				return nil, fmt.Errorf("🚨 SECURITY: zip bomb detected - extraction aborted at %d bytes", totalExtractedBytes+bytesWritten)
			}
			return nil, fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}

		totalExtractedBytes += bytesWritten

		// SECURITY CHECK 3: Real-time Zip Bomb Protection
		// Check after each file extraction
		if totalExtractedBytes > MaxDecompressedSize {
			return nil, fmt.Errorf("🚨 SECURITY: total decompressed size exceeded (%d > %d bytes) - possible zip bomb", totalExtractedBytes, MaxDecompressedSize)
		}

		result.ExtractedFiles = append(result.ExtractedFiles, f.Name)
		result.FileCount++
	}

	result.TotalBytes = totalExtractedBytes
	return result, nil
}

// extractFileSecure extracts a single file with byte counting and size limit enforcement
// Returns the number of bytes written and any error
func extractFileSecure(f *zip.File, destPath string, remainingBudget int64) (int64, error) {
	rc, err := f.Open()
	if err != nil {
		return 0, fmt.Errorf("open zip entry: %w", err)
	}
	defer rc.Close()

	// Force safe permissions (never trust archive permissions)
	mode := os.FileMode(0644)

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}
	defer outFile.Close()

	// Use a counting writer to track bytes in real-time
	// Also enforce the remaining budget
	if remainingBudget <= 0 {
		return 0, fmt.Errorf("size limit exceeded: no budget remaining")
	}

	// Limit the reader to the remaining budget
	limitedReader := io.LimitReader(rc, remainingBudget+1) // +1 to detect overflow

	bytesWritten, err := io.Copy(outFile, limitedReader)
	if err != nil {
		return bytesWritten, fmt.Errorf("write file: %w", err)
	}

	// Check if we hit the limit
	if bytesWritten > remainingBudget {
		os.Remove(destPath) // Clean up partial file
		return bytesWritten, fmt.Errorf("size limit exceeded during extraction")
	}

	return bytesWritten, nil
}

// StorageUploadZipHandler handles secure ZIP file uploads with extraction
// @Summary Upload and extract ZIP archive
// @Description Securely upload and extract a ZIP file with anti-bomb and anti-slip protection
// @Tags storage
// @Security BearerAuth
// @Security CSRFToken
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "ZIP file to upload"
// @Param path formData string false "Target directory path"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/storage/upload-zip [post]
func StorageUploadZipHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		targetPath := c.PostForm("path")
		if targetPath == "" {
			targetPath = "/"
		}

		// Get uploaded file
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		// Open the uploaded file
		src, err := fileHeader.Open()
		if err != nil {
			logger.WithError(err).Error("Failed to open uploaded file")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read upload"})
			return
		}
		defer src.Close()

		// SECURITY CHECK: Verify Magic Bytes (Trust No Extension)
		magicBuf := make([]byte, 4)
		if _, err := io.ReadFull(src, magicBuf); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file too small or unreadable"})
			return
		}

		if !bytes.HasPrefix(magicBuf, []byte(ZipMagicBytes)) {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"magic":      fmt.Sprintf("%x", magicBuf),
			}).Warn("Invalid ZIP magic bytes - possible attack")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid archive format"})
			return
		}

		// Reset reader to beginning
		if _, err := src.Seek(0, io.SeekStart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process file"})
			return
		}

		// Read entire file into memory for zip processing
		// Limit read to prevent memory exhaustion
		zipData, err := io.ReadAll(io.LimitReader(src, MaxDecompressedSize))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read archive"})
			return
		}

		// Get target directory
		targetDir, err := storage.GetFullPath(targetPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target path"})
			return
		}

		// Use the secure unzip function
		result, err := UnzipSecure(zipData, targetDir, logger)
		if err != nil {
			// Check for security violations
			if strings.Contains(err.Error(), "SECURITY") {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"error":      err.Error(),
				}).Error("Security violation during ZIP extraction")
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "security violation",
					"details": err.Error(),
				})
				return
			}

			// Other errors
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("ZIP extraction failed")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "extraction failed",
				"details": err.Error(),
			})
			return
		}

		logger.WithFields(logrus.Fields{
			"request_id":      requestID,
			"files_extracted": result.FileCount,
			"bytes_extracted": result.TotalBytes,
			"target_path":     targetPath,
		}).Info("ZIP archive extracted successfully")

		c.JSON(http.StatusOK, gin.H{
			"status":          "extracted",
			"files_extracted": result.FileCount,
			"bytes_extracted": result.TotalBytes,
			"target_path":     targetPath,
		})
	}
}
