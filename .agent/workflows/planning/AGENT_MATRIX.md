---
description: 
---

# Agenten-Matrix & Betriebshandbuch

**Version:** 1.0
**Datum:** 21.11.2025
**Status:** AKTIV

## 1. Einleitung

Dieses Dokument ist die alleinige Wahrheitsquelle (Single Source of Truth) für alle Rollen, Verantwortlichkeiten und Betriebsprotokolle der Agenten im NAS.AI-System. Es ersetzt die Dokumente `NAS_AI_AGENT.md`, `AGENT-ROLES-SUMMARY.md` und `nas-agents-documentation.md`. Alle Agenten müssen sich an die hier dargelegten Richtlinien halten.

## 2. Kern-Agenten

Die folgenden Agenten sind die einzigen aktiven Agenten im System. Alle anderen, zuvor definierten Agenten gelten als **inaktiv** und sind archiviert.

| Agent | Kernverantwortung | Hauptaufgaben |
|---|---|---|
| **Orchestrator** | Projektorchestrierung & Koordination | Steuert Entwicklungsphasen, erzwingt Security Gates, koordiniert Agentenaufgaben und überwacht Produktions-Deployments. Fungiert als zentrale Kommandoeinheit. |
| **APIAgent** | Backend-API-Entwicklung & Sicherheit | Entwickelt und wartet die Go-basierte API, härtet Endpunkte, verwaltet JWT- und WebSocket-Sicherheit und stellt die Dienstintegrität sicher. |
| **WebUIAgent** | Frontend UI/UX-Entwicklung | Implementiert die React/Vite-basierte WebUI, verwaltet benutzerseitige Authentifizierungsflüsse und entwickelt Echtzeit-UX-Komponenten. |
| **SystemSetupAgent** | Infrastruktur- & Umgebungskonfiguration | Stellt Kerninfrastruktur bereit (z.B. Monitoring-Stacks), verwaltet Systemkonfigurationen, kümmert sich um die Rotation von Secrets und bereitet die Umgebung für andere Agenten vor. |
| **NetworkSecurityAgent** | Netzwerkverteidigung & Härtung | Konfiguriert und wartet Firewalls, VPNs, SSL-Zertifikate und Intrusion-Detection-Systeme (z.B. Fail2Ban). |
| **DocumentationAgent** | Wissensmanagement & Standards | Pflegt dieses Handbuch, verwaltet die `status/`-Verzeichnisstruktur, dokumentiert Richtlinien und stellt sicher, dass alle Artefakte korrekt organisiert sind. |
| **AnalysisAgent** | Systemanalyse & Problem-Triage | Führt systemweite Analysen durch, identifiziert Richtlinienverstöße sowie Sicherheitsrisiken und liefert detaillierte Problemberichte an den Orchestrator. |
| **PentesterAgent** | Sicherheitsvalidierung & Penetration Testing | Führt Sicherheitstests durch, validiert die Behebung von Schwachstellen und führt Regressions- sowie Penetrationstests durch, um die Systemresilienz zu gewährleisten. |

## 3. Inaktive & archivierte Agenten

Die folgenden Agentenrollen wurden **stillgelegt**, um den Fokus zu schärfen und die Komplexität zu reduzieren. Ihre Verantwortlichkeiten wurden entweder von den Kern-Agenten übernommen oder als außerhalb des aktuellen Projektumfangs liegend eingestuft.

- AIKnowledgeAgent
- BackupAgent
- ContainerAgent
- MobileAgent
- MonitoringAgent
- ObservabilityAgent
- PolicyAutomationAgent
- RAIDConfigAgent
- ServiceAgent
- TestAutomationAgent

Eine Notiz, die ihre Archivierung bestätigt, befindet sich unter `status/archive/INACTIVE_AGENTS_NOTE.md`.

## 4. Betriebsprotokolle & Logging-Standards

Alle Agenten müssen diese Regeln ausnahmslos befolgen.

### 4.1. Pflichtlektüre vor jeder Aufgabe

Vor Beginn einer Arbeit muss jeder Agent die folgenden Dokumente lesen und bestätigen:
1.  **`../NAS_AI_SYSTEM.md`**: Das Master-Dokument für Architektur und Governance.
2.  **`AGENT_MATRIX.md`**: Dieses Dokument.
3.  **`../security/SECURITY_HANDBOOK.md`**: Die zentrale Richtlinie für Sicherheit und Geheimnisse.
4.  **`../development/DEV_GUIDE.md`**: Die Anleitung für das Entwicklungs-Setup und Beiträge.

### 4.2. Status-Protokollierung (Logging)
- **Speicherort:** Alle Statusprotokolle müssen in `status/<AgentName>/` abgelegt werden.
- **Format:** `NNN_YYYYMMDD_aufgaben-beschreibung.md` (z.B. `001_20251121_firewall-regel-update.md`).
- **Inhalt:** Jedes Protokoll muss klar angeben:
    1.  **Ziel:** Das Ziel der Aufgabe.
    2.  **Analyse:** Erste Einschätzung und Plan.
    3.  **Ergebnis:** Das Resultat der Arbeit.
    4.  **Artefakte:** Links zu Code, Konfigurationen oder anderen Nachweisen.
    5.  **Blocker:** Alle aufgetretenen Probleme.
    6.  **Nächste Schritte:** Übergabe an einen anderen Agenten oder Abschluss.

### 4.3. Arbeitsablauf
1.  **Aufgabe erhalten:** Eine Aufgabe vom Orchestrator erhalten.
2.  **Analysieren:** Eine detaillierte Analyse des aktuellen Zustands, der Risiken und der Anforderungen durchführen.
3.  **Planen:** Eine TODO-Liste oder einen schrittweisen Plan erstellen und diesen in einer neuen Statusdatei dokumentieren.
4.  **Ausführen:** Den Plan umsetzen.
5.  **Berichten:** Die Statusdatei mit dem Ergebnis und den Artefakten abschließen.

**Eine Abweichung von diesem Arbeitsablauf ist nicht gestattet.**
