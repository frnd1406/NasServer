# 🩺 SRP Analyse Report: Code Health Check

Dieser Bericht identifiziert Dateien, die gegen das Single Responsibility Principle (SRP) verstoßen und refactored werden sollten.

### 1. `infrastructure/ai_knowledge_agent/src/main.py`
- **Zeilen (geschätzt):** ~1274
- **Hauptaufgaben:** Stack-Overflow! Flask-Server, DB-Verbindung, Schema-Check, API-Endpunkte, RAG-Logik, Ollama-Integration, Hintergrund-Tasks, Auth-Middleware.
- **Verstöße:**
    - **God-Object:** Macht *alles*. HTTP, SQL, Business-Logik, Background-Jobs.
    - **Hardcoded SQL:** SQL-Statements (`INSERT`, `SELECT`) direkt im Code verstreut.
    - **Business-Logik in Controller:** RAG-Logik (Prompt-Building) direkt in Route-Handlern.
    - **Global State:** Verlässt sich auf globale Variablen für DB-Pools und Status.
- **Urteil:** 🔴 **SPLIT** (Höchste Priorität!)
    - **Vorschlag:** Aufteilen in:
        1. `routes.py` (Nur HTTP Routing)
        2. `database.py` (DB Connection & Queries)
        3. `rag_service.py` (Core Business Logic)
        4. `ollama_client.py` (External Service Adapter)

### 2. `infrastructure/api/src/services/storage_service.go`
- **Zeilen (geschätzt):** ~808
- **Hauptaufgaben:** Datei-Kopieren, Pfad-Validierung, MIME-Type-Detection, Encryption, Versionierung, Trash.
- **Verstöße:**
    - **Vermischung von Domains:** Core-Storage vs. Encryption vs. Trash.
    - **Komplexität:** `SaveWithEncryption` ist extrem komplex und schwer isoliert zu testen.
    - **Magische Werte:** Magic-Bytes für Dateitypen hardcodiert.
- **Urteil:** 🟡 **REFACTOR**
    - **Vorschlag:** Auslagern von Encryption in `crypto_service.go` und Trash-Logik in `trash_service.go`.

### 3. `infrastructure/webui/src/pages/Files.jsx`
- **Zeilen (geschätzt):** 440+
- **Hauptaufgaben:** View-Controller für Dateimanager, State-Management, Event-Handling.
- **Verstöße:**
    - **Logic in View:** Direkter `fetch` zum Löschen des Papierkorbs (Zeile 289).
    - **Fat Controller:** Zu viel UI-State in einer Datei.
- **Urteil:** 🟡 **REFACTOR**
    - **Vorschlag:** `fetch` in Hook (`useFileStorage`) verschieben. UI in Sub-Komponenten (`FilesHeader`, `FilesActionPanel`) aufbrechen.

### 4. `infrastructure/webui/src/components/EnhancedChatWidget.jsx`
- **Zeilen (geschätzt):** ~210
- **Hauptaufgaben:** Chat-UI, Message-State, API-Kommunikation.
- **Verstöße:**
    - **Logic in View:** `sendMessage` enthält API-Logik.
- **Urteil:** 🟢 **KEEP** (Beobachten)
    - **Vorschlag:** Noch okay. Wenn >300 Zeilen, Custom Hook `useChat` extrahieren.

### 5. `infrastructure/api/src/handlers/settings.go`
- **Zeilen (geschätzt):** ~168
- **Hauptaufgaben:** Request-Validierung, Service-Aufruf.
- **Verstöße:**
    - **Direct IO:** `ValidatePathHandler` greift direkt auf Filesystem zu (statt Service). Vertretbar für Utility-Endpoint.
- **Urteil:** 🟢 **KEEP**
