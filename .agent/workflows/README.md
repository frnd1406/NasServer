# Dokumentation

Übersicht über alle Projektdokumentationen.

---

> **⚠️ WICHTIG FÜR ALLE AGENTEN:**
> Vor Verwendung dieser Dokumentation MUSS die [Agent Pre-Task Checklist](AGENT-CHECKLIST.md) durchgearbeitet werden.
>
> **Pflichtlektüre vor jeder Aufgabe:**
> 1.  `NAS_AI_SYSTEM.md` - System-Architektur & Governance
> 2.  `docs/security/SECURITY_HANDBOOK.md` - Audit & Evidenz
> 3.  `docs/planning/AGENT_MATRIX.md` - Agenten-Arbeitsregeln
>
> **Arbeitsablauf:** Lesen → Analysieren → Planen → Umsetzung (STRIKT!)

---

## Verzeichnisstruktur

```
docs/
├── adr/                # Architecture Decision Records
├── blueprints/         # WebUI Design-Dokumente
├── planning/           # Strategische Planung & Roadmaps
├── policies/           # Richtlinien und Policies
└── manuals/            # Benutzerhandbücher
└── vision/             # Langfristige Vision und Konzepte
```

## Schnellzugriff

### 📐 Architecture Decision Records (ADR)
→ [adr/](./adr/)
- Architekturentscheidungen mit Begründungen
- Beispiel: `004-repository-structure.md`

### 🎨 WebUI Blueprints
→ [blueprints/](./blueprints/)
- Design-Spezifikationen für die WebUI
- **Hauptdokument:** [Blueprint_WebUI.md](./blueprints/Blueprint_WebUI.md)
- **Module:**
  - [Auth](./blueprints/Blueprint_WebUI_Auth.md) - Authentifizierung
  - [Files](./blueprints/Blueprint_WebUI_Files.md) - Dateiverwaltung
  - [Backup](./blueprints/Blueprint_WebUI_Backup.md) - Backup-System
  - [Storage](./blueprints/Blueprint_WebUI_Storage.md) - Speicherverwaltung
  - [Settings](./blueprints/Blueprint_WebUI_Settings.md) - Einstellungen
  - [Users](./blueprints/Blueprint_WebUI_Users.md) - Benutzerverwaltung
  - [Shares](./blueprints/Blueprint_WebUI_Shares.md) - Freigaben
  - [Profile](./blueprints/Blueprint_WebUI_Profile.md) - Benutzerprofil

### 🗺️ Planung & Roadmaps
→ [planning/](./planning/)
- **[MASTER_ROADMAP.md](./planning/MASTER_ROADMAP.md)** - Überblick über aktuelle Phasen und Meilensteine
- **[AGENT_MATRIX.md](./planning/AGENT_MATRIX.md)** - Agenten-Übersicht, Rollen und Verantwortlichkeiten

### 📋 Policies & Richtlinien
→ [policies/](./policies/)
- [orchestrator-collaboration.md](./policies/orchestrator-collaboration.md) - Agent↔Orchestrator Workflow
- [systemsetup-allowlist.md](./policies/systemsetup-allowlist.md) - SystemSetup Allowlist

### 💻 Development & Reference
→ [development/](./development/)
- [DEV_GUIDE.md](./development/DEV_GUIDE.md) - Entwicklungs-Guide
- [REFERENCE_SNIPPETS.md](./reference/REFERENCE_SNIPPETS.md) - Code-Snippets & Patterns

### 📚 Handbücher
→ [manuals/](./manuals/)
- [USER_LOGIN_GUIDE.md](./manuals/USER_LOGIN_GUIDE.md) - Anleitung zur Registrierung und zum Login
- [PASSWORD-MANAGER.md](./manuals/PASSWORD-MANAGER.md) - Passwort-Manager Dokumentation
- [EMAIL_TESTING.md](./manuals/EMAIL_TESTING.md) - E-Mail Testing Guide

### 🔮 Vision & Konzepte
→ [vision/](./vision/)
- [FUTURE_CONCEPTS.md](./vision/FUTURE_CONCEPTS.md) - Langfristige Ideen und Zukunftsperspektiven

## Weitere Dokumentation

- **Status Reports:** `/status/` - Agent-Status und Phase-Reports (siehe auch `docs/planning/AGENT_MATRIX.md`)
- **Lessons Learned:** [LESSONS-LEARNED.md](./status/LESSONS-LEARNED.md)
- **System-Übersicht:** `/NAS_AI_SYSTEM.md`
- **CVE Checklist:** `/CVE_CHECKLIST.md`

## Navigation

```bash
# Alle Blueprints anzeigen
ls docs/blueprints/

# Alle ADRs anzeigen
ls docs/adr/

# Nach Stichwort suchen
grep -r "keyword" docs/
```

---

**Letzte Aktualisierung:** 2025-11-21
**Struktur:** Gemäß `STRUCTURE_PROPOSAL.md` und Konsolidierungsplan