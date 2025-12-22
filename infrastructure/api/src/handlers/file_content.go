package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FileContentHandler returns the raw content of a file for preview
// GET /api/v1/files/content?path=/mnt/data/...
func FileContentHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := c.Query("path")
		if filePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing path parameter"})
			return
		}

		// Security: Ensure path is within /mnt/data
		cleanPath := filepath.Clean(filePath)
		if !strings.HasPrefix(cleanPath, "/mnt/data/") {
			logger.WithField("path", filePath).Warn("Unauthorized file access attempt")
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		// Check if file exists
		info, err := os.Stat(cleanPath)
		if err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
				return
			}
			logger.WithError(err).Error("Failed to stat file")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to access file"})
			return
		}

		// Ensure it's a file, not a directory
		if info.IsDir() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is a directory, not a file"})
			return
		}

		// Security: Limit file size to prevent abuse (max 10MB for preview)
		const maxPreviewSize = 10 * 1024 * 1024 // 10MB
		if info.Size() > maxPreviewSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "file too large for preview",
				"max_size":  maxPreviewSize,
				"file_size": info.Size(),
			})
			return
		}

		// Read file content
		content, err := os.ReadFile(cleanPath)
		if err != nil {
			logger.WithError(err).Error("Failed to read file")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		// Detect content type
		contentType := http.DetectContentType(content)

		// Override for common text formats
		ext := strings.ToLower(filepath.Ext(cleanPath))
		switch ext {
		case ".txt", ".log":
			contentType = "text/plain; charset=utf-8"
		case ".json":
			contentType = "application/json; charset=utf-8"
		case ".md":
			contentType = "text/markdown; charset=utf-8"
		case ".html":
			contentType = "text/html; charset=utf-8"
		case ".js":
			contentType = "application/javascript; charset=utf-8"
		case ".css":
			contentType = "text/css; charset=utf-8"
		case ".xml":
			contentType = "application/xml; charset=utf-8"
		case ".go":
			contentType = "text/plain; charset=utf-8"
		case ".py":
			contentType = "text/plain; charset=utf-8"
		}

		c.Header("Content-Type", contentType)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(http.StatusOK, contentType, content)
	}
}
