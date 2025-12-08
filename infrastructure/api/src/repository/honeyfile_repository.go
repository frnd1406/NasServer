package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Honeyfile represents a trap file for intrusion detection
type Honeyfile struct {
	ID              uuid.UUID    `db:"id"`
	FilePath        string       `db:"file_path"`
	FileType        string       `db:"file_type"` // 'finance', 'it', 'private', 'general'
	TriggerCount    int          `db:"trigger_count"`
	LastTriggeredAt sql.NullTime `db:"last_triggered_at"`
	CreatedAt       time.Time    `db:"created_at"`
	CreatedBy       *uuid.UUID   `db:"created_by"`
}

// HoneyfileRepository handles DB operations for honeyfiles
type HoneyfileRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

// NewHoneyfileRepository creates a new repository instance
func NewHoneyfileRepository(db *sqlx.DB, logger *logrus.Logger) *HoneyfileRepository {
	return &HoneyfileRepository{db: db, logger: logger}
}

// EnsureTable creates the honeyfiles table if it doesn't exist
func (r *HoneyfileRepository) EnsureTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS honeyfiles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			file_path VARCHAR(512) NOT NULL UNIQUE,
			file_type VARCHAR(50) NOT NULL DEFAULT 'general',
			trigger_count INT DEFAULT 0,
			last_triggered_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_by UUID REFERENCES users(id),
			CONSTRAINT honeyfiles_type_check CHECK (file_type IN ('finance', 'it', 'private', 'general'))
		)
	`
	if _, err := r.db.ExecContext(ctx, query); err != nil {
		r.logger.WithError(err).Error("failed to ensure honeyfiles table")
		return fmt.Errorf("ensure honeyfiles table: %w", err)
	}
	return nil
}

// IsHoneyfile checks if a path is marked as a honeyfile (uses canonical path)
func (r *HoneyfileRepository) IsHoneyfile(ctx context.Context, rawPath string) bool {
	cleanPath := filepath.Clean(rawPath)

	var exists bool
	err := r.db.GetContext(ctx, &exists,
		"SELECT EXISTS(SELECT 1 FROM honeyfiles WHERE file_path = $1)", cleanPath)
	if err != nil {
		r.logger.WithError(err).Warn("honeyfile check failed")
		return false
	}
	return exists
}

// ListAll returns all honeyfiles
func (r *HoneyfileRepository) ListAll(ctx context.Context) ([]Honeyfile, error) {
	var honeyfiles []Honeyfile
	err := r.db.SelectContext(ctx, &honeyfiles,
		"SELECT id, file_path, file_type, trigger_count, last_triggered_at, created_at, created_by FROM honeyfiles ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("list honeyfiles: %w", err)
	}
	return honeyfiles, nil
}

// Create adds a new honeyfile marker
func (r *HoneyfileRepository) Create(ctx context.Context, rawPath, fileType string, createdBy *uuid.UUID) (*Honeyfile, error) {
	cleanPath := filepath.Clean(rawPath)

	honeyfile := &Honeyfile{
		FilePath:  cleanPath,
		FileType:  fileType,
		CreatedBy: createdBy,
	}

	query := `
		INSERT INTO honeyfiles (file_path, file_type, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, cleanPath, fileType, createdBy).Scan(&honeyfile.ID, &honeyfile.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create honeyfile: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"id":        honeyfile.ID,
		"file_path": cleanPath,
		"file_type": fileType,
	}).Info("Honeyfile created")

	return honeyfile, nil
}

// Delete removes a honeyfile marker
func (r *HoneyfileRepository) Delete(ctx context.Context, rawPath string) error {
	cleanPath := filepath.Clean(rawPath)

	result, err := r.db.ExecContext(ctx, "DELETE FROM honeyfiles WHERE file_path = $1", cleanPath)
	if err != nil {
		return fmt.Errorf("delete honeyfile: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("honeyfile not found: %s", cleanPath)
	}

	r.logger.WithField("file_path", cleanPath).Info("Honeyfile deleted")
	return nil
}

// IncrementTrigger updates the trigger count and timestamp
func (r *HoneyfileRepository) IncrementTrigger(ctx context.Context, rawPath string) error {
	cleanPath := filepath.Clean(rawPath)

	_, err := r.db.ExecContext(ctx, `
		UPDATE honeyfiles 
		SET trigger_count = trigger_count + 1, 
		    last_triggered_at = NOW() 
		WHERE file_path = $1
	`, cleanPath)

	if err != nil {
		return fmt.Errorf("increment trigger: %w", err)
	}

	r.logger.WithField("file_path", cleanPath).Warn("🚨 HONEYFILE TRIGGER RECORDED")
	return nil
}

// GetAllPaths returns just the paths for cache loading
func (r *HoneyfileRepository) GetAllPaths(ctx context.Context) ([]string, error) {
	var paths []string
	err := r.db.SelectContext(ctx, &paths, "SELECT file_path FROM honeyfiles")
	if err != nil {
		return nil, fmt.Errorf("get all paths: %w", err)
	}
	return paths, nil
}
