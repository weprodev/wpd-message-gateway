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

Use `@/` for cross-layer imports; **relative paths within the same feature**.

## shadcn / styling

- Semantic tokens only (`bg-background`, `text-primary`) — see [shadcn rules](./shadcn/rules/).
- Layout spacing: `flex` + `gap-*`, not `space-x` / `space-y`.
- Buttons with icons: `data-icon="inline-start"` or `inline-end`.
- Dialogs require `DialogTitle` (or `sr-only`).

## Quality gate

From repo root: **[verification chain](../agents/verification.md)** (`/smell develop` → `make audit`). See [frontend/README.md](../../frontend/README.md) for scripts.
