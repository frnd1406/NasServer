package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

// AI-indexable MIME types (text-based files only)
var aiIndexableMimeTypes = map[string]bool{
	"text/plain":      true,
	"application/pdf": true,
	"text/markdown":   true,
	"text/csv":        true,
}

// AIAgentPayload represents the data sent to the AI knowledge agent
type AIAgentPayload struct {
	FilePath string `json:"file_path"`
	FileID   string `json:"file_id"`
	MimeType string `json:"mime_type"`
	Content  string `json:"content,omitempty"` // Optional: if set, agent uses this instead of reading from disk
}

// notifyAIAgent sends a fire-and-forget notification to the AI knowledge agent
func notifyAIAgent(filePath, fileID, mimeType string, content string, honeySvc *services.HoneyfileService, internalSecret string, logger *logrus.Logger) {
	// Run asynchronously to avoid blocking the upload response
	go func() {
		// SECURITY: Never index monitored resources
		if honeySvc != nil && honeySvc.CheckAndTrigger(context.Background(), filePath, services.RequestMetadata{Action: "index_scan", IPAddress: "internal"}) {
			logger.WithField("file_path", filePath).Info("Skipping AI indexing for monitored resource")
			return
		}

		// Check if file is eligible for AI indexing
		if !aiIndexableMimeTypes[mimeType] {
			logger.WithFields(logrus.Fields{
				"file_path": filePath,
				"mime_type": mimeType,
			}).Info("Skipping AI indexing for non-text file")
			return
		}

		// If content is empty but mime type is text, log a debug message (Legacy Mode / Disk Read)
		if content == "" {
			logger.WithField("file_id", fileID).Debug("Indexing triggered without inline content (Legacy Mode)")
		}

		payload := AIAgentPayload{
			FilePath: filePath,
			FileID:   fileID,
			MimeType: mimeType,
			Content:  content,
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
		// SECURITY: Internal Auth
		if internalSecret != "" {
			req.Header.Set("X-Internal-Secret", internalSecret)
		}

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

// notifyAIAgentDelete sends a SYNCHRONOUS deletion notification to the AI knowledge agent.
// This prevents "ghost knowledge" by removing embeddings when files are deleted.
// Returns error if AI agent is unreachable, but caller should soft-fail (log, don't block user).
func notifyAIAgentDelete(filePath, fileID string, internalSecret string, logger *logrus.Logger) error {
	payload := map[string]string{
		"file_path": filePath,
		"file_id":   fileID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.WithError(err).Warn("Failed to marshal AI agent delete payload")
		return fmt.Errorf("marshal payload: %w", err)
	}

	// AI agent delete endpoint
	aiAgentURL := "http://ai-knowledge-agent:5000/delete"

	req, err := http.NewRequest("POST", aiAgentURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.WithError(err).Warn("Failed to create AI agent delete request")
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if internalSecret != "" {
		req.Header.Set("X-Internal-Secret", internalSecret)
	}

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
		return fmt.Errorf("AI agent unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"file_id":   fileID,
		}).Info("AI agent deletion completed successfully")
		return nil
	}

	// Non-2xx status - log but treat as soft error
	logger.WithFields(logrus.Fields{
		"status_code": resp.StatusCode,
		"file_path":   filePath,
		"file_id":     fileID,
	}).Warn("AI agent delete returned non-2xx status")
	return fmt.Errorf("AI agent returned status %d", resp.StatusCode)
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

func StorageUploadHandler(storage *services.StorageService, honeySvc *services.HoneyfileService, cfg *config.Config, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.PostForm("path")

		// ==== PHASE 3: Hybrid Encryption Support ====
		// Read encryption parameters from form
		encryptionModeStr := c.PostForm("encryption_mode")      // NONE, SYSTEM, USER
		encryptionPassword := c.PostForm("encryption_password") // Required for USER mode

		// Default to NONE if not specified (backward compatibility)
		var encryptionMode models.EncryptionMode
		switch strings.ToUpper(encryptionModeStr) {
		case "USER":
			encryptionMode = models.EncryptionUser
		case "SYSTEM":
			encryptionMode = models.EncryptionSystem
		case "NONE", "":
			encryptionMode = models.EncryptionNone
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid encryption_mode: must be NONE, SYSTEM, or USER",
			})
			return
		}

		// Validate USER mode has password
		if encryptionMode == models.EncryptionUser && encryptionPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "encryption_password is required when encryption_mode is USER",
			})
			return
		}

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

		// ==== SAVE FILE (with or without encryption) ====
		var result *services.SaveResult

		if encryptionMode == models.EncryptionNone {
			// Legacy path: No encryption
			result, err = storage.Save(path, src, fileHeader)
		} else {
			// New path: Hybrid encryption
			result, err = storage.SaveWithEncryption(path, src, fileHeader, encryptionMode, encryptionPassword)
		}

		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		// Log encryption status
		logger.WithFields(logrus.Fields{
			"request_id":      requestID,
			"filename":        fileHeader.Filename,
			"encryption_mode": encryptionMode,
			"size_bytes":      result.SizeBytes,
			"checksum":        result.Checksum,
		}).Info("File upload completed")

		// ==== AI AGENT NOTIFICATION ====
		// Only index UNENCRYPTED files (can't index encrypted content!)
		if encryptionMode == models.EncryptionNone {
			var extractedText string
			if aiIndexableMimeTypes[result.MimeType] {
				if _, err := src.Seek(0, 0); err == nil {
					const MaxIndexSize = 2 * 1024 * 1024
					buf := new(bytes.Buffer)
					io.CopyN(buf, src, MaxIndexSize)
					extractedText = buf.String()
				}
			}
			notifyAIAgent(result.Path, result.FileID, result.MimeType, extractedText, honeySvc, cfg.InternalAPISecret, logger)
		} else {
			logger.WithField("filename", fileHeader.Filename).Debug("Skipping AI indexing for encrypted file")
		}

		// Return enhanced response with encryption metadata
		c.JSON(http.StatusOK, gin.H{
			"status":            "ok",
			"encryption_status": encryptionMode,
			"size_bytes":        result.SizeBytes,
			"checksum":          result.Checksum,
		})
	}
}

func StorageDownloadHandler(storage *services.StorageService, honeySvc *services.HoneyfileService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.Query("path")
		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}

		// SECURITY: Check for integrity checkpoint BEFORE serving
		fullPath := filepath.Join("/mnt/data", path)

		// Capture Forensic Metadata
		meta := services.RequestMetadata{
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			// UserID:    nil, // TODO: Extract from Auth Context if available
			Action: "download",
		}

		if honeySvc != nil && honeySvc.CheckAndTrigger(c.Request.Context(), fullPath, meta) {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       path,
				"ip":         meta.IPAddress,
			}).Error("🔒 INTEGRITY VIOLATION - VAULT LOCKED")

			// ACTIVE DEFENSE: Return 403 with misleading error or just 403
			c.JSON(http.StatusForbidden, gin.H{"error": "file corrupted: integrity check failed"})
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

func StorageDeleteHandler(storage *services.StorageService, cfg *config.Config, logger *logrus.Logger) gin.HandlerFunc {
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
		// SYNCHRONOUS call - waits for AI agent response before returning 200
		// Soft-fail: log error but don't block user if AI agent is down
		if err := notifyAIAgentDelete(fullPath, fileID, cfg.InternalAPISecret, logger); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"file_path":  fullPath,
				"file_id":    fileID,
				"error":      err.Error(),
			}).Error("SOFT-FAIL: AI agent deletion failed, ghost knowledge may persist")
			// Continue - don't block user, file is deleted from disk
		}

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

// StorageTrashEmptyHandler permanently deletes ALL items from trash
func StorageTrashEmptyHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		// Get all trash items
		items, err := storage.ListTrash()
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		deletedCount := 0
		for _, item := range items {
			if err := storage.DeleteFromTrash(item.ID); err != nil {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"item_id":    item.ID,
					"error":      err.Error(),
				}).Warn("Failed to delete trash item")
				continue
			}
			deletedCount++
		}

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"count":      deletedCount,
		}).Info("Trash emptied")

		c.JSON(http.StatusOK, gin.H{
			"status":  "emptied",
			"deleted": deletedCount,
		})
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

// moveRequest represents the request for moving a file or folder
type moveRequest struct {
	SourcePath      string `json:"sourcePath" binding:"required"`
	DestinationPath string `json:"destinationPath" binding:"required"`
}

// StorageMoveHandler moves a file or folder to a new location
func StorageMoveHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		var req moveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Validate source and destination paths
		if req.SourcePath == "" || req.DestinationPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "source and destination paths are required"})
			return
		}

		// Prevent moving to same location
		if req.SourcePath == req.DestinationPath {
			c.JSON(http.StatusBadRequest, gin.H{"error": "source and destination are the same"})
			return
		}

		// Get full paths
		sourceFull, err := storage.GetFullPath(req.SourcePath)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		destFull, err := storage.GetFullPath(req.DestinationPath)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		// Check source exists
		if _, err := os.Stat(sourceFull); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "source file not found"})
			return
		}

		// Check destination parent directory exists
		destDir := filepath.Dir(destFull)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "destination directory does not exist"})
			return
		}

		// Check destination doesn't already exist
		if _, err := os.Stat(destFull); err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "destination already exists"})
			return
		}

		// Perform the move
		if err := os.Rename(sourceFull, destFull); err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"request_id": requestID,
				"source":     req.SourcePath,
				"dest":       req.DestinationPath,
			}).Error("Failed to move file")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move file"})
			return
		}

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"source":     req.SourcePath,
			"dest":       req.DestinationPath,
		}).Info("File moved successfully")

		c.JSON(http.StatusOK, gin.H{
			"status": "moved",
			"source": req.SourcePath,
			"dest":   req.DestinationPath,
		})
	}
}

// StorageDownloadZipHandler downloads a directory as a ZIP file
func StorageDownloadZipHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.Query("path")
		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}

		// Get the full path for the directory
		fullPath, err := storage.GetFullPath(path)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		// Check if it's a directory
		info, err := os.Stat(fullPath)
		if err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		if !info.IsDir() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path must be a directory"})
			return
		}

		// Create a buffer to write the ZIP to
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		// Walk the directory and add files to ZIP
		err = filepath.Walk(fullPath, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip the root directory itself
			if filePath == fullPath {
				return nil
			}

			// Get relative path for ZIP entry
			relPath, err := filepath.Rel(fullPath, filePath)
			if err != nil {
				return err
			}

			// Skip hidden files and .trash
			if strings.HasPrefix(filepath.Base(relPath), ".") {
				if fileInfo.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if fileInfo.IsDir() {
				// Add directory entry
				_, err := zipWriter.Create(relPath + "/")
				return err
			}

			// Add file to ZIP
			header, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				return err
			}
			header.Name = relPath
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       path,
				"error":      err.Error(),
			}).Error("storage: failed to create ZIP")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ZIP"})
			return
		}

		zipWriter.Close()

		// Get folder name for the ZIP filename
		folderName := filepath.Base(fullPath)
		if folderName == "" || folderName == "." {
			folderName = "download"
		}

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", folderName))
		c.Data(http.StatusOK, "application/zip", buf.Bytes())
	}
}

// batchDownloadRequest represents the request for batch download
type batchDownloadRequest struct {
	Paths []string `json:"paths" binding:"required"`
}

// StorageBatchDownloadHandler downloads multiple files as a ZIP
func StorageBatchDownloadHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		var req batchDownloadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: paths array required"})
			return
		}

		if len(req.Paths) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no files selected"})
			return
		}

		// Create a buffer for ZIP
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		for _, path := range req.Paths {
			fullPath, err := storage.GetFullPath(path)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"path":       path,
					"error":      err.Error(),
				}).Warn("storage: skipping invalid path in batch download")
				continue
			}

			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}

			if info.IsDir() {
				// Add directory contents to ZIP
				baseName := filepath.Base(fullPath)
				err = filepath.Walk(fullPath, func(filePath string, fileInfo os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if filePath == fullPath {
						return nil
					}

					relPath, err := filepath.Rel(fullPath, filePath)
					if err != nil {
						return err
					}

					// Skip hidden files
					if strings.HasPrefix(filepath.Base(relPath), ".") {
						if fileInfo.IsDir() {
							return filepath.SkipDir
						}
						return nil
					}

					zipPath := filepath.Join(baseName, relPath)

					if fileInfo.IsDir() {
						_, err := zipWriter.Create(zipPath + "/")
						return err
					}

					header, err := zip.FileInfoHeader(fileInfo)
					if err != nil {
						return err
					}
					header.Name = zipPath
					header.Method = zip.Deflate

					writer, err := zipWriter.CreateHeader(header)
					if err != nil {
						return err
					}

					file, err := os.Open(filePath)
					if err != nil {
						return err
					}
					defer file.Close()

					_, err = io.Copy(writer, file)
					return err
				})

				if err != nil {
					logger.WithFields(logrus.Fields{
						"request_id": requestID,
						"path":       path,
						"error":      err.Error(),
					}).Warn("storage: error adding directory to batch ZIP")
				}
			} else {
				// Add single file
				header, err := zip.FileInfoHeader(info)
				if err != nil {
					continue
				}
				header.Name = info.Name()
				header.Method = zip.Deflate

				writer, err := zipWriter.CreateHeader(header)
				if err != nil {
					continue
				}

				file, err := os.Open(fullPath)
				if err != nil {
					continue
				}
				_, err = io.Copy(writer, file)
				file.Close()
				if err != nil {
					continue
				}
			}
		}

		zipWriter.Close()

		c.Header("Content-Disposition", "attachment; filename=\"download.zip\"")
		c.Data(http.StatusOK, "application/zip", buf.Bytes())
	}
}

// StorageMkdirHandler creates a new directory
func StorageMkdirHandler(storage *services.StorageService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		var req struct {
			Path string `json:"path" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
			return
		}

		if err := storage.Mkdir(req.Path); err != nil {
			handleStorageError(c, err, logger, requestID)
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "created"})
	}
}
