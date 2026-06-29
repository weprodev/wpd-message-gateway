# Data Model: Message Dispatch Modes & Request Retention

**Feature**: WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config  
**Date**: 2026-06-28

## Entity Relationship (conceptual)

```text
Workspace
├── workspace_settings (1:N key/value)
│   └── message_dispatch_mode → memory | memory_database | provider | provider_database
├── integrations (provider credentials)
├── inbox / message content (Memory + Database only)
└── message_request_logs (1:N append-only)
        └── retained: boolean
```

## Workspace Setting — Message Dispatch Mode

| Field | Type | Constraints |
| ----- | ---- | ----------- |
| `workspace_id` | UUID | FK → workspaces |
| `key` | TEXT | `message_dispatch_mode` (canonical) |
| `value` | TEXT | One of: `memory`, `memory_database`, `provider`, `provider_database` |

**Legacy value normalization on read** (not stored on new writes; mapped via `SettingValueToDispatchMode`):

| Legacy value | Canonical value |
| ------------ | --------------- |
| `both` | `memory_database` |
| `providers` | `provider` |

## Gateway Dispatch Mode (domain enum)

Internal runtime enum `MessageDispatchMode` (stored in gateway meta / logs, not in `workspace_settings`):

| Constant | String value | Maps from setting |
| -------- | ------------ | ----------------- |
| `DispatchMemoryOnly` | `memory_only` | `memory` |
| `DispatchMemoryAndProvider` | `memory_and_provider` | `memory_database` |
| `DispatchProviderOnly` | `provider_only` | `provider` |
| `DispatchProviderAndDatabase` | `provider_and_database` | `provider_database` |

**Default**: `memory_only` when setting missing or invalid.

## Message Request Log (`message_request_logs`)

Existing columns unchanged except new field:

| Column | Type | Nullable | Description |
| ------ | ---- | -------- | ----------- |
| `id` | UUID | NO | PK |
| `workspace_id` | UUID | NO | FK |
| `api_key_id` | UUID | YES | FK |
| `channel_type` | TEXT | NO | email \| sms \| push \| chat |
| `provider_name` | VARCHAR(64) | YES | |
| `http_method` | VARCHAR(16) | NO | |
| `status_code` | SMALLINT | NO | |
| `endpoint` | VARCHAR(512) | NO | |
| `request_id` | VARCHAR(64) | YES | |
| `duration_ms` | INT | YES | ≥ 0 |
| `error_message` | TEXT | YES | |
| **`retained`** | **BOOLEAN** | **NO** | **DEFAULT false** — see matrix below |
| `created_at` | TIMESTAMPTZ | NO | DEFAULT NOW() |

### Domain struct addition

```go
type MessageRequestLog struct {
    // ... existing fields ...
    Retained bool `json:"retained"`
}
```

## Mode Behavior Matrix

| Setting value | Gateway mode | Message content in DB | Provider dispatch | `retained` on log insert |
| ------------- | ------------ | ----------------------- | ----------------- | ------------------------ |
| `memory` | `memory_only` | No (in-process inbox only) | No | `false` |
| `memory_database` | `memory_and_provider` | Yes (inbox/DB path) | Yes (if integration connected) | `false` |
| `provider` | `provider_only` | No | Yes | `false` |
| `provider_database` | `provider_and_database` | No | Yes (identical to provider) | `true` |

## State Transitions

**Setting change**: Takes effect on next outbound request. Mode read at dispatch time; existing log rows keep their original `retained` value.

**Invalid setting value**: Fall back to `memory` / `memory_only`; log with `retained = false`.

## Validation Rules

- `ParseMessageDispatchMode` accepts only the four gateway string values.
- Portal PATCH rejects unknown setting values for `message_dispatch_mode`.
- `retained` is never NULL on insert; default `false` at DB level for safety.
- Recent Requests queries MUST NOT add `WHERE retained = true`.
- Long-term export queries MUST use `WHERE retained = true`.

## Out of Scope (data model)

- TTL purge job for `retained = false`
- Storing dispatch mode on each log row (derivable from `retained` + workspace setting at send time; add column later if audit requires point-in-time mode snapshot)
