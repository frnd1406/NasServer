#!/bin/sh
set -e

# Adjust ownership of mounted volumes for app user (nobody:nogroup)
chown -R nobody:nogroup /mnt/data /mnt/backups 2>/dev/null || true

exec /app/api
