# Feature Specification: Data Retention Modes & API Key Modals

**Feature Branch**: `001-memory-database-retention`

**Created**: 2025-06-23

**Status**: Draft

**Input**: User description: "Add Provider + Database retention mode; refactor Memory Only and Provider Only to skip database persistence; gate request logging to Memory + Database and Provider + Database only; replace browser prompts with Create, Regenerate, Delete, and Credentials modals for API keys (backend API key actions unchanged)."

## Clarifications

### Session 2025-06-23

- Q: What does Provider + Database persist compared to Provider Only? → A: Same outbound dispatch behavior as Provider Only; only request metadata is written to the database (`message_request_logs`), not message content.
- Q: Which retention modes write request logs? → A: Only Memory + Database and Provider + Database.
- Q: Should Memory Only and Provider Only write to the database? → A: No — neither message content nor request logs.
- Q: Are backend API key create/regenerate/delete endpoints changed? → A: No — frontend modal UX only; existing backend and database behavior for API keys stays as-is.
- Q: What are the canonical retention mode values? → A: `memory`, `memory_database`, `provider`, `provider_database` (legacy `both` and `providers` accepted as read-time aliases only).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select Data Retention Policy (Priority: P1)

A workspace admin opens **Settings → Data Retention** and chooses how outbound messages and request metadata are stored.

**Why this priority**: Retention policy drives all gateway persistence behavior; without it, other stories cannot be validated.

**Independent Test**: Change each retention option, save, send a test message, and verify persistence matches the selected mode.

**Acceptance Scenarios**:

1. **Given** a workspace with no saved policy, **When** the admin opens Data Retention, **Then** **Memory only** is selected by default.
2. **Given** **Memory only** is active, **When** a message is sent, **Then** it is captured in process memory only and nothing is written to the database (no message content, no request logs).
3. **Given** **Memory + Database** is active, **When** a message is sent, **Then** message content is persisted for portal inbox/database access and a request log row is created.
4. **Given** **Provider only** is active, **When** a message is sent, **Then** it is dispatched through the connected integration with no database persistence (no message content, no request logs).
5. **Given** **Provider + Database** is active, **When** a message is sent, **Then** dispatch behavior matches **Provider only** and a request log row is created (no message content persistence).
6. **Given** a workspace saved with legacy value `both`, **When** settings are loaded, **Then** the UI shows **Memory + Database** selected.
7. **Given** a workspace saved with legacy value `providers`, **When** settings are loaded, **Then** the UI shows **Provider only** selected.

---

### User Story 2 - Create API Key with Credentials Modal (Priority: P1)

A developer creates a new API key from **Settings → Developer** using a guided modal flow instead of a browser prompt.

**Why this priority**: API keys are required to exercise gateway endpoints; the create flow must surface one-time credentials safely.

**Independent Test**: Click **Generate key**, enter a name, confirm creation, and verify the credentials modal appears with copyable client ID and secret.

**Acceptance Scenarios**:

1. **Given** the Developer tab, **When** the user clicks **Generate key**, **Then** a modal opens with the prompt "Please add API key name", a text input (placeholder e.g. "Production"), an **X** close control, and a **Generate Key** button.
2. **Given** the create modal is open, **When** the user clicks **X**, **Then** the modal closes and no key is created.
3. **Given** the create modal is open with a non-empty name, **When** the user clicks **Generate Key**, **Then** the key is created via the existing API and a credentials modal opens.
4. **Given** the credentials modal is open, **When** displayed, **Then** it shows a warning that credentials are shown once, read-only fields for API client ID and API secret, and a copy button beside each field.
5. **Given** a copy button, **When** clicked, **Then** the value is copied to the clipboard and the button icon transitions smoothly from a clipboard icon to a checkmark icon.

---

### User Story 3 - Regenerate API Key with Confirmation (Priority: P2)

A developer regenerates an existing API key and receives new one-time credentials.

**Why this priority**: Regeneration is security-sensitive and must replace ad-hoc browser confirms with an explicit modal.

**Independent Test**: Click **Regenerate** on a key row, confirm in the modal, and verify new credentials appear once.

**Acceptance Scenarios**:

1. **Given** an existing API key, **When** the user clicks **Regenerate**, **Then** a confirmation modal opens with header "Are you sure you want to Regenerate the API Key?" and **no X** close button.
2. **Given** the regenerate modal, **When** the user clicks **Cancel**, **Then** the modal closes and the key is unchanged.
3. **Given** the regenerate modal, **When** the user clicks **Generate Key**, **Then** the key is regenerated via the existing API and the credentials modal opens with the same layout and copy behavior as create.

---

### User Story 4 - Delete API Key with Confirmation (Priority: P2)

A developer deletes an API key after explicit confirmation.

**Why this priority**: Prevents accidental key deletion.

**Independent Test**: Click **Delete**, cancel once (key remains), then confirm delete (key removed).

**Acceptance Scenarios**:

1. **Given** an existing API key, **When** the user clicks **Delete**, **Then** a confirmation modal opens with header "Are you sure you want to Delete this API Key?" and **no X** close button.
2. **Given** the delete modal, **When** the user clicks **Cancel**, **Then** the modal closes and the key is not deleted.
3. **Given** the delete modal, **When** the user clicks **Delete** (danger variant), **Then** the key is deleted via the existing API and the modal closes.

---

### Edge Cases

- Empty API key name on create: **Generate Key** remains disabled or shows inline validation; no API call is made.
- Create/regenerate API failure: credentials modal does not open; user sees an error message and can retry.
- Copy when clipboard API is unavailable: user sees a fallback error; copy button does not falsely show success.
- Workspace with unknown legacy retention value: falls back to **Memory only** and can be re-saved with a canonical value.
- Concurrent retention save while a message is in flight: the mode effective at dispatch time is the value read from settings for that request.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST expose four data retention modes: **Memory only**, **Memory + Database**, **Provider only**, and **Provider + Database**.
- **FR-002**: System MUST persist retention selection per workspace using canonical values `memory`, `memory_database`, `provider`, and `provider_database`.
- **FR-003**: System MUST treat legacy values `both` as `memory_database` and `providers` as `provider` when reading settings; new writes MUST use canonical values only.
- **FR-004**: **Memory only** MUST capture messages in process memory only and MUST NOT write message content or request logs to the database.
- **FR-005**: **Memory + Database** MUST persist message content for portal/database access AND write request logs.
- **FR-006**: **Provider only** MUST dispatch through the connected integration with the same behavior as today’s provider-only path and MUST NOT write message content or request logs to the database.
- **FR-007**: **Provider + Database** MUST dispatch identically to **Provider only** and MUST write request logs only (no message content persistence).
- **FR-008**: Request logging MUST occur only when retention is `memory_database` or `provider_database`.
- **FR-009**: API key create, regenerate, and delete MUST continue using existing backend endpoints and database operations unchanged.
- **FR-010**: Create API key MUST use a modal (not `window.prompt`) collecting the key name before submission.
- **FR-011**: Regenerate API key MUST use a confirmation modal without an **X** button; **Cancel** and **Generate Key** actions only.
- **FR-012**: Delete API key MUST use a confirmation modal without an **X** button; **Cancel** and **Delete** (danger) actions only.
- **FR-013**: After successful create or regenerate, a credentials modal MUST display API client ID and API secret once with a warning banner and per-field copy buttons with smooth clipboard-to-checkmark transition.

### Key Entities

- **Data Retention Policy**: Per-workspace setting (`memory` | `memory_database` | `provider` | `provider_database`) controlling message capture, provider dispatch, and request log persistence.
- **Message Request Log**: Audit row for an outbound gateway request (workspace, API key, channel, status, timing, provider metadata) stored only for database-backed retention modes.
- **API Key**: Named workspace credential with client ID and one-time-display secret on create/regenerate.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of retention mode combinations produce the correct persistence outcome (memory-only, message DB, request logs) in manual test scenarios.
- **SC-002**: Admins can create an API key with name entry and credential copy in under 30 seconds without browser-native dialogs.
- **SC-003**: Regenerate and delete flows require an explicit confirmation step before any mutating API call.
- **SC-004**: Zero changes to existing API key backend contract (request/response shapes and HTTP status codes remain stable).

## Assumptions

- Retention policy is stored in workspace settings (`data_retention` key) and mapped to gateway dispatch behavior server-side.
- Message content persistence for **Memory + Database** reuses the existing inbox/database path.
- Request logs are written to the `message_request_logs` table.
- Credentials modal checkmark uses the design-system icon component (not an emoji).
- Regenerate and delete modals intentionally omit the default dialog **X** close control; users must choose **Cancel** or the primary action.
