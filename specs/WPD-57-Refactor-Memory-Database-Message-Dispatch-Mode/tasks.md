# Tasks: Refactor Memory + Database Message Dispatch Mode

**Input**: Design documents from `/specs/WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/dispatch-mode-refactor.md, quickstart.md

**Tests**: Domain and gateway unit tests included where they lock dispatch path, naming, and `retained` behavior (not full TDD).

**Organization**: Tasks grouped by user story. `pkg/` changes require explicit approval per plan.md — out of scope unless naming parity is required.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

**Purpose**: Confirm environment and design alignment before code changes

- [x] T001 Review design artifacts in `specs/WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode/` (spec, plan, research, data-model, contracts, quickstart) and confirm no DB migration is required

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Domain rename, legacy aliases, and `retained` rules required by all user stories

**⚠️ CRITICAL**: No user story work until this phase is complete

- [x] T002 Rename `DispatchMemoryAndProvider` → `DispatchMemoryAndDatabase` with gateway string `memory_and_database` in `internal/core/domain/dispatch_mode.go`
- [x] T003 Add read-only legacy alias `memory_and_provider` → `memory_and_database` in `ParseMessageDispatchMode` and `SettingValueToDispatchMode` in `internal/core/domain/dispatch_mode.go`
- [x] T004 Update `ShouldRetainRequestLog` to return `true` for `memory_and_database` and `provider_and_database` only in `internal/core/domain/dispatch_mode.go`
- [x] T005 [P] Update table-driven tests in `internal/core/domain/dispatch_mode_test.go` for renamed constant, legacy aliases (`both`, `memory_and_provider`), and full `retained` matrix (Provider only `false`, Provider + Database `true`, Memory only `false`, Memory + Database `true`)

**Checkpoint**: Domain primitives ready — gateway and frontend work can proceed

---

## Phase 3: User Story 1 — Memory + Database Captures Locally Without Provider (Priority: P1) 🎯 MVP

**Goal**: **Memory + Database** uses inbox-only dispatch (same as **Memory only**) with `retained = true` on request logs; no provider lookup or send

**Independent Test**: Set `message_dispatch_mode` to `memory_database`, send on each channel — message in inbox, zero provider calls, `message_request_logs.retained = true`, `provider_name = 'memory'`, response meta `dispatch_mode: memory_and_database`

### Implementation for User Story 1

- [x] T006 [US1] Replace `case DispatchMemoryAndProvider` with `case DispatchMemoryAndDatabase` reusing `DispatchMemoryOnly` dispatch body (inbox capture only; no `activeIntegration`, no `sendViaProvider`) in `internal/core/service/gateway_service.go`
- [x] T007 [US1] Ensure send result meta stamps `dispatch_mode: memory_and_database` (not `memory_and_provider`) in `internal/core/service/gateway_service.go`
- [x] T008 [P] [US1] Rewrite memory+database gateway tests: assert inbox capture, zero provider sender invocations, and `provider_name = memory` in `internal/core/service/gateway_service_test.go`
- [x] T009 [P] [US1] Add test that **Memory + Database** succeeds with no connected integration in `internal/core/service/gateway_service_test.go`

**Checkpoint**: Core behavior change complete — inbox only + `retained = true` for Memory + Database

---

## Phase 4: User Story 2 — Consistent Naming Across Settings and Gateway (Priority: P1)

**Goal**: Portal and gateway use `memory_and_database` / `memory_database` vocabulary; legacy values resolve on read only; **Memory + Database** label describes memory capture with DB-backed log retention

**Independent Test**: Grep shows no canonical `memory_and_provider` or `DispatchMemoryAndProvider` in active code; settings round-trip persists `memory_database`; legacy `both` and `memory_and_provider` still resolve

### Implementation for User Story 2

- [x] T010 [P] [US2] Rename portal internal mode id `memory_and_provider` → `memory_and_database` in `frontend/src/features/settings/settings.types.ts`
- [x] T011 [P] [US2] Update portal ↔ API mappers and legacy read aliases in `frontend/src/features/settings/message-dispatch-mode.ts`
- [x] T012 [US2] Revise **Memory + Database** radio label and description (memory capture + request log retention, not provider send) in `frontend/src/features/settings/pages/settings.page.tsx`
- [x] T013 [P] [US2] Update legacy value display mapping (`both` → Memory + Database) in `frontend/src/features/settings/hooks/use-settings.hook.ts` if references remain
- [x] T014 [P] [US2] Update expectations for renamed mode id in `frontend/src/features/settings/message-dispatch-mode.test.ts`
- [x] T015 [US2] Grep codebase and remove remaining canonical `memory_and_provider` / `DispatchMemoryAndProvider` references outside intentional legacy read paths (exclude WPD-33 spec artifacts and this feature spec)

**Checkpoint**: Naming consistent across backend, frontend, and tests

---

## Phase 5: User Story 3 — Retained Request Logs for Database-Backed Modes (Priority: P2)

**Goal**: `retained = true` filter returns **Memory + Database** and **Provider + Database** rows only; **Recent Requests** lists all modes regardless of `retained`

**Independent Test**: Send under all four modes; SQL `WHERE retained = true` includes Memory + Database and Provider + Database only; Portal Recent Requests shows all sends

### Implementation for User Story 3

- [x] T016 [US3] Confirm `ShouldRetainForWorkspace` (or equivalent) delegates to updated `ShouldRetainRequestLog` in `internal/core/service/gateway_service.go`
- [x] T017 [US3] Confirm `PortalService` log list path does not filter by `retained` in `internal/core/service/portal_service.go` — add clarifying comment if already correct
- [x] T018 [P] [US3] Add or extend handler/service tests asserting `retained` false for Memory only and Provider only, `true` for Memory + Database and Provider + Database in `internal/presentation/handler/send_helper_test.go` (create if missing)

**Checkpoint**: Retention semantics match four-mode matrix; operational list unfiltered

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation sync, validation, and quality gate

- [x] T019 [P] Update dispatch mode matrix and rename table in `docs/backend/architecture.md`
- [x] T020 [P] Update inbox behavior for `memory_and_database` in `docs/backend/portal-inbox.md`
- [x] T021 [P] Update dispatch mode usage table in `docs/backend/usage.md`
- [ ] T022 Execute `specs/WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode/quickstart.md` scenarios 1–8 (manual or scripted) — skipped per user request
- [ ] T023 Run verification chain: `golangci-lint run ./cmd/... ./internal/... ./pkg/...`, `go test -race ./internal/core/domain/... ./internal/core/service/...`, `cd frontend && npm run lint && npm run test`, `make audit` — skipped per user request

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS** all user stories
- **User Stories (Phase 3–5)**: All depend on Foundational completion
  - US1 (P1) should complete before or in parallel with US2 naming (US2 frontend can start after T002–T004)
  - US3 (P2) depends on T004 (retention rules) and T006 (dispatch path)
- **Polish (Phase 6)**: Depends on US1–US3 completion

### User Story Dependencies

- **User Story 1 (P1)**: Requires Foundational (T002–T005); no dependency on US2/US3
- **User Story 2 (P1)**: Requires Foundational; can run in parallel with US1 after T002–T004
- **User Story 3 (P2)**: Requires T004 and T006; independently verifiable via SQL + Portal

### Within Each User Story

- Domain before gateway service
- Gateway service before gateway tests
- Frontend types before mappers before page labels
- Story checkpoint before Polish phase

### Parallel Opportunities

- T005, T008, T009 can run in parallel after T006 (different test files)
- T010, T011, T013, T014 can run in parallel (different frontend files)
- T019, T020, T021 can run in parallel (different doc files)
- US1 backend (T006–T009) and US2 frontend (T010–T014) can proceed in parallel once Foundational completes

---

## Parallel Example: User Story 2

```bash
# Launch frontend renames together after T002–T004:
Task: "Rename portal internal mode id in frontend/src/features/settings/settings.types.ts"
Task: "Update portal ↔ API mappers in frontend/src/features/settings/message-dispatch-mode.ts"
Task: "Update legacy display mapping in frontend/src/features/settings/hooks/use-settings.hook.ts"
Task: "Update expectations in frontend/src/features/settings/message-dispatch-mode.test.ts"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002–T005) — **CRITICAL**
3. Complete Phase 3: User Story 1 (T006–T009)
4. **STOP and VALIDATE**: quickstart.md Scenarios 1–2
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → domain ready
2. User Story 1 → inbox-only dispatch + `retained = true` (MVP)
3. User Story 2 → naming and Portal copy aligned
4. User Story 3 → retention filter and Recent Requests verified
5. Polish → docs + `make audit`

### Parallel Team Strategy

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (gateway)
   - Developer B: User Story 2 (frontend)
3. Developer A or B: User Story 3 verification
4. Either: Polish + audit

---

## Notes

- `retained` matrix: Provider only `false`, Provider + Database `true`, Memory only `false`, Memory + Database `true`
- **Memory only** dispatch and `retained = false` are unchanged — no code changes unless regression found
- `pkg/` untouched unless naming parity is required — propose with reason and wait for approval
- No Bruno collection updates expected (`tests/bruno/` has no `memory_and_provider` references)
- Legacy read aliases (`both`, `memory_and_provider`) must never be written on PATCH
- Commit after each task or logical group; stop at any checkpoint to validate independently
