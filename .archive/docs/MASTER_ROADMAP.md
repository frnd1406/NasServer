# Master-Roadmap für NAS.AI

**Version:** 1.0
**Datum:** 21.11.2025
**Status:** AKTIV

## Einleitung
Diese Master-Roadmap beschreibt die komplette Neuimplementierung des NAS.AI Multi-Agenten-Systems, mit starkem Fokus auf Security-by-Design, testbarer Architektur und automatisierbarer Auslieferung. Sie ersetzt alle vorherigen Roadmap-Dokumente, einschließlich `agent-roadmap.md`.

## Phasen & Meilensteine

### Aktueller Sync (Q4 2025)
- **AnalysisAgent:** Phase 3 (Alerting) abgeschlossen – Alerts werden erzeugt und müssen in UI/UX eingelastet werden; Phase 4 darf erst nach UI-Gate weiter.
- **PentesterAgent:** Phase 4 (Red Team) bereits gestartet – muss mit Analysis/Observability rückgekoppelt werden.
- **Orchestrator:** Phase-Übergänge sind asynchron; Phase 3 wird für AnalysisAgent abgeschlossen, während PentesterAgent in Phase 4 läuft. Roadmap und Tickets sind entsprechend angepasst.

### Phase 0 – Discovery, Guardrails & Telemetrie (KW 1)
- **Inventar des bestehenden Codes:** Aufnahme des Codes, insbesondere Monolithen wie `infrastructure/api/src/handlers/files.go:1` und Shell-Orchestrator `infrastructure/orchestrator/bin/orchestrator-event-loop.sh:4`.
- **Definition der Zielarchitektur:** Klärung von Domain-Layer, Event-Bus und Deployment-Topologie, inklusive Architectural Decision Records (ADRs).
- **Security-/Compliance-Gates:** Festlegung von Secrets-Policy, Code-Scanning und Sign-Off-Prozessen als Blocker für nachfolgende Phasen.
- **Minimal-Monitoring:** Sofortiges Deployment von `node_exporter` + Prometheus zur Erfassung von CPU/RAM/IO-Baselines.

### Phase 1 – Plattform- und API-Fundament (KW 2–4)
- **Neues Backend:** Aufsetzen eines neuen Go- oder Rust-Backends mit klaren Packages (Auth, Files, Shares, Backups, System) und jeweiligen Services & Repository-Interfaces.
- **File-Service Extraktion:** Extrahierung des File-Service aus `files.go`, Implementierung von Storage-Abstraktion und Pfad-Sandbox, Ergänzung von Unit-Tests und Benchmarks.
- **Orchestrator-Rewrite Start:** Beginn des Event-Dienstes in Go/Rust mit echter Queue/State-DB, schrittweiser Ersatz des Shell-Workflows.
- **Konfigurations- und Secrets-Layer:** Einführung eines 12-Factor-basierten Layers (ENV + Secret Manager), Entfernung aller Hardcodierungen.

### Phase 2 – Security, Data & Orchestrator-GA (KW 5–6)
- **Umfassende Security:** Neuaufbau von CSRF/WebSocket/Auth (Token Exchange, Rotating Refresh, WebSocket-Handshake mit JWT, Origin-Checks).
- **RBAC/ABAC Implementierung:** Implementierung, Test und Dokumentation von Role-Based / Attribute-Based Access Control.
- **Backup- und Restore-Jobs:** Bereitstellung von Snapshots + Offsite-Backups mit Monitoring und Drill-Runbooks, inklusive Regularien für DB-Dumps, Docker-Volumes und Config-Backups.
- **Infrastruktur-Scans:** Verankerung von SAST/DAST und IaC-Scans in CI/CD, automatisiertes Tracking von Findings.
- **Orchestrator GA:** Produktivsetzung des neuen Event-Dienstes, Abschaltung des Shell-Orchestrators.
- **Observability-Baseline:** Produktive Nutzung von strukturierten Logs, Prometheus/Grafana und OTel-Traces; Alert-Routing und SLOs als Gatekeeper für Phase 3.
- **Test-Gates:** API-Contract-, Auth-, Files- und Backup-Tests laufen in CI und blockieren Deployments bei Regressionen; Coverage-Report ≥80 % für neue Services.

### Phase 3 – Experience Layer & Client Rewrite (KW 7–9)
- **Frontend-Umstellung:** Migration des Frontends auf neues Designsystem & State-Layer (React + Zustand/TanStack Query), Virtualisierung des Filebrowsers, Kapselung des API-Clients.
- **Mobile/PWA-Client:** Aufsetzen des Mobile/PWA-Clients auf denselben API-Contract.
- **Echtzeit-/Toast-/Offline-Pfade:** Vereinheitlichung über authentifizierte WebSocket-Events.
- **End-to-End-Messbarkeit:** Verankerung von Observability- und Test-Hooks des Backends über Shared SDKs auch im Client (Tracing-Header, QA-Hooks).

### Phase 4 – Observability, Testing & Cut-over (KW 10–11)
- **Umfassende Tests:** Vollständige E2E-, UI- und Chaos-Test-Suiten (inkl. Mobile/WebUI) laufen parallel zum neuen Observability-Stack und simulieren Disaster-Szenarien.
- **Load-/Chaos-Tests:** Validierung von Failover- und Rollback-Fähigkeiten (inkl. Orchestrator + Automation-Agenten).
- **Datenmigration & Parallelbetrieb:** Planung der Datenmigration, Schreiben von Backfill-Skripten, Etablierung von Blue/Green- oder Canary-Deployment mit automatischem Rollback.
- **Cut-over:** Hypercare (1–2 Wochen), Abschaltung des Altsystems.

## Rollen, Aufgaben & Prioritäten

### Systemarchitekt / Lead-Entwickler
- **Aufgaben:** Architektur-Blueprints & ADRs erstellen, Schnittstellendefinitionen für Event-Bus & Service-Verträge, Technologie-Stack-Freigabe, Tech-Debt-Backlog und Migrationsstrategie.
- **Fokus:** Altschulden auflösen, Security-Gates verankern, klare Grenzen zwischen Agenten schaffen.

### UI/UX-Designer
- **Aufgaben:** Informationsarchitektur & User-Flows neu denken, High-Fidelity-Designsystem ("Nebula 2.0"), Design-Tokens + Motion-Guidelines.
- **Fokus:** Usability-Probleme vermeiden, Loading/Offline-Zustände visualisieren, Barrierefreiheit sicherstellen.

### Backend-Entwickler
- **Aufgaben:** Domain-Services implementieren, eventischer Orchestrator und Scheduler, Config-/Secrets-Layer, Performance-/Load-Tests.
- **Fokus:** Saubere Layer-Trennung, Idempotente File-Operationen, Observability-Hooks, sichere Pfad- und Quota-Prüfungen.

### Frontend-Entwickler
- **Aufgaben:** API-Client als Typed SDK neu schreiben, Filebrowser + Management-Views reimplementieren, Echtzeit- und Offline-Szenarien stabilisieren, Komponentenbibliothek + Storybook aufsetzen.
- **Fokus:** Kein direkter Storage-Zugriff im UI, Zustände teilen statt duplizieren, Progressive Enhancement und Performance.

### Security-Entwickler / DevSecOps
- **Aufgaben:** Secret-Management, JWT-Rotation, Policy Enforcement, Automatisierte Scans, Hardening-Guides, Runbooks.
- **Fokus:** Security-by-Default, Zero-Trust-CORS, sichere Default-Konfiguration, Nachweisbare Remediation der CVEs.

### QA-Manager / Tester
- **Aufgaben:** Teststrategie & Abdeckung definieren, Automatisierte Regressionstests, Load/Stress-Tests, Release-Checklisten.
- **Fokus:** Alle historischen Bugs/CVEs mit Tests abdecken, Smoke-Tests in CI erzwingen, Nichtfunktionale Anforderungen messen.

### Observability / SRE Agent
- **Aufgaben:** Logging-, Metrics- und Tracing-Standards definieren, SLOs/SLA, Alert-Routing und Runbooks erstellen, Chaos-/Failure-Tests planen.
- **Fokus:** End-to-End-Transparenz, Kapazitätsplanung, Self-Healing-Strategien.

### Documentation & Knowledge Agent
- **Aufgaben:** Living-Docs für Architektur, APIs, Security-Policies und Runbooks pflegen, Developer Experience (Onboarding-Guides, CLI-Handbücher) bereitstellen, Change-Logs, Migration-Guides und User-Facing-Dokumente publizieren.
- **Fokus:** Konsistenter Wissensstand, Nachvollziehbarkeit von Entscheidungen, Dokumentation der Backup-/Restoreprozesse.

### AIKnowledgeAgent
- **Aufgaben:** Semantic/Visual Search Pipelines, Voice-/Assistant-APIs, RDR-Analyse.
- **Fokus:** Semantic API GA, Visual Lens Beta, Voice-Assistent Pilot mit Duress Detection.

### TestAutomationAgent
- **Aufgaben:** E2E-/Regression-/Load-Tests für Auth, Files, Backup-Restore, WebSocket; Chaos/Stress-Szenarien automatisieren; Testdaten & Fixtures verwalten.
- **Fokus:** CI-Gates für Kernflows aktiv, Load-Test-Reports verfügbar, Chaos-Tests in Observability Alerts integriert.

## Priorisierte Nächste Schritte
1.  **Architektur- & Security-Gate (Phase 0):** Zielbild, Threat-Model und Freigabe-Kriterien dokumentieren.
2.  **Core-Service-Refactor (Phase 1):** Neues File-/Auth-Service implementieren, sichere Config-Layer & Event-Orchestrator; CI-Pipeline mit Tests & Scans einrichten.
3.  **Security- & Data-Hardening (Phase 2):** CSRF/WebSocket/Auth finalisieren, Backups & Restore verifizieren, Security-Gates automatisieren.
4.  **Client Rewrite & Observability (Phase 3/4):** UI/Mobile neu liefern, Observability-Stack produktiv nehmen, E2E-Tests + Cut-over planen.
5.  **Automation & AI Enablement (Phase 2+):** PolicyAutomationAgent etablieren, DocumentationAgent aufsetzen und AIKnowledgeAgent als eigenständige Deliverable planen; TestAutomationAgent liefert verbindliche Gates (E2E/Load/Chaos).
