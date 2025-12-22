package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/repository"
	"github.com/nas-ai/api/src/scheduler"
	"github.com/nas-ai/api/src/services"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type BackupSettingsRequest struct {
	Schedule  string `json:"schedule" binding:"required"`
	Retention int    `json:"retention" binding:"required"`
	Path      string `json:"path" binding:"required"`
}

type ValidatePathRequest struct {
	Path string `json:"path" binding:"required"`
}

func SystemSettingsHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"backup": gin.H{
				"schedule":  cfg.BackupSchedule,
				"retention": cfg.BackupRetentionCount,
				"path":      cfg.BackupStoragePath,
			},
		})
	}
}

// ValidatePathHandler checks if a path is absolute, exists, and is writable by creating a temp file.
func ValidatePathHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ValidatePathRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		path := filepath.Clean(strings.TrimSpace(req.Path))
		response := gin.H{
			"valid":    false,
			"exists":   false,
			"writable": false,
			"message":  "",
		}

		if path == "" {
			response["message"] = "path is required"
			c.JSON(http.StatusOK, response)
			return
		}
		if !filepath.IsAbs(path) {
			response["message"] = "path must be absolute"
			c.JSON(http.StatusOK, response)
			return
		}

		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				response["message"] = "path does not exist"
			} else {
				logger.WithError(err).Warn("validate path: stat failed")
				response["message"] = "unable to read path metadata"
			}
			c.JSON(http.StatusOK, response)
			return
		}

		response["exists"] = true

		if !info.IsDir() {
			response["message"] = "path must be a directory"
			c.JSON(http.StatusOK, response)
			return
		}

		tmp, err := os.CreateTemp(path, ".nas-path-check-*")
		if err != nil {
			logger.WithError(err).Warn("validate path: write check failed")
			response["message"] = "path is not writable"
			c.JSON(http.StatusOK, response)
			return
		}
		tmp.Close()
		os.Remove(tmp.Name())

		response["writable"] = true
		response["valid"] = true
		response["message"] = "path is valid"

		c.JSON(http.StatusOK, response)
	}
}

func UpdateBackupSettingsHandler(cfg *config.Config, backupSvc *services.BackupService, settingsRepo *repository.SystemSettingsRepository, logger *logrus.Logger) gin.HandlerFunc {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	return func(c *gin.Context) {
		var req BackupSettingsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		req.Schedule = strings.TrimSpace(req.Schedule)
		req.Path = filepath.Clean(strings.TrimSpace(req.Path))

		if _, err := parser.Parse(req.Schedule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule format"})
			return
		}
		if req.Retention < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "retention must be >= 1"})
			return
		}
		if req.Path == "" || req.Path == "." || req.Path == string(os.PathSeparator) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid backup path"})
			return
		}

		if err := backupSvc.SetBackupPath(req.Path); err != nil {
			logger.WithError(err).Warn("backup settings: failed to set backup path")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cfg.BackupSchedule = req.Schedule
		cfg.BackupRetentionCount = req.Retention
		cfg.BackupStoragePath = req.Path

		if err := settingsRepo.UpsertMany(c.Request.Context(), map[string]string{
			repository.SystemSettingBackupSchedule:  cfg.BackupSchedule,
			repository.SystemSettingBackupRetention: fmt.Sprintf("%d", cfg.BackupRetentionCount),
			repository.SystemSettingBackupPath:      cfg.BackupStoragePath,
		}); err != nil {
			logger.WithError(err).Error("backup settings: failed to persist settings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist settings"})
			return
		}

		if err := scheduler.RestartScheduler(); err != nil {
			logger.WithError(err).Error("backup settings: failed to restart scheduler")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to restart scheduler"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"backup": gin.H{
				"schedule":  cfg.BackupSchedule,
				"retention": cfg.BackupRetentionCount,
				"path":      cfg.BackupStoragePath,
			},
		})
	}
}
