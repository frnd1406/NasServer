# NAS.AI – WebUI Settings (Überblick & Expertenmodus)

Dieses Dokument beschreibt ausschließlich die **Expert Settings** für Administratoren. Alle nutzerbezogenen Einstellungen (Theme, Sprache, Benachrichtigungen, Passwort etc.) wurden in den Profil-Blueprint integriert (`Blueprint_WebUI_Profile.md`).

## 1. Scope & Responsibilities
- Verwaltung der systemweiten Policies (Security, Encryption, Notifications Defaults, System-Konfiguration).
- Zugriff nur für Admins; Wechsel in den Expert-Mode löst Warnung + Reauth aus.

## 2. UX & Layout (Expert Mode)
```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Header: Settings (Expert)                        [Mode: Expert ▼] [Save All]│
│ Breadcrumb: Home / Settings   • Dirty Indicator   • EXPERT WARNING ⚠️       │
├──────────────┬──────────────────────────────────────────────────────────────┤
│ Sidebar      │ Main Panels                                                   │
│ [← Back]     │ ┌──────────────────────────────────────────────────────────┐ │
│ [Home]       │ │ Tabs: Security | Notifications | Docs | System |          │ │
│ Tabs & Info  │ │       Encryption | Audit                                   │ │
│  • Linked CVEs│ │ ┌───────────────┬──────────────┬───────────────┐        │ │
│  • Pending    │ │ │ Basic Panels  │ Advanced Sec │ Live Preview  │        │ │
│  • Mode Desc  │ │ │ (read-only)   │ (TLS, Ports,  │ (diff view)   │        │ │
│ Warning Box   │ │ └───────────────┴──────────────┴───────────────┘        │ │
│  „Nur Experten│ │ Banner: „Expert Mode aktiviert – Änderungen wirken       │ │
│   …“          │ │ sofort. Reauth erforderlich.“                             │ │
├──────────────┴──────────────────────────────────────────────────────────────┤
│ Footer: Unsaved changes • Last saved timestamp • Revert • Save All • Logs ▼ │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 3. Data Flow (Expert Endpoints)
- `GET/PUT /settings/security` – Session-/Password-Policies, Login Limits, 2FA Defaults.
- `GET/PUT /settings/notifications` – Globale Event→Channel Defaults.
- `GET/PUT /settings/system` – Hostname, Ports, Logging, Maintenance Mode.
- `GET/PUT /settings/encryption` – KMS Alias, Key Rotation, Cipher Suites.
- `GET/PUT /docs-settings` – (falls DocsTerminal reaktiviert wird) Lock/Password Policies.
- `GET /settings/audit` – Änderungsprotokolle.
- Mode Toggle: `?mode=expert`; erfordert Reauth + Warnbanner.

## 4. Components (Expert)
- **ModeToggle** (Admin-only) + **ExpertWarningModal**.
- **SettingTabs** für Security, Notifications, Docs, System, Encryption, Audit.
- **AdvancedPanel** pro Tab (TLS, Ports, Rotation etc.).
- **EncryptionCard** (KMS, Rotation, Offline Rules).
- **NotificationMatrix (Org Scope)** – Events vs Channels default.
- **AuditSummary** + Download.
- **SaveBar** + Logs Dropdown.

## 5. Validation & Guardrails
- Admin Role + Reauth Pflicht; alle Aktionen loggen `settings_change`.
- Kritische Felder (Ports, TLS, KMS) erfordern Double-Confirm und verweisen auf `systemsetup-allowlist`.
- Dependencies: Files EncryptDialog disabled wenn `/settings/encryption` unvollständig ist.

## 6. Tests & Telemetrie
- **Unit:** Validators, `useSettingsExpertStore`, NotificationMatrix Org.
- **Integration:** Playwright (Admin toggles Expert Mode, ändert TLS), diff preview & save.
- **Events:** `settings_saved`, `settings_mode_switch`, `settings_advanced_change`, `settings_port_change_warning`.
- Audit Log: `/var/log/settings-changes.log` (append-only) + `settings_audit_view` Telemetrie.

## 7. References
- `Blueprint_WebUI_Profile.md` – User-Level Settings (Theme, Notifications personal etc.).
- `NAS_AI_SYSTEM.md` – Gates & Reporting.
- `Blueprint_WebUI_Auth.md` – Security dependencies.
- `CVE_CHECKLIST.md` – relevante Findings.
- `docs/planning/MASTER_ROADMAP.md` – Owner (WebUIAgent + APIAgent).
- `docs/policies/systemsetup-allowlist.md`, `docs/security/SECURITY_HANDBOOK.md` – Governance.