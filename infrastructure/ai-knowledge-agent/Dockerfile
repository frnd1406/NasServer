# =================================================================
# AI Knowledge Agent - Semantic Search & Embedding Service
# =================================================================
# Purpose: Load ML models, generate embeddings, perform semantic search
# Phase: 2.2 - AI Core Infrastructure
# =================================================================

FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy requirements first for better caching
COPY requirements.txt .

# Install Python dependencies
# IMPORTANT: Use CPU-optimized PyTorch to save space
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create model cache directory
RUN mkdir -p /app/models

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV TRANSFORMERS_CACHE=/app/models
ENV HF_HOME=/app/models

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD curl -f http://localhost:8000/health || exit 1

# Expose port for health checks and future API
EXPOSE 8000

# Run the agent
ENV PYTHONPATH=/app
CMD ["python", "agent.py"]
