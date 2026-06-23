# Quickstart: Validate Memory + Database Retention

## Prerequisites

- Gateway server running with PostgreSQL (`make dev` or equivalent)
- Workspace with Portal access
- Migrations applied (including `both` → `memory_database` normalization)

## 1. Set retention mode

**Portal**: Settings → Data retention → **Memory + Database** → Save.

**API**:

```bash
curl -X PATCH "$BASE/api/v1/workspaces/$WID/settings" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"data_retention":"memory_database"}'
```

**Verify**:

```bash
curl "$BASE/api/v1/workspaces/$WID/settings" -H "Authorization: Bearer $TOKEN"
# Expect: "data_retention": "memory_database"  (NOT "both")
```

## 2. Send a test message

```bash
curl -X POST "$BASE/api/v1/workspaces/$WID/messages/email" \
  -H "X-Client-Id: $CLIENT_ID" \
  -H "X-Client-Secret: $CLIENT_SECRET" \
  -H "Content-Type: application/json" \
  -d '{"to":["test@example.com"],"subject":"retention test","html":"<p>hi</p>"}'
```

**Expected response metadata**:

- `dispatch_mode` = `memory_and_database`
- `inbox_message_id` present
- `stored_message_id` present
- No `integration_id` / provider dispatch indicators

## 3. Verify RAM inbox (Portal)

Open Portal inbox for the workspace — message should appear (same as Memory only).

## 4. Verify durable storage

```sql
SELECT id, channel_type, dispatch_status, provider_message_id
FROM stored_messages
WHERE workspace_id = '$WID'
ORDER BY created_at DESC
LIMIT 1;
```

**Expected**: row exists, `dispatch_status` = `sent`, `dispatched_at` IS NOT NULL, `provider_message_id` IS NULL.

## 5b. Verify memory mode does NOT persist

Set `data_retention` = `memory`, send a message, then:

```sql
SELECT COUNT(*) FROM stored_messages WHERE workspace_id = '$WID';
```

Count must not increase from the memory-only send.

## 5. Verify no provider invocation

- No provider API logs for the send
- Mailgun/SMS/etc. integration not called (even if connected)

## 6. Regression checks

| Mode | Inbox | stored_messages | Provider |
|------|-------|-----------------|----------|
| `memory` | Yes | No | No |
| `memory_database` | Yes | Yes | No |
| `provider` | No* | No | Yes |
| `provider_database` | No* | Yes | Yes |

## 7. Legacy alias read

If DB still has `both` before migration:

```sql
SELECT value FROM workspace_settings WHERE workspace_id = '$WID' AND key = 'data_retention';
```

After migration + API read, response must normalize to `memory_database`.

## Automated tests

```bash
go test -race ./internal/core/domain/... ./internal/core/service/...
make audit
```
