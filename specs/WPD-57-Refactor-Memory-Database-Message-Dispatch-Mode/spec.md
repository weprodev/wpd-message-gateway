# Feature Specification: Refactor Memory + Database Message Dispatch Mode

**Feature Branch**: `WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode`

**Created**: 2026-06-29

**Status**: In Progress

**Input**: Refactor **Memory + Database** dispatch so it behaves like **Memory only** (no provider send) while marking request logs as retained in the database; rename legacy identifiers (`both`, `memory_and_provider`, etc.) to consistent `memory_database` / `memory_and_database` naming across the codebase.

**Supersedes (partial)**: WPD-33 clarifications that set `retained = false` for **Memory + Database** — this feature changes that to `retained = true` and removes provider dispatch from the mode.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Memory + Database Captures Locally Without Provider (Priority: P1)

A workspace admin selects **Memory + Database** so outbound messages are captured in the portal inbox (in-process memory) and request logs are marked for long-term retention — without calling any external integration.

**Why this priority**: This is the core behavior change; naming and retention rules depend on the corrected dispatch path.

**Independent Test**: Set dispatch mode to **Memory + Database**, send a message on each channel, and verify the message appears in the inbox, no provider API is invoked, and the request log row has `retained = true`.

**Acceptance Scenarios**:

1. **Given** **Memory + Database** is active and a real integration is connected, **When** a message is sent, **Then** it is captured in process memory (portal inbox) only and is **not** sent through the integration.
2. **Given** **Memory + Database** is active, **When** a message is sent, **Then** dispatch behavior matches **Memory only** (no provider lookup required for success).
3. **Given** **Memory + Database** is active, **When** a message is sent, **Then** a request log row is created with `retained = true`.
4. **Given** **Memory only** is active, **When** a message is sent, **Then** request log rows still have `retained = false` (unchanged).

---

### User Story 2 - Consistent Naming Across Settings and Gateway (Priority: P1)

Developers and operators see one coherent vocabulary: user-facing setting values use `memory_database`; internal gateway mode strings use `memory_and_database`, replacing misleading `provider` terminology for this mode.

**Why this priority**: Incorrect names caused the mode to be implemented as memory **and** provider; renaming prevents recurrence and aligns code with behavior.

**Independent Test**: Grep codebase and run settings round-trip tests — no remaining canonical references to `memory_and_provider` or `both` as active identifiers; legacy values still resolve correctly on read where specified.

**Acceptance Scenarios**:

1. **Given** a workspace setting stored as `memory_database`, **When** the gateway resolves dispatch mode, **Then** it uses the gateway mode `memory_and_database` (not `memory_and_provider`).
2. **Given** a workspace setting stored with legacy value `both`, **When** settings are read, **Then** the mode resolves to **Memory + Database** behavior and maps to `memory_and_database` at dispatch time.
3. **Given** new settings are saved, **When** **Memory + Database** is selected, **Then** the persisted value is `memory_database` (not `both`).

---

### User Story 3 - Retained Request Logs for Database-Backed Modes (Priority: P2)

A compliance reviewer filters request logs by `retained = true` and sees entries from both **Memory + Database** and **Provider + Database** workspaces.

**Why this priority**: Retention semantics must stay consistent now that two modes mark logs as retained.

**Independent Test**: Send messages under all four modes; query `message_request_logs` where `retained = true` and verify only **Memory + Database** and **Provider + Database** rows match.

**Acceptance Scenarios**:

1. **Given** messages sent under **Provider + Database**, **When** filtered by `retained = true`, **Then** those rows are included (unchanged from WPD-33).
2. **Given** messages sent under **Memory + Database**, **When** filtered by `retained = true`, **Then** those rows are included.
3. **Given** messages sent under **Memory only** or **Provider only**, **When** filtered by `retained = true`, **Then** no rows are returned.
4. **Given** any dispatch mode, **When** **Recent Requests** is opened, **Then** operational request logs appear regardless of `retained` value.

---

### Edge Cases

- Workspace with legacy gateway string `memory_and_provider` in settings or logs: must resolve to `memory_and_database` behavior on read during transition.
- **Memory + Database** with no connected integration: must succeed (same as **Memory only**) — provider absence must not block dispatch.
- Concurrent mode change while a message is in flight: mode effective at dispatch time governs behavior and `retained` value.
- All four channels (email, SMS, push, chat) follow the same dispatch, logging, and `retained` rules.
- Provider dispatch failure is N/A for **Memory + Database** after this change (provider is never called).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: **Memory + Database** MUST capture messages in process memory (portal inbox) and MUST NOT invoke external provider integrations.
- **FR-002**: **Memory + Database** dispatch behavior MUST be equivalent to **Memory only** for the outbound send path (memory capture only).
- **FR-003**: **Memory + Database** MUST mark request log rows in `message_request_logs` with `retained = true`.
- **FR-004**: **Provider + Database** MUST remain unchanged: provider dispatch identical to **Provider only** with `retained = true`.
- **FR-005**: **Memory only** and **Provider only** MUST continue to set `retained = false` on request log rows.
- **FR-006**: System MUST rename gateway dispatch mode `memory_and_provider` to `memory_and_database` everywhere it is a canonical identifier (constants, enums, metadata, tests, docs).
- **FR-007**: System MUST rename domain constant `DispatchMemoryAndProvider` to `DispatchMemoryAndDatabase` (or equivalent) mapping to string value `memory_and_database`.
- **FR-008**: Setting/API canonical value `memory_database` MUST map to gateway mode `memory_and_database`; legacy setting value `both` MUST continue to resolve to **Memory + Database** on read only.
- **FR-009**: New writes of dispatch mode MUST use canonical setting values (`memory`, `memory_database`, `provider`, `provider_database`) — not `both` or `memory_and_provider`.
- **FR-010**: `ShouldRetainRequestLog` (or equivalent) MUST return `true` for both `memory_and_database` and `provider_and_database` gateway modes.
- **FR-011**: Request log creation MUST follow the existing save flow for all modes — no gating or skipping of inserts; only the `retained` value and dispatch path change for **Memory + Database**.
- **FR-012**: Portal **Recent Requests** MUST list request logs from all dispatch modes without filtering by `retained`.
- **FR-013**: Portal UI labels for **Memory + Database** MUST describe memory capture with database-backed request log retention (not provider send).

### Naming Migration Reference

| Context | Before | After |
| ------- | ------ | ----- |
| Setting / API value | `both` | `memory_database` (legacy `both` accepted on read) |
| Setting / API value | `memory_provider` (if present) | `memory_database` |
| Gateway mode string | `memory_and_provider` | `memory_and_database` |
| Gateway constant | `DispatchMemoryAndProvider` | `DispatchMemoryAndDatabase` |
| Portal internal mode id | `memory_and_provider` | `memory_and_database` |

### Key Entities

- **Message Dispatch Mode**: Per-workspace setting (`memory` | `memory_database` | `provider` | `provider_database`) under `message_dispatch_mode`; gateway runtime modes (`memory_only`, `memory_and_database`, `provider_only`, `provider_and_database`).
- **Message Request Log**: Row in `message_request_logs`; `retained = true` for **Memory + Database** and **Provider + Database**.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of channel tests under **Memory + Database** show inbox capture, zero provider calls, and `retained = true` on the request log row.
- **SC-002**: 100% of channel tests under **Memory only** and **Provider only** still show `retained = false`.
- **SC-003**: No canonical code references remain to `memory_and_provider` or `DispatchMemoryAndProvider` after refactor (legacy read aliases excepted).
- **SC-004**: **Recent Requests** shows entries for all four modes with no regression.
- **SC-005**: Retention filter (`retained = true`) returns rows from **Memory + Database** and **Provider + Database** only.

## Assumptions

- “Saved in Database” for **Memory + Database** means request logs with `retained = true`, not persisting full message content to a separate messages table (same distinction as **Provider + Database**).
- Message content for inbox display remains in-process memory capture, identical to **Memory only**.
- Setting key remains `message_dispatch_mode` (unchanged from WPD-33).
- Legacy gateway string `memory_and_provider` may be accepted on read during migration but must not be written by new code.
- Documentation, Bruno tests, and architecture docs are updated as part of implementation (detailed in planning phase).
- Feature spec directory and branch share the name `WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode`.

## Out of Scope

- TTL purge job for non-retained request logs.
- Persisting message body content to PostgreSQL (future inbox persistence).
- Changes to `pkg/gateway` embedded SDK dispatch (unless required for string parity).
- GET/PATCH settings value normalization (rewriting legacy DB rows to canonical values on read).
