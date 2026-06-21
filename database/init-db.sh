#!/bin/bash
set -e

echo "=== Running Database Init Script ==="

echo "Applying migrations..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  -f /docker-entrypoint-initdb.d/migrations/20260318000000_init_gateway.up.sql

echo "Applying seeds..."
export POSTGRES_USER POSTGRES_DB
bash /docker-entrypoint-initdb.d/apply-seeds.sh /docker-entrypoint-initdb.d/seeds

echo "=== Database Init Script Completed ==="
