# Cloudflare Setup für NAS API

## Situation

Du hast bereits:
- ✅ Domain bei Cloudflare
- ✅ SSL/TLS Zertifikate (Cloudflare managed)
- ✅ Cloudflare API Token im Code

Das macht es VIEL einfacher! 🎉

## Option A: Cloudflare Proxy (EINFACHSTE LÖSUNG)

### So funktioniert es:

```
Browser → HTTPS → Cloudflare (SSL/TLS) → HTTP → Dein Raspberry Pi (Port 8080)
                    ^^^^^^^^                           ^^^^^^^
                    Cloudflare              Deine Go API (kein SSL nötig!)
                    handhabt SSL
```

### Setup Schritte:

#### 1. Cloudflare Dashboard Setup

**DNS Records:**
```
Type: A
Name: api (oder @ für Root)
Content: <Deine Raspberry Pi IP oder DynDNS>
Proxy status: Proxied (Orange Cloud ☁️) ← WICHTIG!
TTL: Auto
```

**SSL/TLS Settings:**
- Gehe zu: SSL/TLS → Overview
- Wähle: **"Flexible"** oder **"Full"**

**Flexible:** Cloudflare → Server = HTTP (einfacher, für Start ok)
```
Browser --HTTPS--> Cloudflare --HTTP--> Dein Server
```

**Full:** Cloudflare → Server = HTTPS (besser, aber braucht self-signed cert)
```
Browser --HTTPS--> Cloudflare --HTTPS--> Dein Server
```

**Empfehlung für Start: Flexible**

#### 2. Port Forwarding auf deinem Router

```
External Port: 80
Internal IP: <Raspberry Pi IP>
Internal Port: 8080
Protocol: TCP
```

**ODER** wenn du Port 8080 direkt freigeben willst:
```
External Port: 8080
Internal IP: <Raspberry Pi IP>
Internal Port: 8080
Protocol: TCP
```

#### 3. Deine Go API läuft auf HTTP (kein SSL!)

```bash
export PORT=8080
export CORS_ORIGINS=https://api.your-domain.com,https://your-domain.com
export FRONTEND_URL=https://your-domain.com
./bin/api
```

**Das ist alles!** Cloudflare handhabt SSL automatisch.

---

## Option B: Cloudflare Tunnel (KEINE PORT FORWARDING!)

**Das ist DIE Lösung für Homelab/Raspberry Pi!**

### Vorteile:
- ✅ Kein Port Forwarding nötig!
- ✅ Kein DynDNS nötig!
- ✅ Automatisches SSL
- ✅ DDoS Protection
- ✅ Sicherer (keine offenen Ports)
- ✅ Kostenlos!

### Setup:

#### 1. Install Cloudflared auf Raspberry Pi

```bash
# Download für ARM64 (Raspberry Pi 4/5)
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm64

# Oder für ARM (ältere Pis)
# wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm

# Make executable
chmod +x cloudflared-linux-arm64
sudo mv cloudflared-linux-arm64 /usr/local/bin/cloudflared

# Verify
cloudflared --version
```

#### 2. Login zu Cloudflare

```bash
cloudflared tunnel login
```

Das öffnet Browser → Du authorisierst → Certificate wird gespeichert

#### 3. Create Tunnel

```bash
# Create tunnel
cloudflared tunnel create nas-api

# Output zeigt dir:
# Created tunnel nas-api with id: <TUNNEL-ID>
# Credentials written to: /home/user/.cloudflared/<TUNNEL-ID>.json
```

#### 4. Configure Tunnel

```bash
# Create config file
nano ~/.cloudflared/config.yml
```

```yaml
tunnel: nas-api
credentials-file: /home/user/.cloudflared/<TUNNEL-ID>.json

ingress:
  # API Endpoint
  - hostname: api.your-domain.com
    service: http://localhost:8080
    originRequest:
      noTLSVerify: true
      connectTimeout: 30s

  # Falls du auch Frontend über Tunnel laufen lassen willst
  - hostname: your-domain.com
    service: http://localhost:5173

  # Catchall (required)
  - service: http_status:404
```

#### 5. Add DNS Record in Cloudflare

```bash
cloudflared tunnel route dns nas-api api.your-domain.com
```

Oder manuell im Dashboard:
```
Type: CNAME
Name: api
Content: <TUNNEL-ID>.cfargotunnel.com
Proxy: Proxied (Orange Cloud)
```

#### 6. Run Tunnel

```bash
# Test run
cloudflared tunnel run nas-api

# Wenn alles funktioniert, install als Service
sudo cloudflared service install
sudo systemctl start cloudflared
sudo systemctl enable cloudflared
sudo systemctl status cloudflared
```

#### 7. Start deine Go API (auf localhost)

```bash
export PORT=8080
export CORS_ORIGINS=https://api.your-domain.com,https://your-domain.com
./bin/api
```

**FERTIG!** API ist jetzt erreichbar via HTTPS ohne Port Forwarding!

---

## Cloudflare Security Settings

### In Cloudflare Dashboard

#### SSL/TLS Settings
```
Encryption mode: Full (Strict) oder Full
  - Full: Cloudflare ↔ Server encrypted (self-signed ok)
  - Full (Strict): Cloudflare ↔ Server encrypted (valid cert needed)
  - Für Tunnel: Beide funktionieren

Always Use HTTPS: ON
Automatic HTTPS Rewrites: ON
Minimum TLS Version: 1.2
TLS 1.3: Enabled
```

#### Firewall Rules (optional)
```
Rule 1: Block bad bots
  - If: Known Bots = ON
  - Then: Block

Rule 2: Rate Limiting
  - If: Requests > 100/min from single IP
  - Then: Challenge or Block
```

#### Security Settings
```
Security Level: Medium
Challenge Passage: 30 minutes
Browser Integrity Check: ON
```

---

## Code Anpassungen (Optional)

### Trusted Proxies für Cloudflare

In `src/main.go`:

```go
func main() {
    // ... setup code ...

    // Trust Cloudflare IPs
    // Cloudflare sendet echte Client IP in CF-Connecting-IP Header
    router.SetTrustedProxies([]string{
        "173.245.48.0/20",
        "103.21.244.0/22",
        "103.22.200.0/22",
        "103.31.4.0/22",
        "141.101.64.0/18",
        "108.162.192.0/18",
        "190.93.240.0/20",
        "188.114.96.0/20",
        "197.234.240.0/22",
        "198.41.128.0/17",
        "162.158.0.0/15",
        "104.16.0.0/13",
        "104.24.0.0/14",
        "172.64.0.0/13",
        "131.0.72.0/22",
    })

    // Get real client IP from Cloudflare
    router.Use(func(c *gin.Context) {
        // Cloudflare sends real IP in CF-Connecting-IP
        if cfIP := c.GetHeader("CF-Connecting-IP"); cfIP != "" {
            c.Request.RemoteAddr = cfIP
        }
        c.Next()
    })
}
```

### CORS mit Cloudflare

```go
// In deiner CORS Config
corsConfig := cors.Config{
    AllowOrigins: []string{
        "https://your-domain.com",
        "https://api.your-domain.com",
    },
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

---

## Monitoring & Debugging

### Cloudflare Analytics

Dashboard → Analytics → Traffic
- Request count
- Bandwidth
- Response status codes
- Top countries
- Threats blocked

### Tunnel Logs

```bash
# Live logs
sudo journalctl -u cloudflared -f

# Recent logs
sudo journalctl -u cloudflared -n 100
```

### Test Setup

```bash
# Test from outside
curl https://api.your-domain.com/health

# Should return 200 OK
```

---

## Vergleich der Optionen

### Option A: Cloudflare Proxy (mit Port Forward)

**Pro:**
- Simple setup
- Cloudflare CDN caching
- DDoS protection

**Con:**
- Port Forwarding nötig
- DynDNS setup (falls keine statische IP)
- Port 80/8080 muss offen sein

### Option B: Cloudflare Tunnel (EMPFOHLEN!)

**Pro:**
- ✅ Kein Port Forwarding
- ✅ Kein DynDNS
- ✅ Sicherer
- ✅ Einfacher zu warten
- ✅ Funktioniert auch hinter CGNAT

**Con:**
- Cloudflared Daemon muss laufen
- Leicht höhere Latency (minimal)

---

## Meine Empfehlung für dich

**Cloudflare Tunnel** ist perfekt für dein Setup!

### Quick Start:

```bash
# 1. Install cloudflared
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm64
chmod +x cloudflared-linux-arm64
sudo mv cloudflared-linux-arm64 /usr/local/bin/cloudflared

# 2. Login
cloudflared tunnel login

# 3. Create tunnel
cloudflared tunnel create nas-api

# 4. Configure (siehe config.yml oben)
nano ~/.cloudflared/config.yml

# 5. Route DNS
cloudflared tunnel route dns nas-api api.your-domain.com

# 6. Install service
sudo cloudflared service install
sudo systemctl start cloudflared
sudo systemctl enable cloudflared

# 7. Start API
export PORT=8080
export CORS_ORIGINS=https://your-domain.com
./bin/api
```

**Fertig! HTTPS läuft!** 🚀

---

## Was du schon hast

Ich sehe du hast bereits:
```go
// In deinem Code
export CLOUDFLARE_API_TOKEN="GjKJMQiS998conpswEJhOwQ5b-fKSGjVmmFsofJf"
```

Das ist dein **Cloudflare API Token** - der ist **NICHT** für Tunnel Login!

**Für Tunnel:** Du machst `cloudflare tunnel login` im Browser
**Für API Calls:** Du nutzt den API Token (für DNS updates, etc.)

---

## Next Steps

Soll ich dir helfen:
1. ✅ Cloudflare Tunnel einzurichten?
2. ✅ Systemd Service für API + Tunnel zu erstellen?
3. ✅ Die CORS/Proxy Settings in deinem Code anzupassen?

Was möchtest du zuerst angehen?
