#!/bin/bash
set -e

echo "=== Running Database Init Script ==="

# 1. Run migrations
echo "Applying migrations..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/migrations/20260318000000_init_gateway.up.sql

# 2. Run seeds in order: permissions -> mailgun config -> demo workspace
echo "Applying seeds..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/seeds/001_seed_permissions.sql
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/seeds/002_seed_mailgun_config.sql
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/seeds/003_demo_workspace.sql

echo "=== Database Init Script Completed ==="
