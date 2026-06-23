# Contract: Data Retention Modes

## Portal / API setting

**Key**: `data_retention`

**Canonical enum**:

| Value | Label | Description |
|-------|-------|-------------|
| `memory` | Memory only | RAM inbox capture only; **no** `stored_messages` row |
| `memory_database` | Memory + Database | RAM inbox + durable `stored_messages` with `dispatch_status` = `sent`; no provider |
| `provider` | Providers only | Provider dispatch only |
| `provider_database` | Provider + Database | Provider dispatch + durable `stored_messages` |

## Database (existing init migration — no new file)

Edit `database/migrations/20260318000000_init_gateway.up.sql` in place. Seed `memory_database` directly; correct any `both` values to `memory_database`.

## `message_dispatch_mode` (legacy key — name unchanged)

The setting key `message_dispatch_mode` is **not** renamed. Files that read or map this key are updated so Memory + Database uses `memory_and_database` (not `memory_and_provider`).

| `message_dispatch_mode` value | `data_retention` |
|-------------------------------|------------------|
| `memory_only` | `memory` |
| `memory_and_database` | `memory_database` |
| `provider_only` | `provider` |
| `provider_and_database` | `provider_database` |

## Dispatch mode (internal / response metadata)

| `data_retention` | `dispatch_mode` meta |
|------------------|----------------------|
| `memory` | `memory_only` |
| `memory_database` | `memory_and_database` |
| `provider` | `provider_only` |
| `provider_database` | `provider_and_database` |

Removed: `memory_and_provider` — must not appear in code, logs, or API responses.

## Send result metadata by mode

### `memory_database`

```json
{
  "id": "<inbox_message_id>",
  "meta": {
    "dispatch_mode": "memory_and_database",
    "channel": "email",
    "inbox_message_id": "<uuid>",
    "stored_message_id": "<uuid>"
  }
}
```

Stored row: `dispatch_status` = `sent`, `dispatched_at` set, provider fields null.

### `memory` (no database persistence)

```json
{
  "id": "<inbox_message_id>",
  "meta": {
    "dispatch_mode": "memory_only",
    "channel": "email"
  }
}
```

No `stored_message_id`. No `stored_messages` row created.

### `provider_database` (unchanged reference)

```json
{
  "id": "<provider_message_id>",
  "meta": {
    "dispatch_mode": "provider_and_database",
    "channel": "email",
    "integration_id": "<uuid>",
    "provider_name": "mailgun",
    "stored_message_id": "<uuid>"
  }
}
```

## Breaking change notice

- `dispatch_mode` meta changes from `memory_and_provider` to `memory_and_database`.
- Behavior change: `memory_database` no longer invokes providers; persists to `stored_messages` instead.
- `memory` mode explicitly does not write to `stored_messages`.
