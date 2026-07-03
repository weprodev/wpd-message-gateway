---
description: Review git diff vs develop for code smells against WPD Message Gateway architecture and conventions
argument-hint: "[target-branch]"
allowed-tools: Bash(git:*), Read, Grep, Glob
---

# /smell — Pre-PR quality review

Structured review of **your changes** against this repository’s **architecture**, **engineering principles** (Clean Code, KISS, DRY, DDD, SOLID), **security**, and **conventions**.
Applies to **Go backend**, **React portal**, **embedded SDK** (`pkg/`), and **migrations**.

Run before opening a PR — after local lint/test, before or alongside `make audit`.

**Default base branch:** `develop` · override: `/smell main`

**Authoritative docs:** `docs/backend/architecture.md`, `docs/backend/code-conventions.md`, `docs/frontend/architecture.md`, `docs/frontend/conventions.md`, `docs/agents/review-agent.md`

---

## Step 1 — Ingest diff

!`bash -c '
BASE="$1"
if [ -z "$BASE" ]; then
  BASE="develop"
  git rev-parse "$BASE" >/dev/null 2>&1 || BASE="origin/develop"
  git rev-parse "$BASE" >/dev/null 2>&1 || BASE="main"
fi
if ! git rev-parse "$BASE" >/dev/null 2>&1; then
  echo "ERROR: base $BASE not found. Pass: /smell <branch>"
  exit 1
fi
echo "===== BASE: $BASE ====="
echo
echo "----- Stat (committed vs $BASE) -----"
git diff --stat "$BASE"...HEAD || true
echo
echo "----- Stat (working tree) -----"
git diff --stat HEAD || true
echo
echo "===== Committed diff (vs $BASE, -U10) ====="
git diff -U10 "$BASE"...HEAD || true
echo
echo "===== Working-tree diff (-U10) ====="
git diff -U10 HEAD || true
' -- "$ARGUMENTS"`

---

## Step 2 — Classify change

State one category with a one-sentence reason:

`feature` · `refactor` · `bugfix` · `test` · `docs` · `config` · `mixed`

State which surfaces the diff touches: **backend** · **frontend** · **pkg (SDK)** · **database** · **multiple**

---

## Step 3 — Apply review lens

Pick the **primary** lens and justify in one sentence:

- **Architecture** — layer boundaries, dual-mode gateway, feature-sliced frontend, DDD bounded contexts
- **Clean Code** — naming, function size, clarity, dead code, comment hygiene
- **Security** — auth, secrets, injection, PII in logs
- **Mixed** — when the diff spans more than one lens

Review **only what changed**. Do not demand unrelated refactors.

### Engineering principles pass (always run)

After the primary lens, scan the diff for these principles. Flag violations with catalog IDs below.

| Principle | What to look for in this repo | Prefer |
| --------- | ------------------------------ | ------ |
| **Clean Code** | Unclear names; functions/pages doing many jobs; comments that restate code or lie about behavior; dead imports, exports, files, or docs left after a refactor | Small units; intent-revealing names; comments only for non-obvious *why*; delete unused code |
| **KISS** | Speculative abstractions, premature helpers, config/env wired nowhere, duplicate seed/doc paths “for later”, over-generic types | Simplest correct diff; one clear path; extend when a second use case exists |
| **DRY** | Copy-pasted upsert/API/modal/seed logic; repeated magic strings; parallel docs saying the same thing differently | Extract shared helper or component; domain constants + typed FE mirror (`INTEGRATION_STATUS`, `domain.IntegrationStatus*`) |
| **DDD** | Business rules in handlers, middleware, or UI; domain importing framework/infra; feature-to-feature imports; wire/DB shapes leaking without boundary | Domain vocabulary in `internal/core/domain/`; services orchestrate; ports define contracts; `features/*` slices stay isolated |
| **SOLID** | **SRP:** handler/page owns orchestration + rules + I/O. **OCP:** editing stable registry/core to avoid adding a provider. **LSP:** mocks break interface contracts. **ISP:** bloated ports. **DIP:** core/service importing concrete SQL/HTTP instead of ports | Router → Handler → Service → Port ← Repository; thin handlers; extend via registry/providers, not forks |
| **Frontend imports** | Deep `../../` paths; cross-feature `@/features/<other>` | `@/components/ui/*`, `@/lib/*`, `@/features/<same-feature>/*` |
| **shadcn / Atomic** | Raw `<button>`/`<input>`/custom overlays in features; atoms duplicated outside `components/ui` | Compose `@/components/ui` (shadcn); feature components in `features/<name>/components/` only |
| **React 19** | `useState`+`useEffect` mirroring props/data; fetch without cleanup; conditional hooks; unnecessary global state | Derive during render; side effects in `useEffect` with deps; data in hooks |

**Agent conduct during smell**

- Run smell **yourself** after implementation; do not ask the user to run it.
- **Fix BLOCKER and HIGH** on the diff before marking work done (delivery) or approving merge (review).
- Cite **file:line** and a **catalog ID** for each finding; do not invent issues.
- Do not expand scope (“while we’re here”) unless the diff introduced the debt.
- If the diff is clean, say so — do not nitpick for the sake of a report.

### Backend (when Go changes)

Dependency flow — never skip a layer:

```
Router → Middleware → Handler → Service → Port ← Repository / Provider
```

| Area | Violation |
| ---- | --------- |
| `internal/core/domain/` | Imports Echo, HTTP, infra, presentation, or contains struct tags (`json`, `db`) |
| `internal/core/port/` | Contains implementation instead of interfaces |
| `internal/core/service/` | Direct SQL/HTTP; missing `ctx`; orchestration left in handler |
| `internal/infrastructure/` | Business rules; bypasses port interfaces |
| `internal/presentation/` | Business logic; handler → repository without service; RBAC in handler instead of middleware |
| `pkg/gateway`, `pkg/contracts` | Imports `internal/*`; embedded SDK requires PostgreSQL |
| Providers | No `init()` registration or missing blank import in `internal/app/imports.go` |
| Auth | Portal mutations without JWT + `RequirePermission`; internal ingest open outside local dev |
| Errors | Raw DB errors to clients; swallowed errors without justification |
| Concurrency | Blocking I/O without `ctx`; goroutines without cancellation |
| Secrets | Hardcoded credentials; JWT/passwords/provider keys in logs or responses |
| Migrations | `.up.sql` without `.down.sql`; empty or out-of-sequence files |
| Tests | Logic changed without tests; concurrent code without `-race` consideration |
| DDD / SOLID | Domain rules in handler; stringly-typed status/enums without `domain` constants; service bypasses port |

Go packages for verification:

```bash
golangci-lint run ./cmd/... ./internal/... ./pkg/...
go test -race ./cmd/... ./internal/... ./pkg/...
```

### Frontend (when `frontend/` changes)

**Stack:** React 19 — use modern patterns (no class components, no legacy lifecycle). Prefer hooks; avoid redundant state.

| Area | Violation |
| ---- | --------- |
| **Imports** | Deep relatives (`../../`, `../../../`) when `@/` alias exists — use `@/components/ui/*`, `@/lib/*`, `@/features/<same-feature>/*` |
| **Atomic / UI** | Raw `<button>`, `<input>`, hand-rolled overlays — compose **`@/components/ui`** (shadcn) primitives; feature components stay in `features/<name>/components/` |
| **shadcn** | Bypassing installed UI kit; `space-*` instead of `flex` + `gap-*`; missing `DialogTitle`; semantic tokens violated — see `docs/frontend/shadcn/` |
| **React 19 state** | Derived data stored in `useState` + synced via `useEffect`; fetch-on-render without cleanup; duplicate sources of truth; missing `useEffect` deps / stale closures |
| Routing | Routes or path helpers outside `src/core/router/` |
| Features | **Any** cross-feature import (`@/features/<other>`, relative `../<other>/`, barrels) — features are fully independent; compose in `core/router` only |
| Barrels | `index.ts` exports internals the router does not need |
| API | Token handling outside `core/api/client.ts`; portal calling unauthenticated `/internal/*` |
| Storybook | New or changed components without colocated `.stories.tsx` |
| Security | Tokens in console/logs; unsanitized `dangerouslySetInnerHTML` |
| Clean Code / KISS | Page > ~200 lines with extractable modals/forms; duplicated modal logic inlined after deleting a component |
| DDD / DRY | Magic strings for domain statuses; duplicated API result handling |

**Atomic design in this repo**

| Layer | Location | Rule |
| ----- | -------- | ---- |
| Atoms / primitives | `src/components/ui/` | shadcn components only — Button, Input, Dialog, Modal wrapper, etc. |
| Feature molecules/organisms | `features/<name>/components/` | Compose shadcn primitives; no duplicate atom implementations |
| Pages | `features/<name>/pages/` | Orchestrate feature components + hooks; thin JSX |

```bash
cd frontend && npm run lint && npm run test
```

### Dual-mode gateway

- **Embedded SDK** (`pkg/gateway`): config in code, no server, no DB — must not pull in server-only dependencies.
- **HTTP server**: config in PostgreSQL, portal UI for credentials — handlers stay thin; services own orchestration.

---

## Step 4 — Record findings

One finding per issue. Use a catalog ID, the smallest relevant excerpt, and one sentence each for **why** and **fix**.

Do not invent issues. If the diff is clean, say so.

### Project catalog (prefer these)

**Backend**

| ID | Meaning |
| -- | ------- |
| MGW.LAYER-VIOLATION | Skips or inverts a layer |
| MGW.DOMAIN-DEPS | Domain imports framework or infra |
| MGW.PKG-INTERNAL | Public `pkg/` imports `internal/` |
| MGW.HANDLER-BIZ | Business logic in handler or middleware |
| MGW.SERVICE-DIRECT-DB | Service talks to SQL without a port |
| MGW.RBAC-BYPASS | Workspace mutation without permission middleware |
| MGW.INTERNAL-INGEST-OPEN | Internal ingest exposed without local guard |
| MGW.PROVIDER-REG | Provider missing registry wiring |
| MGW.MIGRATION | Migration hygiene failure |
| MGW.LOG-SECRET | Secrets or PII in logs |
| MGW.N+1-QUERY | Query inside a loop |
| MGW.GOROUTINE-LEAK | Goroutine without lifecycle control |

**Frontend**

| ID | Meaning |
| -- | ------- |
| MGW.FE-CROSS-FEATURE | Feature imports another feature |
| MGW.FE-ROUTER-SCATTER | Routes outside `core/router` |
| MGW.FE-API-BYPASS | Portal uses wrong or unauthenticated API |
| MGW.FE-SPACE-Y | Uses `space-*` instead of `gap-*` |
| MGW.FE-NO-STORY | UI change without Storybook story |
| MGW.FE-PLACEHOLDER | Stub route/page/nav for unimplemented work |
| MGW.FE-TOKEN-LEAK | Token handling outside `core/api/client` |
| MGW.FE-IMPORT-ALIAS | Deep relative import (`../..`+) — use `@/` alias |
| MGW.FE-NO-SHADCN | Feature UI reimplements atoms — use `@/components/ui` shadcn primitives |
| MGW.FE-STATE-DERIVED | Redundant state + effect for values derivable from props/hook data |
| MGW.FE-STATE-EFFECT | Data fetch or side effect with missing cleanup, wrong deps, or race risk |

**General — Clean Code, KISS, DRY**

| ID | Meaning |
| -- | ------- |
| CC.G5 | Duplication — extract or consolidate (DRY) |
| CC.KISS | Over-engineered or speculative abstraction; simpler design suffices |
| CC.N1 | Name does not reveal intent |
| CC.F1 | Function/component does more than one thing (SRP) |
| CC.DEAD-CODE | Unused code, import, export, file, or stale doc left behind |
| CC.NOISE-COMMENT | Comment restates code, is outdated, or explains *what* instead of *why* |

**General — DDD**

| ID | Meaning |
| -- | ------- |
| DDD.LAYER | Business rule or orchestration in the wrong layer (handler/UI/infra) |
| DDD.BOUNDARY | Bounded context violated (feature↔feature, `pkg`→`internal`, shared→feature) |
| DDD.DOMAIN-LEAK | Domain object returned directly from handler without DTO mapping, or containing presentation tags (`json`) |
| DDD.VALUE-TYPE | Domain concept as raw string/number — use typed constant in domain + API/FE mirror |

**General — SOLID**

| ID | Meaning |
| -- | ------- |
| SOLID.SRP | Type or module has multiple reasons to change |
| SOLID.DIP | High-level module depends on low-level detail instead of port/interface |
| SOLID.OCP | Stable code modified where extension/registry would suffice |

**General — Go / TypeScript (use when no MGW ID fits)**

| ID | Meaning |
| -- | ------- |
| CC.G5 | Duplication — extract or consolidate |
| CC.N1 | Name does not reveal intent |
| CC.F1 | Function does more than one thing |
| GO.ERR-IGNORED | Error not handled |
| GO.ERR-NO-WRAP | Error chain broken (`%w` missing) |
| GO.CONTEXT-OMIT | Missing or misused `context.Context` |
| GO.GOROUTINE-LEAK | Goroutine leak risk |
| GO.SQL-CONCAT | SQL built by string concatenation |
| TS.UNHANDLED-PROMISE | Async error not handled |
| TS.HOOK-RULES | React hooks used incorrectly (Rules of Hooks, React 19) |
| TS.DERIVED-STATE | State mirrors props/other state — compute during render or useReducer |

---

## Step 5 — Report

Severity: **BLOCKER** · **HIGH** · **MEDIUM** · **LOW** · **NIT**

**BLOCKER** and **HIGH** must be fixed before merge.

````markdown
# Smell Report

**Base:** `<branch>`
**Classification:** `<category>`
**Surfaces:** `<backend | frontend | …>`
**Lens:** `<primary lens>`
**Principles checked:** Clean Code · KISS · DRY · DDD · SOLID

## Summary
- Files changed: N
- Findings: N (BLOCKER: n, HIGH: n, …)
- Top risk: …

## Findings

### [SEVERITY] `ID` — `path:line`
**Why:** …
**Fix:** …

## Synthesis
One paragraph: overall quality, merge readiness, principles adherence (KISS/DRY/DDD/SOLID), and any doc/test gaps.
````

If no issues: **"No findings on this diff — ready for `make audit`."**

---

## Gate

Part of **`docs/agents/verification.md`**. Agents run this themselves.

1. Fix compile/lint failures on touched paths first (commands above).
2. Complete this smell report; fix **BLOCKER** and **HIGH**.
3. Run **`make audit`** from repo root.
4. Repeat from step 2 if the diff changed materially.

Do not mark work done or open a PR until smell is clean (no BLOCKER/HIGH) and `make audit` passes.

Begin with Step 2.
