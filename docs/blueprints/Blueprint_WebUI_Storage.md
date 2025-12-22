# NAS.AI – WebUI Storage Blueprint

## 1. Scope & Responsibilities
- Storage Monitoring Dashboard (`/storage`).
- Multi-path metrics, alerts, directory analysis.

## 2. UX & Layout
```
┌─────────────────────────────────────────────────────────────────┐
│ Header: Storage Status             [Refresh]  [Filters]        │
├──────────────┬─────────────────────────────────────────────────┤
│ Sidebar      │ Tabs: [Overview][Alerts][Trends][Directory]     │
│  [← Back]    │                                                 │
│  [Home]      │ Overview → Cards + charts                       │
│  Paths list  │ Alerts Tab → table                              │
│  Thresholds  │ Trends → line charts                            │
├──────────────┴─────────────────────────────────────────────────┤
│ Footer: Last Probe • Node status • Alert count                 │
└─────────────────────────────────────────────────────────────────┘
```
## 3. Data Flow
- `GET /storage/overview`, `/storage/alerts`, `/storage/trends`, `/storage/directory-size`.
- WebSocket `storage:alerts` for realtime warnings (published by **Orchestrator**).
- Directory modal `GET /storage/directory-size?path=...`.

## 4. Components
- StorageOverviewCard, AlertsTable, TrendsChart, DirectoryModal, PathFilter.
- hooks: `useStorageMetrics`, `useStorageAlerts`.
- state: `storage.store.ts` (selected path, thresholds).

## 5. Validation
- Alert thresholds from `/settings/storage`.
- Directory size requests debounced.
- Large path lists virtualized.

## 6. Tests/Telemetry
- Unit tests for hooks.
- Playwright for filter, alert ack, directory scan.
- Events: `storage_alert_ack`, `storage_directory_scan`.

## 7. References
- `NAS_AI_SYSTEM.md`, `Blueprint_WebUI.md`, `CVE_CHECKLIST.md`.