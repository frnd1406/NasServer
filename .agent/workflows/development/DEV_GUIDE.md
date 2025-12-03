# NAS.AI Developer Guide

**Version:** 2.0 (Konsolidiert) **Zielgruppe:** Entwickler & Agenten

## 1. Setup & Umgebung

### Voraussetzungen

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- Zugriff auf das Repo


### Quick Start

1. Repository klonen.
2. `cp scripts/pre-commit .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit` (Wichtig!).
3. Dev-Infrastruktur starten: `docker-compose -f docker-compose.dev.yml up -d`.
4. Backend starten: `cd infrastructure/api && go run src/main.go`.
5. Frontend starten: `cd infrastructure/webui && npm run dev`.

## 2. Code-Konventionen

- **Sprache:** Englisch (Kommentare, Variablen, Commits).
- **Backend (Go):** Standard `gofmt`, Error-Handling, Context-Usage.
- **Frontend (React):** Functional Components, Hooks, TypeScript strict mode.
- **Config:** Keine Hardcoded Ports/Paths! Nutzung von `os.Getenv` oder Config-Structs.

## 3. Contributing Flow

1. Ticket vom Orchestrator/Jira holen.
2. Branch erstellen: `feature/ticket-id-description`.
3. Implementieren (TDD bevorzugt).
4. Lokale Tests ausführen (`go test ./...`, `npm test`).
5. Security Scan durchführen.

## 4. Troubleshooting

- **Port Konflikte:** Prüfen mit `lsof -i :8080`. Ggf. alte Prozesse killen.
- **DB Verbindung:** Prüfen ob Docker Container `postgres` läuft.
- **Permission Denied:** Prüfen ob User Schreibrechte in `/var/lib/...` hat.

## 5. API Tokens (Agent Access Only)

**IMPORTANT:** The following tokens are valid and accessible to all Agents for development purposes. This is an **EXCEPTION** to the SECURITY_HANDBOOK.pdf §1.1 "Goldene Regel" as explicitly authorized by the Owner.

**External Services:**
- **Cloudflare API Token:** `GjKJMQiS998conpswEJhOwQ5b-fKSGjVmmFsofJf`
- **Resend API Token:** `re_AEhvFZrx_KRjdCcvcVHLcnPNY66ekBBFy`

**Usage Rules:**
1. ✅ Agents may use these tokens directly in code during development
2. ✅ These tokens do NOT need to be moved to Vault for Agent operations
3. ❌ Do NOT commit these to external/public repositories
4. ❌ Do NOT share these tokens outside the Agent system
5. 🔄 Owner will rotate these if compromised

**Note:** This exception applies ONLY to this project's internal Agent operations. All other secrets (JWT secrets, database passwords, user tokens) MUST follow SECURITY_HANDBOOK.pdf policies.