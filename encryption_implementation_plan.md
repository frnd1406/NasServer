# 🔐 Enterprise Encryption System - Implementation Plan

> **Ziel**: Zero-Knowledge Encryption für NAS Server
> **Standard**: AES-256-GCM + Argon2id Key Derivation
> **Priorität**: KRITISCH - Security-Fundament

---

## 📐 Architektur-Übersicht

```
┌────────────────────────────────────────────────────────────────────────────┐
│                           USER LAYER                                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                       │
│  │   Web UI     │  │   SSH CLI    │  │  Mobile App  │                       │
│  │   (Unlock)   │  │   (Unlock)   │  │   (Future)   │                       │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                       │
└─────────┼─────────────────┼─────────────────┼───────────────────────────────┘
          │                 │                 │
          ▼                 ▼                 ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                           API LAYER (Go)                                    │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     Encryption Service                               │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │   │
│  │  │ Key Manager │  │ Vault API   │  │ Unlock API  │  │ Status API │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └────────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
          │
          ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                           ENCRYPTION LAYER                                  │
│  ┌──────────────────────────┐  ┌──────────────────────────────────────┐    │
│  │   File Encryption        │  │   Database Encryption                │    │
│  │   (AES-256-GCM)          │  │   (SQLCipher)                        │    │
│  │                          │  │                                      │    │
│  │   - Encrypt on Upload    │  │   - All tables encrypted             │    │
│  │   - Decrypt on Download  │  │   - Key from Key Manager             │    │
│  │   - Streaming support    │  │   - Transparent to application       │    │
│  └──────────────────────────┘  └──────────────────────────────────────┘    │
└────────────────────────────────────────────────────────────────────────────┘
          │
          ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                           STORAGE LAYER                                     │
│  ┌──────────────────────────┐  ┌──────────────────────────────────────┐    │
│  │   /mnt/data/encrypted/   │  │   /var/lib/nas/db.encrypted          │    │
│  │   (Encrypted Files)      │  │   (SQLCipher Database)               │    │
│  └──────────────────────────┘  └──────────────────────────────────────┘    │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │   /var/lib/nas/vault/                                                 │  │
│  │   ├── encrypted_dek.bin   (DEK verschlüsselt mit Master Key)         │  │
│  │   ├── salt.bin            (Salt für Key Derivation)                  │  │
│  │   └── config.json         (Verschlüsselungs-Konfiguration)           │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────────┘
```

---

## 🔑 Key Hierarchy

```
MASTER KEY (nur beim User)
    │
    │  Argon2id(password, salt, hardware_id)
    │
    ▼
KEY ENCRYPTION KEY (KEK)
    │
    │  AES-256-GCM Decrypt
    │
    ▼
DATA ENCRYPTION KEY (DEK)
    │
    │  Nur im RAM während System läuft
    │
    ├────────────────┬────────────────┐
    ▼                ▼                ▼
FILE_KEY         DB_KEY          BACKUP_KEY
(für Dateien)  (für SQLCipher)  (für Backups)
```

---

## 📦 Komponenten

### 1. Encryption Service (Go)

**Datei**: `infrastructure/api/src/services/encryption_service.go`

```go
type EncryptionService struct {
    isUnlocked    bool
    dek           []byte        // Data Encryption Key (nur im RAM)
    dbKey         []byte        // Database Key
    vaultPath     string
    mu            sync.RWMutex
}

// Methoden
- NewEncryptionService(vaultPath string) *EncryptionService
- Setup(masterPassword string) error          // Erstmaliges Setup
- Unlock(masterPassword string) error         // System entsperren
- Lock() error                                // System sperren
- IsUnlocked() bool                           // Status prüfen
- EncryptFile(plaintext []byte) ([]byte, error)
- DecryptFile(ciphertext []byte) ([]byte, error)
- GetDatabaseKey() ([]byte, error)
- RotateKeys(newMasterPassword string) error  // Key-Rotation
```

### 2. Vault Handler (API Endpoints)

**Datei**: `infrastructure/api/src/handlers/vault.go`

| Endpoint | Method | Beschreibung |
|----------|--------|-------------|
| `/api/v1/vault/status` | GET | Locked/Unlocked Status |
| `/api/v1/vault/unlock` | POST | Mit Master-Passwort entsperren |
| `/api/v1/vault/lock` | POST | System sperren |
| `/api/v1/vault/setup` | POST | Erstmaliges Setup |
| `/api/v1/vault/change-password` | POST | Master-Passwort ändern |

### 3. File Encryption Layer

**Datei**: `infrastructure/api/src/services/encrypted_storage.go`

```go
type EncryptedStorageService struct {
    baseStorage      *StorageService
    encryptionSvc    *EncryptionService
}

// Wrapper um bestehenden StorageService
- Save() → Verschlüsselt vor dem Speichern
- Open() → Entschlüsselt beim Lesen
- Stream-basiert für große Dateien (Chunk-Encryption)
```

### 4. Database Encryption (SQLCipher)

**Migration von SQLite zu SQLCipher:**

```go
import "github.com/mutecomm/go-sqlcipher/v4"

// Connection mit Encryption
db, err := sql.Open("sqlite3", "file:data.db?_pragma_key=x'HEX_KEY'&_pragma_cipher=aes-256-gcm")
```

### 5. Unlock UI (Frontend)

**Datei**: `infrastructure/webui/src/pages/Unlock.jsx`

```
┌─────────────────────────────────────────────┐
│                                             │
│         🔒 NAS Server Locked                │
│                                             │
│   ┌─────────────────────────────────────┐   │
│   │  Master Password                    │   │
│   │  ●●●●●●●●●●●●                       │   │
│   └─────────────────────────────────────┘   │
│                                             │
│   ┌─────────────────────────────────────┐   │
│   │  2FA Code (optional)                │   │
│   │  ______                             │   │
│   └─────────────────────────────────────┘   │
│                                             │
│         [ 🔓 Unlock Server ]                │
│                                             │
└─────────────────────────────────────────────┘
```

---

## 📋 Implementierungs-Reihenfolge

### Phase 1: Core Encryption (Priorität: KRITISCH)

- [ ] **1.1** Encryption Service Grundstruktur
  - Argon2id Key Derivation
  - AES-256-GCM Encrypt/Decrypt
  - Vault-Dateistruktur

- [ ] **1.2** Vault API Endpoints
  - `/vault/status`
  - `/vault/setup`
  - `/vault/unlock`
  - `/vault/lock`

- [ ] **1.3** Startup-Flow
  - System startet im "Locked" Modus
  - Alle Storage-APIs blockiert bis Unlock

### Phase 2: File Encryption (Priorität: HOCH)

- [ ] **2.1** EncryptedStorageService
  - Wrapper um StorageService
  - Encrypt-on-Write
  - Decrypt-on-Read

- [ ] **2.2** Streaming Encryption
  - Chunk-basiert für große Dateien
  - Memory-effizient

- [ ] **2.3** Migration bestehender Dateien
  - Tool zum Verschlüsseln existierender Daten

### Phase 3: Database Encryption (Priorität: HOCH)

- [ ] **3.1** SQLCipher Integration
  - go-sqlcipher Dependency
  - Connection mit Encryption Key

- [ ] **3.2** Migration
  - Bestehende DB → SQLCipher konvertieren
  - Backup vor Migration

### Phase 4: Frontend (Priorität: MITTEL)

- [ ] **4.1** Unlock Page
  - Schönes UI für Master-Passwort
  - Loading States
  - Error Handling

- [ ] **4.2** Locked State Handling
  - Redirect zu Unlock wenn locked
  - Status-Anzeige in Header

- [ ] **4.3** Setup Wizard
  - Erstmaliges Passwort setzen
  - Backup-Key generieren

### Phase 5: Hardening (Priorität: HOCH)

- [ ] **5.1** Memory Protection
  - DEK im RAM schützen (mlock)
  - Secure Memory Wipe bei Lock

- [ ] **5.2** Brute-Force Protection
  - Rate Limiting auf Unlock
  - Account Lockout nach X Versuchen

- [ ] **5.3** Audit Logging
  - Unlock/Lock Events loggen
  - Failed Attempts loggen

---

## 🔒 Sicherheits-Überlegungen

### Was ist geschützt:

| Komponente | Verschlüsselung | Schlüssel |
|------------|-----------------|-----------|
| Dateien | AES-256-GCM | DEK (RAM) |
| Datenbank | SQLCipher (AES-256) | DB_KEY (RAM) |
| DEK auf Disk | AES-256-GCM | KEK (von Master-Passwort) |
| Backups | AES-256-GCM | BACKUP_KEY (RAM) |

### Attack Vectors:

| Angriff | Schutz |
|---------|--------|
| Physischer Zugriff (Festplatte) | ✅ Alles verschlüsselt |
| Netzwerk-Angriff | ✅ HTTPS + Auth |
| Brute-Force Passwort | ✅ Argon2id (langsam) + Rate Limit |
| Memory Dump | ⚠️ mlock + schnelles Wipe |
| Cold Boot Attack | ⚠️ RAM verschwindet nach Shutdown |
| Insider mit Root-Zugang | ⚠️ Schwer zu verhindern wenn System läuft |

### Recovery:

```
┌─────────────────────────────────────────────────────────────────┐
│  RECOVERY KEY (beim Setup generiert)                             │
│                                                                  │
│  24 Wörter BIP39 Mnemonic:                                       │
│  "apple banana cherry dog elephant ..."                          │
│                                                                  │
│  → User muss sicher aufbewahren (offline!)                       │
│  → Kann Master-Passwort zurücksetzen                             │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📁 Neue Dateien

```
infrastructure/api/src/
├── services/
│   ├── encryption_service.go      [NEU]
│   └── encrypted_storage.go       [NEU]
├── handlers/
│   └── vault.go                   [NEU]
└── middleware/
    └── vault_guard.go             [NEU] (blockiert APIs wenn locked)

infrastructure/webui/src/
├── pages/
│   ├── Unlock.jsx                 [NEU]
│   └── VaultSetup.jsx             [NEU]
└── components/
    └── LockedOverlay.jsx          [NEU]
```

---

## ⚙️ Konfiguration

```yaml
# config.yaml
encryption:
  enabled: true
  algorithm: "aes-256-gcm"
  key_derivation: "argon2id"
  argon2_memory: 65536      # 64 MB
  argon2_iterations: 3
  argon2_parallelism: 4
  vault_path: "/var/lib/nas/vault"
  auto_lock_timeout: 0      # 0 = nie (nur bei Shutdown)
```

---

## ✅ Akzeptanzkriterien

1. **Locked State**: Server startet immer im gesperrten Zustand
2. **Zero Knowledge**: Ohne Master-Passwort sind alle Daten Müll
3. **Performance**: Maximale Latenz +10ms pro File-Operation
4. **Streaming**: Dateien >1GB funktionieren ohne Memory-Explosion
5. **Recovery**: Mit Recovery-Key kann Passwort zurückgesetzt werden
6. **Audit**: Alle Security-Events werden geloggt

---

## 🚀 Nächster Schritt

**Phase 1.1 starten**: Encryption Service Grundstruktur implementieren

Soll ich beginnen?
