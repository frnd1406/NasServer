#!/bin/bash
# ============================================================
# Secret Generator für NAS.AI Installation
# ============================================================
# Generiert sichere Secrets für .env.prod
# ============================================================

set -e

echo "============================================================"
echo "🔐 NAS.AI Secret Generator"
echo "============================================================"
echo ""

# Farben
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Prüfe ob .env.prod existiert
if [ -f ".env.prod" ]; then
    echo -e "${YELLOW}⚠️  .env.prod existiert bereits!${NC}"
    read -p "Möchten Sie fortfahren? Dies überschreibt die Datei. (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Abgebrochen."
        exit 1
    fi
    # Backup erstellen
    cp .env.prod .env.prod.backup.$(date +%Y%m%d_%H%M%S)
    echo "Backup erstellt: .env.prod.backup.*"
fi

echo ""
echo "Generiere sichere Secrets..."
echo ""

# Secrets generieren
POSTGRES_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)
MONITORING_TOKEN=$(openssl rand -base64 32)

echo -e "${GREEN}✅ Secrets generiert${NC}"
echo ""

# Domain abfragen
echo "============================================================"
echo "Domain-Konfiguration"
echo "============================================================"
read -p "Ihre Domain (z.B. example.com): " DOMAIN
read -p "API Sub-Domain (z.B. api.example.com): " API_DOMAIN
read -p "Email-Absender (z.B. noreply@example.com): " EMAIL_FROM

# .env.prod erstellen
cat > .env.prod << EOF
# ============================================================
# NAS.AI Production Environment Configuration
# Automatisch generiert am $(date)
# ============================================================

# ============================================================
# DATABASE SECRETS (PostgreSQL)
# ============================================================
POSTGRES_DB=nas_db
POSTGRES_USER=nas_user
POSTGRES_PASSWORD=$POSTGRES_PASSWORD

# ============================================================
# SECURITY SECRETS
# ============================================================
JWT_SECRET=$JWT_SECRET
MONITORING_TOKEN=$MONITORING_TOKEN

# ============================================================
# DOMAIN CONFIGURATION
# ============================================================
DOMAIN=$DOMAIN
API_DOMAIN=$API_DOMAIN
CORS_ORIGINS=https://$DOMAIN,https://$API_DOMAIN
FRONTEND_URL=https://$DOMAIN

# ============================================================
# EMAIL CONFIGURATION
# ============================================================
EMAIL_FROM=$EMAIL_FROM

# ============================================================
# OPTIONAL: Advanced Configuration
# ============================================================
LOG_LEVEL=INFO
RATE_LIMIT_PER_MIN=100
BACKUP_SCHEDULE="0 2 * * *"
BACKUP_RETENTION_COUNT=7

# ============================================================
# AI Model Configuration
# ============================================================
AI_MODEL_TYPE=sentence-transformers
AI_MODEL_NAME=sentence-transformers/all-MiniLM-L6-v2
MODEL_LOAD_RETRIES=3
INFERENCE_TIMEOUT=30
EOF

# Berechtigungen setzen
chmod 600 .env.prod

echo ""
echo -e "${GREEN}✅ .env.prod erfolgreich erstellt!${NC}"
echo ""
echo "============================================================"
echo "📋 Generierte Secrets (NUR EINMAL ANGEZEIGT!):"
echo "============================================================"
echo "POSTGRES_PASSWORD: $POSTGRES_PASSWORD"
echo "JWT_SECRET: $JWT_SECRET"
echo "MONITORING_TOKEN: $MONITORING_TOKEN"
echo "============================================================"
echo ""
echo -e "${YELLOW}⚠️  WICHTIG: Speichern Sie diese Werte sicher!${NC}"
echo "Diese werden nur einmal angezeigt und sind danach nur noch in .env.prod sichtbar."
echo ""
echo "Nächste Schritte:"
echo "1. docker compose -f docker-compose.prod.yml build"
echo "2. docker compose -f docker-compose.prod.yml up -d"
echo "3. docker compose -f docker-compose.prod.yml logs -f"
echo ""
