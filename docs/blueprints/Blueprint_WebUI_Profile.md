# NAS.AI – WebUI Profile Blueprint

## 1. Scope & Responsibilities
- User-spezifische Einstellungen, Sicherheitsoptionen und Aktivitätsnachweise.
- Konsolidiert frühere „User Settings“ (Theme, Sprache, Notifications) direkt im Profil.
- Schnittstelle für 2FA-Status, API-Tokens, Sessions und Geräteverwaltung.

## 2. UX & Layout
```
┌──────────────────────────────────────────────┐
│ Header: Profile                              │
├────────────┬─────────────────────────────────┤
│ Sidebar    │ Main Panels                     │
│ [← Back]   │ ┌───────────────┬─────────────┐ │
│ [Home]     │ │ Info Card     │ Security    │ │
│ Sections   │ │ Activity Log  │ Sessions    │ │
│  • Profile │ │ API Tokens    │ Preferences │ │
│  • Security│ └───────────────┴─────────────┘ │
├────────────┴─────────────────────────────────┤
│ Footer: last login • 2FA status • logout all │
└──────────────────────────────────────────────┘
```
## 3. Data Flow (Profile + User Settings)
1. **Profile Info & Preferences**  
   - `GET /profile` → {username, display_name, email, avatar_url, last_login, devices, compliance, preferences}.  
   - `PUT /profile` (patch) aktualisiert Felder; Avatar Upload via `POST /profile/avatar`.  
   - `PUT /profile/preferences` (oder `/settings/user`) für Theme, Sprache, Startansicht, Notification Flags.  
2. **Password / Security**  
   - `PUT /profile/password` (current, new, confirmation).  
   - `/profile/2fa/*` Endpunkte aus Auth-Blueprint (enable/disable, backup codes).  
3. **Sessions & Activity**  
   - `GET /profile/sessions` (laufende Geräte) + `DELETE /profile/sessions/:id`.  
   - `GET /profile/activity?limit=50` für Audit-Events.  
4. **API Tokens**  
   - `GET /profile/tokens`, `POST` (name, permissions, expiry), `DELETE /profile/tokens/:id`.  
5. **Notifications (User Scope)**  
   - `GET/PUT /settings/notifications?scope=user` – persönliche Events/Channels.  
6. **Telemetry Hooks**  
   - Frontend sendet `profile_view`, `profile_update_success`, `profile_token_created`, `profile_logout_all`, `profile_preference_saved`.

## 4. Components
- **ProfileCard** – Avatar Upload, display name, email, last login.
- **SecurityPanel** – Password change form, 2FA status, passkey enrollment link.
- **SessionTable** – Device/IP, location (optional), terminate button.
- **ActivityList** – Timeline mit Icons (login, file upload, settings change).
- **TokenManager** – Table + modal für Token-Erstellung (permissions dropdown).
- **PreferenceForm** – Theme selector, language, start page, view mode, notification digest.
- **DevicesPanel** – Liste vertrauenswürdiger Geräte (Browser Fingerprint, OS, IP, Trust-Status, „Mark as trusted/untrust“).
- **NotificationsPanel** – E-Mail/SMS/Webhook-Einstellungen für Login-Alerts, Token-Events, Profile-Änderungen inkl. Channel-Verifizierung.
- **CompliancePanel** – Anzeigen von Policies/Trainings (z. B. „Security Training 2025 ✓“), NDA Status, last acknowledgement.
- **RecoveryInfo** – Backup E-Mail, Recovery Codes vorhanden?, Buttons „Update Recovery Info“.

## 5. Validation
- Avatar ≤ 2 MB, Formate PNG/JPG/WebP.
- Password strength meter (entropy ≥ 60 bits), kein reuse der letzten N Passwörter.
- Token permissions (enum: read/write/admin) nur für Rollen ≥ admin editierbar.
- Session Terminate erfordert Confirm Dialog; Admin kann „Logout All“ erzwingen.
- Language/theme Werte müssen in Allowlist liegen.
- Trusted Devices erfordern Fingerprint-Hash + Ablaufdatum (z. B. 30 Tage).
- Notifications nur, wenn Kontaktmethoden verifiziert sind (Verifizierungsstatus anzeigen).
- Preferences-Änderungen haben Autosave; jede Mutation wird trotzdem im **Audit-Log des Orchestrators** festgehalten.

## 6. Tests/Telemetry
- **Unit Tests:** `useProfileStore`, TokenManager logic, Activity filter utils.
- **Integration/Playwright:** avatar upload + save, password change (happy/error), token CRUD, session terminate, preference update, notification toggle.
- **Telemetry Events:** `profile_update_success`, `profile_token_created`, `profile_logout_all`, `profile_preference_saved`, `profile_notification_toggle`.
- **Audit Logging:** Aktionen schreiben in `agent-orchestrator/events.log` (payload: user_id, action, metadata).
- Notifications/Device-Updates lösen zusätzliche Events `profile_notifications_updated`, `profile_device_trust_changed` aus.

## 7. Security & Hardening
- Sensitive Aktionen (password/token/session) nur mit re-auth (enter password/2FA) – optional prompt.
- API responses minimieren (kein Hash/secret exposure); tokens nur einmalig anzeigen.
- Rate Limit auf token creation + session terminate; Alerts bei Missbrauch (via `security:sessions`).
- All profile mutations folgen Statuslog-Pflicht (`status/<agent>/...` wenn Agent involviert).
- Gerätevertrauen (trusted device) speichert nur Fingerprint-Hash, niemals Rohdaten.
- Compliance-Panel zieht Daten aus signierten Quellen (`/compliance/status`), read-only für normale User.

## 8. References
- `NAS_AI_SYSTEM.md` (Status/Reporting, Gates).
- `Blueprint_WebUI_Auth.md` (2FA, Session Modelle).
- `CVE_CHECKLIST.md` (Auth/Profile Findings).
- `docs/planning/MASTER_ROADMAP.md` (Owner: WebUIAgent + APIAgent).
- `Blueprint_WebUI_Settings_Expert.md` – Admin-only Einstellungen (System/Security global).