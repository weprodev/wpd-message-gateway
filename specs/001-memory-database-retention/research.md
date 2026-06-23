# Research: Align Memory + Database Data Retention

## Decision 1: Rename dispatch mode constant

**Decision**: Rename `DispatchMemoryAndProvider` (`memory_and_provider`) to `DispatchMemoryAndDatabase` (`memory_and_database`).

**Rationale**: Portal UI already labels this mode "Memory + Database" and stores `memory_database`. Backend naming drift caused incorrect implementation (provider send instead of DB persist).

**Alternatives considered**:
- Keep `memory_and_provider` internally, only fix logic — rejected; perpetuates naming mismatch in logs, metadata, and docs.
- Remove dispatch mode entirely, use only `data_retention` — rejected; dispatch mode is the runtime routing key throughout `GatewayService`.

## Decision 2: Memory + Database dispatch behavior

**Decision**: `memory_and_database` writes to RAM inbox **and** `stored_messages`; does **not** call any provider.

**Rationale**: Matches Portal description ("Persist messages in the portal inbox and database") and user requirement. Provider + Database pattern reused only for **persistence** semantics (fail-closed on archive write failure).

**Alternatives considered**:
- RAM + provider (current broken behavior) — rejected.
- RAM + DB + provider — rejected; contradicts spec and frontend copy.
- DB only (no RAM) — rejected; Portal inbox requires RAM capture for preview.

## Decision 3: Remove provider-named symbols (no runtime aliasing)

**Decision**: Rename all `memory_and_provider` / `memory_provider` code identifiers to `memory_and_database` / `memory_database`. Do not keep provider-named constants or functions with normalization shims.

**Rationale**: Developers reading `memory_and_provider` assume provider dispatch occurs; renaming eliminates confusion. Stale DB value `both` is fixed in the existing init migration, not a new migration file.

**Alternatives considered**:
- Keep `memory_and_provider` as parse alias with normalization — rejected; perpetuates misleading names in codebase.
- Accept `memory_and_provider` on API PATCH — rejected; only canonical `memory_database` accepted going forward.

## Decision 4: `message_dispatch_mode` key name unchanged; values updated

**Decision**: Keep the workspace setting key `message_dispatch_mode` as-is. Update value mappings and any related files where Memory + Database was represented as `memory_and_provider` or `both` — those references become `memory_and_database` / `memory_database`.

**Rationale**: User wants the key name preserved but all Memory + Database–related code and docs corrected. Renaming the key itself is unnecessary churn.

**Alternatives considered**:
- Leave all `message_dispatch_mode` code untouched — rejected; stale `memory_and_provider` mappings would remain.
- Rename the key to `data_retention` everywhere — rejected; key name stays `message_dispatch_mode` where it already exists alongside `data_retention`.

## Decision 5: stored_messages dispatch_status for memory_and_database

**Decision**: On successful Memory + Database capture, set `dispatch_status` = `sent` and `dispatched_at` via `RecordDispatchOutcome` (provider fields remain empty).

**Rationale**: User requires `sent` after successful send; `sent` here means capture-to-RAM-and-DB succeeded, not provider delivery.

**Alternatives considered**:
- Leave as `pending` — rejected per user clarification.
- New status `stored` — rejected; existing `sent` value is sufficient with empty provider fields.

## Decision 6: Error handling

**Decision**: Fail-closed if either inbox or archive write fails (no partial success metadata).

**Rationale**: Consistent with Provider + Database persistence failure behavior; prevents silent data loss.

**Alternatives considered**:
- Best-effort RAM if DB fails — rejected; spec requires durable storage for this mode.

## Decision 7: Frontend scope

**Decision**: No frontend changes.

**Rationale**: `settings.page.tsx` already uses `memory_database`; `settings.api.ts` already maps legacy `both` and `memory_and_provider`.
