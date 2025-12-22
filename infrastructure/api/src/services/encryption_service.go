package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"golang.org/x/crypto/argon2"
)

// Encryption errors
var (
	ErrVaultLocked       = errors.New("vault is locked")
	ErrVaultNotSetup     = errors.New("vault is not configured")
	ErrInvalidPassword   = errors.New("invalid master password")
	ErrAlreadyUnlocked   = errors.New("vault is already unlocked")
	ErrAlreadyLocked     = errors.New("vault is already locked")
	ErrVaultAlreadySetup = errors.New("vault is already configured")
)

// Argon2id parameters (OWASP recommended)
const (
	argon2Memory      = 64 * 1024 // 64 MB
	argon2Iterations  = 3
	argon2Parallelism = 4
	argon2KeyLen      = 32 // 256 bits for AES-256
	saltLen           = 32
	nonceLen          = 12 // 96 bits for GCM
)

// VaultConfig stores vault configuration
type VaultConfig struct {
	Algorithm     string `json:"algorithm"`
	KeyDerivation string `json:"keyDerivation"`
	Version       int    `json:"version"`
}

// EncryptionService manages the encryption vault and keys
type EncryptionService struct {
	vaultPath  string
	isUnlocked bool
	dek        []byte // Data Encryption Key (only in RAM when unlocked)
	mu         sync.RWMutex
}

// NewEncryptionService creates a new encryption service with the specified vault path
func NewEncryptionService(vaultPath string) *EncryptionService {
	return &EncryptionService{
		vaultPath:  vaultPath,
		isUnlocked: false,
		dek:        nil,
	}
}

// GetVaultPath returns the current vault path
func (e *EncryptionService) GetVaultPath() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.vaultPath
}

// SetVaultPath updates the vault path (only when locked)
func (e *EncryptionService) SetVaultPath(path string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.isUnlocked {
		return errors.New("cannot change vault path while unlocked")
	}

	e.vaultPath = path
	return nil
}

// IsConfigured checks if the vault has been set up
func (e *EncryptionService) IsConfigured() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	configPath := filepath.Join(e.vaultPath, "config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

// IsUnlocked returns the current lock state
func (e *EncryptionService) IsUnlocked() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isUnlocked
}

// Setup initializes the vault with a master password
func (e *EncryptionService) Setup(masterPassword string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.IsConfiguredUnsafe() {
		return ErrVaultAlreadySetup
	}

	// Create vault directory
	if err := os.MkdirAll(e.vaultPath, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Generate random salt
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate random DEK (Data Encryption Key)
	dek := make([]byte, argon2KeyLen)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return fmt.Errorf("failed to generate DEK: %w", err)
	}

	// Derive KEK (Key Encryption Key) from master password
	kek := argon2.IDKey([]byte(masterPassword), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLen)

	// Encrypt DEK with KEK
	encryptedDEK, err := e.encryptWithKey(dek, kek)
	if err != nil {
		return fmt.Errorf("failed to encrypt DEK: %w", err)
	}

	// Save salt
	saltPath := filepath.Join(e.vaultPath, "salt.bin")
	if err := os.WriteFile(saltPath, salt, 0600); err != nil {
		return fmt.Errorf("failed to save salt: %w", err)
	}

	// Save encrypted DEK
	dekPath := filepath.Join(e.vaultPath, "encrypted_dek.bin")
	if err := os.WriteFile(dekPath, encryptedDEK, 0600); err != nil {
		return fmt.Errorf("failed to save encrypted DEK: %w", err)
	}

	// Save config
	config := VaultConfig{
		Algorithm:     "aes-256-gcm",
		KeyDerivation: "argon2id",
		Version:       1,
	}
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	configPath := filepath.Join(e.vaultPath, "config.json")
	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Auto-unlock after setup
	e.dek = dek
	e.isUnlocked = true

	return nil
}

// Unlock decrypts the DEK using the master password
func (e *EncryptionService) Unlock(masterPassword string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.IsConfiguredUnsafe() {
		return ErrVaultNotSetup
	}

	if e.isUnlocked {
		return ErrAlreadyUnlocked
	}

	// Read salt
	saltPath := filepath.Join(e.vaultPath, "salt.bin")
	salt, err := os.ReadFile(saltPath)
	if err != nil {
		return fmt.Errorf("failed to read salt: %w", err)
	}

	// Read encrypted DEK
	dekPath := filepath.Join(e.vaultPath, "encrypted_dek.bin")
	encryptedDEK, err := os.ReadFile(dekPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted DEK: %w", err)
	}

	// Derive KEK from master password
	kek := argon2.IDKey([]byte(masterPassword), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLen)

	// Decrypt DEK
	dek, err := e.decryptWithKey(encryptedDEK, kek)
	if err != nil {
		return ErrInvalidPassword
	}

	e.dek = dek
	e.isUnlocked = true

	return nil
}

// Lock securely wipes the DEK from memory and locks the vault.
// Uses multi-pass overwrite to defeat forensic recovery.
func (e *EncryptionService) Lock() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.isUnlocked {
		return ErrAlreadyLocked
	}

	// SECURITY: Multi-pass secure wipe (Operation Ironclad)
	// Pass 1: Set all bits to 1 (0xFF)
	// Pass 2: Set all bits to 0 (0x00)
	// This defeats simple forensic recovery techniques
	if e.dek != nil {
		for i := range e.dek {
			e.dek[i] = 0xFF
		}
		for i := range e.dek {
			e.dek[i] = 0x00
		}
		e.dek = nil
	}

	e.isUnlocked = false

	// SECURITY: Force garbage collection to clear any lingering copies
	runtime.GC()

	return nil
}

// EncryptData encrypts data using the DEK
func (e *EncryptionService) EncryptData(plaintext []byte) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.isUnlocked {
		return nil, ErrVaultLocked
	}

	return e.encryptWithKey(plaintext, e.dek)
}

// DecryptData decrypts data using the DEK
func (e *EncryptionService) DecryptData(ciphertext []byte) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.isUnlocked {
		return nil, ErrVaultLocked
	}

	return e.decryptWithKey(ciphertext, e.dek)
}

// encryptWithKey performs AES-256-GCM encryption
func (e *EncryptionService) encryptWithKey(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Prepend nonce to ciphertext
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decryptWithKey performs AES-256-GCM decryption
func (e *EncryptionService) decryptWithKey(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < nonceLen {
		return nil, errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := ciphertext[:nonceLen]
	ciphertext = ciphertext[nonceLen:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// IsConfiguredUnsafe checks configuration without locking (for internal use)
func (e *EncryptionService) IsConfiguredUnsafe() bool {
	configPath := filepath.Join(e.vaultPath, "config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

// GetStatus returns the current vault status
func (e *EncryptionService) GetStatus() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return map[string]interface{}{
		"locked":     !e.isUnlocked,
		"configured": e.IsConfiguredUnsafe(),
		"vaultPath":  e.vaultPath,
	}
}

// VaultBackupFile represents a file for backup export
type VaultBackupFile struct {
	Filename string
	Content  []byte
}

// GetVaultConfigFiles returns salt.bin and config.json for backup export.
// SECURITY: Does NOT include encrypted_dek.bin - user must remember their password!
func (e *EncryptionService) GetVaultConfigFiles() ([]VaultBackupFile, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.IsConfiguredUnsafe() {
		return nil, ErrVaultNotSetup
	}

	files := make([]VaultBackupFile, 0, 2)

	// Read salt.bin (required for password derivation)
	saltPath := filepath.Join(e.vaultPath, "salt.bin")
	saltData, err := os.ReadFile(saltPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read salt.bin: %w", err)
	}
	files = append(files, VaultBackupFile{Filename: "salt.bin", Content: saltData})

	// Read config.json (vault metadata)
	configPath := filepath.Join(e.vaultPath, "config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.json: %w", err)
	}
	files = append(files, VaultBackupFile{Filename: "config.json", Content: configData})

	return files, nil
}
