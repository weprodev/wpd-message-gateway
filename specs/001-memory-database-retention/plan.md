# Implementation Plan: Data Retention Modes

**Branch**: `001-memory-database-retention` | **Date**: 2026-06-24 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-memory-database-retention/spec.md`

## Summary

Add a fourth data retention mode (**Provider + Database**) that mirrors Provider Only dispatch but persists request metadata on successful sends only. Refactor **Memory only** and **Provider only** to skip all database writes (including `message_request_logs`). Gate request logging to **Memory + Database** and **Provider + Database** only, and only on the success path in `SendHelper`. Keep existing `DispatchMemoryAndProvider` (`memory_and_provider`) enum; add `DispatchProviderAndDatabase` only.

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript 5.x / React 19 (portal)

**Primary Dependencies**: Echo v4, PostgreSQL, existing design-system components

**Storage**: PostgreSQL `workspace_settings`, `message_request_logs` (gated inserts); in-process inbox for memory capture

**Testing**: Go table-driven tests (`dispatch_mode_test.go`, gateway/handler tests), Vitest for settings UI, Bruno E2E optional

**Target Platform**: WPD Message Gateway HTTP server + React portal

**Project Type**: Dual-mode gateway (embedded SDK + HTTP server); this feature touches server portal + `internal/` only

**Performance Goals**: No additional latency on send path beyond one enum check before log insert

**Constraints**: Canonical retention values on write; legacy alias support on read; success-only logging

**Scale/Scope**: 4 retention enums, settings page fourth radio, ~6 backend files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
| --------- | ------ | ----- |
| Layer boundaries (handler → service → port) | Pass | Log gating in handler; dispatch in service |
| `pkg/` does not import `internal/` | Pass | Changes stay in `internal/` + `frontend/` |
| Verification chain (lint → smell → audit) | Required post-impl | Per `docs/agents/verification.md` |
| No secrets in logs/UI persistence | Pass | Request logs store metadata only |
| KISS / DRY | Pass | Reuse provider dispatch path for `provider_and_database` |

**Post-design re-check**: Pass — no new packages or cross-layer violations.

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
internal/core/domain/
├── dispatch_mode.go          # enums, mapping, ShouldPersistRequestLog
└── dispatch_mode_test.go     # table-driven mapping tests

internal/core/service/
└── gateway_service.go        # provider_and_database case; keep memory_and_provider

internal/presentation/handler/
└── send_helper.go            # gate RecordLog by dispatch mode; success path only

internal/presentation/handler/ (or portal settings handler)
└── settings sync data_retention ↔ message_dispatch_mode on PATCH

frontend/src/features/settings/
├── settings.types.ts         # RetentionMode union (4 canonical values)
├── pages/settings.page.tsx     # four radios including Provider + Database
└── hooks/use-settings.hook.ts  # legacy alias normalization on load
```

**Structure Decision**: Domain mapping centralized in `dispatch_mode.go`; handler owns log gating and success-only rule.

## Phase 0: Research — Complete

See [research.md](./research.md). Success-only logging resolved via spec clarifications (2026-06-24).

## Phase 1: Design — Complete

| Artifact | Path |
| -------- | ---- |
| Data model | [data-model.md](./data-model.md) |
| Contracts | [contracts/retention-modes.md](./contracts/retention-modes.md) |
| Quickstart | [quickstart.md](./quickstart.md) |

### Implementation phases (for `/speckit-tasks`)

**Phase A — Domain & gateway (backend)**

1. Keep `DispatchMemoryAndProvider` (`memory_and_provider`); add `DispatchProviderAndDatabase`
2. Implement retention ↔ dispatch mapping helpers and legacy alias normalization
3. Add `ShouldPersistRequestLog`; gate `SendHelper.RecordLog` by mode and success path
4. Extract shared provider-send path for `provider_only` and `provider_and_database`
5. Sync `data_retention` and `message_dispatch_mode` on settings PATCH
6. Table-driven tests for all mappings

**Phase B — Portal retention UI**

1. Update `RetentionMode` type and four `RadioOption` entries (add Provider + Database)
2. Normalize legacy values on load in hook

**Phase C — Verification**

1. `go test -race` domain + service + handler
2. `cd frontend && npm run lint && npm run test`
3. Manual quickstart scenarios
4. `make audit`

## Complexity Tracking

No constitution violations requiring justification.
