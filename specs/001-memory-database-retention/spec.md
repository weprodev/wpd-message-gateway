# Feature Specification: Data Retention Modes

**Feature Branch**: `001-memory-database-retention`

**Created**: 2025-06-23

**Status**: Draft

**Input**: User description: "Add Provider + Database retention mode; separate operational request logging from long-term retention using a retained flag on message_request_logs; keep Recent Requests working for all modes."

## Clarifications

### Session 2025-06-23

- Q: What does Provider + Database persist compared to Provider Only? → A: Same outbound dispatch behavior as Provider Only; request metadata is marked **retained** in the database (`message_request_logs.retained = true`), not message content.
- Q: Which retention modes mark request logs as retained? → A: Only Memory + Database and Provider + Database.
- Q: Should Memory Only and Provider Only write message content to the database? → A: No — message content is not persisted.
- Q: What are the canonical retention mode values? → A: `memory`, `memory_database`, `provider`, `provider_database` (legacy `both` and `providers` accepted as read-time aliases only).

### Session 2026-06-24

- Q: When should a request log row be created? → A: Only after a **successful send** — dispatch completes without error and the gateway returns success to the caller. Failed validation, auth, or dispatch attempts MUST NOT be written to `message_request_logs`.
- Q: What data is stored in request logs for database-backed modes? → A: All available request metadata fields for every channel (email, SMS, push, chat) on successful sends; `retained = true` only for Memory + Database and Provider + Database.
- Q: Should `DispatchMemoryAndProvider` be renamed to `DispatchMemoryAndDatabase`? → A: No — keep `DispatchMemoryAndProvider` (`memory_and_provider`); map portal `memory_database` to it in domain helpers.

### Session 2026-06-25

- Q: Should gating inserts to `message_request_logs` by retention mode break Recent Requests? → A: No — **operational logging** (Recent Requests) and **retention policy** are separate concerns on the same table via a `retained` flag (Idea 3).
- Q: What does `retained` mean? → A: `retained = true` means the row is kept per long-term data-retention policy; `retained = false` means operational/monitoring only (eligible for short TTL purge).
- Q: Does Memory only / Provider only create request log rows? → A: On successful sends, yes — with `retained = false` so Recent Requests works; message content is still not persisted.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select Data Retention Policy (Priority: P1)

A workspace admin opens **Settings → Data Retention** and chooses how outbound messages and long-term request metadata are stored.

**Why this priority**: Retention policy drives gateway persistence behavior; without it, other stories cannot be validated.

**Independent Test**: Change each retention option, save, send a test message, and verify persistence matches the selected mode. Confirm **Recent Requests** shows successful sends for all modes.

**Acceptance Scenarios**:

1. **Given** a workspace with no saved policy, **When** the admin opens Data Retention, **Then** **Memory only** is selected by default.
2. **Given** **Memory only** is active, **When** a message is sent successfully, **Then** it is captured in process memory only (no message content in DB) and a request log row is created with `retained = false`; Recent Requests shows the send.
3. **Given** **Memory + Database** is active, **When** a message is sent successfully, **Then** message content is persisted for portal inbox/database access and a request log row is created with `retained = true`.
4. **Given** **Memory + Database** is active, **When** a send fails (validation, auth, or dispatch error), **Then** no request log row is created.
5. **Given** **Provider only** is active, **When** a message is sent successfully, **Then** it is dispatched through the connected integration with no message content in DB and a request log row with `retained = false`; Recent Requests shows the send.
6. **Given** **Provider + Database** is active, **When** a message is sent successfully, **Then** dispatch behavior matches **Provider only** and a request log row is created with `retained = true` (no message content persistence).
7. **Given** **Provider + Database** is active, **When** a send fails, **Then** no request log row is created.
8. **Given** a workspace saved with legacy value `both`, **When** settings are loaded, **Then** the UI shows **Memory + Database** selected.
9. **Given** a workspace saved with legacy value `providers`, **When** settings are loaded, **Then** the UI shows **Provider only** selected.

---

### Edge Cases

- Workspace with unknown legacy retention value: falls back to **Memory only** and can be re-saved with a canonical value.
- Concurrent retention save while a message is in flight: the mode effective at dispatch time is the value read from settings for that request.
- Provider dispatch failure in **Provider + Database** mode: no request log row; caller receives error response.
- All four channels (email, SMS, push, chat) follow the same retention and logging rules.
- Non-retained rows (`retained = false`) may be purged after a configurable operational TTL without affecting retained rows.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose four data retention modes: **Memory only**, **Memory + Database**, **Provider only**, and **Provider + Database**.
- **FR-002**: System MUST persist retention selection per workspace using canonical values `memory`, `memory_database`, `provider`, and `provider_database`.
- **FR-003**: System MUST treat legacy values `both` as `memory_database` and `providers` as `provider` when reading settings; new writes MUST use canonical values only.
- **FR-004**: **Memory only** MUST capture messages in process memory only and MUST NOT persist message content to the database.
- **FR-005**: **Memory + Database** MUST persist message content for portal/database access on successful sends AND create request log rows with `retained = true` on successful sends only.
- **FR-006**: **Provider only** MUST dispatch through the connected integration with the same behavior as today's provider-only path, MUST NOT persist message content, and MUST create request log rows with `retained = false` on successful sends only.
- **FR-007**: **Provider + Database** MUST dispatch identically to **Provider only** and MUST create request log rows with `retained = true` on successful sends only (no message content persistence).
- **FR-008**: On every successful send, system MUST insert one `message_request_logs` row with `retained` set according to retention mode (`true` for `memory_database` and `provider_database`; `false` otherwise).
- **FR-009**: Request logs MUST populate all available metadata fields (workspace, API key, channel, HTTP method, status, endpoint, provider, request ID, duration) for every successful send across all supported channels.
- **FR-010**: Failed requests (invalid JSON, missing auth, dispatch errors) MUST NOT create `message_request_logs` rows regardless of retention mode.
- **FR-011**: Portal **Recent Requests** MUST list operational request logs from `message_request_logs` for all retention modes (not filtered by `retained`).
- **FR-012**: Long-term retention queries and exports MUST filter `message_request_logs` where `retained = true` only.

### Key Entities

- **Data Retention Policy**: Per-workspace setting (`memory` | `memory_database` | `provider` | `provider_database`) controlling message capture, provider dispatch, and the `retained` flag on request logs.
- **Message Request Log**: Operational row for a successful outbound gateway request; `retained` distinguishes monitoring-only rows from policy-backed long-term retention.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of retention mode combinations produce the correct outcome (message content, `retained` flag) in manual test scenarios.
- **SC-002**: Recent Requests shows successful sends for all four retention modes.
- **SC-003**: Successful sends under Memory + Database and Provider + Database create exactly one request log row per request with `retained = true` and all metadata fields populated.
- **SC-004**: Successful sends under Memory only and Provider only create request log rows with `retained = false` and zero message content in the database.

## Assumptions

- Retention policy is stored in workspace settings (`data_retention` key) and mapped to gateway dispatch behavior server-side (`memory_database` → `memory_and_provider`).
- Message content persistence for **Memory + Database** reuses the existing inbox/database path.
- Operational and retained request metadata share `message_request_logs`; separation is via the `retained` column (Idea 3).
- A **successful send** means dispatch completed without error and the gateway responds with success to the API caller.
- Optional background purge of non-retained rows after an operational window is out of scope for v1 unless explicitly scheduled later.
