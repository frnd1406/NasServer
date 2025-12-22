package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/handlers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type healthyChecker struct{}

func (healthyChecker) HealthCheck(ctx context.Context) error {
	return nil
}

func TestHealthEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create router
	router := gin.New()
	router.GET("/health", handlers.Health(healthyChecker{}, healthyChecker{}, logger))

	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.Equal(t, "ok", response["status"])
	assert.NotEmpty(t, response["timestamp"])
	assert.Equal(t, "nas-api", response["service"])
	assert.Contains(t, response["version"], "phase")
}

func TestHealthEndpoint_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	router := gin.New()
	router.GET("/health", handlers.Health(healthyChecker{}, healthyChecker{}, logger))

	// Test multiple requests to ensure consistency
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d failed", i+1)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	}
}
