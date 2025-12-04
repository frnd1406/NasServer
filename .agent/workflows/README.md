# NAS.AI

**Version:** 2.0  
**Status:** Phase 2.2 - AI Core  
**Updated:** 2025-12-04

Secure, automated storage with semantic AI search and glassmorphism UI.

---

## Features

- 🔒 **Security** - JWT/CSRF, CORS, Rate Limiting
- 🧠 **AI Search** - Semantic embeddings with pgvector
- 💾 **Auto-Backup** - Scheduled with retention policies
- 🎨 **Nebula UI** - Glassmorphism design

---

## Quick Start

```bash
# Start production stack
cd infrastructure
docker compose -f docker-compose.prod.yml up -d

# Access
# API:   http://localhost:8080
# WebUI: http://localhost:8080 (via nginx)
# AI:    http://localhost:5000
```

---

## Services

| Service | Port | Purpose |
|---------|------|---------|
| **API** | 8080 | Go backend |
| **WebUI** | 80 | Vite frontend |
| **AI Agent** | 5000 | Embeddings & search |
| **Orchestrator** | 9000 | Health monitoring |
| **PostgreSQL** | 5432 | pgvector database |
| **Redis** | 6379 | Caching |

---

## Structure

```
f1406/
├── infrastructure/       # Docker services (API, WebUI, AI, etc.)
├── orchestrator/         # Health monitoring service
├── scripts/              # CLI tools (nas-cli.sh)
├── docs/                 # Blueprints, policies
├── contrib/              # Trivy report templates
└── .archive/             # Historical reports
```

---

## CLI

```bash
./scripts/nas-cli.sh
```

Features: Deployment, Logs, API Testing, Forensics, Database backup.

---

## Documentation

- **API Endpoints:** `API_ENDPOINTS_COMPREHENSIVE.md`
- **CVE Checklist:** `CVE_CHECKLIST.md`
- **Blueprints:** `docs/blueprints/`

---

**License:** See LICENSE
