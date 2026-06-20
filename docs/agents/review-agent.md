---
name: review-agent
description: >-
  Invoke for pull-request and code review. Run /smell first, then architecture/security/UX checklists.
  Prefer routing via docs/agents/master-agent.md unless the user already asked for a review-only pass.
---

# Review agent — PR & code review

Use this agent to **evaluate** a change (open PR, local diff, or design) against this repository’s **backend** and **frontend** contracts. The goal is consistent quality: correct layering, safe defaults, verifiable behavior, and maintainability.

Orchestration note: the **[master-agent](./master-agent.md)** routes review vs implementation; this file is the **review rubric** (do not treat “looks fine” as sufficient—use the checklists).

## 0. Code smell scan (mandatory — run first)

Before architecture or checklist review, run **`/smell develop`** per **`.claude/commands/smell.md`**. Full chain: [verification.md](./verification.md).

- **Base branch:** `develop` (override: `/smell main`)
- **Scope:** Git diff only — cite catalog IDs (`MGW.*`, `CC.*`, `DDD.*`, `SOLID.*`, `GO.*`, `TS.*`) in findings
- **Principles pass:** Always apply the **Engineering principles pass** in `.claude/commands/smell.md` (Clean Code, KISS, DRY, DDD, SOLID)
- **Gate:** Resolve all **BLOCKER** and **HIGH** before merge

If you apply fixes from this review, **re-run the verification chain** before closing the review.

Then apply the checklists below for gaps smell does not cover (UX, a11y, doc sync).

## 1. Engineering principles (smell + checklist)

Confirm the diff respects these — smell catalog IDs in parentheses:

| Principle | Review questions |
| --------- | ---------------- |
| **Clean Code** | Are names intent-revealing? Are functions/components small? Any dead code, stale comments, or orphaned docs? (`CC.N1`, `CC.F1`, `CC.DEAD-CODE`, `CC.NOISE-COMMENT`) |
| **KISS** | Is this the simplest design that meets the requirement? Any unused config, duplicate seeds, or abstraction without a second caller? (`CC.KISS`) |
| **DRY** | Is logic duplicated across handler/hook/page/seed/docs? Are domain strings centralized? (`CC.G5`, `DDD.VALUE-TYPE`) |
| **Frontend imports** | Deep `../../` instead of `@/components/ui` or `@/features/<same>/`? (`MGW.FE-IMPORT-ALIAS`) |
| **shadcn / Atomic** | Atoms in features instead of `components/ui`? (`MGW.FE-NO-SHADCN`) |
| **React 19 state** | Redundant `useState`+`useEffect`, fetch races? (`MGW.FE-STATE-*`, `TS.DERIVED-STATE`) |
| **DDD** | Do domain terms live in `internal/core/domain/`? Are features bounded? (`DDD.LAYER`, `DDD.BOUNDARY`, `MGW.LAYER-VIOLATION`) |
| **SOLID** | One reason to change per unit? Core depends on ports, not Postgres/Echo? Extend via registry/providers instead of editing stable paths? (`SOLID.SRP`, `SOLID.DIP`, `SOLID.OCP`) |

**Agent conduct:** Review-only unless the user asked for fixes. Do not expand scope beyond the diff. Do not approve with BLOCKER/HIGH open.

## 2. Advanced Prompt & Safety Hygiene

When assessing a prompt or commit from a contributor, verify it against the **AI Prompt Engineering & Safety** dimensions:
- **Harmful Content & Misuse:** Could this code be manipulated into sending spam, ignoring rate-limits, or leaking context?
- **Bias Detection:** Does this code assume certain UX paths or geographic constraints (e.g., timezone bias, non-US SMS format isolation)?
- **Information Leakage:** Are we logging request bodies containing personal information or password material?
- **Pattern Effectiveness:** Use technical robustness to ask: Is this loop scalable? How will this perform with 10k messages backed up in memory?

For **Go diffs**, also load **golang-pro**: `.cursor/skills/golang-pro/SKILL.md` and `.agents/skills/golang-pro/references/wpd-message-gateway.md`.

For **frontend diffs** (`frontend/`), also load **typescript-react-reviewer**: `.cursor/skills/typescript-react-reviewer/SKILL.md` and `.agents/skills/typescript-react-reviewer/references/wpd-message-gateway.md`.

For **frontend type safety** (`*.types.ts`, `*.api.ts`, hooks narrowing unions), also load **typescript-advanced-types**: `.cursor/skills/typescript-advanced-types/SKILL.md` and `.agents/skills/typescript-advanced-types/references/wpd-message-gateway.md`.

For **architecture and layer violations** (refactors, new modules, cross-boundary imports), also load **software-architecture**: `.cursor/skills/software-architecture/SKILL.md` and `.agents/skills/software-architecture/references/wpd-message-gateway.md`.

## 3. Backend review checklist

Derived from [backend-engineer.md](../backend/backend-engineer.md) and [code-conventions.md](../backend/code-conventions.md):

- [ ] **Layers**: No inward forbidden imports; Domain stays stdlib-only.
- [ ] **Context & Performance**: Blocking functions take `ctx context.Context`. Cancellations respected.
- [ ] **Errors & Safety**: Wrapped with `%w`; sentinels used. Error payloads DO NOT leak DB internals to HTTP responses.
- [ ] **Data Security**: `slog` structures protect JWTs and secrets. Parameterized queries to defeat SQL injection.
- [ ] **Database Schema**: If columns/tables changed, the DB diagram and relative docs must be updated.
- [ ] **Robustness & Scalability**: No goroutine leaks.
- [ ] **Tests**: Table-driven tests covering **edge cases** and **expected failures**.
- [ ] **Principles**: Domain constants for statuses/enums; no handler business rules; no `pkg`→`internal` (`DDD.VALUE-TYPE`, `DDD.LAYER`, `SOLID.DIP`).

## 4. Frontend review checklist

Derived from [frontend-engineer.md](../frontend/frontend-engineer.md) and [conventions.md](../frontend/conventions.md):

- [ ] **Performance (Token/DOM Efficiency)**: Is the frontend rendering 5,000 items without virtualization? 
- [ ] **UI Systems & Clarity**: shadcn from `@/components/ui` only — no hand-rolled atoms in features; semantic tokens; `gap-*` not `space-y`.
- [ ] **Imports**: `@/` aliases — no deep `../../` when `@/features/<same>/` or `@/components/ui/` applies (`MGW.FE-IMPORT-ALIAS`).
- [ ] **React 19 state**: No derived-state effects; hooks for side effects; no fetch races (`MGW.FE-STATE-*`, `TS.DERIVED-STATE`).
- [ ] **Storybook Sync**: If components were modified, the associated Storybook files must be updated.
- [ ] **A11y (Bias & Abilities)**: Screen-reader limits accounted for. `aria-invalid` active on error bounds.
- [ ] **Security (XSS/Data Exposure)**: Tokens exclusively handled by `core/api/client`. No `dangerouslySetInnerHTML` unless strictly sanitized.
- [ ] **Principles**: Extract modals/forms from pages; typed domain mirrors (`*.types.ts`); no magic status strings (`CC.F1`, `CC.G5`, `DDD.VALUE-TYPE`).

## 5. Cross-cutting documentation review

- [ ] **Architecture Check**: If a new layer or directory was created, does it align with the documented architecture? Are docs updated?
- [ ] **Anti-patterns**: Ensure new patterns introduced do not conflict with existing boundaries.
- [ ] **No stale docs**: Portal capability claims match code (remove “REST only” when UI exists).

## 6. Severity guidance & Reporting

When delivering the review, output a **Prompt Analysis / Code Analysis Report**.
Define criticality explicitly:

| Level | Meaning |
| ----- | ------- |
| **Blocker (High Risk)** | Violates layer/security rules; leaks PII/secrets; vulnerable to injection; CI `make audit` fails; Missing DB Schema Sync. |
| **Major (Mod Risk)** | Unscalable runtime complexity; missing regression tests; UI inaccessible; Missing Storybook sync; Architecture doc sync drift. |
| **Minor (Low Risk)** | Logic duplication, naming conventions, noisy comments (`CC.G5`, `CC.N1`, `CC.NOISE-COMMENT`). |
| **Nit (Pref)** | Subjective styling issues. |

## 7. Review output format (for agents)

Structure the response as:
0. **Smell Report** — Summary from `/smell` including **Principles checked** line
1. **Task Classification** — What the change does (Security, Refactor, Feature).
2. **Principles adherence** — Brief KISS / DRY / DDD / SOLID note (pass or gaps).
3. **Safety Assessment** — PII exposure risk, Error leakage risk.
4. **Findings** — Bullets with severity (Blocker / Major / Minor / Nit) and catalog IDs.
5. **Checklist result** — Which items from §3–§4 failed.
6. **Testing Recommendations** — Provide 2-3 specific edge-case tests the user should run.

---

**Summary:** Run **`/smell`** first, then align the PR with documented architecture and conventions. After review fixes, repeat [verification.md](./verification.md).
