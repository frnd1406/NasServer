# NAS.AI API Documentation Index

This directory contains comprehensive API documentation for the NAS.AI Agent Orchestration Platform.

## Documentation Files

### Primary Documentation (Single Source of Truth)

#### **API_ENDPOINTS_COMPREHENSIVE.md** ⭐ (START HERE)
- Complete API documentation
- Detailed request/response examples  
- Authentication requirements
- Query parameters and request bodies
- Service specifications table
- WebSocket topics for real-time events
- **Use for:** All API integration and development

---

### Archived Documentation

The following files have been consolidated and moved to `.archive/deprecated-docs/`:
- `API_ENDPOINTS.md` (German, superseded)
- `API_ENDPOINTS_SUMMARY.md` (merged into comprehensive)
- `API_ENDPOINTS_QUICK_REFERENCE.txt` (merged into comprehensive)

---

## Quick Start

### For New Developers
1. Start with **API_ENDPOINTS_SUMMARY.md** to understand the architecture
2. Review the service architecture diagram
3. Read the key features section relevant to your work
4. Check **API_ENDPOINTS_COMPREHENSIVE.md** for specific endpoint details

### For API Integration
1. Use **API_ENDPOINTS_QUICK_REFERENCE.txt** to find your endpoint
2. Check the endpoint in **API_ENDPOINTS_COMPREHENSIVE.md** for full details
3. Look up authentication requirements in the summary
4. Check authentication section for token handling

### For System Administration
1. Review the service specifications in API_ENDPOINTS_SUMMARY.md
2. Check the deployment services section
3. Review environment variables and configuration
4. Check the database schema

---

## Key Statistics

| Metric | Value |
|--------|-------|
| Total API Endpoints | 42 active + 20+ planned |
| Services | 9 (API, AI Agent, Orchestrator, etc.) |
| Primary Language (API) | Go (Gin framework) |
| AI/ML Integration | Python (Flask) with Sentence-Transformers |
| Authentication Methods | 3 (JWT, CSRF, Monitoring Token) |
| Database | PostgreSQL with pgvector |
| Cache | Redis |
| API Port | 8080 |
| AI Agent Port | 5000 |
| Orchestrator Port | 9000 |

---

## API Services

### Main API (Port 8080)
- **Framework:** Go Gin
- **Endpoints:** 34 REST endpoints
- **Features:** File management, backups, authentication, search
- **Location:** `/home/freun/Agent/infrastructure/api/src/main.go`

### AI Knowledge Agent (Port 5000)
- **Framework:** Python Flask
- **Endpoints:** 3 endpoints
- **Features:** Semantic embeddings, query processing
- **Model:** sentence-transformers/all-MiniLM-L6-v2
- **Location:** `/home/freun/Agent/infrastructure/ai_knowledge_agent/src/main.py`

### Orchestrator (Port 9000)
- **Framework:** Go
- **Endpoints:** 4 endpoints
- **Features:** Service health checks, registry management
- **Location:** `/home/freun/Agent/orchestrator/api.go`

---

## Endpoint Categories

### Public (3 endpoints)
- Health checks
- Swagger documentation
- Monitoring ingestion

### Authentication (9 endpoints)
- Register, Login, Logout
- Token refresh
- Email verification
- Password reset
- CSRF token generation

### Protected (2 endpoints)
- User profile
- Monitoring data

### AI/ML (4 endpoints)
- File embedding generation
- Query embedding
- Semantic search
- Health monitoring

### Storage (8 endpoints)
- File CRUD operations
- Trash management
- Upload/Download
- Renaming

### Backups (4 endpoints)
- List, Create, Restore, Delete
- Admin-only restore/delete

### System (5 endpoints)
- Metrics submission and retrieval
- Alert management
- Settings management

### Orchestrator (4 endpoints)
- Health check
- Service status
- Registry listing
- Prometheus metrics

---

## Authentication Overview

### JWT Tokens
- Issued by: `POST /auth/login`
- Stored in: Redis
- Usage: Bearer token in Authorization header
- Rate limit: 5 requests/minute

### CSRF Tokens
- Issued by: `GET /api/v1/auth/csrf`
- Stored in: Redis
- Usage: X-CSRF-Token header
- Required for: State-changing operations (POST, PUT, DELETE)

### Monitoring Token
- Type: X-Monitoring-Token header
- Usage: Metrics submission from agents
- Validation: Token-based (not rate-limited per IP)

---

## Database

### PostgreSQL Connection
- Host: postgres:5432
- Database: nas_db
- Extensions: pgvector (for embeddings)

### Key Tables
- **users** - User accounts and authentication
- **system_metrics** - CPU, RAM, Disk metrics
- **system_alerts** - System warnings
- **file_embeddings** - Vector embeddings for semantic search
- **system_settings** - Configuration storage
- **monitoring** - Agent monitoring data

---

## Security Features

- JWT token-based authentication
- CSRF protection on state-changing operations
- Rate limiting (5-100 requests/minute depending on endpoint)
- Path traversal protection
- Security headers (CSP, X-Frame-Options, HSTS)
- Email verification workflow
- Password reset with secure tokens
- Monitoring token validation

---

## Configuration

### Required Environment Variables
```
JWT_SECRET
DATABASE_URL
REDIS_URL
POSTGRES_PASSWORD
MONITORING_TOKEN
```

### Optional Environment Variables
```
PORT=8080
ENV=production
LOG_LEVEL=info
CORS_ORIGINS=<origins>
RATE_LIMIT_PER_MIN=100
AI_SERVICE_URL=http://ai-knowledge-agent:5000
```

---

## Common Use Cases

### User Authentication Flow
1. `POST /auth/register` - Create account
2. `POST /auth/login` - Get JWT token
3. `GET /api/v1/auth/csrf` - Get CSRF token
4. Use JWT + CSRF for protected endpoints

### File Upload & Search
1. `POST /api/v1/storage/upload` - Upload file
2. AI agent automatically generates embeddings
3. `GET /api/v1/search?q=query` - Semantic search

### System Monitoring
1. Monitoring agent sends metrics to `POST /api/v1/system/metrics`
2. Analysis agent checks metrics and creates alerts
3. Frontend retrieves metrics and alerts

### Backup & Restore
1. `POST /api/v1/backups` - Create backup
2. Scheduled via cron (configurable)
3. `POST /api/v1/backups/:id/restore` - Restore (admin only)

---

## Development

### Testing
- Integration tests in: `/infrastructure/api/test/integration/`
- Handler tests in: `/infrastructure/api/src/handlers/`
- Swagger docs available at: `/swagger/index.html` (dev only)

### Health Checks
```bash
# Main API
curl http://localhost:8080/health

# AI Agent
curl http://localhost:5000/health

# Orchestrator
curl http://localhost:9000/health
```

---

## Performance

### Timeouts
- HTTP Read: 15 seconds
- HTTP Write: 15 seconds
- HTTP Idle: 60 seconds
- Database Query: 5 seconds
- AI Service Call: 8 seconds

### Limits
- Max header size: 1 MB
- Embedding dimension: 384
- Rate limit: 100 req/min (standard), 5 req/min (auth)

---

## File Locations

| Path | Purpose |
|------|---------|
| `/mnt/data/` | File storage |
| `/mnt/backups/` | Backup storage |
| `/infrastructure/api/` | Main API code |
| `/infrastructure/ai_knowledge_agent/` | AI service code |
| `/orchestrator/` | Orchestrator code |
| `/webui/` | Frontend code |

---

## Deployment

### Docker Compose Files
- **Production:** `docker-compose.prod.yml`
- **Development:** `docker-compose.dev.yml`

### Services
1. PostgreSQL (pgvector)
2. Redis
3. API (Gin)
4. WebUI (Vite/React)
5. AI Knowledge Agent
6. Orchestrator
7. Monitoring Agent
8. Analysis Agent
9. Pentester Agent

---

## Future Endpoints

See **API_ENDPOINTS_SUMMARY.md** section "Future Endpoints (Planned)" for:
- User management
- File favorites
- Storage analytics
- Share management (SMB/NFS/FTP)
- Enhanced AI services
- Documentation terminal
- Security audit logs

---

## Support & References

### Main Configuration
- `/home/freun/Agent/agents-config.yaml`
- `/home/freun/Agent/infrastructure/.env.prod`

### System Documentation
- `/home/freun/Agent/NAS_AI_SYSTEM.md` - System architecture and design

### Source Code
- Main API: `/home/freun/Agent/infrastructure/api/src/main.go`
- AI Agent: `/home/freun/Agent/infrastructure/ai_knowledge_agent/src/main.py`
- Orchestrator: `/home/freun/Agent/orchestrator/api.go`

---

## Document Information

- **Created:** 2025-12-01
- **System:** NAS.AI Agent Orchestration Platform
- **Version:** 1.0
- **Scope:** Complete API endpoint documentation

---

**Navigation Tips:**
- Use Ctrl+F to search across documents
- Check headers for quick navigation
- Reference the service architecture diagram in SUMMARY
- Cross-reference COMPREHENSIVE for detailed examples
