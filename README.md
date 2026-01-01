# NAS AI Server System

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://reactjs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Secure, Self-Hosted NAS with Integrated AI Knowledge Base.**

Ein sicheres, selbst-gehostetes NAS-System mit integrierter RAG-KI, Hybrid-Verschlüsselung und moderner Web-Oberfläche.

---

## 🏗 Architecture

```mermaid
graph TD
    Client[WebUI / Mobile] --> Gateway[NGINX Reverse Proxy]
    Gateway --> API[Go API Core]
    API --> Auth[Auth Middleware]
    API --> Svc[Services Layer]
    Svc --> DB[(PostgreSQL)]
    Svc --> Cache[(Redis)]
    Svc --> FS[File System / Storage]
    Svc --> AI[AI Knowledge Agent]
    Orch[Orchestrator] -.-> API
    Orch -.-> AI
```

---

## ✨ Key Features

*   **🧠 AI Core**: Lokale Wissensdatenbank (RAG) mit Ollama-Integration für intelligente Dokumentenanalyse.
*   **🔒 Security**: ChaCha20-Poly1305 Hybrid-Verschlüsselung, Zip-Slip Protection und strikte Validierung.
*   **🚀 Performance**: Go-Backend mit asynchronen Job-Queues (Redis) für schnelle Antwortzeiten.

---

## 🚀 Quick Start

Starten Sie das gesamte System mit Docker Compose:

```bash
git clone https://github.com/frnd1406/NasServer.git
cd NasServer/infrastructure
docker-compose up -d --build
```

---

## 📚 Documentation

*   [**Backend API**](./infrastructure/api/README.md)
*   [**Web Dashboard**](./infrastructure/webui/README.md)
*   [**Orchestrator**](./orchestrator/README.md)
*   [**Developer Guide**](./docs/development/DEV_GUIDE.md)
