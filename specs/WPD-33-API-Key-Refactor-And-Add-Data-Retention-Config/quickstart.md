# Quickstart: Validate Message Dispatch Modes & Request Retention

**Feature**: WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config  
**Prerequisites**: Local Postgres, `configs/local.yml`, demo seed (`make seed-demo` or init-db)

## Setup

```bash
# From repo root
make init-db          # or equivalent: migrations + seeds
make run              # starts API + Portal
```

Demo credentials (from `database/seeds/004_demo_workspace.sql`):
- Portal: `demo@weprodev.com` / `secret`
- API key: `demo-client-id` / `demo-secret`
- Workspace slug: `demo`

## Scenario 1 — Default mode (Memory only)

1. Open Portal → **Settings → Data Retention**.
2. Confirm **Memory only** selected (or set via API):

```bash
curl -s -X PATCH "http://localhost:8080/api/v1/workspaces/<WID>/settings" \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"message_dispatch_mode":"memory"}'
```

3. Send test email via API.
4. **Expected**: Message in Portal inbox (memory); `message_request_logs.retained = false` for the new row.

Verify:

```sql
SELECT retained, status_code, channel_type
FROM message_request_logs
WHERE workspace_id = '<WID>'
ORDER BY created_at DESC LIMIT 1;
```

## Scenario 2 — Memory + Database

1. PATCH `message_dispatch_mode` → `memory_database`.
2. Send test message.
3. **Expected**: Message content persisted; log row `retained = true`.
4. Portal **Recent Requests** shows the entry.

## Scenario 3 — Provider only

1. Ensure email integration connected (demo seed includes Mailgun stub).
2. PATCH `message_dispatch_mode` → `provider`.
3. Send test message.
4. **Expected**: Provider dispatch path; no message content in DB; log `retained = false`.

## Scenario 4 — Provider + Database (new)

1. PATCH `message_dispatch_mode` → `provider_database`.
2. Send test message.
3. **Expected**:
   - Same HTTP/dispatch outcome as Scenario 3 (provider path).
   - Log row `retained = true`.
   - No message content persisted.

Compare Scenarios 3 and 4 side-by-side — dispatch metadata should match except `retained`.

## Scenario 5 — Legacy value alias read

1. Insert legacy **value** directly (one-time test):

```sql
INSERT INTO workspace_settings (workspace_id, key, value)
VALUES ('<WID>', 'message_dispatch_mode', 'both')
ON CONFLICT (workspace_id, key) DO UPDATE SET value = EXCLUDED.value;
```

2. GET settings via API.
3. **Expected**: Response contains `"message_dispatch_mode": "both"` (pass-through). Gateway resolves `both` → `memory_and_provider` at dispatch time via `SettingValueToDispatchMode`.

## Scenario 6 — Recent Requests regression

For each mode in Scenarios 1–4:

1. Open Portal message logs / Recent Requests.
2. **Expected**: Send appears in list regardless of `retained` value.

## Automated Checks

```bash
# Domain mapping tests
go test -race ./internal/core/domain/... -run DispatchMode

# Full gate (after implementation)
golangci-lint run ./cmd/... ./internal/... ./pkg/...
make audit

# Bruno E2E (when collections updated)
cd tests/bruno && bru run --env memory
```

## Pass Criteria

| Check | Pass |
| ----- | ---- |
| Four modes selectable in Portal | ☐ |
| `retained` true only for `both`/`memory_database` and `provider_database` | ☐ |
| Provider + Database dispatch matches Provider only | ☐ |
| Recent Requests works for all modes (no `retained` filter) | ☐ |
| List API / SQL can read `retained` column on each row | ☐ |
| Settings GET/PATCH use `message_dispatch_mode` key | ☐ |

See [spec.md](./spec.md) success criteria SC-001–SC-004.
