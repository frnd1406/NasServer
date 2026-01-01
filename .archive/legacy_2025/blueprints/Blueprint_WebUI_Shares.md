# NAS.AI – WebUI Shares Blueprint

## 1. Scope & Responsibilities
- SMB/NFS/FTP share management (`/shares`).

## 2. UX & Layout
```
┌──────────────────────────────────────────────────────────────┐
│ Header: Shares           [Search] [Filters] [Add Share]      │
├─────────────┬────────────────────────────────────────────────┤
│ Sidebar     │ Cards List                                     │
│ [← Back]    │ ┌───────────────────────────────────────────┐  │
│ [Home]      │ │ ShareCard (name, path, protocol)         │  │
│ Filters     │ │ Status badge + action buttons            │  │
│  Protocol   │ └───────────────────────────────────────────┘  │
│  Enabled    │ Modal overlay for create/edit.                 │
├─────────────┴────────────────────────────────────────────────┤
│ Footer: export config • service status (SMB/NFS/FTP)         │
└──────────────────────────────────────────────────────────────┘
```
## 3. Data Flow
- `/shares` CRUD, `POST /shares/test`, `/shares/export`.
- Service status via `/services/status` or WS topic `service:status` (published by **Orchestrator**).

## 4. Components
- ShareCard, ShareFormModal, ProtocolBadge, QuickActions.
- hooks: `useShares`, `useShareForm`.

## 5. Validation
- Path allowlist, protocol-specific constraints.
- Enabled toggle requires confirmation.

## 6. Tests/Telemetry
- Unit tests for form validation.
- Playwright add/edit/delete share.
- Events: `share_created`, `share_service_restart`.

## 7. References
- `NAS_AI_SYSTEM.md`, `Blueprint_WebUI_Files.md`, `CVE_CHECKLIST.md`.