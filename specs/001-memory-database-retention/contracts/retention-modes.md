# Contract: Data Retention Modes

## Workspace settings (portal)

### GET / PATCH `/api/v1/workspaces/:wid/settings`

**Retention field**: `data_retention`

| Value | Meaning |
| ----- | ------- |
| `memory` | Memory only — no message content in DB; request logs operational (`retained = false`) |
| `memory_database` | Memory + Database — inbox + request logs with `retained = true` |
| `provider` | Provider only — dispatch only; request logs operational (`retained = false`) |
| `provider_database` | Provider + Database — dispatch + request logs with `retained = true` |

**Legacy read aliases** (GET normalizes via `NormalizeRetentionValue`):

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

## Request log persistence rule (Idea 3)

### Insert (all modes)

`message_request_logs` INSERT on **successful send only**:

1. Outbound send completed without error; gateway responds with success.
2. Populate all metadata fields.
3. Set `retained = ShouldRetainRequestLog(message_dispatch_mode)`:
   - `memory_and_provider`, `provider_and_database` → `true`
   - `memory_only`, `provider_only` → `false`

Failed validation (4xx), auth errors, and dispatch failures MUST NOT insert rows.

### Read paths

| Consumer | Query |
| -------- | ----- |
| **Recent Requests** (portal inbox/logs) | All rows for workspace (existing `ListWithSource`; no `retained` filter) |
| **Retention export / compliance** (future) | `WHERE retained = true` |

## Channels

Rule applies uniformly to `email`, `sms`, `push`, and `chat` send endpoints.
