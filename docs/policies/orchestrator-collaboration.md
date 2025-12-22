# Richtlinie: Zusammenarbeit Agenten ↔ Orchestrator

**Version:** 1.1  
**Stand:** 21. November 2025  
**Verbindlich für:** alle aktiven technischen Agenten und den Orchestrator  
**Quellen:** `NAS_AI_SYSTEM.md` §4, `docs/AGENT-CHECKLIST.md`, `docs/planning/AGENT_MATRIX.md`, `docs/security/SECURITY_HANDBOOK.md`

---

## 1. Ziel & Geltung

Diese Richtlinie präzisiert, wie Agenten und der Orchestrator Aufgaben abstimmen, Statuslogs pflegen und Eskalationen handeln. Sie ergänzt die [Agent Pre-Task Checklist](docs/AGENT-CHECKLIST.md) sowie das [Sicherheits-Handbuch](docs/security/SECURITY_HANDBOOK.md) und gilt für **jede** Interaktion vom Ticket-Start bis zum Abschluss. Abweichungen benötigen eine schriftliche Owner-Freigabe.

---

## 2. Rollen & Verantwortlichkeiten

| Rolle | Kernaufgaben | Verpflichtende Artefakte |
|-------|--------------|---------------------------|
| **Orchestrator** | Tickets & Gates verwalten, Backlog priorisieren, Freigaben dokumentieren, Evidence prüfen, Owner informieren | `status/Orchestrator/...`, `status/backlog/*.md`, Incident-/Package-Queue unter `/var/lib/orchestrator/` |
| **Aktiver Agent** | Aufgabenverständnis bestätigen, TODO-Plan erstellen, Umsetzung/Evidence dokumentieren, Risiken melden | `status/<Agent>/<NNN>_*`, Evidence-Files, ggf. `status/<Agent>/Clarification/*` |
| **Support-Agenten** (SystemSetupAgent, NetworkSecurityAgent, PentesterAgent, AnalysisAgent, DocumentationAgent) | Allowlist-/Security-Freigaben, Audit-Logs, Telemetrie bereitstellen | `status/<Agent>/...`, `/var/log/package-*.log`, Observability Dashboards |

---

## 3. Arbeitszyklus (Handshake)

1. **Init (Ticket aktivieren)**
   - Orchestrator prüft `status/backlog/` und offene Incidents. Nur wenn alte Tasks abgeschlossen → neues Ticket auf „terminal_active“ setzen.
   - Agent bestätigt Pflichtlektüre (`NAS_AI_SYSTEM.md`, `docs/security/SECURITY_HANDBOOK.md`, `docs/planning/AGENT_MATRIX.md`, eigenes Log + Backlog) **im ersten Logabschnitt**.
   - Neues Statuslog erzeugen (`./scripts/new-statuslog.sh`). Abschnitt „TODO/Plan“ enthält klare Schritte und Gate-Checks.

2. **Analyse & Umsetzung**
   - Agent dokumentiert Ist-Zustand, Risiken, Ressourcen. Ohne Plan beginnt keine Umsetzung.
   - Jede relevante Aktion → Evidence (Logs, Hashes, Screenshots). Pfade im Statuslog referenzieren.
   - Package-/Allowlist-Anfragen laufen immer über `/var/lib/orchestrator/package-queue/` + Statuslog-Link.
   - Fortschritt ≥4h oder nach jedem Teilabschnitt im Log ergänzen („Checkpoint“). Orchestrator validiert und setzt Ticketstatus (`in_analysis`, `executing`, `waiting_owner`).

3. **Review & Abschluss**
   - Agent markiert Umsetzung als abgeschlossen, führt Tests, aktualisiert Dokumentation.
   - Abschlussabschnitt enthält: Ergebnisse, Tests, offene Punkte, „Terminal freigegeben“.
   - Orchestrator prüft Evidence. Bei OK → Ticket `reviewed` → Backlog nächstes Item. Bei Lücken → Status `needs_evidence`, Agent ergänzt.

---

## 4. Kommunikationskanäle & Artefakte

- `status/<Agent>/Phase|Policy|Evidence/` – Primäres Log, Evidence-Unterordner Pflicht bei Security/Infra-Themen.
- `status/backlog/<YYYYMMDD>_<agent>_<task>.md` – Offene oder pausierte Tasks. Wenn Verzeichnis fehlt → Orchestrator legt es vor Ticketstart an.
- `status/<Agent>/Clarification/*.md` – Fragen/Antworten bei Unklarheiten. Jeder „clarification_needed“ Eintrag enthält Referenzen & gewünschte Entscheidung.
- `status/Orchestrator/Phase/...` – Gate-/Phasenprotokolle. Orchestrator verlinkt auf Agentenlogs.
- Incident-/Package-Verzeichnisse unter `/var/lib/orchestrator/` – maschinenlesbare JSONs für Automationen (Allowlist, Privilege-Escalation, Alerts).
- Observability Hooks (`files:progress`, `system:alerts`) – **Orchestrator** liefert Telemetrie, Orchestrator mappt Topics auf Tickets.

---

## 5. Task Lifecycle States

| State | Bedeutung | Verantwortlich für Wechsel |
|-------|-----------|-----------------------------|
| `queued` | Ticket wartet auf freien Terminalslot / ungeklärtes Backlog | Orchestrator |
| `in_analysis` | Pflichtlektüre, Ist-Analyse, Plan erstellt | Agent bestätigt Logeintrag, Orchestrator prüft |
| `executing` | Umsetzung läuft | Agent aktualisiert Log-Einträge mindestens alle 4h |
| `waiting_owner` | Entscheidung/Sicherheitsfreigabe offen | Agent formuliert Entscheidungsblock, Orchestrator informiert Owner |
| `needs_evidence` | Evidence/Test fehlt | Orchestrator markiert, Agent ergänzt |
| `reviewed` | Orchestrator hat Log + Evidence akzeptiert | Orchestrator |
| `blocked` | Offenes Risiko/CVE/Backlog verhindert Fortsetzung | Orchestrator setzt Flag, Agent verschiebt Kontext nach `status/backlog/` |

State-Änderungen werden im Statuslog dokumentiert (z. B. „Ticket → executing @ 2025-11-17T10:05Z“) und vom Orchestrator im eigenen Phase-Log gespiegelt.

---

## 6. Eskalationen & Sonderfälle

- **Clarifications:** Agent erstellt `clarification-needed` Log, verlinkt betroffene Dateien, Orchestrator liefert Antwort im passenden Folder und aktualisiert Ticketstate.
- **Risiken/CVEs:** Sobald Sicherheitsrisiko erkannt wurde → `status/security/OPEN-RISKS.md` + Agentenlog + Hinweis an Orchestrator. Kein Weiterarbeiten, bis Gate geschlossen.
- **Privilege Elevation:** Follows Break-Glass Flow aus `NAS_AI_SYSTEM.md`. Jede Anforderung verlangt vorab ein Statuslog-Update + Ticket-Referenz.
- **Fehlende Backlog-Items:** Orchestrator legt Template-Datei an und verweist im Task-Log. Agents dokumentieren Lektüre (auch wenn Liste leer).
- **Automatisierte Alerts:** **Orchestrator** pusht Alerts -> Orchestrator → `status/backlog` oder direktes Ticket (`incident_<id>.json`).

---

## 7. Integration mit bestehenden Policies

- **Checklist Hooks:** Punkt „Arbeitsabfolge“ bleibt unverändert, erhält jedoch Verweis auf dieses Dokument.
- **Audit-Logging:** Terminal-Workflow (Ticket aktivieren → Arbeit → Evidence → Review) bleibt identisch, wird hier detailliert.
- **Roadmap-Gates:** Orchestrator aktualisiert `docs/planning/MASTER_ROADMAP.md` nach jedem Review und markiert Gate-Status (Security, Telemetry, Experience, AI).

---

## 8. Umsetzung & Review

- Änderungen an diesem Dokument werden vom Orchestrator vorgeschlagen und durch Owner `user` freigegeben.
- Jeder Agent bestätigt im nächsten Statuslog, dass Version 1.0 gelesen wurde (Rubrik „Pflichtlektüre“).
- Feedback fließt über `status/DocumentationAgent/...` in die nächste Revision.

---

**Kontakt:** Orchestrator-Team (`status/Orchestrator/Phase/...`) oder Owner `user`.

---

## Nachtrag: Definition of Done (Git Commit Policy)

Um den Reifegrad des Systems von der Planung zur Implementierung zu verschieben, gilt ab sofort folgende Direktive für den Orchestrator:

1. Code vor Konzept
Der Fortschritt wird ausschließlich anhand von Git Commits gemessen.
- Akzeptiert: Änderungen an .go, .sh, .sql, .yaml Dateien.
- Nicht akzeptiert: Reine .md Statusberichte ohne begleitenden Code.

2. Arbeitsnachweis
Eine Aufgabe gilt erst als abgeschlossen ("Done"), wenn funktionierender oder zumindest kompilierbarer Code in das Repository committet wurde.

3. Modus Operandi
Anstatt Pläne zu erstellen, wie ein Agent funktionieren könnte, wird der Agent-Code (MVP) implementiert. Fehlerhafter Code (WIP) ist besser als perfekte Dokumentation.
