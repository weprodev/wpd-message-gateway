# Contract: Data Retention Modes

## Workspace settings (portal)

### GET / PATCH `/api/v1/workspaces/:wid/settings`

**Retention field**: `data_retention`

| Value | Meaning |
| ----- | ------- |
| `memory` | Memory only — no DB writes |
| `memory_database` | Memory + Database — inbox + request logs (success only) |
| `provider` | Provider only — dispatch only, no DB writes |
| `provider_database` | Provider + Database — dispatch + request logs (success only) |

**Legacy read aliases** (GET may return stored legacy value; UI normalizes):

| Stored | Normalized |
| ------ | ---------- |
| `both` | `memory_database` |
| `providers` | `provider` |

**PATCH example**:

```json
{ "data_retention": "provider_database" }
```

Server MUST persist canonical value and sync `message_dispatch_mode` per mapping table in `data-model.md`.

## Gateway dispatch (internal)

**Setting key**: `message_dispatch_mode`

| Value | Retention source |
| ----- | ---------------- |
| `memory_only` | `memory` |
| `memory_and_provider` | `memory_database` |
| `provider_only` | `provider` |
| `provider_and_database` | `provider_database` |

**Removed values** (MUST NOT be written; parse rejects or maps on read once):

- `both` (as dispatch mode; maps to `memory_and_provider` on read)

## Request log persistence rule

`message_request_logs` INSERT occurs **if and only if**:

1. `message_dispatch_mode` ∈ `{ memory_and_provider, provider_and_database }`, **and**
2. The outbound send completed successfully (dispatch returned without error; gateway responds with success).

Failed validation (4xx), auth errors, and dispatch failures MUST NOT insert rows.

## Channels

Rule applies uniformly to `email`, `sms`, `push`, and `chat` send endpoints.
