# LESSONS LEARNED - NAS WebUI Project
**Date:** 2025-11-21
**Session:** Dashboard & Files Page Development

---

## 🎯 Was funktioniert hat

### Backend API (Go)
1. **File-API erfolgreich implementiert**
   - Endpoints: List, Upload, Download, Delete, Rename, CreateDirectory
   - Implementierung: Siehe `docs/development/REFERENCE_SNIPPETS.md` (Go File Handler)
   - Base directory konfigurierbar via `FILES_BASE_DIR` env var

2. **CORS Middleware funktioniert**
   - Implementierung: Siehe `docs/development/REFERENCE_SNIPPETS.md` (Go CORS Middleware)

3. **Path Sanitization gegen Directory Traversal**
   - Implementierung: Siehe `docs/development/REFERENCE_SNIPPETS.md` (Go File Handler)

### Frontend (React + TypeScript)
1. **Zustand State Management** funktioniert gut für Files
   - Implementierung: Siehe `docs/development/REFERENCE_SNIPPETS.md` (TypeScript Files Store)

2. **Vite Dev Server** - automatisches Hot Reload

---

## ❌ Probleme & Erkenntnisse

### 1. **Auth-Problem - KRITISCH**
**Problem:** Neue API hatte keine Auth-Endpoints (Login/Register)
- Login versuchte Port 8080 zu erreichen
- Neue API hatte nur File-Endpoints
- Alte Production API lief bereits auf Port 8080 mit Auth

**Lesson:**
- **Immer zuerst Auth-Endpoints implementieren!** (Siehe `docs/blueprints/Blueprint_WebUI_Auth.md`)
- Oder Auth temporär deaktivieren für Development

### 2. **Port-Konflikte**
- Mehrere API-Instanzen liefen gleichzeitig
- Port 8080, 8081, 5173, 5174, 5175 alle belegt
- Schwer zu debuggen welche API wohin zeigt

**Lesson:**
- Alle Prozesse vor Neustart killen
- Ein einziger Port für API nutzen
- `.env.local` für Frontend API-URL nutzen

### 3. **Frontend-Pages Chaos**
- Dashboard wurde mehrmals umgeschrieben
- Files-Page mit/ohne Design durcheinander
- Zu viele Iterationen ohne klaren Plan

**Lesson:**
- **Erst minimale Version → dann Schritt für Schritt erweitern**
- Ein Feature nach dem anderen
- Nicht gleichzeitig Design + Funktionen

### 4. **Database nicht notwendig für Files**
- API versuchte DB zu connecten
- Files brauchen keine DB
- Auth würde DB brauchen

**Lesson:**
- Optional dependencies klar trennen
- DB nur wenn wirklich nötig

---

## 💡 Best Practices für Neustart

### Development Workflow
1. **Start with minimal versions:**
   - Empty Dashboard → Add username → Add logout → Add links
   - Empty Files page → List files → Add folder navigation → Add create folder

2. **One feature at a time:**
   - ✅ Files List working → dann erst Upload
   - ✅ Login working → dann erst Files

3. **Backend first, then Frontend:**
   - API Endpoints testen mit `curl` BEVOR Frontend
   - Swagger/Postman Collection wäre hilfreich

### Architecture
1. **API Structure:**
   - Siehe `docs/blueprints/Blueprint_WebUI.md` und `NAS_AI_SYSTEM.md` für die API-Struktur.

2. **Frontend Structure:**
   - Siehe `docs/blueprints/Blueprint_WebUI.md` für die Frontend-Struktur.

3. **Environment:**
   ```bash
   # API
   PORT=8080
   JWT_SECRET=xxx
   FILES_BASE_DIR=/data

   # Frontend .env.local
   VITE_API_URL=http://localhost:8080
   ```

---

## 🔧 Code Snippets die gut funktionieren

### File API Handler (Go)
```go
// Siehe docs/development/REFERENCE_SNIPPETS.md für Details
```

### Zustand Store (TypeScript)
```typescript
const useFilesStore = create<FilesStore>((set, get) => ({
    files: [],
    currentPath: '/',
    fetchFiles: async (path?: string) => {
        const response = await apiClient.get('/api/v1/files', {
            params: { path: path || '/' }
        });
        set({ files: response.data.files });
    }
}));
```

---

## 📋 Empfehlungen für nächste Session

### Phase 1: Setup
1. Clean environment - alle Prozesse stoppen
2. Eine API-Instanz auf Port 8080
3. WebUI Dev Server auf Port 5173

### Phase 2: Backend
1. Auth-Endpoints implementieren (Login, Register, Refresh) - siehe `docs/blueprints/Blueprint_WebUI_Auth.md`
2. JWT Middleware für protected routes
3. File-Endpoints (bereits vorhanden)

### Phase 3: Frontend
1. Login-Page (bereits vorhanden mit Design)
2. Minimal Dashboard (nur Text, Link, Logout)
3. Minimal Files (nur List + Navigation)
4. Dann erweitern: Upload, Delete, etc.

### Phase 4: Integration
1. Test Login → Dashboard flow
2. Test Dashboard → Files flow
3. Test File operations

---

## 🗂️ File Structure (BEHALTEN)
```
/
├── infrastructure/
│   ├── api/
│   │   ├── src/
│   │   │   ├── handlers/
│   │   │   │   ├── files.go      ✅ KEEP (Referenz in docs/development/REFERENCE_SNIPPETS.md)
│   │   │   │   ├── system.go
│   │   │   │   └── websocket.go
│   │   │   ├── main.go           ✅ KEEP (mit CORS)
│   │   │   ├── vault/
│   │   │   ├── logger/
│   │   │   └── health/
│   │   └── bin/
│   └── webui/
│       ├── src/
│       │   ├── pages/
│       │   │   ├── Login.tsx     ✅ KEEP (hat Design)
│       │   │   ├── Dashboard.tsx (minimal)
│   │   │   │   └── Files.tsx     (minimal)
│       │   ├── state/
│       │   │   ├── auth.ts
│       │   │   └── files.ts      ✅ KEEP
│       │   └── services/api/
│       │       ├── client.ts     ✅ KEEP (hat CORS Config)
│       │       └── files.ts      ✅ KEEP
│       └── .env.local            ✅ KEEP
└── docs/
    └── LESSONS-LEARNED.md        ✅ DIESE DATEI
```

---

## ⚠️ Was beim Cleanup LÖSCHEN
- Alle .backup Dateien
- Temporäre Daten in `/home/user/Agent/infrastructure/data/`
- `node_modules/.cache/**`
- Alte API binaries (rebuild fresh)

## ✅ Was BEHALTEN

1. Source Code Structure
2. Go modules (go.mod, go.sum)
3. package.json, package-lock.json
4. Dieser LESSONS-LEARNED.md
5. Referenzierte Dokumente (z.B. `docs/planning/MASTER_ROADMAP.md`, `docs/planning/AGENT_MATRIX.md`)

---

**Ende der Dokumentation**