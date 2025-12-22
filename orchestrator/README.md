# Orchestrator

**Port:** 9000  
**Status:** ✅ Production  
**Updated:** 2025-12-04

---

## Purpose

Service health monitoring and coordination for NAS.AI infrastructure:
- **Health Check Loop** - Parallel checks every 30s
- **Prometheus Metrics** - `/metrics` endpoint
- **Service Registry** - JSON-based service tracking
- **Thread-Safe** - Race condition fixes applied

---

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Orchestrator health |
| GET | `/metrics` | Prometheus metrics |
| GET | `/api/services` | Service status JSON |
| GET | `/api/registry` | Service registry |

---

## Configuration

| Env Variable | Default | Description |
|--------------|---------|-------------|
| `REGISTRY_PATH` | `./data/registry.json` | Service registry file |
| `API_URL` | `http://localhost:8080` | API base URL |
| `API_ADDR` | `:9000` | Listen address |

---

## Usage

```bash
# Build
make build

# Run
make run

# Custom port
API_ADDR=:9001 make run
```

---

## Architecture

```
orchestrator/
├── orchestrator_loop.go   # Main loop, health checks (parallel)
├── config.go              # Centralized configuration
├── metrics.go             # Prometheus /metrics (thread-safe)
├── registry.go            # Service registry
├── api.go                 # HTTP handlers
└── data/registry.json     # Service definitions
```

---

## Recent Fixes (2025-12-04)

- ✅ Fixed race condition in metrics.go
- ✅ Parallelized health checks with sync.WaitGroup
- ✅ Extracted config to config.go
- ✅ Idiomatic Go error handling

---

**Maintained by:** NAS.AI Team
