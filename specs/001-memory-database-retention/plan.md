# Implementation Plan: Align Memory + Database Data Retention

**Branch**: `001-memory-database-retention` | **Date**: 2026-06-22 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-memory-database-retention/spec.md`

## Summary

Align backend dispatch naming and behavior with the Portal's **Memory + Database** retention mode. Rename the internal dispatch mode from `memory_and_provider` to `memory_and_database`, normalize workspace `data_retention` values (`both` → `memory_database`), and change dispatch logic so messages are captured in RAM (Portal inbox) **and** persisted to `stored_messages` — without invoking any external provider. Frontend remains unchanged.

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript/React 19 (frontend — no changes)

**Primary Dependencies**: Echo HTTP server, PostgreSQL (pgx), existing `GatewayService` dispatch pipeline, `StoredMessageWriter` / `InboxWriter` ports

**Storage**: PostgreSQL `workspace_settings` (key `data_retention`), `stored_messages` (durable payloads), in-process RAM inbox

**Testing**: `go test -race ./internal/...`, Bruno E2E (`tests/bruno/`), `make audit`

**Target Platform**: WPD Message Gateway HTTP server mode (`cmd/server`)

**Project Type**: Dual-mode gateway (HTTP server + embedded SDK); change scoped to server dispatch path

**Performance Goals**: No additional provider latency for `memory_database` mode (provider call removed)

**Constraints**: `pkg/*` must not import `internal/*`; minimal diff; preserve other three retention modes; remove all `memory_and_provider` / `memory_provider` symbols; keep `message_dispatch_mode` key name, update its Memory + Database value mappings

**Scale/Scope**: ~10 files (domain, service, migration, docs, tests); no frontend changes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Status | Notes |
|------|--------|-------|
| Layer boundaries (Handler → Service → Port) | PASS | Changes stay in `domain` + `service` + migration |
| Test coverage for behavior change | PASS | New unit tests for `memory_and_database` dispatch |
| No secrets in logs/responses | PASS | Unchanged |
| Verification chain (lint → smell → audit) | PASS | Required before PR |
| KISS / minimal diff | PASS | Reuse existing `writeToInbox` + `writeToArchive` helpers |

**Post-design re-check**: PASS — no new abstractions; single new dispatch branch pattern.

## Project Structure

### Documentation (this feature)

```text
specs/001-memory-database-retention/
├── plan.md              # This file
├── research.md          # Phase 0 decisions
├── data-model.md        # Entity/value mappings
├── quickstart.md        # Validation guide
├── contracts/           # Retention mode contract
└── tasks.md             # (/speckit-tasks — not created here)
```

### Source Code (repository root)

```text
internal/core/domain/
├── dispatch_mode.go           # Rename constant, update mappings

internal/core/service/
├── gateway_service.go         # New memory_and_database dispatch branch
├── gateway_service_test.go    # Tests for memory_and_database

database/migrations/
├── 20260318000000_init_gateway.up.sql    # Edit in place (seeds/defaults for data_retention)
└── 20260318000000_init_gateway.down.sql  # Update only if up.sql schema changes require it

docs/backend/
├── architecture.md            # Update mode table
├── usage.md                   # Update retention mapping table
└── portal-inbox.md            # Update inbox capture modes
```

**Structure Decision**: Backend-only change in existing Clean Architecture layers. No new packages.

## Complexity Tracking

No constitution violations requiring justification.

## Implementation Phases

### Phase A — Domain rename (no provider-named symbols)

1. Rename `DispatchMemoryAndProvider` → `DispatchMemoryAndDatabase` (`memory_and_database`).
2. Remove `memory_and_provider` from `ParseMessageDispatchMode`, `normalizeRetentionValue`, and all mappings.
3. Update `retentionValueForMode`, `DataRetentionValueToDispatchMode`.
4. Update `message_dispatch_mode` value mappings in `dispatch_mode.go`, `portal_service.go`, `settings.api.ts`, and docs — key name unchanged; replace `memory_and_provider` → `memory_and_database` where Memory + Database is referenced.

### Phase B — Dispatch logic (`gateway_service.go`)

Replace removed `memory_and_provider` branch:

```text
memory_and_database:
  1. writeToInbox()   → fail closed on error
  2. writeToArchive() → fail closed on error
  3. RecordDispatchOutcome(status=sent, dispatched_at=now, no provider fields)
  4. attachMeta: dispatch_mode, channel, inbox_message_id, stored_message_id
  5. NO sendViaProvider
```

`memory_only` branch unchanged — no `writeToArchive()` call.

### Phase C — Existing migration (no new file)

Edit `database/migrations/20260318000000_init_gateway.up.sql` in place — do **not** add a separate migration file. Apply the `both` → `memory_database` fix wherever workspace settings are seeded or defaulted in that file:

```sql
-- Example: correct seed/default value directly, or normalize inline where rows are inserted
UPDATE workspace_settings
SET value = 'memory_database'
WHERE key = 'data_retention' AND value = 'both';
```

Prefer seeding `memory_database` from the start over inserting `both` and updating. No `.down.sql` change unless the up migration schema is modified.

### Phase D — Tests & docs

1. Add `TestGatewayService_SendEmail_memoryAndDatabase` — verifies inbox + stored writes, `dispatch_status` = `sent`, no provider.
2. Add `TestGatewayService_SendEmail_memoryOnly_noStoredMessage` — verifies no `stored_messages` row.
3. Add domain mapping tests for `memory_and_database` ↔ `memory_database`.
4. Update `message_dispatch_mode` value maps in `settings.api.ts`, `portal_service.go`, and backend docs (`architecture.md`, `usage.md`, `portal-inbox.md`).
5. Run verification chain.

### Phase E — Out of scope

- Portal UI label/copy changes (already correct)
- Renaming the `message_dispatch_mode` setting key itself
- New `stored_messages` dispatch status enum value
- SDK embedded mode changes
