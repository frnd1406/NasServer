# NAS AI Backend API

Das Herzstück des Systems, geschrieben in **Go (Golang)**.

## 🏛 Core Philosophy

**"Thin Handlers, Fat Services"**

*   **Handlers**: Zuständig für HTTP-Parsing, Validierung und Response-Mapping. **Keine Business-Logik!**
*   **Services**: Enthalten die eigentliche Logik, Sicherheits-Checks und Datenverarbeitung.

## 📦 Service Catalog

*   **[ArchiveService](./src/services/archive_service.go)**
    *   Sicherer Upload und Download von Archiven.
    *   Schutz gegen Zip-Slip und Zip-Bomb Angriffe.

*   **[ContentDeliveryService](./src/services/content_delivery_service.go)**
    *   Streaming von Multimedia-Inhalten.
    *   Behandelt Range-Requests (auch für verschlüsselte Dateien).
    *   On-the-fly Entschlüsselung.

*   **[AIAgentService](./src/services/ai_agent_service.go)**
    *   Gateway zum Python-basierten AI Knowledge Agent.

*   **[EncryptionService](./src/services/encryption_service.go)**
    *   Verwaltet Master-Keys und Datei-Verschlüsselung (XChaCha20-Poly1305).

## 🛠 Development

### Setup

```bash
# Kopiere Beispiel-Konfiguration
cp .env.example .env
```

### Commands

```bash
make run      # Startet den Server lokal
make test     # Führt Unittests aus
make swagger  # Aktualisiert API-Dokumentation
```
