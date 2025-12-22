//go:build examples
// +build examples

package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/nas-ai/api/src/logger"
)

// Example demonstrating slog structured logging
func main() {
	// Create production logger (JSON output)
	log := logger.NewSlogLogger("info", "production")

	// Basic logging
	log.Info("Application started",
		slog.String("version", "2.0.0"),
		slog.String("environment", "production"),
	)

	// Log with multiple attributes
	log.Info("User registered",
		slog.String("user_id", "user_123abc"),
		slog.String("email", "user@example.com"),
		slog.Time("timestamp", time.Now()),
		slog.Bool("email_verified", false),
	)

	// Log with error
	err := performOperation()
	if err != nil {
		log.Error("Operation failed",
			slog.String("error", err.Error()),
			slog.String("operation", "send_email"),
			slog.Duration("retry_after", 5*time.Minute),
		)
	}

	// Log with grouping
	log.Info("Database query executed",
		slog.Group("query",
			slog.String("sql", "SELECT * FROM users WHERE id = $1"),
			slog.Int("rows_affected", 1),
			slog.Duration("duration", 25*time.Millisecond),
		),
		slog.Group("connection",
			slog.String("host", "localhost"),
			slog.Int("port", 5433),
			slog.String("database", "nas_db"),
		),
	)

	// Context-aware logging
	ctx := context.Background()
	ctxLog := logger.ContextWithLogger(ctx, log)
	logFromContext := logger.LoggerFromContext(ctxLog)
	logFromContext.Info("Processing request with context")

	// Different log levels
	log.Debug("Debug information", slog.Any("data", map[string]int{"count": 42}))
	log.Warn("Warning message", slog.String("reason", "rate_limit_approaching"))
	log.Error("Error occurred", slog.String("error", "connection refused"))
}

func performOperation() error {
	// Simulate operation
	return nil
}
