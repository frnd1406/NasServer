package testutils

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// TestEnv holds all dependencies needed for integration tests.
// It separates mock initialization from router wiring (SRP).
type TestEnv struct {
	// Mocks
	UserRepo     *MockUserRepository
	JWTService   *MockJWTService
	PasswordSvc  *MockPasswordService
	TokenService *MockTokenService

	// Real Services (with fake backends)
	RedisClient *database.RedisClient
	Config      *config.Config
	Logger      *logrus.Logger

	// Internal: for cleanup
	miniredis *miniredis.Miniredis
}

// NewTestEnv creates a fully initialized TestEnv.
// Call t.Cleanup() is handled automatically.
func NewTestEnv(t *testing.T) *TestEnv {
	// Setup miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(func() { mr.Close() })

	// Redis client pointing to miniredis
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	redisClient := &database.RedisClient{Client: rdb}

	// Logger (silent for tests)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Config with test defaults
	cfg := &config.Config{
		JWTSecret:    "test-secret-at-least-32-characters-long",
		ResendAPIKey: "re_test_123",
		EmailFrom:    "test@nas.ai",
		FrontendURL:  "http://localhost:3000",
		InviteCode:   "TEST_INVITE",
	}

	return &TestEnv{
		// Mocks
		UserRepo:     new(MockUserRepository),
		JWTService:   new(MockJWTService),
		PasswordSvc:  new(MockPasswordService),
		TokenService: new(MockTokenService),

		// Real with fakes
		RedisClient: redisClient,
		Config:      cfg,
		Logger:      logger,
		miniredis:   mr,
	}
}

// ResetMocks clears all mock expectations (useful for sub-tests).
func (e *TestEnv) ResetMocks() {
	e.UserRepo = new(MockUserRepository)
	e.JWTService = new(MockJWTService)
	e.PasswordSvc = new(MockPasswordService)
	e.TokenService = new(MockTokenService)
}
