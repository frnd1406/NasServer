package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

// AI-indexable MIME types (text-based files only)
var aiIndexableMimeTypes = map[string]bool{
	"text/plain":       true,
	"application/pdf":  true,
	"text/markdown":    true,
	"text/csv":         true,
}

// AIAgentPayload represents the data sent to the AI knowledge agent
type AIAgentPayload struct {
	FilePath string `json:"file_path"`
	FileID   string `json:"file_id"`
	MimeType string `json:"mime_type"`
}

// notifyAIAgent sends a fire-and-forget notification to the AI knowledge agent
func notifyAIAgent(filePath, fileID, mimeType string, logger *logrus.Logger) {
	// Run asynchronously to avoid blocking the upload response
	go func() {
		// Check if file is eligible for AI indexing
		if !aiIndexableMimeTypes[mimeType] {
			logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"mime_type": mimeType,
			}).Info("Skipping AI indexing for non-text file")
			return
		}

		payload := AIAgentPayload{
			FilePath: filePath,
			FileID:   fileID,
			MimeType: mimeType,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			logger.WithError(err).Warn("Failed to marshal AI agent payload")
			return
		}

		// AI agent endpoint (internal Docker DNS)
		aiAgentURL := "http://ai-knowledge-agent:5000/process"

		// Create HTTP request with timeout
		req, err := http.NewRequest("POST", aiAgentURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			logger.WithError(err).Warn("Failed to create AI agent request")
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// Use a client with timeout to prevent hanging
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"url":       aiAgentURL,
				"file_path": filePath,
				"error":     err.Error(),
			}).Warn("Failed to trigger AI agent")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"file_id":   fileID,
				"mime_type": mimeType,
			}).Info("Triggered AI agent successfully")
		} else {
			logger.WithFields(logrus.Fields{
				"status_code": resp.StatusCode,
				"file_path":   filePath,
			}).Warn("AI agent returned non-200 status")
		}
	}()
}

// notifyAIAgentDelete sends a fire-and-forget deletion notification to the AI knowledge agent
// This prevents "ghost knowledge" by removing embeddings when files are deleted
func notifyAIAgentDelete(filePath, fileID string, logger *logrus.Logger) {
	go func() {
		payload := map[string]string{
			"file_path": filePath,
			"file_id":   fileID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			logger.WithError(err).Warn("Failed to marshal AI agent delete payload")
			return
		}

		// AI agent delete endpoint
		aiAgentURL := "http://ai-knowledge-agent:5000/delete"

		req, err := http.NewRequest("POST", aiAgentURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			logger.WithError(err).Warn("Failed to create AI agent delete request")
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"url":       aiAgentURL,
				"file_path": filePath,
				"file_id":   fileID,
				"error":     err.Error(),
			}).Warn("Failed to notify AI agent of deletion")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"file_id":   fileID,
			}).Info("AI agent deletion triggered successfully")
		} else {
			logger.WithFields(logrus.Fields{
				"status_code": resp.StatusCode,
				"file_path":   filePath,
				"file_id":     fileID,
			}).Warn("AI agent delete returned non-2xx status")
		}
	}()
}

func handleStorageError(c *gin.Context, err error, logger *logrus.Logger, requestID string) {
	status := http.StatusBadRequest
	message := "storage operation failed"

	// Map specific errors to appropriate HTTP status codes and messages
	if errors.Is(err, services.ErrPathTraversal) {
		status = http.StatusForbidden
		message = "access denied: path traversal detected"
	} else if errors.Is(err, services.ErrInvalidFileType) {
		status = http.StatusBadRequest
		message = "invalid file type: only images, documents, videos, and archives are allowed"
	} else if errors.Is(err, services.ErrFileTooLarge) {
		status = http.StatusBadRequest
		message = "file too large: maximum upload size is 100MB"
	} else if os.IsNotExist(err) {
		status = http.StatusNotFound
		message = "file or directory not found"
	}

	logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"error":      err.Error(),
		"status":     status,
	}).Warn("storage: request failed")

	c.JSON(status, gin.H{
		"error": gin.H{
			"code":       "storage_error",
			"message":    message,
			"request_id": requestID,
		},
	})
}

type renameRequest struct {
	OldPath string `json:"oldPath" binding:"required"`
	NewName string `json:"newName" binding:"required"`
}

func StorageListHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.Query("path")

		items, err := storage.List(path)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": items,
		})
	}
}

func StorageUploadHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.PostForm("path")

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		src, err := fileHeader.Open()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("storage: open upload file failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read upload"})
			return
		}
		defer src.Close()

		// Save file and get metadata
		result, err := storage.Save(path, src, fileHeader)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		// Trigger AI agent notification (fire & forget)
		// This happens AFTER successful save to disk
		notifyAIAgent(result.Path, result.FileID, result.MimeType, logger)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func StorageDownloadHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.Query("path")
		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}

		file, info, ctype, err := storage.Open(path)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}
		defer file.Close()

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", info.Name()))
		c.DataFromReader(http.StatusOK, info.Size(), ctype, file, nil)
	}
}

func StorageDeleteHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.Query("path")
		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}

		// Extract fileID for AI agent notification (before deletion!)
		fileID := filepath.Base(path)
		// Construct full path for AI agent (assuming /mnt/data base path)
		fullPath := filepath.Join("/mnt/data", path)

		if err := storage.Delete(path); err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		// Notify AI agent to delete embeddings (prevents ghost knowledge)
		// This happens AFTER successful deletion from filesystem
		notifyAIAgentDelete(fullPath, fileID, logger)

		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

func StorageTrashListHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		items, err := storage.ListTrash()
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	}
}

func StorageTrashRestoreHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		if err := storage.RestoreFromTrash(filepath.ToSlash(id)); err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "restored"})
	}
}

func StorageTrashDeleteHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		if err := storage.DeleteFromTrash(filepath.ToSlash(id)); err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

func StorageRenameHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		var req renameRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		if err := storage.Rename(req.OldPath, req.NewName); err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "renamed"})
	}
}
