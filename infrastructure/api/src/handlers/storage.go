package handlers

import (
	"archive/zip"
	"bytes"
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
func notifyAIAgent(filePath, fileID, mimeType string, content string, logger *logrus.Logger) {
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

		// Extract content for AI indexing (RAM-Push)
		var extractedText string
		if aiIndexableMimeTypes[result.MimeType] {
			// Rewind source stream to beginning
			if _, err := src.Seek(0, 0); err == nil {
				// Limit content extraction to 2MB to avoid OOM
				const MaxIndexSize = 2 * 1024 * 1024

				buf := new(bytes.Buffer)
				io.CopyN(buf, src, MaxIndexSize)
				extractedText = buf.String()
			} else {
				logger.Warn("Could not seek upload stream for AI indexing")
			}
		}

		notifyAIAgent(result.Path, result.FileID, result.MimeType, extractedText, logger)

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
