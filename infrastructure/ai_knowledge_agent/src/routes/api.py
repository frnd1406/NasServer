import os
import logging
from flask import Blueprint, jsonify, request, current_app

logger = logging.getLogger("ai_knowledge_agent")

api_bp = Blueprint('api', __name__)

# SECURITY: Shared Secret for Internal Communication
INTERNAL_API_SECRET = os.getenv("INTERNAL_API_SECRET")

@api_bp.before_request
def check_internal_auth():
    """Global middleware to require X-Internal-Secret header."""
    if request.endpoint == 'api.health': # Allow health check
        return None
    if request.method == 'OPTIONS':
        return None

    secret = request.headers.get('X-Internal-Secret')
    if not INTERNAL_API_SECRET:
        logger.error("INTERNAL_API_SECRET not set! Rejecting request.")
        return jsonify({"error": "Server misconfiguration"}), 500
            
    if not secret or secret != INTERNAL_API_SECRET:
        logger.warning("Unauthorized access attempt from %s", request.remote_addr)
        return jsonify({"error": "Unauthorized"}), 403

def get_service():
    """Helper to get RAG Service from app config."""
    return current_app.config['RAG_SERVICE']

@api_bp.route("/health", methods=["GET"])
def health():
    svc = get_service()
    db_ok = svc.db.db_pool is not None
    ollama_ok = svc.check_ollama_health()
    
    status = "ok" if (svc.model_loaded and db_ok) else "degraded"
    return jsonify({
        "status": status,
        "ollama_ready": ollama_ok,
        "db_ready": db_ok,
        "message": "System operational" if status == "ok" else "System degraded"
    }), 200 if status == "ok" else 503

@api_bp.route("/status", methods=["GET"])
def status():
    svc = get_service()
    
    # Get models (simplified check)
    ollama_models = []
    try:
        if svc.check_ollama_health():
            ollama_models = [svc.embedding_model, svc.llm_model]
    except Exception: pass

    return jsonify({
        "ollama": {"connected": svc.model_loaded, "models": ollama_models},
        "models": {"embedding": svc.embedding_model, "llm": svc.llm_model},
        "index": {"indexed_files": svc.db.get_index_stats()},
        "settings": {"auto_index": svc.auto_index_enabled}
    })

@api_bp.route("/process", methods=["POST"])
def process_file():
    svc = get_service()
    if not svc.model_loaded:
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json() or {}
        file_id = data.get("file_id")
        mime_type = data.get("mime_type")
        direct_content = data.get("content") 
        file_path = data.get("file_path")

        if not file_id or not mime_type:
            return jsonify({"error": "Missing required fields"}), 400

        content = ""
        if direct_content:
            content = direct_content
        elif file_path:
            if not os.path.exists(file_path):
                return jsonify({"error": "File not found"}), 404
            with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
                content = f.read()
        else:
            return jsonify({"error": "No content provided"}), 400

        if not content.strip():
            return jsonify({"status": "skipped", "reason": "empty content"}), 200

        svc.process_file(content, mime_type, file_id, file_path)
        
        return jsonify({"status": "success", "file_id": file_id})

    except Exception as e:
        logger.error("Error processing file: %s", e)
        return jsonify({"error": str(e)}), 500

@api_bp.route("/query", methods=["POST"])
def unified_query():
    svc = get_service()
    if not svc.model_loaded:
        return jsonify({"error": "Ollama not loaded"}), 503

    try:
        data = request.get_json() or {}
        query = data.get("query", "").strip()
        if not query:
            return jsonify({"error": "Missing query"}), 400

        result = svc.unified_query(query)
        return jsonify(result)
        
    except Exception as e:
        logger.error("Query error: %s", e)
        return jsonify({"error": str(e)}), 500

@api_bp.route("/delete", methods=["POST"])
def delete_embeddings():
    svc = get_service()
    try:
        data = request.get_json() or {}
        file_id = data.get("file_id")
        file_path = data.get("file_path")
        
        deleted = svc.db.delete_embeddings(file_id, file_path)
        
        if deleted > 0:
            return jsonify({"status": "success", "deleted_count": deleted})
        else:
            return jsonify({"status": "not_found", "message": "No embeddings found"}), 404

    except Exception as e:
        logger.error("Delete error: %s", e)
        return jsonify({"error": str(e)}), 500

# Legacy Endpoints for compatibility
@api_bp.route("/rag", methods=["POST"])
def rag_query():
    # Maps to unified query for now, or could implement strict RAG
    # Reusing unified query logic as it handles RAG
    return unified_query()

@api_bp.route("/list_vectors", methods=["GET"])
def list_vectors():
    svc = get_service()
    try:
        ids = svc.db.list_all_vectors()
        return jsonify({"file_ids": ids, "count": len(ids)})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
