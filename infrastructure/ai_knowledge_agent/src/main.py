import logging
import os
import time
import json
from contextlib import contextmanager

import requests
from flask import Flask, jsonify, request
import psycopg2
from psycopg2 import pool, OperationalError
from pgvector.psycopg2 import register_vector
import numpy as np


logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("ai_knowledge_agent")

app = Flask(__name__)

# Ollama Configuration
OLLAMA_URL = os.getenv("OLLAMA_URL", "http://host.docker.internal:11434")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "mxbai-embed-large")
LLM_MODEL = os.getenv("LLM_MODEL", "llama3.2")
EMBEDDING_DIM = 1024  # mxbai-embed-large dimension

model_loaded = False
db_pool = None

# FIX [BUG-PY-006]: Module-level constant for DSN
DSN = {
    "host": os.getenv("PGHOST", "postgres"),
    "port": os.getenv("PGPORT", "5432"),
    "dbname": os.getenv("PGDATABASE", "nas_db"),
    "user": os.getenv("PGUSER", "nas_user"),
    "password": os.getenv("PGPASSWORD", ""),
}


def get_ollama_embedding(text: str) -> list:
    """Get embedding from Ollama mxbai-embed-large model."""
    try:
        response = requests.post(
            f"{OLLAMA_URL}/api/embeddings",
            json={"model": EMBEDDING_MODEL, "prompt": text},
            timeout=30
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
            timeout=120
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
        response = requests.get(f"{OLLAMA_URL}/api/tags", timeout=5)
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
                    cur.execute(f"""
                        CREATE TABLE IF NOT EXISTS file_embeddings (
                            id SERIAL PRIMARY KEY,
                            file_id TEXT NOT NULL,
                            file_path TEXT NOT NULL,
                            mime_type TEXT NOT NULL,
                            content TEXT NOT NULL,
                            embedding vector({EMBEDDING_DIM}),
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE(file_id)
                        )
                    """)
                    conn.commit()
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
    if not db_pool:
        raise Exception("Database pool not initialized")
    
    conn = db_pool.getconn()
    try:
        yield conn
    finally:
        db_pool.putconn(conn)


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


@app.route("/process", methods=["POST"])
def process_file():
    """Process an uploaded file and generate embeddings via Ollama."""
    if not model_loaded:
        logger.error("Ollama not ready")
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "Missing JSON payload"}), 400

        file_path = data.get("file_path")
        file_id = data.get("file_id")
        mime_type = data.get("mime_type")

        if not all([file_path, file_id, mime_type]):
            return jsonify({"error": "Missing required fields: file_path, file_id, mime_type"}), 400

        logger.info("Processing file: %s (ID: %s)", file_path, file_id)

        if not os.path.exists(file_path):
            return jsonify({"error": f"File not found: {file_path}"}), 404

        with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()

        if not content.strip():
            return jsonify({"status": "skipped", "reason": "empty file"}), 200

        # Get embedding from Ollama
        logger.info("Generating Ollama embedding for %s (%d chars)", file_id, len(content))
        embedding = get_ollama_embedding(content[:8000])  # Limit content size
        
        with get_db_connection() as conn:
            register_vector(conn)
            with conn.cursor() as cur:
                cur.execute("""
                    INSERT INTO file_embeddings (file_id, file_path, mime_type, content, embedding)
                    VALUES (%s, %s, %s, %s, %s)
                    ON CONFLICT (file_id) DO UPDATE SET
                        content = EXCLUDED.content,
                        embedding = EXCLUDED.embedding,
                        created_at = CURRENT_TIMESTAMP
                """, (file_id, file_path, mime_type, content, embedding))
                conn.commit()

        logger.info("File indexed: %s", file_id)
        return jsonify({
            "status": "success",
            "file_id": file_id,
            "content_length": len(content),
            "embedding_dim": len(embedding)
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
                    SELECT file_id, file_path, content, 
                           1 - (embedding <=> %s::vector) as similarity
                    FROM file_embeddings
                    ORDER BY embedding <=> %s::vector
                    LIMIT %s
                """, (query_embedding, query_embedding, limit))
                
                results = [{
                    "file_id": row[0],
                    "file_path": row[1],
                    "content": row[2][:500],
                    "similarity": float(row[3])
                } for row in cur.fetchall()]

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
                    SELECT file_id, file_path, content,
                           1 - (embedding <=> %s::vector) as similarity
                    FROM file_embeddings
                    ORDER BY embedding <=> %s::vector
                    LIMIT %s
                """, (query_embedding, query_embedding, top_k))
                
                docs = [{
                    "file_id": row[0],
                    "file_path": row[1],
                    "content": row[2],
                    "similarity": float(row[3])
                } for row in cur.fetchall()]

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


def main():
    check_ollama_health()
    init_db_pool()
    app.run(host="0.0.0.0", port=5000, debug=False)


# === GUNICORN COMPATIBILITY ===
check_ollama_health()
init_db_pool()


if __name__ == "__main__":
    main()
