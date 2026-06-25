# Implementation Plan: Data Retention Modes

**Branch**: `001-memory-database-retention` | **Date**: 2026-06-25 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-memory-database-retention/spec.md`

## Summary

Add a fourth data retention mode (**Provider + Database**) with dispatch parity to Provider Only. Separate **operational request logging** (Recent Requests) from **long-term retention** using a `retained` flag on `message_request_logs` (Idea 3): always insert on successful send; set `retained = true` only for **Memory + Database** and **Provider + Database**. Keep `DispatchMemoryAndProvider`; add `DispatchProviderAndDatabase`.

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript 5.x / React 19 (portal)

**Primary Dependencies**: Echo v4, PostgreSQL, existing design-system components

**Storage**: PostgreSQL `workspace_settings`, `message_request_logs` (+ `retained` column migration); in-process inbox for memory capture

**Testing**: Go table-driven tests (`dispatch_mode_test.go`, gateway/handler tests), Vitest for settings UI, Bruno E2E optional

**Target Platform**: WPD Message Gateway HTTP server + React portal

**Project Type**: Dual-mode gateway (embedded SDK + HTTP server); this feature touches server portal + `internal/` only

**Performance Goals**: One INSERT per successful send; one boolean check for `retained`

**Constraints**: Canonical retention values on write; legacy alias on GET; success-only inserts; Recent Requests unfiltered

**Scale/Scope**: 1 migration column, 4 retention enums, settings fourth radio, ~8 backend files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
| --------- | ------ | ----- |
| Layer boundaries (handler → service → port) | Pass | Handler sets `retained`; dispatch in service |
| `pkg/` does not import `internal/` | Pass | Changes stay in `internal/` + `frontend/` |
| Verification chain (lint → smell → audit) | Required post-impl | Per `docs/agents/verification.md` |
| No secrets in logs/UI persistence | Pass | Request logs store metadata only |
| KISS / DRY | Pass | One table + flag vs dual table |

**Post-design re-check**: Pass — Idea 3 reduces complexity vs insert gating.

## Project Structure

### Documentation (this feature)

```text
specs/001-memory-database-retention/
├── plan.md
├── research.md
├── data-model.md
├── contracts/retention-modes.md
├── quickstart.md
└── spec.md
```

### Source Code (repository root)

```text
database/migrations/
└── *_add_retained_to_message_request_logs.up.sql

internal/core/domain/
├── dispatch_mode.go          # ShouldRetainRequestLog (rename from ShouldPersistRequestLog)
├── message_request_log.go  # Retained field
└── dispatch_mode_test.go

internal/core/service/
└── gateway_service.go        # provider_and_database; keep memory_and_provider

internal/presentation/handler/
└── send_helper.go            # always RecordLog on success; pass retained flag

internal/infrastructure/repository/postgres/
└── message_request_log_repository.go  # INSERT retained column

internal/core/service/
└── portal_service.go         # GetSettings normalize; PatchSettings sync

frontend/src/features/settings/
├── settings.types.ts
├── pages/settings.page.tsx
└── hooks/use-settings.hook.ts
```

**Structure Decision**: Idea 3 — no second table; `retained` on existing log entity; Recent Requests query unchanged.

## Phase 0: Research — Complete

See [research.md](./research.md). Idea 3 adopted 2026-06-25 (operational vs retention split).

## Phase 1: Design — Complete

| Artifact | Path |
| -------- | ---- |
| Data model | [data-model.md](./data-model.md) |
| Contracts | [contracts/retention-modes.md](./contracts/retention-modes.md) |
| Quickstart | [quickstart.md](./quickstart.md) |

### Implementation phases (for `/speckit-tasks`)

**Phase A — Schema & domain**

1. Migration: add `retained BOOLEAN NOT NULL DEFAULT false` to `message_request_logs`
2. Domain: `ShouldRetainRequestLog`; rename/clarify helper vs insert gating
3. Update `MessageRequestLog` struct and repository INSERT
4. Table-driven mapping tests

**Phase B — Gateway & handler**

1. Keep `DispatchMemoryAndProvider`; add `DispatchProviderAndDatabase`
2. Extract shared provider-send path
3. `SendHelper`: always log successful sends; set `retained` from `ShouldRetainRequestLog` (remove insert gating)
4. Sync `data_retention` ↔ `message_dispatch_mode` on PATCH; normalize on GET

**Phase C — Portal UI**

1. Four retention radios including Provider + Database
2. Canonical `RetentionMode` types

**Phase D — Verification**

1. Assert Recent Requests shows sends for memory/provider modes (`retained = false`)
2. Assert `retained = true` only for database-backed modes
3. `go test -race`, frontend lint/test, `make audit`

## Complexity Tracking

No constitution violations requiring justification.
