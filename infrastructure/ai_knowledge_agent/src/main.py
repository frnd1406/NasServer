import logging
import os
import time
import json
import threading
from contextlib import contextmanager

import requests
from flask import Flask, jsonify, request
import psycopg2
from psycopg2 import pool, OperationalError
from pgvector.psycopg2 import register_vector
import numpy as np
import redis as redis_lib

# Support both container (src.intent_classifier) and local execution
try:
    from src.intent_classifier import classify_intent, get_limit_for_intent
except ImportError:
    from intent_classifier import classify_intent, get_limit_for_intent


logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("ai_knowledge_agent")

app = Flask(__name__)

# Ollama Configuration
OLLAMA_URL = os.getenv("OLLAMA_URL", "http://host.docker.internal:11434")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "mxbai-embed-large")
LLM_MODEL = os.getenv("LLM_MODEL", "llama3.2")
EMBEDDING_DIM = 1024  # mxbai-embed-large dimension

# SECURITY: Auto-indexing disabled by default for Zero-Trust
# Only push-based ingestion via /process endpoint is allowed
AUTO_INDEX_ENABLED = os.getenv("AUTO_INDEX_ENABLED", "false").lower() == "true"

# Redis Configuration for Job Queue
REDIS_URL = os.getenv("REDIS_URL", "redis://redis:6379")
JOB_STREAM = "ai:jobs"
CONSUMER_GROUP = "ai-workers"
RESULT_KEY_PREFIX = "ai:results:"
JOB_TIMEOUT = 120  # seconds
WORKER_ENABLED = os.getenv("WORKER_ENABLED", "true").lower() == "true"

model_loaded = False
db_pool = None
redis_client = None

def get_db_password():
    """Read database password from secret file or env var."""
    pwd_file = os.getenv("PGPASSWORD_FILE") or os.getenv("POSTGRES_PASSWORD_FILE")
    if pwd_file and os.path.exists(pwd_file):
        try:
            with open(pwd_file, 'r') as f:
                return f.read().strip()
        except Exception as e:
            logger.warning("Failed to read password file: %s", e)
    return os.getenv("PGPASSWORD", "")

# FIX [BUG-PY-006]: Module-level constant for DSN
DSN = {
    "host": os.getenv("PGHOST", "postgres"),
    "port": os.getenv("PGPORT", "5432"),
    "dbname": os.getenv("PGDATABASE", "nas_db"),
    "user": os.getenv("PGUSER", "nas_user"),
    "password": get_db_password(),
}


def get_ollama_embedding(text: str) -> list:
    """Get embedding from Ollama mxbai-embed-large model."""
    try:
        response = requests.post(
            f"{OLLAMA_URL}/api/embeddings",
            json={"model": EMBEDDING_MODEL, "prompt": text},
            timeout=300
        )
        response.raise_for_status()
        return response.json()["embedding"]
    except Exception as e:
        logger.error("Ollama embedding failed: %s", e)
        raise


def get_llama_response(prompt: str, documents: list) -> dict:
    """
    Get intelligent response from Llama 3.2.
    AI decides which documents are relevant and how many to cite.
    Returns: {"answer": str, "cited_sources": list, "confidence": str}
    """
    system_prompt = """Du bist ein intelligenter KI-Assistent für ein NAS-Dokumentensystem.

DEINE AUFGABE:
1. Analysiere die bereitgestellten Dokumente und entscheide SELBST, welche relevant sind
2. Zitiere NUR die Dokumente, die wirklich zur Frage passen (0-5, je nach Relevanz)
3. Gib eine DIREKTE Antwort wenn die Information vorhanden ist (z.B. "Der Server kostet 149,99€")
4. Wenn keine passenden Dokumente existieren, sag das ehrlich

ANTWORT-FORMAT (strikt einhalten):
---
RELEVANTE QUELLEN: [Liste der wirklich relevanten Dateinamen, oder "Keine"]
KONFIDENZ: [HOCH/MITTEL/NIEDRIG]
ANTWORT: [Deine direkte Antwort mit Quellenverweisen wie [Dok1]]
---

BEISPIELE:
- Frage "Was kostet der Server?" → Wenn Rechnung 149,99€ zeigt: "Der Server kostet 149,99€ [rechnung_xyz.txt]"
- Frage nach etwas Unbekanntem → "Dazu habe ich keine Informationen in den Dokumenten gefunden."
- Mehrere relevante Docs → "Basierend auf [dok1] und [dok2]: ..."

Antworte IMMER auf Deutsch. Sei präzise und direkt."""

    # Build context with document info
    context_parts = []
    for i, doc in enumerate(documents, 1):
        sim_percent = int(doc['similarity'] * 100)
        context_parts.append(f"[Dokument {i}: {doc['file_id']} (Ähnlichkeit: {sim_percent}%)]\n{doc['content'][:2000]}")
    
    context = "\n\n---\n\n".join(context_parts)

    full_prompt = f"""Hier sind {len(documents)} Dokumente aus der Datenbank:

{context}

---

FRAGE DES BENUTZERS: {prompt}

Analysiere die Dokumente und antworte im vorgegebenen Format:"""

    try:
        response = requests.post(
            f"{OLLAMA_URL}/api/generate",
            json={
                "model": LLM_MODEL,
                "prompt": full_prompt,
                "system": system_prompt,
                "stream": False,
                "options": {"temperature": 0.2, "num_predict": 800}
            },
            timeout=600
        )
        response.raise_for_status()
        raw_response = response.json()["response"]
        
        # Parse the structured response
        result = parse_llama_response(raw_response, documents)
        return result
        
    except Exception as e:
        logger.error("Llama generation failed: %s", e)
        raise


def parse_llama_response(raw: str, documents: list) -> dict:
    """Parse Llama's structured response into components."""
    result = {
        "answer": raw,
        "cited_sources": [],
        "confidence": "MITTEL",
        "raw_response": raw
    }
    
    try:
        # Extract RELEVANTE QUELLEN
        if "RELEVANTE QUELLEN:" in raw:
            sources_line = raw.split("RELEVANTE QUELLEN:")[1].split("\n")[0].strip()
            if sources_line.lower() != "keine" and sources_line != "[]":
                # Find matching documents
                for doc in documents:
                    if doc['file_id'] in sources_line or doc['file_id'].replace('.txt', '') in sources_line:
                        result["cited_sources"].append({
                            "file_id": doc['file_id'],
                            "file_path": doc['file_path'],
                            "similarity": doc['similarity']
                        })
        
        # Extract KONFIDENZ
        if "KONFIDENZ:" in raw:
            conf_line = raw.split("KONFIDENZ:")[1].split("\n")[0].strip().upper()
            if conf_line in ["HOCH", "MITTEL", "NIEDRIG"]:
                result["confidence"] = conf_line
        
        # Extract ANTWORT
        if "ANTWORT:" in raw:
            answer_part = raw.split("ANTWORT:")[1].strip()
            # Remove trailing --- if present
            if "---" in answer_part:
                answer_part = answer_part.split("---")[0].strip()
            result["answer"] = answer_part
            
    except Exception as e:
        logger.warning("Could not parse structured response: %s", e)
    
    return result


def check_ollama_health():
    """Check if Ollama is available and models are loaded."""
    global model_loaded
    try:
        response = requests.get(f"{OLLAMA_URL}/api/tags", timeout=300)
        if response.status_code == 200:
            models = [m["name"] for m in response.json().get("models", [])]
            embed_ok = any(EMBEDDING_MODEL in m for m in models)
            llm_ok = any(LLM_MODEL in m for m in models)
            model_loaded = embed_ok and llm_ok
            if model_loaded:
                logger.info("Ollama ready: %s + %s", EMBEDDING_MODEL, LLM_MODEL)
            else:
                logger.warning("Missing models. Available: %s", models)
            return model_loaded
    except Exception as e:
        logger.error("Ollama health check failed: %s", e)
        model_loaded = False
    return False


def prewarm_models():
    """Pre-warm Ollama models by making dummy requests at startup.
    This forces Ollama to load models into memory so first user request is fast.
    """
    global model_loaded
    logger.info("Pre-warming Ollama models...")
    
    # Import classifier model name
    try:
        from src.intent_classifier import CLASSIFIER_MODEL
    except ImportError:
        from intent_classifier import CLASSIFIER_MODEL
    
    try:
        # Pre-warm embedding model
        logger.info("Loading embedding model: %s", EMBEDDING_MODEL)
        _ = get_ollama_embedding("Warmup test for embedding model.")
        logger.info("✅ Embedding model loaded and ready")
        
        # Pre-warm classifier model (small, fast)
        logger.info("Loading classifier model: %s", CLASSIFIER_MODEL)
        response = requests.post(
            f"{OLLAMA_URL}/api/generate",
            json={
                "model": CLASSIFIER_MODEL,
                "prompt": "Hello",
                "stream": False,
                "options": {"num_predict": 1}
            },
            timeout=120
        )
        response.raise_for_status()
        logger.info("✅ Classifier model loaded and ready")
        
        # Pre-warm LLM model with a simple prompt
        logger.info("Loading LLM model: %s", LLM_MODEL)
        response = requests.post(
            f"{OLLAMA_URL}/api/generate",
            json={
                "model": LLM_MODEL,
                "prompt": "Hello",
                "stream": False,
                "options": {"num_predict": 1}  # Minimal response
            },
            timeout=300
        )
        response.raise_for_status()
        logger.info("✅ LLM model loaded and ready")
        
        model_loaded = True
        logger.info("🚀 All 3 models pre-warmed and ready for instant responses!")
        return True
        
    except Exception as e:
        logger.error("Model pre-warming failed: %s", e)
        model_loaded = False
        return False


def background_health_check():
    """Run periodic health checks and auto-index new files every 30 seconds.
    - Monitors Ollama availability
    - Scans for new files and indexes them automatically
    """
    global model_loaded
    while True:
        time.sleep(30)  # Check every 30 seconds
        
        # Health check Ollama
        try:
            response = requests.get(f"{OLLAMA_URL}/api/tags", timeout=10)
            if response.status_code == 200:
                models = [m["name"] for m in response.json().get("models", [])]
                embed_ok = any(EMBEDDING_MODEL in m for m in models)
                llm_ok = any(LLM_MODEL in m for m in models)
                model_loaded = embed_ok and llm_ok
                if not model_loaded:
                    logger.warning("Background check: Models missing, attempting re-warm...")
                    prewarm_models()
        except Exception as e:
            logger.warning("Background health check failed: %s", e)
            model_loaded = False
        
        # Auto-index new files (DISABLED for Zero-Trust Security)
        # Only explicit push via /process endpoint should ingest data
        if AUTO_INDEX_ENABLED and model_loaded and db_pool:
            try:
                new_files = index_all_files()
                if new_files > 0:
                    logger.info("🆕 Auto-indexed %d new files", new_files)
            except Exception as e:
                logger.warning("Background auto-index failed: %s", e)
        # else: Auto-index disabled - secure push-only mode


def start_background_health_check():
    """Start the background health check thread."""
    thread = threading.Thread(target=background_health_check, daemon=True)
    thread.start()
    logger.info("Background health check thread started")


def index_all_files(data_dir="/mnt/data"):
    """Scan and index all text files in the data directory.
    This runs at startup to ensure all files are searchable.
    """
    if not model_loaded:
        logger.warning("Cannot index files - models not loaded")
        return 0
    
    if not os.path.exists(data_dir):
        logger.warning("Data directory does not exist: %s", data_dir)
        return 0
    
    # Supported file extensions
    TEXT_EXTENSIONS = {'.txt', '.md', '.json', '.csv', '.log', '.xml', '.html', '.py', '.js', '.go', '.sh'}
    
    indexed_count = 0
    skipped_count = 0
    
    logger.info("📂 Starting auto-indexing of files in %s...", data_dir)
    
    # Get already indexed files from database
    existing_files = set()
    try:
        with get_db_connection() as conn:
            with conn.cursor() as cur:
                cur.execute("SELECT file_id FROM file_embeddings")
                existing_files = {row[0] for row in cur.fetchall()}
        logger.info("Found %d already indexed files", len(existing_files))
    except Exception as e:
        logger.warning("Could not check existing files: %s", e)
    
    # Walk through all files
    for root, dirs, files in os.walk(data_dir):
        # Skip hidden directories and trash
        dirs[:] = [d for d in dirs if not d.startswith('.') and d != '.trash']
        
        for filename in files:
            # Skip hidden files
            if filename.startswith('.'):
                continue
                
            file_path = os.path.join(root, filename)
            file_ext = os.path.splitext(filename)[1].lower()
            
            # Skip non-text files
            if file_ext not in TEXT_EXTENSIONS:
                continue
            
            # Skip already indexed files
            if filename in existing_files:
                skipped_count += 1
                continue
            
            try:
                # Read file content
                with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
                    content = f.read()
                
                if not content.strip():
                    continue
                
                # Get embedding
                logger.info("Indexing: %s", filename)
                embedding = get_ollama_embedding(content[:8000])
                
                # Store in database
                metadata = {
                    "file_path": file_path,
                    "mime_type": "text/plain",
                    "content_length": len(content)
                }
                
                with get_db_connection() as conn:
                    register_vector(conn)
                    with conn.cursor() as cur:
                        cur.execute("""
                            INSERT INTO file_embeddings (file_id, chunk_index, content, embedding, metadata)
                            VALUES (%s, %s, %s, %s, %s)
                            ON CONFLICT (file_id, chunk_index) DO UPDATE SET
                                content = EXCLUDED.content,
                                embedding = EXCLUDED.embedding,
                                metadata = EXCLUDED.metadata,
                                created_at = CURRENT_TIMESTAMP
                        """, (filename, 0, content, embedding, json.dumps(metadata)))
                        conn.commit()
                
                indexed_count += 1
                logger.info("✅ Indexed: %s", filename)
                
            except Exception as e:
                logger.error("Failed to index %s: %s", filename, e)
    
    logger.info("📊 Auto-indexing complete: %d new files indexed, %d skipped (already indexed)", 
                indexed_count, skipped_count)
    return indexed_count


def init_db_pool(max_retries=5, base_delay=1.0):
    """Initialize database connection pool with retry logic."""
    global db_pool
    
    for attempt in range(1, max_retries + 1):
        try:
            db_pool = psycopg2.pool.SimpleConnectionPool(
                minconn=1,
                maxconn=10,
                **DSN
            )
            logger.info("Database connection pool initialized (attempt %d/%d).", attempt, max_retries)
            
            # Create schema with 1024D for mxbai-embed-large
            conn = db_pool.getconn()
            try:
                register_vector(conn)
                with conn.cursor() as cur:
                    # Schema is already created by SQL migration - just verify it exists
                    cur.execute("""
                        SELECT EXISTS (
                            SELECT FROM information_schema.tables
                            WHERE table_name = 'file_embeddings'
                        )
                    """)
                    exists = cur.fetchone()[0]
                    if not exists:
                        logger.error("file_embeddings table does not exist!")
                        return False
                logger.info("Database schema verified (vector dim: %d).", EMBEDDING_DIM)
            finally:
                db_pool.putconn(conn)
            
            return True
            
        except OperationalError as e:
            delay = base_delay * (2 ** (attempt - 1))
            logger.warning(
                "DB connection failed (attempt %d/%d): %s. Retrying in %.1fs...",
                attempt, max_retries, e, delay
            )
            if attempt < max_retries:
                time.sleep(delay)
            else:
                logger.error("Failed to connect to database after %d attempts", max_retries)
                return False
                
        except Exception as e:
            logger.error("Unexpected error initializing DB pool: %s", e)
            return False
    
    return False


@contextmanager
def get_db_connection():
    global db_pool
    if not db_pool:
        logger.warning("Database pool not initialized. Attempting lazy initialization...")
        if not init_db_pool(max_retries=3):
             raise Exception("Database pool not initialized (lazy init failed)")
    
    conn = db_pool.getconn()
    try:
        yield conn
    finally:
        db_pool.putconn(conn)


# SECURITY: Shared Secret for Internal Communication
INTERNAL_API_SECRET = os.getenv("INTERNAL_API_SECRET")

@app.before_request
def check_internal_auth():
    """Global middleware to require X-Internal-Secret header."""
    # Allow health checks without auth (optional, but good for Docker/K8s)
    if request.endpoint == 'health':
        return None
        
    # Skip auth for static files or OPTIONS if needed (none here)
    if request.method == 'OPTIONS':
        return None

    secret = request.headers.get('X-Internal-Secret')
    if not INTERNAL_API_SECRET:
            # Fail open or closed? Closed is safer.
            logger.error("INTERNAL_API_SECRET not set on server! Rejecting all requests.")
            return jsonify({"error": "Server misconfiguration"}), 500
            
    if not secret or secret != INTERNAL_API_SECRET:
        logger.warning("Unauthorized access attempt from %s to %s", request.remote_addr, request.path)
        return jsonify({"error": "Unauthorized"}), 403

@app.route("/health", methods=["GET"])
def health():

    """Health check endpoint."""
    db_ok = db_pool is not None
    ollama_ok = check_ollama_health()
    
    if not model_loaded:
        return jsonify({
            "status": "degraded",
            "ollama_ready": ollama_ok,
            "db_ready": db_ok,
            "message": "Ollama models not ready"
        }), 503
    
    return jsonify({
        "status": "ok" if (model_loaded and db_ok) else "degraded",
        "ollama_ready": ollama_ok,
        "db_ready": db_ok,
        "embedding_model": EMBEDDING_MODEL,
        "llm_model": LLM_MODEL,
        "embedding_dim": EMBEDDING_DIM
    })


@app.route("/status", methods=["GET"])
def status():
    """Comprehensive status endpoint for Settings page."""
    try:
        from src.intent_classifier import CLASSIFIER_MODEL
    except ImportError:
        from intent_classifier import CLASSIFIER_MODEL
    
    # Get Ollama models
    ollama_models = []
    ollama_connected = False
    try:
        resp = requests.get(f"{OLLAMA_URL}/api/tags", timeout=5)
        if resp.status_code == 200:
            ollama_connected = True
            ollama_models = [m["name"] for m in resp.json().get("models", [])]
    except Exception:
        pass
    
    # Get index stats from database
    total_files = 0
    indexed_files = 0
    try:
        with get_db_connection() as conn:
            with conn.cursor() as cur:
                cur.execute("SELECT COUNT(*) FROM file_embeddings")
                indexed_files = cur.fetchone()[0]
    except Exception:
        pass
    
    return jsonify({
        "ollama": {
            "connected": ollama_connected,
            "url": OLLAMA_URL,
            "models": ollama_models
        },
        "models": {
            "embedding": EMBEDDING_MODEL,
            "classifier": CLASSIFIER_MODEL,
            "llm": LLM_MODEL
        },
        "index": {
            "total_files": total_files,
            "indexed_files": indexed_files
        },
        "settings": {
            "auto_index": True,
            "embedding_dim": EMBEDDING_DIM
        }
    })


@app.route("/reindex", methods=["POST"])
def reindex():
    """Trigger re-indexing of all files."""
    try:
        # Run auto-indexing in background
        import threading
        def do_reindex():
            auto_index_files()
        
        thread = threading.Thread(target=do_reindex)
        thread.start()
        
        return jsonify({
            "status": "started",
            "message": "Re-indexing started in background"
        })
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/process", methods=["POST"])
def process_file():
    """Process an uploaded file.
    Supports two modes:
    1. Direct Text (Push): content is sent in JSON (for encrypted files)
    2. Disk Read (Pull): agent reads file from disk (legacy/public files)
    """
    if not model_loaded:
        logger.error("Ollama not ready")
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "Missing JSON payload"}), 400

        file_id = data.get("file_id")
        mime_type = data.get("mime_type")
        
        # NEU: Wir schauen, ob uns der Text direkt geschickt wurde
        direct_content = data.get("content") 
        file_path = data.get("file_path")

        if not file_id or not mime_type:
            return jsonify({"error": "Missing required fields: file_id, mime_type"}), 400

        content = ""

        # MODUS A: Push (Sicher für Verschlüsselung)
        if direct_content:
            logger.info("Processing direct content push for %s", file_id)
            content = direct_content

        # MODUS B: Pull (Fallback für unverschlüsselte Dateien auf Disk)
        elif file_path:
            logger.info("Processing file from disk: %s", file_path)
            if not os.path.exists(file_path):
                return jsonify({"error": f"File not found: {file_path}"}), 404
            
            with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
                content = f.read()
        else:
            return jsonify({"error": "No content and no file_path provided"}), 400

        if not content.strip():
            return jsonify({"status": "skipped", "reason": "empty content"}), 200

        # Ab hier bleibt alles gleich (Embedding erzeugen)
        logger.info("Generating Ollama embedding for %s (%d chars)", file_id, len(content))
        embedding = get_ollama_embedding(content[:8000]) 
        
        metadata = {
            "file_path": file_path or "memory", # Merken, woher es kam
            "mime_type": mime_type,
            "content_length": len(content)
        }

        with get_db_connection() as conn:
            register_vector(conn)
            with conn.cursor() as cur:
                cur.execute("""
                    INSERT INTO file_embeddings (file_id, chunk_index, content, embedding, metadata)
                    VALUES (%s, %s, %s, %s, %s)
                    ON CONFLICT (file_id, chunk_index) DO UPDATE SET
                        content = EXCLUDED.content,
                        embedding = EXCLUDED.embedding,
                        metadata = EXCLUDED.metadata,
                        created_at = CURRENT_TIMESTAMP
                """, (file_id, 0, content, embedding, json.dumps(metadata)))
                conn.commit()

        logger.info("File indexed successfully: %s", file_id)
        return jsonify({
            "status": "success",
            "file_id": file_id,
            "mode": "push" if direct_content else "disk"
        })

    except Exception as e:
        logger.error("Error processing file: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500




@app.route("/embed_query", methods=["POST"])
def embed_query():
    """Generate embedding for a search query via Ollama."""
    if not model_loaded:
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json()
        text = data.get("text", "")
        
        if not text:
            return jsonify({"error": "Missing 'text' field"}), 400

        embedding = get_ollama_embedding(text)
        
        return jsonify({
            "embedding": embedding,
            "dimension": len(embedding)
        })

    except Exception as e:
        logger.error("Embed query error: %s", str(e))
        return jsonify({"error": str(e)}), 500


@app.route("/search", methods=["POST"])
def vector_search():
    """Semantic vector search."""
    if not model_loaded:
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json()
        query = data.get("query", "")
        limit = data.get("limit", 10)

        if not query:
            return jsonify({"error": "Missing 'query' field"}), 400

        # Get query embedding
        query_embedding = get_ollama_embedding(query)

        with get_db_connection() as conn:
            register_vector(conn)
            with conn.cursor() as cur:
                cur.execute("""
                    SELECT file_id, content, metadata,
                           1 - (embedding <=> %s::vector) as similarity
                    FROM file_embeddings
                    ORDER BY embedding <=> %s::vector
                    LIMIT %s
                """, (query_embedding, query_embedding, limit))

                results = []
                for row in cur.fetchall():
                    metadata = row[2] or {}
                    results.append({
                        "file_id": row[0],
                        "file_path": metadata.get("file_path", "unknown"),
                        "content": row[1][:500],
                        "similarity": float(row[3])
                    })

        return jsonify({"results": results, "query": query})

    except Exception as e:
        logger.error("Search error: %s", str(e))
        return jsonify({"error": str(e)}), 500


@app.route("/rag", methods=["POST"])
def rag_query():
    """
    RAG (Retrieval Augmented Generation) endpoint.
    1. Find relevant documents via vector search
    2. Build context from documents
    3. Generate answer with Llama 3.2
    """
    if not model_loaded:
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json()
        query = data.get("query", "")
        top_k = data.get("top_k", 5)

        if not query:
            return jsonify({"error": "Missing 'query' field"}), 400

        logger.info("RAG query: %s", query[:100])

        # Step 1: Get query embedding
        query_embedding = get_ollama_embedding(query)

        # Step 2: Find relevant documents
        with get_db_connection() as conn:
            register_vector(conn)
            with conn.cursor() as cur:
                cur.execute("""
                    SELECT file_id, content, metadata,
                           1 - (embedding <=> %s::vector) as similarity
                    FROM file_embeddings
                    ORDER BY embedding <=> %s::vector
                    LIMIT %s
                """, (query_embedding, query_embedding, top_k))

                docs = []
                for row in cur.fetchall():
                    metadata = row[2] or {}
                    docs.append({
                        "file_id": row[0],
                        "file_path": metadata.get("file_path", "unknown"),
                        "content": row[1],
                        "similarity": float(row[3])
                    })

        if not docs:
            return jsonify({
                "answer": "Keine relevanten Dokumente gefunden.",
                "sources": [],
                "query": query
            })

        # Step 4: Generate intelligent answer with Llama
        logger.info("Generating intelligent RAG response with %d candidate documents", len(docs))
        llama_result = get_llama_response(query, docs)

        return jsonify({
            "answer": llama_result["answer"],
            "cited_sources": llama_result["cited_sources"],
            "confidence": llama_result["confidence"],
            "all_candidates": len(docs),
            "query": query,
            "model": LLM_MODEL
        })

    except Exception as e:
        logger.error("RAG error: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500


@app.route("/query", methods=["POST"])
def unified_query():
    """
    Unified Query Endpoint - AI decides routing and limit.
    
    Flow:
    1. Intent Classification (Llama or heuristics)
    2. Dynamic Search with variable limit
    3. Return either Search-Results OR RAG-Answer
    
    Request:
        {"query": "user input string"}
    
    Response for search mode:
        {
            "mode": "search",
            "intent": {...},
            "files": [...],
            "query": "original query"
        }
    
    Response for question mode:
        {
            "mode": "answer",
            "intent": {...},
            "answer": "AI generated answer",
            "sources": [...],
            "confidence": "HOCH|MITTEL|NIEDRIG",
            "query": "original query"
        }
    """
    if not model_loaded:
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json()
        query = data.get("query", "").strip()

        if not query:
            return jsonify({"error": "Missing 'query' field"}), 400

        logger.info("Unified query: %s", query[:100])

        # Step 1: Intent Classification
        intent = classify_intent(query)
        logger.info("Intent classified: type=%s, count_hint=%s, limit=%d", 
                   intent["type"], intent["count_hint"], intent["limit"])

        # Use refined query if available, otherwise original
        search_query = intent.get("refined_query") or query
        limit = intent["limit"]

        # Step 2: Get query embedding
        query_embedding = get_ollama_embedding(search_query)

        # Step 3: Vector search with dynamic limit
        with get_db_connection() as conn:
            register_vector(conn)
            with conn.cursor() as cur:
                cur.execute("""
                    SELECT file_id, content, metadata,
                           1 - (embedding <=> %s::vector) as similarity
                    FROM file_embeddings
                    WHERE (metadata->>'file_path') LIKE '/mnt/data/%%'
                      AND (metadata->>'file_path') NOT LIKE '%%/.trash/%%'
                    ORDER BY embedding <=> %s::vector
                    LIMIT %s
                """, (query_embedding, query_embedding, limit))

                docs = []
                for row in cur.fetchall():
                    metadata = row[2] or {}
                    docs.append({
                        "file_id": row[0],
                        "file_path": metadata.get("file_path", "unknown"),
                        "content": row[1],
                        "similarity": float(row[3])
                    })

        # Step 4: Route based on intent type
        if intent["type"] == "question":
            # RAG mode: Generate AI answer
            if not docs:
                return jsonify({
                    "mode": "answer",
                    "intent": intent,
                    "answer": "Keine relevanten Dokumente gefunden.",
                    "sources": [],
                    "confidence": "NIEDRIG",
                    "query": query
                })

            logger.info("Generating RAG response with %d documents", len(docs))
            llama_result = get_llama_response(query, docs)

            return jsonify({
                "mode": "answer",
                "intent": intent,
                "answer": llama_result["answer"],
                "sources": llama_result["cited_sources"],
                "confidence": llama_result["confidence"],
                "all_candidates": len(docs),
                "query": query,
                "model": LLM_MODEL
            })

        else:
            # Search mode: Return file results
            # Truncate content for response
            files = []
            for doc in docs:
                files.append({
                    "file_id": doc["file_id"],
                    "file_path": doc["file_path"],
                    "content": doc["content"][:500],
                    "similarity": doc["similarity"]
                })

            return jsonify({
                "mode": "search",
                "intent": intent,
                "files": files,
                "total_found": len(files),
                "query": query
            })

    except Exception as e:
        logger.error("Unified query error: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500


@app.route("/delete", methods=["POST"])
def delete_embeddings():
    """
    Delete embeddings for a specific file from the database.
    Prevents ghost knowledge by removing vector data when files are deleted.
    """
    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "Missing JSON payload"}), 400

        file_id = data.get("file_id")
        file_path = data.get("file_path")

        # Support deletion by either file_id or file_path
        if not file_id and not file_path:
            return jsonify({"error": "Either 'file_id' or 'file_path' is required"}), 400

        logger.info("Deleting embeddings for file_id=%s, file_path=%s", file_id, file_path)

        with get_db_connection() as conn:
            with conn.cursor() as cur:
                if file_id:
                    cur.execute("DELETE FROM file_embeddings WHERE file_id = %s", (file_id,))
                else:
                    cur.execute("DELETE FROM file_embeddings WHERE file_path = %s", (file_path,))

                deleted_count = cur.rowcount
                conn.commit()

        if deleted_count > 0:
            logger.info("Deleted %d embedding(s) for file_id=%s, file_path=%s", deleted_count, file_id, file_path)
            return jsonify({
                "status": "success",
                "deleted_count": deleted_count,
                "file_id": file_id,
                "file_path": file_path
            })
        else:
            logger.warning("No embeddings found for file_id=%s, file_path=%s", file_id, file_path)
            return jsonify({
                "status": "not_found",
                "deleted_count": 0,
                "message": "No embeddings found for this file"
            }), 404

    except Exception as e:
        logger.error("Delete error: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500


@app.route("/list_vectors", methods=["GET"])
def list_vectors():
    """
    List all file IDs in the vector database.
    Used for garbage collection / reconciliation to find orphaned embeddings.
    """
    try:
        with get_db_connection() as conn:
            with conn.cursor() as cur:
                cur.execute("SELECT DISTINCT file_id FROM file_embeddings")
                file_ids = [row[0] for row in cur.fetchall()]
        
        return jsonify({
            "file_ids": file_ids,
            "count": len(file_ids)
        })
    
    except Exception as e:
        logger.error("List vectors error: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500

def main():
    prewarm_models()
    init_db_pool()
    # SECURITY: Auto-index disabled for Zero-Trust (push-only via /process)
    if AUTO_INDEX_ENABLED:
        logger.info("⚠️ Auto-indexing enabled by env var - not recommended for production")
        index_all_files()
    else:
        logger.info("🔒 Auto-indexing disabled (Zero-Trust mode) - use /process endpoint")
    start_background_health_check()
    if WORKER_ENABLED:
        start_redis_worker()
    app.run(host="0.0.0.0", port=5000, debug=False)


# === REDIS JOB WORKER ===
def init_redis():
    """Initialize Redis connection for job queue."""
    global redis_client
    try:
        redis_client = redis_lib.from_url(REDIS_URL, decode_responses=True)
        redis_client.ping()
        logger.info("✅ Redis connection established for job queue")
        return True
    except Exception as e:
        logger.warning("Redis connection failed (worker disabled): %s", e)
        redis_client = None
        return False


def ensure_consumer_group():
    """Create consumer group if it doesn't exist."""
    if not redis_client:
        return
    try:
        redis_client.xgroup_create(JOB_STREAM, CONSUMER_GROUP, id="0", mkstream=True)
        logger.info("Created consumer group: %s", CONSUMER_GROUP)
    except redis_lib.ResponseError as e:
        if "BUSYGROUP" not in str(e):
            logger.error("Failed to create consumer group: %s", e)


def process_job_internal(query: str) -> dict:
    """
    Process a single AI query job.
    Returns the result dict to be stored in Redis.
    """
    if not model_loaded:
        return {"status": "failed", "error": "Ollama not loaded"}
    
    try:
        # Intent Classification
        intent = classify_intent(query)
        search_query = intent.get("refined_query") or query
        limit = intent["limit"]
        
        # Get query embedding
        query_embedding = get_ollama_embedding(search_query)
        
        # Vector search
        with get_db_connection() as conn:
            register_vector(conn)
            with conn.cursor() as cur:
                cur.execute("""
                    SELECT file_id, content, metadata,
                           1 - (embedding <=> %s::vector) as similarity
                    FROM file_embeddings
                    WHERE (metadata->>'file_path') LIKE '/mnt/data/%%'
                      AND (metadata->>'file_path') NOT LIKE '%%/.trash/%%'
                    ORDER BY embedding <=> %s::vector
                    LIMIT %s
                """, (query_embedding, query_embedding, limit))
                
                docs = []
                for row in cur.fetchall():
                    metadata = row[2] or {}
                    docs.append({
                        "file_id": row[0],
                        "file_path": metadata.get("file_path", "unknown"),
                        "content": row[1],
                        "similarity": float(row[3])
                    })
        
        # Route based on intent type
        if intent["type"] == "question":
            if not docs:
                return {
                    "status": "completed",
                    "mode": "answer",
                    "intent": intent,
                    "answer": "Keine relevanten Dokumente gefunden.",
                    "sources": [],
                    "confidence": "NIEDRIG",
                    "query": query
                }
            
            llama_result = get_llama_response(query, docs)
            return {
                "status": "completed",
                "mode": "answer",
                "intent": intent,
                "answer": llama_result["answer"],
                "sources": llama_result["cited_sources"],
                "confidence": llama_result["confidence"],
                "all_candidates": len(docs),
                "query": query,
                "model": LLM_MODEL
            }
        else:
            # Search mode
            files = []
            for doc in docs:
                files.append({
                    "file_id": doc["file_id"],
                    "file_path": doc["file_path"],
                    "content": doc["content"][:500],
                    "similarity": doc["similarity"]
                })
            
            return {
                "status": "completed",
                "mode": "search",
                "intent": intent,
                "files": files,
                "total_found": len(files),
                "query": query
            }
            
    except Exception as e:
        logger.error("Job processing error: %s", str(e), exc_info=True)
        return {"status": "failed", "error": str(e)}


def worker_loop():
    """Redis Stream consumer loop - processes AI jobs from queue."""
    if not redis_client:
        logger.warning("Redis not available, worker disabled")
        return
    
    consumer_name = f"worker-{os.getpid()}"
    logger.info("🚀 Starting Redis worker: %s", consumer_name)
    
    while True:
        try:
            # Read from stream with blocking (5 second timeout)
            messages = redis_client.xreadgroup(
                CONSUMER_GROUP, consumer_name,
                {JOB_STREAM: ">"},
                count=1, block=5000
            )
            
            if not messages:
                continue
            
            for stream_name, entries in messages:
                for entry_id, data in entries:
                    job_id = data.get("job_id")
                    query = data.get("query", "")
                    
                    if not job_id:
                        logger.warning("Job without ID, skipping")
                        redis_client.xack(JOB_STREAM, CONSUMER_GROUP, entry_id)
                        continue
                    
                    logger.info("Processing job %s: %s", job_id, query[:50])
                    
                    # Update status to processing
                    redis_client.setex(
                        f"{RESULT_KEY_PREFIX}{job_id}",
                        3600,
                        json.dumps({
                            "job_id": job_id,
                            "status": "processing",
                            "query": query
                        })
                    )
                    
                    # Process the job
                    try:
                        result = process_job_internal(query)
                        result["job_id"] = job_id
                        result["completed_at"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
                    except Exception as e:
                        logger.error("Job %s failed: %s", job_id, e)
                        result = {
                            "job_id": job_id,
                            "status": "failed",
                            "error": str(e)
                        }
                    
                    # Store result
                    redis_client.setex(
                        f"{RESULT_KEY_PREFIX}{job_id}",
                        3600,  # 1 hour TTL
                        json.dumps(result)
                    )
                    
                    # Acknowledge message
                    redis_client.xack(JOB_STREAM, CONSUMER_GROUP, entry_id)
                    logger.info("✅ Job %s completed: %s", job_id, result.get("status"))
                    
        except Exception as e:
            logger.error("Worker error: %s", e, exc_info=True)
            time.sleep(1)


def start_redis_worker():
    """Start the Redis worker thread."""
    if not init_redis():
        logger.warning("Redis worker not started (connection failed)")
        return
    
    ensure_consumer_group()
    
    thread = threading.Thread(target=worker_loop, daemon=True)
    thread.start()
    logger.info("🎯 Redis worker thread started")


# === GUNICORN COMPATIBILITY ===
prewarm_models()
init_db_pool()
# SECURITY: Auto-index disabled for Zero-Trust (push-only via /process)
if AUTO_INDEX_ENABLED:
    logger.info("⚠️ Auto-indexing enabled by env var - not recommended for production")
    index_all_files()
else:
    logger.info("🔒 Auto-indexing disabled (Zero-Trust mode) - use /process endpoint")
start_background_health_check()
if WORKER_ENABLED:
    start_redis_worker()


if __name__ == "__main__":
    main()

