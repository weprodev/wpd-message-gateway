# Feature Specification: Data Retention Modes

**Feature Branch**: `001-memory-database-retention`

**Created**: 2025-06-23

**Status**: Draft

**Input**: User description: "Add Provider + Database retention mode; refactor Memory Only and Provider Only to skip database persistence; gate request logging to Memory + Database and Provider + Database only; log only successful sends to message_request_logs."

## Clarifications

### Session 2025-06-23

- Q: What does Provider + Database persist compared to Provider Only? → A: Same outbound dispatch behavior as Provider Only; only request metadata is written to the database (`message_request_logs`), not message content.
- Q: Which retention modes write request logs? → A: Only Memory + Database and Provider + Database.
- Q: Should Memory Only and Provider Only write to the database? → A: No — neither message content nor request logs.
- Q: What are the canonical retention mode values? → A: `memory`, `memory_database`, `provider`, `provider_database` (legacy `both` and `providers` accepted as read-time aliases only).

### Session 2026-06-24

- Q: When should a request log row be created? → A: Only after a **successful send** — dispatch completes without error and the gateway returns success to the caller. Failed validation, auth, or dispatch attempts MUST NOT be written to `message_request_logs`.
- Q: What data is stored in request logs for database-backed modes? → A: All available request metadata fields for every channel (email, SMS, push, chat) on successful sends only.
- Q: Should `DispatchMemoryAndProvider` be renamed to `DispatchMemoryAndDatabase`? → A: No — keep `DispatchMemoryAndProvider` (`memory_and_provider`); map portal `memory_database` to it in domain helpers.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select Data Retention Policy (Priority: P1)

A workspace admin opens **Settings → Data Retention** and chooses how outbound messages and request metadata are stored.

**Why this priority**: Retention policy drives all gateway persistence behavior; without it, other stories cannot be validated.

**Independent Test**: Change each retention option, save, send a test message, and verify persistence matches the selected mode.

**Acceptance Scenarios**:

1. **Given** a workspace with no saved policy, **When** the admin opens Data Retention, **Then** **Memory only** is selected by default.
2. **Given** **Memory only** is active, **When** a message is sent successfully, **Then** it is captured in process memory only and nothing is written to the database (no message content, no request logs).
3. **Given** **Memory + Database** is active, **When** a message is sent successfully, **Then** message content is persisted for portal inbox/database access and a request log row is created.
4. **Given** **Memory + Database** is active, **When** a send fails (validation, auth, or dispatch error), **Then** no request log row is created.
5. **Given** **Provider only** is active, **When** a message is sent successfully, **Then** it is dispatched through the connected integration with no database persistence (no message content, no request logs).
6. **Given** **Provider + Database** is active, **When** a message is sent successfully, **Then** dispatch behavior matches **Provider only** and a request log row is created (no message content persistence).
7. **Given** **Provider + Database** is active, **When** a send fails, **Then** no request log row is created.
8. **Given** a workspace saved with legacy value `both`, **When** settings are loaded, **Then** the UI shows **Memory + Database** selected.
9. **Given** a workspace saved with legacy value `providers`, **When** settings are loaded, **Then** the UI shows **Provider only** selected.

---

### Edge Cases

- Workspace with unknown legacy retention value: falls back to **Memory only** and can be re-saved with a canonical value.
- Concurrent retention save while a message is in flight: the mode effective at dispatch time is the value read from settings for that request.
- Provider dispatch failure in **Provider + Database** mode: no request log row; caller receives error response.
- All four channels (email, SMS, push, chat) follow the same retention and logging rules.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose four data retention modes: **Memory only**, **Memory + Database**, **Provider only**, and **Provider + Database**.
- **FR-002**: System MUST persist retention selection per workspace using canonical values `memory`, `memory_database`, `provider`, and `provider_database`.
- **FR-003**: System MUST treat legacy values `both` as `memory_database` and `providers` as `provider` when reading settings; new writes MUST use canonical values only.
- **FR-004**: **Memory only** MUST capture messages in process memory only and MUST NOT write message content or request logs to the database.
- **FR-005**: **Memory + Database** MUST persist message content for portal/database access on successful sends AND write request logs on successful sends only.
- **FR-006**: **Provider only** MUST dispatch through the connected integration with the same behavior as today's provider-only path and MUST NOT write message content or request logs to the database.
- **FR-007**: **Provider + Database** MUST dispatch identically to **Provider only** and MUST write request logs on successful sends only (no message content persistence).
- **FR-008**: Request logging MUST occur only when retention is `memory_database` or `provider_database` AND the send completed successfully.
- **FR-009**: Request logs MUST populate all available metadata fields (workspace, API key, channel, HTTP method, status, endpoint, provider, request ID, duration, error message when applicable) for every successful send across all supported channels.
- **FR-010**: Failed requests (invalid JSON, missing auth, dispatch errors) MUST NOT create `message_request_logs` rows regardless of retention mode.

### Key Entities

- **Data Retention Policy**: Per-workspace setting (`memory` | `memory_database` | `provider` | `provider_database`) controlling message capture, provider dispatch, and request log persistence.
- **Message Request Log**: Audit row for a successful outbound gateway request (workspace, API key, channel, status, timing, provider metadata) stored only for database-backed retention modes.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of retention mode combinations produce the correct persistence outcome (memory-only, message DB, request logs) in manual test scenarios.
- **SC-002**: Zero request log rows are created for Memory only, Provider only, or any failed send path.
- **SC-003**: Successful sends under Memory + Database and Provider + Database create exactly one request log row per request with all metadata fields populated.

## Assumptions

- Retention policy is stored in workspace settings (`data_retention` key) and mapped to gateway dispatch behavior server-side (`memory_database` → `memory_and_provider`).
- Message content persistence for **Memory + Database** reuses the existing inbox/database path.
- Request logs are written to the `message_request_logs` table.
- A **successful send** means dispatch completed without error and the gateway responds with success to the API caller.
