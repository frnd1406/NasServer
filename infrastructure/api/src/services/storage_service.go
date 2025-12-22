package services

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var ErrPathTraversal = errors.New("path escapes base directory")
var ErrInvalidFileType = errors.New("file type not allowed")
var ErrFileTooLarge = errors.New("file exceeds maximum size")

// Security constants
const MaxUploadSize = 100 * 1024 * 1024 // 100 MB

// AllowedMimeTypes defines the whitelist of permitted file types
var AllowedMimeTypes = map[string]bool{
	// Images
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,

	// Documents
	"application/pdf": true,
	"text/plain":      true,
	"text/csv":        true,
	"text/markdown":   true,

	// Archives (careful - could contain malware, but needed for backups)
	"application/zip":              true,
	"application/x-zip-compressed": true,
	"application/gzip":             true,
	"application/x-gzip":           true,
	"application/x-tar":            true,

	// Video
	"video/mp4":  true,
	"video/mpeg": true,
	"video/webm": true,

	// Audio
	"audio/mpeg": true,
	"audio/mp3":  true,
	"audio/wav":  true,
	"audio/ogg":  true,
}

// Magic number signatures for common file types (first 16 bytes)
var magicNumbers = map[string][]byte{
	"image/jpeg":      {0xFF, 0xD8, 0xFF},
	"image/png":       {0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	"image/gif":       {0x47, 0x49, 0x46, 0x38},
	"application/pdf": {0x25, 0x50, 0x44, 0x46},
	"application/zip": {0x50, 0x4B, 0x03, 0x04},
	"video/mp4":       {0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70}, // ftyp box
}

// StorageEntry represents a file or directory item within the storage root.
type StorageEntry struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"isDir"`
	ModTime  time.Time `json:"modTime"`
	MimeType string    `json:"mimeType,omitempty"`
	IsImage  bool      `json:"isImage,omitempty"`
}

// TrashEntry represents a soft-deleted file.
type TrashEntry struct {
	ID           string    `json:"id"`           // bucket/relative/path
	Name         string    `json:"name"`         // file name
	OriginalPath string    `json:"originalPath"` // original relative path
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"modTime"`
}

// StorageService provides basic file operations within a confined base directory.
type StorageService struct {
	basePath  string
	trashPath string
	logger    *logrus.Logger
}

// NewStorageService initializes the service and ensures the base path exists.
func NewStorageService(basePath string, logger *logrus.Logger) (*StorageService, error) {
	if basePath == "" {
		return nil, fmt.Errorf("base path is required")
	}

	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("resolve base path: %w", err)
	}

	if err := os.MkdirAll(absBase, 0o755); err != nil {
		return nil, fmt.Errorf("ensure base path: %w", err)
	}

	trash := filepath.Join(absBase, ".trash")
	if err := os.MkdirAll(trash, 0o755); err != nil {
		return nil, fmt.Errorf("ensure trash path: %w", err)
	}

	return &StorageService{
		basePath:  absBase,
		trashPath: trash,
		logger:    logger,
	}, nil
}

func (s *StorageService) sanitizePath(rel string) (string, error) {
	if strings.Contains(rel, "..") {
		return "", ErrPathTraversal
	}
	// Prepend slash so Clean treats it as absolute, then trim to avoid breaking out.
	cleaned := filepath.Clean("/" + rel)
	trimmed := strings.TrimPrefix(cleaned, "/")
	full := filepath.Join(s.basePath, trimmed)

	abs, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}

	if abs != s.basePath && !strings.HasPrefix(abs, s.basePath+string(os.PathSeparator)) {
		return "", ErrPathTraversal
	}

	return abs, nil
}

// List returns the entries for the given relative path.
func (s *StorageService) List(relPath string) ([]StorageEntry, error) {
	target, err := s.sanitizePath(relPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}

	var items []StorageEntry
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			s.logger.WithError(err).Warn("storage: failed to read entry info")
			continue
		}

		mimeType := ""
		isImage := false
		if !info.IsDir() {
			mimeType = mime.TypeByExtension(filepath.Ext(e.Name()))
			if strings.HasPrefix(mimeType, "image/") {
				isImage = true
			}
		}

		items = append(items, StorageEntry{
			Name:     e.Name(),
			Size:     info.Size(),
			IsDir:    info.IsDir(),
			ModTime:  info.ModTime(),
			MimeType: mimeType,
			IsImage:  isImage,
		})
	}

	return items, nil
}

// ValidateFileType checks if the file type is allowed based on magic numbers and MIME type
func (s *StorageService) ValidateFileType(file multipart.File, filename string) error {
	// Read first 512 bytes for magic number detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Detect MIME type from magic numbers (more reliable than extension)
	detectedType := http.DetectContentType(buffer[:n])

	// Strip charset parameter if present (e.g., "text/plain; charset=utf-8" -> "text/plain")
	if idx := strings.Index(detectedType, ";"); idx != -1 {
		detectedType = strings.TrimSpace(detectedType[:idx])
	}

	// Log detection for debugging
	s.logger.WithFields(logrus.Fields{
		"filename":      filename,
		"detected_type": detectedType,
		"bytes_read":    n,
	}).Debug("File type detection")

	// Check against whitelist
	// EXCEPTION: .enc files are encrypted and have random bytes (no magic number)
	if strings.ToLower(filepath.Ext(filename)) == ".enc" {
		s.logger.WithField("filename", filename).Debug("Allowing .enc file without magic number check")
		return nil
	}

	if !AllowedMimeTypes[detectedType] {
		// Special case: some files might be detected as octet-stream
		// Check magic number signatures manually
		if detectedType == "application/octet-stream" {
			for mimeType, magic := range magicNumbers {
				if len(buffer) >= len(magic) && bytesMatch(buffer[:len(magic)], magic) {
					detectedType = mimeType
					break
				}
			}
		}

		// Final check after magic number verification
		if !AllowedMimeTypes[detectedType] {
			s.logger.WithFields(logrus.Fields{
				"filename":      filename,
				"detected_type": detectedType,
			}).Warn("File type not allowed")
			return fmt.Errorf("%w: %s (detected as %s)", ErrInvalidFileType, filename, detectedType)
		}
	}

	// Additional check: Reject executable extensions even if MIME type passes
	ext := strings.ToLower(filepath.Ext(filename))
	dangerousExtensions := []string{
		".exe", ".bat", ".cmd", ".com", ".pif", ".scr", ".vbs", ".js", ".jar",
		".sh", ".bash", ".zsh", ".fish", ".ps1", ".app", ".deb", ".rpm",
		".php", ".jsp", ".asp", ".aspx", ".cgi", ".pl", ".py", ".rb",
	}

	for _, dangerous := range dangerousExtensions {
		if ext == dangerous {
			s.logger.WithFields(logrus.Fields{
				"filename":  filename,
				"extension": ext,
			}).Warn("Dangerous file extension blocked")
			return fmt.Errorf("%w: executable or script file extension not allowed (%s)", ErrInvalidFileType, ext)
		}
	}

	return nil
}

// bytesMatch compares two byte slices
func bytesMatch(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ValidateFileSize checks if the file size is within allowed limits
func (s *StorageService) ValidateFileSize(file multipart.File, fileHeader *multipart.FileHeader) error {
	// Get file size from FileHeader
	size := fileHeader.Size

	if size > MaxUploadSize {
		s.logger.WithFields(logrus.Fields{
			"size":     size,
			"max_size": MaxUploadSize,
		}).Warn("File exceeds maximum upload size")
		return fmt.Errorf("%w: file size %d bytes exceeds maximum of %d bytes", ErrFileTooLarge, size, MaxUploadSize)
	}

	return nil
}

// SaveResult contains metadata about the saved file
type SaveResult struct {
	Path     string
	MimeType string
	FileID   string
}

// Save stores the provided file into the given relative directory.
// If the file already exists, automatic versioning creates backup copies (.v1.bak, .v2.bak, .v3.bak)
func (s *StorageService) Save(dir string, file multipart.File, fileHeader *multipart.FileHeader) (*SaveResult, error) {
	filename := fileHeader.Filename

	if filename == "" {
		return nil, fmt.Errorf("filename is required")
	}

	// SECURITY: Validate file size FIRST (before reading content)
	if err := s.ValidateFileSize(file, fileHeader); err != nil {
		return nil, err
	}

	// Read first 512 bytes for MIME detection (before validation)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Detect MIME type
	detectedMimeType := http.DetectContentType(buffer[:n])

	// SECURITY: Validate file type (magic numbers + extension check)
	if err := s.ValidateFileType(file, filename); err != nil {
		return nil, err
	}

	targetDir, err := s.sanitizePath(dir)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return nil, fmt.Errorf("create target dir: %w", err)
	}

	destPath, err := s.sanitizePath(filepath.Join(dir, filepath.Base(filename)))
	if err != nil {
		return nil, err
	}

	// FILE VERSIONING: If file exists, rotate versions before overwriting
	if _, err := os.Stat(destPath); err == nil {
		s.rotateVersions(destPath, 3) // Keep max 3 backup versions
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"filename": filename,
		"path":     destPath,
	}).Info("File uploaded successfully")

	// Create result with metadata
	result := &SaveResult{
		Path:     destPath,
		MimeType: detectedMimeType,
		FileID:   filepath.Base(filename), // Use filename as ID
	}

	return result, nil
}

// rotateVersions implements "Time Machine Light" - keeps last N versions of a file
// Before overwriting a file, this function:
// 1. Deletes oldest version if > maxVersions exist
// 2. Renames existing versions: .v2.bak -> .v3.bak, .v1.bak -> .v2.bak
// 3. Renames current file to .v1.bak
func (s *StorageService) rotateVersions(filePath string, maxVersions int) {
	// Delete oldest version if it exceeds max
	oldestPath := fmt.Sprintf("%s.v%d.bak", filePath, maxVersions)
	if _, err := os.Stat(oldestPath); err == nil {
		if err := os.Remove(oldestPath); err != nil {
			s.logger.WithError(err).Warn("Failed to remove oldest version")
		}
	}

	// Shift existing versions: v2 -> v3, v1 -> v2
	for i := maxVersions - 1; i >= 1; i-- {
		oldPath := fmt.Sprintf("%s.v%d.bak", filePath, i)
		newPath := fmt.Sprintf("%s.v%d.bak", filePath, i+1)
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				s.logger.WithError(err).WithFields(logrus.Fields{
					"from": oldPath,
					"to":   newPath,
				}).Warn("Failed to rotate version")
			}
		}
	}

	// Move current file to .v1.bak
	v1Path := fmt.Sprintf("%s.v1.bak", filePath)
	if err := os.Rename(filePath, v1Path); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"from": filePath,
			"to":   v1Path,
		}).Warn("Failed to create backup version")
	} else {
		s.logger.WithFields(logrus.Fields{
			"file":   filepath.Base(filePath),
			"backup": filepath.Base(v1Path),
		}).Info("Created backup version before overwrite")
	}
}

// Open returns a file handle and metadata for download.
func (s *StorageService) Open(relPath string) (*os.File, os.FileInfo, string, error) {
	target, err := s.sanitizePath(relPath)
	if err != nil {
		return nil, nil, "", err
	}

	info, err := os.Stat(target)
	if err != nil {
		return nil, nil, "", err
	}
	if info.IsDir() {
		return nil, nil, "", fmt.Errorf("cannot download a directory")
	}

	f, err := os.Open(target)
	if err != nil {
		return nil, nil, "", err
	}

	ctype := mime.TypeByExtension(filepath.Ext(info.Name()))
	if ctype == "" {
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		ctype = http.DetectContentType(buf[:n])
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			f.Close()
			return nil, nil, "", err
		}
	}

	return f, info, ctype, nil
}

// Delete removes a file or directory (recursively) within the storage root.
func (s *StorageService) Delete(relPath string) error {
	if relPath == "" || relPath == "/" || relPath == "." {
		return fmt.Errorf("refusing to delete storage root")
	}

	target, err := s.sanitizePath(relPath)
	if err != nil {
		return err
	}

	relClean := strings.TrimPrefix(filepath.Clean(relPath), string(os.PathSeparator))
	bucket := time.Now().UTC().Format("20060102T150405Z")
	destRel := filepath.Join(bucket, relClean)
	destAbs := filepath.Join(s.trashPath, destRel)

	if err := os.MkdirAll(filepath.Dir(destAbs), 0o755); err != nil {
		return fmt.Errorf("create trash dir: %w", err)
	}

	if err := os.Rename(target, destAbs); err != nil {
		return fmt.Errorf("move to trash: %w", err)
	}
	return nil
}

// ListTrash lists all soft-deleted files.
func (s *StorageService) ListTrash() ([]TrashEntry, error) {
	var entries []TrashEntry
	err := filepath.Walk(s.trashPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(s.trashPath, path)
		if err != nil {
			return nil
		}
		parts := strings.SplitN(rel, string(os.PathSeparator), 2)
		original := ""
		if len(parts) == 2 {
			original = parts[1]
		}
		entries = append(entries, TrashEntry{
			ID:           filepath.ToSlash(rel),
			Name:         info.Name(),
			OriginalPath: filepath.ToSlash(original),
			Size:         info.Size(),
			ModTime:      info.ModTime(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// RestoreFromTrash moves a file from trash back to its original location.
func (s *StorageService) RestoreFromTrash(id string) error {
	if id == "" {
		return fmt.Errorf("invalid trash id")
	}
	source := filepath.Join(s.trashPath, filepath.FromSlash(id))
	if _, err := os.Stat(source); err != nil {
		return err
	}

	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		return fmt.Errorf("missing original path info")
	}
	originalRel := parts[1]
	dest, err := s.sanitizePath(originalRel)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.Rename(source, dest)
}

// DeleteFromTrash removes a trashed file permanently.
func (s *StorageService) DeleteFromTrash(id string) error {
	if id == "" {
		return fmt.Errorf("invalid trash id")
	}
	target := filepath.Join(s.trashPath, filepath.FromSlash(id))
	if !strings.HasPrefix(target, s.trashPath) {
		return ErrPathTraversal
	}
	return os.RemoveAll(target)
}

// Rename renames a file within the same directory.
func (s *StorageService) Rename(oldRel, newName string) error {
	if newName == "" {
		return fmt.Errorf("new name required")
	}
	oldAbs, err := s.sanitizePath(oldRel)
	if err != nil {
		return err
	}
	info, err := os.Stat(oldAbs)
	if err != nil {
		return err
	}
	dir := filepath.Dir(oldAbs)
	newAbs := filepath.Join(dir, filepath.Base(newName))
	if newAbs == oldAbs {
		return nil
	}
	if info.IsDir() {
		return os.Rename(oldAbs, newAbs)
	}
	return os.Rename(oldAbs, newAbs)
}

// GetFullPath returns the full filesystem path for a relative path (with security checks)
func (s *StorageService) GetFullPath(relPath string) (string, error) {
	return s.sanitizePath(relPath)
}

// Mkdir creates a new directory at the given relative path
func (s *StorageService) Mkdir(relPath string) error {
	if relPath == "" || relPath == "/" {
		return fmt.Errorf("invalid directory path")
	}

	target, err := s.sanitizePath(relPath)
	if err != nil {
		return err
	}

	// Check if already exists
	if info, err := os.Stat(target); err == nil {
		if info.IsDir() {
			return nil // Already exists as directory
		}
		return fmt.Errorf("path exists but is not a directory")
	}

	return os.MkdirAll(target, 0o755)
}
