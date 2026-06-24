# Research: Data Retention Modes

## R1 — Unify portal retention setting with gateway dispatch mode

**Decision**: Map portal `data_retention` values to `message_dispatch_mode` in `internal/core/domain/dispatch_mode.go` with bidirectional helpers; settings PATCH writes both keys atomically (or a single canonical key with read-time alias support).

**Rationale**: Frontend already persists `data_retention`; gateway reads `message_dispatch_mode`. A single domain mapping table avoids drift.

**Alternatives considered**:
- Frontend writes `message_dispatch_mode` directly — rejected; breaks existing settings API contract.
- Separate retention service — rejected; YAGNI for four enum values.

## R2 — Provider + Database dispatch semantics

**Decision**: Add `DispatchProviderAndDatabase` (`provider_and_database`) that reuses the `DispatchProviderOnly` code path in `gateway_service.dispatch` (extract shared provider-send helper).

**Rationale**: User requires identical outbound behavior to Provider Only; only request logging differs. Sharing dispatch logic prevents behavioral drift.

**Alternatives considered**:
- Duplicate switch case — rejected; violates DRY.
- Log gating solely from retention string in handler — acceptable but dispatch mode enum is cleaner for tests.

## R3 — Request log gating location and success-only rule

**Decision**: Gate `SendHelper.RecordLog` (presentation layer) using resolved dispatch mode — log only for `memory_and_provider` and `provider_and_database`, and only on the success path (after `send(ctx)` returns without error). Remove `RecordLog` calls from validation, auth, and dispatch-error branches when retention would otherwise allow logging.

**Rationale**: `RecordLog` is centralized; gating here avoids touching every channel handler. Success-only rule matches 2026-06-24 clarification and reduces noise in audit table.

**Alternatives considered**:
- Gate inside `GatewayService.RecordLog` — viable but mixes transport audit with domain dispatch.
- Log failures too — rejected per spec FR-008/FR-010.

## R4 — Memory Only / Provider Only database persistence

**Decision**: No schema changes. Stop writing `message_request_logs` for `memory_only` and `provider_only`. Confirm inbox writer for memory modes remains in-process only (no Postgres message tables for memory capture).

**Rationale**: Spec explicitly removes DB persistence for these modes; current unconditional `RecordLog` insert is the fix target.

## R5 — Legacy retention value migration

**Decision**: Read-time aliases only (`both` → `memory_database`, `providers` → `provider`); PATCH always writes canonical values. No SQL migration required.

**Rationale**: `workspace_settings` is key-value; alias normalization in domain layer is sufficient per spec FR-003.

## R6 — Keep DispatchMemoryAndProvider enum name

**Decision**: Retain `DispatchMemoryAndProvider` (`memory_and_provider`) as the gateway dispatch mode for portal `memory_database` retention. Do not rename to `memory_and_database`. Map `memory_database` ↔ `memory_and_provider` in domain helpers only.

**Rationale**: Existing code and tests already use `DispatchMemoryAndProvider`; renaming adds churn without functional benefit. Portal uses user-facing `memory_database`; gateway keeps stable runtime enum.

**Alternatives considered**:
- Rename to `memory_and_database` — rejected per product decision (2026-06-24).
