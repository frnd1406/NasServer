package testutils

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/nas-ai/api/src/domain/auth"
	"github.com/nas-ai/api/src/domain/files"
	"github.com/nas-ai/api/src/drivers/storage"
	files_repo "github.com/nas-ai/api/src/repository/files"
	"github.com/nas-ai/api/src/services/content"
	"github.com/nas-ai/api/src/services/security"
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

func (m *MockJWTService) ValidateToken(tokenString string) (*security.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*security.TokenClaims), args.Error(1)
}

func (m *MockJWTService) ExtractClaims(tokenString string) (*security.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*security.TokenClaims), args.Error(1)
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

func (m *MockTokenService) IsTokenRevoked(ctx context.Context, userID string, iat int64) bool {
	args := m.Called(ctx, userID, iat)
	return args.Bool(0)
}

// ============================================================
// Mock: EncryptionPolicyService
// ============================================================

type MockEncryptionPolicyService struct {
	mock.Mock
}

func (m *MockEncryptionPolicyService) DetermineMode(filename string, size int64, override string) files.EncryptionMode {
	args := m.Called(filename, size, override)
	return args.Get(0).(files.EncryptionMode)
}

// ============================================================
// Mock: HoneyfileService
// ============================================================

type MockHoneyfileService struct {
	mock.Mock
}

// Implementation of content.HoneyfileServiceInterface
func (m *MockHoneyfileService) CheckAndTrigger(ctx context.Context, path string, meta content.RequestMetadata) bool {
	args := m.Called(ctx, path, meta)
	return args.Bool(0)
}

// Helper (optional, for direct service testing if needed)
func (m *MockHoneyfileService) IsHoneyfile(path string) bool {
	args := m.Called(path)
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

func (m *MockAIAgentService) NotifyDelete(ctx context.Context, path, fileID string) error {
	args := m.Called(ctx, path, fileID)
	return args.Error(0)
}

// ============================================================
// Mock: StorageService (Implements content.StorageService)
// ============================================================

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) Save(dir string, file multipart.File, fileHeader *multipart.FileHeader) (*content.SaveResult, error) {
	args := m.Called(dir, file, fileHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*content.SaveResult), args.Error(1)
}

func (m *MockStorageService) SaveWithEncryption(ctx context.Context, dir string, file multipart.File, fileHeader *multipart.FileHeader, mode files.EncryptionMode, password string) (*content.SaveResult, error) {
	args := m.Called(ctx, dir, file, fileHeader, mode, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*content.SaveResult), args.Error(1)
}

// Open implementation signature must match interface
func (m *MockStorageService) Open(relPath string) (*os.File, os.FileInfo, string, error) {
	args := m.Called(relPath)
	// Return types are complex. Mock must return matching types.
	// os.File can be nil if error.
	f, _ := args.Get(0).(*os.File)
	fi, _ := args.Get(1).(os.FileInfo)
	s := args.String(2)
	return f, fi, s, args.Error(3)
}

func (m *MockStorageService) List(relPath string) ([]content.StorageEntry, error) {
	args := m.Called(relPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]content.StorageEntry), args.Error(1)
}

func (m *MockStorageService) Delete(relPath string) error {
	args := m.Called(relPath)
	return args.Error(0)
}

func (m *MockStorageService) ListTrash() ([]content.TrashEntry, error) {
	args := m.Called()
	return args.Get(0).([]content.TrashEntry), args.Error(1)
}

func (m *MockStorageService) RestoreFromTrash(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorageService) DeleteFromTrash(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorageService) Rename(oldRel, newName string) error {
	args := m.Called(oldRel, newName)
	return args.Error(0)
}

func (m *MockStorageService) Move(srcRel, dstRel string) error {
	args := m.Called(srcRel, dstRel)
	return args.Error(0)
}

func (m *MockStorageService) Mkdir(relPath string) error {
	args := m.Called(relPath)
	return args.Error(0)
}

func (m *MockStorageService) GetFullPath(relPath string) (string, error) {
	args := m.Called(relPath)
	return args.String(0), args.Error(1)
}

// User requested convenience methods (mapped to Listeners/Actions if needed, or kept)
func (m *MockStorageService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (interface{}, error) {
	args := m.Called(ctx, fileHeader)
	return args.Get(0), args.Error(1)
}
func (m *MockStorageService) DownloadFile(path string) (*os.File, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}
func (m *MockStorageService) DeleteFile(path string) error {
	return m.Delete(path)
}
func (m *MockStorageService) ListFiles(path string) ([]interface{}, error) {
	// Adaptation needed if return type mismatch
	args := m.Called(path)
	return args.Get(0).([]interface{}), args.Error(1)
}

// ============================================================
// Mock: FileEmbeddingRepository
// ============================================================

type MockFileEmbeddingRepository struct {
	mock.Mock
}

func (m *MockFileEmbeddingRepository) DeleteByFileID(ctx context.Context, fileID string) (int64, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileEmbeddingRepository) IndexFile(ctx context.Context, fileID string, content string) error {
	args := m.Called(ctx, fileID, content)
	return args.Error(0)
}

// ============================================================
// Mock: HoneyfileRepository
// ============================================================

type MockHoneyfileRepository struct {
	mock.Mock
}

func (m *MockHoneyfileRepository) IsHoneyfile(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockHoneyfileRepository) Create(ctx context.Context, filePath, fileType string, createdBy *uuid.UUID) (*files_repo.Honeyfile, error) {
	args := m.Called(ctx, filePath, fileType, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*files_repo.Honeyfile), args.Error(1)
}

func (m *MockHoneyfileRepository) CreateHoneyfile(ctx context.Context, filePath, fileType string) error {
	// Alias for Create? Or specific signature requested
	args := m.Called(ctx, filePath, fileType)
	return args.Error(0)
}

func (m *MockHoneyfileRepository) GetAllPaths(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// ============================================================
// Mock: EncryptionService
// ============================================================

type MockEncryptionService struct {
	mock.Mock
}

func (m *MockEncryptionService) EncryptData(plaintext []byte) ([]byte, error) {
	args := m.Called(plaintext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncryptionService) DecryptData(ciphertext []byte) ([]byte, error) {
	args := m.Called(ciphertext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncryptionService) IsUnlocked() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEncryptionService) Unlock(masterPassword string) error {
	args := m.Called(masterPassword)
	return args.Error(0)
}

func (m *MockEncryptionService) Lock() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockEncryptionService) Setup(masterPassword string) error {
	args := m.Called(masterPassword)
	return args.Error(0)
}

func (m *MockEncryptionService) IsConfigured() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEncryptionService) GetVaultPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEncryptionService) SetVaultPath(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// ============================================================
// Mock: EncryptedStorageService
// ============================================================

type MockEncryptedStorageService struct {
	mock.Mock
}

func (m *MockEncryptedStorageService) SaveEncrypted(dir string, file multipart.File, fileHeader *multipart.FileHeader) (*content.SaveResult, error) {
	args := m.Called(dir, file, fileHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*content.SaveResult), args.Error(1)
}

func (m *MockEncryptedStorageService) OpenEncrypted(relPath string) (io.ReadCloser, os.FileInfo, string, error) {
	args := m.Called(relPath)
	if args.Get(0) == nil {
		return nil, nil, "", args.Error(3)
	}
	return args.Get(0).(io.ReadCloser), args.Get(1).(os.FileInfo), args.String(2), args.Error(3)
}

func (m *MockEncryptedStorageService) ListEncrypted(relPath string) ([]storage.StorageEntry, error) {
	args := m.Called(relPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]storage.StorageEntry), args.Error(1)
}

func (m *MockEncryptedStorageService) DeleteEncrypted(relPath string) error {
	args := m.Called(relPath)
	return args.Error(0)
}

func (m *MockEncryptedStorageService) IsEncryptionEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockEncryptedStorageService) GetEncryptedBasePath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEncryptedStorageService) SetEncryptedBasePath(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
