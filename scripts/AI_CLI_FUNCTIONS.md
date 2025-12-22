# AI Knowledge Agent - CLI Funktionen

Die nas-cli.sh wurde um umfassende AI Knowledge Agent Management-Funktionen erweitert.

## Neue Menü-Optionen

### 11) 🗑️ AI DB Clear (Embeddings löschen)
Löscht die komplette `file_embeddings` Tabelle und erstellt sie neu.

**Anwendungsfall:**
- Bereinigung alter/fehlerhafter Embeddings
- Reset vor Neuindizierung aller Dateien

**Ablauf:**
1. Sicherheitsabfrage (y/N)
2. DROP TABLE file_embeddings CASCADE
3. CREATE TABLE mit korrektem Schema (vector(384))

---

### 12) 🔄 AI Agent Restart
Startet den AI Knowledge Agent Container neu.

**Anwendungsfall:**
- Nach Konfigurationsänderungen
- Bei hängenden Prozessen
- Neuladen des ML-Models

**Ablauf:**
1. docker compose restart ai-knowledge-agent
2. Wartet 20s für Model-Loading
3. Zeigt die letzten 15 Log-Zeilen

---

### 13) 🔨 AI Agent Rebuild (neu bauen)
Kompletter Rebuild des AI Agent Images und Containers.

**Anwendungsfall:**
- Nach Code-Änderungen in `ai_knowledge_agent/src/`
- Update von Dependencies (requirements.txt)
- Nach Dockerfile-Anpassungen

**Ablauf:**
1. Stoppe & lösche Container
2. Build mit --no-cache
3. Starte neuen Container
4. Wartet 25s für Model-Loading
5. Zeigt Status & Logs

---

### 14) 🧠 AI Full Reset (DB + Restart)
Kombiniert DB-Bereinigung + Agent-Neustart.

**Anwendungsfall:**
- Kompletter Neustart der AI-Funktionalität
- Nach größeren Änderungen
- Troubleshooting bei Problemen

**Ablauf:**
1. Führt ai_db_clear() aus
2. Führt ai_agent_restart() aus

---

### 15) 📊 AI Embeddings anzeigen
Zeigt alle gespeicherten Embeddings aus der Datenbank.

**Ausgabe:**
```
file_id        | mime_type  | content_length | embedding_dim | created_at
--------------+------------+----------------+---------------+------------
integration_test.txt | text/plain | 370 | 384 | 2025-11-29 20:20:16
test_ai_knowledge.txt | text/plain | 551 | 384 | 2025-11-29 20:18:28
```

**Anwendungsfall:**
- Übersicht über indexierte Dateien
- Debugging
- Monitoring

---

### 16) 🧪 AI Endpoint Test
Testet den AI Agent `/process` Endpoint mit einer Datei.

**Interaktiv:**
- Fragt nach Dateiname in /mnt/data/
- Default: test.txt

**Ablauf:**
1. Erstellt JSON Payload
2. Sendet POST zu http://ai-knowledge-agent:5000/process
3. Zeigt Response
4. Zeigt AI Agent Logs (letzte 10 Zeilen)

**Beispiel Payload:**
```json
{
  "file_path": "/mnt/data/test.txt",
  "file_id": "test.txt",
  "mime_type": "text/plain"
}
```

---

## Technische Details

### Datenbankschema
```sql
CREATE TABLE file_embeddings (
  id SERIAL PRIMARY KEY,
  file_id TEXT NOT NULL,
  file_path TEXT NOT NULL,
  mime_type TEXT NOT NULL,
  content TEXT NOT NULL,
  embedding vector(384),  -- pgvector extension
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(file_id)
);
```

### Container-Abhängigkeiten
- **Volume Mount**: `nas_data:/mnt/data:ro` (read-only)
- **Network**: `nas-network`
- **Depends on**: postgres (healthy)
- **Model**: sentence-transformers/all-MiniLM-L6-v2
- **Port**: 5000 (intern)

### Environment Variables
```bash
PGHOST=postgres
PGPORT=5432
PGDATABASE=nas_db
PGUSER=nas_user
PGPASSWORD=<from .env.prod>
```

---

## Verwendung

```bash
cd /home/freun/Agent
./scripts/nas-cli.sh

# Oder direkt:
bash /home/freun/Agent/scripts/nas-cli.sh
```

## Workflow-Beispiele

### Nach Code-Änderungen am AI Agent:
```
1. Option 13 wählen (AI Agent Rebuild)
2. Option 15 wählen (Embeddings prüfen)
3. Option 16 wählen (Endpoint testen)
```

### Kompletter Neustart:
```
1. Option 14 wählen (AI Full Reset)
2. Option 15 wählen (Verifizieren dass DB leer ist)
```

### Nur DB bereinigen:
```
1. Option 11 wählen (AI DB Clear)
```

---

## Integration mit dem System

Die AI Agent Funktionen nutzen die bestehende Infrastruktur:

- **Docker Compose**: `$INFRA_DIR/docker-compose.prod.yml`
- **Postgres**: nas-api-postgres Container
- **Storage**: Shared volume mit API Container

Alle Funktionen sind non-breaking und beeinflussen nicht:
- API Server
- WebUI
- Andere Services (monitoring, analysis, pentester)


---

## Neue Features - Semantic Search (v2)

### 17) 🔍 AI /embed_query Test
Testet den neuen `/embed_query` Endpoint des AI Agents.

**Funktion:**
- Generiert Embeddings für beliebigen Text
- Ohne File-Upload erforderlich
- Direkte Embedding-Generierung

**Interaktiv:**
- Fragt nach Suchtext
- Sendet POST zu http://ai-knowledge-agent:5000/embed_query
- Zeigt das generierte Embedding (384-dimensional)

**Beispiel Request:**
```json
{
  "text": "Wie hoch sind die Server-Kosten?"
}
```

**Beispiel Response:**
```json
{
  "embedding": [0.123, -0.456, 0.789, ...]
}
```

---

### 18) 🎯 Semantic Search Test
Testet die End-to-End semantische Suche über die API.

**Workflow:**
1. User gibt Suchanfrage ein (z.B. "Server Kosten")
2. API sendet Query an AI Agent für Embedding
3. API führt pgvector Similarity Search durch
4. Zeigt Top-10 Ergebnisse sortiert nach Ähnlichkeit

**Ausgabe:**
```
Suchergebnisse:
{
  "query": "Server Kosten",
  "results": [
    {
      "file_path": "/mnt/data/rechnung.txt",
      "content": "Die monatlichen Server-Kosten betragen...",
      "similarity": 0.87
    }
  ]
}

Top 3 Resultate:
Datei: /mnt/data/rechnung.txt
Ähnlichkeit: 0.87
Inhalt: Die monatlichen Server-Kosten betragen...
```

**SQL Query (intern):**
```sql
SELECT 
  file_path, 
  content, 
  1 - (embedding <=> $query_embedding::vector) as similarity
FROM file_embeddings
ORDER BY embedding <=> $query_embedding::vector
LIMIT 10;
```

---

### 19) 🚀 AI+API Full Rebuild
Kompletter Rebuild beider Services mit allen neuen Features.

**Was wird gemacht:**
1. Stoppe ai-knowledge-agent und api Container
2. Rebuild AI Agent (--no-cache) mit /embed_query
3. Rebuild API mit Semantic Search Handler
4. Starte beide Container neu
5. Wartet 30s für Initialization
6. Zeigt Logs beider Services

**Anwendungsfall:**
- Nach größeren Code-Änderungen
- Deployment neuer Features
- Full Integration Test

**Neue API Config:**
```go
AI_SERVICE_URL=http://ai-knowledge-agent:5000  // Default
```

---

## API Endpoints (neu)

### GET /api/v1/search?q={query}
Semantische Suche über alle indizierten Dokumente.

**Request:**
```bash
curl "http://localhost:8080/api/v1/search?q=server%20kosten"
```

**Response:**
```json
{
  "query": "server kosten",
  "results": [
    {
      "file_path": "/mnt/data/rechnung.txt",
      "content": "Die monatlichen Server-Kosten...",
      "similarity": 0.87
    }
  ]
}
```

**Status Codes:**
- 200: Success
- 400: Missing query parameter
- 502: AI agent unavailable
- 500: Database error

---

## AI Agent Endpoints (aktualisiert)

### POST /process
Indexiert eine Datei (bestehendes Feature)

### POST /embed_query (NEU)
Generiert Embedding für beliebigen Text

**Request:**
```json
{
  "text": "Meine Suchanfrage"
}
```

**Response:**
```json
{
  "embedding": [0.123, -0.456, ...]
}
```

### GET /health
Health Check

---

## Technische Architektur

```
┌─────────────┐
│   WebUI     │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│           API Server                │
│  • GET /api/v1/search?q=...        │
│  • SearchHandler                    │
└────┬─────────────────────┬─────────┘
     │                     │
     ▼                     ▼
┌──────────────┐    ┌────────────────┐
│ AI Agent     │    │   PostgreSQL   │
│ /embed_query │    │   + pgvector   │
└──────────────┘    └────────────────┘
```

**Datenfluss:**
1. User → API: `GET /search?q=server kosten`
2. API → AI: `POST /embed_query {"text": "server kosten"}`
3. AI → API: `{"embedding": [...]}`
4. API → DB: Similarity search mit pgvector
5. DB → API: Top 10 ähnliche Dokumente
6. API → User: JSON mit Ergebnissen

---

## Performance

- **Embedding Generation**: ~200-600ms (CPU-only)
- **DB Similarity Search**: ~50-200ms (abhängig von Datenmenge)
- **Total Response Time**: ~300-800ms

**Optimierungen:**
- AI Agent Model wird beim Start geladen (einmalig)
- Shared HTTP Client für AI requests
- pgvector Index für schnelle Similarity Search
- Read-only Volume Mount für Storage

