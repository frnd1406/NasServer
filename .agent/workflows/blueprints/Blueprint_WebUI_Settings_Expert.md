# NAS.AI – WebUI Settings Blueprint

## 1. Scope & Responsibilities
- Zentraler Ort für System-, User-, Security- und Notification-Einstellungen (`/settings`).
- Unterstützt zwei Modi: **User Mode** (vereinfachte Oberfläche) und **Expert Mode** (vollständige Parameter, Admin-only).

## 2. UX & Layout (Expert Mode)
```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Header: Settings (Expert)                        [Mode: Expert ▼] [Save All]│
│ Breadcrumb: Home / Settings   • Dirty Indicator   • EXPERT WARNING ⚠️       │
├──────────────┬──────────────────────────────────────────────────────────────┤
│ Sidebar      │ Main Panels                                                   │
│ [← Back]     │ ┌──────────────────────────────────────────────────────────┐ │
│ [Home]       │ │ Tabs: Security | Notifications | Docs | System |          │ │
│ Info Panel   │ │       Encryption | Audit                                   │ │
│  • Linked CVEs│ │ ┌───────────────┬──────────────┬───────────────┐        │ │
│  • Pending    │ │ │ Basic Panels  │ Advanced Sec │ Live Preview  │        │ │
│  • Mode Desc  │ │ │ (read-only)   │ (TLS, Ports,  │ (diff view)   │        │ │
│ Warning Box   │ │ └───────────────┴──────────────┴───────────────┘        │ │
│  "Nur Experten│ │ Banner: "Expert Mode aktiviert – Änderungen wirken       │ │
│   …"          │ │ sofort. Reauth erforderlich."                             │ │
├──────────────┴──────────────────────────────────────────────────────────────┤
│ Footer: Unsaved changes • Last saved timestamp • Revert • Save All • Logs ▼ │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 3. Data Flow
- `GET/PUT /settings/user` – Profilpräferenzen (Theme, Sprache, Notifications).
- `GET/PUT /settings/security` – Policies (session timeout, password rules, 2FA defaults, File-Encryption-Regeln).
- `GET/PUT /settings/system` – Hostname, ports, logging level (Expert Mode only).
- `GET/PUT /settings/notifications` – Channels, webhooks, thresholds.
- `GET/PUT /docs-settings` – DocsTerminal/lock policies (falls reaktiviert).
- `GET /settings/audit` – Änderungsprotokolle.
- Mode selection (`?mode=user|expert`) gesteuert via query param + local preference; Expert Mode erfordert Admin-Reauth + Warnbanner („Änderungen wirken sofort – nur erfahrenes Personal“). Im Expert Mode werden Basic Panels read-only, Advanced Panels aktiv.

## 4. Components
- **ModeToggle** – Dropdown/User vs Expert, Admin-only, löst Warnmodal & Reauth aus.
- **SettingTabs** – Render Tabs basierend auf Mode + Role.
- **FormSection** – Card mit Title, description, inputs (Input/Textarea/Toggle/Select).
- **AdvancedPanel** – Akkordeon für Expert-Optionen (z. B. custom TLS config).
- **NotificationMatrix** – Tabelle (Event vs Channel) mit toggles.
- **EncryptionCard** – Eigene Section mit Feldern für KMS, Key Rotation, Allowed Algorithms; zeigt Warnungen, wenn Files EncryptDialog aktiv ohne Settings.
- **SaveBar** – Sticky footer mit Dirty-State + revert/back buttons.
- **AuditSummary** – Zeigt letzte Änderungen + Link zu Status-Log.

## 5. Validation
- Role gating: Expert Mode nur für Admin; System Settings nur Admin.
- Schema validation (z. B. JSON schema pro tab) + inline error messages.
- Confirmation dialogs für kritische Optionen (Password policy change, system restart).
- Auto-save optional nur im User Mode; Expert Mode erfordert double-confirm.
- Each change loggt `settings_change` Event (user, tab, fields).

## 6. Tests/Telemetry
- **Unit:** `useSettingsStore`, validators, NotificationMatrix logic.
- **Integration:** Playwright scenario (user mode -> change theme), (Admin -> expert mode change TLS).
- **Events:** `settings_saved`, `settings_mode_switch`, `settings_advanced_change`.
- Telemetry: record time to complete update, errors per tab.
- Audit logs stored in `/var/log/settings-changes.log` (append-only).

## 7. References
- `NAS_AI_SYSTEM.md` (Gates, incidents), `Blueprint_WebUI_Auth.md` (security dependencies), `CVE_CHECKLIST.md` (policy-related CVEs), `docs/planning/MASTER_ROADMAP.md` (owners).
- `Blueprint_WebUI_Encryption.md` (falls KMS/Keys separat dokumentiert werden).
- `docs/security/SECURITY_HANDBOOK.md` (Audit-Log-Richtlinien und Änderungsverlauf).

## 8. Tab- & Field-Details

### 8.1 Account (User Mode)
- **Inputs:** Display Name (text), Email (read-only), Avatar uploader, Language select, Theme toggle.  
- **Actions:** `PUT /settings/user` (partial). Avatar via `POST /profile/avatar`.  
- **Validation:** Display name 3–50 Zeichen; Avatar ≤2 MB PNG/JPG/WebP.  
- **User vs Expert:** In Expert Mode ist dieser Tab read-only, Änderungen laufen über Profile.

### 8.2 Security (Basic + Expert)
- **User Mode:** Password change button (opens auth modal), 2FA status, passkey info.  
- **Expert Mode:** Policies (Session timeout min/max, login attempt limits, forced 2FA).  
- **API:** `GET/PUT /settings/security`.  
- **Warnings:** Changing session timeout <5 min zeigt Warndialog; Expert Mode muss CVE-Liste referenzieren (z. B. Auth Hardening).

### 8.3 Preferences
- **Fields:** Default landing page, view mode, time format, notification digest.  
- **Storage:** Persist via `/settings/user`; einige (wie local layout) optional im Browser.  
- **User vs Expert:** Beide Modi identisch, aber Expert Mode zeigt zusätzliche „apply to org“-Option (Admin only).

### 8.4 Notifications
- **Matrix:** Events (Login Failure, Backup Failure, Storage Alert) vs Channels (Email, Push, Webhook).  
- **Expert Mode:** Admin kann globale defaults definieren; User Mode nur personal toggles.  
- **API:** `/settings/notifications` + verification endpoints (e.g., `/notifications/email/verify`).  
- **Validation:** Webhook URLs per regex, require success ping.

### 8.5 Docs (optional)
- Wenn DocsTerminal deaktiviert → read-only. Falls reaktiviert: Lock toggle, password settings, invite codes.  
- API: `/docs-settings`.  
- Mode: Only Expert sees doc lock controls.

### 8.6 System (Expert)
- **Fields:** Hostname, domain, API port, WebUI port, logging level, maintenance mode toggle.  
- **Actions:** `PUT /settings/system`; show diff preview + require confirmation.  
- **Warnung:** Changing ports triggert modal mit „Service restart required“.  
- **Dependencies:** writes to `agents-config.yaml` via Orchestrator.

### 8.7 Encryption (Expert)
- **Fields:** KMS alias, Key rotation interval, Allowed cipher suites, Offline passphrase requirements.  
- **Integrationen:** Files EncryptDialog liest diese Werte; wenn nicht gesetzt, EncryptDialog disabled.  
- **API:** `/settings/encryption`.  
- **Validation:** rotation interval ≥ 7 Tage; alias muss existieren (KMS health check).  
- **Security:** require double confirmation + reference to `systemsetup-allowlist`.

### 8.8 Audit (Expert)
- **View:** Table of recent setting changes (user, tab, fields, timestamp) fetched from `/settings/audit`.  
- **Actions:** Download log, filter by user/tab.  
- **Mode:** Expert only; read-only.  
- **Telemetry:** `settings_audit_view` event when opened.

### 8.9 Mode Toggle Behavior
- Switching von Profil/Standard-Ansicht → Expert Settings blendet Warnbanner ein („Änderungen wirken sofort“), verlangt Passwort/2FA/Passkey.  
- Log Event `settings_mode_switch` (mit old/new Mode, user id); Hinweis, dass persönliche Settings im Profil verbleiben.

## 9. History & Rollback
- Jede Expert-Änderung wird in `/var/log/settings-changes.log` protokolliert (siehe `docs/security/SECURITY_HANDBOOK.md`). Optionaler Mirror kann deaktiviert werden (`mirror_enabled: false`), aber Primary Log bleibt Pflicht.
- UI enthält „Recent Changes“-Panel mit Diff-Link und Button „Rollback“. Der Button erzeugt automatisch Status-Datei + ruft `orchestratorctl rollback-settings <id>` auf.
- Rollback ist verpflichtend verfügbar; wenn Logging deaktiviert oder fehlerhaft ist, blockiert die UI weitere Saves, bis der Zustand behoben wurde.
- Warnhinweis wird erneut angezeigt, bevor ein diff angewendet oder zurückgerollt wird, damit Nutzende die Auswirkungen verstehen.