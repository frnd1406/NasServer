# NAS.AI – WebUI Files Blueprint

## 1. Scope & Responsibilities
- Hauptansicht für Datei-/Ordner-Verwaltung (`/files`).
- Multiplattform (Desktop/Tablet/Mobile) mit Drag&Drop, Kontextmenüs, WS-Realtime.

## 2. UX & Layout
- Zweispaltiges Layout: Sidebar (tree/favorites), Hauptbereich (grid/list).
- View Modes: Grid (thumbnails) & List (detailliert).
- Breadcrumb Navigation + Search Bar oben.
- Status-Leiste unten (ausgewählte Elemente, Gesamtgröße).

```
┌─────────────────────────────────────────────────────────────────────┐
│ Header: [Back] [Home] | Breadcrumb: /mnt/data/projects | Search 🔍  │
├──────────────┬──────────────────────────────────────────────────────┤
│ Sidebar      │ View Switch [Grid][List]                             │
│  [← Back]    │                                                      │
│  [Home]      │ Main Panel                                           │
│  Favorites   │ ┌────────────────────────────────────────────────┐   │
│   • /media   │ │ FileAction Toolbar (Upload, New Folder, Trash) │   │
│   • /docs    │ ├────────────────────────────────────────────────┤   │
│  Tree        │ │ Grid/List of files with icons, tags, badges    │   │
│   /          │ │ Context menu on right click, drag handles      │   │
│   └─ data    │ │ Sticky upload panel (progress) bottom right    │   │ 
│ Trash        │ └────────────────────────────────────────────────┘   │ 
├──────────────┴──────────────────────────────────────────────────────┤
│ Footer: Selected 3 items • Total 4.2 GB • WebSocket: Connected      │
└─────────────────────────────────────────────────────────────────────┘
```

## 3. Datenfluss
1. Initial Load → `GET /files?path=/...` (server returns entries, metadata, permissions).
2. Realtime Updates via WebSocket `files:progress`, `files:favorites`.
3. Favoriten/Clipboard/Trash via API (`/favorites`, `/files/trash`).
4. Uploads → `POST /files/upload` (multipart) + progress events; folder uploads nutzen HTML5 directory API.
5. Downloads → direct `/files/download?path=...` (signed URLs) + optional ZIP.
6. Actions (rename, delete, encrypt) → `POST /files/actions` (body: op, path).

## 4. Komponentenstruktur
```
modules/files/
├── components/
│   ├── FileActions.tsx (toolbar)
│   ├── FileList.tsx (virtualized list)
│   ├── FileGrid.tsx (thumbnails)
│   ├── Breadcrumbs.tsx
│   ├── FileModals/
│   │   ├── RenameModal.tsx
│   │   ├── ShareDialog.tsx
│   │   ├── EncryptDialog.tsx
│   │   └── TrashModal.tsx
│   ├── UploadPanel.tsx
│   ├── DownloadPanel.tsx
│   └── ContextMenu.tsx
├── hooks/
│   ├── useFiles.ts (state machine)
│   ├── useClipboard.ts
│   └── useDragDrop.ts
├── state/
│   └── files.store.ts (Zustand/persist)
└── tests/
    ├── files.spec.tsx (Playwright)
    └── useFiles.test.ts
```

## 5. State & Behavior
- `useFiles` kapselt: currentPath, files[], loading, selection, viewMode.
- Persistente States (viewMode, lastPath) in `files.store.ts`.
- Clipboard (copy/cut) stored in memory, expire after session.
- Drag & Drop: HTML5 events + custom overlay; drop triggers `POST /files/move`.
- Favorites: server-side via `/favorites`; UI cached lokal, re-sync per WS.

## 6. Validierung & Fehlerbehandlung
- Pfad-Sanitizing via API (server rejects invalid). UI zeigt Toast + revert.
- Delete/Encrypt nur wenn `userRole` erlaubt (aus API response).
- Upload Limit (max size, allowed extensions) aus `/settings/files`.
- Offline Mode: `useOnlineStatus` → Buttons disabled, Banner angezeigt.

## 7. Modals & Aktionen
- RenameModal: inline validation (no slash, length).
- ShareDialog: create link via `/shares` (password optional).
- EncryptDialog: KMS Integration (passphrase). Einstellungen via `/settings/encryption`.
- TrashModal: list deleted items, restore/delete permanently.
- PreviewModal: Bilder/PDF Quick View.

## 8. Tests & Telemetrie
- Unit Tests für hooks (filtering, sorting, search).
- Playwright: upload, rename, share, favorite.
- Telemetry Events: `files_upload_start/success/fail`, `files_action_error`.
- WebSocket reconnect logging.

## 9. Roadmap-Verknüpfung
- Referenziert `NAS_AI_SYSTEM.md` (Next Features) für Access-Control & Alerts.
- CVE-relevante Tasks (Path Traversal) → `CVE_CHECKLIST.md`.
- Owner laut `docs/planning/MASTER_ROADMAP.md`: WebUIAgent + APIAgent.

## 10. Referenzen
- `Blueprint_WebUI.md` (global layout, alerts).
- `Blueprint_WebUI_Auth.md` (session context).
- `CVE_CHECKLIST.md` (File-Service Findings).
- `NAS_AI_SYSTEM.md` (Status/Reporting Policies).

## 11. AI-Integration (Vision)
- **Semantic Search Panel:** Optional rechter Drawer, der Suchbegriffe an das AI-Modul sendet (`/api/v1/ai/search?path=...`). Ergebnisse als Facetten (Ähnlichkeit, Tags, OCR-Text) → Klick navigiert in Files-Ansicht.
- **Auto-Tagging Badges:** WebSocket Topic `ai:tags` liefert neue Tags pro Datei; UI zeigt Chips unterhalb der Einträge und erlaubt Filter „Tag = Vertrag“.
- **Smart Suggestions:** Suggestion-Bar (unter Search) zeigt häufige Anfragen, ähnlich wie „Zuletzt gesucht“, basierend auf anonymisierten AI-Lens Statistiken.
- **Security Constraints:** AI-Endpunkte sind read-only; Mutationen (z. B. auto-organisieren) werden erst nach explizitem Gate eingeführt. Alle AI-Aufrufe liefern Trace-IDs, damit Telemetrie (AI ↔ Files) nachvollziehbar bleibt.