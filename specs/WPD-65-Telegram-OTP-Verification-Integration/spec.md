# Feature Specification: Telegram OTP Chat Provider

**Feature Branch**: `WPD-65-Telegram-OTP-Verification-Integration`

**Created**: 2026-07-04

**Status**: Draft

**Input**: User description: "Add a new Provider for the Chat Channel Type named Telegram that sends OTP Verification Messages through the Chat channel. Update Front End, Config, and Documents. Portal Connect follows the same Integrations flow Mailgun uses for Email: the provider is listed in the catalog under its channel section with icon, name, and description; when status is active the row shows a Connect button (not Coming soon); Connect opens a modal that loads config fields from the provider catalog, collects credentials (API key as a masked password field), and on submit saves an encrypted workspace integration via the Portal API; Connected rows show a badge and Disconnect. Recipients are identified by phone number (normalized to E.164). Provider-specific fields such as sender_username are passed via message metadata. API requests use the `message` field; the Telegram provider maps it to the provider's `code` field when set. If `message` is empty and `metadata.code_length` is set, Telegram generates a random OTP of that length. Concurrent sends are processed in batches of five."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Send OTP via Chat API (Priority: P1)

A backend application sends one-time verification codes to users through the existing Chat channel API, using the Telegram provider and the recipient's phone number.

**Why this priority**: Delivers the core business value — OTP delivery — without requiring any Portal UI changes.

**Independent Test**: Connect Telegram with a valid API key in a workspace; call the Chat send endpoint with a phone number in `to` and either (a) a numeric code in `message`, or (b) empty `message` with `metadata.code_length` set; confirm a successful delivery response.

**Acceptance Scenarios**:

1. **Given** a workspace with an active Telegram integration, **When** a client sends a chat message with `to` containing a valid phone number and `message` containing a 4–8 digit numeric code, **Then** the gateway delivers the OTP through Telegram and returns a success result with a message identifier.
2. **Given** a chat request with an invalid or unparsable phone number, **When** the gateway processes the send, **Then** the request is rejected with a clear validation error and no provider call is made.
3. **Given** a chat request with a phone number in a non-E.164 format (e.g., local format with country context available), **When** the gateway processes the send, **Then** the number is normalized to E.164 before delivery to Telegram.
4. **Given** a chat request where `message` is empty and `metadata.code_length` has a value (4–8), **When** the gateway processes the send, **Then** Telegram generates a random numeric OTP of that length and delivers it (no `code` field is sent to the provider).
5. **Given** a chat request where `message` is provided but is not a valid OTP (non-numeric or outside 4–8 digits), **When** the gateway processes the send, **Then** the request is rejected with a clear validation error.
6. **Given** a chat request where `message` is empty and `metadata.code_length` is missing or invalid, **When** the gateway processes the send, **Then** the request is rejected with a clear validation error.

---

### User Story 2 - Connect Telegram in Portal (Priority: P2)

A workspace administrator connects Telegram from the Integrations page under the Chat section. The flow matches how Mailgun is connected under Email today: on Integrations the admin sees providers grouped by channel; Telegram appears as a row in the Chat group with its catalog icon, display name, and description; because Telegram is active in the catalog, the row shows a **Connect** button instead of a "Coming soon" badge. Clicking Connect opens the shared Connect modal, which fetches Telegram's config field definitions from the provider catalog API and renders a form (API key as a masked password field). The admin submits the form; the Portal saves a workspace integration with `channel_type: chat`, `provider_name: telegram`, status **connected**, and the API key stored encrypted—same upsert path Mailgun uses. On success the row shows a **Connected** badge and a **Disconnect** button; invalid credentials show an error in the modal without changing integration status.

**Why this priority**: Enables self-service credential setup without manual database or config edits.

**Independent Test**: Open Integrations → Chat, click Connect on Telegram, enter a valid API key in the masked field, submit, and confirm the provider shows as Connected.

**Acceptance Scenarios**:

1. **Given** Telegram is an active provider in the catalog, **When** an admin opens the Chat section on Integrations, **Then** Telegram appears with its icon, name, and a Connect button (not "Coming soon").
2. **Given** the Connect modal is open for Telegram, **When** the admin views the form, **Then** the API Key field is masked (password-style), the description states this is for OTP verification (not a chat bot), and dismissal is via a Cancel button (no header close/X control).
3. **Given** valid Telegram credentials, **When** the admin submits Connect, **Then** the integration is saved encrypted and Telegram shows as Connected with Disconnect available.
4. **Given** invalid or rejected credentials, **When** the admin submits Connect, **Then** an error is shown in the modal and the integration is not marked Connected.

---

### User Story 3 - Batch OTP delivery under load (Priority: P3)

When multiple OTP sends are requested at once (e.g., bulk or concurrent API calls), the gateway throttles outbound Telegram requests so no more than five are in flight at a time.

**Why this priority**: Protects against provider rate limits and keeps delivery predictable under burst traffic.

**Independent Test**: Submit 20 chat send requests concurrently; observe that at most five provider calls run simultaneously and all complete successfully when credentials and numbers are valid.

**Acceptance Scenarios**:

1. **Given** 20 valid OTP send requests arrive concurrently, **When** the gateway processes them, **Then** sends proceed in waves of at most five concurrent provider calls until all are handled.
2. **Given** a batch of five is in progress, **When** one send fails, **Then** remaining requests in that wave still complete and failures are reported per message without blocking unrelated valid sends in later waves.

---

### User Story 4 - Provider-specific options via metadata (Priority: P3)

An integrator passes Telegram-specific options (such as `sender_username`) through the Chat message `metadata` field without changing the public Chat API shape.

**Why this priority**: Supports optional Telegram Gateway features while keeping the channel contract stable.

**Independent Test**: Send a chat message with `metadata.sender_username` set; confirm the value is forwarded to Telegram on the outbound request.

**Acceptance Scenarios**:

1. **Given** a chat request with `metadata` containing `sender_username`, **When** the Telegram provider sends the OTP, **Then** the sender username is included in the provider request.
2. **Given** a chat request with no `metadata` and a non-empty `message`, **When** the Telegram provider sends the OTP, **Then** only required fields (phone number and code) are sent.
3. **Given** a chat request with empty `message` and `metadata.code_length` set, **When** the Telegram provider sends the OTP, **Then** `code_length` is forwarded to Telegram and no `code` field is included (Telegram generates the OTP).

---

### Edge Cases

- What happens when the recipient phone number is valid E.164 but not registered on Telegram? The provider returns a delivery failure; the gateway surfaces a non-success result with a safe, generic client message.
- What happens when the workspace has no Telegram integration or Telegram is deactivated? The send fails with an integration/configuration error before calling the provider.
- What happens when the Telegram API key is missing or invalid? Connect fails in Portal; sends fail with an authentication/configuration error.
- What happens when `to` contains multiple phone numbers in one request? Each recipient receives the same OTP code (or a Telegram-generated code per send when `message` is empty and `code_length` is set); phone validation and E.164 normalization apply per recipient; batch concurrency rules apply across all provider calls spawned by the request.
- What happens when `message` is empty but `metadata.code_length` is set? Telegram generates a random OTP of the requested length; the gateway does not require or validate a client-supplied code for that request.
- What happens when both `message` and `metadata.code_length` are set? The explicit `message` value takes precedence; `code_length` is ignored (Telegram Gateway behavior when `code` is provided).
- What happens when optional metadata contains unsupported or malformed keys? Unsupported keys are ignored; malformed values for known keys cause validation failure before the provider call.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST register Telegram as an **active** Chat channel provider in the provider catalog (replacing the current "coming soon" / not_supported state).
- **FR-002**: System MUST send OTP verification messages through Telegram using the recipient's phone number as the delivery target.
- **FR-003**: System MUST validate each recipient phone number and reject sends when the number cannot be parsed or is not a valid number.
- **FR-004**: System MUST normalize valid recipient phone numbers to **E.164** format before calling Telegram.
- **FR-005**: When `message` is non-empty, system MUST map the Chat API `message` field to the Telegram provider's `code` field on outbound OTP requests.
- **FR-006**: When `message` is non-empty, system MUST accept OTP codes only as fully numeric strings between 4 and 8 characters in length (inclusive); otherwise the request is rejected.
- **FR-006a**: When `message` is empty and `metadata.code_length` has a valid value (4–8), system MUST omit `code` on the provider request and forward `code_length` so Telegram generates a random OTP of that length.
- **FR-006b**: When `message` is empty and `metadata.code_length` is missing or invalid, system MUST reject the request with a clear validation error.
- **FR-007**: System MUST pass Telegram-specific optional parameters (including `sender_username` and other documented Gateway options) from the Chat message `metadata` map to the provider request.
- **FR-008**: System MUST limit concurrent outbound Telegram OTP requests to **five in flight** at a time; additional requests wait until a slot in the current wave is free, processing in sequential waves until all are complete.
- **FR-009**: System MUST store Telegram connection credentials (API key) encrypted per workspace integration, consistent with other providers.
- **FR-010**: Portal Integrations MUST list Telegram under the **Chat** section the same way Mailgun is listed under Email: catalog-driven row (icon, name, description from `providers` seed), **Connect** when unconnected and active, **Connected** badge plus **Disconnect** when integrated, and **Coming soon** badge only when catalog status is not active.
- **FR-011**: The Telegram Connect modal MUST include a masked API Key field, explain that the integration is for **OTP verification messages** (not a general chat bot), and use a **Cancel** button for dismissal without a header close (X) control.
- **FR-012**: System MUST support configuring Telegram for embedded SDK usage via static provider config (alongside Portal-stored credentials for server mode).
- **FR-013**: All documentation that describes channel types, provider catalogs, configuration, and integration flows MUST be updated to include Telegram OTP for Chat.

### Key Entities

- **Telegram Provider (catalog)**: Chat-channel provider entry with name `telegram`, active status, icon, and OTP-focused description.
- **Telegram Integration (workspace)**: Workspace-scoped connection storing encrypted API key and integration status (connected / disconnected).
- **Chat OTP Request**: Outbound message with `to` (phone numbers), `message` (explicit OTP code, optional when `metadata.code_length` is set), optional `metadata` (e.g., `code_length`, `sender_username`, `ttl`, `payload`, `callback_url`).
- **Send Result**: Per-message outcome including provider request identifier and delivery status, consistent with other Chat providers.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A workspace admin can connect Telegram from the Portal Chat section in under 2 minutes without operator assistance.
- **SC-002**: 95% of valid single-recipient OTP send requests (correct phone + either a 4–8 digit `message` or empty `message` with valid `metadata.code_length` + active integration) complete with a success response on the first attempt under normal provider availability.
- **SC-003**: When 20 concurrent valid OTP sends are submitted, all complete with no more than five simultaneous provider calls observed at any instant.
- **SC-004**: 100% of documentation pages that enumerate Chat providers or integration setup steps list Telegram OTP with accurate status (available) and setup instructions.
- **SC-005**: Invalid phone numbers, invalid explicit OTP codes, and empty `message` without a valid `metadata.code_length` are rejected before any Telegram API call in 100% of tested negative scenarios.

## Assumptions

- Telegram **Gateway API** (verification messages to phone numbers) is the delivery mechanism; general Telegram Bot API chat messaging is out of scope.
- OTP **verification** (checking whether the user entered the correct code via `checkVerificationStatus`) is out of scope for this feature; only **sending** OTP codes is in scope.
- The Connect flow requires only an **API key** (Gateway token) at integration time; optional per-message parameters live in `metadata`.
- Phone numbers may arrive in local or international formats; normalization uses standard libphonenumber-style parsing with a default region only when the spec/plan defines one (otherwise numbers must include country code or `+` prefix).
- Existing Chat channel API contract (`POST /v1/chat`, Portal send-test for `chat`) remains unchanged; Telegram is a new provider implementation behind the same contract.
- Telegram Connect reuses the existing Integrations page and Portal integration API—the same path Mailgun uses today: provider config fields are defined in the catalog, the Connect modal renders them dynamically (API key as a masked field), and submitted credentials are stored encrypted per workspace integration.
- Memory provider remains the default Chat provider for local development unless a workspace explicitly connects Telegram.

## Out of Scope

- WhatsApp or other Chat providers beyond catalog status already seeded.
- Telegram delivery webhooks / `callback_url` processing inside the gateway (metadata may forward `callback_url` to Telegram, but the gateway does not host webhook handlers in this feature).
- Portal UI changes beyond Telegram Connect availability, description, and modal dismiss behavior.
