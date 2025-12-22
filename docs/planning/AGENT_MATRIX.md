# Agenten-Matrix & Betriebshandbuch

**Version:** 2.0  
**Updated:** 2025-12-04

## Aktive Services

| Service | Typ | Port | Status |
|---------|-----|------|--------|
| **api** | Go Backend | 8080 | ✅ Aktiv |
| **webui** | Vite Frontend | 80 | ✅ Aktiv |
| **postgres** | pgvector DB | 5432 | ✅ Aktiv |
| **redis** | Cache | 6379 | ✅ Aktiv |
| **orchestrator** | Go Health Monitor | 9000 | ✅ Aktiv |
| **monitoring** | Go Metrics Agent | - | ✅ Aktiv |
| **analysis-agent** | Go Analytics | - | ✅ Aktiv |
| **pentester-agent** | Go Security | - | ✅ Aktiv |
| **ai-knowledge-agent** | Python ML | 5000 | ✅ Aktiv |

## Service-Verantwortlichkeiten

| Service | Kernverantwortung |
|---------|-------------------|
| **API** | REST Endpoints, Auth, File Management, Backup |
| **WebUI** | React Frontend, User Interface |
| **Orchestrator** | Health Checks, Prometheus Metrics, Service Registry |
| **Monitoring** | System Metrics Collection |
| **Analysis Agent** | Alert Analysis, Metric Evaluation |
| **Pentester Agent** | Security Scanning, Header Validation |
| **AI Knowledge Agent** | Embeddings, Semantic Search, pgvector |

## Entwicklungsrichtlinien

### Pflichtlektüre
1. `README.md` - Quick Start
2. `NAS_AI_SYSTEM.md` - Architektur
3. `docs/development/DEV_GUIDE.md` - Setup

### Code-Konventionen
- **Go:** `gofmt`, Error Handling, Context Usage
- **Python:** PEP8, Type Hints
- **React:** Functional Components, Hooks
- **Config:** Environment Variables, keine Hardcodes

### Commit-Format
```
type(scope): description

feat(api): add /embed endpoint
fix(auth): resolve JWT refresh bug
docs(readme): update setup instructions
chore(deps): upgrade dependencies
```

## Docker Commands

```bash
# Alle Services starten
docker compose -f infrastructure/docker-compose.prod.yml up -d

# Einzelnen Service neustarten
docker compose restart api

# Logs anzeigen
docker compose logs -f ai-knowledge-agent

# Health Check
curl http://localhost:8080/health
curl http://localhost:5000/health
curl http://localhost:9000/health
```

---

**Maintained by:** NAS.AI Team
