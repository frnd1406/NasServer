package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nas-ai/api/src/domain/auth"
	"github.com/nas-ai/api/src/services/security"
	"github.com/nas-ai/api/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestLogin_Success_WithMocks demonstrates how to use the DI router with Mocks
// directly, verifying that Real Handlers interact with services correctly.
func TestLogin_Success_WithMocks(t *testing.T) {
	// 1. Setup Mocks
	mockUserRepo := new(testutils.MockUserRepository)
	mockJWTService := new(testutils.MockJWTService)
	mockPasswordService := new(testutils.MockPasswordService)
	// We don't need TokenService for Login, but SetupTestRouter requires it
	mockTokenService := new(testutils.MockTokenService)

	// 2. Setup Expectations
	// User found
	mockUser := &auth.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

	// Password valid
	mockPasswordService.On("ComparePassword", "hashed_password", "password123").Return(nil)

	// Tokens generated
	mockJWTService.On("GenerateAccessToken", "user-1", "test@example.com").Return("mock_access_token", nil)
	mockJWTService.On("GenerateRefreshToken", "user-1", "test@example.com").Return("mock_refresh_token", nil)

	// Note: LoginHandler (Real) also uses Redis for CSRF token generation (logic.GeneratCSRFToken).
	// Currently SetupTestRouter uses env.RedisClient (Real or Miniredis).
	// We can pass a "Real" Miniredis env for the Redis part while mocking services.
	env := testutils.NewTestEnv(t) // This provides Miniredis

	// 3. Initialize Router with Mocks
	router := testutils.SetupTestRouter(
		env, // Provides Redis + Config + Logger
		mockUserRepo,
		mockJWTService,
		mockPasswordService,
		mockTokenService,
		nil, nil, nil, nil, // File services
		nil, nil, // Encryption services
	)

	// 4. Execute Request
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 5. Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check user in response
	userData, ok := response["user"].(map[string]interface{})
	assert.True(t, ok, "User field should be present")
	assert.Equal(t, "user-1", userData["id"])
	assert.Equal(t, "test@example.com", userData["email"])

	// Check Cookies for tokens (since handlers set HttpOnly cookies)
	cookies := w.Result().Cookies()
	var accessTokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "access_token" {
			accessTokenCookie = c
			break
		}
	}
	assert.NotNil(t, accessTokenCookie, "Access token cookie should be set")
	assert.Equal(t, "mock_access_token", accessTokenCookie.Value)

	// Verify Mocks
	mockUserRepo.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
}

func TestLogin_Failure_WithMocks(t *testing.T) {
	mockUserRepo := new(testutils.MockUserRepository)
	mockJWTService := new(testutils.MockJWTService)
	mockPasswordService := new(testutils.MockPasswordService)
	mockTokenService := new(testutils.MockTokenService)
	env := testutils.NewTestEnv(t)

	// Expectation: User found, but password invalid
	mockUser := &auth.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

	// Password INVALID
	mockPasswordService.On("ComparePassword", "hashed_password", "wrong_password").Return(security.ErrInvalidPassword) // Assuming security package has this error or just generic error

	router := testutils.SetupTestRouter(
		env,
		mockUserRepo,
		mockJWTService,
		mockPasswordService,
		mockTokenService,
		nil, nil, nil, nil, // File services
		nil, nil, // Encryption services
	)

	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "wrong_password",
	}
	jsonData, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockPasswordService.AssertExpectations(t)
}
