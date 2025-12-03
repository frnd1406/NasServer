import logging
import os
import time
from contextlib import contextmanager

from flask import Flask, jsonify, request
import psycopg2
from psycopg2 import pool, OperationalError
from pgvector.psycopg2 import register_vector
from sentence_transformers import SentenceTransformer


logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("ai_knowledge_agent")

app = Flask(__name__)

MODEL_NAME = "sentence-transformers/all-MiniLM-L6-v2"
model = None
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

def load_model():
    global model, model_loaded
    logger.info("Loading model %s ...", MODEL_NAME)
    try:
        model = SentenceTransformer(MODEL_NAME)
        model_loaded = True
        logger.info("Model loaded.")
    except Exception as err:
        model = None
        model_loaded = False
        logger.error("Model load failed: %s", err)


def init_db_pool():
    global db_pool
    try:
        # FIX [BUG-PY-008]: Use connection pool to prevent churn
        db_pool = psycopg2.pool.SimpleConnectionPool(
            minconn=1,
            maxconn=10,
            **DSN
        )
        logger.info("Database connection pool initialized.")
        return True
    except Exception as e:
        logger.error("Failed to initialize DB pool: %s", e)
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
    db_ok = False
    if db_pool:
        try:
            with get_db_connection() as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT 1")
                db_ok = True
        except Exception as e:
            logger.error("Health check DB failure: %s", e)
            db_ok = False
    
    # Return 503 if model is not ready (crash prevention)
    if not model_loaded or model is None:
        return jsonify(
            {
                "status": "loading",
                "model_loaded": False,
                "db_ok": db_ok,
                "message": "Model is still loading, service not ready"
            }
        ), 503

    return jsonify(
        {
            "status": "ok" if (model_loaded and db_ok) else "degraded",
            "model_loaded": model_loaded,
            "db_ok": db_ok,
        }
    ), 200


@app.route("/process", methods=["POST"])
def process_file():
    """
    Process an uploaded file and generate embeddings.
    """
    # Null-pointer prevention: Check both model_loaded AND model is not None
    if not model_loaded or model is None:
        logger.error("Model not loaded yet or is None")
        return jsonify({"error": "Model not loaded"}), 503

    try:
        # FIX [BUG-PY-014]: Add explicit JSON parse error logging
        try:
            data = request.get_json()
        except Exception as json_err:
            logger.error("JSON parse error: %s", json_err)
            return jsonify({"error": "Invalid JSON payload"}), 400

        if not data:
            return jsonify({"error": "Missing JSON payload"}), 400

        file_path = data.get("file_path")
        file_id = data.get("file_id")
        mime_type = data.get("mime_type")

        if not all([file_path, file_id, mime_type]):
            return jsonify({"error": "Missing required fields: file_path, file_id, mime_type"}), 400

        logger.info("Processing file: %s (ID: %s, MIME: %s)", file_path, file_id, mime_type)

        # Read file content
        if not os.path.exists(file_path):
            logger.error("File not found: %s", file_path)
            return jsonify({"error": f"File not found: {file_path}"}), 404

        with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
            content = f.read()

        if not content.strip():
            logger.warning("File is empty: %s", file_path)
            return jsonify({"status": "skipped", "reason": "empty file"}), 200

        # Generate embeddings (BEFORE opening DB connection to avoid leak)
        logger.info("Generating embeddings for %s (%d chars)", file_id, len(content))
        embedding = model.encode(content)

        # FIX [BUG-PY-002]: Use connection pool context manager to ensure cleanup
        with get_db_connection() as conn:
            # CRITICAL FIX: Register pgvector type adapter
            register_vector(conn)

            with conn.cursor() as cur:
                # Create table if it doesn't exist
                cur.execute("""
                    CREATE TABLE IF NOT EXISTS file_embeddings (
                        id SERIAL PRIMARY KEY,
                        file_id TEXT NOT NULL,
                        file_path TEXT NOT NULL,
                        mime_type TEXT NOT NULL,
                        content TEXT NOT NULL,
                        embedding vector(384),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        UNIQUE(file_id)
                    )
                """)

                # Insert or update embedding
                cur.execute("""
                    INSERT INTO file_embeddings (file_id, file_path, mime_type, content, embedding)
                    VALUES (%s, %s, %s, %s, %s)
                    ON CONFLICT (file_id)
                    DO UPDATE SET
                        file_path = EXCLUDED.file_path,
                        mime_type = EXCLUDED.mime_type,
                        content = EXCLUDED.content,
                        embedding = EXCLUDED.embedding,
                        created_at = CURRENT_TIMESTAMP
                """, (file_id, file_path, mime_type, content, embedding.tolist()))

                conn.commit()
                logger.info("Successfully stored embeddings for %s", file_id)

        return jsonify({
            "status": "success",
            "file_id": file_id,
            "content_length": len(content),
            "embedding_dim": len(embedding)
        }), 200

    except Exception as e:
        logger.error("Error processing file: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500


@app.route("/embed_query", methods=["POST"])
def embed_query():
    """
    Generate an embedding for an arbitrary query text.
    """
    # Null-pointer prevention: Check both model_loaded AND model is not None
    if not model_loaded or model is None:
        logger.error("Model not loaded yet or is None")
        return jsonify({"error": "Model not loaded"}), 503

    try:
        data = request.get_json(silent=True) or {}
        text = data.get("text", "")
        if not text or not str(text).strip():
            return jsonify({"error": "Missing or empty 'text'"}), 400

        logger.info("Generating embedding for query (%d chars)", len(text))
        embedding = model.encode(text)

        return jsonify({"embedding": embedding.tolist()}), 200
    except Exception as e:
        logger.error("Error generating query embedding: %s", str(e), exc_info=True)
        return jsonify({"error": str(e)}), 500


def main():
    load_model()
    init_db_pool()
    # FIX [BUG-PY-015]: Explicitly disable Flask debug mode for production safety
    app.run(host="0.0.0.0", port=5000, debug=False)


if __name__ == "__main__":
    main()
