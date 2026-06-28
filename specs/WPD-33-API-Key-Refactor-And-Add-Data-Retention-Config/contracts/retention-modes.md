# Contract: Message Dispatch Mode Settings

**Feature**: WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config  
**Version**: 1.0 (draft)  
**Applies to**: Portal REST API, workspace settings storage, gateway dispatch

## Settings API

### GET `/api/v1/workspaces/:wid/settings`

Response returns workspace settings as stored (pass-through):

```json
{
  "message_dispatch_mode": "memory_database",
  "owner_email": "admin@example.com"
}
```

**Normalization rules (GET)** — v1:
- Return keys and values from `workspace_settings` as persisted.
- Legacy **value** alias mapping (`both` → `memory_database`, etc.) on read is **deferred to a future refactor**.

### PATCH `/api/v1/workspaces/:wid/settings`

Request body (partial update):

```json
{
  "message_dispatch_mode": "provider_database"
}
```

**Write rules (PATCH)** — v1:
- Accept any setting key/value pair the client sends; persist each entry to `workspace_settings` as given.
- Dispatch mode MUST use the key `message_dispatch_mode` only.
- Canonical value set: `memory`, `memory_database`, `provider`, `provider_database` (legacy values `both`, `providers` accepted on read via gateway mapping only).
- Stricter validation and legacy value normalization on PATCH are **deferred to a future refactor**.

## Canonical Values

| UI label | Setting value | Gateway dispatch |
| -------- | ------------- | ---------------- |
| Memory only | `memory` | `memory_only` |
| Memory + Database | `memory_database` | `memory_and_provider` |
| Provider only | `provider` | `provider_only` |
| Provider + Database | `provider_database` | `provider_and_database` |

Portal UI may also persist gateway strings directly (`memory_only`, `provider_and_database`, etc.).

## Gateway Send Side Effects

For all modes, on each API send handled by `SendHelper.DispatchAndLog`:

1. Resolve dispatch mode from workspace settings (`message_dispatch_mode` key).
2. Map setting value via `SettingValueToDispatchMode` (handles portal aliases and gateway strings).
3. Execute dispatch per resolved mode (unchanged insert conditions for logs).
4. Insert `message_request_logs` row with `retained = ShouldRetainRequestLog(mode)`:

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

## SDK / Public Package

No contract change — `pkg/gateway` and `pkg/contracts` unaffected.
