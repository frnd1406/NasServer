# NAS.AI – WebUI Auth Blueprint

## 1. Auth Module Deep Dive

### 1.1 Use Cases
1. **Login (Username/Password + optional Passkey)**  
2. **Register / Invite Acceptance (Admin-gated)**  
3. **Password Recovery (Forgot/Reset)**  
4. **Session Listing & Terminate (Security Control Center)**  
5. **Passkey/WebAuthn Enrollment (v2)**  
6. **2FA (TOTP) On/Off & Backup Codes**

### 1.2 UI Struktur
```
┌──────────────────────────────┐
│  NAS Login                   │
│  --------------------------  │
│  [ Username input       ]    │
│  [ Password input       ]    │
│  ( ) Use Passkey   ( ) 2FA   │
│  [ Login ]   [ Forgot? ]     │
│                              │
│  - Admin contact info        │
│  - Security notices          │
└──────────────────────────────┘
```


```
modules/auth/
├── components/
│   ├── AuthCard.tsx
│   ├── CredentialForm.tsx
│   ├── PasskeyButton.tsx
│   ├── TwoFactorSetup.tsx
│   └── SessionList.tsx
├── pages/
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── ForgotPassword.tsx
│   ├── ResetPassword.tsx
│   └── SecurityControlCenter.tsx
├── hooks/
│   ├── useAuthForm.ts
│   └── useSessions.ts
└── tests/
    ├── login.spec.tsx
    └── sessions.spec.tsx
```

### 1.3 Datenfluss
1. **Login**  
   - `POST /auth/login` → `{ access_token, refresh_token, token_type }`  
   - Tokens in secure storage (`IndexedDB` + `crypto.subtle`) via `AuthContext`.  
   - `AuthContext` broadcastet `login_success` Event (BroadcastChannel) → Sidebar/Alert Pill aktualisieren.  
2. **Refresh**  
   - `POST /auth/refresh` (silent) alle 10 min; Access Token (15 min TTL) wird rotiert.  
   - Refresh Token (7 Tage) Single-Use; beim Fail: Logout + Alert.  
3. **Logout**  
   - `POST /auth/logout` + local token purge.  
4. **Passkey/WebAuthn**  
   - Registration: `GET /auth/passkey/challenge` → `navigator.credentials.create`.  
   - Login: `POST /auth/passkey/verify`.  
5. **Session Listing**  
   - Admin `/auth/sessions` (list, terminate).  
   - Benutzer `/auth/sessions/self` (eigene Sessions).  
6. **Registration & Invite Acceptance**  
   - Admin erzeugt Invite (`POST /auth/invite`) → Token + optional Mail.  
   - Register Page lädt Details via `GET /auth/invite/:token` (Rolle, Ablauf, Hinweis); Seite lädt nur, wenn zusätzlich ein gültiger Admin-Passkey aus dem Provisioning-Tool vorliegt (`navigator.credentials.get`).  
   - Submit `POST /auth/register` mit `{invite_token, admin_passkey_assertion, username, email, password, optional passkey_publicKey}`.  
   - Nach Erfolg: optional Auto-Login (`POST /auth/login`), sonst Redirect + Success Toast. Invite wird sofort invalidiert.  
   - Ohne Invite: Button disabled bzw. Hinweis „nur per Einladung“.  

### 1.4 Komponentenverhalten
- **AuthCard**: Nebula Styling, responsive (max 480px). Slots für Forms/Alerts.  
- **CredentialForm**: Input Validation (Debounce 300ms), Show/Hide Password, Strength Meter.  
- **PasskeyButton**: Feature flag (browser support). Fallback bei Fehlern → Toast + Logging (category `auth/passkey`).  
- **TwoFactorSetup**:  
  1. `GET /auth/2fa/qrcode` → base64 PNG.  
  2. User scannt, gibt Code ein (`POST /auth/2fa/enable`).  
  3. Backup Codes (`GET /auth/2fa/backup`).  
- **SessionList**: Table + Filter (Device, IP, Location). Terminate Buttons (admin or self).  

### 1.5 Validierungen & States
- Login Button disabled, bis Formular valide; Loader während Request.  
- Fehler (z. B. falsches Passwort) → Toast + Field Feedback (rot).  
- Bruteforce Schutz: Bei `429` zeigt UI Countdown + Link zu Support.  
- Forgot/Reset:  
  - `POST /auth/forgot` → success message (keine Details).  
  - Reset Page validiert token vor Anzeige (Loading Spinner).  
- Register:  
  - Invite Token (JWT) wird lokal validiert (exp, signature).  
  - Admin-Passkey (aus Provisioning-Tool) erforderlich; ohne erfolgreiche `navigator.credentials.get` bleibt Seite gesperrt.  
  - Passwort-Policy (Length, entropy) direkt im Formular anzeigen.  
  - 2FA Setup wird direkt im Anschluss erzwungen (QR + Backup Codes).  
  - Bei ungültigem oder abgelaufenem Token → Fehlerseite + Link zum Admin.

### 1.6 Tests
- **Unit**:  
  - `useAuthForm` (validation, disabled states).  
  - `AuthContext` (token storage, refresh).  
- **Integration**:  
  - Playwright `auth.spec.ts`: login, logout, refresh fail, 2FA path.  
  - `sessions.spec.ts`: list/terminate.  
- **Security**:  
  - Ensure tokens never in LocalStorage.  
  - Verify `SameSite=strict` cookies (if used).  
  - CSP headers enforced bei Auth Pages (no inline scripts).  

### 1.7 Metriken & Logging
- `AuthContext` sendet Telemetry Events: `auth_login_success`, `auth_login_failure`, `auth_passkey_error`.  
- Session Terminate → Audit Log (`agent-orchestrator/events.log`) mit User/Device/IP.  
- Rate Limit Hits → Alert Drawer (severity warning) via `security:sessions` Topic.

### 1.8 Roadmap-Anknüpfung
- Security Control Center (S2) baut auf diesen Komponenten auf.  
- Settings → Security Tab (O1) nutzt dieselben APIs (`/settings/security`, `/auth/audit`).  
- CVE-Verweise: Offene Auth-relevante CVEs aus `CVE_CHECKLIST.md` steuern Priorität.

### 1.9 Access Control & Hardening
- **Registrierung nur durch Admin:** Invite-API ist admin-only; Register UI ist standardmäßig deaktiviert und wird nur über Invite-Link + Admin-Passkey freigeschaltet.  
- **Domain Lock:** Selbst bei späterem Public Domain Setup wird `/register` hinter Reverse-Proxy-ACL (BasicAuth/ClientCert) verborgen; Environment `REGISTRATION_ALLOWED_HOSTS` definiert Whitelist.  
- **Passkey Requirement:** Ein dediziertes Provisioning-Tool erstellt den Admin-Passkey (Hardware Token). Ohne erfolgreiche Assertion verweigert die UI sämtliche Register-Actions.  
- **2FA Pflicht:** Nach erfolgreicher Registrierung muss 2FA sofort aktiviert werden (ohne Skip).  
- **Audit Trail:** Jeder Invite/Registration Flow schreibt Events in `agent-orchestrator/events.log` + `CVE_CHECKLIST.md` Referenzen, damit Missbrauch sofort auffällt.
- **Provisioning Tool (in Planung):**  
  - CLI/GUI Utility generiert Admin-Passkeys (WebAuthn) und verwaltet 2FA-Secrets.  
  - Exportiert QR + Backup Codes, signiert Invite Links und synchronisiert mit `systemsetup-allowlist`.  
  - Wird als separater Agent (z.B. vom `DocumentationAgent` oder `SystemSetupAgent`) entwickelt und hier dokumentiert, sobald implementiert.

## 1.10 Referenzen & Kontext
- `NAS_AI_SYSTEM.md` – Gesamtarchitektur, Gates und Incident-Flows.
- `CVE_CHECKLIST.md` – Relevante Auth-/Security-Findings, die vor Änderungen geprüft werden müssen.
- `docs/planning/MASTER_ROADMAP.md` – Verantwortlichkeiten/Owners (APIAgent, WebUIAgent) für Auth-Features.
- `Blueprint_WebUI.md` – Übergeordnete UI-Layer, Alert-Center, Modul-Landschaft.
- `docs/security/SECURITY_HANDBOOK.md` & `docs/policies/systemsetup-allowlist.md` – Governance rund um Logs, Pakete und Secrets.