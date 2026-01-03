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
// Refactored Integration Tests using TestEnv
// ============================================================

func TestLogin_Success(t *testing.T) {
	// 1. Setup
	env := testutils.NewTestEnv(t)
	router := testutils.NewTestRouter(env)

	// 2. Define mocked user
	testUser := &auth.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	// 3. Define Expectations
	env.UserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	env.PasswordSvc.On("ComparePassword", "hashed_password", "SecurePassword123!").Return(nil)
	env.JWTService.On("GenerateAccessToken", "user-123", "test@example.com").Return("mock-access-token", nil)
	env.JWTService.On("GenerateRefreshToken", "user-123", "test@example.com").Return("mock-refresh-token", nil)

	// 4. Execute
	body := `{"email": "test@example.com", "password": "SecurePassword123!"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 5. Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "mock-access-token", resp["access_token"])
	assert.Equal(t, "mock-refresh-token", resp["refresh_token"])

	// 6. Verify all expectations were met
	env.UserRepo.AssertExpectations(t)
	env.PasswordSvc.AssertExpectations(t)
	env.JWTService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// 1. Setup
	env := testutils.NewTestEnv(t)
	router := testutils.NewTestRouter(env)

	testUser := &auth.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	// 2. Define Expectations
	env.UserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
	env.PasswordSvc.On("ComparePassword", "hashed_password", "WrongPassword").Return(assert.AnError) // Fail

	// 3. Execute
	body := `{"email": "test@example.com", "password": "WrongPassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 4. Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	// 1. Setup
	env := testutils.NewTestEnv(t)
	router := testutils.NewTestRouter(env)

	// 2. Define Expectations - no user found
	env.UserRepo.On("FindByEmail", mock.Anything, "missing@example.com").Return(nil, nil)

	// 3. Execute
	body := `{"email": "missing@example.com", "password": "SomePassword123!"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 4. Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRegister_Success(t *testing.T) {
	// 1. Setup
	env := testutils.NewTestEnv(t)
	router := testutils.NewTestRouter(env)

	newUser := &auth.User{
		ID:       "user-new",
		Username: "newuser",
		Email:    "new@example.com",
	}

	// 2. Define Expectations
	env.UserRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil) // No existing
	env.PasswordSvc.On("ValidatePasswordStrength", "SecurePassword123!").Return(nil)
	env.UserRepo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, nil) // No existing
	env.PasswordSvc.On("HashPassword", "SecurePassword123!").Return("hashed_SecurePassword123!", nil)
	env.UserRepo.On("CreateUser", mock.Anything, "newuser", "new@example.com", "hashed_SecurePassword123!").Return(newUser, nil)
	env.JWTService.On("GenerateAccessToken", "user-new", "new@example.com").Return("mock-access", nil)
	env.JWTService.On("GenerateRefreshToken", "user-new", "new@example.com").Return("mock-refresh", nil)

	// 3. Execute
	body := `{"username": "newuser", "email": "new@example.com", "password": "SecurePassword123!", "invite_code": "TEST_INVITE"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 4. Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "mock-access", resp["access_token"])

	// 5. Verify
	env.UserRepo.AssertExpectations(t)
	env.PasswordSvc.AssertExpectations(t)
	env.JWTService.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	// 1. Setup
	env := testutils.NewTestEnv(t)
	router := testutils.NewTestRouter(env)

	existingUser := &auth.User{ID: "existing", Email: "existing@example.com"}

	// 2. Expectations
	env.UserRepo.On("FindByUsername", mock.Anything, "newuser").Return(nil, nil)
	env.PasswordSvc.On("ValidatePasswordStrength", "SecurePassword123!").Return(nil)
	env.UserRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil) // Exists!

	// 3. Execute
	body := `{"username": "newuser", "email": "existing@example.com", "password": "SecurePassword123!", "invite_code": "TEST_INVITE"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 4. Assert
	assert.Equal(t, http.StatusConflict, w.Code)
}
