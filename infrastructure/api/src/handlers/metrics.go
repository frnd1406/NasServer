package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/models"
	"github.com/nas-ai/api/src/repository"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
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

// SystemMetricsLiveHandler returns real-time system stats (CPU, RAM, Disk) using gopsutil.
// Used for the Admin Dashboard "Health Card".
func SystemMetricsLiveHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. CPU
		cpuPercent, err := cpu.Percent(0, false)
		if err != nil {
			logger.WithError(err).Warn("Failed to get CPU stats")
		}
		cpuVal := 0.0
		if len(cpuPercent) > 0 {
			cpuVal = cpuPercent[0]
		}

		// 2. RAM
		vm, err := mem.VirtualMemory()
		if err != nil {
			logger.WithError(err).Warn("Failed to get RAM stats")
		}
		ramVal := 0.0
		ramTotal := uint64(0)
		if vm != nil {
			ramVal = vm.UsedPercent
			ramTotal = vm.Total
		}

		// 3. Disk (Root)
		diskStat, err := disk.Usage("/")
		if err != nil {
			logger.WithError(err).Warn("Failed to get Disk stats")
		}
		diskVal := 0.0
		diskTotal := uint64(0)
		if diskStat != nil {
			diskVal = diskStat.UsedPercent
			diskTotal = diskStat.Total
		}

		c.JSON(http.StatusOK, gin.H{
			"cpu_percent":  math.Round(cpuVal*100) / 100,
			"ram_percent":  math.Round(ramVal*100) / 100,
			"disk_percent": math.Round(diskVal*100) / 100,
			"ram_total":    ramTotal,
			"disk_total":   diskTotal,
			"timestamp":    time.Now(),
		})
	}
}
