package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nas-ai/api/src/domain/auth"
	"github.com/nas-ai/api/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ============================================================
// REAL INTEGRATION TESTS
// These test the ACTUAL handler code paths with mocked backends
// NO inline fake handlers - we test the real business logic
// ============================================================

func TestIntegration_Register_RealFlow(t *testing.T) {
	// 1. Setup - uses real testable handlers
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	// 2. Expectations (The Mocking)
	// We expect the REAL handler to call these methods on our mocks
	env.UserRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil)
	env.PasswordSvc.On("ValidatePasswordStrength", "SecurePassword123!").Return(nil)
	env.UserRepo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, nil)
	env.PasswordSvc.On("HashPassword", "SecurePassword123!").Return("hashed_password", nil)

	createdUser := &auth.User{
		ID:       "user-new-123",
		Username: "newuser",
		Email:    "new@example.com",
	}
	env.UserRepo.On("CreateUser", mock.Anything, "newuser", "new@example.com", "hashed_password").
		Return(createdUser, nil)
	env.JWTService.On("GenerateAccessToken", "user-new-123", "new@example.com").
		Return("real-access-token", nil)
	env.JWTService.On("GenerateRefreshToken", "user-new-123", "new@example.com").
		Return("real-refresh-token", nil)

	// 3. Execution (The HTTP Request)
	body := `{"username": "newuser", "email": "new@example.com", "password": "SecurePassword123!", "invite_code": "TEST_INVITE"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 4. Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "real-access-token", resp["access_token"])
	assert.Equal(t, "real-refresh-token", resp["refresh_token"])

	// 5. Verify mocks were called (CRITICAL!)
	env.UserRepo.AssertExpectations(t)
	env.PasswordSvc.AssertExpectations(t)
	env.JWTService.AssertExpectations(t)
}

func TestIntegration_Register_DuplicateEmail(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	existingUser := &auth.User{ID: "existing", Email: "existing@example.com"}

	env.UserRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil)
	env.PasswordSvc.On("ValidatePasswordStrength", "SecurePassword123!").Return(nil)
	env.UserRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	body := `{"username": "newuser", "email": "existing@example.com", "password": "SecurePassword123!", "invite_code": "TEST_INVITE"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "email_exists")

	env.UserRepo.AssertExpectations(t)
}

func TestIntegration_Register_WeakPassword(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	env.UserRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil)
	env.PasswordSvc.On("ValidatePasswordStrength", "weak").Return(assert.AnError)

	body := `{"username": "newuser", "email": "new@example.com", "password": "weak", "invite_code": "TEST_INVITE"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "weak_password")

	env.UserRepo.AssertExpectations(t)
	env.PasswordSvc.AssertExpectations(t)
}

func TestIntegration_Register_InvalidInviteCode(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	body := `{"username": "newuser", "email": "new@example.com", "password": "SecurePassword123!", "invite_code": "WRONG_CODE"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid_invite_code")
}

func TestIntegration_Login_RealFlow(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	testUser := &auth.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	env.UserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	env.PasswordSvc.On("ComparePassword", "hashed_password", "SecurePassword123!").Return(nil)
	env.JWTService.On("GenerateAccessToken", "user-123", "test@example.com").Return("login-access-token", nil)
	env.JWTService.On("GenerateRefreshToken", "user-123", "test@example.com").Return("login-refresh-token", nil)

	body := `{"email": "test@example.com", "password": "SecurePassword123!"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "login-access-token", resp["access_token"])

	env.UserRepo.AssertExpectations(t)
	env.PasswordSvc.AssertExpectations(t)
	env.JWTService.AssertExpectations(t)
}

func TestIntegration_Login_InvalidCredentials(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	testUser := &auth.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	env.UserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	env.PasswordSvc.On("ComparePassword", "hashed_password", "WrongPassword").Return(assert.AnError)

	body := `{"email": "test@example.com", "password": "WrongPassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid_credentials")
}

func TestIntegration_Login_UserNotFound(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	env.UserRepo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, nil)

	body := `{"email": "nonexistent@example.com", "password": "AnyPassword123!"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid_credentials")
}

// ============================================================
// PROTECTED ENDPOINT TESTS (uses REAL AuthMiddleware)
// ============================================================

func TestIntegration_ProtectedEndpoint_ValidToken(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	// Generate a REAL JWT token
	jwtService, err := testutils.NewRealJWTService(env)
	require.NoError(t, err)

	token, err := jwtService.GenerateAccessToken("user-protected", "protected@example.com")
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user-protected")
}

func TestIntegration_ProtectedEndpoint_NoToken(t *testing.T) {
	env := testutils.NewTestEnv(t)
	router := testutils.NewRealRouter(env)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
