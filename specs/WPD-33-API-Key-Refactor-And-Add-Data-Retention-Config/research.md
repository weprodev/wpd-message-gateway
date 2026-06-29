# Research: Message Dispatch Modes & Request Retention

**Feature**: WPD-33-API-Key-Refactor-And-Add-Data-Retention-Config  
**Date**: 2026-06-28

## R1 — Operational logging vs long-term retention

**Decision**: Single table (`message_request_logs`) with boolean `retained` column.

**Rationale**:
- Recent Requests already reads from `message_request_logs` via `ListWithSource` with no mode filter.
- Gating inserts by mode would break Recent Requests for operational-only workspaces.
- A flag allows future TTL purge of `retained = false` rows without schema churn.

**Alternatives considered**:
| Alternative | Rejected because |
| ----------- | ---------------- |
| Separate `retained_request_logs` table | Duplicates metadata columns; dual-write risk |
| Skip logging for memory-only modes | Breaks FR-010 / Recent Requests regression |
| Derive `retained` at query time from workspace setting | Historical rows wrong after mode change |

---

## R2 — Provider + Database dispatch behavior

**Decision**: New gateway mode `provider_and_database` executes identical code path to `provider_only`; only `retained` on the log row differs.

**Rationale**: Spec FR-007/FR-008 require identical outbound behavior; duplicating switch arms violates DRY — share provider dispatch helper or fall through same case with mode passed to logging layer.

**Alternatives considered**:
| Alternative | Rejected because |
| ----------- | ---------------- |
| Separate handler/service | Unnecessary duplication |
| Persist message content like Memory + Database | Contradicts spec (no content persistence) |

---

## R3 — Canonical naming (portal vs gateway)

**Decision**: Two-layer naming with explicit mapping in `internal/core/domain/dispatch_mode.go`:

| Portal / API setting value | Gateway `MessageDispatchMode` | `retained` |
| -------------------------- | ----------------------------- | ---------- |
| `memory` | `memory_only` | `false` |
| `memory_database` | `memory_and_provider` | `false` |
| `provider` | `provider_only` | `false` |
| `provider_database` | `provider_and_database` | `true` |

Legacy read aliases: `both` → `memory_database`, `providers` → `provider`.

**Rationale**: Portal uses short user-facing values; gateway uses explicit dispatch semantics. Single mapping module prevents drift. Gateway mode `memory_and_provider` is unchanged for this feature.

**Alternatives considered**:
| Alternative | Rejected because |
| ----------- | ---------------- |
| One string everywhere | Portal already uses `memory`/`both`/`providers`; breaking change |

---

## R4 — Settings key: `message_dispatch_mode`

**Decision**: Backend, Portal, and API use `message_dispatch_mode` as the sole setting key. Portal UI PATCH/GET uses this key with canonical portal values (`memory`, `memory_database`, `provider`, `provider_database`) or gateway strings (`memory_only`, `provider_and_database`, etc.). Legacy **values** (`both`, `providers`) are mapped on read via `SettingValueToDispatchMode`; no alternate setting key is supported.

**Rationale**: FR-003/SC-004; backend domain defines `SettingKeyMessageDispatchMode`. Single key avoids dual-source bugs between Portal and gateway.

**Alternatives considered**:
| Alternative | Rejected because |
| ----------- | ---------------- |
| Accept multiple setting keys | Perpetuates dispatch-mode resolution bugs |
| DB migration renaming legacy keys in `workspace_settings` | Not needed — greenfield uses `message_dispatch_mode` only |

---

## R5 — Schema change strategy

**Decision**: Add `retained BOOLEAN NOT NULL DEFAULT false` to `message_request_logs` in existing init migration `20260318000000_init_gateway.up.sql` (and down migration). No new timestamped migration file.

**Rationale**: Spec assumption; project appears pre-production with editable init migration.

**Alternatives considered**:
| Alternative | Rejected because |
| ----------- | ---------------- |
| New forward-only migration | Rework and adding extra files |

---

## R6 — When to set `retained` on insert

**Decision**: Set at insert time from dispatch mode resolved for that request (`ShouldRetainRequestLog(mode)`). Do not change which requests produce a log row.

**Rationale**: FR-010 — current `send_helper.DispatchAndLog` logs auth failures, bad JSON, dispatch errors, and successes. All continue to insert; `retained` reflects mode at dispatch time (false for operational-only modes even on error rows, unless product later defines error-specific rules — out of scope v1).

**Alternatives considered**:
| Alternative | Rejected because |
| ----------- | ---------------- |
| Success-only logging | Contradicts clarified spec (flow unchanged) |
| `retained` only on HTTP 200 | Mode distinction lost on failed provider attempts |

---

## R7 — SDK mode impact

**Decision**: No changes to `pkg/gateway` — embedded SDK has no DB or request-log pipeline.

**Rationale**: Feature is server/Portal scoped; SDK uses in-process config, not `workspace_settings`.
