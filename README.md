# NAS.AI

**Version:** 2.1  
**Status:** Phase 2.3 - Encryption  
**Updated:** 2025-12-07

Secure, AI-powered NAS with semantic search, zero-knowledge encryption, and glassmorphism UI.

---

## ✨ Features

### Core
- 🔒 **Security** - JWT/CSRF, CORS, Rate Limiting, Path Traversal Protection
- 💾 **File Manager** - Upload, Download, Rename, Delete, Trash/Restore
- 📦 **Batch Operations** - Multi-select, ZIP download, Folder download
- 🗂️ **Auto-Backup** - Scheduled with retention policies

### AI
- 🧠 **Semantic Search** - 1024D embeddings via mxbai-embed-large
- 💬 **RAG Chat** - Qwen2.5 LLM for intelligent Q&A
- 🤖 **AI Assistant** - Dedicated page with split layout
- 🔍 **Knowledge Base** - pgvector for vector similarity

### Security
- 🔐 **Zero-Knowledge Encryption** - AES-256-GCM + Argon2id
- 🔑 **Vault System** - Master password, lock/unlock
- 📁 **Encrypted Storage** - Files only readable via WebUI

### UI/UX
- 🎨 **Nebula Theme** - Glassmorphism design
- 🌙 **Dark/Light Mode** - Theme toggle
- ⌨️ **Keyboard Shortcuts** - Power user features
- 🔔 **Toast Notifications** - Real-time feedback
- 🔎 **Live Search** - Filter files in real-time

---

## Quick Start

```bash
# Start production stack
cd infrastructure
docker compose -f docker-compose.prod.yml up -d

# Access
# WebUI: http://localhost:8080
# API:   http://localhost:8080/api/v1
# AI:    http://localhost:5000
```

---

## Services

| Service | Port | Purpose |
|---------|------|---------|
| **API** | 8080 | Go backend |
| **WebUI** | 80 | Vite + React frontend |
| **AI Agent** | 5000 | Embeddings, RAG, Search |
| **Orchestrator** | 9000 | Health monitoring |
| **PostgreSQL** | 5432 | pgvector database |
| **Redis** | 6379 | Caching & sessions |

---

## Recent Changes

| Date | Feature |
|------|---------|
| 2025-12-07 | 🔐 Zero-Knowledge Encryption System |
| 2025-12-07 | 📦 Batch Download, Folder ZIP, Search Filter |
| 2025-12-06 | 🤖 Dedicated AI Assistant Page |
| 2025-12-05 | 🧠 1024D Vector Embeddings Upgrade |
| 2025-12-04 | 💬 RAG Endpoint with Qwen2.5 LLM |
| 2025-12-03 | ⚙️ Settings Page (Profile, Appearance, Admin) |

---

## Structure

```
f1406/
├── infrastructure/       # Docker services
│   ├── api/              # Go Backend
│   ├── webui/            # React Frontend
│   ├── ai_knowledge_agent/ # Python ML
│   └── db/               # PostgreSQL
├── orchestrator/         # Health monitoring
├── scripts/              # CLI tools
├── docs/                 # Blueprints
└── BACKLOG.md            # Feature tracking
```

---

## CLI

```bash
./scripts/nas-cli.sh
```

Features: Deployment, Logs, API Testing, Forensics, Database backup.

---

## Documentation

- **Backlog:** `BACKLOG.md`
- **API Endpoints:** `API_ENDPOINTS_COMPREHENSIVE.md`
- **CVE Checklist:** `CVE_CHECKLIST.md`
- **Encryption Plan:** `encryption_implementation_plan.md`

---

**License:** See LICENSE
