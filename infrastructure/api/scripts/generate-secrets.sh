#!/bin/bash

# Generate Secure Secrets for NAS API
# Creates random secure secrets for JWT and other services

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}=== NAS API Secret Generator ===${NC}\n"

# Generate JWT Secret (64 characters, base64)
echo -e "${YELLOW}Generating JWT Secret...${NC}"
JWT_SECRET=$(openssl rand -base64 48 | tr -d '\n')
echo -e "  ${GREEN}✓${NC} Generated (64 characters)\n"

# Generate Session Secret (for cookies, CSRF)
echo -e "${YELLOW}Generating Session Secret...${NC}"
SESSION_SECRET=$(openssl rand -base64 32 | tr -d '\n')
echo -e "  ${GREEN}✓${NC} Generated (43 characters)\n"

# Generate CSRF Secret
echo -e "${YELLOW}Generating CSRF Secret...${NC}"
CSRF_SECRET=$(openssl rand -base64 32 | tr -d '\n')
echo -e "  ${GREEN}✓${NC} Generated (43 characters)\n"

# Create .env file
# Auto-detect base path
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DIR="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$API_DIR/.env"

echo -e "${BLUE}Creating .env file...${NC}"

cat > "$ENV_FILE" <<EOF
# NAS API Environment Variables
# Generated: $(date)
# IMPORTANT: Keep this file SECRET! Add to .gitignore!

# === Server Configuration ===
PORT=8080
ENV=production
LOG_LEVEL=info

# === Security Secrets ===
# CRITICAL: Never commit these to git!
JWT_SECRET=${JWT_SECRET}
SESSION_SECRET=${SESSION_SECRET}
CSRF_SECRET=${CSRF_SECRET}

# === CORS Configuration ===
CORS_ORIGINS=https://your-domain.com,https://api.your-domain.com

# === Frontend URL ===
FRONTEND_URL=https://your-domain.com

# === Rate Limiting ===
RATE_LIMIT_PER_MIN=100

# === Database Configuration ===
DB_HOST=localhost
DB_PORT=5433
DB_USER=nas_user
DB_PASSWORD=nas_dev_password
DB_NAME=nas_db

# === Redis Configuration ===
REDIS_HOST=localhost
REDIS_PORT=6380

# === Email Configuration (Resend) ===
RESEND_API_KEY=re_AEhvFZrx_KRjdCcvcVHLcnPNY66ekBBFy
EMAIL_FROM=NAS.AI <noreply@your-domain.com>

# === Cloudflare Configuration ===
CLOUDFLARE_API_TOKEN=GjKJMQiS998conpswEJhOwQ5b-fKSGjVmmFsofJf
CLOUDFLARE_R2_BUCKET=nas-ai-storage
EOF

chmod 600 "$ENV_FILE"
echo -e "  ${GREEN}✓${NC} Created: $ENV_FILE\n"

# Display secrets for manual export
echo -e "${GREEN}=== Generated Secrets ===${NC}\n"

echo -e "${BLUE}JWT Secret (for authentication):${NC}"
echo -e "export JWT_SECRET='${JWT_SECRET}'"
echo ""

echo -e "${BLUE}Session Secret (for cookies):${NC}"
echo -e "export SESSION_SECRET='${SESSION_SECRET}'"
echo ""

echo -e "${BLUE}CSRF Secret (for CSRF protection):${NC}"
echo -e "export CSRF_SECRET='${CSRF_SECRET}'"
echo ""

# Create export script
EXPORT_FILE="$SCRIPT_DIR/export-env.sh"
cat > "$EXPORT_FILE" <<EOF
#!/bin/bash
# Export environment variables from .env file
# Source this file: source scripts/export-env.sh

export JWT_SECRET='${JWT_SECRET}'
export SESSION_SECRET='${SESSION_SECRET}'
export CSRF_SECRET='${CSRF_SECRET}'
export PORT=8080
export ENV=production
export LOG_LEVEL=info
export CORS_ORIGINS=https://your-domain.com,https://api.your-domain.com
export FRONTEND_URL=https://your-domain.com
export RATE_LIMIT_PER_MIN=100

echo "✓ Environment variables exported"
EOF

chmod +x "$EXPORT_FILE"
echo -e "${GREEN}✓ Created export script: $EXPORT_FILE${NC}\n"

# Update .gitignore
GITIGNORE="$API_DIR/.gitignore"
if [ -f "$GITIGNORE" ]; then
    if ! grep -q "^\.env$" "$GITIGNORE"; then
        echo ".env" >> "$GITIGNORE"
        echo -e "${GREEN}✓ Added .env to .gitignore${NC}\n"
    fi
else
    echo ".env" > "$GITIGNORE"
    echo "scripts/export-env.sh" >> "$GITIGNORE"
    echo -e "${GREEN}✓ Created .gitignore${NC}\n"
fi

# Instructions
echo -e "${YELLOW}=== Usage Instructions ===${NC}\n"

echo -e "${BLUE}Option 1: Source the export script${NC}"
echo -e "  source scripts/export-env.sh"
echo -e "  ./bin/api"
echo ""

echo -e "${BLUE}Option 2: Use .env file with a tool${NC}"
echo -e "  # Install dotenv (if not already)"
echo -e "  go install github.com/joho/godotenv/cmd/godotenv@latest"
echo -e "  godotenv -f .env ./bin/api"
echo ""

echo -e "${BLUE}Option 3: Export manually${NC}"
echo -e "  export JWT_SECRET='${JWT_SECRET}'"
echo -e "  ./bin/api"
echo ""

echo -e "${BLUE}Option 4: Load in systemd service${NC}"
echo -e "  Edit /etc/systemd/system/nas-api.service"
echo -e "  Add: EnvironmentFile=$ENV_FILE"
echo ""

echo -e "${GREEN}=== Security Reminders ===${NC}\n"
echo -e "${YELLOW}⚠${NC}  NEVER commit .env to git!"
echo -e "${YELLOW}⚠${NC}  Keep JWT_SECRET secure - if leaked, regenerate it"
echo -e "${YELLOW}⚠${NC}  Use different secrets for dev/staging/production"
echo -e "${YELLOW}⚠${NC}  Rotate secrets regularly (every 90 days)"
echo ""

echo -e "${GREEN}✓ All secrets generated successfully!${NC}\n"
