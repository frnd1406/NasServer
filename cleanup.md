

##  Final Refactoring Map: Middleware, Models, Drivers &
## Repositories
Ziel: Vollendung der Clean Architecture durch Trennung von Business-Objekten, technischer
Middleware, Infrastruktur-Treibern und Datenzugriff.
 1. Middleware (package middleware)
Trennung von technischem "Rauschen" und Business-Regeln.
Ziel-Struktur:
infrastructure/api/src/middleware/
├── core/           # Technisch: Logging, CORS, Panic, RequestID
└── logic/          # Fachlich: Auth, VaultGuard, RateLimit
## Zuordnung:
logging.go, headers.go, panic.go, requestid.go, cors.go -> middleware/core/
auth.go, auth_test.go, csrf.go, ratelimit.go, vault_guard.go ->
middleware/logic/
Hinweis: Imports in server/routes.go müssen angepasst werden.
 2. Files / Storage Drivers (package storage_drivers)
Isolierung der physischen Speicher-Logik aus dem Service-Layer.
Aktuell liegt das in src/services/storage. Das ist verwirrend neben src/services/content.
Ziel-Struktur:
infrastructure/api/src/drivers/
└── storage/        # Implementierungen: Local, S3, etc.
## Zuordnung:
src/services/storage/provider.go (Interface) ->
src/drivers/storage/interface.go
src/services/storage/local_store.go -> src/drivers/storage/local.go
Vorteil: Der ContentService importiert dann den Driver, nicht umgekehrt.
 3. Repositories (package repository)
Gruppierung der Datenbank-Zugriffe nach Domain.
02.01.26, 12:45Gemini
https://gemini.google.com/gem/6fc0c5ffba13/ddf6f6fb0a95b73b1/3

Ziel-Struktur:
infrastructure/api/src/repository/
├── auth/           # User & Rollen
├── files/          # Datei-Metadaten
├── system/         # Monitoring & Alerts
└── settings/       # System-Config
## Zuordnung:
user_repository.go, user_repository_sqlx.go -> repository/auth/
file_repository.go, file_embeddings_repository.go, honeyfile_repository.go ->
repository/files/
monitoring_repository.go, system_alerts_repository.go,
system_metrics_repository.go -> repository/system/
system_settings_repository.go -> repository/settings/
Hinweis: Passe die Initialisierung in server/server.go an (Server.UserRepo =
auth_repo.New(...)).
 4. Models -> Domain (package domain)
Die atomaren Bausteine der Applikation.
ACHTUNG: Dies ist ein Breaking Change für das gesamte Projekt. Alle Imports von
.../api/src/models werden ungültig. Führe diesen Schritt ZULETZT aus.
Ziel-Struktur: Anstatt eines flachen models Ordners nutzen wir Domain-Pakete.
infrastructure/api/src/domain/
├── auth/           # User, Claims, Role
├── files/          # File, Share, Honeyfile
├── system/         # Metrics, Alert, Monitoring
└── config/         # SystemSettings, Setup
## Zuordnung:
models/user.go -> domain/auth/entity.go
models/file.go -> domain/files/entity.go
models/share.go -> domain/files/share.go
models/system_metrics.go -> domain/system/metrics.go
models/system_alert.go -> domain/system/alert.go
models/monitoring.go -> domain/system/monitoring.go
烙 Strategie für den Agenten (Reihenfolge ist kritisch!)
02.01.26, 12:45Gemini
https://gemini.google.com/gem/6fc0c5ffba13/ddf6f6fb0a95b73b2/3

- Middleware & Drivers: Verschiebe diese zuerst. Das ist risikoarm.
- Repositories: Verschiebe Repositories in Unterordner. Repariere Imports in services.
- Models (The Big Bang):
Verschiebe Models in domain/.
Benenne Packages um (package auth, package files).
Search & Replace: Suche global (grep) nach "github.com/.../models" und
ersetze es durch die neuen Pfade.
Ersetze im Code models.User durch auth.User, models.File durch files.File.
- Finalize:go mod tidy und go build.
02.01.26, 12:45Gemini
https://gemini.google.com/gem/6fc0c5ffba13/ddf6f6fb0a95b73b3/3