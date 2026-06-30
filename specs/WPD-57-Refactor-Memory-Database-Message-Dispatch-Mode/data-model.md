# Data Model: Refactor Memory + Database Message Dispatch Mode

**Feature**: WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode  
**Date**: 2026-06-29

## Entity Relationship (conceptual)

```text
Workspace
├── workspace_settings (1:N key/value)
│   └── message_dispatch_mode → memory | memory_database | provider | provider_database
├── integrations (unused by memory_and_database dispatch)
├── inbox (in-process RAM — memory_only and memory_and_database)
└── message_request_logs (1:N append-only)
        └── retained: boolean
```

No schema migration required — `retained` column exists from WPD-33.

## Workspace Setting — Message Dispatch Mode

| Field | Type | Constraints |
| ----- | ---- | ----------- |
| `workspace_id` | UUID | FK → workspaces |
| `key` | TEXT | `message_dispatch_mode` (canonical) |
| `value` | TEXT | One of: `memory`, `memory_database`, `provider`, `provider_database` |

**Legacy value normalization on read** (mapped via `SettingValueToDispatchMode`; not stored on new writes):

| Legacy value | Resolves to setting | Gateway mode |
| ------------ | ------------------- | ------------ |
| `both` | `memory_database` | `memory_and_database` |
| `providers` | `provider` | `provider_only` |
| `memory_and_provider` (gateway string in settings) | `memory_database` behavior | `memory_and_database` |

## Gateway Dispatch Mode (domain enum)

| Constant | String value | Maps from setting | Legacy read alias |
| -------- | ------------ | ----------------- | ----------------- |
| `DispatchMemoryOnly` | `memory_only` | `memory` | — |
| `DispatchMemoryAndDatabase` | `memory_and_database` | `memory_database` | `memory_and_provider` |
| `DispatchProviderOnly` | `provider_only` | `provider` | — |
| `DispatchProviderAndDatabase` | `provider_and_database` | `provider_database` | — |

**Default**: `memory_only` when setting missing or invalid.

## Message Request Log (`message_request_logs`)

Schema unchanged from WPD-33. Relevant column:

| Column | Type | Description |
| ------ | ---- | ----------- |
| `retained` | BOOLEAN NOT NULL DEFAULT false | `true` for database-backed retention modes |

## Mode Behavior Matrix (updated)

| Setting value | Gateway mode | Inbox capture | Provider dispatch | `provider_name` on log | `retained` |
| ------------- | ------------ | ------------- | ----------------- | ---------------------- | ---------- |
| `memory` | `memory_only` | Yes | No | `memory` | `false` |
| `memory_database` | `memory_and_database` | Yes | **No** | `memory` | **`true`** |
| `provider` | `provider_only` | No | Yes | integration name | `false` |
| `provider_database` | `provider_and_database` | No | Yes | integration name | `true` |

### Delta from WPD-33

| Field | WPD-33 (`memory_database`) | WPD-57 (`memory_database`) |
| ----- | -------------------------- | ---------------------------- |
| Gateway mode string | `memory_and_provider` | `memory_and_database` |
| Provider dispatch | Yes | **No** |
| `retained` | `false` | **`true`** |
| Message content in DB | Documented as yes | **No** — in-process inbox only (same as memory only) |

## State Transitions

**Setting change**: Takes effect on next outbound request. Existing log rows keep their original `retained` value.

**Invalid setting value**: Fall back to `memory_only`; log with `retained = false`.

**Rename transition**: Historical log rows may contain `dispatch_mode: memory_and_provider` in response metadata; new sends stamp `memory_and_database`.

## Validation Rules

- `ParseMessageDispatchMode` accepts four canonical gateway strings plus legacy `memory_and_provider` (maps to `memory_and_database`).
- `ShouldRetainRequestLog(mode)` returns `true` for `memory_and_database` and `provider_and_database`.
- Portal PATCH rejects unknown setting values for `message_dispatch_mode`.
- Recent Requests queries MUST NOT filter by `retained`.
- Long-term export queries MUST use `WHERE retained = true`.

## Out of Scope (data model)

- TTL purge job for `retained = false`
- Persisting inbox message bodies to PostgreSQL
- Storing dispatch mode snapshot on each log row
