package handlers

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

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

func StorageListHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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

func StorageUploadHandler(storage *services.StorageManager, policyService *services.EncryptionPolicyService, honeySvc *services.HoneyfileService, aiService *services.AIAgentService, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		path := c.PostForm("path")

		// ==== PHASE 3B: Hybrid Encryption with Policy Support ====
		// Read encryption override from form (AUTO, FORCE_USER, FORCE_NONE)
		encryptionOverride := c.PostForm("encryption_override")
		encryptionPassword := c.PostForm("encryption_password") // Required for USER mode

		// Get file header first to determine encryption mode
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		// Use policy service to determine encryption mode intelligently
		encryptionMode := policyService.DetermineMode(
			fileHeader.Filename,
			fileHeader.Size,
			encryptionOverride,
		)

		// Validate USER mode has password
		if encryptionMode == models.EncryptionUser && encryptionPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "🔐 Verschlüsselung erforderlich",
				"message": "Diese Datei muss verschlüsselt werden (PDF, Dokumente, etc.). Bitte richte zuerst den Vault ein unter Einstellungen → Vault, oder lade die Datei ohne Verschlüsselung hoch.",
				"code":    "VAULT_SETUP_REQUIRED",
				"action":  "Gehe zu Einstellungen → Vault und richte ein Master-Passwort ein.",
			})
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
			result, err = storage.SaveWithEncryption(c.Request.Context(), path, src, fileHeader, encryptionMode, encryptionPassword)
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
			if _, err := src.Seek(0, 0); err == nil {
				const MaxIndexSize = 2 * 1024 * 1024
				buf := new(bytes.Buffer)
				io.CopyN(buf, src, MaxIndexSize)
				extractedText = buf.String()
			}
			aiService.NotifyUpload(result.Path, result.FileID, result.MimeType, extractedText)
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

func StorageDownloadHandler(storage *services.StorageManager, honeySvc *services.HoneyfileService, logger *logrus.Logger) gin.HandlerFunc {
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

func StorageDeleteHandler(storage *services.StorageManager, aiService *services.AIAgentService, logger *logrus.Logger) gin.HandlerFunc {
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
		if err := aiService.NotifyDelete(fullPath, fileID); err != nil {
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

func StorageTrashListHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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

func StorageTrashRestoreHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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

func StorageTrashDeleteHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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
func StorageTrashEmptyHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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

func StorageRenameHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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
func StorageMoveHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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
func StorageDownloadZipHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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
func StorageBatchDownloadHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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
func StorageMkdirHandler(storage *services.StorageManager, logger *logrus.Logger) gin.HandlerFunc {
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
