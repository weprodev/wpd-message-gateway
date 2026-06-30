# Contract: Memory + Database Dispatch Refactor

**Feature**: WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode  
**Version**: 1.0 (in progress)  
**Applies to**: Portal REST API, workspace settings, gateway dispatch, request logging  
**Supersedes**: WPD-33 `contracts/retention-modes.md` rows for `memory_database` / `memory_and_provider`

## Settings API

Unchanged key: `message_dispatch_mode`.

### GET `/api/v1/workspaces/:wid/settings`

Pass-through from `workspace_settings` (legacy values may appear as stored).

### PATCH `/api/v1/workspaces/:wid/settings`

```json
{
  "message_dispatch_mode": "memory_database"
}
```

**Write rules**:
- Canonical values: `memory`, `memory_database`, `provider`, `provider_database`
- MUST NOT write `both`, `memory_and_provider`, or `memory_provider`

## Canonical Values (updated)

| UI label | Setting value | Gateway dispatch | Legacy read aliases |
| -------- | ------------- | ---------------- | ------------------- |
| Memory only | `memory` | `memory_only` | — |
| Memory + Database | `memory_database` | `memory_and_database` | `both`, `memory_and_provider` |
| Provider only | `provider` | `provider_only` | `providers` |
| Provider + Database | `provider_database` | `provider_and_database` | — |

Portal internal mode id: `memory_and_database` (replaces `memory_and_provider`).

## Gateway Send Side Effects

For each API send via `SendHelper.DispatchAndLog`:

1. Resolve `message_dispatch_mode` from workspace settings.
2. Map via `SettingValueToDispatchMode` (includes legacy aliases).
3. Dispatch per resolved mode:

| Resolved gateway mode | Dispatch action | Integration required? |
| --------------------- | --------------- | --------------------- |
| `memory_only` | Inbox capture | No |
| `memory_and_database` | Inbox capture only (same as `memory_only`) | **No** |
| `provider_only` | Provider send | Yes (or memory fallback if integration is memory) |
| `provider_and_database` | Provider send (identical to `provider_only`) | Yes |

4. Insert `message_request_logs` with `retained = ShouldRetainRequestLog(mode)`:

| Resolved gateway mode | `retained` |
| --------------------- | ---------- |
| `memory_only` | `false` |
| `memory_and_database` | **`true`** |
| `provider_only` | `false` |
| `provider_and_database` | `true` |

## Response Metadata

Send result `meta.dispatch_mode` stamps canonical gateway string (`memory_and_database` for Memory + Database).

Legacy value `memory_and_provider` MUST NOT appear in new responses.

## Recent Requests API

Unchanged — returns all logs regardless of `retained`.

## Long-Term Retention Query

```sql
SELECT * FROM message_request_logs
WHERE workspace_id = $1 AND retained = true
ORDER BY created_at DESC;
```

Now includes rows from **Memory + Database** and **Provider + Database** workspaces.

## Breaking Changes

| Before (WPD-33) | After (WPD-57) |
| --------------- | -------------- |
| `memory_database` → provider + inbox dual send | `memory_database` → inbox only |
| `memory_and_provider` gateway string | `memory_and_database` |
| `retained = false` for Memory + Database | `retained = true` |
| Portal mode id `memory_and_provider` | `memory_and_database` |

**Compatibility**: Read aliases `both` and `memory_and_provider` continue to resolve to new behavior.

## SDK / Public Package

No contract change — `pkg/gateway` and `pkg/contracts` unaffected (out of scope v1).
