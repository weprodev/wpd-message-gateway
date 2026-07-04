#!/bin/bash
# Apply idempotent SQL seeds in order. Invoked by init-db.sh and `make seed-demo` only —
# must NOT live at the docker-entrypoint-initdb.d root (Postgres runs *.sh there alphabetically
# before init-db.sh, which breaks first-time migration).
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
