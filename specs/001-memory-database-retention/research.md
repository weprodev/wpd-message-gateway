# Research: Data Retention Modes & API Key Modals

## R1 — Unify portal retention setting with gateway dispatch mode

**Decision**: Map portal `data_retention` values to `message_dispatch_mode` in `internal/core/domain/dispatch_mode.go` with bidirectional helpers; settings PATCH writes both keys atomically (or a single canonical key with read-time alias support).

**Rationale**: Frontend already persists `data_retention`; gateway reads `message_dispatch_mode`. A single domain mapping table avoids drift and matches existing `dispatch_mode.go` task scope.

**Alternatives considered**:
- Frontend writes `message_dispatch_mode` directly — rejected; breaks existing settings API contract and mixes user-facing naming with runtime keys.
- Separate retention service — rejected; YAGNI for four enum values.

## R2 — Provider + Database dispatch semantics

**Decision**: Add `DispatchProviderAndDatabase` (`provider_and_database`) that reuses the `DispatchProviderOnly` code path in `gateway_service.dispatch` (extract shared provider-send helper) but sets meta flag `persist_request_log=true` for the handler layer.

**Rationale**: User requires identical outbound behavior to Provider Only; only request logging differs. Sharing dispatch logic prevents behavioral drift.

**Alternatives considered**:
- Duplicate switch case — rejected; violates DRY.
- Log gating solely from retention string in handler — acceptable fallback but dispatch meta is cleaner for tests.

## R3 — Request log gating location

**Decision**: Gate `SendHelper.RecordLog` (presentation layer) using resolved dispatch mode — log only for `memory_and_database` and `provider_and_database`.

**Rationale**: `RecordLog` is already centralized; gating here avoids touching every channel handler and matches spec FR-008.

**Alternatives considered**:
- Gate inside `GatewayService.RecordLog` — viable but mixes transport audit with domain dispatch; handler already owns HTTP status mapping.

## R4 — Memory Only / Provider Only database persistence

**Decision**: No schema changes. Stop writing `message_request_logs` for `memory_only` and `provider_only`. Confirm inbox writer for memory modes remains in-process only (no Postgres message tables today for memory capture).

**Rationale**: Spec explicitly removes DB persistence for these modes; current `RecordLog` unconditional insert is the fix target.

## R5 — API key modal UX

**Decision**: New feature-scoped modal components under `frontend/src/features/settings/components/` using existing Radix `Dialog` primitives. Create modal uses default close **X**; regenerate/delete use `DialogContent` variant without close button (hide via `showClose={false}` prop or custom content wrapper).

**Rationale**: Matches user flows; reuses design-system dialog, button, input, and icon components.

**Alternatives considered**:
- Browser `prompt`/`confirm` — rejected by spec.
- Global modal provider — rejected; scope limited to settings page.

## R6 — Credentials one-time display

**Decision**: `CredentialsModal` receives `{ clientId, clientSecret }` from create/regenerate API responses; dismiss clears secret from React state (never store in localStorage).

**Rationale**: Aligns with existing API shapes (`client_secret` on create, regenerate returns `{ client_secret }` plus row `client_id`).

## R7 — Legacy retention value migration

**Decision**: Read-time aliases only (`both` → `memory_database`, `providers` → `provider`); PATCH always writes canonical values. No SQL migration required.

**Rationale**: `workspace_settings` is key-value; alias normalization in domain layer is sufficient per spec FR-003.

## R8 — Copy button feedback

**Decision**: Local `copiedField` state per field; `navigator.clipboard.writeText`; CSS `transition` on icon swap (`content_copy` → `check` Icon component).

**Rationale**: Meets smooth transition requirement without new dependencies.
