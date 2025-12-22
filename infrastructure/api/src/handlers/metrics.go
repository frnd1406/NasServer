package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/repository"
	"github.com/sirupsen/logrus"
)

type SystemMetricsRequest struct {
	AgentID   string  `json:"agent_id" binding:"required"`
	CPUUsage  float64 `json:"cpu_usage" binding:"required"`
	RAMUsage  float64 `json:"ram_usage" binding:"required"`
	DiskUsage float64 `json:"disk_usage" binding:"required"`
}

// SystemMetricsHandler nimmt Metriken entgegen und schützt per API-Key-Header.
func SystemMetricsHandler(repo *repository.SystemMetricsRepository, apiKey string, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		if apiKey == "" || c.GetHeader("X-Monitoring-Token") != apiKey {
			logger.WithField("request_id", requestID).Warn("unauthorized system metrics call")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":       "unauthorized",
					"message":    "invalid or missing monitoring token",
					"request_id": requestID,
				},
			})
			return
		}

		var req SystemMetricsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Warn("invalid system metrics payload")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":       "invalid_request",
					"message":    "invalid payload",
					"request_id": requestID,
				},
			})
			return
		}

		metric := &models.SystemMetric{
			AgentID:   req.AgentID,
			CPUUsage:  req.CPUUsage,
			RAMUsage:  req.RAMUsage,
			DiskUsage: req.DiskUsage,
		}

		if err := repo.Insert(c.Request.Context(), metric); err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("failed to store system metrics")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":       "internal_error",
					"message":    "failed to store metrics",
					"request_id": requestID,
				},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":     "ok",
			"request_id": requestID,
		})
	}
}

// SystemMetricsListHandler liefert die neuesten Metriken (öffentlich; read-only).
func SystemMetricsListHandler(repo *repository.SystemMetricsRepository, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")

		limit := 10
		if raw := c.Query("limit"); raw != "" {
			if n, err := strconv.Atoi(raw); err == nil {
				limit = n
			}
		}

		items, err := repo.List(c.Request.Context(), limit)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"error":      err.Error(),
			}).Error("failed to list system metrics")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":       "internal_error",
					"message":    "failed to load metrics",
					"request_id": requestID,
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": items,
		})
	}
}
