package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the API server
type Config struct {
	// Server
	Port string

	// CORS (Whitelist - NO WILDCARD!)
	CORSOrigins []string

	// Rate Limiting
	RateLimitPerMin int

	// Logging
	LogLevel string

	// JWT (Phase 2 - but validate now!)
	JWTSecret     string
	JWTSecretFile string

	// Database (Phase 2)
	DatabaseURL  string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string

	// Redis (Phase 2)
	RedisURL  string
	RedisHost string
	RedisPort string

	// Email (Phase 3 - Resend)
	ResendAPIKey string
	EmailFrom    string
	FrontendURL  string

	// Cloudflare (Phase 3)
	CloudflareAPIToken string
	CloudflareR2Bucket string

	// Environment
	Environment string

	// Monitoring (agent ingestion)
	MonitoringToken string

	// AI/semantic search
	AIServiceURL string

	// LLM for RAG
	OllamaURL string
	LLMModel  string

	// Internal Security (Shared Secret)
	InternalAPISecret string

	// Abuse Prevention
	InviteCode string

	// Backup configuration
	BackupSchedule       string
	BackupRetentionCount int
	BackupStoragePath    string

	// Consistency check (orphan cleanup)
	ConsistencyCheckIntervalMin int
}

// LoadConfig loads configuration using Viper (supports .env, config.yaml, and env vars)
// CRITICAL: Fails fast if required secrets are missing!
func LoadConfig() (*Config, error) {
	return LoadConfigWithViper()
}

// LoadConfigFromEnv is the legacy configuration loader (kept for backward compatibility)
// Use LoadConfig() instead, which now uses Viper
func LoadConfigFromEnv() (*Config, error) {
	cfg := &Config{
		// Defaults
		Port:                        getEnv("PORT", "8080"),
		LogLevel:                    getEnv("LOG_LEVEL", "info"),
		Environment:                 getEnv("ENV", "development"),
		RateLimitPerMin:             getEnvInt("RATE_LIMIT_PER_MIN", 100),
		BackupSchedule:              getEnv("BACKUP_SCHEDULE", "0 3 * * *"),
		BackupRetentionCount:        getEnvInt("BACKUP_RETENTION_COUNT", 7),
		BackupStoragePath:           getEnv("BACKUP_STORAGE_PATH", "/mnt/backups"),
		ConsistencyCheckIntervalMin: getEnvInt("CONSISTENCY_CHECK_INTERVAL_MIN", 5),
		AIServiceURL:                getEnv("AI_SERVICE_URL", "http://ai-knowledge-agent:5000"),
		OllamaURL:                   getEnv("OLLAMA_URL", "http://localhost:11434"),
		LLMModel:                    getEnv("LLM_MODEL", "qwen2.5:3b"),
	}

	// CORS Origins (Whitelist)
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:5173")
	cfg.CORSOrigins = strings.Split(corsOrigins, ",")
	for i := range cfg.CORSOrigins {
		cfg.CORSOrigins[i] = strings.TrimSpace(cfg.CORSOrigins[i])
	}

	// JWT Secret - REQUIRED (even if not used in Phase 1)
	// Fail-fast principle: Better fail now than at runtime!
	if secretFile := os.Getenv("JWT_SECRET_FILE"); secretFile != "" {
		secret, err := readSecretFromFile(secretFile)
		if err != nil {
			return nil, err
		}
		cfg.JWTSecret = secret
		cfg.JWTSecretFile = secretFile
	} else {
		cfg.JWTSecret = strings.TrimSpace(os.Getenv("JWT_SECRET"))
		if cfg.JWTSecret == "" {
			return nil, fmt.Errorf("CRITICAL: JWT_SECRET environment variable is required (no defaults allowed)")
		}
	}

	// Validate JWT secret strength (min 32 chars)
	if err := ValidateJWTSecret(cfg.JWTSecret); err != nil {
		return nil, err
	}

	// Database Configuration (Phase 2)
	// Support both DATABASE_URL (single string) or individual components
	cfg.DatabaseURL = getEnv("DATABASE_URL", "")
	if cfg.DatabaseURL == "" {
		// Build from components (for docker-compose dev)
		cfg.DatabaseHost = getEnv("DB_HOST", "localhost")
		cfg.DatabasePort = getEnv("DB_PORT", "5433")
		cfg.DatabaseUser = getEnv("DB_USER", "nas_user")
		cfg.DatabasePass = getEnv("DB_PASSWORD", "nas_dev_password")
		cfg.DatabaseName = getEnv("DB_NAME", "nas_db")
		cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
	}

	// Redis Configuration (Phase 2)
	cfg.RedisURL = getEnv("REDIS_URL", "")
	if cfg.RedisURL == "" {
		cfg.RedisHost = getEnv("REDIS_HOST", "localhost")
		cfg.RedisPort = getEnv("REDIS_PORT", "6380")
		cfg.RedisURL = fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	}

	// Email Configuration (Phase 3)
	cfg.ResendAPIKey = getEnv("RESEND_API_KEY", "")
	cfg.EmailFrom = getEnv("EMAIL_FROM", "NAS.AI <noreply@your-domain.com>")
	cfg.FrontendURL = getEnv("FRONTEND_URL", "https://your-domain.com")

	// Cloudflare Configuration (Phase 3)
	cfg.CloudflareAPIToken = getEnv("CLOUDFLARE_API_TOKEN", "")
	cfg.CloudflareR2Bucket = getEnv("CLOUDFLARE_R2_BUCKET", "nas-ai-storage")

	// Monitoring
	cfg.MonitoringToken = strings.TrimSpace(getEnv("MONITORING_TOKEN", ""))
	if len(cfg.MonitoringToken) < 16 {
		return nil, fmt.Errorf("CRITICAL: MONITORING_TOKEN must be at least 16 characters")
	}

	if strings.TrimSpace(cfg.BackupSchedule) == "" {
		return nil, fmt.Errorf("CRITICAL: BACKUP_SCHEDULE is required")
	}
	if cfg.BackupRetentionCount < 1 {
		return nil, fmt.Errorf("CRITICAL: BACKUP_RETENTION_COUNT must be >= 1")
	}
	if strings.TrimSpace(cfg.BackupStoragePath) == "" {
		return nil, fmt.Errorf("CRITICAL: BACKUP_STORAGE_PATH is required")
	}

	return cfg, nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvInt gets environment variable as int with fallback
func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return fallback
}
