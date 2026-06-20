# Portal frontend architecture

Canonical reference for `frontend/src/` structure. For naming, shadcn rules, and workflow see [conventions.md](./conventions.md) and [frontend-engineer.md](./frontend-engineer.md).

## Layer map

```
src/
├── core/          # Infrastructure — no feature UI
│   ├── api/
│   ├── router/    # ROUTES + router (only module that imports features)
│   └── services/
├── features/      # Independent bounded contexts
├── shared/        # Generic UI (no feature imports)
├── components/ui/
└── lib/
```

## Dependency rules (ESLint)

| Layer | May import | Must not import |
| ----- | ---------- | --------------- |
| `features/*` | `@/core/*`, `@/shared/*`, `@/components/*`, `@/lib/*`, `@/features/<same>/*` | **Any** `@/features/<other>/*` or relative import into another feature folder |
| `shared/*` | `@/core/*`, `@/components/*`, `@/lib/*` | `@/features/*` |
| `core/*` except `router/` | `@/lib/*` | `@/features/*` |
| `core/router` | feature barrels (`@/features/<name>`) | — |

Route composition lives in **`core/router/router.tsx`** only.

## Feature slice

```
features/workspaces/
├── index.ts                 # public barrel
├── pages/
├── layouts/                 # domain shells (not shared/)
├── hooks/use-*.hook.ts
├── *.api.ts
└── *.types.ts
```

Use **`@/features/<name>/…`** for imports inside a feature (not `../../`). Only **`core/router`** imports feature barrels.

## UI composition (Atomic + shadcn)

- **Atoms:** `components/ui/` — shadcn only; never duplicate Button/Input/Dialog in features.
- **Feature components:** `features/<name>/components/` — compose shadcn; domain-specific only.
- **Pages:** orchestrate hooks + feature components; React 19 — no redundant derived state.

## Routing

All paths: **`core/router/routes.ts`**.

```typescript
import { ROUTES } from "@/core/router/routes"

ROUTES.workspaces
ROUTES.workspace.overview(workspaceId)
```

## Auth

Session guards are deferred. Pages are reachable without login; JWT is attached by `core/api/client` when present. When auth ships, add a guard in `core/router/router.tsx` only.

## Adding a feature

1. Create `features/<name>/` with pages, hooks, `*.api.ts`, `*.types.ts`
2. Export public API from `index.ts`
3. Extend `core/router/routes.ts` and `core/router/router.tsx`
4. Colocate Storybook stories
