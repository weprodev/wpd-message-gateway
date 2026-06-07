#!/bin/sh
# Runs golang-migrate then (on "up" only) idempotent seeds from database/seeds/*.sql in order.
# Shared by migration jobs. Set DATABASE_DIR, DATABASE_URL (or DB_*), and for seeds
# either DATABASE_URL + psql, or SQL_RUNNER with {} placeholder for the seed file path.
#
# Staging-only seeds: any seed file whose name contains "staging" is run only when
# ENVIRONMENT=staging. Set ENVIRONMENT in your deployment (e.g. Cloud Run, K8s job,
# or CI) so that staging runs get ENVIRONMENT=staging and production does not.
set -e

MIGRATE_CMD="${MIGRATE_CMD:-/migrate-tool}"
DATABASE_DIR="${DATABASE_DIR:-/database}"
# SQL_RUNNER: command to run a seed file; {} is replaced by the file path.
# In-container: psql "$DATABASE_URL" -f "{}"
# Local (Makefile): docker compose exec -T <db-container> psql ... -f - < "{}"

# 1. Construct DATABASE_URL if not set (for K8s compatibility)
if [ -z "$DATABASE_URL" ]; then
    echo "🔧 Constructing DATABASE_URL..."
    if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
        echo "⚠️  Missing DB env vars. Skipping DATABASE_URL."
    else
        DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require"
        export DATABASE_URL
        echo "✅ DATABASE_URL set."
    fi
fi

# 2. Prepare Arguments
ARGS="$@"
if echo "$ARGS" | grep -v -q "\-database"; then
    if [ -n "$DATABASE_URL" ]; then
        ARGS="-database $DATABASE_URL $ARGS"
    fi
fi
if echo "$ARGS" | grep -v -q "\-path"; then
    ARGS="-path $DATABASE_DIR/migrations $ARGS"
fi

# 3. Run migrate
echo "🔄 Running migrations..."
$MIGRATE_CMD $ARGS
echo "✅ Migrations applied."

# 4. Execute Seeds (only when 'up' was passed).
# Seeds run in lexicographic order; each SQL file must be idempotent (e.g. ON CONFLICT DO NOTHING).
if echo "$ARGS" | grep -q "up"; then
    RUN_SEEDS=false
    if [ ! -d "$DATABASE_DIR/seeds" ]; then
        echo "ℹ️  No seeds directory — skipping."
    elif [ -n "$SQL_RUNNER" ]; then
        RUN_SEEDS=true
    elif [ -n "$DATABASE_URL" ]; then
        RUN_SEEDS=true
    else
        echo "⚠️  DATABASE_URL / SQL_RUNNER not set — skipping seeds."
    fi

    if [ "$RUN_SEEDS" = true ]; then
        echo "🌱 Running Seeds..."
        SEED_COUNT=0
        for file in "$DATABASE_DIR"/seeds/*.sql; do
            [ -f "$file" ] || continue
            basename_file=$(basename "$file")
            case "$basename_file" in
                *staging*)
                    if [ "${ENVIRONMENT:-}" != "staging" ]; then
                        echo "   ⏭️  Skipping $basename_file (ENVIRONMENT is not staging)"
                        continue
                    fi
                    ;;
            esac
            SEED_COUNT=$((SEED_COUNT + 1))
            echo "   → $basename_file"
            if [ -n "$SQL_RUNNER" ]; then
                CMD=$(echo "$SQL_RUNNER" | sed "s|{}|$file|g")
                eval "$CMD"
            else
                psql "$DATABASE_URL" -f "$file"
            fi
        done
        if [ $SEED_COUNT -gt 0 ]; then
            echo "✅ Seeds applied."
        else
            echo "ℹ️  No seed files — skipping."
        fi
    fi
fi

echo "✅ Done."
