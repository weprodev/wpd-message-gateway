# Quickstart: Validate Memory + Database Dispatch Refactor

**Feature**: WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode  
**Prerequisites**: Local Postgres, `configs/local.yml`, demo seed (`make seed-demo` or init-db)

## Setup

```bash
# From repo root
make init-db
make run
```

Demo credentials:
- Portal: `demo@weprodev.com` / `secret`
- API key: `demo-client-id` / `demo-secret`
- Workspace slug: `demo`

## Scenario 1 — Memory + Database: inbox only, retained log

1. PATCH dispatch mode:

```bash
curl -s -X PATCH "http://localhost:8080/api/v1/workspaces/<WID>/settings" \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"message_dispatch_mode":"memory_database"}'
```

2. Ensure a real email integration is connected (e.g. Mailgun stub from demo seed).
3. Send test email via API.
4. **Expected**:
   - Message appears in Portal inbox.
   - **No** outbound call to external provider (stub receives nothing / no provider error when integration exists).
   - `message_request_logs.retained = true` for the new row.
   - `provider_name = 'memory'` on the log row.
   - Response meta includes `dispatch_mode: memory_and_database`.

Verify:

```sql
SELECT retained, provider_name, status_code, channel_type
FROM message_request_logs
WHERE workspace_id = '<WID>'
ORDER BY created_at DESC LIMIT 1;
```

## Scenario 2 — Memory + Database without integration

1. Disconnect or remove active integration for a channel (or use workspace with none).
2. PATCH `message_dispatch_mode` → `memory_database`.
3. Send test message.
4. **Expected**: Success (same as Memory only) — dispatch does not require integration.

## Scenario 3 — Memory only unchanged

1. PATCH `message_dispatch_mode` → `memory`.
2. Send test message.
3. **Expected**: Inbox capture; `retained = false`.

## Scenario 4 — Provider + Database unchanged

1. PATCH `message_dispatch_mode` → `provider_database`.
2. Send test message.
3. **Expected**: Provider dispatch; `retained = true`; no inbox message.

## Scenario 5 — Legacy alias `both`

```sql
UPDATE workspace_settings
SET value = 'both'
WHERE workspace_id = '<WID>' AND key = 'message_dispatch_mode';
```

Send test message.

**Expected**: Same as Scenario 1 (inbox only, `retained = true`, gateway mode `memory_and_database`).

## Scenario 6 — Legacy gateway alias `memory_and_provider`

PATCH settings with gateway string directly (if portal allows) or insert:

```sql
UPDATE workspace_settings
SET value = 'memory_and_provider'
WHERE workspace_id = '<WID>' AND key = 'message_dispatch_mode';
```

Send test message.

**Expected**: Resolves to `memory_and_database` behavior; inbox only; `retained = true`.

## Scenario 7 — Recent Requests regression

For Scenarios 1, 3, and 4:

1. Open Portal **Recent Requests**.
2. **Expected**: All sends visible regardless of `retained`.

## Scenario 8 — Retention filter

```sql
SELECT COUNT(*) FROM message_request_logs
WHERE workspace_id = '<WID>' AND retained = true;
```

After Scenarios 1 and 4, count MUST include both Memory + Database and Provider + Database sends; MUST NOT include Memory only or Provider only sends.

## Automated Checks

```bash
go test -race ./internal/core/domain/... -run DispatchMode
go test -race ./internal/core/service/... -run Gateway
golangci-lint run ./cmd/... ./internal/... ./pkg/...
make audit
cd tests/bruno && bru run --env memory
```

## Pass Criteria

| Check | Pass |
| ----- | ---- |
| Memory + Database: inbox capture, no provider send | ☐ |
| Memory + Database: `retained = true` | ☐ |
| Memory only / Provider only: `retained = false` | ☐ |
| Provider + Database: unchanged (`retained = true`, provider send) | ☐ |
| Legacy `both` and `memory_and_provider` resolve correctly | ☐ |
| No canonical `memory_and_provider` in new code (grep) | ☐ |
| Recent Requests unfiltered by `retained` | ☐ |

See [spec.md](./spec.md) success criteria SC-001–SC-005.
