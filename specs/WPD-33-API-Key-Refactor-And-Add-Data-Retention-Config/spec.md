# Feature Specification: Message Dispatch Modes & Request Retention

**Feature Branch**: `WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config`

**Created**: 2026-06-27

**Status**: In Progress

**Input**: Add **Provider + Database** message dispatch mode; mark database-backed modes on request logs via a `retained` flag; rename backend `data_retention` to `message_dispatch_mode`; keep operational logging flow unchanged.

## Clarifications

### Session 2026-06-27

- Q: What does **Provider + Database** do compared to **Provider Only**? → A: Identical outbound dispatch and request-logging behavior as **Provider Only**; the only difference is `retained = true` on the request log row (no message content persistence).
- Q: Which modes set `retained = true`? → A: **Memory + Database** and **Provider + Database** only.
- Q: Should the request-log insert flow change by mode? → A: No — all modes continue saving request logs exactly as today; only the new `retained` column distinguishes database-backed retention from operational-only modes.
- Q: What is the canonical backend setting key? → A: `message_dispatch_mode` (replace all backend uses of `data_retention`).
- Q: What are the canonical mode values? → A: `memory`, `memory_database`, `provider`, `provider_database` (legacy `both` → `memory_database`, `providers` → `provider` accepted on read only).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select Message Dispatch Mode (Priority: P1)

A workspace admin opens **Settings → Data Retention** and chooses how outbound messages are handled and which request logs are marked for long-term retention.

**Why this priority**: Dispatch mode drives gateway behavior and the `retained` flag; without it, other stories cannot be validated.

**Independent Test**: Change each mode, save, send test messages, and verify dispatch behavior, request-log presence, and `retained` values match the selection. Confirm **Recent Requests** continues to show entries for all modes.

**Acceptance Scenarios**:

1. **Given** a workspace with no saved policy, **When** the admin opens Data Retention, **Then** **Memory only** is selected by default.
2. **Given** **Memory only** is active, **When** a message is sent, **Then** it is captured in process memory only (no message content in the database) and any request log row created follows today's logging rules with `retained = false`.
3. **Given** **Memory + Database** is active, **When** a message is sent, **Then** message content is persisted for portal inbox/database access and any request log row created has `retained = true`.
4. **Given** **Provider only** is active, **When** a message is sent, **Then** it is dispatched through the connected integration with no message content in the database and any request log row created has `retained = false`.
5. **Given** **Provider + Database** is active, **When** a message is sent, **Then** dispatch behavior is identical to **Provider only** and any request log row created has `retained = true` (no message content persistence).
6. **Given** a workspace saved with legacy value `both`, **When** settings are loaded, **Then** the UI shows **Memory + Database** selected.
7. **Given** a workspace saved with legacy value `providers`, **When** settings are loaded, **Then** the UI shows **Provider only** selected.

---

### User Story 2 - Distinguish Operational vs Retained Request Logs (Priority: P2)

An operator or compliance reviewer needs to tell which request logs belong to database-backed dispatch modes versus operational-only modes.

**Why this priority**: The `retained` flag enables future retention policy and purge rules without changing Recent Requests.

**Independent Test**: Send messages under each mode and query request logs; verify `retained` is set correctly while Recent Requests remains unfiltered by mode.

**Acceptance Scenarios**:

1. **Given** request logs from **Memory + Database** or **Provider + Database**, **When** filtered by `retained = true`, **Then** all matching rows are returned.
2. **Given** request logs from **Memory only** or **Provider only**, **When** filtered by `retained = true`, **Then** no rows are returned.
3. **Given** any dispatch mode, **When** Recent Requests is opened, **Then** operational request logs appear regardless of `retained` value.

---

### Edge Cases

- Workspace with unknown legacy dispatch value: falls back to **Memory only** and can be re-saved with a canonical value.
- Concurrent mode save while a message is in flight: the mode effective at dispatch time is the value read from settings for that request.
- Provider dispatch failure in **Provider + Database** mode: same error handling as **Provider only**; request logging follows today's rules for failed sends.
- All four channels (email, SMS, push, chat) follow the same dispatch, logging, and `retained` rules.
- Non-retained rows (`retained = false`) may be purged after a configurable operational TTL without affecting retained rows (purge job out of scope for v1).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose four message dispatch modes: **Memory only**, **Memory + Database**, **Provider only**, and **Provider + Database**.
- **FR-002**: System MUST persist dispatch mode per workspace using canonical values `memory`, `memory_database`, `provider`, and `provider_database` under the backend setting key `message_dispatch_mode`.
- **FR-003**: System MUST replace all backend references to `data_retention` with `message_dispatch_mode`.
- **FR-004**: System MUST treat legacy values `both` as `memory_database` and `providers` as `provider` when reading settings; new writes MUST use canonical values only.
- **FR-005**: **Memory only** MUST capture messages in process memory only and MUST NOT persist message content to the database.
- **FR-006**: **Memory + Database** MUST persist message content for portal/database access and MUST mark request log rows with `retained = true`.
- **FR-007**: **Provider only** MUST dispatch through the connected integration with the same behavior as today's provider-only path, MUST NOT persist message content, and MUST mark request log rows with `retained = false`.
- **FR-008**: **Provider + Database** MUST dispatch identically to **Provider only**, MUST NOT persist message content, and MUST mark request log rows with `retained = true`.
- **FR-009**: System MUST add a boolean `retained` column to request logs (`message_request_logs`): `true` for **Memory + Database** and **Provider + Database**; `false` for **Memory only** and **Provider only**.
- **FR-010**: Request log creation MUST follow the existing save flow for all modes — no gating or skipping of inserts by dispatch mode; only the `retained` value differs.
- **FR-011**: Portal **Recent Requests** MUST list request logs from all dispatch modes (not filtered by `retained`).
- **FR-012**: Long-term retention queries and exports MUST filter request logs where `retained = true` only.

### Key Entities

- **Message Dispatch Mode**: Per-workspace setting (`memory` | `memory_database` | `provider` | `provider_database`) stored as `message_dispatch_mode`, controlling message capture, provider dispatch, and the `retained` flag on request logs.
- **Message Request Log**: Row recording a gateway API request in `message_request_logs`; `retained` distinguishes database-backed retention modes from operational-only modes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of manual test scenarios across all four dispatch modes produce the correct combination of message content persistence, dispatch path, and `retained` flag.
- **SC-002**: Recent Requests shows request log entries for all four dispatch modes without regression.
- **SC-003**: Request logs under **Memory + Database** and **Provider + Database** have `retained = true`; logs under **Memory only** and **Provider only** have `retained = false`.
- **SC-004**: No backend code or API contract continues to use `data_retention` as the setting key or identifier.

## Assumptions

- Message content persistence for **Memory + Database** reuses the existing inbox/database path.
- Operational and retained request metadata share the same request-log table; separation is via the `retained` column only.
- Existing migration and seed files will be updated in place (no new migration file required for v1 schema changes).
- Portal UI may continue to label the settings tab "Data Retention" for users; backend storage and domain naming use `message_dispatch_mode`.
- Optional background purge of non-retained rows after an operational window is out of scope for v1.
