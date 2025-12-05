# NAS.AI API Dokumentation

Automatisch generiert am: 2025-11-27 17:30:30

## Übersicht

Diese Dokumentation beschreibt alle verfügbaren API-Endpunkte der NAS.AI Plattform.

**Base URL:** `https://felix-freund.com`

## Authentifizierung

Die meisten Endpunkte erfordern Authentifizierung mittels JWT Token.

### Login

```bash
curl -X POST https://felix-freund.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "your-email@example.com"
  }
}
```

### CSRF Token abrufen

```bash
curl -X GET https://felix-freund.com/api/v1/auth/csrf \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Authentifizierte Requests

```bash
curl -X GET https://felix-freund.com/api/v1/ENDPOINT \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

---

## Public Endpoints

### Health Check
**GET** `/health`

Überprüft den Status der API und Datenbank-Verbindungen.

```bash
curl https://felix-freund.com/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "postgres": "connected",
  "redis": "connected"
}
```

### System Metrics
**GET** `/api/v1/system/metrics`

Ruft System-Metriken ab (CPU, RAM, Disk).

**Query Parameter:**
- `limit` (optional): Anzahl der zurückzugebenden Einträge

```bash
curl https://felix-freund.com/api/v1/system/metrics?limit=5
```

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": "uuid",
      "cpu_usage": 45.2,
      "ram_usage": 62.8,
      "disk_usage": 38.5,
      "timestamp": "2025-11-27T16:00:00Z"
    }
  ]
}
```

### System Alerts
**GET** `/api/v1/system/alerts`

Ruft System-Warnungen ab.

```bash
curl https://felix-freund.com/api/v1/system/alerts
```

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": "uuid",
      "severity": "warning",
      "message": "High CPU usage detected",
      "resolved": false,
      "created_at": "2025-11-27T16:00:00Z"
    }
  ]
}
```

---

## Authentication Endpoints

### Register
**POST** `/auth/register`

Registriert einen neuen Benutzer.

```bash
curl -X POST https://felix-freund.com/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "new-user@example.com",
    "password": "secure-password"
  }'
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "user_id": "uuid"
}
```

### Login
**POST** `/auth/login`

Authentifiziert einen Benutzer.

```bash
curl -X POST https://felix-freund.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }'
```

**Response (200 OK):**
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": "uuid",
    "email": "user@example.com"
  }
}
```

### Logout
**POST** `/auth/logout`

Meldet den Benutzer ab (erfordert Authentifizierung).

```bash
curl -X POST https://felix-freund.com/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Logged out successfully"
}
```

### Refresh Token
**POST** `/auth/refresh`

Erneuert den JWT Token.

```bash
curl -X POST https://felix-freund.com/auth/refresh \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "token": "new-jwt-token-here"
}
```

---

## Protected Endpoints (Require Authentication)

### System Settings

#### Get Settings
**GET** `/api/v1/system/settings`

Ruft System-Einstellungen ab (z.B. Backup-Konfiguration).

```bash
curl https://felix-freund.com/api/v1/system/settings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "backup": {
    "schedule": "0 3 * * *",
    "retention": 7,
    "path": "/mnt/backups"
  }
}
```

#### Update Backup Settings
**PUT** `/api/v1/system/settings/backup`

Aktualisiert Backup-Einstellungen.

```bash
curl -X PUT https://felix-freund.com/api/v1/system/settings/backup \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "schedule": "0 2 * * *",
    "retention": 14,
    "path": "/mnt/backups"
  }'
```

**Response (200 OK):**
```json
{
  "backup": {
    "schedule": "0 2 * * *",
    "retention": 14,
    "path": "/mnt/backups"
  }
}
```

---

### Backups

#### List Backups
**GET** `/api/v1/backups`

Listet alle verfügbaren Backups auf.

```bash
curl https://felix-freund.com/api/v1/backups \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": "backup-20251127T030000Z.tar.gz",
      "name": "backup-20251127T030000Z.tar.gz",
      "size": 1024000,
      "created_at": "2025-11-27T03:00:00Z"
    }
  ]
}
```

#### Create Backup
**POST** `/api/v1/backups`

Erstellt ein neues Backup.

```bash
curl -X POST https://felix-freund.com/api/v1/backups \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (201 Created):**
```json
{
  "message": "Backup created successfully",
  "backup_id": "backup-20251127T160000Z.tar.gz"
}
```

#### Restore Backup
**POST** `/api/v1/backups/:id/restore`

Stellt ein Backup wieder her.

```bash
curl -X POST https://felix-freund.com/api/v1/backups/backup-20251127T030000Z.tar.gz/restore \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Backup restored successfully"
}
```

#### Delete Backup
**DELETE** `/api/v1/backups/:id`

Löscht ein Backup.

```bash
curl -X DELETE https://felix-freund.com/api/v1/backups/backup-20251127T030000Z.tar.gz \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Backup deleted successfully"
}
```

---

### Storage / File Management

#### List Files
**GET** `/api/v1/storage/files`

Listet Dateien in einem Verzeichnis auf.

**Query Parameter:**
- `path` (required): Pfad zum Verzeichnis

```bash
curl "https://felix-freund.com/api/v1/storage/files?path=/" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "items": [
    {
      "name": "documents",
      "isDir": true,
      "size": 0,
      "modTime": "2025-11-27T12:00:00Z"
    },
    {
      "name": "photo.jpg",
      "isDir": false,
      "size": 204800,
      "modTime": "2025-11-27T14:30:00Z"
    }
  ]
}
```

#### Upload File
**POST** `/api/v1/storage/upload`

Lädt eine Datei hoch.

```bash
curl -X POST https://felix-freund.com/api/v1/storage/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -F "file=@/path/to/local/file.txt" \
  -F "path=/"
```

**Response (200 OK):**
```json
{
  "message": "File uploaded successfully",
  "path": "/file.txt"
}
```

#### Download File
**GET** `/api/v1/storage/download`

Lädt eine Datei herunter.

**Query Parameter:**
- `path` (required): Pfad zur Datei

```bash
curl "https://felix-freund.com/api/v1/storage/download?path=/photo.jpg" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -o photo.jpg
```

#### Rename File
**POST** `/api/v1/storage/rename`

Benennt eine Datei oder ein Verzeichnis um.

```bash
curl -X POST https://felix-freund.com/api/v1/storage/rename \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "oldPath": "/photo.jpg",
    "newPath": "/vacation.jpg"
  }'
```

**Response (200 OK):**
```json
{
  "message": "File renamed successfully"
}
```

#### Delete File
**DELETE** `/api/v1/storage/delete`

Verschiebt eine Datei in den Papierkorb.

**Query Parameter:**
- `path` (required): Pfad zur Datei

```bash
curl -X DELETE "https://felix-freund.com/api/v1/storage/delete?path=/old-file.txt" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "File moved to trash"
}
```

#### List Trash
**GET** `/api/v1/storage/trash`

Listet gelöschte Dateien im Papierkorb auf.

```bash
curl https://felix-freund.com/api/v1/storage/trash \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": "uuid",
      "name": "old-file.txt",
      "originalPath": "/old-file.txt",
      "deletedAt": "2025-11-27T15:00:00Z"
    }
  ]
}
```

#### Restore from Trash
**POST** `/api/v1/storage/trash/restore/:id`

Stellt eine Datei aus dem Papierkorb wieder her.

```bash
curl -X POST https://felix-freund.com/api/v1/storage/trash/restore/uuid \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "File restored successfully"
}
```

#### Delete from Trash Permanently
**DELETE** `/api/v1/storage/trash/:id`

Löscht eine Datei permanent aus dem Papierkorb.

```bash
curl -X DELETE https://felix-freund.com/api/v1/storage/trash/uuid \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "File permanently deleted"
}
```

---

## Status Codes

| Code | Bedeutung |
|------|-----------|
| 200  | OK - Request erfolgreich |
| 201  | Created - Ressource erstellt |
| 204  | No Content - Request erfolgreich, keine Daten zurück |
| 400  | Bad Request - Ungültige Anfrage |
| 401  | Unauthorized - Authentifizierung erforderlich |
| 403  | Forbidden - Keine Berechtigung |
| 404  | Not Found - Ressource nicht gefunden |
| 500  | Internal Server Error - Server-Fehler |

---

## Error Response Format

Alle Fehler werden im folgenden Format zurückgegeben:

```json
{
  "error": "Beschreibung des Fehlers"
}
```

---

## Rate Limiting

- **Standard:** 100 Requests/Minute
- **Auth Endpoints:** 5 Requests/Minute

Bei Überschreitung wird Status Code `429 Too Many Requests` zurückgegeben.

---

## CORS

Erlaubte Origins werden über die Umgebungsvariable `CORS_ORIGINS` konfiguriert.

---

**Generiert am:** 2025-11-27 17:30:30
**Version:** 1.0.0
