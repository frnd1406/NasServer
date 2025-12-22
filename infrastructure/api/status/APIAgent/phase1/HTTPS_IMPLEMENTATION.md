# HTTPS Implementation Guide - NAS API

## Übersicht

Für Production muss die API mit HTTPS laufen. Es gibt mehrere Ansätze:

## Option 1: Reverse Proxy (EMPFOHLEN für Production)

### Mit Nginx als Reverse Proxy

```nginx
# /etc/nginx/sites-available/nas-api
server {
    listen 80;
    server_name api.your-domain.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.your-domain.com;

    # SSL Certificate (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/api.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.your-domain.com/privkey.pem;

    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Proxy to Go API
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (falls benötigt)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Rate Limiting (zusätzlich zu Go Rate Limiter)
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    limit_req zone=api_limit burst=20 nodelay;
}
```

### Let's Encrypt Setup

```bash
# Install Certbot
sudo apt-get update
sudo apt-get install certbot python3-certbot-nginx

# Get Certificate
sudo certbot --nginx -d api.your-domain.com

# Auto-renewal (certbot creates cron job automatically)
sudo certbot renew --dry-run
```

### Systemd Service für API

```ini
# /etc/systemd/system/nas-api.service
[Unit]
Description=NAS API Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=api
WorkingDirectory=/opt/nas-api
Environment="PORT=8080"
Environment="JWT_SECRET=your-secret-from-env"
Environment="DATABASE_URL=postgres://user:pass@localhost/nas"
Environment="REDIS_URL=localhost:6379"
Environment="CORS_ORIGINS=https://your-domain.com"
Environment="FRONTEND_URL=https://your-domain.com"
ExecStart=/opt/nas-api/bin/api
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl enable nas-api
sudo systemctl start nas-api
sudo systemctl status nas-api
```

---

## Option 2: Native HTTPS in Go (für Development/Testing)

### Code Änderungen in `src/main.go`

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/nas-ai/api/src/database"
	"github.com/nas-ai/api/src/handlers"
	"github.com/nas-ai/api/src/middleware"
	"github.com/nas-ai/api/src/services"
	"github.com/sirupsen/logrus"
)

func main() {
	// ... existing setup code ...

	// HTTPS Configuration
	useHTTPS := os.Getenv("USE_HTTPS") == "true"
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,

		// Security timeouts
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Optional: Configure TLS
	if useHTTPS {
		srv.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}

	// Start server in goroutine
	go func() {
		if useHTTPS {
			logger.WithFields(logrus.Fields{
				"port":  cfg.Port,
				"https": true,
			}).Info("Starting HTTPS server...")

			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				logger.WithError(err).Fatal("Failed to start HTTPS server")
			}
		} else {
			logger.WithFields(logrus.Fields{
				"port":  cfg.Port,
				"https": false,
			}).Info("Starting HTTP server...")

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.WithError(err).Fatal("Failed to start HTTP server")
			}
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	// Close database connections
	if err := db.Close(); err != nil {
		logger.WithError(err).Error("Error closing database")
	}
	if err := redis.Close(); err != nil {
		logger.WithError(err).Error("Error closing Redis")
	}

	logger.Info("Server exited gracefully")
}
```

### Self-Signed Certificate (nur für Development!)

```bash
# Generate self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
  -subj "/C=DE/ST=Bavaria/L=Munich/O=NAS/CN=localhost"

# Set environment variables
export USE_HTTPS=true
export TLS_CERT_FILE=./cert.pem
export TLS_KEY_FILE=./key.pem
export PORT=8443

# Run server
./bin/api
```

---

## Option 3: Cloudflare Tunnel (Einfachste Lösung)

### Setup

```bash
# Install cloudflared
wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm64.deb
sudo dpkg -i cloudflared-linux-arm64.deb

# Login to Cloudflare
cloudflared tunnel login

# Create tunnel
cloudflared tunnel create nas-api

# Configure tunnel
cat > ~/.cloudflared/config.yml <<EOF
tunnel: nas-api
credentials-file: /home/user/.cloudflared/<TUNNEL-ID>.json

ingress:
  - hostname: api.your-domain.com
    service: http://localhost:8080
  - service: http_status:404
EOF

# Run tunnel
cloudflared tunnel run nas-api

# Or install as service
sudo cloudflared service install
sudo systemctl start cloudflared
```

**Vorteile:**
- Automatisches HTTPS (Cloudflare managed)
- DDoS Protection
- Kein Port Forwarding nötig
- Kostenlos

---

## Option 4: Docker mit Traefik (für Container Deployments)

### docker-compose.yml

```yaml
version: '3.8'

services:
  traefik:
    image: traefik:v2.10
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.httpchallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web"
      - "--certificatesresolvers.letsencrypt.acme.email=admin@your-domain.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "./letsencrypt:/letsencrypt"

  api:
    build: .
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(\`api.your-domain.com\`)"
      - "traefik.http.routers.api.entrypoints=websecure"
      - "traefik.http.routers.api.tls.certresolver=letsencrypt"
      - "traefik.http.services.api.loadbalancer.server.port=8080"
    environment:
      - PORT=8080
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - postgres
      - redis
```

---

## Empfehlung für dein Setup

### Development (Lokal)
```bash
# HTTP reicht
export PORT=8080
./bin/api
```

### Staging/Production auf Raspberry Pi

**Beste Option: Nginx Reverse Proxy + Let's Encrypt**

```bash
# 1. Setup Nginx
sudo apt-get install nginx

# 2. Get SSL Certificate
sudo certbot --nginx -d api.your-domain.com

# 3. Configure Nginx (siehe oben)
sudo nano /etc/nginx/sites-available/nas-api
sudo ln -s /etc/nginx/sites-available/nas-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# 4. Run API on localhost:8080
export PORT=8080
export CORS_ORIGINS=https://your-domain.com
./bin/api

# 5. Setup Systemd service (siehe oben)
```

### Warum Nginx + Let's Encrypt?

**Vorteile:**
- ✅ Kostenlose SSL Zertifikate (Let's Encrypt)
- ✅ Automatische Renewal
- ✅ Bessere Performance (Nginx ist sehr schnell)
- ✅ Zusätzliche Security Layer
- ✅ Rate Limiting auf Nginx-Ebene
- ✅ Static File Serving (falls benötigt)
- ✅ Load Balancing (für Zukunft)
- ✅ Einfacher zu debuggen
- ✅ Standard in Production

**Nachteile von nativem HTTPS in Go:**
- ❌ Mehr Code zu warten
- ❌ Certificate Management in App
- ❌ Neustart bei Certificate Renewal
- ❌ Kein Caching/Compression aus der Box

---

## Security Checklist für Production

- [ ] HTTPS erzwungen (HTTP → HTTPS redirect)
- [ ] TLS 1.2+ only
- [ ] Strong cipher suites
- [ ] HSTS Header
- [ ] Certificate Auto-Renewal
- [ ] Firewall configured (nur 80, 443 offen)
- [ ] Rate Limiting (Nginx + Go)
- [ ] Security Headers
- [ ] CORS richtig konfiguriert
- [ ] Secrets in Environment Variables (nicht im Code!)
- [ ] Database Credentials gesichert
- [ ] Redis Password gesetzt
- [ ] Logging für Security Events
- [ ] Monitoring (Uptime, Errors)

---

## Nächste Schritte

1. Entscheide dich für einen Ansatz (Empfehlung: Nginx Reverse Proxy)
2. Teste Setup lokal
3. Deploy auf Staging
4. Security Audit
5. Production Deployment
