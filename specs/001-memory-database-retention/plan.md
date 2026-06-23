# Implementation Plan: Data Retention Modes & API Key Modals

**Branch**: `001-memory-database-retention` | **Date**: 2025-06-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-memory-database-retention/spec.md`

## Summary

Add a fourth data retention mode (**Provider + Database**) that mirrors Provider Only dispatch but persists request metadata. Refactor **Memory only** and **Provider only** to skip all database writes (including `message_request_logs`). Gate request logging to **Memory + Database** and **Provider + Database** only. Replace browser dialogs on the Settings Developer tab with Create, Regenerate, Delete, and Credentials modals; API key backend endpoints remain unchanged.

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript 5.x / React 19 (portal)

**Primary Dependencies**: Echo v4, PostgreSQL, Radix Dialog, existing design-system components

**Storage**: PostgreSQL `workspace_settings`, `message_request_logs` (gated inserts); in-process inbox for memory capture

**Testing**: Go table-driven tests (`dispatch_mode_test.go`, gateway/handler tests), Vitest for modal components, Bruno E2E optional

**Target Platform**: WPD Message Gateway HTTP server + React portal

**Project Type**: Dual-mode gateway (embedded SDK + HTTP server); this feature touches server portal + `internal/` only

**Performance Goals**: No additional latency on send path beyond one enum check before log insert

**Constraints**: No API key handler/DB changes; canonical retention values on write; legacy alias support on read

**Scale/Scope**: 4 retention enums, 4 modal components, ~6 backend files, settings page refactor

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
| --------- | ------ | ----- |
| Layer boundaries (handler → service → port) | Pass | Log gating in handler; dispatch in service |
| `pkg/` does not import `internal/` | Pass | Changes stay in `internal/` + `frontend/` |
| Verification chain (lint → smell → audit) | Required post-impl | Per `docs/agents/verification.md` |
| No secrets in logs/UI persistence | Pass | Credentials modal clears state on dismiss |
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
└── gateway_service.go        # provider_and_database case, rename memory_and_database

internal/presentation/handler/
└── send_helper.go            # gate RecordLog by dispatch mode

internal/presentation/handler/ (or portal settings handler)
└── settings sync data_retention ↔ message_dispatch_mode on PATCH

frontend/src/features/settings/
├── settings.types.ts         # RetentionMode union
├── pages/settings.page.tsx     # four radios, wire modals
├── hooks/use-settings.hook.ts  # credentials state helpers
└── components/
    ├── create-api-key-modal/
    ├── regenerate-api-key-modal/
    ├── delete-api-key-modal/
    └── credentials-modal/
```

**Structure Decision**: Feature-local modal components under `frontend/src/features/settings/components/`; domain mapping centralized in `dispatch_mode.go`.

## Phase 0: Research — Complete

See [research.md](./research.md). All NEEDS CLARIFICATION items resolved via spec clarifications (2025-06-23).

## Phase 1: Design — Complete

| Artifact | Path |
| -------- | ---- |
| Data model | [data-model.md](./data-model.md) |
| Contracts | [contracts/retention-modes.md](./contracts/retention-modes.md) |
| Quickstart | [quickstart.md](./quickstart.md) |

### Implementation phases (for `/speckit-tasks`)

**Phase A — Domain & gateway (backend)**

1. Rename `DispatchMemoryAndProvider` → `DispatchMemoryAndDatabase`; add `DispatchProviderAndDatabase`
2. Implement retention ↔ dispatch mapping helpers and legacy alias normalization
3. Add `ShouldPersistRequestLog`; gate `SendHelper.RecordLog`
4. Extract shared provider-send path for `provider_only` and `provider_and_database`
5. Sync `data_retention` and `message_dispatch_mode` on settings PATCH
6. Table-driven tests for all mappings

**Phase B — Portal retention UI**

1. Update `RetentionMode` type and four `RadioOption` entries (add Provider + Database)
2. Normalize legacy values on load in hook

**Phase C — API key modals (frontend only)**

1. `CreateApiKeyModal` with name input + **X** close
2. `RegenerateApiKeyModal` / `DeleteApiKeyModal` without **X**
3. `CredentialsModal` with warning, copy buttons, icon transition
4. Remove `window.prompt` / `window.confirm` from `settings.page.tsx`
5. Component tests for modal interactions

**Phase D — Verification**

1. `go test -race` domain + service + handler
2. `cd frontend && npm run lint && npm run test`
3. Manual quickstart scenarios
4. `make audit`

## Complexity Tracking

No constitution violations requiring justification.
