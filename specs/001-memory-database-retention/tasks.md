---

## description: "Task list for Align Memory + Database Data Retention"

# Tasks: Align Memory + Database Data Retention

**Input**: Design documents from `/specs/001-memory-database-retention/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/retention-modes.md, quickstart.md

**Tests**: Included — plan Phase D and success criteria SC-001–SC-006 require unit tests for new and regression behavior.

**Organization**: Tasks grouped by user story (US1–US3) per spec.md priorities.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: US1, US2, or US3 mapping to spec user stories

## Phase 1: Setup

**Purpose**: Confirm scope and validation path before code changes

- [x] T001 Review design artifacts (`spec.md`, `plan.md`, `data-model.md`, `contracts/retention-modes.md`, `quickstart.md`) in `specs/001-memory-database-retention/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Domain-layer rename and mappings — MUST complete before US1 dispatch logic

**⚠️ CRITICAL**: No user story implementation until this phase is complete

- [ ] T002 Rename `DispatchMemoryAndProvider` → `DispatchMemoryAndDatabase` (`memory_and_database`) and remove all `memory_and_provider` symbols in `internal/core/domain/dispatch_mode.go`
- [ ] T003 Update `retentionValueForMode`, `DataRetentionValueToDispatchMode`, `ParseMessageDispatchMode`, and `normalizeRetentionValue` in `internal/core/domain/dispatch_mode.go` (map `memory_database` ↔ `memory_and_database`; remove `both` / `memory_and_provider` runtime aliases)
- [ ] T004 Rename `RetentionProviders` (`"providers"`) → `RetentionProvider` (`"provider"`) and update all references in `internal/core/domain/dispatch_mode.go`
- [ ] T005 [P] Add table-driven mapping tests in `internal/core/domain/dispatch_mode_test.go` for `memory_and_database` ↔ `memory_database` and `provider_only` ↔ `provider`

**Checkpoint**: Domain constants and mappings compile; tests pass for `go test ./internal/core/domain/...`

---

## Phase 3: User Story 1 — Memory + Database captures without provider send (Priority: P1) 🎯 MVP

**Goal**: `memory_database` writes to RAM inbox + `stored_messages`, sets `dispatch_status` = `sent`, no provider invocation

**Independent Test**: Set `data_retention` = `memory_database`, send email — inbox row + `stored_messages` row with `sent`; zero provider calls (see `quickstart.md` §2–5)

### Implementation for User Story 1

- [ ] T006 [US1] Replace `DispatchMemoryAndProvider` branch with `captureMemoryAndDatabase` in `internal/core/service/gateway_service.go`: `writeToInbox()` → `writeToArchive()` → `RecordDispatchOutcome(status=sent)` → attach `inbox_message_id` + `stored_message_id` meta; no `sendViaProvider`
- [ ] T007 [US1] Add `TestGatewayService_SendEmail_memoryAndDatabase` in `internal/core/service/gateway_service_test.go` — asserts inbox + stored writes, `dispatch_mode` = `memory_and_database`, `stored_message_id` in meta, `lastOutcome.Status` = `sent`, no provider integration invoked
- [ ] T008 [US1] Add `TestGatewayService_SendEmail_memoryOnly_noStoredMessage` in `internal/core/service/gateway_service_test.go` — asserts RAM capture only, `storedMessages` stub not called

**Checkpoint**: `go test -race ./internal/core/service/...` passes US1 tests; `memory_database` send behavior matches spec FR-001, FR-002, FR-006, FR-007, FR-008

---

## Phase 4: User Story 2 — Consistent naming across Portal and backend (Priority: P1)

**Goal**: No `memory_and_provider` / `memory_provider` / `both` in code or stored settings; `message_dispatch_mode` key unchanged, values updated

**Independent Test**: `rg memory_and_provider memory_provider` in `internal/` returns zero matches; settings API returns `data_retention` = `memory_database` after Portal save

### Implementation for User Story 2

- [ ] T009 [US2] Edit `database/migrations/20260318000000_init_gateway.up.sql` in place — seed `memory_database` directly; correct any `both` values to `memory_database` (no new migration file)
- [ ] T010 [P] [US2] Update `DISPATCH_TO_RETENTION` in `frontend/src/features/settings/settings.api.ts`: `memory_and_database` → `memory_database`, `provider_only` → `provider`; remove `memory_and_provider` entry
- [ ] T011 [P] [US2] Update canonical `RetentionMode` from `"providers"` to `"provider"` in `frontend/src/features/settings/settings.types.ts` and `frontend/src/features/settings/pages/settings.page.tsx`
- [ ] T012 [P] [US2] Update retention mode tables and `memory_and_database` semantics in `docs/backend/architecture.md`
- [ ] T013 [P] [US2] Update retention mapping table and examples in `docs/backend/usage.md`
- [ ] T014 [P] [US2] Update inbox capture modes in `docs/backend/portal-inbox.md` (`memory_and_database` replaces `memory_and_provider`)
- [ ] T015 [US2] Update existing test fixtures from `RetentionProviders` to `RetentionProvider` in `internal/core/service/gateway_service_test.go` and any other `internal/` references

**Checkpoint**: `rg 'memory_and_provider|memory_provider|\\bboth\\b' internal/ frontend/src/features/settings/` returns zero misleading matches; init migration seeds `memory_database`

---

## Phase 5: User Story 3 — Other retention modes remain stable (Priority: P2)

**Goal**: `memory`, `provider`, and `provider_database` behavior unchanged

**Independent Test**: Existing and updated unit tests pass for all three modes; compare to pre-change baselines per `quickstart.md` §6

### Implementation for User Story 3

- [ ] T016 [US3] Verify and fix if needed existing `TestGatewayService_SendEmail_memoryOnly` and `TestGatewayService_SendEmail_providerOnly_memoryIntegration` in `internal/core/service/gateway_service_test.go` after domain rename
- [ ] T017 [US3] Verify and fix if needed existing `TestGatewayService_SendEmail_providerAndDatabase_`* tests in `internal/core/service/gateway_service_test.go` (unchanged provider + DB behavior)
- [ ] T018 [P] [US3] Run full service test suite: `go test -race ./internal/core/service/...`

**Checkpoint**: All gateway service tests green; no regression in provider-only or provider+database paths

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Verification chain and end-to-end validation

- [ ] T019 Run `golangci-lint run ./cmd/... ./internal/... ./pkg/...`
- [ ] T020 Run `/smell develop` and fix any BLOCKER/HIGH findings
- [ ] T021 Run `make audit` (full quality gate)
- [ ] T022 Validate scenarios in `specs/001-memory-database-retention/quickstart.md` manually or via Bruno if server available

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS all user stories**
- **US1 (Phase 3)**: Depends on Foundational (T002–T005)
- **US2 (Phase 4)**: Depends on Foundational; can run parallel with US1 after T002–T004 (docs/frontend/migration independent of gateway dispatch)
- **US3 (Phase 5)**: Depends on US1 + US2 code changes being in place
- **Polish (Phase 6)**: Depends on US1–US3 complete

### User Story Dependencies


| Story | Priority | Depends on   | Can parallel with     |
| ----- | -------- | ------------ | --------------------- |
| US1   | P1       | Foundational | US2 (after T002–T004) |
| US2   | P1       | Foundational | US1 (after T002–T004) |
| US3   | P2       | US1, US2     | —                     |


### Within Each User Story

- Domain mappings (Foundational) before dispatch logic (US1)
- Implementation before regression verification (US3)
- All stories before Polish verification chain

### Parallel Opportunities

- **T005** ∥ **T002–T004** only after T002 starts (same file — actually T005 is new file, parallel with T002-T004 if different people: T005 is dispatch_mode_test.go, T002-T004 are dispatch_mode.go — T005 depends on T002-T004 content). Fix: T005 depends on T002-T004.
- **T010, T011, T012, T013, T014** can run in parallel (different files) during US2
- **T012, T013, T014** parallel with **T006** if US2 docs work starts while US1 dispatch is implemented (different files)

---

## Parallel Example: User Story 2

```bash
# After T009 (migration), launch doc + frontend updates together:
Task T010: frontend/src/features/settings/settings.api.ts
Task T011: frontend/src/features/settings/settings.types.ts + settings.page.tsx
Task T012: docs/backend/architecture.md
Task T013: docs/backend/usage.md
Task T014: docs/backend/portal-inbox.md
```

---

## Parallel Example: User Story 1 + US2 docs

```bash
# After Foundational phase, two developers:
Developer A: T006–T008 (gateway_service.go + tests)
Developer B: T012–T014 (backend docs only)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002–T005)
3. Complete Phase 3: User Story 1 (T006–T008)
4. **STOP and VALIDATE**: Run US1 tests + quickstart §2–5
5. Deploy/demo if ready

### Incremental Delivery

1. Foundational → US1 (MVP: correct Memory + Database behavior)
2. US2 (naming, migration, frontend maps, docs)
3. US3 (regression confirmation)
4. Polish (lint → smell → audit)

### Suggested MVP Scope

**User Story 1 only** (T001–T008): Delivers the core fix — RAM + DB persistence, `sent` status, no provider. Naming/docs (US2) can follow in the same PR or immediately after.

---

## Notes

- Do **not** create a new migration file; edit `database/migrations/20260318000000_init_gateway.up.sql` only
- Do **not** rename the `message_dispatch_mode` setting **key**; update value mappings only
- `sent` for Memory + Database means capture succeeded, not provider delivery
- Portal UI labels unchanged; only `settings.api.ts` / `settings.types.ts` value strings update (`provider` not `providers`)
- Total tasks: **22** (T001–T022)

