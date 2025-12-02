"""
=================================================================
Database Connection Module with Retry Logic
=================================================================
Purpose: Robust PostgreSQL connection with automatic reconnection
Phase: 2.2 - AI Core Infrastructure
=================================================================
"""

import os
import time
from loguru import logger
from sqlalchemy import create_engine, text
from sqlalchemy.pool import NullPool
import psycopg2
from pgvector.psycopg2 import register_vector


class DatabaseConnection:
    """
    Handles PostgreSQL connection with automatic retry logic.
    Ensures vector extension is available before proceeding.
    """

    def __init__(self, max_retries=10, retry_delay=5):
        self.database_url = os.getenv("DATABASE_URL",
            "postgresql://nas_user:nas_password@postgres:5432/nas_db")

        # FIX [BUG-PY-013]: Validate DB password is not default in production
        if "nas_password" in self.database_url or "nas_dev_password" in self.database_url:
            logger.warning("⚠️  WARNING: Using default database password! Change DATABASE_URL in production!")
            if os.getenv("ENVIRONMENT", "development") == "production":
                raise ValueError("Default database password not allowed in production environment")

        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self.engine = None
        self.connection = None

    def connect(self):
        """
        Connect to PostgreSQL with retry logic.
        Waits for database to be ready if it's still starting up.
        """
        # Create SQLAlchemy engine once
        if not self.engine:
            self.engine = create_engine(
                self.database_url,
                poolclass=NullPool,  # No connection pooling for simplicity
                echo=False
            )

        for attempt in range(1, self.max_retries + 1):
            try:
                logger.info(f"Attempt {attempt}/{self.max_retries}: Connecting to PostgreSQL...")

                # Test connection
                with self.engine.connect() as conn:
                    result = conn.execute(text("SELECT version();"))
                    version = result.fetchone()[0]
                    logger.info(f"✅ Connected to PostgreSQL: {version}")

                    # Verify pgvector extension is installed
                    result = conn.execute(text(
                        "SELECT extversion FROM pg_extension WHERE extname = 'vector';"
                    ))
                    vector_version = result.fetchone()

                    if vector_version:
                        logger.info(f"✅ pgvector extension found: v{vector_version[0]}")
                    else:
                        logger.error("❌ pgvector extension not installed!")
                        raise Exception("pgvector extension required but not found")

                # Create raw psycopg2 connection for pgvector
                self.connection = psycopg2.connect(self.database_url)
                register_vector(self.connection)
                logger.info("✅ pgvector registered for vector operations")

                return True

            except Exception as e:
                logger.warning(f"Connection attempt {attempt} failed: {str(e)}")

                if attempt < self.max_retries:
                    logger.info(f"Retrying in {self.retry_delay} seconds...")
                    time.sleep(self.retry_delay)
                else:
                    logger.error(f"❌ Failed to connect after {self.max_retries} attempts")
                    # Cleanup engine on final failure
                    if self.engine:
                        self.engine.dispose()
                    raise

        return False

    def get_engine(self):
        """Get SQLAlchemy engine for queries"""
        if not self.engine:
            self.connect()
        return self.engine

    def get_connection(self):
        """Get raw psycopg2 connection for vector operations"""
        if not self.connection:
            self.connect()
        return self.connection

    def close(self):
        """Close database connections"""
        if self.connection:
            self.connection.close()
            logger.info("Closed psycopg2 connection")

        if self.engine:
            self.engine.dispose()
            logger.info("Disposed SQLAlchemy engine")
