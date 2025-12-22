# NAS.AI WebUI Frontend

**Owner:** WebUIAgent
**Technologie:** React 18 + Vite + TypeScript + Tailwind CSS
**Status:** Phase 1 - Foundation Setup

---

## OVERVIEW

Dies ist das React-basierte Frontend für das NAS.AI-System. Es kommuniziert mit dem Backend API via REST und WebSocket für Echtzeitfunktionen.

---

## STRUKTUR

```
webui/
├── src/
│   ├── main.tsx               # Application entry point
│   ├── App.tsx                # Root component
│   ├── pages/                 # Page components
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Dashboard.tsx
│   │   └── Files.tsx
│   ├── components/            # Reusable components
│   │   ├── auth/             # Auth-related components
│   │   ├── files/            # File browser components
│   │   └── common/           # Common UI components
│   ├── hooks/                 # Custom React hooks
│   ├── services/              # API client & WebSocket
│   │   ├── api.ts            # Axios API client
│   │   └── websocket.ts      # WebSocket client
│   ├── store/                 # Zustand state management
│   │   ├── auth.store.ts
│   │   ├── files.store.ts
│   │   └── ui.store.ts
│   ├── types/                 # TypeScript types
│   ├── utils/                 # Helper functions
│   └── styles/                # Global styles
├── public/                    # Static assets
├── tests/
│   ├── unit/                  # Component tests
│   └── e2e/                   # Playwright E2E tests
├── docs/
│   └── components.md          # Component documentation
├── package.json               # Dependencies
├── vite.config.ts             # Vite configuration
├── tsconfig.json              # TypeScript configuration
├── tailwind.config.js         # Tailwind CSS configuration
└── Dockerfile                 # Container image
```

---

## PHASE 1 TASKS (WebUIAgent Assignment)

### Epic 1: Project Setup & Foundation (3-5 Tage)

**Tasks:**
1. ⏳ Vite + React + TypeScript initialization
2. ⏳ Tailwind CSS setup
3. ⏳ ESLint + Prettier configuration
4. ⏳ Folder structure creation
5. ⏳ Basic routing (React Router v6)
6. ⏳ Build & dev scripts (package.json)

**Deliverables:**
- `package.json` with all dependencies
- Working dev server on `:5173`
- Basic routing (/, /login, /dashboard)
- Tailwind CSS working

**Acceptance Criteria:**
- `npm run dev` starts server without errors
- Navigate to http://localhost:5173
- Hot reload working
- TypeScript strict mode enabled

---

### Epic 2: API SDK & Authentication (5-7 Tage)

**Tasks:**
1. ⏳ Axios API client setup
   - Base URL configuration
   - Request/Response interceptors
   - Error handling
   - JWT token injection
2. ⏳ Auth store (Zustand)
   - Login state
   - Token management (access + refresh)
   - User profile
   - Persistence (localStorage with encryption)
3. ⏳ Auth API endpoints integration
   - POST /auth/register
   - POST /auth/login
   - POST /auth/refresh
   - POST /auth/logout
4. ⏳ Protected route wrapper
5. ⏳ Login page UI
6. ⏳ Register page UI

**Deliverables:**
- Working API client (`src/services/api.ts`)
- Auth store with all actions
- Login/Register pages
- Protected route component
- Token refresh logic

**Acceptance Criteria:**
- User can register new account
- User can login (JWT stored securely)
- Protected routes redirect to /login if unauthenticated
- Token auto-refresh before expiry
- All TypeScript types defined

**References:**
- NAS_AI_SYSTEM.md → §7.1 (Auth API Contract)
- Blueprint_WebUI_Auth.md
- DEV_GUIDE.md → §5 (API Tokens)

---

### Epic 3: File Browser Core (7-10 Tage)

**Tasks:**
1. ⏳ Files API integration
   - GET /files (list)
   - POST /files (upload)
   - GET /files/:path (download)
   - DELETE /files/:path
2. ⏳ Files store (Zustand)
   - Current directory
   - File list
   - Selection state
   - Upload progress
3. ⏳ File browser component
   - File list (table view)
   - Breadcrumb navigation
   - Context menu (right-click)
   - File/folder icons
4. ⏳ Upload component
   - Drag & drop
   - Progress bar
   - Multi-file support
5. ⏳ Basic file operations
   - Upload
   - Download
   - Delete (with confirmation)

**Deliverables:**
- Working file browser
- Upload with progress
- Download files
- Delete files

**Acceptance Criteria:**
- Can navigate folders
- Can upload files (with progress)
- Can download files
- Can delete files
- Context menu works
- All operations update UI in real-time

**References:**
- NAS_AI_SYSTEM.md → §7.1 (Files API Contract)
- Blueprint_WebUI_Files.md

---

### Epic 4: UI/UX Polish (3-5 Tage)

**Tasks:**
1. ⏳ Toast notifications (success, error, info)
2. ⏳ Loading states & skeletons
3. ⏳ Error boundaries
4. ⏳ 404 page
5. ⏳ Responsive design (mobile-first)
6. ⏳ Dark mode support (optional)

**Deliverables:**
- Toast system working
- All loading states handled
- Error boundaries catching errors
- Responsive on mobile/tablet/desktop

**Acceptance Criteria:**
- API errors show user-friendly toasts
- Loading skeletons during data fetch
- App doesn't crash on errors
- Works on 320px width (mobile)

---

## DEVELOPMENT SETUP

### Prerequisites
- Node.js 20+ installed
- npm or pnpm
- Access to Backend API (http://localhost:8080)

### Quick Start

```bash
# 1. Navigate to WebUI directory
cd /home/user/Agent/infrastructure/webui

# 2. Install dependencies
npm install

# 3. Setup environment variables
cat > .env.local <<EOF
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_BASE_URL=ws://localhost:8080
EOF

# 4. Run development server
npm run dev

# Server runs on http://localhost:5173
```

### Environment Variables

Required:
- `VITE_API_BASE_URL` - Backend API URL (e.g., http://localhost:8080)
- `VITE_WS_BASE_URL` - WebSocket URL (e.g., ws://localhost:8080)

Optional:
- `VITE_ENV` - Environment (development, production)

### Testing

```bash
# Run unit tests (Vitest)
npm run test

# Run with coverage
npm run test:coverage

# Run E2E tests (Playwright)
npm run test:e2e

# Lint code
npm run lint

# Format code
npm run format
```

### Build

```bash
# Production build
npm run build

# Preview production build
npm run preview
```

---

## TECHNOLOGY STACK

### Core
- **React 18** - UI framework
- **Vite** - Build tool & dev server
- **TypeScript** - Type safety
- **React Router v6** - Routing
- **Zustand** - State management

### Styling
- **Tailwind CSS** - Utility-first CSS
- **shadcn/ui** (Phase 2) - Component library
- **Lucide React** - Icon library

### API & Data
- **Axios** - HTTP client
- **TanStack Query** (Phase 2) - Data fetching & caching
- **WebSocket** (Phase 2) - Real-time updates

### Testing
- **Vitest** - Unit testing
- **React Testing Library** - Component testing
- **Playwright** - E2E testing

### Developer Experience
- **ESLint** - Linting
- **Prettier** - Code formatting
- **Husky** (Phase 2) - Git hooks

---

## API CLIENT STRUCTURE

```typescript
// src/services/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor (inject JWT)
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor (handle errors, refresh token)
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    // Token refresh logic here
    return Promise.reject(error);
  }
);

export default api;
```

---

## STATE MANAGEMENT

### Auth Store Example

```typescript
// src/store/auth.store.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,

      login: async (email, password) => {
        // Call API, set tokens
      },

      logout: () => {
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false });
      },

      refreshAccessToken: async () => {
        // Refresh token logic
      },
    }),
    {
      name: 'auth-storage',
      // Encrypt sensitive data in localStorage (Phase 2)
    }
  )
);
```

---

## COMPONENT GUIDELINES

### File Structure
```
components/
├── auth/
│   ├── LoginForm.tsx          # Login form component
│   ├── RegisterForm.tsx       # Register form
│   └── ProtectedRoute.tsx     # Route guard
├── files/
│   ├── FileBrowser.tsx        # Main file browser
│   ├── FileList.tsx           # File list table
│   ├── FileUpload.tsx         # Upload component
│   └── FileContextMenu.tsx    # Right-click menu
└── common/
    ├── Button.tsx             # Reusable button
    ├── Input.tsx              # Form input
    ├── Toast.tsx              # Toast notifications
    └── Loading.tsx            # Loading spinner
```

### Naming Conventions
- **Components:** PascalCase (e.g., `FileBrowser.tsx`)
- **Hooks:** camelCase with `use` prefix (e.g., `useAuth.ts`)
- **Stores:** camelCase with `.store` suffix (e.g., `auth.store.ts`)
- **Types:** PascalCase with `.types` suffix (e.g., `auth.types.ts`)

### Code Style
- Functional components only (no class components)
- TypeScript strict mode
- Props interfaces defined above component
- Hooks at top of component
- Early returns for guards

---

## SECURITY REQUIREMENTS

### MUST HAVE (Phase 1)
- ✅ JWT tokens stored securely (localStorage with encryption later)
- ✅ Token refresh before expiry
- ✅ Protected routes redirect to /login
- ✅ No sensitive data in console.log
- ✅ Input validation & sanitization

### NICE TO HAVE (Phase 2)
- CSRF token handling
- XSS prevention (sanitize user content)
- Content Security Policy headers
- HTTPS only in production

---

## BLUEPRINTS REFERENCE

The following blueprints define detailed UI/UX requirements:

- `Blueprint_WebUI.md` - Overall architecture
- `Blueprint_WebUI_Auth.md` - Login/Register flows
- `Blueprint_WebUI_Files.md` - File browser design
- `Blueprint_WebUI_Profile.md` - User profile
- `Blueprint_WebUI_Settings.md` - Settings page

Location: `/home/user/Agent/docs/blueprints/`

---

## STATUS TRACKING

**Statuslogs:** `/home/user/Agent/status/WebUIAgent/phase1/`

**Format:** `NNN_YYYYMMDD_task-description.md`

**Required for each Epic:**
1. Analysis (Plan, Dependencies, Risks)
2. Implementation logs (step-by-step progress)
3. Screenshots/Demo (evidence)
4. Final summary with artefacts

---

## NEXT STEPS

1. **WebUIAgent:** Read this README + AGENT_CHECKLIST.md
2. **WebUIAgent:** Create `status/WebUIAgent/phase1/001_20251121_epic1-setup.md`
3. **WebUIAgent:** Start Epic 1 (Project Setup)
4. **WebUIAgent:** Report to Orchestrator when Epic 1 complete

---

## REFERENCES

- **System Architecture:** `/home/user/Agent/NAS_AI_SYSTEM.md`
- **Security Handbook:** `/home/user/Agent/docs/security/SECURITY_HANDBOOK.pdf`
- **Agent Matrix:** `/home/user/Agent/docs/planning/AGENT_MATRIX.md`
- **Dev Guide:** `/home/user/Agent/docs/development/DEV_GUIDE.md`
- **Blueprints:** `/home/user/Agent/docs/blueprints/Blueprint_WebUI*.md`

---

**Assigned:** 2025-11-21
**Target:** Phase 1 Complete by 2025-11-28
**Owner:** WebUIAgent

Terminal freigegeben.
