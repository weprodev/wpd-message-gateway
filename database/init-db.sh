#!/bin/bash
set -e

echo "=== Running Database Init Script ==="

TABLE_EXISTS=$(psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -tAc \
  "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users');")

if [ "$TABLE_EXISTS" = "t" ]; then
  echo "Schema already present — skipping migration."
  exit 0
fi

echo "Applying migrations..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  -f /docker-entrypoint-initdb.d/migrations/20260318000000_init_gateway.up.sql

echo "=== Database Init Script Completed (seeds run on gateway startup) ==="
