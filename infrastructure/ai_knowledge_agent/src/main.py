import logging
import os
from flask import Flask

from src.database import Database
from src.services.rag_service import RAGService
from src.routes.api import api_bp

# Configure Logging
logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("ai_knowledge_agent")

def create_app():
    """Application Factory & Composition Root"""
    app = Flask(__name__)
    
    # 1. Initialize Infrastructure
    db = Database()
    
    # 2. Initialize Application Services (Dependency Injection)
    rag_service = RAGService(db)
    
    # 3. Store Service in App Context (for Routes)
    app.config['RAG_SERVICE'] = rag_service
    
    # 4. Register Interface (Routes)
    app.register_blueprint(api_bp)
    
    # 5. Startup Logic
    with app.app_context():
        success = db.init_pool()
        if not success:
            logger.error("Failed to initialize database pool on startup")
        
        # Start background tasks
        rag_service.start_background_threads()
        rag_service.prewarm_models()
        
    return app

if __name__ == "__main__":
    app = create_app()
    port = int(os.getenv("PORT", 5000))
    app.run(host="0.0.0.0", port=port, debug=False)
