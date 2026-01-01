package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// EncryptedStorageService wraps StorageService to provide transparent encryption/decryption
// It intercepts Save/Open operations to encrypt files before storage and decrypt on retrieval
type EncryptedStorageService struct {
	storage    *StorageManager
	encryption *EncryptionService
	logger     *logrus.Logger
	// encryptedBasePath is the directory where encrypted files are stored
	// This can be different from the main storage path for demo purposes
	encryptedBasePath string
}

// NewEncryptedStorageService creates a new encrypted storage wrapper
func NewEncryptedStorageService(
	storage *StorageManager,
	encryption *EncryptionService,
	encryptedBasePath string,
	logger *logrus.Logger,
) (*EncryptedStorageService, error) {
	if storage == nil {
		return nil, fmt.Errorf("storage service is required")
	}
	if encryption == nil {
		return nil, fmt.Errorf("encryption service is required")
	}

	// Ensure encrypted base path exists
	if encryptedBasePath != "" {
		if err := os.MkdirAll(encryptedBasePath, 0700); err != nil {
			return nil, fmt.Errorf("create encrypted base path: %w", err)
		}
	}

	return &EncryptedStorageService{
		storage:           storage,
		encryption:        encryption,
		encryptedBasePath: encryptedBasePath,
		logger:            logger,
	}, nil
}

// IsEncryptionEnabled checks if the vault is unlocked and encryption is available
func (e *EncryptedStorageService) IsEncryptionEnabled() bool {
	return e.encryption != nil && e.encryption.IsUnlocked()
}

// SaveEncrypted stores a file with encryption
// The file is encrypted in memory before being written to disk
// For large files, consider using streaming encryption (Phase 3.2)
func (e *EncryptedStorageService) SaveEncrypted(dir string, file multipart.File, fileHeader *multipart.FileHeader) (*SaveResult, error) {
	if !e.IsEncryptionEnabled() {
		return nil, ErrVaultLocked
	}

	filename := fileHeader.Filename
	if filename == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// Validate file size
	if err := e.storage.ValidateFileSize(file, fileHeader); err != nil {
		return nil, err
	}

	// Validate file type
	if err := e.storage.ValidateFileType(file, filename); err != nil {
		return nil, err
	}

	// Read entire file into memory (for encryption)
	// Note: For large files, streaming encryption should be implemented
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	// Encrypt the data
	encryptedData, err := e.encryption.EncryptData(data)
	if err != nil {
		return nil, fmt.Errorf("encrypt file: %w", err)
	}

	// Determine target path
	targetDir := filepath.Join(e.encryptedBasePath, dir)
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return nil, fmt.Errorf("create target dir: %w", err)
	}

	// Add .enc extension to indicate encrypted file
	encryptedFilename := filename + ".enc"
	destPath := filepath.Join(targetDir, encryptedFilename)

	// Write encrypted data
	if err := os.WriteFile(destPath, encryptedData, 0600); err != nil {
		return nil, fmt.Errorf("write encrypted file: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"filename":      filename,
		"encryptedPath": destPath,
		"originalSize":  len(data),
		"encryptedSize": len(encryptedData),
	}).Info("File encrypted and saved")

	// Securely wipe plaintext from memory
	for i := range data {
		data[i] = 0
	}

	return &SaveResult{
		Path:     destPath,
		MimeType: "application/octet-stream", // Encrypted files have no content type
		FileID:   encryptedFilename,
	}, nil
}

// OpenEncrypted retrieves and decrypts a file
// Returns a reader with the decrypted content
func (e *EncryptedStorageService) OpenEncrypted(relPath string) (io.ReadCloser, os.FileInfo, string, error) {
	if !e.IsEncryptionEnabled() {
		return nil, nil, "", ErrVaultLocked
	}

	// Ensure path ends with .enc
	if !strings.HasSuffix(relPath, ".enc") {
		relPath = relPath + ".enc"
	}

	fullPath := filepath.Join(e.encryptedBasePath, relPath)

	// Check file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, nil, "", err
	}

	// Read encrypted file
	encryptedData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, nil, "", fmt.Errorf("read encrypted file: %w", err)
	}

	// Decrypt
	decryptedData, err := e.encryption.DecryptData(encryptedData)
	if err != nil {
		return nil, nil, "", fmt.Errorf("decrypt file: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"path":          relPath,
		"encryptedSize": len(encryptedData),
		"decryptedSize": len(decryptedData),
	}).Debug("File decrypted for reading")

	// Create a reader from decrypted data
	reader := io.NopCloser(bytes.NewReader(decryptedData))

	// Determine original MIME type from filename (without .enc)
	originalName := strings.TrimSuffix(filepath.Base(relPath), ".enc")
	mimeType := "application/octet-stream"
	// Could detect from magic numbers here if needed

	return reader, info, mimeType + "; original=" + originalName, nil
}

// ListEncrypted lists files in the encrypted storage directory
func (e *EncryptedStorageService) ListEncrypted(relPath string) ([]StorageEntry, error) {
	targetDir := filepath.Join(e.encryptedBasePath, relPath)

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []StorageEntry{}, nil
		}
		return nil, err
	}

	var items []StorageEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		// Strip .enc extension for display
		displayName := strings.TrimSuffix(name, ".enc")
		isEncrypted := strings.HasSuffix(name, ".enc")

		items = append(items, StorageEntry{
			Name:     displayName,
			Size:     info.Size(),
			IsDir:    info.IsDir(),
			ModTime:  info.ModTime(),
			MimeType: "encrypted",
			IsImage:  false, // Encrypted files don't reveal content type
		})

		if isEncrypted {
			e.logger.WithField("file", displayName).Debug("Listed encrypted file")
		}
	}

	return items, nil
}

// DeleteEncrypted removes an encrypted file
func (e *EncryptedStorageService) DeleteEncrypted(relPath string) error {
	// Ensure path ends with .enc
	if !strings.HasSuffix(relPath, ".enc") {
		relPath = relPath + ".enc"
	}

	fullPath := filepath.Join(e.encryptedBasePath, relPath)

	// Security check
	if !strings.HasPrefix(fullPath, e.encryptedBasePath) {
		return ErrPathTraversal
	}

	return os.RemoveAll(fullPath)
}

// GetEncryptedBasePath returns the base path for encrypted files
func (e *EncryptedStorageService) GetEncryptedBasePath() string {
	return e.encryptedBasePath
}

// SetEncryptedBasePath updates the encrypted files base path
func (e *EncryptedStorageService) SetEncryptedBasePath(path string) error {
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}
	e.encryptedBasePath = path
	return nil
}

// GetUnderlyingStorage returns the wrapped storage service for non-encrypted operations
func (e *EncryptedStorageService) GetUnderlyingStorage() *StorageManager {
	return e.storage
}

// GetEncryptionService returns the encryption service for status checks
func (e *EncryptedStorageService) GetEncryptionService() *EncryptionService {
	return e.encryption
}
