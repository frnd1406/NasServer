package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StorageSettingsRequest for saving storage monitoring settings
type StorageSettingsRequest struct {
	WarningThreshold  int  `json:"warningThreshold"`
	CriticalThreshold int  `json:"criticalThreshold"`
	AutoCleanup       bool `json:"autoCleanup"`
	CleanupAgeDays    int  `json:"cleanupAgeDays"`
}

// StorageSettingsGetHandler returns the current storage monitoring settings
func StorageSettingsGetHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		config, err := loadSetupConfig()
		if err != nil {
			logger.WithError(err).Error("Failed to load storage settings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load settings"})
			return
		}

		if config == nil {
			// Return defaults
			c.JSON(http.StatusOK, StorageSettingsRequest{
				WarningThreshold:  80,
				CriticalThreshold: 95,
				AutoCleanup:       false,
				CleanupAgeDays:    90,
			})
			return
		}

		c.JSON(http.StatusOK, StorageSettingsRequest{
			WarningThreshold:  config.StorageMonitor.WarningThreshold,
			CriticalThreshold: config.StorageMonitor.CriticalThreshold,
			AutoCleanup:       config.StorageMonitor.AutoCleanup,
			CleanupAgeDays:    config.StorageMonitor.CleanupAgeDays,
		})
	}
}

// StorageSettingsSaveHandler saves storage monitoring settings to setup.json
func StorageSettingsSaveHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req StorageSettingsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Load existing config
		config, err := loadSetupConfig()
		if err != nil {
			logger.WithError(err).Error("Failed to load config for storage settings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load config"})
			return
		}

		if config == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "setup not complete"})
			return
		}

		// Update storage settings
		config.StorageMonitor = StorageMonitor{
			WarningThreshold:  req.WarningThreshold,
			CriticalThreshold: req.CriticalThreshold,
			AutoCleanup:       req.AutoCleanup,
			CleanupAgeDays:    req.CleanupAgeDays,
		}

		if err := saveSetupConfig(config); err != nil {
			logger.WithError(err).Error("Failed to save storage settings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save settings"})
			return
		}

		logger.Info("Storage settings updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": "saved"})
	}
}
