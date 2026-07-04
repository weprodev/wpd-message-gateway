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

Portal routes under `/workspaces/:wid/*` require a JWT (`ProtectedRoute` in `core/router/protected-route.tsx`). The API client in `core/api/client.ts` attaches the token to requests.

### Workspace authorization (UI)

Inside `WorkspaceLayout`, `WorkspaceAuthorizationProvider` exposes the active workspace `role` and resolved `permissions[]` from `GET /api/v1/workspaces`. Use this for **UI gating only** — the backend enforces permissions via wpd-gogate middleware on every route.

```tsx
import { Can, Permission, Role, useWorkspaceAuthorization } from "@/core/auth"

// Declarative
<Can permission={Permission.APIKeysWrite}>
  <Button>Generate key</Button>
</Can>

// Imperative
const { can, hasRole } = useWorkspaceAuthorization()
if (can(Permission.SettingsWrite)) { /* ... */ }
```

Constants mirror `internal/core/domain/permission.go`. Static role matrices (e.g. invitation role picker) can use `hasRolePermission()` from `core/auth/permissions.ts`.

**Never rely on UI checks for security** — they only hide controls; unauthorized API calls still return 403.

## Adding a feature

1. Create `features/<name>/` with pages, hooks, `*.api.ts`, `*.types.ts`
2. Export public API from `index.ts`
3. Extend `core/router/routes.ts` and `core/router/router.tsx`
4. Colocate Storybook stories
