# NAS.AI Data Flow (Upload → AI Agent → DB → Search)

## Overview
1) **WebUI**: User lädt Datei hoch (`/api/v1/storage/upload`).
2) **API**: Speichert unter `/mnt/data/...`, erstellt (geplant) Embedding-Job.
3) **AI Knowledge Agent**: Holt Content, erzeugt Embedding (all-MiniLM-L6-v2), speichert in `file_embeddings` (pgvector).
4) **Search**: `/api/v1/search` holt Embedding vom AI-Agent und fragt pgvector via `embedding <=> $1` (cosine).

## Integration Test (manuell/CI)
```bash
# 1) Upload Test-Dokument
API=https://felix-freund.com
TOKEN="<JWT>"
CSRF="<CSRF>"
curl -s -X POST "$API/api/v1/storage/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-CSRF-Token: $CSRF" \
  -F "file=@test_rechnung.txt" \
  -F "path=/"

# 2) Wait for AI Agent to index (poll /health or sleep 5s)

# 3) Validate in DB (pgvector)
docker compose -f infrastructure/docker-compose.prod.yml exec -T postgres \
  psql -U nas_user -d nas_db -c "SELECT file_path, embedding IS NOT NULL AS has_vec FROM file_embeddings WHERE file_path LIKE '%test_rechnung.txt';"

# 4) Search should return the doc
curl -s "$API/api/v1/search?q=Rechnung" -H "Authorization: Bearer $TOKEN" -H "X-CSRF-Token: $CSRF"
```

## Notes / Safety
- Search is scoped to `/mnt/data/**` and excludes `/.trash/`.
- AI Agent startup: loads model + DB connectivity check; embeddings stored in pgvector.
- Fail-fast: If AI offline, `/health` on agent returns `status: degraded`.
