# Tasks: Data Retention Modes (Idea 3 — `retained` flag)

**Input**: Design documents from `/specs/001-memory-database-retention/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/retention-modes.md

## Format: `[ID] [P?] [Story] Description`

## Phase 1: Schema & Domain

- [ ] T001 Add migration `retained BOOLEAN NOT NULL DEFAULT false` to `message_request_logs` in `database/migrations/`
- [ ] T002 Add `Retained bool` to `MessageRequestLog` in `internal/core/domain/message_request_log.go`
- [ ] T003 Rename/clarify `ShouldPersistRequestLog` → `ShouldRetainRequestLog` in `internal/core/domain/dispatch_mode.go` (true for `memory_and_provider` and `provider_and_database` only)
- [ ] T004 [P] Update mapping tests in `internal/core/domain/dispatch_mode_test.go` for `ShouldRetainRequestLog`

## Phase 2: Repository & Handler

- [ ] T005 [US1] Include `retained` in INSERT in `internal/infrastructure/repository/postgres/message_request_log_repository.go`
- [ ] T006 [US1] Remove insert gating from `RecordLog`; always log successful sends; set `entry.Retained` from `ShouldRetainRequestLog` in `internal/presentation/handler/send_helper.go`
- [ ] T007 [US1] Update `send_helper_test.go` — all modes insert on success; assert `retained` true/false per mode
- [ ] T008 [US1] Add `provider_and_database` dispatch + shared provider path in `internal/core/service/gateway_service.go`
- [ ] T009 [US1] Sync retention on PATCH; normalize on GET in `internal/core/service/portal_service.go`

## Phase 3: Portal UI

- [ ] T010 [P] [US1] Four canonical retention values in `frontend/src/features/settings/settings.types.ts` and `settings.page.tsx`

## Phase 4: Verification

- [ ] T011 Run quickstart scenarios in `specs/001-memory-database-retention/quickstart.md` (Recent Requests + `retained` column)
- [ ] T012 Run verification chain per `docs/agents/verification.md`

## Notes

- Recent Requests (`ListWithSource`) — no `retained` filter.
- Do **not** gate INSERT by retention mode; gate **retained column** only.
- Optional TTL purge for `retained = false` — out of scope v1.
