# WPD Message Gateway — architecture overlay

Project-specific rules for **software-architecture** in this repository. Complements domain skills — load **with** golang-pro (backend) or typescript-react-reviewer + typescript-advanced-types (frontend).

**Canonical docs:** `docs/backend/architecture.md`, `docs/frontend/architecture.md`, `docs/agents/verification.md`

---

## System shape (two modes, one codebase)

| Mode | Surface | Config source | Architecture goal |
| ---- | ------- | ------------- | ----------------- |
| Embedded SDK | `pkg/gateway`, `pkg/contracts`, `pkg/provider/*` | In-process `gateway.Config` | Library-first, zero infra |
| HTTP server | `cmd/server`, `internal/*`, `frontend/` | PostgreSQL + Portal UI | Hexagonal + feature-sliced UI |

**Iron rule:** `pkg/*` never imports `internal/*`. SDK consumers depend only on public packages.

---

## Backend — Clean Architecture / DDD

```
Router → Middleware → Handler → Service → Port ← Repository / Provider
```

| Bounded context | Location | Responsibility |
| --------------- | -------- | -------------- |
| Domain | `internal/core/domain/` | Entities, value objects — stdlib only |
| Ports | `internal/core/port/` | Repository + inbox interfaces |
| Use cases | `internal/core/service/` | Orchestration, no SQL/HTTP |
| Adapters (out) | `internal/infrastructure/` | Postgres, inbox store, providers |
| Adapters (in) | `internal/presentation/` | Echo handlers, DTO decode |
| Public contracts | `pkg/contracts/` | Sender interfaces, message DTOs |
| Provider plugins | `pkg/provider/<name>/` | `register.go` → `pkg/registry` |

### Separation of concerns (block merge)

| Smell ID | Violation |
| -------- | --------- |
| `MGW.ARCH-PKG-INTERNAL` | `pkg/` imports `internal/` |
| `MGW.ARCH-HANDLER-SQL` | SQL or provider SDK in handlers |
| `MGW.ARCH-SERVICE-HTTP` | Echo/HTTP types in services |
| `MGW.ARCH-PROVIDER-INTERNAL` | New provider under `internal/` instead of `pkg/provider/<name>/` |
| `MGW.ARCH-FAT-HANDLER` | Handler > ~200 lines or business logic in handler |

**Prefer:** split handlers by portal concern (`portal_auth_handler`, `portal_workspace_handler`, …). Keep handlers thin — decode, call service, map response.

### Library-first (this repo)

- **Use:** Echo, `slog`, `pgx`, registry pattern, existing `pkg/contracts` senders
- **Don't:** Custom retry/rate-limit when a small service method + existing patterns suffice
- **Don't:** Duplicate sender contracts in `internal/core/port/` — they live in `pkg/contracts/`

---

## Frontend — feature-sliced Clean Architecture

```
core/ (infra) → features/ (bounded contexts) → shared/ + components/ui/
```

| Layer | Role | Boundary |
| ----- | ---- | -------- |
| `core/api`, `core/router` | Infrastructure | Token handling only in `client.ts` |
| `features/*` | Domain UI + use cases | No cross-feature imports |
| `shared/`, `components/ui/` | Presentation primitives | No feature imports |

### Separation of concerns (block merge)

| Smell ID | Violation |
| -------- | --------- |
| `MGW.ARCH-FE-CROSS-FEATURE` | Feature imports another feature |
| `MGW.ARCH-FE-LOGIC-IN-UI` | Fetch/business rules in presentational components |
| `MGW.ARCH-FE-ROUTER-SCATTER` | Route paths outside `core/router/routes.ts` |
| `MGW.ARCH-FE-FAT-PAGE` | Page > ~200 lines — extract hooks/subcomponents |

**Prefer:** `*.api.ts` for I/O, `use-*.hook.ts` for state, pages compose layout + hooks.

---

## Shared principles (upstream + WPD)

- **Early return** over deep nesting (max ~3 levels)
- **Single responsibility** per file/module — split at ~200 lines
- **Domain naming** over generic dumps (`OrderCalculator` > `utils/helpers.go` with unrelated funcs)
- **KISS** — smallest layer-legal change; no speculative abstractions
- **DRY** — reuse `pkg/contracts`, feature barrels, shared UI — don't copy DTOs
- **Plan before large edits** — state layer impact in 3–7 bullets (master-agent protocol)

### `shared/` and `lib/` in this repo

These names are **allowed** with **narrow scope**:

- `frontend/src/shared/` — generic shell/layout only, no feature imports
- `frontend/src/lib/` — small cross-cutting helpers (e.g. `cn()`)
- Go `pkg/` — public SDK surface, not a junk drawer

Do not add catch-all `utils` packages with unrelated functions.

---

## When to load this skill

| Task | Load with |
| ---- | --------- |
| New feature, refactor, layer move | Domain skill + **software-architecture** |
| PR / design review | **review-agent** + **software-architecture** + domain skill |
| Provider or SDK design | **golang-pro** + **software-architecture** |
| New feature slice or portal page | **typescript-react-reviewer** + **software-architecture** |

---

## Verification

After architectural changes:

```bash
/smell develop    # fix BLOCKER/HIGH (layer violations)
make audit        # full gate
```

Update `docs/backend/architecture.md` or `docs/frontend/architecture.md` when boundaries or package layout change.
