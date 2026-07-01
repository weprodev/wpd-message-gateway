# Implementation Plan: Refactor Memory + Database Message Dispatch Mode

**Branch**: `WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode` | **Date**: 2026-06-29 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode/spec.md`

## Summary

Correct **Memory + Database** dispatch so it matches **Memory only** (inbox capture, no provider call) while marking request logs `retained = true` (same retention semantics as **Provider + Database**). Rename misleading identifiers: gateway mode `memory_and_provider` → `memory_and_database`, constant `DispatchMemoryAndProvider` → `DispatchMemoryAndDatabase`, portal internal id `memory_and_provider` → `memory_and_database`. Legacy strings (`both`, `memory_and_provider`) remain read-only aliases.

**Technical approach**: Replace the `DispatchMemoryAndProvider` gateway branch (memory + provider dual path) with a `DispatchMemoryAndDatabase` branch that reuses the `DispatchMemoryOnly` dispatch body; extend `ShouldRetainRequestLog` to include `memory_and_database`; mechanical rename across domain, service, frontend, docs, and Bruno.

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript / React 19 (Portal UI)

**Primary Dependencies**: Echo (HTTP), PostgreSQL, Bruno E2E (`tests/bruno/`)

**Storage**: PostgreSQL — `workspace_settings` (`message_dispatch_mode`), `message_request_logs` (`retained` column exists from WPD-33); inbox remains in-process RAM

**Testing**: `go test -race ./internal/core/domain/... ./internal/core/service/...`, Bruno (`bru run --env memory`), `make audit`

**Target Platform**: HTTP gateway server (`cmd/server`) + Portal SPA (`frontend/`)

**Project Type**: Dual-mode gateway — embedded SDK (`pkg/`) + HTTP server (`internal/`)

**Performance Goals**: **Memory + Database** removes provider lookup and HTTP send — strictly fewer operations than current `memory_and_provider` path

**Constraints**:
- `pkg/*` must not import `internal/*`
- No new DB migration — schema unchanged; behavior + naming only
- Recent Requests must not filter by `retained`
- Legacy read aliases (`both`, `memory_and_provider`) must not break existing workspaces or log metadata

**Scale/Scope**: One dispatch-mode refactor × four channels; ~12 source/doc files; `retained` by mode — **Provider only** `false`, **Provider + Database** `true`, **Memory only** `false`, **Memory + Database** `true` (supersedes WPD-33 `retained = false` for `memory_database` only)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Status | Notes |
| ---- | ------ | ----- |
| Layer boundaries (Handler → Service → Port ← Repository) | Pass | Dispatch change in `gateway_service`; retention in `dispatch_mode` domain |
| `pkg/` isolation | Pass | Change `pkg/` only if naming must be updated; otherwise propose with a valid reason and wait for approval before editing |
| Verification chain (lint → smell → audit) | Pass | Required before merge |
| Docs sync with code | Pass | Update `architecture.md`, `portal-inbox.md`, `usage.md` |
| KISS / minimal diff | Pass | Reuse `memory_only` dispatch body for **Memory + Database**; rename + set `retained = true` for that mode only — **Memory only** stays `retained = false` (no change) |

**Post-design re-check**: Pass — no new abstractions; delete provider branch from memory+database case.

## Project Structure

### Documentation (this feature)

```text
specs/WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode/
├── plan.md              # This file
├── spec.md
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1 validation guide
├── contracts/
│   └── dispatch-mode-refactor.md
└── checklists/
    └── requirements.md
```

(`tasks.md` generated via `/speckit-tasks` when implementation starts.)

### Source Code (repository root)

```text
internal/core/domain/
├── dispatch_mode.go           # Rename constant; legacy alias; ShouldRetainRequestLog
└── dispatch_mode_test.go

internal/core/service/
├── gateway_service.go         # Replace memory_and_provider case → memory_and_database (memory_only body)
└── gateway_service_test.go

frontend/src/features/settings/
├── settings.types.ts          # memory_and_provider → memory_and_database
├── settings.page.tsx            # Radio ids, labels, dispatch mode state
└── hooks/use-settings.hook.ts   # Portal ↔ API value mapping

docs/backend/
├── architecture.md
├── portal-inbox.md
└── usage.md

tests/bruno/                   # Update env/mode references if present
```

**Structure Decision**: Web application layout — all changes localized to domain dispatch mode, gateway dispatch switch, Portal settings UI, and docs. No repository or migration changes.

## Phase Overview

### Phase 0 — Research (`research.md`)

Resolve rename matrix, legacy alias strategy, dispatch branch reuse, and WPD-33 supersession scope.

### Phase 1 — Design

| Artifact | Purpose |
| -------- | ------- |
| `data-model.md` | Updated mode ↔ `retained` matrix; rename table |
| `contracts/dispatch-mode-refactor.md` | API/gateway contract deltas |
| `quickstart.md` | Validation scenarios for new behavior |

### Phase 2 — Implementation (deferred)

Generate `tasks.md` via `/speckit-tasks`. Expected workstreams:

1. **Domain** — `DispatchMemoryAndDatabase` / `memory_and_database`; `ParseMessageDispatchMode` accepts legacy `memory_and_provider` on read; `ShouldRetainRequestLog` true for `memory_and_database` + `provider_and_database`
2. **Gateway** — Replace `case DispatchMemoryAndProvider` with `DispatchMemoryAndDatabase` using `memory_only` logic (no `activeIntegration`, no `sendViaProvider`)
3. **Frontend** — Rename portal mode id; update settings hook mappers; revise **Memory + Database** label copy
4. **Tests** — Update table-driven domain/service tests; assert zero provider calls under `memory_database`
5. **Docs** — Sync architecture, inbox, usage tables
6. **Verification** — `golangci-lint`, `go test -race`, Bruno, `make audit`

## Complexity Tracking

No constitution violations requiring justification.
