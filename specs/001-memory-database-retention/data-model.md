# Data Model: Data Retention Modes

## Entities

### DataRetentionMode (portal / workspace_settings)

| Canonical value | UI label | Legacy alias (read only) |
| --------------- | -------- | ------------------------ |
| `memory` | Memory only | — |
| `memory_database` | Memory + Database | `both` |
| `provider` | Provider only | `providers` |
| `provider_database` | Provider + Database | — |

**Storage**: `workspace_settings.key = 'data_retention'`, `value` = canonical enum string.

**Validation**: Reject unknown values on PATCH; normalize aliases on GET.

### MessageDispatchMode (gateway runtime)

| Value | Maps from retention | Message capture | Provider dispatch | `retained` on success |
| ----- | ------------------- | --------------- | ----------------- | --------------------- |
| `memory_only` | `memory` | In-process inbox | No | `false` |
| `memory_and_provider` | `memory_database` | In-process inbox | Per integration | `true` |
| `provider_only` | `provider` | No | Yes | `false` |
| `provider_and_database` | `provider_database` | No | Yes (same as provider_only) | `true` |

**Storage**: `workspace_settings.key = 'message_dispatch_mode'` (synced when retention saved).

**Deprecated removed**: runtime alias `both` as a dispatch mode value (read-time alias for retention only).

### MessageRequestLog (schema change — Idea 3)

Table: `message_request_logs`

**Purpose split (single table)**:

| Concern | Rule |
| ------- | ---- |
| **Operational monitoring** (Recent Requests) | Insert on every **successful** send; all modes |
| **Long-term retention** | `retained = true` only for `memory_and_provider` and `provider_and_database` |

**New column**:

| Column | Type | Default | Meaning |
| ------ | ---- | ------- | ------- |
| `retained` | `BOOLEAN NOT NULL` | `false` | `true` = kept per data-retention policy; `false` = operational only |

**Existing fields** (unchanged): `workspace_id`, `api_key_id`, `channel_type`, `http_method`, `status_code`, `endpoint`, `provider_name`, `request_id`, `duration_ms`, `error_message`, `created_at`.

**Indexes** (recommended):

- Existing: `(workspace_id, created_at DESC)` — Recent Requests
- New: `(workspace_id, retained, created_at DESC) WHERE retained = true` — retention exports

**Lifecycle (optional v2)**: Purge rows where `retained = false` and `created_at < now() - operational_ttl`.

### API Key (unchanged)

No schema or handler changes for this feature.

## State transitions

### Retention policy change

```
[any mode] --user saves new mode--> [updated mode]
```

Effective for the next outbound request after settings persist. No retroactive update of `retained` on existing log rows.

## Domain helpers (`dispatch_mode.go`)

- `RetentionValueForMode(MessageDispatchMode) string`
- `DataRetentionValueToDispatchMode(string) (MessageDispatchMode, bool)`
- `ParseMessageDispatchMode(string) (MessageDispatchMode, bool)`
- `NormalizeRetentionValue(string) string` — applies legacy aliases
- `ShouldRetainRequestLog(MessageDispatchMode) bool` — sets `retained` column on insert (`true` for `memory_and_provider` and `provider_and_database`)

## Frontend types

```typescript
type RetentionMode = "memory" | "memory_database" | "provider" | "provider_database"
```

Legacy `both` / `providers` normalized on settings GET (backend).
