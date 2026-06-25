# Research: Data Retention Modes

## R1 — Unify portal retention setting with gateway dispatch mode

**Decision**: Map portal `data_retention` values to `message_dispatch_mode` in `internal/core/domain/dispatch_mode.go` with bidirectional helpers; settings PATCH writes both keys atomically (or a single canonical key with read-time alias support).

**Rationale**: Frontend already persists `data_retention`; gateway reads `message_dispatch_mode`. A single domain mapping table avoids drift.

**Alternatives considered**:
- Frontend writes `message_dispatch_mode` directly — rejected; breaks existing settings API contract.
- Separate retention service — rejected; YAGNI for four enum values.

## R2 — Provider + Database dispatch semantics

**Decision**: Add `DispatchProviderAndDatabase` (`provider_and_database`) that reuses the `DispatchProviderOnly` code path in `gateway_service.dispatch` (extract shared provider-send helper).

**Rationale**: User requires identical outbound behavior to Provider Only; only the `retained` flag on request logs differs. Sharing dispatch logic prevents behavioral drift.

**Alternatives considered**:
- Duplicate switch case — rejected; violates DRY.

## R3 — Operational logging vs retention (Idea 3)

**Decision**: Use a single table `message_request_logs` with a `retained BOOLEAN` column. Always insert on successful send (all modes) for **Recent Requests** / monitoring. Set `retained = true` only when `message_dispatch_mode` is `memory_and_provider` or `provider_and_database`.

**Rationale**: Gating inserts by retention mode broke Recent Requests, which reads from `message_request_logs`. Operational visibility and long-term retention are different concerns; one row with a flag avoids duplicate tables and duplicate writes.

**Alternatives considered**:
- Gate inserts by mode (no rows for memory/provider only) — rejected; breaks Recent Requests.
- Second table for retained rows only — viable (Idea 2) but more schema and dual-write complexity.
- Recent Requests from in-memory source — rejected; no durable source for provider sends.

## R4 — Message content vs request metadata

**Decision**: Retention policy continues to gate **message content** (inbox/DB) by mode. Request metadata always logs operationally; `retained` flag gates long-term policy storage.

**Rationale**: Memory only and Provider only must not persist message bodies; operational request rows with `retained = false` satisfy portal monitoring without long-term retention commitment.

## R5 — Legacy retention value migration

**Decision**: Read-time aliases only (`both` → `memory_database`, `providers` → `provider`); PATCH always writes canonical values. Normalize on `GetSettings` in portal service.

**Rationale**: `workspace_settings` is key-value; alias normalization in domain layer is sufficient per spec FR-003.

## R6 — Keep DispatchMemoryAndProvider enum name

**Decision**: Retain `DispatchMemoryAndProvider` (`memory_and_provider`) as the gateway dispatch mode for portal `memory_database` retention.

**Rationale**: Existing code and tests already use `DispatchMemoryAndProvider`; renaming adds churn without functional benefit.

## R7 — Non-retained row lifecycle

**Decision**: v1 ships with `retained` column only; optional TTL purge job for `retained = false` rows deferred unless product requests fixed retention window for operational logs.

**Rationale**: Purge policy (7 vs 30 days) is a product/ops decision; schema supports it without blocking MVP.
