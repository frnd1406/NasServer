# Repository-Struktur – Vorschlag (Codex)

> **⚠️ PFLICHTLEKTÜRE BEACHTEN:** Dieses Dokument darf nur verwendet werden, nachdem folgende Unterlagen gemäß [Agent Pre-Task Checklist](docs/AGENT-CHECKLIST.md) durchgearbeitet wurden:
> 1. `NAS_AI_SYSTEM.md` – System-Architektur & Governance
> 2. `docs/security/SECURITY_HANDBOOK.md` – Audit & Evidenz
> 3. `docs/planning/AGENT_MATRIX.md` – Agenten-Arbeitsregeln
>
> **Arbeitsablauf:** Lesen → Analysieren → Planen → Umsetzung (STRIKT!)

> Stand: 21. November 2025 – Dieser Vorschlag gruppiert vorhandene Dateien in logischere Bereiche, ohne bestehende Backups anzutasten.

## 1. Zielbild (Übersicht)
```
/
├── docs/
│   ├── adr/                 # Architecture Decision Records
│   ├── blueprints/          # WebUI Design-Dokumente
│   ├── planning/            # Projekt-Planung & Roadmaps
│   ├── policies/            # Richtlinien und Policies
│   ├── manuals/             # Benutzerhandbücher
│   └── vision/              # Langfristige Vision und Konzepte
├── backend/                 # Go-Quellcode + configs
├── frontend/                # React-Quellcode
├── orchestrator/            # Go Orchestrator
├── agents/                  # Agent Scripts
│   ├── systemsetup/
│   ├── documentation/
│   └── ...
├── infrastructure/          # Docker, Configs
│   ├── docker-compose.yml
│   └── caddy/
├── CVE_CHECKLIST.md
├── NAS_AI_SYSTEM.md
├── README.md
└── status/                  # Agent Status-Updates
    ├── Orchestrator/
    ├── AnalysisAgent/
    ├── APIAgent/
    ├── SystemSetupAgent/
    ├── NetworkSecurityAgent/
    ├── WebUIAgent/
    ├── DocumentationAgent/
    ├── PentesterAgent/
    └── archive/
```

## 2. Details nach Bereich

### docs/
- **adr/**: Architekturentscheidungen mit Begründungen (z.B. `004-repository-structure.md`).
- **blueprints/**: alle `Blueprint_WebUI_*.md` (Auth, Storage, Settings etc.).
- **planning/**: `MASTER_ROADMAP.md`, `AGENT_MATRIX.md` (Strategische Planung & Roadmaps).
- **policies/**: Compliance & Allowlist-Dokumente (`orchestrator-collaboration.md`, `systemsetup-allowlist.md`, `SECURITY_HANDBOOK.md`).
- **manuals/**: `USER_LOGIN_GUIDE.md` (Benutzerhandbücher).
- **vision/**: `FUTURE_CONCEPTS.md` (Langfristige Vision und Konzepte).

### infrastructure/
- Fokussiert auf Anwendungs-Code (Go-API etc.). Observability-Dateien sind jetzt im `docker-compose.yml` integriert.

### status/
- Bisherige Struktur bleibt bestehen (Agent/Phase).

## 3. ToDo bei Umsetzung (bereits erledigt)
1. Ordner `docs/planning`, `docs/manuals`, `docs/vision` anlegen. (ERLEDIGT)
2. Dateien in die passenden Unterordner verschieben (unter Beachtung von Referenzen in Markdown – ggf. relative Pfade aktualisieren). (ERLEDIGT)
3. `docs/README.md` aktualisieren. (ERLEDIGT)
4. Systempfade, die von Diensten genutzt werden (z. B. `/var/lib/observability/prometheus/config/prometheus.yml`), bleiben unverändert; hier nur Kopien bzw. Dokumentationsquellen im Repo verwalten. (BEIBEHALTEN)

## 4. Hinweise
- **Backups** unter `/mnt/raid/backups` bleiben unangetastet.
- Vor dem Verschieben prüfen, ob Automationsskripte/Referenzen absolute Pfade nutzen; ggf. anpassen.
- Dieser Plan ist ein Vorschlag – Anpassungen je nach Team-Präferenz.