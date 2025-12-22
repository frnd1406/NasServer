# NAS.AI – WebUI Blueprint

## 1. Layered Composition
```
┌─────────────────────────────────────────────────────────────┐
│ App Shell (Layout, Theming, Router, Error Boundaries)       │
├─────────────────────────────────────────────────────────────┤
│ Experience Modules                                          │
│ ┌────────────┬────────────┬────────────┬────────────┬──────┐ │
│ │ Files      │ Backups    │ Monitoring │ Security   │ AI   │ │
│ │ Dashboard  │ Scheduler  │ & Alerts   │ Sessions   │ Lens │ │
│ └────────────┴────────────┴────────────┴────────────┴──────┘ │
├─────────────────────────────────────────────────────────────┤
│ Shared Services (API SDK, WS Client, Auth Context, i18n,    │
│ Feature Flags, Telemetry Hooks)                             │
├─────────────────────────────────────────────────────────────┤
│ UI System (“Nebula 2.0”: Tokens, Primitives, Components,    │
│ Motion)                                                     │
└─────────────────────────────────────────────────────────────┘
```
> Die Module "Monitoring", "Security" und "AI" stellen hier High-Level-Konzepte dar und sind nicht direkt an die früheren spezifischen Agenten gebunden. Die zugrunde liegende Logik wird von den Core-Agenten und dem Orchestrator bereitgestellt.

## 2. Modul- und Verzeichnisstruktur (`/srv/webui/src`)
```
src
├── app/
│   ├── main.tsx            # Vite entry
│   ├── routes.tsx          # Route definitions (lazy)
│   ├── AppShell.tsx        # Layout, Suspense boundaries
│   └── providers/          # QueryClient, Auth, Theme, Telemetry
├── core/
│   ├── components/         # Button, Card, Input, Modal, Toast, Skeleton
│   ├── hooks/              # useDebounce, useHotkeys, useMediaQuery
│   └── theme/              # tokens.ts, global.css (Nebula)
├── services/
│   ├── api/
│   │   ├── client.ts       # Axios/Fetch wrapper with interceptors
│   │   ├── queries/        # TanStack Query definitions
│   │   └── ws.ts           # Authenticated WebSocket hub
│   ├── analytics/otel.ts   # Frontend traces/log events
│   └── feature-flags.ts
├── modules/
│   ├── files/
│   │   ├── routes.tsx
│   │   ├── components/     # Tree, Grid, Preview, Actions
│   │   └── state/          # Zustand slice (drag/drop, selection)
│   ├── backups/
│   ├── monitoring/
│   ├── security-center/
│   └── ai-lens/            # Semantic/visual search UI
├── widgets/
│   ├── status-bar/
│   ├── toast-center/
│   └── realtime-indicator/
└── assets/
    └── icons, illustrations, fonts
```

## 3. Datenfluss & Kommunikation
```
Component → useQuery/useMutation ──► API SDK ──► REST/GraphQL
Component → useRealtimeChannel ──► WS Client ──► Event Hub (favorites, backups, AI jobs)
Component → useTelemetry() ──► OTel exporter ──► /var/lib/observability/tempo
Auth Context
│
├─ Issues access/refresh tokens via Auth Service
├─ Persists session in secure storage (WebCrypto)
└─ Broadcasts session state to tabs via BroadcastChannel
Feature Flags
│
├─ Source of truth: `/srv/policy-repo/feature-flags.json`
└─ Cached in IndexedDB with signature verification
```

## 4. Build-, Test- & Deploy-Pipeline
```
Developer Commit
   │
   ├─ Lint (eslint, stylelint) + Unit Tests (vitest) + Storybook snapshots
   ├─ e2e smoke (Playwright) against API staging
   └─ Bundle Analyzer (≤ 500 KB initial JS target)
       ▼
CI Build (Vite) → artifact in `/srv/webui/dist`
       ▼
Security Scan (npm audit, dependency-check)
       ▼
Upload to `/var/www/webui/` via Orchestrator/SystemSetupAgent
       ▼
Post-deploy hook publishes version + integrity hash to `agent-orchestrator/events.log`
```

## 5. Integration mit AI & Observability
- AI Lens Modul nutzt den Semantic-Search Endpoint (`/api/v1/ai/search`) sowie Bild-/OCR-Pipelines; Ergebnisse werden im Files-Modul als Facetten angezeigt.
- Telemetry Hooks fügen jedem Request Trace-IDs hinzu (`x-trace-id` aus OTel Context) und schicken Frontend-Logs an Loki (label `source=webui`).
- WebSocket-Kanäle sind nach Feature gruppiert (`files:*`, `backups:*`, `ai:*`). Auth erfolgt über kurzlebige WS-Tokens, die der API SDK rotiert.
- Dark-/Light-Theme sowie Accessibility-Einstellungen werden in `/srv/webui/state/preferences.json` persistiert und beim Login mit dem Benutzerprofil synchronisiert.

## 6. Page Blueprints & Verbesserungen
| Page | Zweck & Hauptdaten | Muss-Komponenten | Verbesserungsbedarf |
|------|-------------------|------------------|---------------------|
| **Auth (Login/Register/Forgot/Reset)** (siehe `Blueprint_WebUI_Auth.md`) | Einstieg, Session-Aufbau via `/auth/*` | Nebula Auth Card, Passkey/2FA Widgets, Passwordless Option | Strengere Validierung, Passkey/WebAuthn Flows, Rate-Limit/Lockout UI, Kontextmeldungen aus Security Agents |
| **Files (`pages/Files.jsx`)** (siehe `Blueprint_WebUI_Files.md`) | Primäre Dateiverwaltung via `/files`, `/favorites`, `/shares`, WS Realtime | FileActions, Upload/Download Panels, Virtualized FileList, Modal Stack (Rename/Share/Encrypt/Trash) | Deep-Linking per URL (`?path=/foo`), Server-seitige Favoriten statt `localStorage`, WebSocket Sync für progress, Quick Look Preview (PDF/media), Access-Control Surface (role badges), Breadcrumb + search suggestions |
| **Favorites (`pages/Favorites.jsx`)** (TODO blueprint) | Schneller Zugriff auf Ordner | Favoriten-Grid, Empty-State CTA | Persistenz via API (`/favorites`), Filter/Tags, Integration mit Files (direktes Öffnen), Multi-Device Sync, Audit-Log der Änderungen |
| **Backup (`pages/Backup.jsx`)** (siehe `Blueprint_WebUI_Backup.md`) | Job- und History-Management (`/backups`, `/history`) | Job Cards/Table, Create/Edit Modal, History Timeline, Stats Widget | Validation für Cron/Schedule, Job Templates (Full/Incremental/Diff, Docker/DB Presets), Inline Status (progress via WS), Runbooks-Verlinkung, DR-Drill Score, Graphen (Success/Failure Rate) |
| **Storage (`pages/Storage.jsx`)** (siehe `Blueprint_WebUI_Storage.md`) | Storage KPIs (`/storage/*`) | Alerts Panel, Usage Cards, Trend Bars, Directory-Size Modal | Charting Lib (Sparkline/Area), Filter per mount, SMART/RAID Status Widget, Integration mit Backup Alerts, Export (CSV/PDF) |
| **Users (`pages/Users.jsx`)** (siehe `Blueprint_WebUI_Users.md`) | User CRUD (`/users`) | Table/Grid, Create/Edit Modal, Role Selector, MFA Badge | Such-/Filterzeile, Pagination, Role-Based Access Matrix, Forced Password Reset, Hardware Key enrollment, Activity Feed (last login/IP) |
| **Shares (`pages/Shares.jsx`)** (siehe `Blueprint_WebUI_Shares.md`) | SMB/NFS/FTP Freigaben (`/shares`) | Share Cards, Protocol badges, Modal mit Form | Validation gegen Pfad Allowlist, Preview ACL (allowed hosts/users), Quick Links (Mount commands, QR), Versionierung (history), Integration mit **Orchestrator** Telemetrie |
| **Profile (`pages/Profile.jsx`)** (siehe `Blueprint_WebUI_Profile.md`) | User-spezifische Infos | Avatar Card, Stats, Security Block, Activity List | Daten aus `/profile` statt Mock, 2FA Management CTA, Device Sessions, API Tokens, Storage Quota Nutzungsanzeige |
| **Settings (`pages/Settings.jsx`)** (siehe `Blueprint_WebUI_Settings_Expert.md`) | Personalisierte + System-Settings (`/settings`, `/docs-settings`) | Tabbed Panel (Account/Security/Preferences/Notifications/Tools/Docs) | Persistenz via API (derzeit TODO), Feature Flags UI (AB Tests), Notification Channels (Alert routing), Secrets (App Passwords), Import/Export Settings |
| **Backup/Storage Widgets** | Cross-page components (status cards, badges) | Shared widget lib unter `src/widgets/` | Zentralisieren & Storybook-Doku, Skeleton States, Observability hooks (component mount/unmount logs) |
| **Mobile Views (`*-Mobile.css`)** | Responsive Styles | CSS Modules (Grid→Stack) | Consolidate via CSS-in-JS or tokens, ensure accessible hit targets, offline banners für PWA |
> Empfehlung: Jede Seite erhält eine kurze ADR + QA-Checkliste (Accessibility, Performance, Security) und künftige Realtime-Events (WebSocket Topics) werden im API SDK beschrieben, damit die WebUI strukturiert mit Orchestrator interagiert.

## 7. Global Alert Surface
```
┌────────────────────────────────────────────────────────────┐
│ Header Alert Pill  ──┐                                     │
│                      │                                     │
│ Alert Center Drawer ◄┴─ aggregates                         │
│                      │    - storage:alerts (Orchestrator)│
│                      │    - backups:jobs (Orchestrator)      │
│                      │    - security:sessions (AuthService) │
│ Toast Stream ◄───────┘    - service-status (Orchestrator)   │
└────────────────────────────────────────────────────────────┘
```
- **Alert Pill** im App-Header zeigt Anzahl kritischer/offener Meldungen (Severity-Farben). Klick öffnet Drawer.
- **Alert Center Drawer** (`/srv/webui/src/modules/alert-center/`) listet Alerts chronologisch mit Quelle, Agent, Aktion (Runbook Link). Filter nach Severity, Agent, Zeit.
- **Toast Stream** für Sofortmeldungen; persistente Alerts landen zusätzlich im Drawer, bis vom User quittiert oder automatisch resolved (per API).
- **Data Source**: WebSocket Topics `storage:alerts`, `backups:jobs`, `security:sessions`, `service:status`. Fallback Poll über `/alerts/feed` Endpoint falls WS offline.
- **Display Policy**: Alle Alerts mit `level ∈ {warning, critical}` müssen global sichtbar sein, unabhängig davon, welche Page geöffnet ist. Informative (`info`) Alerts erscheinen als Badge im jeweiligen Modul, dürfen aber wegfilterbar sein.
- **Acknowledgement Flow**: User quittiert Alerts via `POST /alerts/:id/ack`; WebUI protokolliert User/Time und blendet Meldung aus. Ungelesene Alerts werden beim Logout gespeichert und beim nächsten Login erneut angezeigt.

## 8. Auth Module Deep Dive
Details ausgelagert nach `Blueprint_WebUI_Auth.md` für Use Cases, Komponenten, Tests und Telemetrie.