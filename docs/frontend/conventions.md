# Frontend conventions (Portal)

Rules for [`frontend/src/`](../../frontend/src/). **Structure & layers:** [architecture.md](./architecture.md). **Workflow:** [frontend-engineer.md](./frontend-engineer.md). **UI components:** [shadcn/SKILL.md](./shadcn/SKILL.md).

## Principles

- **KISS** — small modules; hooks for side effects; one router, one routes file.
- **DRY** — shared primitives in `components/ui` and `lib/`; routes only in `core/router/routes.ts`.
- **SOLID** — one concern per hook/page; domain layouts live in the owning feature.
- **Isolation** — no cross-feature imports (ESLint); compose in `core/router`.

## Naming

| Kind | Convention | Example |
| ---- | ---------- | ------- |
| Routes | `core/router/routes.ts` | `ROUTES.workspace.overview(wid)` |
| Feature barrel | `index.ts` | `features/inbox/index.ts` |
| Page | `*.page.tsx` | `overview.page.tsx` |
| Layout | `features/<name>/layouts/` | `workspace-layout.tsx` |
| Hook | `hooks/use-*.hook.ts` | `use-inbox-logs.hook.ts` |
| API / types | `*.api.ts`, `*.types.ts` | `inbox.api.ts` |

Use `@/` for all imports — one alias covers everything:

| Target | Import as |
| ------ | --------- |
| shadcn / atoms | `@/components/ui/<component>` |
| lib, shared, core | `@/lib/*`, `@/shared/*`, `@/core/*` |
| Same feature | `@/features/<feature>/…` — **not** `../../` |
| Router (barrels) | `@/features/<feature>` |

**Do not** import another feature — not via `@/features/<other>`, not via relative paths into another slice. ESLint enforces this in `eslint.config.js`. Route composition stays in `core/router` only.

## Atomic design + shadcn

| Layer | Path | Rule |
| ----- | ---- | ---- |
| Atoms | `components/ui/` | shadcn primitives — add via `npx shadcn@latest add` in `frontend/` |
| Feature UI | `features/<name>/components/` | Compose atoms; business-specific layout only |
| Pages | `features/<name>/pages/` | Orchestration; no raw HTML form controls |

Follow [shadcn/SKILL.md](./shadcn/SKILL.md): semantic tokens, `flex`+`gap-*`, `DialogTitle`, `cn()`.

## React 19

- Derive values during render — do not `useState` + `useEffect` to mirror props or hook data.
- Side effects in `useEffect` only; include cleanup for subscriptions/fetches; respect dependency arrays.
- Colocate data fetching in hooks (`use-*.hook.ts`); pages stay thin.
- React Compiler–friendly code: no conditional hooks; stable keys in lists.

## shadcn / styling

- Semantic tokens only (`bg-background`, `text-primary`) — see [shadcn rules](./shadcn/rules/).
- Layout spacing: `flex` + `gap-*`, not `space-x` / `space-y`.
- Buttons with icons: `data-icon="inline-start"` or `inline-end`.
- Dialogs require `DialogTitle` (or `sr-only`).

## Quality gate

From repo root: **[verification chain](../agents/verification.md)** (`/smell develop` → `make audit`). See [frontend/README.md](../../frontend/README.md) for scripts.

**Agent skills:** [typescript-react-reviewer](../../.cursor/skills/typescript-react-reviewer/SKILL.md) + [overlay](../../.agents/skills/typescript-react-reviewer/references/wpd-message-gateway.md); [typescript-advanced-types](../../.cursor/skills/typescript-advanced-types/SKILL.md) + [overlay](../../.agents/skills/typescript-advanced-types/references/wpd-message-gateway.md); [software-architecture](../../.cursor/skills/software-architecture/SKILL.md) + [overlay](../../.agents/skills/software-architecture/references/wpd-message-gateway.md).
