# Agent Pre-Task Checklist

**Version:** 1.1
**Datum:** 2025-11-21
**Quelle:** `NAS_AI_SYSTEM.md`, `docs/security/SECURITY_HANDBOOK.pdf`, `docs/planning/AGENT_MATRIX.md`

---

## ZWECK

Diese Checkliste MUSS von jedem Agenten vor Beginn JEDER neuen Aufgabe durchgearbeitet werden. Sie stellt sicher, dass alle Governance-Anforderungen eingehalten werden.

---

## ✅ PFLICHTLEKTÜRE (VOR JEDER AUFGABE)

### 1. Kern-Dokumente lesen

- [ ] **NAS_AI_SYSTEM.md** gelesen
  - Pfad: `/home/user/Agent/NAS_AI_SYSTEM.md`
  - Fokus: Grundlegendes Verständnis der Systemarchitektur und Governance.

- [ ] **AGENT_MATRIX.md** gelesen
  - Pfad: `docs/planning/AGENT_MATRIX.md`
  - Fokus: Verständnis der eigenen Rolle und der Rollen anderer Agenten.

- [ ] **DEV_GUIDE.md** gelesen
  - Pfad: `docs/development/DEV_GUIDE.md`
  - Fokus: Verständnis der Entwicklungsprozesse, Werkzeuge und Konventionen.

- [ ] **SECURITY_HANDBOOK.md** gelesen
  - Pfad: `docs/security/SECURITY_HANDBOOK.pdf`
  - Fokus: Verständnis der Sicherheitsrichtlinien, insbesondere im Umgang mit Secrets.

### 2. Agent-spezifische Dokumente

- [ ] **Eigene Status-Logs** durchgegangen
  - Pfad: `status/<AgentName>/`
  - Letzte 3 Einträge gelesen
  - Offene Blocker identifiziert

- [ ] **Backlog geprüft**
  - Pfad: `status/backlog/`
  - Relevante Einträge für meinen Bereich gefunden
  - Kontext im neuen Statuslog vermerkt

### 3. Relevante Referenzen (bei Bedarf)

- [ ] **CVE_CHECKLIST.md** (bei Security-Tasks)
  - Pfad: `CVE_CHECKLIST.md` (im Root-Verzeichnis)
  - Offene CVEs für meinen Bereich geprüft

- [ ] **Blueprint-Dokumente** (bei UI/API-Tasks)
  - Pfad: `docs/blueprints/`
  - Relevante Blueprints identifiziert

- [ ] **ADR-Dokumente** (bei Architektur-Tasks)
  - Pfad: `docs/adr/`
  - Relevante Entscheidungen geprüft

---

## 📋 ARBEITSABLAUF (STRIKT EINHALTEN)

### Phase 1: VERSTEHEN ✅

- [ ] **Aufgabe gelesen und verstanden**
  - Owner-Anweisung klar
  - Scope definiert
  - Deliverables identifiziert

- [ ] **Unklarheiten geklärt**
  - Keine Widersprüche in Dokumentation
  - Keine "Halluzinationen" vermutet
  - Bei Unsicherheit: `clarification_needed` an Orchestrator

### Phase 2: ANALYSIEREN ✅

- [ ] **Ist-Zustand erfasst**
  - Aktueller Service-Status bekannt
  - Bestehende Dateien/Konfigurationen geprüft
  - Dependencies identifiziert

- [ ] **Risiken identifiziert**
  - Security-Risiken bewertet
  - Breaking Changes erkannt
  - Rollback-Strategie überlegt

- [ ] **Ressourcen ermittelt**
  - Benötigte Dateien bekannt
  - Erforderliche Zugriffsrechte klar
  - Zeitbedarf geschätzt

### Phase 3: PLANEN ✅

- [ ] **TODO-Liste erstellt**
  - Konkrete Schritte definiert
  - Reihenfolge festgelegt
  - Verifikationspunkte eingeplant

- [ ] **Statuslog angelegt**
  - Datei: `status/<AgentName>/<NNN>_YYYYMMDD_<task>.md`
  - Schema: NNN fortlaufend
  - Inhalt: Aufgabe, Ist-Zustand, Risiken, TODO-Plan

- [ ] **Entscheidungspunkte dokumentiert**
  - Owner-Approval benötigt (wo?)
  - Alternativen aufgezeigt
  - Empfehlungen ausgesprochen

### Phase 4: UMSETZEN (ERST NACH APPROVAL)

- [ ] **Owner-Approval eingeholt**
  - Bei kritischen Änderungen
  - Bei Security-relevanten Tasks
  - Bei neuen Dependencies

- [ ] **Schrittweise Umsetzung**
  - TODO-Liste abarbeiten
  - Nach jedem Schritt: Statuslog aktualisieren
  - Bei Problemen: Sofort dokumentieren

- [ ] **Evidenz sammeln**
  - Logs sichern
  - Screenshots/Output speichern
  - Artefakte verlinken

### Phase 5: VERIFIZIEREN ✅

- [ ] **Funktionstest durchgeführt**
  - Happy Path getestet
  - Error Cases geprüft
  - Performance akzeptabel

- [ ] **Dokumentation aktualisiert**
  - README-Dateien angepasst
  - Pfade korrigiert
  - Beispiele funktional

- [ ] **Statuslog finalisiert**
  - Alle Schritte dokumentiert
  - Evidenz-Links eingetragen
  - "Terminal freigegeben" Zeile hinzugefügt

---

## 🚫 VERBOTEN

### Diese Aktionen sind NICHT erlaubt ohne explizite Freigabe:

- ❌ **Raten statt Klären**
  - Bei Unsicherheit: `clarification_needed`
  - NIEMALS eigene Annahmen treffen

- ❌ **Umsetzung vor Analyse**
  - IMMER: Lesen → Analysieren → TODO → Umsetzung
  - NIEMALS direkt mit Code-Änderungen beginnen

- ❌ **Eigene Ordner erstellen**
  - IMMER: `docs/STRUCTURE_PROPOSAL.md` konsultieren
  - NIEMALS: Wilde Ordnerstrukturen

- ❌ **Secrets committen**
  - IMMER: Platzhalter verwenden
  - NIEMALS: Echte Tokens/Keys in Docs

- ❌ **Systemd/Services ohne Approval**
  - IMMER: Owner fragen bei Service-Restarts
  - NIEMALS: Produktions-Services ohne Absprache ändern

- ❌ **Backups anfassen**
  - IMMER: `/mnt/raid/backups` unangetastet lassen
  - NIEMALS: Backup-Dateien verschieben/ändern

---

## ❓ FAQ

**Q: Muss ich wirklich ALLE Kern-Dokumente vor JEDER Aufgabe lesen?**
A: Ja. Die Pflichtlektüre (`NAS_AI_SYSTEM.md`, `AGENT_MATRIX.md`, `DEV_GUIDE.md`, `SECURITY_HANDBOOK.pdf`) ist bei jeder neuen Anweisung erforderlich.

**Q: Was wenn ich etwas nicht verstehe?**
A: SOFORT `clarification_needed` an den Orchestrator senden. NIEMALS raten oder Annahmen treffen.

**Q: Kann ich die Analyse-Phase überspringen wenn die Aufgabe einfach ist?**
A: Nein. Der Arbeitsablauf (Lesen → Analysieren → Planen → Umsetzen) ist STRIKT einzuhalten. Keine Abkürzungen.

**Q: Wo erstelle ich neue Dateien?**
A: IMMER erst `docs/STRUCTURE_PROPOSAL.md` konsultieren. Das Erstellen eigener Ordnerstrukturen ohne Absprache ist verboten.