# NAS.AI – WebUI Backup Blueprint

## 1. Scope & Responsibilities
- Verwaltung aller Backup-Jobs (`/backup`).
- Echtzeit-Status (WS), Historie, DR-Workflows (run, restore).
- Interaktion mit dem **Orchestrator** (jobs queue).

## 2. UX & Layout
- Hauptansicht: Tabbed Cards (Jobs, History, Stats).
- Jobs-Grid mit Status-Badges (running, success, failed).
- Floating Action Button „New Backup“.
- Timeline-View für Historie (Filter Status/Job).
- Sidebar Quick Filters (Backup-Typ, Zielpfad).

```
┌──────────────────────────────────────────────────────────────────┐
│ Header:  "Backup Control Center"  [Search]   [Filters]  [New+]   │
├──────────────┬───────────────────────────────────────────────────┤
│ Sidebar      │ Tabs: [Jobs] [History] [Stats]                    │
│ • Job Type   │                                                   | 
│ • Destination|                                             ┌─────▼────┐
│ • Status     │ Jobs Grid                                   │  Jobs    │
│              │ ┌──────────────────────────────────────────┐│  Tab     │
│              │ │ JobCard #1   [Running][Next: 02:00]      │└──────────┘
│              │ │ JobCard #1   [Running][Next: 02:00]      │      |
│              │ │ Buttons: Run Now | Edit | Disable        │      |
│              │ ├──────────────────────────────────────────┤      |
│              │ │ JobCard #2   [Failed][Retry CTA]         │      |
│              │ └──────────────────────────────────────────┘      |
│              │                                                   |
│              │ History Tab → Timeline Table                      |
│              │ Stats Tab  → KPIs + charts                        |
│ [← Back]     │                                                   |
│ [Home]       │                                                   |
├──────────────┴───────────────────────────────────────────────────┤
│ Footer: Last Sync • WebSocket Status • Selected Job Summary      │
└──────────────────────────────────────────────────────────────────┘
```

## 3. Datenfluss
1. `GET /backups` → Liste Jobs + Metadata (schedule, enabled, type).
2. `GET /backups/history?limit=...` → letzte Läufe mit Status.
3. `POST /backups` / `PUT /backups/:id` / `DELETE ...` → Job CRUD.
4. `POST /backups/:id/run` → Trigger + Orchestrator Event.
5. `POST /backups/:id/toggle` → enable/disable.
6. WebSocket `backups:jobs` liefert Progress & Finalstatus.
7. `POST /backups/history/:id/restore` → Restore Flow (modal).

## 4. Komponentenstruktur
```
modules/backups/
├── components/
│   ├── BackupJobCard.tsx
│   ├── BackupHistoryTable.tsx
│   ├── BackupStats.tsx
│   ├── JobFormModal.tsx
│   ├── RestoreModal.tsx
│   └── JobTimeline.tsx
├── hooks/
│   ├── useBackups.ts (query + websocket)
│   └── useBackupForm.ts
├── state/
│   └── backups.store.ts (filters, selected job)
└── tests/
    ├── backups.spec.tsx
    └── useBackups.test.ts
```

## 5. Scheduler & Validation
- Cron Templates (Daily, Weekly, Custom) via Dropdown + „advanced“ Textfeld.
- Validate Source/Destination Paths (auto suggestions, read-only).
- Backup Types (Full, Incremental, Differential) mit Tooltips.
- Retention Days Numeric Input (min/max).
- Enable/Disable toggle (auto confirm).

## 6. Realtime & Alerts
- WebSocket Progress Bar (percent, files processed).
- Failed Run → Alert Drawer + toast.
- Running Jobs heben sich visuell ab (glow + spinner).
- Stats Widget: „Last Success“, „Next Scheduled“, „Failed (24h)“.

## 7. History & Restore
- Timeline Table mit Filter (Status, Job, Date Range).
- Each row: run duration, size, files count, logs link.
- Restore Modal: Pfad-Input, Conflict Strategy (overwrite/new folder).
- After restore → summary toast + link zu Logs.

## 8. Tests & Telemetrie
- Unit: `useBackups` (transform API data, merging WS updates).
- Integration (Playwright): create job, run now, observe progress, restore.
- Telemetry Events: `backup_job_created`, `backup_run_started`, `backup_restore_failed`.
- Logs: Orchestrator Events verlinken in UI.

## 9. Roadmap & CVE Bezug
- `NAS_AI_SYSTEM.md` Next Features S1-S4 (Security/Observability) beeinflussen Backup UI (Alert Hub, Clarification Flow).
- CVE: Offene Backup/Path Issues laut `CVE_CHECKLIST.md` blockieren Deploy.
- Owner: **Orchestrator** + **WebUIAgent** (siehe `docs/planning/MASTER_ROADMAP.md`).

## 10. Referenzen
- `Blueprint_WebUI.md` (global layout, alert center).
- `Blueprint_WebUI_Auth.md` (Session Context, Admin Roles).
- `Blueprint_WebUI_Files.md` (File interactions, favorites).
- `CVE_CHECKLIST.md`, `NAS_AI_SYSTEM.md`.