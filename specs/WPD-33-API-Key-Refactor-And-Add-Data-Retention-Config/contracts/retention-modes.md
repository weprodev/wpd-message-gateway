# Contract: Message Dispatch Mode Settings

**Feature**: WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config  
**Version**: 1.0 (draft)  
**Applies to**: Portal REST API, workspace settings storage, gateway dispatch

## Settings API

### GET `/api/v1/workspaces/:wid/settings`

Response returns workspace settings as stored (pass-through). Either key may appear depending on what was saved:

```json
{
  "data_retention": "both",
  "owner_email": "admin@example.com"
}
```

or:

```json
{
  "message_dispatch_mode": "memory_database",
  "owner_email": "admin@example.com"
}
```

**Normalization rules (GET)** — v1:
- No key or value migration on read. Return keys and values from `workspace_settings` as persisted.
- Legacy value alias mapping (`both` → `memory_database`, etc.) and `data_retention` → `message_dispatch_mode` consolidation are **deferred to a future refactor**.

### PATCH `/api/v1/workspaces/:wid/settings`

Request body (partial update). Either dispatch key is accepted:

```json
{
  "data_retention": "provider_database"
}
```

or:

```json
{
  "message_dispatch_mode": "provider_database"
}
```

**Write rules (PATCH)** — v1:
- Accept any setting key/value pair the client sends; persist each entry to `workspace_settings` as given.
- Both `data_retention` and `message_dispatch_mode` are valid for this feature — no rejection or silent ignore of either key.
- Canonical value set for new **Provider + Database** mode: `provider_database` (plus existing portal values `memory`, `both`/`memory_database`, `providers`/`provider` as used today).
- Stricter validation, single canonical key, and legacy alias normalization are **deferred to a future refactor**.

## Canonical Values

| UI label | Setting value | Gateway dispatch |
| -------- | ------------- | ---------------- |
| Memory only | `memory` | `memory_only` |
| Memory + Database | `memory_database` | `memory_and_provider` |
| Provider only | `provider` | `provider_only` |
| Provider + Database | `provider_database` | `provider_and_database` |

## Gateway Send Side Effects

For all modes, on each API send handled by `SendHelper.DispatchAndLog`:

1. Resolve dispatch mode from workspace settings (gateway reads `message_dispatch_mode`; portal may still persist `data_retention` until a future settings refactor).
2. Execute dispatch per mode (unchanged insert conditions for logs).
3. Insert `message_request_logs` row with `retained = ShouldRetainRequestLog(mode)`:

| Resolved gateway mode | `retained` |
| --------------------- | ---------- |
| `memory_only` | `false` |
| `memory_and_provider` | `true` |
| `provider_only` | `false` |
| `provider_and_database` | `true` |

## Recent Requests API

Existing list endpoint behavior unchanged — returns all logs for workspace regardless of `retained`.

Optional future response field (non-breaking):

```json
{
  "id": "...",
  "retained": true,
  "channel_type": "email"
}
```

Not required for v1 acceptance if column is queryable internally.

## Long-Term Retention Query (internal / future export)

```sql
SELECT * FROM message_request_logs
WHERE workspace_id = $1 AND retained = true
ORDER BY created_at DESC;
```

## Breaking Changes

| Before | After |
| ------ | ----- |
| Three portal retention options | Four options including **Provider + Database** (`provider_database`) |
| No `retained` column | Required on all new inserts |
| No `provider_and_database` gateway mode | New mode available |

**Deferred (future refactor)**: consolidating `data_retention` → `message_dispatch_mode`, GET normalization, and PATCH validation.

## SDK / Public Package

No contract change — `pkg/gateway` and `pkg/contracts` unaffected.
