# Implementation Plan: Message Dispatch Modes & Request Retention

**Branch**: `WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config` | **Date**: 2026-06-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config/spec.md`

## Summary

Add a fourth message dispatch mode (**Provider + Database**), introduce a `retained` boolean on `message_request_logs` to distinguish database-backed modes from operational-only logging, unify backend naming under `message_dispatch_mode`, and extend domain mapping for four canonical portal values (`memory`, `memory_database`, `provider`, `provider_database`). Operational logging flow stays unchanged — every request continues to insert a log row; only the `retained` flag differs by mode.

**Technical approach **: One append-only request-log table serves both **Recent Requests** (all rows) and long-term retention policy (`retained = true` filter). No insert gating by dispatch mode.

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript / React 19 (Portal UI)

**Primary Dependencies**: Echo (HTTP), PostgreSQL (`pgsql` client), Bruno E2E (`tests/bruno/`)

**Storage**: PostgreSQL — `workspace_settings` (key/value), `message_request_logs` (append-only audit), inbox tables for **Memory + Database** message content

**Testing**: `go test -race ./internal/...`, table-driven domain tests, Bruno (`bru run --env memory`), `make audit`

**Target Platform**: HTTP gateway server (`cmd/server`) + Portal SPA (`frontend/`)

**Project Type**: Dual-mode gateway — embedded SDK (`pkg/`) + HTTP server (`internal/`)

**Performance Goals**: No additional DB round-trips per send beyond existing settings read + log insert; `retained` derived in-process from resolved dispatch mode

**Constraints**:
- `pkg/*` must not import `internal/*`
- Update existing init migration in place (no new migration file for v1)
- Recent Requests must not filter by `retained`
- Provider + Database dispatch path must share Provider Only logic (DRY)

**Scale/Scope**: Four dispatch modes × four channels (email, SMS, push, chat); backend domain + handler + repository + Portal settings UI

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Status | Notes |
| ---- | ------ | ----- |
| Layer boundaries (Handler → Service → Port ← Repository) | Pass | Changes stay in existing layers; no cross-feature imports |
| `pkg/` isolation | Pass | Dispatch-mode logic remains in `internal/core/domain` |
| Verification chain (lint → smell → audit) | Pass | Required before merge |
| Docs sync with code | Pass | Update `docs/backend/architecture.md`, `portal-inbox.md`, `usage.md` when modes change |
| KISS / minimal diff | Pass | Reuse Provider Only path for Provider + Database; single `retained` column vs separate tables |

**Post-design re-check**: Pass — no new services or cross-cutting abstractions required.

## Project Structure

### Documentation (this feature)

```text
specs/WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config/
├── plan.md              # This file
├── spec.md
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1 validation guide
├── contracts/
│   └── retention-modes.md
└── checklists/
    └── requirements.md
```

(`tasks.md` generated via `/speckit-tasks`.)

### Source Code (repository root)

```text
internal/core/domain/
├── dispatch_mode.go           # Modes, parsing, retention mapping helpers
├── dispatch_mode_test.go
└── message_request_log.go     # Add Retained field

internal/core/service/
├── gateway_service.go         # Add provider_and_database case (memory_and_provider)
└── portal_service.go          # Normalize message_dispatch_mode on GET/PATCH

internal/infrastructure/repository/postgres/
└── message_request_log_repository.go   # INSERT/SELECT retained

internal/presentation/handler/
└── send_helper.go               # Set entry.Retained from domain helper

database/migrations/
└── 20260318000000_init_gateway.up.sql  # retained column on message_request_logs

frontend/src/features/settings/
├── settings.types.ts
├── settings.page.tsx
└── hooks/use-settings.hook.ts
```

**Structure Decision**: Web application layout — Go backend under `internal/`, React Portal under `frontend/`. All feature changes are localized to domain dispatch mode, gateway dispatch switch, request-log persistence, settings API normalization, and Portal retention panel.

## Phase Overview

### Phase 0 — Research (`research.md`)

Resolve naming (portal vs gateway dispatch strings), `retained` semantics, migration-in-place strategy, and settings key migration (`data_retention` → `message_dispatch_mode`).

### Phase 1 — Design

| Artifact | Purpose |
| -------- | ------- |
| `data-model.md` | Entities, fields, mode ↔ `retained` matrix |
| `contracts/retention-modes.md` | Portal/API settings contract + gateway dispatch mapping |
| `quickstart.md` | Runnable validation scenarios |

### Phase 2 — Implementation (deferred)

Generate `tasks.md` via `/speckit-tasks` when implementation starts. Expected workstreams:

1. **Schema** — `retained BOOLEAN NOT NULL DEFAULT false` on `message_request_logs`
2. **Domain** — four modes, legacy alias normalization, `ShouldRetainRequestLog(mode) bool`
3. **Gateway** — `DispatchProviderAndDatabase` (same body as `DispatchProviderOnly`); keep existing `DispatchMemoryAndProvider` / `memory_and_provider`
4. **Logging** — populate `Retained` in `send_helper.RecordLog` without changing insert conditions
5. **Portal API** — GET returns canonical `message_dispatch_mode`; PATCH accepts legacy keys/values on read, writes canonical only
6. **Portal UI** — four radio options; save `message_dispatch_mode` with canonical values
7. **Verification** — domain tests, Bruno, `make audit`

## Complexity Tracking

No constitution violations requiring justification.
