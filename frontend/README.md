# Message Gateway — frontend

React 19 + TypeScript + Vite portal. Dev server proxies `/api` to the Go backend (`vite.config.ts`).

## Commands

Run from **repository root** when possible:

| Make target | Purpose |
| ----------- | ------- |
| `make ui` | Vite dev server (port **10104**) |
| `make storybook` | Storybook (port **6006**) |
| `make audit` | Format, lint, test, build (Go + frontend + Storybook) |
| `make build` | Compile Go + `npm run build:all` |

Inside `frontend/`:

| npm script | Purpose |
| ---------- | ------- |
| `npm run dev` | Vite dev server |
| `npm run lint` / `format` | ESLint |
| `npm run test` | Vitest |
| `npm run build` | Production bundle → `dist/` |
| `npm run build:all` | App + Storybook static build |

## Source layout

| Path | Role |
| ---- | ---- |
| `src/core/router/` | `ROUTES` + route tree (composition) |
| `src/features/*/` | Independent feature slices (`index.ts` barrels) |
| `src/shared/` | Generic UI — no feature imports |
| `src/components/ui/` | Independent shared UI components |

### Component Development Rules
To keep components isolated and testable, every React component must live in its own dedicated directory with the following structure:
```
[ComponentName]/
├── [ComponentName].tsx         # Implementation
├── index.ts                     # Clean re-export (e.g., export * from "./[ComponentName]")
├── [ComponentName].stories.tsx # Storybook stories
└── [ComponentName].test.tsx    # Vitest + React Testing Library unit test
```

## Storybook in Docker
Storybook is integrated into our containerized local development stack:
- Spin up Storybook in Docker: `docker compose up -d --build storybook`
- Access Storybook in your browser: `http://localhost:6006`

## Docs

- [Architecture](../docs/frontend/architecture.md)
- [Conventions](../docs/frontend/conventions.md)
- [Frontend engineer role](../docs/frontend/frontend-engineer.md)
