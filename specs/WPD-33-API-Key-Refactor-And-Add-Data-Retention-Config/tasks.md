# Tasks: Message Dispatch Modes & Request Retention

**Input**: Design documents from `/specs/WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/retention-modes.md, quickstart.md

**Tests**: Domain and handler unit tests included where they lock mapping and `retained` behavior (not full TDD).

**Organization**: Tasks grouped by user story. Settings key consolidation (`data_retention` → `message_dispatch_mode` on GET/PATCH) is **out of scope** — deferred per contracts/retention-modes.md.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup

**Purpose**: Confirm environment and design alignment before code changes

- [x] T001 Review design artifacts in `specs/WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config/` (spec, plan, research, data-model, contracts)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema and domain primitives required by all user stories

**⚠️ CRITICAL**: No user story work until this phase is complete

- [x] T002 Add `retained BOOLEAN NOT NULL DEFAULT false` to `message_request_logs` in `database/migrations/20260318000000_init_gateway.up.sql` and mirror in `database/migrations/20260318000000_init_gateway.down.sql`
- [x] T003 Add `Retained bool` field to `MessageRequestLog` in `internal/core/domain/message_request_log.go`
- [x] T004 Add `DispatchProviderAndDatabase` (`provider_and_database`) constant and extend `ParseMessageDispatchMode` in `internal/core/domain/dispatch_mode.go`
- [x] T005 Add portal-value mappers in `internal/core/domain/dispatch_mode.go`: `SettingValueToDispatchMode` (`memory`, `both`, `memory_database`, `providers`, `provider`, `provider_database` → gateway modes) and `ShouldRetainRequestLog(mode) bool` (`true` for `memory_and_provider` and `provider_and_database` only)
- [x] T006 [P] Add table-driven tests in `internal/core/domain/dispatch_mode_test.go` for setting-value mapping, legacy aliases (`both`, `providers`), and `ShouldRetainRequestLog`

**Checkpoint**: Domain and schema ready — gateway, repository, and handler work can proceed

---

## Phase 3: User Story 1 — Select Message Dispatch Mode (Priority: P1) 🎯 MVP

**Goal**: Four retention options in Portal; correct dispatch path and `retained` flag on every request log insert; **Provider + Database** behaves like Provider Only with `retained = true`

**Independent Test**: Set each mode in Settings → Data Retention, send a test message, verify dispatch behavior and `message_request_logs.retained` per quickstart.md Scenarios 1–4; Recent Requests shows all sends

### Implementation for User Story 1

- [x] T007 [US1] Extract shared provider-dispatch helper and add `DispatchProviderAndDatabase` case (same logic as `DispatchProviderOnly`) in `internal/core/service/gateway_service.go` — keep `DispatchMemoryAndProvider` / `memory_and_provider` unchanged
- [x] T008 [US1] Update `resolveDispatchMode` in `internal/core/service/gateway_service.go` to read `message_dispatch_mode` first, fall back to `data_retention`, and map portal values via `SettingValueToDispatchMode`
- [x] T009 [US1] Include `retained` in INSERT in `internal/infrastructure/repository/postgres/message_request_log_repository.go`
- [x] T010 [US1] Resolve workspace dispatch mode and set `entry.Retained` from `ShouldRetainRequestLog` in `internal/presentation/handler/send_helper.go` (add service helper on `GatewayService` if needed — e.g. `ResolveDispatchMode` / `ShouldRetainForWorkspace`)
- [x] T011 [P] [US1] Add `provider_database` to `RetentionMode` in `frontend/src/features/settings/settings.types.ts`
- [x] T012 [P] [US1] Add **Provider + Database** radio option in `frontend/src/features/settings/pages/settings.page.tsx` (continue saving via `data_retention` key per v1 contract)
- [x] T013 [US1] Update legacy value display mapping (`both` → Memory + Database, `providers` → Provider only) in `frontend/src/features/settings/hooks/use-settings.hook.ts`
- [x] T014 [P] [US1] Add gateway tests for `provider_and_database` dispatch in `internal/core/service/gateway_service_test.go`
- [x] T015 [P] [US1] Assert `retained` true/false per mode in `internal/presentation/handler/send_helper_test.go`

**Checkpoint**: All four modes selectable; sends produce correct dispatch + `retained`; Recent Requests unchanged

---

## Phase 4: User Story 2 — Distinguish Operational vs Retained Request Logs (Priority: P2)

**Goal**: `retained` column queryable; Recent Requests lists all rows; retained-only filter available for future export

**Independent Test**: SQL filter `WHERE retained = true` returns only Memory + Database and Provider + Database sends; Recent Requests API returns rows regardless of `retained`

### Implementation for User Story 2

- [x] T016 [US2] Include `l.retained` in SELECT/Scan for `ListWithSource` in `internal/infrastructure/repository/postgres/message_request_log_repository.go` (no `WHERE retained` filter)
- [x] T017 [US2] Confirm `PortalService.ListLogs` in `internal/core/service/portal_service.go` and `portal_handler` log list path pass through all rows without `retained` filter — add comment if already correct
- [x] T018 [P] [US2] Document retained-only query pattern in `specs/WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config/quickstart.md` Scenario pass criteria if implementation differs

**Checkpoint**: Operational list unfiltered; retained rows identifiable via column

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Docs, validation, quality gate

- [x] T019 [P] Update dispatch mode tables in `docs/backend/architecture.md`, `docs/backend/portal-inbox.md`, and `docs/backend/usage.md` (four modes, `retained` column, `provider_and_database`)
- [ ] T020 Run manual scenarios in `specs/WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config/quickstart.md`
- [x] T021 Run verification chain per `docs/agents/verification.md` (`golangci-lint`, `/smell develop`, `make audit`)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on T001 — **blocks all user stories**
- **User Story 1 (Phase 3)**: Depends on Phase 2 completion
- **User Story 2 (Phase 4)**: Depends on T009 (retained persisted); can overlap late US1 tasks
- **Polish (Phase 5)**: Depends on US1 + US2 complete

### User Story Dependencies

- **US1 (P1)**: Requires foundational domain + schema; no dependency on US2
- **US2 (P2)**: Requires `retained` column written (T009/T010); list/read path only

### Parallel Opportunities

```text
After T005:
  T006 [P] dispatch_mode_test.go

After T010:
  T011 [P] settings.types.ts
  T012 [P] settings.page.tsx
  T014 [P] gateway_service_test.go
  T015 [P] send_helper_test.go

US2:
  T016 + T018 [P] can run in parallel after T009
```

---

## Parallel Example: User Story 1

```bash
# Frontend (after T010):
frontend/src/features/settings/settings.types.ts
frontend/src/features/settings/pages/settings.page.tsx

# Tests (after T010):
internal/core/service/gateway_service_test.go
internal/presentation/handler/send_helper_test.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1–2 (schema + domain)
2. Complete Phase 3 (T007–T015)
3. **STOP and VALIDATE** quickstart Scenarios 1–4
4. Demo four-mode retention + `retained` column

### Incremental Delivery

1. Foundation → US1 (MVP: modes + logging)
2. US2 → retained queryable in list/SQL
3. Polish → docs + audit

### Out of Scope (v1 — do not implement)

- Renaming `memory_and_provider` → `memory_and_database`
- GET/PATCH settings normalization or rejecting `data_retention`
- TTL purge job for `retained = false`
- Changes to `pkg/gateway` SDK

---

## Notes

- Portal continues using `data_retention` setting key; gateway must fall back to that key (T008).
- Do **not** gate log inserts by mode — only set `retained` (FR-010).
- `ListWithSource` must never filter by `retained` (Recent Requests).
