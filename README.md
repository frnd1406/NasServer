# NAS AI Server System

[![Go Report Card](https://goreportcard.com/badge/github.com/frnd1406/NasServer)](https://goreportcard.com/report/github.com/frnd1406/NasServer)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

**Secure, Self-Hosted NAS with Integrated AI Knowledge Base.**

This system combines robust file storage with a local AI RAG (Retrieval-Augmented Generation) agent, offering intelligent search, content analysis, and secure hybrid encryption for your private data.

---

## 🏗 Architecture Overview

The system follows a modular microservices architecture:

*   **Frontend**: [React 18](https://reactjs.org/) + [Vite](https://vitejs.dev/) + [TailwindCSS](https://tailwindcss.com/)
    *   Responsive "Glassmorphism" UI.
    *   Client-side encryption helpers.
*   **Backend**: [Go (Golang)](https://go.dev/)
    *   High-performance REST API.
    *   "Thin Handlers, Fat Services" design pattern.
*   **Database**:
    *   **PostgreSQL**: Metadata, Users, Vector Embeddings (pgvector).
    *   **Redis**: Caching, Background Job Queues.
*   **AI Engine**: Python-based RAG Agent using [Ollama](https://ollama.ai/).
    *   Locally hosted LLMs (Llama 3, Mistral, etc.).

---

## 🚀 Quickstart

Get your instance running in minutes using Docker.

```bash
# 1. Clone the repository
git clone https://github.com/frnd1406/NasServer.git
cd NasServer

# 2. Navigate to infrastructure
cd infrastructure

# 3. Start the stack
docker-compose -f docker-compose.dev.yml up -d --build
```

Access the Web UI at `http://localhost:5173` (default).

---

## 📚 Documentation

*   [**Backend API**](./infrastructure/api/README.md): API architecture, setup, and development.
*   [**Web UI**](./infrastructure/webui/README.md): Frontend structure and components.
*   [**Orchestrator**](./orchestrator/README.md): Service health and management.
*   [**Developer Guide**](./docs/development/DEV_GUIDE.md): Detailed contribution guidelines.
