# Contract: Data Retention Modes

## Workspace settings (portal)

### GET / PATCH `/api/v1/workspaces/:wid/settings`

**Retention field**: `data_retention`

| Value | Meaning |
| ----- | ------- |
| `memory` | Memory only — no DB writes |
| `memory_database` | Memory + Database — inbox + request logs |
| `provider` | Provider only — dispatch only, no DB writes |
| `provider_database` | Provider + Database — dispatch + request logs |

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
| `memory_and_database` | `memory_database` |
| `provider_only` | `provider` |
| `provider_and_database` | `provider_database` |

**Removed values** (MUST NOT be written; parse rejects or maps on read once):

- `memory_and_provider`
- `both` (as dispatch mode)

## Request log persistence rule

`message_request_logs` INSERT occurs **if and only if** `message_dispatch_mode` ∈ `{ memory_and_database, provider_and_database }`.

Applies to success and error paths in `SendHelper.RecordLog`.

## API keys (unchanged)

| Method | Path | Body | Notes |
| ------ | ---- | ---- | ----- |
| POST | `/api/v1/workspaces/:wid/api-keys` | `{ "name": string }` | Response includes `client_secret` once |
| POST | `/api/v1/workspaces/:wid/api-keys/:keyId/regenerate` | — | Response `{ "client_secret": string }` |
| DELETE | `/api/v1/workspaces/:wid/api-keys/:keyId` | — | 204/200 per existing handler |

Frontend modals MUST NOT alter these contracts.

## UI modal contract (frontend-only)

### CreateApiKeyModal

| Element | Behavior |
| ------- | -------- |
| Prompt | "Please add API key name" |
| Input placeholder | e.g. "Production" |
| Close | **X** dismisses without create |
| Primary | "Generate Key" → `createApiKey(wid, name)` → open CredentialsModal |

### RegenerateApiKeyModal

| Element | Behavior |
| ------- | -------- |
| Title | "Are you sure you want to Regenerate the API Key?" |
| Close | No **X** |
| Cancel | Dismiss |
| Primary | "Generate Key" → `regenerateApiKey` → CredentialsModal |

### DeleteApiKeyModal

| Element | Behavior |
| ------- | -------- |
| Title | "Are you sure you want to Delete this API Key?" |
| Close | No **X** |
| Cancel | Dismiss |
| Primary | "Delete" (danger) → `deleteApiKey` → dismiss |

### CredentialsModal

| Element | Behavior |
| ------- | -------- |
| Warning | One-time display banner |
| Fields | API client ID, API secret (read-only) |
| Copy | Clipboard + icon transition `content_copy` → `check` |
