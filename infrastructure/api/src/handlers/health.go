package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthChecker is implemented by dependencies that can be probed.
type HealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// getDiskStats returns disk usage for the given path
func getDiskStats(path string) (total, used, free uint64) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0, 0
	}
	total = stat.Blocks * uint64(stat.Bsize)
	free = stat.Bavail * uint64(stat.Bsize)
	used = total - free
	return
}

// countFilesAndFolders counts files and folders in path
func countFilesAndFolders(root string) (files, folders int) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			folders++
		} else {
			files++
		}
		return nil
	})
	return
}

// Health godoc
// @Summary Health check endpoint
// @Description Returns API health status and dependency information
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status information"
// @Failure 503 {object} map[string]interface{} "Dependency unavailable"
// @Router /health [get]
func Health(db HealthChecker, redis HealthChecker, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dependencies := gin.H{}
		healthy := true

		if db == nil {
			logger.Error("PostgreSQL health check skipped: dependency not provided")
			dependencies["database"] = "unhealthy"
			healthy = false
		} else if err := db.HealthCheck(ctx); err != nil {
			logger.WithError(err).Error("PostgreSQL health check failed")
			dependencies["database"] = "unhealthy"
			healthy = false
		} else {
			dependencies["database"] = "ok"
		}

		if redis == nil {
			logger.Error("Redis health check skipped: dependency not provided")
			dependencies["redis"] = "unhealthy"
			healthy = false
		} else if err := redis.HealthCheck(ctx); err != nil {
			logger.WithError(err).Error("Redis health check failed")
			dependencies["redis"] = "unhealthy"
			healthy = false
		} else {
			dependencies["redis"] = "ok"
		}

		// Get disk stats for /mnt/data
		diskTotal, diskUsed, diskFree := getDiskStats("/mnt/data")
		fileCount, folderCount := countFilesAndFolders("/mnt/data")

		status := gin.H{
			"status":       "ok",
			"timestamp":    time.Now().Format(time.RFC3339),
			"service":      "nas-api",
			"version":      "1.0.0-phase1",
			"dependencies": dependencies,
			"disk_total":   diskTotal,
			"disk_used":    diskUsed,
			"disk_free":    diskFree,
			"file_count":   fileCount,
			"folder_count": folderCount,
		}

		if !healthy {
			status["status"] = "degraded"
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}

		c.JSON(http.StatusOK, status)
	}
}
