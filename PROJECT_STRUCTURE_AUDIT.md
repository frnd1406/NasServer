# 📊 Analyse-Bericht: Project Cleanup & Restructuring

Dieser Bericht analysiert die aktuelle Projektstruktur und gibt Empfehlungen zur Bereinigung und Neuorganisation.

#### Legende
- 🟢 **Safe**: Kann ohne Bedenken gelöscht/verschoben werden (Build-Artefakte, Cache).
- 🟡 **Verify**: Sollte kurz geprüft werden (alte Backups, Test-Skripte).
- 🔴 **Critical**: **VORSICHT!** Sicherheitsrisiko oder destruktiv (Secrets, DB-Daten).

---

## 1. Junk & Cleanup Candidates (Lösch-Vorschläge)

Diese Dateien sind meist generiert, veraltet oder temporär und blähen das Repository unnötig auf.

| Datei / Pfad | Aktion | Level | Begründung |
| :--- | :--- | :--- | :--- |
| `infrastructure/webui/dist/` | 🗑️ DELETE | 🟢 | Build-Output (Frontend). Kann jederzeit neu generiert werden (`npm run build`). |
| `infrastructure/api/bin/` | 🗑️ DELETE | 🟢 | Kompilierte Go-Binaries. Sollten nicht im Repo liegen. |
| `infrastructure/analysis/analysis-agent` | 🗑️ DELETE | 🟢 | Kompiliertes Binary. |
| `infrastructure/monitoring/monitoring-agent`| 🗑️ DELETE | 🟢 | Kompiliertes Binary. |
| `infrastructure/ai_knowledge_agent/src/__pycache__` | 🗑️ DELETE | 🟢 | Python Bytecode Cache. |
| `infrastructure/db/backup_pre_vector_...sql` | 🗑️ DELETE | 🟡 | Altes manuelles Backup vom 29.11.2025. Wenn nicht mehr benötigt -> Weg. |
| `infrastructure/REINDEX_READY.md` | 🗑️ DELETE | 🟡 | Wahrscheinlich ein temporärer "Flag"-File oder Notiz. Inhalt prüfen. |
| `infrastructure/VECTOR_UPGRADE_COMPLETE.md`| 🗑️ DELETE | 🟡 | Status-Flag/Notiz nach Upgrade. Wahrscheinlich obsolet. |
| `infrastructure/api/DEPLOYMENT_SUCCESS.md` | 🗑️ DELETE | 🟢 | Temporäres Deployment-Log. |
| `infrastructure/api/PRODUCTION_LIVE.md` | 🗑️ DELETE | 🟢 | Temporäres Status-Log. |

---

## 2. Security & Critical Findings (Sofort handeln!)

Hier liegen Dateien, die **niemals** im Versionskontrollsystem (Git) liegen sollten.

| Datei / Pfad | Aktion | Level | Begründung |
| :--- | :--- | :--- | :--- |
| `infrastructure/secrets/` | 🛡️ MOVE/IGNORE | 🔴 | **CRITICAL!** Enthält Secrets (`jwt_secret`, `postgres_password`). **Empfehlung:** Ordner in `.gitignore` aufnehmen! Secrets sollten via Environment-Variablen oder Vault injectet werden. |

---

## 3. Restructuring & Organization (Ordnung schaffen)

Das Projekt hat viele Dokumentations-Dokumente im Root-Verzeichnis und in Unterordnern verstreut. Eine Konsolidierung in einem `docs/`-Ordner wird empfohlen.

### Vorzuschlagende Struktur für `docs/`
- `docs/planning/` (für Backlog, Master-Plan, Phasen)
- `docs/architecture/` (für System-Diagramme, Konzepte)
- `docs/api/` (für Endpunkte, Swagger, API-Status)
- `docs/incidents/` (für Post-Mortems, Fix-Berichte)
- `docs/security/` (für Security-Handbücher, Policies)

| Datei / Pfad | Ziel-Pfad (Vorschlag) | Level | Begründung |
| :--- | :--- | :--- | :--- |
| `Master-Plan.md` | `docs/planning/` | 🟢 | Projektplanung. |
| `Phase3b.md` | `docs/planning/` | 🟢 | Phasenplanung. |
| `BACKLOG.md` | `docs/planning/` | 🟢 | Aufgabenliste. |
| `NAS_AI_SYSTEM.md` | `docs/architecture/` | 🟢 | Architekturdokumentation. |
| `API_ENDPOINTS_COMPREHENSIVE.md` | `docs/api/` | 🟢 | API Referenz. |
| `infrastructure/GHOST_KNOWLEDGE_FIX.md` | `docs/incidents/` | 🟢 | Incident Report. |
| `infrastructure/CHAT_INTERFACE_FIX.md` | `docs/incidents/` | 🟢 | Incident Report. |
| `infrastructure/SECURITY_HANDBOOK.md` | `docs/security/` | 🟡 | Prüfen ob Duplikat zu `docs/security/SECURITY_HANDBOOK.pdf`. |
| `infrastructure/api/CSRF_ENDPOINT.md` | `docs/api/` | 🟢 | API Dokumentation. |
| `infrastructure/api/DOMAIN_CONFIG.md` | `docs/api/` | 🟢 | API Dokumentation. |
| `infrastructure/api/status/` | `docs/api/history/` | 🟡 | Ordner wirkt wie eine manuelle Status-Historie. |
| `infrastructure/webui/test_crypto.js` | `infrastructure/webui/test/` | 🟢 | Testskript lag lose im Source-Ordner. |
| `Gemini.pdf` | `docs/references/` | 🟡 | PDF im Root-Verzeichnis. |
