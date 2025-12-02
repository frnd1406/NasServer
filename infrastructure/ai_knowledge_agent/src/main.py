import logging
import os
import time

from flask import Flask, jsonify, request
import psycopg2
from psycopg2 import OperationalError
from pgvector.psycopg2 import register_vector
from sentence_transformers import SentenceTransformer


logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("ai_knowledge_agent")

app = Flask(__name__)

MODEL_NAME = "sentence-transformers/all-MiniLM-L6-v2"
model = None
model_loaded = False


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


def db_connect_with_retry(retries=10, delay=2.0):
    dsn = {
        "host": os.getenv("PGHOST", "postgres"),
        "port": os.getenv("PGPORT", "5432"),
        "dbname": os.getenv("PGDATABASE", "nas_db"),
        "user": os.getenv("PGUSER", "nas_user"),
        "password": os.getenv("PGPASSWORD", ""),
    }
    for attempt in range(1, retries + 1):
        try:
            conn = psycopg2.connect(**dsn)
            conn.close()
            logger.info("Database connection OK (attempt %d).", attempt)
            return True
        except OperationalError as err:
            logger.warning("DB connection failed (attempt %d/%d): %s", attempt, retries, err)
            time.sleep(delay)
    return False


@app.route("/health", methods=["GET"])
def health():
    db_ok = db_connect_with_retry(retries=1, delay=0)  # quick check on health requests

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

    Expected JSON payload:
    {
        "file_path": "/mnt/data/document.txt",
        "file_id": "document.txt",
        "mime_type": "text/plain"
    }
    """
    # Null-pointer prevention: Check both model_loaded AND model is not None
    if not model_loaded or model is None:
        logger.error("Model not loaded yet or is None")
        return jsonify({"error": "Model not loaded"}), 503

    conn = None  # Initialize for try-finally safety
    try:
        data = request.get_json()
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

        # Store in database with proper resource management
        dsn = {
            "host": os.getenv("PGHOST", "postgres"),
            "port": os.getenv("PGPORT", "5432"),
            "dbname": os.getenv("PGDATABASE", "nas_db"),
            "user": os.getenv("PGUSER", "nas_user"),
            "password": os.getenv("PGPASSWORD", ""),
        }

        # Open connection AFTER embedding generation to minimize leak window
        conn = psycopg2.connect(**dsn)

        # CRITICAL FIX: Register pgvector type adapter
        register_vector(conn)

        try:
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
        finally:
            # CRITICAL: Always close connection even on error
            if conn:
                conn.close()
                logger.debug("DB connection closed for %s", file_id)

        return jsonify({
            "status": "success",
            "file_id": file_id,
            "content_length": len(content),
            "embedding_dim": len(embedding)
        }), 200

    except Exception as e:
        logger.error("Error processing file: %s", str(e), exc_info=True)
        # CRITICAL: Close connection in error path too
        if conn:
            try:
                conn.close()
                logger.debug("DB connection closed after error")
            except:
                pass
        return jsonify({"error": str(e)}), 500


@app.route("/embed_query", methods=["POST"])
def embed_query():
    """
    Generate an embedding for an arbitrary query text.

    Expected JSON payload:
    {
        "text": "Meine Suchanfrage"
    }
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
    db_connect_with_retry()
    app.run(host="0.0.0.0", port=5000)


if __name__ == "__main__":
    main()
