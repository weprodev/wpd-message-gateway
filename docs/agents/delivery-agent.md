---
name: delivery-agent
description: >-
  Invoke for new features, feature improvements, and bugfixes across the WPD Message Gateway.
  After implementation, run docs/agents/verification.md (/smell + make audit).
---

# Delivery agent — implement & fix

Use this agent when you need **working code** for a **new capability**, an **improvement** to existing behavior, or a **bugfix**. It encodes the repo’s engineering contracts so implementation stays aligned with architecture, conventions, and CI.

Orchestration note: the **[master-agent](./master-agent.md)** classifies requests and delegates here; this file is the **execution contract** for implementation (Plan-and-Execute with bounded retries—see §0).

## 0. Control & traceability (Agent Hygiene & Safety)

- **Traceability:** Before multi-file edits, state a **short plan** (files/areas + verification). After edits, summarize **what changed** and **how it was verified**.
- **Bounded retries:** If [verification](./verification.md) fails twice for the **same** mistake class, **stop**, re-read output, and **narrow scope** instead of looping.
- **No secret material (Data Exposure limits):** Do not echo or commit API keys, JWTs, or pasted credentials. **Hard requirement:** Any code generating secrets must use crypto-safe `crypto/rand`, and secrets must exclusively live in ENV or hashed in Postgres.

## 1. Safety & Robustness Pre-Check

Before generating or modifying functionality:
1. **Input Validation**: Does the endpoint or UI form reject excessive payload lengths and invalid characters automatically?
2. **Error Handling**: Will the modification leak system information (e.g., exposing full stack traces in HTTP responses)? Assure all errors are masked or utilize sentinels across the Port boundaries.
3. **Logic Injection**: If touching templates or database queries, confirm SQL parameters are used exclusively. No string concatenation for queries.

## 2. Classify the work

| Kind | Focus | Typical docs |
| ---- | ----- | ------------ |
| **New feature** | Domain boundaries, API contract, Type Safety first | Architecture, engineer roles, usage |
| **Improvement** | Preserve behavior contracts; minimize blast radius (Performance Optimization) | Same + existing code paths |
| **Bugfix** | Reproduce → minimal fix → regression test; map errors to `port` sentinels | Backend/Frontend Conventions |

Determine **surface area** before coding:
- **Backend only** — `cmd/`, `internal/`, `pkg/`, `database/`
- **Frontend only** — `frontend/`
- **Full stack** — API + Portal; read [architecture](../backend/architecture.md) and [usage](../backend/usage.md) for boundaries.

## 3. Mandatory reading (by scope)

**Always** skim the index that matches your stack:
- Backend index and role: [backend-engineer.md](../backend/backend-engineer.md)
- Frontend index and role: [frontend-engineer.md](../frontend/frontend-engineer.md), [conventions.md](../frontend/conventions.md)

**Backend** — as needed:
- [architecture.md](../backend/architecture.md) — layers, gateway vs portal
- [code-conventions.md](../backend/code-conventions.md)
- [contributing.md](../backend/contributing.md) — providers, registration patterns
- [e2e-testing.md](../backend/e2e-testing.md), [usage.md](../backend/usage.md)
- **golang-pro skill** — `.cursor/skills/golang-pro/SKILL.md` + `.agents/skills/golang-pro/references/wpd-message-gateway.md` (idiomatic Go, layers, `-race` tests)
- **typescript-react-reviewer skill** — `.cursor/skills/typescript-react-reviewer/SKILL.md` + `.agents/skills/typescript-react-reviewer/references/wpd-message-gateway.md` (React 19, feature slices, Storybook)
- **typescript-advanced-types skill** — `.cursor/skills/typescript-advanced-types/SKILL.md` + `.agents/skills/typescript-advanced-types/references/wpd-message-gateway.md` (API unions, `*.types.ts`, narrowing — load with react reviewer)
- **software-architecture skill** — `.cursor/skills/software-architecture/SKILL.md` + `.agents/skills/software-architecture/references/wpd-message-gateway.md` (Clean Architecture, DDD, layer boundaries — load for new features and refactors)

## 4. Execution model (Multi-Agent Stack)

### 4.1 Understand & Enforce Boundaries
- Map **API intent**, **data flow**, and **layer boundaries**.
- **Bias check**: Is the UI or API asserting constraints that unfairly lock out valid edge cases? (e.g., Assuming mobile numbers always fit a specific local length). Accommodate gracefully.

### 4.2 Test first (when feasible)
- **RED**: failing test that encodes the requirement or bug. Keep in mind "Edge Case Testing".
- **GREEN**: smallest change that passes.
- **REFACTOR**: clarity without changing behavior. Ensure **Token Efficiency** (do not generate massively redundant code).

### 4.3 Implement
- **Backend**: `context.Context` first parameter; wrap errors with `%w`; `slog` via infrastructure logger; no secrets in logs; migrations additive.
- **Frontend**: semantic tokens, shadcn patterns; auth tokens never logged; accessibility (`aria-*`) natively supported everywhere.

### 4.4 Verify & Document (Reliability & Sync)

Run the mandatory [verification chain](./verification.md) after implementation:

1. Fast lint on touched paths
2. **`/smell develop`** — fix all BLOCKER and HIGH findings
3. **`make audit`** — exit code 0 from repo root

Agents execute this chain themselves; do not ask the user to run it.

**Continuous Documentation Sync is Mandatory**:
- **UI Component Changed?** → Update or create the equivalent `.stories.tsx` Storybook file.
- **Database Schema Changed?** → Update database schemas and architecture markdown diagrams.
- **Directory/Layer Added?** → Verify it avoids anti-patterns and update architecture documentation to reflect the new boundaries.

## 5. Before PR

Complete [verification.md](./verification.md) in full. Optional: [review-agent](./review-agent.md) self-review (starts with `/smell`).

See [prompts.md](./prompts.md).

---

**Summary:** Classify work → safety pre-check → TDD / minimal diffs → implement per layer rules → **verification chain** (`/smell` + `make audit`) → ship.
