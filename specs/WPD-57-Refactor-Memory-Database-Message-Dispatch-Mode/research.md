# Research: Refactor Memory + Database Message Dispatch Mode

**Feature**: WPD-57-Refactor-Memory-Database-Message-Dispatch-Mode  
**Date**: 2026-06-29

## R1 — Correct dispatch semantics for Memory + Database

**Decision**: **Memory + Database** uses the same outbound path as **Memory only** — `writeToInbox` only; no `activeIntegration` lookup and no `sendViaProvider` call.

**Rationale**: WPD-33 implemented `memory_and_provider` as memory capture **plus** provider send, which contradicts the product label "Memory + Database". Users expect local capture with database-backed request log retention, not a dual send. Reusing the `memory_only` branch eliminates provider dependency and matches FR-001/FR-002.

**Alternatives considered**:
- Keep provider send but add DB message persistence — rejected; inbox is in-process RAM only, and spec excludes message-body DB persistence.
- New fifth mode — rejected; correct existing mode instead of proliferating options.

## R2 — Retained flag for Memory + Database

**Decision**: Set `retained = true` for gateway mode `memory_and_database` (via `ShouldRetainRequestLog`).

**Rationale**: Spec defines "saved in database" as request log retention, aligned with **Provider + Database**. Two modes now share long-term retention semantics; operational modes (`memory_only`, `provider_only`) stay `retained = false`.

**Alternatives considered**:
- Keep `retained = false` (WPD-33) — rejected; explicitly superseded by this feature.
- Separate retention column per mode — rejected; existing `retained` boolean is sufficient.

## R3 — Naming migration matrix

**Decision**: Apply pattern-preserving renames:

| Layer | Before | After |
| ----- | ------ | ----- |
| Setting value | `both`, `memory_provider` | `memory_database` (read alias only for legacy) |
| Gateway string | `memory_and_provider` | `memory_and_database` |
| Go constant | `DispatchMemoryAndProvider` | `DispatchMemoryAndDatabase` |
| Portal TS union | `memory_and_provider` | `memory_and_database` |

**Rationale**: Names must reflect behavior (memory + database retention) not provider dispatch. Pattern `memory_and_*` → `memory_and_database` mirrors `provider_and_database`.

**Alternatives considered**:
- Rename setting to `memory_and_database` — rejected; portal/API already standardized on short values (`memory_database`) in WPD-33.
- Big-bang DB migration rewriting all `both` rows — rejected; out of scope; gateway maps on read.

## R4 — Legacy alias strategy

**Decision**: Read-only aliases at gateway parse/mapping layer:
- Setting `both` → `memory_and_database` behavior (existing)
- Gateway string `memory_and_provider` → `memory_and_database` (new alias in `ParseMessageDispatchMode` / `SettingValueToDispatchMode`)
- New writes use canonical values only; no automatic DB rewrite on GET/PATCH

**Rationale**: Clarification session 2026-06-29 — avoids breaking stored settings, historical log metadata (`dispatch_mode` in result meta), and Bruno collections during rollout.

**Alternatives considered**:
- Hard reject `memory_and_provider` — rejected; breaks existing deployments.
- Migrate-on-read rewriting DB — rejected; out of scope per spec.

## R5 — Provider name on request logs

**Decision**: Under `memory_and_database`, persist `provider_name = "memory"` on request logs (same as `memory_only`).

**Rationale**: `ProviderNameForLog` already returns `memory` for `DispatchMemoryOnly`; extend the same branch for `DispatchMemoryAndDatabase`. No integration is consulted.

**Alternatives considered**:
- Empty `provider_name` — rejected; inconsistent with memory-only rows and complicates Recent Requests display.

## R6 — WPD-33 artifact handling

**Decision**: WPD-57 artifacts are authoritative for `memory_database` behavior; WPD-33 `data-model.md` / `contracts/` remain historical. Implementation updates live code and `docs/backend/*` only.

**Rationale**: WPD-33 directory name differs (`WPD-33-Add-Data-Retention-Config`); no need to rewrite shipped spec — this feature documents the correction explicitly.

**Alternatives considered**:
- Amend WPD-33 spec in place — rejected; separate feature branch and spec directory per user request.
