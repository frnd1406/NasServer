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

const (
	// ZipMagicBytes - First 4 bytes of a valid ZIP file
	ZipMagicBytes = "PK\x03\x04"

	// MaxCompressionRatio - Maximum allowed ratio of uncompressed/compressed size
	// Prevents zip bombs (e.g., 42.zip where 42KB expands to petabytes)
	MaxCompressionRatio = 100

	// MaxUnpackedSize - Hard cap on total extracted content (1GB)
	MaxUnpackedSize = 1 << 30

	// MaxSingleFileSize - Maximum size for a single extracted file (500MB)
	MaxSingleFileSize = 500 << 20
)

// StorageUploadZipHandler handles secure ZIP file uploads with extraction
// @Summary Upload and extract ZIP archive
// @Description Securely upload and extract a ZIP file with anti-bomb and anti-slip protection
// @Tags storage
// @Security BearerAuth
// @Security CSRFToken
// @Accept multipart/form-data
// @Produce json
// @Param file formance file true "ZIP file to upload"
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

		// SECURITY CHECK 1: Verify Magic Bytes (Trust No Extension)
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
		// (Required by archive/zip - streaming not supported)
		zipData, err := io.ReadAll(io.LimitReader(src, MaxUnpackedSize))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read archive"})
			return
		}
		compressedSize := int64(len(zipData))

		// Open ZIP reader
		zipReader, err := zip.NewReader(bytes.NewReader(zipData), compressedSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or corrupted archive"})
			return
		}

		// Get target directory
		targetDir, err := storage.GetFullPath(targetPath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target path"})
			return
		}

		// Pre-calculate total uncompressed size for ratio check
		var totalUncompressed uint64
		for _, f := range zipReader.File {
			totalUncompressed += f.UncompressedSize64
		}

		// SECURITY CHECK 2: Anti-Zip-Bomb (Ratio Limit)
		if compressedSize > 0 {
			ratio := float64(totalUncompressed) / float64(compressedSize)
			if ratio > MaxCompressionRatio {
				logger.WithFields(logrus.Fields{
					"request_id":   requestID,
					"compressed":   compressedSize,
					"uncompressed": totalUncompressed,
					"ratio":        ratio,
				}).Error("🚨 ZIP BOMB DETECTED - Compression ratio exceeded")
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "compression ratio exceeded",
					"details": "Archive appears to be a zip bomb",
				})
				return
			}
		}

		// SECURITY CHECK 3: Hard Cap on Total Size
		if totalUncompressed > MaxUnpackedSize {
			logger.WithFields(logrus.Fields{
				"request_id":   requestID,
				"uncompressed": totalUncompressed,
				"limit":        MaxUnpackedSize,
			}).Warn("Archive exceeds maximum uncompressed size")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "archive too large",
				"details": fmt.Sprintf("Maximum uncompressed size is %d bytes", MaxUnpackedSize),
			})
			return
		}

		// Extract files
		var extractedFiles []string
		var extractedBytes int64

		for _, f := range zipReader.File {
			// SECURITY CHECK 4: Anti-Zip-Slip (Path Traversal)
			destPath, err := safeExtractPath(targetDir, f.Name)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"entry":      f.Name,
				}).Error("🚨 ZIP SLIP DETECTED - Path traversal attempt")
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "security violation",
					"details": "path traversal detected",
				})
				return
			}

			// Handle directories
			if f.FileInfo().IsDir() {
				if err := os.MkdirAll(destPath, 0755); err != nil {
					logger.WithError(err).Warn("Failed to create directory")
				}
				continue
			}

			// Single file size check
			if f.UncompressedSize64 > MaxSingleFileSize {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"entry":      f.Name,
					"size":       f.UncompressedSize64,
				}).Warn("Single file exceeds size limit")
				continue // Skip this file
			}

			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				logger.WithError(err).Warn("Failed to create parent directory")
				continue
			}

			// Extract file with safe permissions
			if err := extractFile(f, destPath); err != nil {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"entry":      f.Name,
					"error":      err.Error(),
				}).Warn("Failed to extract file")
				continue
			}

			extractedFiles = append(extractedFiles, f.Name)
			extractedBytes += int64(f.UncompressedSize64)
		}

		logger.WithFields(logrus.Fields{
			"request_id":      requestID,
			"files_extracted": len(extractedFiles),
			"bytes_extracted": extractedBytes,
			"target_path":     targetPath,
		}).Info("ZIP archive extracted successfully")

		c.JSON(http.StatusOK, gin.H{
			"status":          "extracted",
			"files_extracted": len(extractedFiles),
			"bytes_extracted": extractedBytes,
			"target_path":     targetPath,
		})
	}
}

// safeExtractPath validates and returns a safe extraction path
// Prevents Zip-Slip attacks by ensuring path stays within target directory
func safeExtractPath(targetDir, entryName string) (string, error) {
	// Clean the entry name to remove any ../ sequences
	cleanEntry := filepath.Clean(entryName)

	// Construct destination path
	destPath := filepath.Join(targetDir, cleanEntry)

	// Ensure the destination is within the target directory
	cleanTarget := filepath.Clean(targetDir) + string(os.PathSeparator)
	cleanDest := filepath.Clean(destPath)

	// Add separator to cleanDest for proper prefix matching
	if !strings.HasPrefix(cleanDest+string(os.PathSeparator), cleanTarget) && cleanDest != filepath.Clean(targetDir) {
		return "", fmt.Errorf("path traversal detected: %s", entryName)
	}

	return destPath, nil
}

// extractFile extracts a single file from the ZIP archive
func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("open zip entry: %w", err)
	}
	defer rc.Close()

	// Force safe permissions (never trust archive permissions)
	mode := os.FileMode(0644)
	if f.FileInfo().IsDir() {
		mode = 0755
	}

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer outFile.Close()

	// Copy with size limit
	_, err = io.Copy(outFile, io.LimitReader(rc, MaxSingleFileSize))
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
