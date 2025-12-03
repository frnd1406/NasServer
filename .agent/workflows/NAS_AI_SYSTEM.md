# NAS.AI – Struktur & Infrastruktur (ASCII Blueprint)
## Table of Contents
1. [System Overview](#1-system-overview)
2. [Communication & Storage Topology](#2-communication--storage-topology)
3. [Data & Control Flows](#3-data--control-flows)
4. [Security & Governance](#4-security--governance)
5. [Filesystem Layout](#5-filesystem-layout-ascii-map)
6. [WebUI Architecture Blueprint](#6-webui-architecture-blueprint)
7. [API & Event Contracts](#7-api--event-contracts)
8. [State Management, Testing & QA Gates](#8-state-management-testing--qa-gates)
9. [Offline & PWA Strategy](#9-offline--pwa-strategy)
10. [Agent Incident Response Playbook](#10-agent-incident-response-playbook)
11. [Next Features & Initiatives](#11-next-features--initiatives)
12. [References & Mandatory Reading](#12-references--mandatory-reading)
## 1. System Overview
> **Security Notice:** Dieses Dokument dient ausschließlich zur Navigation und
> Steuerung der NAS.AI-Governance. Operative Agenten dürfen keine direkten
> Änderungen vornehmen und erhalten nur die für ihre Rolle freigegebenen Auszüge.
> Alle Zugriffe werden über den AgentOrchestrator protokolliert.
```
                             ┌────────────────────────┐
                             │   Benutzer & Clients   │
                             │  WebUI / Mobile / CLI  │
                             └──────────┬─────────────┘
                                        │ HTTPS (JWT, mTLS)
                                        ▼
┌────────────────────────────────────────────────────────────────────┐
│                        Experience Tier (Phase 3)                    │
│ ┌─────────────┐   ┌───────────────┐   ┌─────────────────────────┐  │
│ │ React/Vite  │<->│ API SDK (TS)  │<->│ Realtime WS Edge (Auth) │  │
│ │ (WebUIAgent)│   │ Mobile SDK    │   │ Toast/Events/Push       │  │
│ └─────────────┘   └───────────────┘   └─────────────────────────┘  │
└───────────┬────────────────────────────────────────────────────────┘
            │ GraphQL/REST (JWT) + WS
            ▼
┌──────────────────────────────────────────────────────────────────┐
│                       Service & Orchestrator Tier                │
│ ┌─────────────┐   ┌────────────────┐   ┌────────────────────────┐│
│ │ APIAgent    │<->│ Event Bus (NATS│<->│ Orchestrator (Go)     │ │
│ │ (Go micro-  │   │  /MQTT)        │   │ Queue + State DB      │ │
│ │ services)   │   └────────────────┘   └──────────┬────────────┘ │
│ └─────┬───────┘                                   │Events        │
│       │ gRPC/REST                                 ▼              │
│ ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌────────┐ │
│ │ SystemSetup  │  │ NetworkSec   │  │ Documentation│  │ Pentester│ │
│ │ Agent        │  │ Agent        │  │ Agent        │  │ Agent  │ │
│ └─────┬────────┘  └────┬─────────┘  └────┬─────────┘  └──┬─────┘ │
│       │sudo/Ansible     │iptables/CAs     │docs-as-code     │sec-tests│
│       │                 │                 │                 │      │
│ ┌─────▼────┐ ┌──────────▼──────┐                                  │
│ │ Analysis │ │ WebUIAgent      │                                  │
│ │ Agent    │ │(React/Vite)     │                                  │
│ └──────────┘ └─────────────────┘                                  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
            │
            ▼
┌───────────────────┐      ┌──────────────────────┐
│ Filesystem Tier   │      │ Observability & AI FS │
│ (/srv, /mnt/raid) │      │ (/var/otel, /var/ai) │
└───────────────────┘      └──────────────────────┘
```
## 2. Communication & Storage Topology
```
   Event Flow                          Telemetrie                     Storage
┌──────────────┐      ┌────────────┐      ┌──────────────┐          ┌─────────────┐
│ Agents       │ ---> │ Event Bus  │ ---> │ Orchestrator │ ---> DR  │ /var/lib/...│
└─────┬────────┘      └────┬───────┘      └──────┬───────┘          └────┬────────┘
      │Status/Logs          │Queues                │Workflows             │Mounts
      ▼                     ▼                      ▼                      ▼
┌──────────────┐   ┌──────────────────┐  ┌──────────────────┐   ┌────────────────┐
│ /var/log/    │   │ /var/lib/queues/ │  │ /etc/nas/*.yaml  │   │ /mnt/raid/data │
│ agent-*.log  │   │ pkg,events,drift │  │ policies, allow  │   │ /mnt/raid/back │
└──────────────┘   └──────────────────┘  └──────────────────┘   └────────────────┘
Telemetry Bus:  node_exporter → Prometheus → Alertmanager/Grafana (OTel Collector ↔ Loki/Tempo)
Observability Storage: `/var/lib/observability/{prometheus,grafana,loki,tempo}`
AI Knowledge Layer (Future):
 - Vector DB (`/var/lib/ai/vector-db`)
 - Embedding Cache (`/srv/ai/cache`)
 - Model Store (`/var/lib/ai/models`, signed binaries)
Backup & Restore Paths (Future):
 - `/mnt/raid/backups/jobs` (primary snapshots)
 - `/mnt/raid/backups/history` (metadata + logs)
 - Offsite staging: `/var/backups/offsite` → rclone/rsync to external target

Read-only Snapshot Policy:
 - Alle Agenten dürfen Snapshots ausschließlich via `orchestratorctl backup-mount --mode=ro` einbinden.
 - WebUI-/Experience-Teams nutzen sie, um ältere HTML/CSS-Layouts oder UX-Flows zu analysieren, bevor neue Designs ausgerollt werden.
 - API-/Service-Agenten vergleichen dort frühere Code- und Config-Stände, um Regressionsfixes vorzubereiten.
 - Jede Mount-Session wird in `/var/lib/orchestrator/backup-access.log` dokumentiert; Schreibversuche oder direkte Änderungen führen zu `policy:backup-ro-violation`.
```
## 3. Data & Control Flows
```
Users → WebUI/Mobile → API Gateway
                   │JWT/CSR
                   ▼
               Auth Service ──► Token Store (Redis/Postgres)
                   │
                   ├─► File Service ──► Storage Abstraction ──► /mnt/raid/data
                   │                 (validatePath, quota, trash sandbox)
                   │
                   ├─► Share/Link Service ──► DB (`shares`, `favorites`)
                   │
                   └─► Backup Scheduler (Future) ──► cron/Queue → Orchestrator triggers jobs
```
## 4. Security & Governance
```
Break-Glass Flow:
   Agent → privilege_elevation_request → Orchestrator Approval
      → sudo scope unlocked (≤2h) → auto rollback → log in /var/log/privilege-override.log
Package Flow:
   Agent request → /var/lib/orchestrator/package-queue/
      → SystemSetupAgent validates (CVE, signature) → decision logged
      → install via allowlisted scripts → Audit log append-only (`chattr +a`)
Audit Stack:
   - /var/log/package-approvals.log (JSONL)
   - /var/log/package-installs.log
   - /var/log/source-fetch.log
   - Incident tickets: /var/lib/orchestrator/incidents/<id>.json
Status & Reporting Policy:
   - Jeder Agent dokumentiert jede abgeschlossene Aufgabe (Task/Story) in einem nummerierten Markdown-Log im Verzeichnis `status/`.
   - Für Details zur Benennungskonvention, Nutzung des Helper-Skripts und anderen Konventionen, siehe den **Entwicklungs-Leitfaden**.
   - Orchestrator prüft auf neue Files; fehlende Updates blockieren Folgeaufgaben.
   - Offene/unterbrochene Tasks wandern automatisch in `status/backlog/<YYYYMMDD>_<agent>_<task>.md`. Dieses Backlog dient allen Agenten als Pflichtlektüre vor dem Start neuer Arbeiten.
   - **Unklarheiten/Halluzinationen:** Sobald ein Agent den Kontext nicht mehr versteht, widersprüchliche Informationen findet oder „Halluzinationen“ vermutet, wird der Task sofort pausiert und ein Request `clarification_needed` an den Orchestrator gesendet. Der Agent darf nicht raten. Stattdessen:
     1. Log-Eintrag in `status/<agent>/<NNN>_<date>_clarification.md` mit präziser Frage.
     2. Orchestrator weist Task neu zu oder liefert Kontext (z. B. Referenz-Doc).
     3. Orchestrator informiert den Owner (`user`) automatisch über den offenen Punkt, sammelt die Klarstellung und verteilt sie an alle betroffenen Agenten.
     4. Erst nach schriftlicher Bestätigung wird weitergearbeitet.
```
## 5. Filesystem Layout (ASCII Map)
```
/
├── etc/
│   ├── nas/
│   │   ├── package-allowlist.yaml
│   │   ├── agents-config.yaml
│   │   └── policies/*.yaml
│   └── sudoers.d/
│       ├── storage-ops
│       ├── network-security
│       └── monitorbot
├── var/
│   ├── lib/
│   │   ├── orchestrator/{events,package-queue,incidents}
│   │   ├── monitoring/{prometheus,grafana,loki,tempo}
│   │   ├── ai/{models,vector-db,cache}
│   │   └── backups/{jobs,history,offsite-staging}
│   ├── log/
│   │   ├── agent-orchestrator/events.log
│   │   ├── package-{approvals,installs}.log
│   │   ├── source-fetch.log
│   │   └── privilege-override.log
│   └── backups/offsite/
├── srv/
│   ├── webui/
│   ├── api/
│   └── policy-repo/
└── mnt/
    └── raid/
        ├── data/
        ├── backups/
        └── snapshots/
```
## 6. WebUI Architecture Blueprint
Details ausgelagert nach `Blueprint_WebUI.md` (Layer, Module, Datenflüsse, Alert Center, Page Blueprints).

## 7. API & Event Contracts
### 7.1 REST/GraphQL Contracts pro Modul
|    Modul     | Endpoints | Payload Highlights | Contract Tests |
|--------------|-----------|--------------------|----------------|
| Auth         | `POST /auth/login`, `POST /auth/register`, `POST /auth/refresh`, `POST /auth/passkey/challenge` | Access/Refresh JWT mit `token_type`, Device Fingerprint, optional Passkey Assertion | Playwright API tests + `contract/auth.postman.json` |
| Files        | `GET/POST/PUT/DELETE /files`, `POST /files/zip`, `POST /files/unzip`, `GET /files/thumbnail` | `path`, `operation`, `sanitizeTrashName`, `etag`, `favorite_id` | Go contract tests (`handlers/files_contract_test.go`), JSON Schema für responses |
| Favorites    | `GET/POST/DELETE /favorites` | `folder_path`, `display_name`, `user_id` | Jest contracts gegen mocked API |
| Backup       | `GET/POST /backups`, `POST /backups/:id/{run,toggle}`, `GET /backups/history`, `POST /backups/history/:id/restore` | Cron/Schedule templates, retention policies, job status | k6 smoke + JSON schema validation |
| Storage      | `GET /storage/{overview,usage,directory-size,alerts,trends}` | `used_percent`, `file_type_distribution`, `alert.level` | Dredd tests to ensure docs ↔ API parity |
| Users        | `GET/POST/PUT/DELETE /users`, `POST /users/:id/reset-password` | Role list, MFA flags, forced reset tokens | Pact contracts (WebUI ↔ API) |
| Shares       | `GET/POST/PUT/DELETE /shares` | `type` enum (smb/nfs/ftp), `allowed_clients`, `share_id` | Jest contract suite |
| DocsTerminal | `GET/POST /docs-terminal/*`, `GET/PUT /docs-settings` | Command stream, session token, remember_session flag | Integration test hitting CLI mock |

### 7.2 Event & WebSocket Topics
|            Topic   |             Publisher           |                  Payload            |              WebUI Subscriber          |
|--------------------|---------------------------------|-------------------------------------|----------------------------------------|
| `files:progress`   | FileService + Orchestrator | `{path, op, percent, user_id}`      | FileUpload/FileDownload widgets        |
| `files:favorites`  | FavoritesService                | `{favorite_id, action}`             | Favorites page + Files quick badges    |
| `backups:jobs`     | BackupAgent Scheduler (Future) | `{job_id, status, started_at, eta}` | Backup page timeline                   |
| `storage:alerts`   | Orchestrator                 | `{level, message, action}`          | Storage page alerts banner             |
| `security:sessions`| AuthService                     | `{user_id, device_id, action}`      | Settings > Security + Profile activity |
| `ai:search`        | AIKnowledgeAgent (Future)       | `{query_id, status, facets}`        | AI Lens module                         |
| `docs:terminal`    | DocumentationAgent              | `stream` lines                      | DocsTerminal websocket                 |
> Alle Topics werden im API SDK (`services/api/ws.ts`) zentral registriert; bei neuen Topics muss eine ADR + Schema-Datei (`/srv/webui/contracts/ws/<topic>.json`) erstellt werden.

## 8. State Management, Testing & QA Gates
### 8.1 Zustand/Store Architektur
```
src/state
├── auth.store.ts        # session, tokens, devices
├── files.store.ts       # currentPath, selection, uploads
├── backups.store.ts     # jobs cache, filters, websocket data
├── notifications.store.ts
└── ui.store.ts          # theme, layout, modals
```
- Jede Store-Datei exportiert Hooks (`useFilesStore`) und Actions; persistente Stores nutzen `zustand/persist` mit IndexedDB (verschlüsselt via WebCrypto).
- Stores loggen State-Transitions (DEV only) an `/var/lib/observability/loki` via console intercept.

### 8.2 Test-Pyramide
| Ebene | Tools | Scope |
|-------|-------|-------|
| Unit | Vitest + React Testing Library | Hooks/Stores/Components |
| Contract/API | Pact + Dredd + Postman | Response Schemas, error codes |
| Integration | Cypress component tests | Module interactions (Files, Backup) |
| E2E | Playwright (desktop/mobile), Lighthouse CI | Auth, Files CRUD, Backup run, Docs Terminal |
| Performance | k6 browser, Web Vitals | P95 TTI < 2.5s, bundle size gates |
| Accessibility | axe-core + Storybook checks | WCAG 2.1 AA |

### 8.3 QA Gates
- **Security Gate**: Static analysis (ESLint security rules), dependency scans, CSP regression tests.
- **Observability Gate**: Build fails if trace/span injection fehlt (CI verifiziert `x-trace-id` pro Route).
- **Offline/PWA Gate**: Workbox-Cache + Playwright Offline-Test müssen grün sein.

## 9. Offline & PWA Strategy
- Service Worker (Workbox) precached Shell + Files Skeleton; Background Sync für Upload-Queue.
- IndexedDB Cache speichert letzte 100 API Responses; Revalidate via ETag/`If-None-Match`.
- Offline Banner + Retry-Actions in Files/Backup; `useOnlineStatus` dispatcht Events an Stores.
- PWA Manifest mit Shortcuts (Files, Backup, AI Lens); `beforeinstallprompt` im App Shell.
- Pending Ops landen in `/srv/webui/state/local-ops.json` und werden nach Reconnect über `POST /pending-operations` synchronisiert.

## 10. Agent Incident Response Playbook
```
Trigger: Agent meldet `status=failed` ODER Orchestrator feuert Alert (critical)
1. Contain
   - Orchestrator setzt Agent-Status auf `paused`.
   - Incident-Datei anlegen: `/var/lib/orchestrator/incidents/<agent>-<timestamp>.json`.
2. Collect
   - Logs sichern: `journalctl -u <agent>` → `/var/log/agent-errors/<agent>/<ts>.log`.
   - Telemetrie-Snapshot (node_exporter metrics, trace IDs).
   - Automatischer Self-Test `healthcheck.sh` bzw. `make test:agent-<name>`.
3. Notify
   - WebUI Alert Center empfängt `agent:failure` Event.
   - Orchestrator verschickt Pager (Telegram/E-Mail) mit Incident-ID.
4. Diagnose & Fix
   - Owner analysiert Incident-JSON (cause, exit code, letzte Aktion).
   - Bei Config-Drift → SystemSetupAgent `drift-check`.
   - Bei Regression → Git Diff + CI-Rerun (inkl. relevanter Tests).
5. Recover
   - Nach Fix `orchestratorctl resume <agent>` ausführen und Self-Test bestehen lassen.
   - Incident abschließen via `/var/lib/orchestrator/incidents/<id>-resolved.json`.
6. Learn
   - ADR/Runbook aktualisieren.
   - Regressionstest (Unit/E2E) hinzufügen, referenziert Incident-ID.
```
- Jeder Agent liefert ein `healthcheck.sh` für automatische Verifikation.
- Incident-Logs sind append-only (`chattr +a`) und werden im täglichen Backup unter `/mnt/raid/backups/history/incidents` archiviert.

## 11. Next Features & Initiatives
> **Security/Triage Gate:** Vor jedem neuen Task `CVE_CHECKLIST.md`, offene Incidents sowie das Security-Gate „AUTH-N-01“ prüfen. Keine Umsetzung, solange kritische CVEs (CVSS ≥ 7) offen sind oder das Telemetry-Gate (Minimal-Monitoring) nicht erfüllt ist.
> **Orchestrator Policy:** Neue Tickets werden nur gestartet, wenn `status/backlog/*.md` abgearbeitet ist und der Agent sein letztes Statuslog aktualisiert hat.

### 11.1 Archivierte Kernfeatures
- Dateifunktionen: ZIP/Unzip, Suche, Sortierung, Bulk-Operationen, Favoriten (inkl. Realtime).
- UX-Grundlage: Toast Center, Liquid-Glass Theme, Upload-/Download-Progress.
- Operations: Backup- und Storage-System (Scheduler, Alerts) laufen stabil.

### 11.2 Priorität 🔴 – Security & Observability
| Nr. | Feature | Owner | Abhängigkeiten | Deliverables / Referenzen |
|-----|---------|-------|----------------|---------------------------|
| S1 | CVE-Hardening Sprint | APIAgent + SystemSetupAgent | `CVE-2025-XXXX`, `CVE-2024-ABC1` (siehe `CVE_CHECKLIST.md`) | Auth/JWT Fixes, Secret-Rotation, Openssl Updates, PentesterAgent Nachweis |
| S2 | Security Control Center | WebUIAgent + APIAgent | S1, `Blueprint_WebUI.md` §7 | Sessions/Token UI, `/auth/sessions`, Device Revoke, Trace Hooks |
| S3 | Monitoring & Alert Hub | Orchestrator + WebUIAgent | Telemetry Gate, node_exporter+Prometheus | Dashboard Widgets, WS Topics (`storage:alerts`, `backups:jobs`, `security:sessions`), Drawer/Acknowledge Flow |
| S4 | Incident Automation & Clarification Flow | Orchestrator | `Blueprint_WebUI.md` (Alert Surface) + NAS Blueprint §10 | Auto `status/backlog/*.md`, `clarification_needed` Workflow, CLI |

### 11.3 Priorität 🟠 – Operations & UX
| Nr. | Feature | Owner | Dependencies | Deliverables |
|-----|---------|-------|--------------|-------------|
| O1 | Settings – Security Tab 2.0 | WebUIAgent + APIAgent | S1 abgeschlossen | `/settings/security` GET/PUT, UI für Policy, Timeout, Rate Limits, 2FA Defaults |
| O2 | Session Timeline Component | WebUIAgent | O1, `/auth/audit` | Reusable Timeline, Filter, CSV Export |
| O3 | Observability Quick Actions | Orchestrator | Telemetry Stack live | Buttons Prometheus/Grafana/Logs, Service Status Cards, Admin Restart |
| O4 | Dashboard Refresh (Phase 1) | WebUIAgent | Monitoring Hub Data | Storage/Activity Widgets, Quick Stats, Alert Pill Integration |

### 11.4 Priorität 🟡 – Enablement & Guided Experience
| Nr. | Feature | Owner | Notes |
|-----|---------|-------|-------|
| E1 | API Token Playground | APIAgent + DocumentationAgent | Dev-Panel mit Beispiel-Calls, JWT Decoder, Snippets |
| E2 | Guided Setup Cards | DocumentationAgent + WebUIAgent | Onboarding Stepper („Harden Auth“, „Monitoring“, „Pentest“), nutzt `status/backlog` |
| E3 | Notification Routing v1 | WebUIAgent + Orchestrator | Mapping neuer Alerts auf Toast/Drawer/Push, Logging im Orchestrator |
| E4 | AI Lens Beta (Future) | (TBD) | `/api/v1/ai/search`, read-only Facetten |

### 11.5 Archiv & Nachzügler
- Legacy-Spezifikationen (Storage Cleanup, Profile 2.0, altes Dashboard) liegen in `docs/archive/legacy-webui.md`. Nur reaktivieren, wenn Roadmap es verlangt.

### 11.6 Umsetzung & Reporting
1. Feature wählen → Gates/CVEs prüfen.
2. Task-File `status/<agent>/<NNN>_<YYYYMMDD>_<feature>.md` anlegen (Links zu CVE/Blueprint).
3. Implementieren + testen + deployen.
4. Statuslog aktualisieren, `CVE_CHECKLIST.md` (falls relevant) pflegen.
5. Orchestrator aktualisiert `docs/planning/MASTER_ROADMAP.md`.

## 12. References & Mandatory Reading
- **`docs/planning/MASTER_ROADMAP.md`** – Aktueller Status je Agent, Phasenplan und Gates.
- **`docs/planning/AGENT_MATRIX.md`** – Rollenbeschreibung, Abhängigkeiten, Hardwarestatus.
- **`docs/development/DEV_GUIDE.md`** – Anleitung für das Entwicklungs-Setup und Beiträge.
- **`docs/security/SECURITY_HANDBOOK.md`** – Zentrale Richtlinie für Sicherheit und Geheimnisse.
- **`docs/policies/orchestrator-collaboration.md`** – Kollaborations- und Ticket-Workflow Agenten ↔ Orchestrator.
- **`docs/policies/systemsetup-allowlist.md`** – Governance & Paketprozesse.
- `CVE_CHECKLIST.md` – Vollständige Liste offener/geschlossener Schwachstellen inkl. Verantwortlicher und Nachweise; muss vor Security-/Release-Gates geprüft werden.
- Jeder Agent muss vor Arbeitsbeginn diese Referenzen prüfen und in seinem Status-Log verlinken, damit nachvollziehbar ist, auf welchem Stand er gearbeitet hat.
> Dieses Markdown-Dokument fasst die geplante Struktur, Kommunikationswege und Dateisystem-Governance des NAS.AI-Systems zusammen und zeigt, wie klassische Agenten mit den KI-Modulen gekoppelt werden.
