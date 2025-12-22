package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/repository"
	"github.com/sirupsen/logrus"
)

type MonitoringIngestRequest struct {
	Source     string  `json:"source" binding:"required"`
	CPUPercent float64 `json:"cpu_percent" binding:"required"`
	RAMPercent float64 `json:"ram_percent" binding:"required"`
}

func MonitoringIngestHandler(repo *repository.MonitoringRepository, monitoringToken string, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		if c.GetHeader("X-Monitoring-Token") != monitoringToken {
			logger.WithField("request_id", requestID).Warn("invalid monitoring token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var req MonitoringIngestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("invalid monitoring payload")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		sample := &models.MonitoringSample{
			Source:     req.Source,
			CPUPercent: req.CPUPercent,
			RAMPercent: req.RAMPercent,
		}

		if err := repo.InsertSample(c.Request.Context(), sample); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("failed to store monitoring sample")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store sample"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func MonitoringListHandler(repo *repository.MonitoringRepository, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		samples, err := repo.ListRecent(c.Request.Context(), 50)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("failed to list monitoring samples")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load monitoring data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": samples,
		})
	}
}
