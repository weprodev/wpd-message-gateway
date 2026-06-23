# Data Model: Align Memory + Database Data Retention

## Workspace Setting

| Field | Value |
|-------|-------|
| Table | `workspace_settings` |
| Key | `data_retention` |
| Canonical values | `memory`, `memory_database`, `provider`, `provider_database` |

### Data fix (existing init migration — no new file)

Edit `database/migrations/20260318000000_init_gateway.up.sql` in place:

| Legacy value | Corrected to |
|--------------|--------------|
| `both` | `memory_database` |

Seed `memory_database` directly where workspace settings are inserted; avoid inserting `both`.

No rename of the `message_dispatch_mode` setting **key**. Value mappings in related files are updated (see FR-012).

### `message_dispatch_mode` value mapping (key unchanged)

| Dispatch mode value | Maps to `data_retention` |
|---------------------|--------------------------|
| `memory_only` | `memory` |
| `memory_and_database` | `memory_database` |
| `provider_only` | `provider` |
| `provider_and_database` | `provider_database` |

Removed value: `memory_and_provider` (replaced by `memory_and_database`).

### API values (no runtime aliasing)

Only canonical `data_retention` values are accepted on write. Legacy `both` is fixed by migration, not by runtime normalization of provider-named strings.

## Dispatch Mode (runtime)

| Dispatch mode | RAM inbox | stored_messages | Provider send |
|---------------|-----------|-----------------|---------------|
| `memory_only` | Yes | **No** | No |
| `memory_and_database` | Yes | Yes | **No** |
| `provider_only` | No* | No | Yes |
| `provider_and_database` | No* | Yes | Yes |

\*Memory integration fallback may write to inbox when active integration is the memory provider — unchanged existing behavior.

## Stored Message (memory_and_database)

| Field | Value on successful send |
|-------|--------------------------|
| `workspace_id` | Sender workspace |
| `channel_type` | `email` / `sms` / `push` / `chat` |
| `payload` | Full message JSON |
| `dispatch_status` | `sent` |
| `provider_message_id` | NULL |
| `provider_status_code` | NULL |
| `dispatch_error` | NULL |
| `dispatched_at` | UTC timestamp of successful capture |

On failure before both RAM and DB writes succeed: no `sent` status; row may not exist or remains uncommitted per transaction boundaries.

`RecordDispatchOutcome` is called with `Status = sent` and empty provider fields after successful inbox + archive writes.

## Stored Message (memory_only)

No row created in `stored_messages`. Sent API requests are **not** persisted to the database.

## API Response Metadata (memory_and_database)

| Meta key | Source |
|----------|--------|
| `dispatch_mode` | `memory_and_database` |
| `channel` | channel string |
| `inbox_message_id` | RAM capture ID |
| `stored_message_id` | `stored_messages.id` |
| `provider_name` | absent |
| `integration_id` | absent |

## Code naming (developer-facing)

| Removed | Replaced with |
|---------|---------------|
| `DispatchMemoryAndProvider` | `DispatchMemoryAndDatabase` |
| `memory_and_provider` | `memory_and_database` |
| `memory_provider` (if present) | `memory_database` |

Zero provider-named symbols for this policy after implementation.

## Entities unchanged

- `Integration` — not consulted for `memory_and_database` dispatch
- `message_dispatch_mode` setting key name unchanged; value mappings updated for Memory + Database
- Other retention modes — no schema changes
