# NAS.AI API Backend

**Owner:** APIAgent
**Technologie:** Go 1.22+
**Status:** Phase 1 - Security-First Foundation
**Security-Posture:** 🔒 ZERO-TRUST from Day 1

---

## ⚠️ CRITICAL SECURITY NOTICE

**NIEMALS die API ohne vollständige Middleware-Kette starten!**

Jeder Handler MUSS durch diese Middleware-Chain laufen:
```
Request → Panic Recovery → CORS → Rate Limit → Auth → CSRF → Handler
```

**Consequence:** API ohne Middleware = Scheunentor offen → Sofortiges Security-Incident!

---

## OVERVIEW

Dies ist das Go-basierte Backend für das NAS.AI-System. Es folgt einem **Security-First Design**:
- Zero-Trust Architecture (Phase 1!)
- Middleware-Chain vor jedem Handler
- Fail-Fast bei fehlenden Dependencies/Secrets
- Path Sanitization by default
- Structured Audit Logging

---

## STRUKTUR

```
api/
├── src/
│   ├── main.go                 # Application entry point + Middleware setup
│   ├── middleware/             # 🔒 SECURITY LAYER (Epic 1 - PRIORITY!)
│   │   ├── panic.go           # Panic recovery
│   │   ├── cors.go            # CORS with strict origin checks
│   │   ├── ratelimit.go       # Rate limiting (per IP + per user)
│   │   ├── auth.go            # JWT validation (REQUIRED!)
│   │   ├── csrf.go            # CSRF token validation
│   │   ├── logging.go         # Request logging + audit trail
│   │   └── chain.go           # Middleware chain builder
│   ├── handlers/               # HTTP handlers (NEVER without middleware!)
│   │   ├── health.go          # Health checks (public endpoint)
│   │   ├── auth.go            # Authentication endpoints
│   │   ├── files.go           # File operations (protected)
│   │   └── websocket.go       # WebSocket (protected)
│   ├── models/                 # Data models
│   ├── services/               # Business logic
│   │   ├── auth_service.go    # JWT generation, validation
│   │   ├── file_service.go    # File ops with path sanitization
│   │   └── storage.go         # Storage abstraction
│   ├── repository/             # Data access layer
│   ├── config/                 # Configuration (12-factor)
│   │   └── config.go          # ENV-based config (NO DEFAULTS for secrets!)
│   └── utils/                  # Helper functions
│       ├── validator.go       # Input validation
│       └── sanitizer.go       # Path sanitization
├── tests/
│   ├── security/              # 🔒 Security tests (MUST PASS!)
│   │   ├── middleware_test.go
│   │   ├── auth_test.go
│   │   └── path_traversal_test.go
│   ├── integration/           # Integration tests
│   └── unit/                  # Unit tests
├── config/
│   ├── config.example.yaml    # Example config (NO SECRETS!)
│   └── middleware.yaml        # Middleware chain definition
├── docs/
│   ├── api-spec.yaml          # OpenAPI specification
│   └── security-gates.md      # Security checklist
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

---

## 🔐 PHASE 1 TASKS - SECURITY-FIRST ORDER

### Epic 1: Security Foundations & Middleware Chain (5-7 Tage) 🚨 PRIORITY #1

**⚠️ CRITICAL:** Diese Epic MUSS vor Epic 2 abgeschlossen sein. Keine API-Endpoints ohne Middleware!

**Tasks:**
1. ✅ Go module initialization (`go mod init github.com/nas-ai/api`)
2. ✅ Basic project structure (middleware/, handlers/, config/, services/)
3. ⏳ Configuration system (12-factor, ENV-based, fail-fast on missing secrets)
4. ⏳ **Middleware Chain Implementation:**
   - `panic.go` - Panic recovery (catch crashes, log, return 500)
   - `cors.go` - CORS with strict origin whitelist (NO wildcard `*`)
   - `ratelimit.go` - Rate limiting (100 req/min per IP, 1000/min per user)
   - `auth.go` - JWT validation middleware (REQUIRED for protected routes)
   - `csrf.go` - CSRF token validation for POST/PUT/DELETE
   - `logging.go` - Structured audit logging (request ID, user ID, IP, method, path, duration)
   - `chain.go` - Middleware chain builder
5. ⏳ Health check endpoint (`GET /health`) - **PUBLIC** (no auth, but rate-limited)
6. ⏳ Logger setup (structured JSON logging to stdout + audit log)
7. ⏳ Makefile (build, test, lint, run, security-scan)

**Deliverables:**
- `go.mod` + `go.sum`
- Complete middleware chain in `src/middleware/`
- Basic HTTP server with middleware chain active
- `/health` endpoint responding with 200 OK (public, but rate-limited)
- Structured logs showing middleware execution
- Security tests for middleware

**Acceptance Criteria:**
- ✅ All middleware tests passing (≥90% coverage)
- ✅ `make run` starts server with middleware chain active
- ✅ `curl http://localhost:8080/health` returns 200 (rate-limited)
- ✅ Requests without JWT to `/api/*` routes return 401 Unauthorized
- ✅ CORS blocks requests from non-whitelisted origins
- ✅ Rate limit kicks in after 100 requests/min
- ✅ Panic in handler caught and logged (returns 500)
- ✅ All logs structured JSON with request IDs

**Security Gates:**
- 🔒 **Gate 0:** Middleware chain enforced BEFORE any business logic
- 🔒 **Gate 1:** No public endpoints except `/health` (and later `/auth/*`)
- 🔒 **Gate 2:** CORS whitelist configured (NO `*`)
- 🔒 **Gate 3:** Rate limiting active

**References:**
- SECURITY_HANDBOOK.pdf → §2 (Security Gates)
- NAS_AI_SYSTEM.md → §4 (Security & Governance)

---

### Epic 2: Authentication & Secrets (5-7 Tage) 🔐 PRIORITY #2

**Dependencies:** Epic 1 COMPLETE (middleware chain ready)

**Tasks:**
1. ⏳ JWT authentication service
   - Token generation (access + refresh)
   - Token validation (integrated with `middleware/auth.go`)
   - Refresh token flow
   - **CRITICAL:** Load JWT secret from ENV (NO DEFAULTS!) → fixes SEC-2025-003
   - **CRITICAL:** Fail-fast on startup if JWT_SECRET missing
2. ⏳ User registration endpoint (`POST /auth/register`)
   - Password hashing (bcrypt, cost 12)
   - Input validation (email format, password strength ≥8 chars)
   - Rate limiting (5 registrations per IP per hour)
3. ⏳ Login endpoint (`POST /auth/login`)
   - Rate limiting (10 failed attempts per IP per hour → Fail2Ban trigger)
   - Audit logging (all attempts, success/failure)
4. ⏳ Logout endpoint (`POST /auth/logout`)
   - Token revocation (blacklist in Redis)
5. ⏳ Refresh token endpoint (`POST /auth/refresh`)
6. ⏳ CSRF token generation endpoint (`GET /auth/csrf-token`)
7. ⏳ Protected route example (`GET /api/profile` requires JWT)

**Deliverables:**
- Working auth endpoints (all behind middleware chain)
- JWT middleware integrated
- CSRF tokens generated and validated
- Unit tests (≥80% coverage)
- Integration tests for auth flow
- Security tests (brute-force, token forgery, CSRF bypass attempts)

**Acceptance Criteria:**
- ✅ User can register (password hashed, never logged)
- ✅ User can login (JWT returned, expires in 1h)
- ✅ Refresh token extends session (new access token)
- ✅ Logout revokes tokens (blacklist in Redis)
- ✅ JWT secret loaded from ENV (fail-fast if missing) → **SEC-2025-003 ✅**
- ✅ Protected routes return 401 without valid JWT
- ✅ CSRF validation blocks requests without token
- ✅ Rate limiting prevents brute-force
- ✅ All auth events logged (audit trail)
- ✅ PentesterAgent validates: No token forgery, no CSRF bypass

**Security Gates:**
- 🔒 **Gate 4:** All endpoints (except `/health`, `/auth/*`) require JWT
- 🔒 **Gate 5:** CSRF validation on POST/PUT/DELETE
- 🔒 **Gate 6:** No default JWT secret (SEC-2025-003)
- 🔒 **Gate 7:** Password hashing with bcrypt (cost ≥ 12)

**References:**
- CVE_CHECKLIST.md → SEC-2025-003 (JWT defaults)
- SECURITY_HANDBOOK.pdf → §1.2 (Secrets from ENV)
- NAS_AI_SYSTEM.md → §7.1 (Auth API Contract)

---

### Epic 3: File Operations API with Path Sanitization (5-7 Tage) 📁

**Dependencies:** Epic 2 COMPLETE (Auth + CSRF working)

**Tasks:**
1. ⏳ Path sanitization utility
   - Block path traversal (`../`, `..\\`, absolute paths)
   - Whitelist allowed directories
   - Canonicalize paths
   - Unit tests for all traversal vectors
2. ⏳ Storage abstraction layer
   - Interface for file operations
   - Quota checks
   - Trash sandbox (deleted files → trash folder, not immediate deletion)
3. ⏳ File listing endpoint (`GET /api/files`)
   - Requires JWT + CSRF token
   - Path sanitization enforced
   - Pagination support
4. ⏳ File upload (`POST /api/files`)
   - Max file size enforcement (100MB default)
   - MIME type validation (no executable uploads)
   - Virus scanning hook (future)
5. ⏳ File download (`GET /api/files/:path`)
   - Path sanitization enforced
   - Range request support
6. ⏳ File delete (`DELETE /api/files/:path`)
   - Move to trash (not immediate deletion)
   - Audit logging (who deleted what, when)
7. ⏳ File rename/move (`PATCH /api/files/:path`)

**Deliverables:**
- CRUD operations for files (all protected by middleware)
- Path sanitization preventing traversal attacks
- Unit + integration + security tests
- API documentation (OpenAPI spec)

**Acceptance Criteria:**
- ✅ All file operations work (list, upload, download, delete, rename)
- ✅ Path traversal attacks BLOCKED (PentesterAgent validated)
- ✅ Unauthorized users get 401
- ✅ CSRF validation on write operations
- ✅ File uploads limited to 100MB
- ✅ Deleted files in trash (not immediate deletion)
- ✅ All operations logged (audit trail)
- ✅ Tests cover happy path + attack vectors

**Security Gates:**
- 🔒 **Gate 8:** Path traversal protection (all vectors blocked)
- 🔒 **Gate 9:** File size limits enforced
- 🔒 **Gate 10:** MIME type validation (no executables)
- 🔒 **Gate 11:** Audit logging for all file ops

**References:**
- NAS_AI_SYSTEM.md → §7.1 (Files API Contract)
- CVE_CHECKLIST.md → SEC-003 (Path Traversal - validates fix)

---

### Epic 4: Observability & Fail-Fast (2-3 Tage) 📊

**Dependencies:** Epic 3 COMPLETE (Core API working)

**Tasks:**
1. ⏳ Dependency health checks on startup
   - PostgreSQL connection check
   - Redis connection check
   - Vault connection check (if used)
   - **CRITICAL:** Fail-fast if ANY critical dependency unreachable → fixes PERF-001
2. ⏳ Prometheus metrics exporter (`/metrics`)
   - Request count (by method, path, status)
   - Request duration (histogram)
   - Active connections
   - Auth failures (counter)
   - File operations (counter)
3. ⏳ Request ID middleware (for distributed tracing)
   - Generate unique ID per request
   - Include in all logs
   - Return in `X-Request-ID` header
4. ⏳ Structured error responses
   - JSON error format: `{"error": {"code": "...", "message": "..."}}`
   - No stack traces in production
   - Request ID in error response

**Deliverables:**
- Startup health checks (fail-fast)
- `/metrics` endpoint with Prometheus metrics
- Request IDs in all logs
- Structured error responses

**Acceptance Criteria:**
- ✅ App exits with clear error if DB unreachable → **PERF-001 ✅**
- ✅ App exits with clear error if Redis unreachable
- ✅ Prometheus can scrape `/metrics` (public, rate-limited)
- ✅ All logs contain request IDs
- ✅ Error responses include request ID (for debugging)
- ✅ No stack traces in production errors

**Security Gates:**
- 🔒 **Gate 12:** No sensitive data in error messages
- 🔒 **Gate 13:** No stack traces exposed to clients

**References:**
- CVE_CHECKLIST.md → PERF-001 (Fail-fast checks)
- NAS_AI_SYSTEM.md → §8 (Observability)

---

## 🔐 ZERO-TRUST REQUIREMENTS (PHASE 1)

### Architectural Principles

1. **Never Trust, Always Verify**
   - Every request authenticated (JWT) except `/health`, `/auth/*`
   - Every write operation validated (CSRF)
   - Every file path sanitized

2. **Least Privilege**
   - Default: Deny all
   - Explicit allow: `/health` (public), `/auth/*` (public)
   - All other routes: Require JWT + CSRF

3. **Defense in Depth**
   - Middleware chain (6 layers)
   - Input validation (all inputs)
   - Output sanitization (no XSS)
   - Audit logging (all actions)

4. **Fail Secure**
   - Missing JWT secret → Fail-fast (don't start)
   - Missing dependency → Fail-fast
   - Panic in handler → Caught, logged, 500 (not crash)
   - Invalid CSRF token → 403 (not ignored)

5. **Assume Breach**
   - All operations logged (audit trail)
   - Token revocation (logout blacklists)
   - Rate limiting (prevent brute-force)
   - Path sanitization (prevent traversal)

### Middleware Chain (ENFORCED!)

```go
// src/main.go
router := mux.NewRouter()

// Build middleware chain (order matters!)
chain := middleware.NewChain(
    middleware.PanicRecovery(),    // Catch panics
    middleware.CORS(),             // CORS with whitelist
    middleware.RateLimit(),        // Rate limiting
    middleware.RequestID(),        // Request ID for tracing
    middleware.Logging(),          // Audit logging
)

// Public endpoints (no auth, but rate-limited)
router.Handle("/health", chain.Then(handlers.Health())).Methods("GET")

// Auth endpoints (public, but rate-limited)
authRouter := router.PathPrefix("/auth").Subrouter()
authRouter.Handle("/register", chain.Then(handlers.Register())).Methods("POST")
authRouter.Handle("/login", chain.Then(handlers.Login())).Methods("POST")
authRouter.Handle("/csrf-token", chain.Then(handlers.CSRFToken())).Methods("GET")

// Protected API (requires JWT + CSRF)
apiChain := chain.Append(
    middleware.Auth(),             // JWT validation
    middleware.CSRF(),             // CSRF validation
)

apiRouter := router.PathPrefix("/api").Subrouter()
apiRouter.Handle("/profile", apiChain.Then(handlers.Profile())).Methods("GET")
apiRouter.Handle("/files", apiChain.Then(handlers.ListFiles())).Methods("GET")
apiRouter.Handle("/files", apiChain.Then(handlers.UploadFile())).Methods("POST")
// ...
```

---

## DEVELOPMENT SETUP

### Prerequisites
- Go 1.22+ installed
- Docker & Docker Compose (for DB, Redis)
- Access to `/home/user/Agent/docs/development/DEV_GUIDE.md`

### Quick Start

```bash
# 1. Navigate to API directory
cd /home/user/Agent/infrastructure/api

# 2. Initialize Go module
go mod init github.com/nas-ai/api

# 3. Install dependencies
go get github.com/gorilla/mux
go get github.com/golang-jwt/jwt/v5
go get github.com/go-redis/redis/v8
go get github.com/lib/pq
go get github.com/prometheus/client_golang/prometheus
go mod tidy

# 4. Setup environment variables (REQUIRED!)
export JWT_SECRET="$(openssl rand -base64 32)"  # NEVER use default!
# or point to a file (Vault/Docker secret)
# export JWT_SECRET_FILE="/run/secrets/jwt_secret"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="nas"
export DB_PASSWORD="$(openssl rand -base64 16)"
export DB_NAME="nas_db"
export REDIS_HOST="localhost"
export REDIS_PORT="6379"
export CORS_ORIGINS="http://localhost:5173"
export RATE_LIMIT_PER_MIN="100"

# 5. Start dependencies (Docker Compose)
docker compose -f ../docker-compose.dev.yml up -d

# 6. Run development server
go run src/main.go

# OR use Makefile
make run
```

### Environment Variables (REQUIRED!)

**Critical (Fail-Fast if Missing):**
- `JWT_SECRET` - JWT signing secret (min 32 bytes, NO default!) or `JWT_SECRET_FILE` pointing to a file
- `MONITORING_TOKEN` - Shared secret für `/monitoring/ingest`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL config
- `REDIS_HOST`, `REDIS_PORT` - Redis config

**Security:**
- `CORS_ORIGINS` - Comma-separated origin whitelist (e.g., `http://localhost:5173,https://nas.example.com`)
- `RATE_LIMIT_PER_MIN` - Rate limit per IP (default: 100)
- `CSRF_SECRET` - CSRF token secret (default: random on startup, but should be set for multi-instance)

**Optional:**
- `PORT` - API port (default: 8080)
- `LOG_LEVEL` - Logging level (default: info)
- `ENV` - Environment (development, production)

### Health Endpoint
- Route: `GET /health` (public, rate-limited)
- Prüft PostgreSQL und Redis (Fail-Fast-Philosophie). OK → `200`, degraded → `503` mit Details.
- Beispiel:
  ```bash
  curl -s http://localhost:8080/health | jq
  # {"status":"ok","dependencies":{"database":"ok","redis":"ok"},"service":"nas-api","version":"1.0.0-phase1",...}
  ```

### JWT-Secret-Rotation
1. Neues Secret erzeugen (mind. 32 Zeichen): `openssl rand -base64 48 > /run/secrets/jwt_secret`
2. Deployment auf neue Variable zeigen lassen (`JWT_SECRET_FILE=/run/secrets/jwt_secret` oder `JWT_SECRET="<value>"`).
3. API neu starten. Alte Access/Refresh Tokens werden dadurch ungültig (neue Signatur).
4. Optional: Nutzer abmelden/Refresh-Store leeren, falls harte Session-Invalidierung gewünscht ist.

### Testing

```bash
# Run all tests
make test

# Run with coverage (must be ≥80%)
make test-coverage

# Run security tests
make test-security

# Run integration tests
make test-integration

# Lint code
make lint

# Security scan (gitleaks + gosec)
make security-scan
```

---

## SECURITY REQUIREMENTS (MANDATORY!)

### MUST HAVE (Phase 1) - NO EXCEPTIONS

- 🔒 **Middleware chain ENFORCED on all routes**
  - Panic Recovery → CORS → Rate Limit → Auth → CSRF
- 🔒 **No hardcoded secrets** (JWT secret from ENV only)
  - Fail-fast if JWT_SECRET missing
- 🔒 **Path traversal protection** in all file operations
  - Whitelist directories
  - Canonicalize paths
  - Block `../`, absolute paths
- 🔒 **JWT validation on all protected routes**
  - No public API endpoints except `/health`, `/auth/*`
- 🔒 **CSRF validation on all write operations**
  - POST, PUT, DELETE require CSRF token
- 🔒 **Rate limiting** (prevent brute-force)
  - 100 req/min per IP (global)
  - 10 failed logins per IP per hour
  - 5 registrations per IP per hour
- 🔒 **Audit logging** (all operations)
  - Who, what, when, result
  - Stored in append-only log
- 🔒 **Fail-fast on missing dependencies**
  - DB unreachable → Exit with error
  - Redis unreachable → Exit with error
- 🔒 **Password hashing** (bcrypt, cost ≥ 12)
- 🔒 **No sensitive data in logs** (passwords, tokens)

### NICE TO HAVE (Phase 2)

- mTLS for service-to-service communication
- Vault integration for secrets
- Advanced rate limiting (per user, per endpoint)
- Audit logging to Loki

---

## 🔐 ENCRYPTION (Phase 2.3 - IMPLEMENTED)

### Zero-Knowledge Architecture

NAS.AI implements application-layer encryption using:
- **AES-256-GCM** for symmetric encryption
- **Argon2id** for key derivation
- **Zero-Knowledge** - Master password never stored

### Key Hierarchy

```
Master Password (user-only, never stored)
    └── KEK (Key Encryption Key, via Argon2id)
        └── DEK (Data Encryption Key, RAM-only)
            └── FILE_KEY, DB_KEY, BACKUP_KEY
```

### Vault API Endpoints

```yaml
# Public (auth required)
GET  /api/v1/vault/status    # { locked, configured, vaultPath }
GET  /api/v1/vault/config    # { vaultPath, configured }

# Protected (auth + CSRF)
POST /api/v1/system/vault/setup   # Initialize with masterPassword
POST /api/v1/system/vault/unlock  # Unlock vault
POST /api/v1/system/vault/lock    # Lock vault (wipes DEK)
PUT  /api/v1/system/vault/config  # Update vault path
```

### Encrypted Storage API

```yaml
GET    /api/v1/encrypted/status   # Encryption status
GET    /api/v1/encrypted/files    # List encrypted files
POST   /api/v1/encrypted/upload   # Encrypt & save file
GET    /api/v1/encrypted/download # Decrypt & download
GET    /api/v1/encrypted/preview  # Decrypt for preview
DELETE /api/v1/encrypted/delete   # Delete encrypted file
```

**Encrypted Storage Path:** `/media/frnd14/DEMO`

### Vault Guard Middleware

The `VaultGuard` middleware blocks access to protected routes when vault is locked:
- Returns `HTTP 423 Locked` with JSON error
- Protects `/api/v1/encrypted/*` endpoints

### Files Created

| File | Purpose |
|------|---------|
| `services/encryption_service.go` | Core encryption logic |
| `services/encrypted_storage.go` | Storage wrapper |
| `handlers/vault.go` | Vault API handlers |
| `handlers/encrypted_storage.go` | Encrypted storage handlers |
| `middleware/vault_guard.go` | Vault state middleware |

---

## API CONTRACTS

### Authentication

```yaml
POST /auth/register
Headers:
  Content-Type: application/json
Request:
  {
    "username": "string",
    "email": "string",
    "password": "string" (min 8 chars)
  }
Response (201 Created):
  {
    "user_id": "uuid",
    "access_token": "jwt",
    "refresh_token": "jwt",
    "expires_in": 3600
  }

POST /auth/login
Headers:
  Content-Type: application/json
Request:
  {
    "email": "string",
    "password": "string"
  }
Response (200 OK):
  {
    "user_id": "uuid",
    "access_token": "jwt",
    "refresh_token": "jwt",
    "expires_in": 3600,
    "csrf_token": "string"
  }
```

Full API specification: `docs/api-spec.yaml` (to be created in Epic 4)

---

## DEPENDENCIES

### Core
- `gorilla/mux` - HTTP router
- `golang-jwt/jwt/v5` - JWT tokens
- `lib/pq` - PostgreSQL driver
- `go-redis/redis/v8` - Redis client
- `bcrypt` - Password hashing

### Security
- `ulule/limiter` - Rate limiting
- `rs/cors` - CORS middleware
- `google/uuid` - Request IDs

### Observability
- `prometheus/client_golang` - Metrics
- `sirupsen/logrus` - Structured logging

### Testing
- `stretchr/testify` - Test assertions
- `httptest` - HTTP testing

---

## STATUS TRACKING

**Statuslogs:** `/home/user/Agent/status/APIAgent/phase1/`

**Format:** `NNN_YYYYMMDD_task-description.md`

**Required for each Epic:**
1. Analysis (Ist-Zustand, Risiken, Security-Check)
2. Implementation logs (step-by-step, middleware first!)
3. Security testing evidence (PentesterAgent validation)
4. Final summary with artefacts

---

## SECURITY GATES CHECKLIST

Before marking ANY Epic as complete, verify:

### Epic 1 Gates:
- [ ] Middleware chain implemented (6 layers)
- [ ] All middleware tests passing (≥90% coverage)
- [ ] CORS whitelist configured (NO `*`)
- [ ] Rate limiting active (verified with load test)
- [ ] Panic recovery catches all crashes
- [ ] Logs structured JSON with request IDs

### Epic 2 Gates:
- [ ] JWT secret loaded from ENV (NO default)
- [ ] Fail-fast on missing JWT_SECRET
- [ ] All auth endpoints protected by middleware
- [ ] CSRF tokens generated and validated
- [ ] Password hashing (bcrypt, cost 12)
- [ ] Token revocation on logout (Redis blacklist)
- [ ] Rate limiting on auth endpoints (brute-force protection)
- [ ] Audit logging for all auth events
- [ ] PentesterAgent validates: No token forgery, no CSRF bypass

### Epic 3 Gates:
- [ ] Path sanitization blocks ALL traversal vectors
- [ ] All file endpoints require JWT + CSRF
- [ ] File size limits enforced
- [ ] MIME type validation (no executables)
- [ ] Deleted files in trash (not immediate deletion)
- [ ] Audit logging for all file ops
- [ ] PentesterAgent validates: No path traversal, no unauthorized access

### Epic 4 Gates:
- [ ] Fail-fast on missing dependencies (DB, Redis)
- [ ] Prometheus metrics exposed on `/metrics`
- [ ] Request IDs in all logs
- [ ] No sensitive data in errors
- [ ] No stack traces in production

---

## NEXT STEPS

1. **APIAgent:** Read AGENT_CHECKLIST.md + this README
2. **APIAgent:** Create `status/APIAgent/phase1/001_20251121_epic1-security-foundations.md`
3. **APIAgent:** Start Epic 1 (Middleware Chain) - **COMPLETE BEFORE ANY ENDPOINTS!**
4. **APIAgent:** Report to Orchestrator after each Epic complete

---

## REFERENCES

- **System Architecture:** `/home/user/Agent/NAS_AI_SYSTEM.md`
- **Security Handbook:** `/home/user/Agent/docs/security/SECURITY_HANDBOOK.pdf`
- **CVE Checklist:** `/home/user/Agent/CVE_CHECKLIST.md`
- **Agent Matrix:** `/home/user/Agent/docs/planning/AGENT_MATRIX.md`
- **Dev Guide:** `/home/user/Agent/docs/development/DEV_GUIDE.md`

---

**Assigned:** 2025-11-21 (UPDATED with Security-First approach)
**Target:** Phase 1 Complete by 2025-11-28
**Owner:** APIAgent
**Priority:** 🔒 SECURITY FIRST - Middleware Chain before ANY business logic!

Terminal freigegeben.
