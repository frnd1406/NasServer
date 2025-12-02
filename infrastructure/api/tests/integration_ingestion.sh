#!/bin/bash
set -euo pipefail

# Usage: JWT_TOKEN=... CSRF_TOKEN=... API_URL=https://felix-freund.com ./integration_ingestion.sh

RED=$'\033[0;31m'; GREEN=$'\033[0;32m'; YELLOW=$'\033[1;33m'; NC=$'\033[0m'

API_URL="${API_URL:-https://felix-freund.com}"
FILE="${FILE:-test_rechnung.txt}"
TOKEN="${JWT_TOKEN:-}"
CSRF="${CSRF_TOKEN:-}"

[ -z "$TOKEN" ] && { echo "${RED}JWT_TOKEN fehlt${NC}"; exit 1; }
[ -z "$CSRF" ] && { echo "${RED}CSRF_TOKEN fehlt${NC}"; exit 1; }

echo "${YELLOW}📤 Upload $FILE ...${NC}"
curl -s -X POST "$API_URL/api/v1/storage/upload" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-CSRF-Token: $CSRF" \
  -F "file=@$FILE" \
  -F "path=/"
echo

echo "${YELLOW}⏳ Warte 5s auf Indexing...${NC}"
sleep 5

echo "${YELLOW}🔎 Suche nach 'Rechnung'...${NC}"
curl -s "$API_URL/api/v1/search?q=Rechnung" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-CSRF-Token: $CSRF" | sed 's/.*"results":/results:/'

echo "${GREEN}✅ Done${NC}"
