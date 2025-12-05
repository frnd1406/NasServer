# NAS.AI API Endpoints - Complete Summary

**Generated:** 2025-12-01  
**Location:** `/home/freun/Agent`  
**System:** Agent Orchestration Platform with AI Knowledge Integration

---

## Executive Summary

The NAS.AI system comprises 9 services across 3 tiers (Experience, Service/Orchestration, Filesystem) with a comprehensive REST API architecture. The system includes:

- **47 Active REST API Endpoints** across multiple services
- **AI/ML Integration** with semantic search capabilities via embeddings
- **Multi-tier Authentication** (JWT, CSRF, Monitoring Tokens)
- **Service Orchestration** with health monitoring and metrics

---

## Service Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      WebUI (Vite/React)                         │
│                      Port: 80/3000                              │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTPS/JWT
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    MAIN API (Gin/Go)                            │
│                    Port: 8080                                   │
│  ├─ 34 REST Endpoints                                           │
│  ├─ JWT Token Management                                       │
│  ├─ File Storage Management                                    │
│  ├─ Backup/Restore Operations                                  │
│  └─ Semantic Search (via AI Agent)                             │
└────┬─────────────────┬─────────────────┬──────────────────────┘
     │                 │                 │
     ▼                 ▼                 ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│AI Knowledge  │ │Orchestrator  │ │PostgreSQL    │
│Agent         │ │(Go)          │ │(pgvector)    │
│(Flask/Python)│ │Port: 9000    │ │Port: 5432    │
│Port: 5000    │ │4 Endpoints   │ │              │
│              │ │Health Check  │ │              │
│3 Endpoints   │ │Services      │ │              │
│- health      │ │Registry      │ │              │
│- process     │ │Metrics       │ │              │
│- embed_query │ │              │ │              │
└──────────────┘ └──────────────┘ └──────────────┘
     │                              │
     └──────────────┬───────────────┘
                    ▼
         ┌─────────────────────┐
         │ Redis Cache (6379)  │
         │ Token Storage       │
         └─────────────────────┘
         
Background Agents:
- Monitoring Agent (Go) → sends metrics to API
- Analysis Agent (Go) → creates alerts from metrics
- Pentester Agent → security testing
```

---

## Endpoint Categories & Count

| Category | Count | Service | Auth |
|----------|-------|---------|------|
| Public | 3 | API | None |
| Authentication | 9 | API | None/JWT |
| Protected APIs | 2 | API | JWT+CSRF |
| AI/ML | 4 | AI Agent + API | Internal/None |
| Monitoring | 5 | API | None/Token |
| Storage | 8 | API | JWT+CSRF |
| Backups | 4 | API | JWT+CSRF/Admin |
| System Settings | 3 | API | JWT+CSRF |
| Orchestrator | 4 | Orchestrator | None |
| **TOTAL** | **42** | | |

---

## Complete Endpoint List

### Public Endpoints (3)
```
GET  /health                                      [API:8080]
GET  /swagger/*any                                [API:8080] (dev only)
POST /monitoring/ingest                           [API:8080]
```

### Authentication Endpoints (9)
```
POST /auth/register                               [API:8080]
POST /auth/login                                  [API:8080]
POST /auth/refresh                                [API:8080]
POST /auth/logout                                 [API:8080]
POST /auth/verify-email                           [API:8080]
POST /auth/resend-verification                    [API:8080]
POST /auth/forgot-password                        [API:8080]
POST /auth/reset-password                         [API:8080]
GET  /api/v1/auth/csrf                            [API:8080]
```

### Protected API Endpoints (2)
```
GET  /api/profile                                 [API:8080]
GET  /api/monitoring                              [API:8080]
```

### AI/ML Endpoints (4)
```
GET  /health                                      [AI Agent:5000] (Flask)
POST /process                                     [AI Agent:5000]
POST /embed_query                                 [AI Agent:5000]
GET  /api/v1/search                               [API:8080]
```

### Monitoring & Metrics (5)
```
POST /api/v1/system/metrics                       [API:8080]
GET  /api/v1/system/metrics                       [API:8080]
GET  /api/v1/system/alerts                        [API:8080]
POST /api/v1/system/alerts                        [API:8080]
POST /api/v1/system/alerts/:id/resolve            [API:8080]
```

### Storage & File Management (8)
```
GET  /api/v1/storage/files                        [API:8080]
POST /api/v1/storage/upload                       [API:8080]
GET  /api/v1/storage/download                     [API:8080]
POST /api/v1/storage/rename                       [API:8080]
DELETE /api/v1/storage/delete                     [API:8080]
GET  /api/v1/storage/trash                        [API:8080]
POST /api/v1/storage/trash/restore/:id            [API:8080]
DELETE /api/v1/storage/trash/:id                  [API:8080]
```

### Backup Management (4)
```
GET  /api/v1/backups                              [API:8080]
POST /api/v1/backups                              [API:8080]
POST /api/v1/backups/:id/restore                  [API:8080]
DELETE /api/v1/backups/:id                        [API:8080]
```

### System Settings (3)
```
GET  /api/v1/system/settings                      [API:8080]
PUT  /api/v1/system/settings/backup               [API:8080]
POST /api/v1/system/validate-path                 [API:8080]
```

### Orchestrator Endpoints (4)
```
GET  /health                                      [Orchestrator:9000]
GET  /api/services                                [Orchestrator:9000]
GET  /api/registry                                [Orchestrator:9000]
GET  /metrics                                     [Orchestrator:9000]
```

---

## Key Features by Domain

### Authentication & Security
- JWT token-based authentication
- CSRF protection on state-changing operations
- Email verification with token-based verification
- Password reset with secure tokens
- Session management and logout
- Rate limiting (5 req/min on auth endpoints)
- Bearer token validation
- Monitoring token authentication for agents

### File Management & Storage
- Full CRUD operations on files
- Trash/recycle bin functionality
- Path traversal protection
- Directory listing with file metadata
- Multipart file upload support
- File renaming and metadata operations
- Automatic AI embedding generation on upload
- Storage mounted at `/mnt/data`

### AI & Semantic Search
- Sentence-Transformers model (all-MiniLM-L6-v2)
- Vector embeddings stored in PostgreSQL pgvector
- Automatic embedding generation on file upload
- Semantic similarity search with cosine distance
- 384-dimensional embeddings
- Model loading status monitoring
- Query embedding API for client-side usage

### Backup & Disaster Recovery
- Scheduled backup creation (cron-based)
- Configurable retention policies
- Backup storage at `/mnt/backups`
- Restore capabilities (admin-only)
- Backup listing and metadata retrieval
- Persistent settings storage
- Deletion with access control

### System Monitoring
- Real-time CPU, RAM, Disk metrics
- Automatic alert generation (thresholds: CPU 80%, RAM 90%)
- Alert resolution tracking
- Severity levels (WARNING, CRITICAL)
- Lookback window analysis
- 10-second default collection interval
- Prometheus metrics export

### Service Orchestration
- Health check monitoring for all services
- Service registry management
- Uptime percentage calculation
- Consecutive failure tracking
- 30-second health check interval
- Service status response with metadata
- Graceful degradation handling

---

## Authentication Methods

| Method | Token Type | Storage | TTL | Usage |
|--------|-----------|---------|-----|-------|
| JWT | Bearer Token | Redis | Session-based | User authentication |
| CSRF | Token Header | Redis | Session-based | State-changing operations |
| Monitoring Token | X-Monitoring-Token | Environment | N/A | Metrics submission |
| API Key | X-Monitoring-Token | Environment | N/A | Service-to-service |

---

## Database Schema

### Key Tables
```
users
├── id (UUID)
├── email (string, unique)
├── password_hash (string)
├── email_verified (boolean)
├── created_at (timestamp)

system_metrics
├── id (serial)
├── agent_id (string)
├── cpu_usage (float)
├── ram_usage (float)
├── disk_usage (float)
├── created_at (timestamp, indexed)

system_alerts
├── id (UUID)
├── severity (enum: WARNING, CRITICAL)
├── message (string)
├── is_resolved (boolean)
├── created_at (timestamp)

file_embeddings (pgvector)
├── id (serial)
├── file_id (string, unique)
├── file_path (string)
├── mime_type (string)
├── content (text)
├── embedding (vector(384))
├── created_at (timestamp)

system_settings
├── key (string, primary key)
├── value (string)
```

---

## Configuration & Environment

### Required Environment Variables
```
JWT_SECRET              Secret key for JWT signing
DATABASE_URL           PostgreSQL connection string
REDIS_URL              Redis connection string
POSTGRES_PASSWORD      PostgreSQL password
MONITORING_TOKEN       Token for metrics submission
```

### Optional Environment Variables
```
PORT                   API port (default: 8080)
ENV                    Environment (development/production)
LOG_LEVEL              Logging level (default: info)
CORS_ORIGINS           Comma-separated CORS origins
RATE_LIMIT_PER_MIN     Rate limit (default: 100)
FRONTEND_URL           Frontend URL for CORS
AI_SERVICE_URL         AI Knowledge Agent URL
BACKUP_SCHEDULE        Cron expression for backups
BACKUP_RETENTION_COUNT Number of backups to keep
BACKUP_STORAGE_PATH    Backup storage directory
```

---

## Rate Limiting

- **Auth Endpoints:** 5 requests/minute (strict)
- **Standard Endpoints:** 100 requests/minute
- **Search:** 100 requests/minute
- **Health Check:** No limit
- **Metrics:** Token-based, no per-IP limit

Returns `429 Too Many Requests` when exceeded.

---

## Security Headers

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Strict-Transport-Security: max-age=31536000
Content-Security-Policy: default-src 'self'
```

---

## CORS Configuration

- Origins configured via `CORS_ORIGINS` environment variable
- Preflight (OPTIONS) requests allowed globally without authentication
- Credentials included when configured
- Standard HTTP methods (GET, POST, PUT, DELETE)

---

## Future Endpoints (Planned)

Based on NAS_AI_SYSTEM.md roadmap:

### User Management
```
GET    /api/v1/users                    User listing
POST   /api/v1/users                    Create user
PUT    /api/v1/users/:id                Update user
DELETE /api/v1/users/:id                Delete user
POST   /api/v1/users/:id/reset-password Admin password reset
GET    /auth/sessions                   List sessions
POST   /auth/sessions/:id/revoke        Revoke session
```

### Files & Favorites
```
GET    /api/v1/files/:id/thumbnail      File thumbnails
POST   /api/v1/files/zip                Create archive
POST   /api/v1/files/unzip              Extract archive
GET    /api/v1/favorites                List favorites
POST   /api/v1/favorites                Add favorite
DELETE /api/v1/favorites/:id            Remove favorite
```

### Storage Analytics
```
GET    /api/v1/storage/overview         Storage overview
GET    /api/v1/storage/usage            Usage statistics
GET    /api/v1/storage/directory-size   Directory size
GET    /api/v1/storage/alerts           Storage alerts
GET    /api/v1/storage/trends           Usage trends
```

### Shares (SMB/NFS/FTP)
```
GET    /api/v1/shares                   List shares
POST   /api/v1/shares                   Create share
PUT    /api/v1/shares/:id               Update share
DELETE /api/v1/shares/:id               Delete share
```

### AI Services
```
GET    /api/v1/ai/search                Enhanced search with facets
POST   /api/v1/ai/documents/:id/analyze Document analysis
GET    /api/v1/ai/suggestions           AI recommendations
```

### Documentation & Settings
```
GET    /api/v1/docs-terminal/*          Terminal commands
GET    /api/v1/docs-settings            Settings
PUT    /api/v1/docs-settings            Update settings
GET    /api/v1/settings/security        Security config
PUT    /api/v1/settings/security        Update security
GET    /api/v1/auth/audit               Audit logs
```

---

## Testing & Development

### Available in Development
- Swagger/OpenAPI documentation at `/swagger/index.html`
- Enhanced JSON logging
- Relaxed CORS for local development
- Debug mode support

### Testing Resources
- Integration tests in `/infrastructure/api/test/integration/`
- Handler tests in `/infrastructure/api/src/handlers/`
- Example implementations in `/infrastructure/api/examples/`

### Health Check Tests
```bash
curl http://localhost:8080/health
curl http://localhost:5000/health        # AI Agent
curl http://localhost:9000/health        # Orchestrator
```

---

## Performance Characteristics

### Timeout Settings
- HTTP Server: 15s read, 15s write, 60s idle
- AI Service: 8s client timeout
- Orchestrator: 5s service check timeout
- Database: 5s query timeout

### Model Loading
- AI model loads asynchronously on startup
- Health endpoint returns 503 until model is ready
- Embedding dimension: 384
- Processing time: ~100ms per document

### Database Performance
- Indexed timestamp columns for metrics queries
- pgvector HNSW index for embedding similarity
- Connection pooling via sqlx
- Prepared statements for security

---

## File Paths & Storage

```
/mnt/data/              Primary data storage (files)
/mnt/backups/           Backup storage
/var/log/               Container logs
/srv/webui/             WebUI artifacts
/var/lib/postgresql/    PostgreSQL data
/var/lib/redis/         Redis data
```

---

## Deployment Services

From `docker-compose.prod.yml`:

| Service | Image | Port | Memory | Dependencies |
|---------|-------|------|--------|--------------|
| postgres | pgvector/pgvector:pg16 | 5432 | unlimited | - |
| redis | redis:7-alpine | 6379 | unlimited | - |
| api | nas-api:1.0.0 | 8080 | unlimited | postgres, redis |
| webui | nas-webui:1.0.0 | 80 | unlimited | api |
| monitoring | nas-monitoring:1.0.0 | - | unlimited | api |
| analysis-agent | nas-analysis:1.0.0 | - | unlimited | postgres, api |
| pentester-agent | nas-pentester:1.0.0 | - | unlimited | api |
| ai-knowledge-agent | nas-ai-knowledge-agent:1.0.0 | 5000 | 4GB limit, 2GB reservation | postgres |

---

## References

- **Main API:** `/home/freun/Agent/infrastructure/api/src/main.go`
- **AI Agent:** `/home/freun/Agent/infrastructure/ai_knowledge_agent/src/main.py`
- **Orchestrator:** `/home/freun/Agent/orchestrator/api.go`
- **Documentation:** `/home/freun/Agent/NAS_AI_SYSTEM.md`
- **Configuration:** `/home/freun/Agent/agents-config.yaml`

---

**Document Version:** 1.0  
**Last Updated:** 2025-12-01  
**Scope:** All API endpoints in Agent codebase (complete)
