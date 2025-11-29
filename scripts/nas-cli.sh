#!/bin/bash
set -euo pipefail

# Colors
RED=$'\033[0;31m'
GREEN=$'\033[0;32m'
YELLOW=$'\033[1;33m'
BLUE=$'\033[0;34m'
CYAN=$'\033[0;36m'
NC=$'\033[0m'

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
INFRA_DIR="$ROOT_DIR/infrastructure"
API_DIR="$INFRA_DIR/api"
DOC_OUTPUT="$ROOT_DIR/API_ENDPOINTS.md"
COMPOSE_FILE="$INFRA_DIR/docker-compose.prod.yml"
DEV_COMPOSE_FILE="$INFRA_DIR/docker-compose.dev.yml"
ENV_FILE="$INFRA_DIR/.env.prod"
PROJECT_NAME="${COMPOSE_PROJECT_NAME:-$(basename "$INFRA_DIR")}"
DC_CMD="docker compose -f \"$COMPOSE_FILE\""
DC_CMD_DEV="docker compose -f \"$DEV_COMPOSE_FILE\""
NETWORK_NAME="${PROJECT_NAME}_nas-network"

API_URL_DEFAULT="https://felix-freund.com"
API_URL="${API_URL:-$API_URL_DEFAULT}"
VERBOSE="${VERBOSE:-false}"
CURL_AUTH_HEADERS=()

print_header() {
  clear
  printf "%s\n" \
    "${BLUE}╔══════════════════════════════════════════════════════════════╗${NC}" \
    "${BLUE}║${NC}   🚀 ${CYAN}NAS.AI Infrastructure CLI (Monolith)${NC}                      ${BLUE}║${NC}" \
    "${BLUE}║${NC}   ${YELLOW}Stabilität · Sicherheit · DX${NC}                              ${BLUE}║${NC}" \
    "${BLUE}╚══════════════════════════════════════════════════════════════╝${NC}" \
    ""
  printf "%s\n" "${CYAN}Root:${NC} ${ROOT_DIR}"
  printf "%s\n" "${CYAN}API URL:${NC} ${API_URL}"
  printf "%s\n" ""
}

fail() { echo -e "${RED}❌ $*${NC}" >&2; exit 1; }
info() { echo -e "${BLUE}ℹ️  $*${NC}"; }
ok()   { echo -e "${GREEN}✅ $*${NC}"; }
warn() { echo -e "${YELLOW}⚠️  $*${NC}"; }

status_explain() {
  case "$1" in
    200) echo "OK" ;;
    201) echo "Created" ;;
    204) echo "No Content" ;;
    400) echo "Bad Request (Payload/Parameter fehlerhaft)" ;;
    401) echo "Unauthorized (Token fehlt/ungültig)" ;;
    402) echo "Payment Required" ;;
    403) echo "Forbidden (keine Berechtigung)" ;;
    404) echo "Not Found (Route/Resource fehlt)" ;;
    409) echo "Conflict (Zustand kollidiert)" ;;
    413) echo "Payload Too Large" ;;
    415) echo "Unsupported Media Type" ;;
    429) echo "Too Many Requests (Rate Limit)" ;;
    500) echo "Internal Server Error" ;;
    502) echo "Bad Gateway" ;;
    503) echo "Service Unavailable" ;;
    504) echo "Gateway Timeout" ;;
    000) echo "Verbindungsfehler" ;;
    *)   echo "Status $1" ;;
  esac
}

status_matches_expected() {
  local code=$1 expected=$2
  IFS='/' read -ra parts <<<"$expected"
  for e in "${parts[@]}"; do
    if [ "$code" = "$e" ]; then
      return 0
    fi
  done
  return 1
}

require_cmd() {
  for cmd in "$@"; do
    command -v "$cmd" >/dev/null 2>&1 || fail "Befehl fehlt: $cmd"
  done
}

check_prereqs() {
  require_cmd docker curl jq git
  [ -f "$COMPOSE_FILE" ] || warn "Hinweis: $COMPOSE_FILE fehlt (Compose-Befehle könnten fehlschlagen)."
}

# --- Utility: curl wrapper ----------------------------------------------------
http_request() {
  local method=$1 path=$2 expected=$3 timeout=${4:-8} data=${5:-}
  local url="${API_URL}${path}"
  local tmp_body
  tmp_body=$(mktemp)
  local code

  if [ -n "$data" ]; then
    code=$(curl -sS -o "$tmp_body" -w '%{http_code}' -X "$method" \
      -H "Content-Type: application/json" \
      "${CURL_AUTH_HEADERS[@]}" \
      --data "$data" \
      --max-time "$timeout" \
      "$url" || echo "000")
  else
    code=$(curl -sS -o "$tmp_body" -w '%{http_code}' -X "$method" \
      "${CURL_AUTH_HEADERS[@]}" \
      --max-time "$timeout" \
      "$url" || echo "000")
  fi

  local body
  body=$(cat "$tmp_body")
  rm -f "$tmp_body"

  local parsed
  parsed=$(echo "$body" | jq -c '.' 2>/dev/null || echo "$body")

  local msg; msg=$(status_explain "$code")
  if status_matches_expected "$code" "$expected"; then
    echo -e "${GREEN}✅${NC} ${method} ${path}  ${GREEN}${code}${NC} ${msg}"
    [ -n "$parsed" ] && [ "$VERBOSE" = "true" ] && echo "$parsed"
    return 0
  else
    echo -e "${RED}❌${NC} ${method} ${path}  ${RED}${code}${NC} ${msg} (expected: ${expected})"
    [ -n "$parsed" ] && echo -e "${YELLOW}Antwort:${NC} $parsed"
    return 1
  fi
}

# --- API Health --------------------------------------------------------------
health_single() {
  info "Starte Einzel-Health-Check gegen ${API_URL}"
  local endpoints=(
    "GET /health 200"
    "GET /api/v1/system/metrics?limit=1 200"
    "GET /api/v1/system/alerts 200"
  )
  local failures=0
  for line in "${endpoints[@]}"; do
    IFS=' ' read -r method path expected <<<"$line"
    http_request "$method" "$path" "$expected" || failures=$((failures+1))
  done
  [ "$failures" -eq 0 ] || fail "$failures Checks fehlgeschlagen"
}

health_monitor() {
  local interval="${CHECK_INTERVAL:-15}"
  info "Monitoring gestartet (Intervall ${interval}s, API ${API_URL}). Strg+C zum Beenden."
  while true; do
    local ts
    ts="$(date '+%H:%M:%S')"
    if http_request GET /health 200 6 >/dev/null; then
      echo -e "${GREEN}[$ts] OK${NC}"
    else
      echo -e "${RED}[$ts] FEHLER${NC}"
    fi
    sleep "$interval"
  done
}

# --- Endpoint Tests ----------------------------------------------------------
test_endpoints() {
  VERBOSE="${VERBOSE:-false}"
  info "Starte Endpoint-Tests (API ${API_URL})"

  local tests=(
    "GET /health 200 Public-Health"
    "GET /api/v1/system/metrics?limit=1 200 Public-Metrics"
    "GET /api/v1/system/alerts 200 Public-Alerts"
    "POST /auth/login 400 Invalid-Login"
    "GET /api/v1/auth/csrf 200 CSRF-Token"
    "GET /api/v1/system/settings 401 Settings-NoAuth"
    "GET /api/v1/backups 401 Backups-NoAuth"
    "GET /api/v1/storage/files?path=/ 401 Storage-NoAuth"
  )

  CURL_AUTH_HEADERS=()
  if [ -n "${JWT_TOKEN:-}" ] && [ -n "${CSRF_TOKEN:-}" ]; then
    CURL_AUTH_HEADERS=(-H "Authorization: Bearer ${JWT_TOKEN}" -H "X-CSRF-Token: ${CSRF_TOKEN}")
    tests+=(
      "GET /api/v1/system/settings 200 Settings"
      "GET /api/v1/backups 200 Backups"
      "GET /api/v1/storage/files?path=/ 200 Storage"
    )
  else
    warn "JWT_TOKEN/CSRF_TOKEN nicht gesetzt – Auth-Tests werden übersprungen."
  fi

  local failures=0
  for t in "${tests[@]}"; do
    IFS=' ' read -r method path expected label <<<"$t"
    echo -e "${BLUE}➜${NC} $label"
    http_request "$method" "$path" "$expected" 8 || failures=$((failures+1))
  done

  [ "$failures" -eq 0 ] && ok "Alle Tests bestanden" || fail "$failures Tests fehlgeschlagen"
}

# --- Fetch tokens via login --------------------------------------------------
login_and_set_tokens() {
  read -r -p "E-Mail: " email
  read -r -s -p "Passwort: " password
  echo ""

  if [ -z "$email" ] || [ -z "$password" ]; then
    warn "E-Mail oder Passwort leer. Abbruch."
    return 1
  fi
  local payload
  payload=$(printf '{"email":"%s","password":"%s"}' "$email" "$password")

  local resp
  resp=$(curl -s -X POST "${API_URL}/auth/login" \
    -H "Content-Type: application/json" \
    --data "$payload")

  local token refresh csrf
  token=$(echo "$resp" | jq -r '.access_token // empty')
  refresh=$(echo "$resp" | jq -r '.refresh_token // empty')
  csrf=$(echo "$resp" | jq -r '.csrf_token // empty')

  if [ -z "$token" ] || [ -z "$csrf" ]; then
    echo -e "${RED}❌ Login fehlgeschlagen. Antwort:${NC} $resp"
    return 1
  fi

  export JWT_TOKEN="$token"
  export REFRESH_TOKEN="$refresh"
  export CSRF_TOKEN="$csrf"

  echo -e "${GREEN}✅ Login OK${NC}"
  echo -e "JWT_TOKEN: $(echo "$JWT_TOKEN" | head -c 24)…"
  echo -e "CSRF_TOKEN: $(echo "$CSRF_TOKEN" | head -c 12)…"
}

# --- Git Savepoint -----------------------------------------------------------
git_savepoint() {
  info "Git Savepoint (Repo: $ROOT_DIR)"
  if [ ! -d "$ROOT_DIR/.git" ]; then
    git -C "$ROOT_DIR" init
    ok "Git init ausgeführt."
  fi

  local status
  status="$(git -C "$ROOT_DIR" status --porcelain)"
  if [ -z "$status" ]; then
    ok "Keine Änderungen – nichts zu tun."
    return 0
  fi

  git -C "$ROOT_DIR" add -A

  read -r -p "Commit-Message: " msg
  if [ -z "$msg" ]; then
    msg="Auto-savepoint $(date '+%Y-%m-%d %H:%M:%S')"
    warn "Leere Message, verwende: $msg"
  fi

  git -C "$ROOT_DIR" commit -m "$msg"

  if git -C "$ROOT_DIR" remote get-url origin >/dev/null 2>&1; then
    local branch
    branch="$(git -C "$ROOT_DIR" rev-parse --abbrev-ref HEAD)"
    git -C "$ROOT_DIR" pull --rebase origin "$branch" || warn "Pull fehlgeschlagen, setze fort."
    git -C "$ROOT_DIR" push -u origin "$branch"
    ok "Commit gepusht."
  else
    warn "Kein Remote konfiguriert – Commit nur lokal."
  fi
}

# --- API Docs ----------------------------------------------------------------
generate_api_docs() {
  info "Generiere API Docs nach $DOC_OUTPUT"
  cat > "$DOC_OUTPUT" <<EOF
# NAS.AI API Dokumentation

Automatisch generiert am: $(date)

**Base URL:** \`${API_URL}\`

## Public Endpoints
- GET /health
- GET /api/v1/system/metrics
- GET /api/v1/system/alerts

## Auth Endpoints
- POST /auth/login
- POST /auth/register
- POST /auth/refresh
- POST /auth/logout

## Geschützte Endpoints (JWT + CSRF)
- GET /api/v1/system/settings
- PUT /api/v1/system/settings/backup
- POST /api/v1/system/validate-path
- GET /api/v1/backups
- POST /api/v1/backups
- POST /api/v1/backups/{id}/restore
- DELETE /api/v1/backups/{id}
- GET /api/v1/storage/files
- POST /api/v1/storage/upload
- DELETE /api/v1/storage/delete
- GET /api/v1/storage/trash
- POST /api/v1/storage/trash/restore/{id}
- DELETE /api/v1/storage/trash/{id}
- POST /api/v1/storage/rename

Hinweis: Weitere Details siehe Swagger (nur in DEV verfügbar) oder Backend-Code unter $API_DIR/src.
EOF
  ok "Dokumentation erstellt."
}

# --- Logs --------------------------------------------------------------------
tail_logs() {
  if [ ! -f "$COMPOSE_FILE" ]; then
    fail "Compose-Datei $COMPOSE_FILE fehlt."
  fi
  read -r -p "Service (leer = alle): " svc
  read -r -p "Ohne Prefix anzeigen? (y/N): " nopfx
  local flags=""
  if [[ "$nopfx" =~ ^[Yy]$ ]]; then
    flags="--no-log-prefix"
  fi
  eval $DC_CMD logs -f --tail=100 $flags ${svc:-}
}

# --- Extended API Health Check (container-free) -----------------------------
api_health_check_full() {
  info "Starte erweiterten API Health Check gegen ${API_URL}"
  local failures=0
  # Erwartete Status-Codes für anonyme Calls
  local checks=(
    "GET /health 200"
    "POST /auth/register 400"
    "POST /auth/login 400"
    "POST /auth/refresh 400"
    "GET /api/v1/auth/csrf 200"
    "GET /api/v1/storage/files 401"
    "POST /api/v1/storage/upload 401"
    "GET /api/v1/backups 401"
    "GET /api/v1/system/settings 401"
  )
  for c in "${checks[@]}"; do
    IFS=' ' read -r m p exp <<<"$c"
    printf "%s " "${BLUE}➜${NC}"
    echo "${m} ${p}"
    http_request "$m" "$p" "$exp" 8 "" || failures=$((failures+1))
  done
  [ "$failures" -eq 0 ] && ok "Alle Basis-Checks bestanden" || fail "$failures Checks fehlgeschlagen"
}

# --- Docker Clean Rebuild ----------------------------------------------------
docker_clean_rebuild() {
  info "Docker Clean Rebuild (Compose: $COMPOSE_FILE)"
  read -r -p "Welcher Service? (all/api/webui/…) [all]: " svc
  svc=${svc:-all}
  read -r -p "Aggressiven Cache prune durchführen? (y/N): " force
  echo ""

  info "[1/5] Cache aufräumen"
  docker image prune -f || true
  docker builder prune -f || true
  if [[ "$force" =~ ^[Yy]$ ]]; then
    warn "Aggressiver prune (--all --volumes)"
    docker builder prune -af || true
    docker system prune -af --volumes || true
  fi

  info "[2/5] Container stoppen"
  if [ "$svc" = "all" ]; then
    eval $DC_CMD down --remove-orphans || true
  else
    eval $DC_CMD stop "$svc" || true
    eval $DC_CMD rm -f "$svc" || true
  fi

  info "[3/5] Alte Images entfernen"
  if [ "$svc" = "all" ]; then
    docker images | grep "nas-" | awk '{print $3}' | xargs -r docker rmi -f || true
  else
    docker images | grep "nas-$svc" | awk '{print $3}' | xargs -r docker rmi -f || true
  fi

  info "[4/5] Neu bauen (--no-cache)"
  (cd "$INFRA_DIR" && if [ "$svc" = "all" ]; then eval $DC_CMD build --no-cache; else eval $DC_CMD build --no-cache "$svc"; fi)

  info "[5/5] Starten"
  if [ "$svc" = "all" ]; then
    eval $DC_CMD up -d
  else
    eval $DC_CMD up -d "$svc"
  fi

  info "Status:"
  if [ "$svc" = "all" ]; then
    eval $DC_CMD ps
  else
    eval $DC_CMD ps "$svc"
  fi
  ok "Rebuild abgeschlossen."
}

# --- Build & Start API (dev compose) ----------------------------------------
build_and_start_api_dev() {
  info "Baue API Image im dev-Modus (Tag: nas-api:lastest)"
  (cd "$API_DIR" && docker build -t nas-api:lastest .)

  if [ ! -f "$DEV_COMPOSE_FILE" ]; then
    warn "Dev-Compose-Datei fehlt: $DEV_COMPOSE_FILE"
    return
  fi

  info "Starte API Service via docker-compose.dev.yml"
  eval $DC_CMD_DEV up -d api

  info "Status:"
  eval $DC_CMD_DEV ps api || true
  ok "API wurde neu gebaut und gestartet."
}

# --- Smart waits for deployment ---------------------------------------------
smart_wait_pg() {
  local retries=30 delay=2
  info "Warte auf Postgres (pg_isready)..."
  for ((i=1; i<=retries; i++)); do
    if eval $DC_CMD exec -T postgres pg_isready -U nas_user -d nas_db -h localhost >/dev/null 2>&1; then
      ok "Postgres ist bereit."
      return 0
    fi
    echo -e "${YELLOW}  Versuch $i/$retries...${NC}"
    sleep "$delay"
  done
  fail "Postgres wurde nicht bereit."
}

smart_wait_api() {
  local retries=30 delay=2
  info "Warte auf API (/health)..."
  for ((i=1; i<=retries; i++)); do
    if docker run --rm --network "$NETWORK_NAME" curlimages/curl:8.9.1 -fsS http://api:8080/health >/dev/null 2>&1; then
      ok "API ist erreichbar."
      return 0
    fi
    echo -e "${YELLOW}  Versuch $i/$retries...${NC}"
    sleep "$delay"
  done
  fail "API (/health) nicht erreichbar."
}

apply_db_seed_if_reset() {
  local do_reset=$1
  if [ "$do_reset" != "true" ]; then
    return 0
  fi
  info "Starte DB-Init (Clean Slate)..."
  if [ -f "$INFRA_DIR/db/init.sql" ]; then
    eval $DC_CMD exec -T postgres psql -U nas_user -d nas_db < "$INFRA_DIR/db/init.sql"
  fi
  if [ -f "$INFRA_DIR/db/migrations/001_add_email_verification.sql" ]; then
    eval $DC_CMD exec -T postgres psql -U nas_user -d nas_db < "$INFRA_DIR/db/migrations/001_add_email_verification.sql"
  fi
  ok "DB-Init abgeschlossen."
}

# --- AI Knowledge Agent Functions --------------------------------------------
ai_db_clear() {
  info "Lösche AI Knowledge Embeddings Tabelle..."

  read -r -p "Wirklich alle AI Embeddings löschen? (y/N): " confirm
  if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    warn "Abgebrochen."
    return 0
  fi

  eval $DC_CMD exec -T postgres psql -U nas_user -d nas_db -c "DROP TABLE IF EXISTS file_embeddings CASCADE;" || fail "DB Drop fehlgeschlagen"
  ok "file_embeddings Tabelle gelöscht."

  info "Erstelle neue leere Tabelle..."
  eval $DC_CMD exec -T postgres psql -U nas_user -d nas_db -c "
    CREATE TABLE IF NOT EXISTS file_embeddings (
      id SERIAL PRIMARY KEY,
      file_id TEXT NOT NULL,
      file_path TEXT NOT NULL,
      mime_type TEXT NOT NULL,
      content TEXT NOT NULL,
      embedding vector(384),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      UNIQUE(file_id)
    );" || fail "Tabellen-Erstellung fehlgeschlagen"

  ok "Neue leere file_embeddings Tabelle erstellt."
}

ai_agent_restart() {
  info "Starte AI Knowledge Agent neu..."

  eval $DC_CMD restart ai-knowledge-agent || fail "Neustart fehlgeschlagen"

  info "Warte auf Model-Loading (20s)..."
  sleep 20

  info "Prüfe AI Agent Logs:"
  eval $DC_CMD logs --tail 15 ai-knowledge-agent

  ok "AI Knowledge Agent neu gestartet."
}

ai_agent_rebuild() {
  info "Rebuild AI Knowledge Agent (mit neuem Code)..."

  info "[1/3] Stoppe Container..."
  eval $DC_CMD stop ai-knowledge-agent || true
  eval $DC_CMD rm -f ai-knowledge-agent || true

  info "[2/3] Rebuild Image..."
  (cd "$INFRA_DIR" && eval $DC_CMD build --no-cache ai-knowledge-agent)

  info "[3/3] Starte Container..."
  eval $DC_CMD up -d ai-knowledge-agent

  info "Warte auf Model-Loading (25s)..."
  sleep 25

  info "Status & Logs:"
  eval $DC_CMD ps ai-knowledge-agent
  eval $DC_CMD logs --tail 15 ai-knowledge-agent

  ok "AI Agent neu gebaut und gestartet."
}

ai_full_reset() {
  echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
  echo -e "${BLUE}║${NC}   🧠 AI Knowledge Agent - Full Reset            ${BLUE}║${NC}"
  echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
  echo ""

  warn "Dies wird durchführen:"
  echo "  1) Alle AI Embeddings aus Datenbank löschen"
  echo "  2) AI Agent Container neu starten"
  echo "  3) Logs zur Kontrolle ausgeben"
  echo ""

  read -r -p "Fortfahren? (y/N): " confirm
  if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    warn "Abgebrochen."
    return 0
  fi

  ai_db_clear
  echo ""
  ai_agent_restart

  ok "Full AI Reset abgeschlossen."
}

ai_show_embeddings() {
  info "Zeige gespeicherte AI Embeddings..."

  eval $DC_CMD exec -T postgres psql -U nas_user -d nas_db -c "
    SELECT
      file_id,
      mime_type,
      LENGTH(content) as content_length,
      384 as embedding_dim,
      created_at
    FROM file_embeddings
    ORDER BY created_at DESC
    LIMIT 20;" || warn "Tabelle existiert möglicherweise nicht"
}

ai_test_endpoint() {
  info "Teste AI Agent /process Endpoint..."

  read -r -p "Dateiname in /mnt/data/ [test.txt]: " filename
  filename=${filename:-test.txt}

  local payload
  payload=$(printf '{"file_path":"/mnt/data/%s","file_id":"%s","mime_type":"text/plain"}' "$filename" "$filename")

  info "Sende Request an AI Agent..."
  eval $DC_CMD exec api sh -c "'curl -X POST http://ai-knowledge-agent:5000/process -H \"Content-Type: application/json\" -d '\''$payload'\'' 2>&1'" || fail "Request fehlgeschlagen"

  echo ""
  info "Prüfe AI Agent Logs:"
  eval $DC_CMD logs --tail 10 ai-knowledge-agent
}

ai_test_embed_query() {
  info "Teste AI Agent /embed_query Endpoint..."

  read -r -p "Suchtext: " searchtext
  if [ -z "$searchtext" ]; then
    warn "Kein Suchtext eingegeben. Abbruch."
    return 0
  fi

  local payload
  payload=$(printf '{"text":"%s"}' "$searchtext")

  info "Sende Request an AI Agent..."
  local response
  response=$(eval $DC_CMD exec api sh -c "'curl -sS -X POST http://ai-knowledge-agent:5000/embed_query -H \"Content-Type: application/json\" -d '\''$payload'\'' 2>&1'")

  echo -e "${GREEN}Response:${NC}"
  echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

ai_test_semantic_search() {
  info "Teste Semantische Suche..."

  read -r -p "Suchanfrage: " query
  if [ -z "$query" ]; then
    warn "Keine Suchanfrage eingegeben. Abbruch."
    return 0
  fi

  info "Sende Suchanfrage an API..."

  # URL-encode the query
  local encoded_query
  encoded_query=$(printf '%s' "$query" | jq -sRr @uri)

  local response
  response=$(eval $DC_CMD exec api sh -c "'curl -sS \"http://localhost:8080/api/v1/search?q=$encoded_query\" 2>&1'")

  echo -e "${GREEN}Suchergebnisse:${NC}"
  echo "$response" | jq '.' 2>/dev/null || echo "$response"

  echo ""
  info "Top 3 Resultate:"
  echo "$response" | jq -r '.results[:3] | .[] | "\nDatei: \(.file_path)\nÄhnlichkeit: \(.similarity)\nInhalt: \(.content[:200])..."' 2>/dev/null || true
}

ai_rebuild_and_restart_all() {
  echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
  echo -e "${BLUE}║${NC}   🚀 AI Agent + API Rebuild & Restart          ${BLUE}║${NC}"
  echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
  echo ""

  warn "Dies wird durchführen:"
  echo "  1) AI Agent neu bauen (mit /embed_query Endpoint)"
  echo "  2) API neu bauen (mit Semantic Search)"
  echo "  3) Beide Container neu starten"
  echo ""

  read -r -p "Fortfahren? (y/N): " confirm
  if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    warn "Abgebrochen."
    return 0
  fi

  info "[1/4] Stoppe Container..."
  eval $DC_CMD stop ai-knowledge-agent api || true

  info "[2/4] Rebuild AI Agent..."
  (cd "$INFRA_DIR" && eval $DC_CMD build --no-cache ai-knowledge-agent)

  info "[3/4] Rebuild API..."
  (cd "$API_DIR" && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o api-new ./src/main.go)
  docker cp "$API_DIR/api-new" nas-api:/app/api

  info "[4/4] Starte Container..."
  eval $DC_CMD up -d ai-knowledge-agent api

  info "Warte auf Services (30s)..."
  sleep 30

  info "Status & Logs:"
  eval $DC_CMD ps ai-knowledge-agent api
  echo ""
  info "AI Agent Logs:"
  eval $DC_CMD logs --tail 10 ai-knowledge-agent
  echo ""
  info "API Logs:"
  eval $DC_CMD logs --tail 10 api

  ok "Rebuild & Restart abgeschlossen."
}

# --- Deploy Prod -------------------------------------------------------------
deploy_prod() {
  [ -f "$ENV_FILE" ] || fail ".env.prod fehlt ($ENV_FILE)"
  [ -f "$COMPOSE_FILE" ] || fail "docker-compose.prod.yml fehlt ($COMPOSE_FILE)"

  echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
  echo -e "${BLUE}║${NC}   🚀 NAS.AI Production Deployment               ${BLUE}║${NC}"
  echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
  echo ""

  read -r -p "Datenbank komplett zurücksetzen? (y/N): " ans
  DO_RESET="false"
  [[ "$ans" =~ ^[Yy]$ ]] && DO_RESET="true"

  info "Stoppe Container..."
  if [ "$DO_RESET" = "true" ]; then
    eval $DC_CMD down -v --remove-orphans
  else
    eval $DC_CMD down --remove-orphans
  fi

  info "Starte Compose (build + up)..."
  eval $DC_CMD up -d --build

  smart_wait_pg
  apply_db_seed_if_reset "$DO_RESET"
  smart_wait_api

  info "Status:"
  eval $DC_CMD ps
  ok "Deployment abgeschlossen."
  echo -e "${YELLOW}Logs:${NC} $DC_CMD logs -f --no-log-prefix"
}

# --- Menu --------------------------------------------------------------------
main_menu() {
  while true; do
    print_header
    printf "%s\n" \
      "${YELLOW}┌──────────────────────────────────────────────┐${NC}" \
      "${YELLOW}│${NC}   ${CYAN}API & System${NC}                              ${YELLOW}│${NC}" \
      "${YELLOW}├──────────────────────────────────────────────┤${NC}" \
      "${YELLOW}│${NC} 1) 🔍 API Health Check (single)             ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 2) 📡 API Monitoring Loop                 ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 3) 🧪 Endpoint Tests                      ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 4) 📚 API Docs generieren                 ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 5) 💾 Git Savepoint (add/commit/push)    ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 6) 📜 Docker Logs (optional no-prefix)   ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 7) ✅ API Health Check (erweitert)       ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 8) 🐳 Docker Clean Rebuild               ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 9) 🚀 Deploy Prod                         ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}10) 🔨 API neu bauen & starten (dev)         ${YELLOW}│${NC}" \
      "${YELLOW}├──────────────────────────────────────────────┤${NC}" \
      "${YELLOW}│${NC}   ${CYAN}AI Knowledge Agent${NC}                        ${YELLOW}│${NC}" \
      "${YELLOW}├──────────────────────────────────────────────┤${NC}" \
      "${YELLOW}│${NC}11) 🗑️  AI DB Clear (Embeddings löschen)      ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}12) 🔄 AI Agent Restart                    ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}13) 🔨 AI Agent Rebuild (neu bauen)       ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}14) 🧠 AI Full Reset (DB + Restart)       ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}15) 📊 AI Embeddings anzeigen             ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}16) 🧪 AI /process Endpoint Test          ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}17) 🔍 AI /embed_query Test               ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}18) 🎯 Semantic Search Test               ${YELLOW}│${NC}" \
      "${YELLOW}│${NC}19) 🚀 AI+API Full Rebuild                ${YELLOW}│${NC}" \
      "${YELLOW}├──────────────────────────────────────────────┤${NC}" \
      "${YELLOW}│${NC} L) 🔐 Login & Tokens setzen               ${YELLOW}│${NC}" \
      "${YELLOW}│${NC} 0) ❌ Beenden                             ${YELLOW}│${NC}" \
      "${YELLOW}└──────────────────────────────────────────────┘${NC}"
    read -r -p "Auswahl: " choice
    case "$choice" in
      1) health_single; read -r -p "Weiter mit Enter..." _ ;;
      2) health_monitor ;;
      3) test_endpoints; read -r -p "Weiter mit Enter..." _ ;;
      4) generate_api_docs; read -r -p "Weiter mit Enter..." _ ;;
      5) git_savepoint; read -r -p "Weiter mit Enter..." _ ;;
      6) tail_logs ;;
      7) api_health_check_full; read -r -p "Weiter mit Enter..." _ ;;
      8) docker_clean_rebuild; read -r -p "Weiter mit Enter..." _ ;;
      9) deploy_prod; read -r -p "Weiter mit Enter..." _ ;;
      10) build_and_start_api_dev; read -r -p "Weiter mit Enter..." _ ;;
      11) ai_db_clear; read -r -p "Weiter mit Enter..." _ ;;
      12) ai_agent_restart; read -r -p "Weiter mit Enter..." _ ;;
      13) ai_agent_rebuild; read -r -p "Weiter mit Enter..." _ ;;
      14) ai_full_reset; read -r -p "Weiter mit Enter..." _ ;;
      15) ai_show_embeddings; read -r -p "Weiter mit Enter..." _ ;;
      16) ai_test_endpoint; read -r -p "Weiter mit Enter..." _ ;;
      17) ai_test_embed_query; read -r -p "Weiter mit Enter..." _ ;;
      18) ai_test_semantic_search; read -r -p "Weiter mit Enter..." _ ;;
      19) ai_rebuild_and_restart_all; read -r -p "Weiter mit Enter..." _ ;;
      [Ll]) login_and_set_tokens; read -r -p "Weiter mit Enter..." _ ;;
      0) exit 0 ;;
      *) warn "Ungültige Auswahl"; sleep 1 ;;
    esac
  done
}

# --- Entry -------------------------------------------------------------------
check_prereqs
main_menu
