# Service Orchestrator

The conductor of the NAS AI system. This service monitors the health of all other components (API, Database, AI Agent) and manages service discovery.

## 🎯 Responsibilities

1.  **Service Discovery**: dynamically registers available services.
2.  **Health Checks**: Periodically pings services to ensure they are responsive.
3.  **Self-Healing**: triggers restarts or alerts if critical services (like the Database) fail.
4.  **Configuration Management**: Loads system-wide service configuration from `registry.json`.

## ⚙️ Configuration

The `registry.json` file defines the services to watch:

```json
{
  "services": [
    {
      "name": "api-server",
      "url": "http://api:8080/health",
      "critical": true
    },
    {
      "name": "ai-agent",
      "url": "http://ai-service:5000/health",
      "critical": false
    }
  ]
}
```
