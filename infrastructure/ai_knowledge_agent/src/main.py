import logging
import os
import time

from flask import Flask, jsonify
import psycopg2
from psycopg2 import OperationalError
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
    model = SentenceTransformer(MODEL_NAME)
    model_loaded = True
    logger.info("Model loaded.")


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
    return jsonify(
        {
            "status": "ok" if (model_loaded and db_ok) else "degraded",
            "model_loaded": model_loaded,
            "db_ok": db_ok,
        }
    ), 200


def main():
    load_model()
    db_connect_with_retry()
    app.run(host="0.0.0.0", port=5000)


if __name__ == "__main__":
    main()
