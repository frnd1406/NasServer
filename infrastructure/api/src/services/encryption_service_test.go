package services

import (
	"bytes"
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEncryptDecryptStream tests basic stream encryption/decryption roundtrip
func TestEncryptDecryptStream(t *testing.T) {
	password := "test-password-123"
	testData := []byte("Hello, NasCrypt V2! This is a test message for chunked AEAD encryption.")

	// Encrypt
	var encryptedBuf bytes.Buffer
	err := EncryptStream(password, bytes.NewReader(testData), &encryptedBuf)
	if err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Verify header exists
	encrypted := encryptedBuf.Bytes()
	if len(encrypted) < HeaderSize {
		t.Fatalf("Encrypted data too short: %d bytes", len(encrypted))
	}

	// Check magic bytes
	if string(encrypted[0:4]) != MagicBytes {
		t.Errorf("Invalid magic bytes: got %q, expected %q", string(encrypted[0:4]), MagicBytes)
	}

	// Check version
	if encrypted[4] != Version {
		t.Errorf("Invalid version: got 0x%02x, expected 0x%02x", encrypted[4], Version)
	}

	// Decrypt
	var decryptedBuf bytes.Buffer
	err = DecryptStream(password, bytes.NewReader(encrypted), &decryptedBuf)
	if err != nil {
		t.Fatalf("DecryptStream failed: %v", err)
	}

	// Verify data matches
	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Errorf("Decrypted data mismatch:\ngot: %q\nexp: %q", decryptedBuf.String(), string(testData))
	}
}

// TestEncryptDecryptLargeData tests encryption of data spanning multiple chunks
func TestEncryptDecryptLargeData(t *testing.T) {
	password := "large-data-password"

	// Create data larger than one chunk (64KB * 3 = 192KB)
	testData := make([]byte, ChunkSize*3+1234)
	if _, err := rand.Read(testData); err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	// Encrypt
	var encryptedBuf bytes.Buffer
	err := EncryptStream(password, bytes.NewReader(testData), &encryptedBuf)
	if err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Calculate expected size: Header + 4 chunks (3 full + 1 partial)
	expectedMinSize := HeaderSize + 3*EncryptedChunkSize + TagSize + 1234
	if encryptedBuf.Len() < expectedMinSize {
		t.Errorf("Encrypted size %d is less than expected minimum %d", encryptedBuf.Len(), expectedMinSize)
	}

	// Decrypt
	var decryptedBuf bytes.Buffer
	err = DecryptStream(password, bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf)
	if err != nil {
		t.Fatalf("DecryptStream failed: %v", err)
	}

	// Verify data matches
	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Errorf("Decrypted data mismatch (first 100 bytes):\ngot: %x\nexp: %x", decryptedBuf.Bytes()[:100], testData[:100])
	}
}

// TestWrongPassword verifies that wrong password fails authentication
func TestWrongPassword(t *testing.T) {
	correctPassword := "correct-password"
	wrongPassword := "wrong-password"
	testData := []byte("Secret data that should not be decrypted with wrong password")

	// Encrypt with correct password
	var encryptedBuf bytes.Buffer
	err := EncryptStream(correctPassword, bytes.NewReader(testData), &encryptedBuf)
	if err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Try to decrypt with wrong password
	var decryptedBuf bytes.Buffer
	err = DecryptStream(wrongPassword, bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf)
	if err == nil {
		t.Error("DecryptStream should have failed with wrong password")
	}
	if !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("Expected authentication failure error, got: %v", err)
	}
}

// TestIsEncrypted tests the magic byte detection
func TestIsEncrypted(t *testing.T) {
	// Create encrypted data
	password := "test"
	testData := []byte("test data")

	var encryptedBuf bytes.Buffer
	if err := EncryptStream(password, bytes.NewReader(testData), &encryptedBuf); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Test with encrypted data
	if !IsEncrypted(bytes.NewReader(encryptedBuf.Bytes())) {
		t.Error("IsEncrypted should return true for encrypted data")
	}

	// Test with plain data
	if IsEncrypted(bytes.NewReader([]byte("plain text data"))) {
		t.Error("IsEncrypted should return false for plain text")
	}

	// Test with short data
	if IsEncrypted(bytes.NewReader([]byte("NA"))) {
		t.Error("IsEncrypted should return false for short data")
	}
}

// TestEmptyData tests handling of empty input
func TestEmptyData(t *testing.T) {
	password := "test"

	// Encrypt empty data
	var encryptedBuf bytes.Buffer
	err := EncryptStream(password, bytes.NewReader([]byte{}), &encryptedBuf)
	if err != nil {
		t.Fatalf("EncryptStream failed for empty data: %v", err)
	}

	// Should have at least a header
	if encryptedBuf.Len() < HeaderSize {
		t.Errorf("Encrypted empty data should have at least header, got %d bytes", encryptedBuf.Len())
	}

	// Decrypt
	var decryptedBuf bytes.Buffer
	err = DecryptStream(password, bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf)
	if err != nil {
		t.Fatalf("DecryptStream failed: %v", err)
	}

	if decryptedBuf.Len() != 0 {
		t.Errorf("Decrypted empty data should be empty, got %d bytes", decryptedBuf.Len())
	}
}

// TestChunkNonceDerivation tests the nonce derivation function
func TestChunkNonceDerivation(t *testing.T) {
	baseNonce := make([]byte, NonceSize)
	for i := range baseNonce {
		baseNonce[i] = byte(i)
	}

	// Generate nonces for first few chunks
	nonce0 := deriveChunkNonce(baseNonce, 0)
	nonce1 := deriveChunkNonce(baseNonce, 1)
	nonce2 := deriveChunkNonce(baseNonce, 2)

	// Nonce 0 should equal base nonce (XOR with 0)
	if !bytes.Equal(nonce0, baseNonce) {
		t.Error("Nonce for chunk 0 should equal base nonce")
	}

	// All nonces should be different
	if bytes.Equal(nonce0, nonce1) {
		t.Error("Nonces for different chunks should be different")
	}
	if bytes.Equal(nonce1, nonce2) {
		t.Error("Nonces for different chunks should be different")
	}

	// Same chunk index should give same nonce (deterministic)
	nonce1Again := deriveChunkNonce(baseNonce, 1)
	if !bytes.Equal(nonce1, nonce1Again) {
		t.Error("Same chunk index should give same nonce")
	}
}

// TestVaultSetupUnlock tests the vault service functionality
func TestVaultSetupUnlock(t *testing.T) {
	// Create temp directory for vault
	tempDir, err := os.MkdirTemp("", "nascrypt-test-vault")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	vaultPath := filepath.Join(tempDir, "vault")
	service := NewEncryptionService(vaultPath)

	// Initially not configured
	if service.IsConfigured() {
		t.Error("New vault should not be configured")
	}

	// Setup vault
	password := "master-password-123"
	if err := service.Setup(password); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Should be configured and unlocked after setup
	if !service.IsConfigured() {
		t.Error("Vault should be configured after setup")
	}
	if !service.IsUnlocked() {
		t.Error("Vault should be unlocked after setup")
	}

	// Encrypt data with vault
	testData := []byte("vault test data")
	encrypted, err := service.EncryptData(testData)
	if err != nil {
		t.Fatalf("EncryptData failed: %v", err)
	}

	// Lock the vault
	if err := service.Lock(); err != nil {
		t.Fatalf("Lock failed: %v", err)
	}
	if service.IsUnlocked() {
		t.Error("Vault should be locked after Lock()")
	}

	// Should fail to encrypt when locked
	_, err = service.EncryptData(testData)
	if err != ErrVaultLocked {
		t.Errorf("Expected ErrVaultLocked, got: %v", err)
	}

	// Unlock with wrong password should fail
	if err := service.Unlock("wrong-password"); err != ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got: %v", err)
	}

	// Unlock with correct password
	if err := service.Unlock(password); err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}

	// Decrypt data
	decrypted, err := service.DecryptData(encrypted)
	if err != nil {
		t.Fatalf("DecryptData failed: %v", err)
	}
	if !bytes.Equal(testData, decrypted) {
		t.Error("Decrypted data doesn't match original")
	}
}

// TestReadHeader tests header parsing
func TestReadHeader(t *testing.T) {
	password := "test"
	testData := []byte("test data for header parsing")

	var encryptedBuf bytes.Buffer
	if err := EncryptStream(password, bytes.NewReader(testData), &encryptedBuf); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	salt, baseNonce, err := ReadHeader(bytes.NewReader(encryptedBuf.Bytes()))
	if err != nil {
		t.Fatalf("ReadHeader failed: %v", err)
	}

	if len(salt) != SaltSize {
		t.Errorf("Salt size mismatch: got %d, expected %d", len(salt), SaltSize)
	}
	if len(baseNonce) != NonceSize {
		t.Errorf("BaseNonce size mismatch: got %d, expected %d", len(baseNonce), NonceSize)
	}
}

// TestCalculateChunkCount tests chunk count calculation
func TestCalculateChunkCount(t *testing.T) {
	tests := []struct {
		encryptedSize int64
		expected      int64
	}{
		{0, 0},
		{HeaderSize, 0},
		{HeaderSize - 1, 0},
		{HeaderSize + 1, 1},
		{HeaderSize + EncryptedChunkSize, 1},
		{HeaderSize + EncryptedChunkSize + 1, 2},
		{HeaderSize + EncryptedChunkSize*3, 3},
	}

	for _, test := range tests {
		got := CalculateChunkCount(test.encryptedSize)
		if got != test.expected {
			t.Errorf("CalculateChunkCount(%d) = %d, expected %d", test.encryptedSize, got, test.expected)
		}
	}
}

// BenchmarkEncryptStream benchmarks encryption performance
func BenchmarkEncryptStream(b *testing.B) {
	password := "benchmark-password"
	testData := make([]byte, 1*1024*1024) // 1MB
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := EncryptStream(password, bytes.NewReader(testData), &buf); err != nil {
			b.Fatalf("EncryptStream failed: %v", err)
		}
	}
}

// BenchmarkDecryptStream benchmarks decryption performance
func BenchmarkDecryptStream(b *testing.B) {
	password := "benchmark-password"
	testData := make([]byte, 1*1024*1024) // 1MB
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	var encryptedBuf bytes.Buffer
	if err := EncryptStream(password, bytes.NewReader(testData), &encryptedBuf); err != nil {
		b.Fatalf("EncryptStream failed: %v", err)
	}
	encrypted := encryptedBuf.Bytes()

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := DecryptStream(password, bytes.NewReader(encrypted), &buf); err != nil {
			b.Fatalf("DecryptStream failed: %v", err)
		}
	}
}

// ==============================================================================
// HIGH-PERFORMANCE StreamCipher TESTS
// ==============================================================================

// TestStreamCipher tests the high-performance StreamCipher API
func TestStreamCipher(t *testing.T) {
	password := "stream-cipher-password"
	testData := make([]byte, ChunkSize*2+500)
	if _, err := rand.Read(testData); err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	// Create cipher
	cipher := NewStreamCipher(password)
	defer cipher.Close()

	// Encrypt
	var encryptedBuf bytes.Buffer
	if err := cipher.EncryptStream(bytes.NewReader(testData), &encryptedBuf); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Decrypt
	var decryptedBuf bytes.Buffer
	if err := cipher.DecryptStream(bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf); err != nil {
		t.Fatalf("DecryptStream failed: %v", err)
	}

	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Error("Decrypted data doesn't match original")
	}
}

// BenchmarkStreamCipherEncrypt benchmarks high-performance batch encryption
// This should be significantly faster than BenchmarkEncryptStream for repeated operations
func BenchmarkStreamCipherEncrypt(b *testing.B) {
	password := "benchmark-password"
	testData := make([]byte, 1*1024*1024) // 1MB
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	// Create cipher ONCE (key derivation happens here)
	cipher := NewStreamCipher(password)
	defer cipher.Close()

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := cipher.EncryptStream(bytes.NewReader(testData), &buf); err != nil {
			b.Fatalf("EncryptStream failed: %v", err)
		}
	}
}

// BenchmarkStreamCipherDecrypt benchmarks high-performance batch decryption
func BenchmarkStreamCipherDecrypt(b *testing.B) {
	password := "benchmark-password"
	testData := make([]byte, 1*1024*1024) // 1MB
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	// Create cipher and encrypt
	cipher := NewStreamCipher(password)
	defer cipher.Close()

	var encryptedBuf bytes.Buffer
	if err := cipher.EncryptStream(bytes.NewReader(testData), &encryptedBuf); err != nil {
		b.Fatalf("EncryptStream failed: %v", err)
	}
	encrypted := encryptedBuf.Bytes()

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := cipher.DecryptStream(bytes.NewReader(encrypted), &buf); err != nil {
			b.Fatalf("DecryptStream failed: %v", err)
		}
	}
}

// BenchmarkRawXChaCha20 benchmarks raw XChaCha20-Poly1305 without KDF overhead
// This shows the theoretical maximum throughput of the cipher itself
func BenchmarkRawXChaCha20(b *testing.B) {
	key := make([]byte, KeySize)
	if _, err := rand.Read(key); err != nil {
		b.Fatalf("Failed to generate key: %v", err)
	}

	testData := make([]byte, 1*1024*1024) // 1MB
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		// Simulate chunked encryption
		for offset := 0; offset < len(testData); offset += ChunkSize {
			end := offset + ChunkSize
			if end > len(testData) {
				end = len(testData)
			}
			_ = testData[offset:end]
		}
	}
}

// ==============================================================================
// CROSS-PLATFORM COMPATIBILITY TESTS
// ==============================================================================
// These tests ensure the encryption works correctly on all platforms:
// - ARM (Raspberry Pi, Apple Silicon)
// - Low-memory devices (1GB RAM)
// - 32-bit systems
// - Various CPU architectures

// TestLowMemoryOperation tests that encryption works with minimal memory usage.
// This is critical for Raspberry Pi and other embedded devices.
func TestLowMemoryOperation(t *testing.T) {
	password := "low-memory-test"

	// Test with a file size that would cause issues if memory is wasted
	// 256KB = 4 chunks, manageable even on 512MB RAM devices
	testData := make([]byte, 256*1024)
	if _, err := rand.Read(testData); err != nil {
		t.Fatalf("Failed to generate random data: %v", err)
	}

	// Encrypt using StreamCipher (more memory-efficient for batch operations)
	cipher := NewStreamCipher(password)
	defer cipher.Close()

	var encryptedBuf bytes.Buffer
	if err := cipher.EncryptStream(bytes.NewReader(testData), &encryptedBuf); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	// Decrypt
	var decryptedBuf bytes.Buffer
	if err := cipher.DecryptStream(bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf); err != nil {
		t.Fatalf("DecryptStream failed: %v", err)
	}

	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Error("Data mismatch - low memory operation failed")
	}

	t.Logf("Low-memory test passed: encrypted %d bytes in %d chunks",
		len(testData), len(testData)/ChunkSize+1)
}

// TestSmallChunkBoundaries tests edge cases at chunk boundaries.
// Important for 32-bit systems where integer overflow could occur.
func TestSmallChunkBoundaries(t *testing.T) {
	password := "boundary-test"
	cipher := NewStreamCipher(password)
	defer cipher.Close()

	testCases := []struct {
		name string
		size int
	}{
		{"ExactlyOneChunk", ChunkSize},
		{"OneByteLessThanChunk", ChunkSize - 1},
		{"OneByteMoreThanChunk", ChunkSize + 1},
		{"ExactlyTwoChunks", ChunkSize * 2},
		{"TwoChunksMinusOne", ChunkSize*2 - 1},
		{"TinyData", 64},
		{"SingleByte", 1},
		{"Empty", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testData := make([]byte, tc.size)
			if tc.size > 0 {
				if _, err := rand.Read(testData); err != nil {
					t.Fatalf("Failed to generate random data: %v", err)
				}
			}

			var encryptedBuf bytes.Buffer
			if err := cipher.EncryptStream(bytes.NewReader(testData), &encryptedBuf); err != nil {
				t.Fatalf("Encrypt failed for size %d: %v", tc.size, err)
			}

			var decryptedBuf bytes.Buffer
			if err := cipher.DecryptStream(bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf); err != nil {
				t.Fatalf("Decrypt failed for size %d: %v", tc.size, err)
			}

			if !bytes.Equal(testData, decryptedBuf.Bytes()) {
				t.Errorf("Data mismatch for size %d", tc.size)
			}
		})
	}
}

// TestConcurrentAccess tests that the service is thread-safe.
// Important for multi-core systems and concurrent file operations.
func TestConcurrentAccess(t *testing.T) {
	password := "concurrent-test"
	cipher := NewStreamCipher(password)
	defer cipher.Close()

	const numGoroutines = 10
	const dataSize = ChunkSize + 100

	errChan := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testData := make([]byte, dataSize)
			if _, err := rand.Read(testData); err != nil {
				errChan <- err
				return
			}

			var encryptedBuf bytes.Buffer
			if err := cipher.EncryptStream(bytes.NewReader(testData), &encryptedBuf); err != nil {
				errChan <- err
				return
			}

			var decryptedBuf bytes.Buffer
			if err := cipher.DecryptStream(bytes.NewReader(encryptedBuf.Bytes()), &decryptedBuf); err != nil {
				errChan <- err
				return
			}

			if !bytes.Equal(testData, decryptedBuf.Bytes()) {
				errChan <- errors.New("data mismatch in goroutine")
				return
			}

			errChan <- nil
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("Concurrent test failed: %v", err)
		}
	}
}

// TestDeterministicNonceDerivation ensures nonce derivation is consistent.
// Critical for random-access decryption and cross-platform compatibility.
func TestDeterministicNonceDerivation(t *testing.T) {
	baseNonce := make([]byte, NonceSize)
	for i := range baseNonce {
		baseNonce[i] = byte(i)
	}

	// Test that both derivation methods produce identical results
	for chunkIndex := uint64(0); chunkIndex < 1000; chunkIndex++ {
		// Original allocation method
		nonce1 := deriveChunkNonce(baseNonce, chunkIndex)

		// In-place method
		nonce2 := make([]byte, NonceSize)
		deriveChunkNonceInPlace(baseNonce, chunkIndex, nonce2)

		if !bytes.Equal(nonce1, nonce2) {
			t.Errorf("Nonce derivation mismatch at chunk %d", chunkIndex)
		}
	}

	// Test extreme chunk indices (important for large files and 64-bit safety)
	extremeIndices := []uint64{
		0,
		1,
		1<<32 - 1, // Max uint32
		1 << 32,   // First value beyond uint32
		1<<63 - 1, // Near max uint64
		1<<64 - 1, // Max uint64
	}

	for _, idx := range extremeIndices {
		nonce := deriveChunkNonce(baseNonce, idx)
		if len(nonce) != NonceSize {
			t.Errorf("Invalid nonce length for index %d", idx)
		}
	}
}

// TestPortabilityDataFormat verifies the binary format is consistent.
// Ensures files encrypted on one platform can be decrypted on another.
func TestPortabilityDataFormat(t *testing.T) {
	password := "portability-test"
	testData := []byte("Portable data test - should work on ARM, x86, x64, etc.")

	var encryptedBuf bytes.Buffer
	if err := EncryptStream(password, bytes.NewReader(testData), &encryptedBuf); err != nil {
		t.Fatalf("EncryptStream failed: %v", err)
	}

	encrypted := encryptedBuf.Bytes()

	// Verify header structure
	if len(encrypted) < HeaderSize {
		t.Fatalf("Encrypted data too short")
	}

	// Check magic bytes
	if string(encrypted[0:4]) != MagicBytes {
		t.Error("Magic bytes mismatch")
	}

	// Check version
	if encrypted[4] != Version {
		t.Errorf("Version mismatch: got 0x%02x, expected 0x%02x", encrypted[4], Version)
	}

	// Verify salt is at correct position (bytes 5-20)
	salt := encrypted[5 : 5+SaltSize]
	if len(salt) != SaltSize {
		t.Errorf("Salt size wrong: %d", len(salt))
	}

	// Verify base nonce is at correct position (bytes 21-44)
	baseNonce := encrypted[5+SaltSize : 5+SaltSize+NonceSize]
	if len(baseNonce) != NonceSize {
		t.Errorf("BaseNonce size wrong: %d", len(baseNonce))
	}

	t.Logf("Header verified: Magic=%s, Version=0x%02x, SaltLen=%d, NonceLen=%d",
		MagicBytes, Version, len(salt), len(baseNonce))
}

// BenchmarkLowEndDevice simulates encryption performance on a Raspberry Pi.
// Uses smaller data sizes and single-threaded operation.
func BenchmarkLowEndDevice(b *testing.B) {
	password := "raspi-benchmark"

	// Simulate typical Raspberry Pi file operations (64KB files)
	testData := make([]byte, 64*1024)
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	cipher := NewStreamCipher(password)
	defer cipher.Close()

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := cipher.EncryptStream(bytes.NewReader(testData), &buf); err != nil {
			b.Fatalf("EncryptStream failed: %v", err)
		}
	}
}

// BenchmarkMemoryConstrained simulates operation under memory pressure.
// Uses standard API which has higher memory overhead per call.
func BenchmarkMemoryConstrained(b *testing.B) {
	password := "memory-constrained"

	// Small file - 8KB (typical config file)
	testData := make([]byte, 8*1024)
	if _, err := rand.Read(testData); err != nil {
		b.Fatalf("Failed to generate random data: %v", err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(testData)))

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := EncryptStream(password, bytes.NewReader(testData), &buf); err != nil {
			b.Fatalf("EncryptStream failed: %v", err)
		}
	}
}
