package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type validatePathResponse struct {
	Valid    bool   `json:"valid"`
	Exists   bool   `json:"exists"`
	Writable bool   `json:"writable"`
	Message  string `json:"message"`
}

func TestValidatePathHandler_ValidWritablePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	logger := logrus.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := strings.NewReader(`{"path":"` + dir + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/system/validate-path", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ValidatePathHandler(logger)(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp validatePathResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.True(t, resp.Valid)
	require.True(t, resp.Exists)
	require.True(t, resp.Writable)
	require.Equal(t, "path is valid", resp.Message)
}

func TestValidatePathHandler_RelativePathInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := logrus.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := strings.NewReader(`{"path":"relative/path"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/system/validate-path", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ValidatePathHandler(logger)(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp validatePathResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.False(t, resp.Valid)
	require.False(t, resp.Exists)
	require.False(t, resp.Writable)
	require.Contains(t, resp.Message, "absolute")
}

func TestValidatePathHandler_NonWritablePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := logrus.New()

	dir := filepath.Join(t.TempDir(), "readonly")
	require.NoError(t, os.MkdirAll(dir, 0o555))
	defer os.Chmod(dir, 0o755) // ensure cleanup works
	defer os.RemoveAll(dir)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := strings.NewReader(`{"path":"` + dir + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/system/validate-path", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ValidatePathHandler(logger)(c)

	require.Equal(t, http.StatusOK, w.Code)

	var resp validatePathResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.False(t, resp.Valid)
	require.True(t, resp.Exists)
	require.False(t, resp.Writable)
	require.Contains(t, resp.Message, "writable")
}
