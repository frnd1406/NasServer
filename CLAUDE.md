# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**NAS.AI** is a secure, AI-powered Network Attached Storage system with semantic search, zero-knowledge encryption, and a glassmorphism UI. The system is built as a microservices architecture using Docker Compose.

**Current Phase:** 2.3 - Encryption Implementation
**Version:** 2.1

## Quick Start

```bash
# Start production stack
cd infrastructure
docker compose -f docker-compose.prod.yml up -d

# Start development stack
docker compose -f docker-compose.dev.yml up -d

# Access services
# WebUI: http://localhost:80
# API:   http://localhost:8080
# AI:    http://localhost:5000
```

## Development Commands

### API Service (Go)

```bash
cd infrastructure/api

# Build
make build

# Run tests
make test

# Run tests with coverage (requires 80% minimum)
make test-coverage

# Run security scans (gosec + gitleaks)
make security-scan

# Run linter
make lint

# Run (requires JWT_SECRET environment variable)
export JWT_SECRET=$(openssl rand -base64 32)
make run
```

### WebUI (React + Vite)

```bash
cd infrastructure/webui

# Install dependencies
npm install

# Development server
npm run dev

# Production build
npm run build

# Preview production build
npm run preview
```

### AI Knowledge Agent (Python)

```bash
cd infrastructure/ai_knowledge_agent

# Install dependencies
pip install -r requirements.txt

# Run service (uses Ollama for embeddings and RAG)
python app.py
```

### Orchestrator (Health Monitoring)

```bash
cd orchestrator

# Build
make build

# Run in development mode
make dev

# Run in production mode
make prod

# Run tests
make test
```

### Unified CLI Tool

The `scripts/nas-cli.sh` provides a menu-driven interface for common operations:

```bash
cd scripts
./nas-cli.sh

# Options include:
# - API health checks and monitoring
# - Endpoint testing (with auth)
# - Docker deployment
# - Git savepoints
# - Login and token management
```

## Architecture

### Service Layer

The system follows a microservices architecture with clear separation of concerns:

```
┌─────────────┐
│   WebUI     │ (React + Vite + TailwindCSS)
│   Port 80   │
└──────┬──────┘
       │ REST/JSON
       ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│ API Gateway │────→│ PostgreSQL   │     │ Redis Cache │
│ Go :8080    │     │ + pgvector   │     │ :6379       │
└──────┬──────┘     └──────────────┘     └─────────────┘
       │
       ├─→ AI Knowledge Agent (Python :5000)
       │   • Ollama-based embeddings
       │   • RAG with Qwen2.5 LLM
       │   • pgvector for semantic search
       │
       └─→ Orchestrator (Go :9000)
           • Health monitoring
           • Prometheus metrics
```

### Go Backend Architecture Pattern

The API follows a layered architecture:

```
Handler Layer (handlers/)
    ↓ validates requests, marshals responses
Service Layer (services/)
    ↓ business logic, encryption, external services
Repository Layer (repository/)
    ↓ database access, queries
Database (PostgreSQL + Redis)
```

**Key Pattern:**
- **Handlers** (`handlers/*.go`) - HTTP request handling, validation, response formatting
- **Services** (`services/*.go`) - Business logic, encryption, file operations, email
- **Repository** (`repository/*.go`) - Database operations, SQL queries
- **Middleware** (`middleware/*.go`) - Auth, CSRF, rate limiting, vault guard

### Frontend Architecture

The WebUI follows a component-based architecture:

```
src/
├── pages/           # Route-level components
├── components/      # Reusable UI components
├── hooks/           # Custom React hooks (useFileStorage, useFileSelection)
├── context/         # React Context providers
├── utils/           # Pure utility functions
└── lib/             # Third-party integrations
```

**Key Patterns:**
- Custom hooks for business logic (e.g., `useFileStorage`, `useFilePreview`)
- GlassCard component for consistent glassmorphism UI
- Context providers for global state (auth, theme)

## Key Architectural Decisions

### Zero-Knowledge Encryption

The system implements AES-256-GCM encryption with Argon2id key derivation:

- **Vault Path:** `/tmp/nas-vault-demo` (non-persistent by default for maximum security)
- **Encrypted Storage:** `/media/frnd14/DEMO` (optional encrypted file storage)
- **Key Hierarchy:** Master Password → KEK → DEK → File/DB/Backup Keys
- **Startup State:** System starts in LOCKED state, requires unlock via master password

Files stored in encrypted storage are only readable through the WebUI when the vault is unlocked.

**Important:** The encryption service (`services/encryption_service.go`) keeps the Data Encryption Key (DEK) only in RAM. On container restart, the vault must be unlocked again.

### AI Knowledge Layer

- **Embeddings:** Uses Ollama with local models (no cloud dependencies)
- **Vector Storage:** pgvector extension in PostgreSQL (1024 dimensions)
- **RAG:** Qwen2.5 LLM for intelligent question answering
- **Secure Feeding:** `SecureAIFeeder` service sends decrypted content to AI without writing plaintext to disk

### Security Features

- **Authentication:** JWT access tokens (15min) + refresh tokens (7 days)
- **CSRF Protection:** Double-submit cookie pattern
- **Rate Limiting:** 100 requests/min per user via Redis
- **Path Traversal Protection:** All file operations validated
- **Security Headers:** CSP, X-Frame-Options, HSTS, X-Content-Type-Options

## Testing

### Running Tests

```bash
# Go API tests
cd infrastructure/api
make test                 # All tests
make test-coverage        # With coverage report
make test-security        # Security tests only

# Orchestrator tests
cd orchestrator
make test
```

### Test Files Location

- API tests: `infrastructure/api/src/services/*_test.go`, `infrastructure/api/src/handlers/*_test.go`
- Orchestrator tests: `orchestrator/orchestrator_test.go`

### API Endpoint Testing

Use the CLI tool for comprehensive endpoint testing:

```bash
./scripts/nas-cli.sh
# Select option for endpoint tests
# Can test with/without authentication
```

## Docker Services

| Service | Port | Purpose |
|---------|------|---------|
| postgres | 5432 | Primary database with pgvector extension |
| redis | 6379 | Session cache, rate limiting |
| api | 8080 | Go REST API backend |
| webui | 80 | React frontend |
| ai-knowledge-agent | 5000 | Python ML service (embeddings, RAG) |
| orchestrator | 9000 | Health monitoring service |

### Docker Volumes

- `postgres_data` → PostgreSQL data
- `redis_data` → Redis persistence
- `nas_data` → Main file storage (`/mnt/data`)
- `nas_backups` → Backup storage (`/mnt/backups`)

## Important Files and Locations

### Configuration
- `infrastructure/.env.prod` - Production environment variables
- `infrastructure/docker-compose.prod.yml` - Production stack
- `infrastructure/docker-compose.dev.yml` - Development stack

### Documentation
- `README.md` - Quick start guide
- `NAS_AI_SYSTEM.md` - Comprehensive system architecture
- `API_ENDPOINTS_COMPREHENSIVE.md` - Complete API reference
- `encryption_implementation_plan.md` - Encryption system design
- `BACKLOG.md` - Feature tracking and completed work

### Database
- `infrastructure/db/` - PostgreSQL migrations and schema

### Scripts
- `scripts/nas-cli.sh` - Unified CLI for deployment, testing, monitoring

## Development Workflow

### Adding a New API Endpoint

1. Create handler in `infrastructure/api/src/handlers/`
2. Add service logic in `infrastructure/api/src/services/`
3. Add repository methods in `infrastructure/api/src/repository/` (if database access needed)
4. Register route in `infrastructure/api/src/main.go`
5. Add tests in corresponding `*_test.go` files
6. Run `make test-coverage` to ensure 80% coverage
7. Test via `./scripts/nas-cli.sh` or direct API calls

### Adding a New Frontend Page

1. Create page component in `infrastructure/webui/src/pages/`
2. Extract reusable components to `infrastructure/webui/src/components/`
3. Create custom hooks in `infrastructure/webui/src/hooks/` for business logic
4. Add route in `infrastructure/webui/src/App.jsx`
5. Use GlassCard component for consistent UI styling

### Deploying Changes

```bash
# Full rebuild and deploy
cd infrastructure
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# Restart specific service
docker compose -f docker-compose.prod.yml restart api

# View logs
docker compose -f docker-compose.prod.yml logs -f api
```

## Common Patterns

### Go Error Handling

The codebase uses idiomatic Go error handling with structured logging:

```go
if err != nil {
    logger.WithError(err).Error("Failed to process request")
    c.JSON(http.StatusInternalServerError, gin.H{
        "status": "error",
        "error": "Internal server error",
    })
    return
}
```

### API Response Format

All API responses follow a consistent format:

```json
{
  "status": "ok|error",
  "data": { ... },
  "error": null | "Error message"
}
```

### Authentication Flow

1. User logs in → receives JWT access token + refresh token
2. Access token stored in memory, refresh token in httpOnly cookie
3. All protected routes require JWT in `Authorization: Bearer <token>` header
4. State-changing operations also require CSRF token in `X-CSRF-Token` header
5. Refresh token used to obtain new access token when expired

### File Operations Security

All file operations go through validation:

1. Path traversal check via `validatePath()` in `storage_service.go`
2. User quota enforcement
3. File type validation for uploads
4. Encryption/decryption for files in encrypted storage

## Environment Variables

### Required for API

- `JWT_SECRET` - Secret for JWT signing (generate with `openssl rand -base64 32`)
- `POSTGRES_PASSWORD` - PostgreSQL password
- `CSRF_SECRET` - Secret for CSRF token generation

### Optional

- `ENCRYPTED_STORAGE_PATH` - Path for encrypted file storage (default: disabled)
- `CORS_ORIGINS` - Allowed CORS origins
- `FRONTEND_URL` - Frontend URL for CORS
- `AI_SERVICE_URL` - AI agent URL (default: `http://ai-knowledge-agent:5000`)

## Notes

- The orchestrator runs parallel health checks every 30 seconds with thread-safe metrics
- The AI service uses Ollama locally - no external API dependencies
- Backup service supports scheduled backups with retention policies
- The encrypted storage feature is optional and can be disabled by not mounting the volume
- All secrets must be provided at startup (fail-fast configuration loading)
