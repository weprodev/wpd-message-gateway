# Feature Specification: Align Memory + Database Data Retention

**Feature Branch**: `001-memory-database-retention`

**Created**: 2026-06-22

**Status**: Draft

**Input**: User description: "Align backend data retention naming and behavior with the Portal. Rename Memory + Provider to Memory + Database. Persist to RAM and stored_messages; do not send via provider. Normalize workspace_settings value from both to memory_database."

## Clarifications

### Session 2026-06-22

- Q: Should `memory` retention persist sent API requests to the database? → A: No — only `memory_database` persists sent API requests to durable storage; `memory` remains RAM-only.
- Q: Should legacy `memory_and_provider` / `memory_provider` names remain in code with normalization? → A: No — rename all developer-facing identifiers to `memory_database` / `memory_and_database`; do not keep provider-named symbols that imply provider dispatch.
- Q: Is `message_dispatch_mode` in scope for this change? → A: Keep the setting **key name** `message_dispatch_mode` unchanged; update any **files and value mappings** tied to Memory + Database (e.g. replace `memory_and_provider` with `memory_and_database` in dispatch-mode maps, docs, and related code).
- Q: What `dispatch_status` should stored messages have after a successful Memory + Database send? → A: `sent` (updated on successful capture, not left as `pending`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Memory + Database captures without provider send (Priority: P1)

A workspace administrator selects **Memory + Database** in Portal settings so outbound messages are available in the Portal inbox for preview and also durably stored for audit/replay — without invoking any external provider.

**Why this priority**: This is the core behavioral fix; the current backend incorrectly sends messages through the provider under a misleading name.

**Independent Test**: Set `data_retention` to `memory_database`, send a message, and verify it appears in the Portal inbox and in durable message storage with no provider dispatch.

**Acceptance Scenarios**:

1. **Given** a workspace with `data_retention` = `memory_database`, **When** an outbound message is sent via the API, **Then** the message is written to the in-process inbox (RAM) **and** persisted to durable message storage (`stored_messages`), and **no** external provider is called.
2. **Given** a workspace with `data_retention` = `memory_database`, **When** the send completes successfully, **Then** the durable stored message has `dispatch_status` = `sent`, and the API response metadata reflects Memory + Database naming (not Memory + Provider).
3. **Given** a workspace with `data_retention` = `memory`, **When** a message is sent via the API, **Then** the message is captured in RAM only — it is **not** written to durable message storage and **no** provider is called.

---

### User Story 2 - Consistent naming across Portal and backend (Priority: P1)

Developers, administrators, and operators see **Memory + Database** naming in code, API responses, and workspace settings — with no remaining `memory_and_provider`, `memory_provider`, or `both` identifiers that imply provider dispatch for this mode.

**Why this priority**: Provider-named symbols mislead developers into assuming messages are sent through an integration when they are not.

**Independent Test**: Search the backend codebase for `memory_and_provider` / `memory_provider` after implementation — zero matches except removed migration comments; settings API returns `memory_database`.

**Acceptance Scenarios**:

1. **Given** a workspace saved with Memory + Database in the Portal, **When** settings are read via API, **Then** `data_retention` is `memory_database` (not `both` or any provider-named value).
2. **Given** existing workspace rows with legacy value `both` for `data_retention`, **When** the database migration runs, **Then** those rows are updated to `memory_database`.
3. **Given** backend dispatch code, constants, functions, and log messages for this retention policy, **When** reviewed by a developer, **Then** only `memory_database` / `memory_and_database` naming is used — no `memory_and_provider` or `memory_provider` symbols remain.

---

### User Story 3 - Other retention modes remain stable (Priority: P2)

Workspaces using Memory only, Providers only, or Provider + Database must continue to behave exactly as today.

**Why this priority**: Prevents regression while fixing the misaligned mode.

**Independent Test**: Send messages under each of the three unaffected modes and compare outcomes to pre-change baselines.

**Acceptance Scenarios**:

1. **Given** `data_retention` = `provider`, **When** a message is sent, **Then** only the active provider is invoked (no RAM inbox capture, no durable storage).
2. **Given** `data_retention` = `provider_database`, **When** a message is sent, **Then** the provider is invoked and the full payload is persisted durably (unchanged behavior).
3. **Given** `data_retention` = `memory`, **When** a message is sent, **Then** only the Portal inbox (RAM) capture occurs — no row is created in durable message storage.

---

### Edge Cases

- What happens when durable storage write fails in Memory + Database mode? The send fails and neither success metadata nor a `sent` status is recorded (fail-closed, same posture as Provider + Database persistence failures).
- What happens when RAM inbox write fails in Memory + Database mode? The send fails; no durable row is committed with `sent` status.
- What happens when no stored-message writer is configured? Memory + Database dispatch returns an error (durable storage is required for this mode).
- What happens when RAM capture succeeds but marking the stored message `sent` fails? The send is reported as failed; the stored row must not remain in a misleading success state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST treat Portal **Memory + Database** (`memory_database`) as: capture message in RAM (Portal inbox) **and** persist the full payload to durable message storage — **without** calling any external provider.
- **FR-002**: The system MUST **not** persist sent API requests to durable message storage when `data_retention` = `memory`; RAM inbox capture only.
- **FR-003**: The system MUST rename all backend identifiers for this policy from `memory_and_provider` / `memory_provider` to `memory_and_database` / `memory_database` — constants, dispatch mode values, function names, and log messages. No provider-named symbols may remain for this retention mode.
- **FR-004**: The system MUST store the canonical workspace setting value `memory_database` under the `data_retention` key in workspace settings (not `both`).
- **FR-005**: Existing `data_retention` rows or seeds with value `both` MUST be corrected to `memory_database` in the existing init migration (`20260318000000_init_gateway.up.sql`); no new migration file.
- **FR-006**: The system MUST map `memory_database` to dispatch behavior that writes to RAM inbox and durable storage (mirroring Provider + Database persistence semantics) — but **must not** invoke provider dispatch.
- **FR-007**: On successful Memory + Database send, the system MUST set the durable stored message `dispatch_status` to `sent` and record `dispatched_at`; provider outcome fields (`provider_message_id`, `provider_status_code`) remain empty.
- **FR-008**: API response metadata for Memory + Database sends MUST expose `stored_message_id` and inbox message identity, use `memory_and_database` as `dispatch_mode`, and MUST NOT indicate provider dispatch occurred.
- **FR-009**: Memory only, Providers only, and Provider + Database modes MUST retain their current behavior without regression.
- **FR-010**: Portal UI labels and retention options MUST remain unchanged (frontend is already correct).
- **FR-011**: Documentation describing retention modes MUST be updated to reflect Memory + Database semantics and remove incorrect Memory + Provider / provider-send descriptions for this mode.
- **FR-012**: The workspace setting key name `message_dispatch_mode` MUST remain unchanged; any files that map or document its values for Memory + Database MUST use `memory_and_database` (not `memory_and_provider` or `both`).

- **Workspace setting (`data_retention`)**: Canonical retention policy for a workspace. Allowed values: `memory`, `memory_database`, `provider`, `provider_database`.
- **Dispatch mode (`memory_and_database`)**: Backend routing mode corresponding to Memory + Database; replaces the removed `memory_and_provider` mode.
- **Inbox message (RAM)**: Ephemeral in-process capture shown in the Portal inbox; written for `memory` and `memory_database` only.
- **Stored message (durable)**: Persisted outbound payload in `stored_messages`; written for `memory_database` and `provider_database` only. After successful Memory + Database capture, `dispatch_status` = `sent`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of test sends with `memory_database` produce both an inbox-visible message and a durable stored message with `dispatch_status` = `sent`, with zero provider invocations.
- **SC-002**: 100% of test sends with `memory` produce an inbox-visible message and **zero** durable stored message rows.
- **SC-003**: 100% of workspace settings reads after Portal save return `data_retention` = `memory_database` (no `both` in stored or returned values post-migration).
- **SC-004**: Post-migration, zero workspace rows retain `both` as the `data_retention` value.
- **SC-005**: Regression suite for the three unaffected retention modes passes with no behavior change.
- **SC-006**: Backend codebase contains no `memory_and_provider` or `memory_provider` symbols for this retention policy after implementation.

## Assumptions

- Durable message storage is the existing `stored_messages` table used by Provider + Database mode; no new storage system is introduced.
- The `message_dispatch_mode` workspace setting **key name** is unchanged; value mappings and related files (e.g. `dispatch_mode.go`, `settings.api.ts`, docs) are updated wherever they reference Memory + Database under the old provider-named values.
- Portal UI labels require no changes; `settings.api.ts` dispatch-mode maps are updated where Memory + Database values are referenced.
- Database changes are made in the existing init migration (`20260318000000_init_gateway.up.sql`); no new migration file.
- For Memory + Database, `sent` means the message was successfully captured to RAM and durable storage — not that an external provider delivered it.
