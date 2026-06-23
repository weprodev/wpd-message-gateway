# Data Model: Data Retention Modes & API Key Modals

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

| Value | Maps from retention | Message capture | Provider dispatch | Request logs |
| ----- | ------------------- | --------------- | ----------------- | ------------ |
| `memory_only` | `memory` | In-process inbox | No | No |
| `memory_and_database` | `memory_database` | In-process inbox | Per integration | Yes |
| `provider_only` | `provider` | No | Yes | No |
| `provider_and_database` | `provider_database` | No | Yes (same as provider_only) | Yes |

**Storage**: `workspace_settings.key = 'message_dispatch_mode'` (synced when retention saved).

**Deprecated removed**: `memory_and_provider`, runtime alias `both` for dispatch mode.

### MessageRequestLog (unchanged schema)

Table: `message_request_logs`

Written **only** when dispatch mode is `memory_and_database` or `provider_and_database`.

Fields used: `workspace_id`, `api_key_id`, `channel_type`, `http_method`, `status_code`, `endpoint`, `provider_name`, `request_id`, `duration_ms`, `error_message`.

### API Key (unchanged)

No schema or handler changes. Modals consume existing REST responses:

- `POST /api/v1/workspaces/:wid/api-keys` → includes `client_secret` once
- `POST /api/v1/workspaces/:wid/api-keys/:keyId/regenerate` → `{ client_secret }`
- `DELETE /api/v1/workspaces/:wid/api-keys/:keyId`

## State transitions

### Retention policy change

```
[any mode] --user saves new mode--> [updated mode]
```

Effective for the next outbound request after settings persist. No retroactive deletion of existing logs or inbox messages.

### API key lifecycle (UI)

```
Generate key --> Create modal (name) --> API create --> Credentials modal --> dismiss
Regenerate --> Confirm modal --> API regenerate --> Credentials modal --> dismiss
Delete --> Confirm modal --> API delete --> list refresh
```

## Domain helpers (`dispatch_mode.go`)

- `retentionValueForMode(MessageDispatchMode) string`
- `DataRetentionValueToDispatchMode(string) (MessageDispatchMode, bool)`
- `ParseMessageDispatchMode(string) (MessageDispatchMode, bool)`
- `normalizeRetentionValue(string) string` — applies legacy aliases
- `ShouldPersistRequestLog(MessageDispatchMode) bool`

## Frontend types

```typescript
type RetentionMode = "memory" | "memory_database" | "provider" | "provider_database"
```

Replace `"both"` / `"providers"` in `settings.types.ts`.
