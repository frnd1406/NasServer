package testutils

import (
	"context"

	"github.com/nas-ai/api/src/domain/auth"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Mock: UserRepository
// ============================================================

// MockUserRepository mocks auth_repo.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (*auth.User, error) {
	args := m.Called(ctx, username, email, passwordHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*auth.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*auth.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	args := m.Called(ctx, userID, newPasswordHash)
	return args.Error(0)
}

func (m *MockUserRepository) VerifyEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// ============================================================
// Mock: JWTService
// ============================================================

// MockJWTService mocks security.JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (interface{}, error) {
	args := m.Called(tokenString)
	return args.Get(0), args.Error(1)
}

// ============================================================
// Mock: PasswordService
// ============================================================

// MockPasswordService mocks security.PasswordService
type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockPasswordService) ValidatePasswordStrength(password string) error {
	args := m.Called(password)
	return args.Error(0)
}

// ============================================================
// Mock: TokenService
// ============================================================

// MockTokenService mocks security.TokenService
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateVerificationToken(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateVerificationToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GeneratePasswordResetToken(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidatePasswordResetToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) InvalidateUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTokenService) IsTokenRevoked(ctx context.Context, userID string, iat int64) (bool, error) {
	args := m.Called(ctx, userID, iat)
	return args.Bool(0), args.Error(1)
}

// ============================================================
// Mock: StorageManager (for file uploads)
// ============================================================

// MockStorageManager mocks content.StorageManager
type MockStorageManager struct {
	mock.Mock
}

// SaveResult mimics content.SaveResult for mock returns
type MockSaveResult struct {
	Path      string
	FileID    string
	MimeType  string
	SizeBytes int64
	Checksum  string
}

func (m *MockStorageManager) Save(path string, reader interface{}, header interface{}) (*MockSaveResult, error) {
	args := m.Called(path, reader, header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockSaveResult), args.Error(1)
}

func (m *MockStorageManager) List(path string) ([]interface{}, error) {
	args := m.Called(path)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockStorageManager) Delete(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// ============================================================
// Mock: EncryptionPolicyService
// ============================================================

type MockEncryptionPolicyService struct {
	mock.Mock
}

func (m *MockEncryptionPolicyService) DetermineMode(filename string, size int64, override string) string {
	args := m.Called(filename, size, override)
	return args.String(0)
}

// ============================================================
// Mock: HoneyfileService
// ============================================================

type MockHoneyfileService struct {
	mock.Mock
}

func (m *MockHoneyfileService) IsHoneyfile(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockHoneyfileService) CheckAndTrigger(ctx context.Context, path string, meta interface{}) bool {
	args := m.Called(ctx, path, meta)
	return args.Bool(0)
}

// ============================================================
// Mock: AIAgentService
// ============================================================

type MockAIAgentService struct {
	mock.Mock
}

func (m *MockAIAgentService) NotifyUpload(path, fileID, mimeType, text string) {
	m.Called(path, fileID, mimeType, text)
}
