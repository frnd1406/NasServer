# NAS.AI Developer Guide

**Version:** 2.1  
**Updated:** 2025-12-04

## 1. Setup & Umgebung

### Voraussetzungen
- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### Quick Start

```bash
# 1. Repository klonen
git clone <repo-url>
cd f1406

# 2. Environment konfigurieren
cp infrastructure/.env.example infrastructure/.env.prod
# Secrets in .env.prod eintragen!

# 3. Dev-Infrastruktur starten
cd infrastructure
docker compose -f docker-compose.dev.yml up -d

# 4. Backend starten
cd api && go run ./src/main.go

# 5. Frontend starten (neues Terminal)
cd webui && npm install && npm run dev
```

## 2. Code-Konventionen

| Bereich | Standard |
|---------|----------|
| **Sprache** | Englisch (Kommentare, Variablen, Commits) |
| **Go** | `gofmt`, Error-Handling, Context-Usage |
| **React** | Functional Components, Hooks |
| **Config** | Keine Hardcoded Werte! `.env` oder Config-Structs |

## 3. Contributing Flow

1. Branch erstellen: `feature/<description>`
2. Implementieren (Tests schreiben)
3. Lokale Tests: `go test ./...`, `npm test`
4. Security Scan
5. Pull Request

## 4. Secrets Management

> ⚠️ **WICHTIG:** Alle Secrets gehören in `.env.prod` oder einen Secret Manager!

```bash
# Secrets NIEMALS im Code!
# Richtig:
JWT_SECRET=${JWT_SECRET}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Falsch:
JWT_SECRET="hardcoded-secret-123"
```

Für API-Tokens (Cloudflare, Resend, etc.):
- Development: `.env.prod`
- Production: Secret Manager / Vault

## 5. Troubleshooting

| Problem | Lösung |
|---------|--------|
| Port Konflikt | `lsof -i :8080` → Prozess killen |
| DB Connection | `docker compose ps` → postgres prüfen |
| Permission Denied | Schreibrechte in `/mnt/data` prüfen |

## 6. Nützliche Commands

```bash
# Alle Services starten
./scripts/nas-cli.sh

# Logs anzeigen
docker compose -f docker-compose.prod.yml logs -f api

# API testen
curl http://localhost:8080/health
```