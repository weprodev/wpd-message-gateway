#!/bin/bash
# Apply idempotent SQL seeds in order. Used by init-db.sh and `make seed-demo`.
set -euo pipefail

SEEDS_DIR="${1:-/docker-entrypoint-initdb.d/seeds}"
PSQL=(psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER:?}" --dbname "${POSTGRES_DB:?}")

for seed in \
  001_seed_permissions.sql \
  002_seed_providers.sql \
  003_seed_mailgun_config.sql \
  004_demo_workspace.sql
do
  echo "→ seeds/$seed"
  "${PSQL[@]}" -f "$SEEDS_DIR/$seed"
done
