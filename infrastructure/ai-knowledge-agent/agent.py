"""
=================================================================
AI Knowledge Agent - Main Application
=================================================================
Purpose: Load ML models, provide health checks, await commands
Phase: 2.2 - AI Core Infrastructure
=================================================================
"""

import os
import sys
import time
from loguru import logger
from flask import Flask, jsonify
from sentence_transformers import SentenceTransformer
from db_connection import DatabaseConnection
import threading

# Configure logging
logger.remove()
logger.add(sys.stderr, format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {message}")

# Initialize Flask app for health checks
app = Flask(__name__)

# Global state
model = None
db_connection = None
model_loaded = False
db_connected = False
model_lock = threading.Lock()


def load_model():
    """
    Load the sentence transformer model into memory.
    Uses all-MiniLM-L6-v2: fast, efficient, good quality embeddings.
    """
    global model, model_loaded

    try:
        model_name = os.getenv("MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2")
        logger.info(f"Loading model: {model_name}")
        logger.info("This may take a few minutes on first run (downloading model)...")

        temp_model = SentenceTransformer(model_name)

        # Test model with a sample encoding
        with model_lock:
            model = temp_model
            test_embedding = model.encode("Hello, world!")
            embedding_dim = len(test_embedding)

        logger.info(f"✅ Model loaded successfully!")
        logger.info(f"   Model: {model_name}")
        logger.info(f"   Embedding dimension: {embedding_dim}")
        logger.info(f"   Max sequence length: {model.max_seq_length}")

        model_loaded = True
        return True

    except Exception as e:
        logger.error(f"❌ Failed to load model: {str(e)}")
        model_loaded = False
        return False


def connect_database():
    """
    Connect to PostgreSQL database with retry logic.
    """
    global db_connection, db_connected

    try:
        logger.info("Connecting to PostgreSQL database...")
        db_connection = DatabaseConnection(max_retries=10, retry_delay=5)
        db_connection.connect()
        db_connected = True
        logger.info("✅ Database connection established")
        return True

    except Exception as e:
        logger.error(f"❌ Failed to connect to database: {str(e)}")
        db_connected = False
        return False


@app.route('/health', methods=['GET'])
def health_check():
    """
    Health check endpoint for Docker healthcheck.
    Returns 200 if model loaded and DB connected, 503 otherwise.
    """
    embedding_dim = None
    if model_loaded and model is not None:
        with model_lock:
            if model is not None: # Double check after lock
                embedding_dim = len(model.encode("test"))

    status = {
        "status": "healthy" if (model_loaded and db_connected) else "unhealthy",
        "model_loaded": model_loaded,
        "database_connected": db_connected,
        "model_name": os.getenv("MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2"),
        "embedding_dimension": embedding_dim
    }

    status_code = 200 if (model_loaded and db_connected) else 503
    return jsonify(status), status_code


@app.route('/status', methods=['GET'])
def status():
    """
    Detailed status endpoint.
    """
    embedding_dim = None
    if model_loaded and model is not None:
        with model_lock:
             if model is not None: # Double check after lock
                embedding_dim = len(model.encode("test"))

    return jsonify({
        "service": "AI Knowledge Agent",
        "version": "1.0.0",
        "phase": "2.2",
        "model_loaded": model_loaded,
        "database_connected": db_connected,
        "model_name": os.getenv("MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2") if model_loaded else None,
        "embedding_dimension": embedding_dim,
        "ready": model_loaded and db_connected
    })


@app.route('/embed', methods=['POST'])
def embed_text():
    """
    Generate embeddings for text (future endpoint).
    Currently returns placeholder.
    """
    if not model_loaded or model is None:
        return jsonify({"error": "Model not loaded"}), 503

    # Example of locking for a real embedding call
    # with model_lock:
    #     data = request.get_json()
    #     text = data.get("text")
    #     embedding = model.encode(text)

    return jsonify({
        "message": "Embedding endpoint - coming soon in Phase 2.3",
        "status": "ready"
    }), 200


def run_flask():
    """Run Flask server in a separate thread"""
    app.run(host='0.0.0.0', port=8000, debug=False)


def main():
    """
    Main entry point.
    1. Load ML model
    2. Connect to database
    3. Start health check server
    4. Keep running (await commands)
    """
    logger.info("=" * 60)
    logger.info("AI Knowledge Agent - Starting")
    logger.info("=" * 60)

    # Step 1: Load model
    logger.info("Step 1/3: Loading ML model...")
    if not load_model():
        logger.error("Failed to load model - exiting")
        sys.exit(1)

    # Step 2: Connect to database
    logger.info("Step 2/3: Connecting to database...")
    if not connect_database():
        logger.error("Failed to connect to database - exiting")
        sys.exit(1)

    # Step 3: Start health check server
    logger.info("Step 3/3: Starting health check server...")
    flask_thread = threading.Thread(target=run_flask, daemon=True)
    flask_thread.start()
    logger.info("✅ Health check server started on port 8000")

    logger.info("=" * 60)
    logger.info("🚀 AI Knowledge Agent is READY")
    logger.info("=" * 60)
    logger.info("Status:")
    logger.info(f"  ✅ Model loaded: {model_loaded}")
    logger.info(f"  ✅ Database connected: {db_connected}")
    logger.info(f"  ✅ Health check: http://localhost:8000/health")
    logger.info("Awaiting commands...")
    logger.info("=" * 60)

    # Keep the main thread alive
    try:
        while True:
            time.sleep(60)
            logger.debug("Agent alive - awaiting commands...")
    except KeyboardInterrupt:
        logger.info("Shutting down gracefully...")
        if db_connection:
            db_connection.close()
        logger.info("Goodbye!")


if __name__ == "__main__":
    main()
