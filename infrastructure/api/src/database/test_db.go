package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// NewTestDatabase creates an in-memory SQLite database for testing
// This allows testing real SQL queries without needing a full Postgres instance
func NewTestDatabase(logger *logrus.Logger) (*DB, error) {
	// Open in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Create users table schema (matching Postgres schema)
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		email_verified BOOLEAN NOT NULL DEFAULT FALSE,
		verified_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create test schema: %w", err)
	}

	logger.Debug("✅ Test database (SQLite in-memory) initialized")

	return &DB{
		DB:     db,
		logger: logger,
	}, nil
}
