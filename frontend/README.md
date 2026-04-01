# Message Gateway — frontend

React 19 + TypeScript + Vite portal UI for the Message Gateway API. The dev server proxies `/api` to the Go backend (see `vite.config.ts`).

## Commands

From the **repository root**, use Make for day-to-day work and full checks:

| Make target | Purpose |
| ----------- | ------- |
| `make ui` | Start the Vite dev server (port **10104**) |
| `make storybook` | Start Storybook (port **6006**) |
| `make audit` | Format, lint, and test **Go and frontend**; then verify **builds** (Go compile, Vite app bundle, static Storybook) |
| `make build` | Compile Go packages and run frontend **`build:all`** (app + Storybook static build) |

`make audit` is the gate to run before pushing or opening a PR: it catches issues that lint alone does not (TypeScript errors, broken bundles, Storybook’s separate Vite graph). `make build` is the lighter “does everything compile?” check without reformatting or running tests.

### npm scripts (run inside `frontend/`)

| Command | Purpose |
| ------- | ------- |
| `npm run dev` | Vite dev server (port **10104**) |
| `npm run format` | ESLint with `--fix` (JS/TS/TSX) |
| `npm run lint` | ESLint (no fixes) |
| `npm run test` | Vitest (CI mode) |
| `npm run build` | `tsc` + production Vite bundle → `dist/` |
| `npm run build-storybook` | Static Storybook → `storybook-static/` |
| `npm run build:all` | `build` then `build-storybook` |
| `npm run storybook` | Storybook dev server (port **6006**) |
| `npm run preview` | Serve `dist/` locally |

## Source layout

| Area | Role |
| ---- | ---- |
| `src/app/` | Composition: router, aggregated `ROUTES`, shell wiring |
| `src/core/` | Cross-cutting infrastructure (HTTP client, React Query factory) |
| `src/features/*/` | **Bounded features** — each owns `paths`, UI pages, and feature-specific API/types (no imports from sibling features) |
| `src/components/ui/` | Shared primitives (button, input, icon, …) |
| `src/shared/` | Layout shells and other non-feature UI shared by the app |
| `src/lib/` | Small utilities (`cn`, etc.) |

Features expose **route path constants** from `paths.ts`. The app layer merges them in `src/app/paths.ts` for `NavLink` / `Navigate` and keeps the router as the single place that mounts feature pages.

## Design system

- **Storybook:** foundations (colors, icons) and component stories live under `src/stories/` and `src/components/**/*.stories.tsx`.
- **Theming:** semantic CSS variables in `src/index.css`; Material Symbols `@font-face` in `src/icons.css` (loaded from `main.tsx`).

## Conventions

- Prefer small, named modules over large pages; keep API calls in `*.api.ts` and types in `*.types.ts` inside the same feature.
- Use `@/` imports (`tsconfig` path alias to `src/`).

## Further reading (repository docs)

- [docs/frontend/README.md](../docs/frontend/README.md) — Portal doc index (Vite, TypeScript, links)
- [docs/frontend/conventions.md](../docs/frontend/conventions.md) — DDD-style features, KISS/DRY/SOLID, naming
- [docs/frontend/frontend-engineer.md](../docs/frontend/frontend-engineer.md) — role, principal workflow, security, Storybook, self-review
- [docs/frontend/shadcn/SKILL.md](../docs/frontend/shadcn/SKILL.md) — shadcn/ui skill (rules, CLI); official [shadcn docs](https://ui.shadcn.com/docs)
