# NAS AI Backend API

The backend is built with **Go (Golang)** and serves as the core orchestrator for file operations, metadata management, and AI interaction.

## 🏛 Core Philosophy

**"Thin Handlers, Fat Services"**

*   **Handlers**: Responsible *only* for HTTP request parsing, validation, and invoking services. No business logic allowed here.
*   **Services**: Encapsulate all business logic, security checks, and data manipulation.

## 🔑 Key Services

*   **[ArchiveService](./src/services/archive_service.go)**
    *   Handles secure ZIP creation and extraction.
    *   **Security**: Prevents "Zip Slip" attacks, enforces quotas (size, file count, compression ratio).

*   **[ContentDeliveryService](./src/services/content_delivery_service.go)**
    *   Centralized logic for file streaming.
    *   **Features**: Smart range requests, transparent decryption streaming, Nginx X-Accel-Redirect support.

*   **[AIAgentService](./src/services/ai_agent_service.go)**
    *   Manages asynchronous communication with the Python AI Knowledge Agent.
    *   Handles upload notifications and index triggers.

*   **[EncryptionService](./src/services/encryption_service.go)**
    *   Provides stream-based encryption using **XChaCha20-Poly1305**.
    *   Manages Vault locking/unlocking and key derivation.

## 🛠 Development

### Prerequisites
*   Go 1.22+
*   PostgreSQL
*   Redis

### Commands

```bash
# Start the API server locally
make run

# Run all unit and integration tests
make test

# Generate/Update Swagger documentation
make swagger
```

## 🔐 Security Standards

*   **Input Validation**: All inputs are validated at the service layer.
*   **Path Traversal**: Strict checks in `StorageManager` before filesystem access.
*   **SRP**: Single Responsibility Principle is strictly enforced to minimize test pollution and logic coupling.
