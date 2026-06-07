# WPD Message Gateway — Portal frontend overlay

Project-specific rules for **typescript-react-reviewer** in this repository. These **override** generic guidance where they conflict.

**Canonical docs:** `docs/frontend/architecture.md`, `docs/frontend/conventions.md`, `docs/frontend/frontend-engineer.md`, `docs/agents/verification.md`

---

## Stack (not Next.js)

| Layer | Choice |
| ----- | ------ |
| Runtime | React **19** + Vite SPA |
| Routing | `react-router-dom` v7 — routes in `core/router/routes.ts` only |
| Server state | **TanStack Query** (`@tanstack/react-query`) where async lists matter |
| Styling | Tailwind + shadcn/ui — semantic tokens only |
| Tests | Vitest + Testing Library (colocated `*.test.tsx`) |
| Docs | Storybook (colocated `*.stories.tsx`) |

**No** React Server Components, **no** `'use server'` — upstream Next.js guidance does not apply.

---

## Feature-sliced architecture (ESLint-enforced)

```
src/
├── core/          # api, router (only router imports features)
├── features/      # auth, workspaces, inbox — isolated bounded contexts
├── shared/        # generic UI shell (no feature imports)
├── components/ui/ # atomic design system
└── lib/
```

| Layer | May import | Must not |
| ----- | ---------- | -------- |
| `features/*` | `@/core/*`, `@/shared/*`, `@/components/*`, relative same-feature | `@/features/*` |
| `shared/*`, `core/*` (except router) | `@/core/*`, `@/components/*`, `@/lib/*` | `@/features/*` |
| `core/router` | feature barrels | — |

**Barrel files (`features/<name>/index.ts`) are required** for public feature APIs — this is intentional, not bundle-bloat anti-pattern. Do not flag feature `index.ts` barrels; flag cross-feature imports instead (`MGW.FE-CROSS-FEATURE`).

---

## WPD-specific review priorities

### Block merge (maps to `/smell` BLOCKER/HIGH)

| ID | Issue |
| -- | ----- |
| `MGW.FE-CROSS-FEATURE` | Feature imports another feature |
| `MGW.FE-ROUTER-SCATTER` | Routes outside `core/router/` |
| `MGW.FE-API-BYPASS` | Token handling outside `core/api/client.ts` |
| `MGW.FE-TOKEN-LEAK` | JWT in console/logs |
| `MGW.FE-SPACE-Y` | `space-x-*` / `space-y-*` instead of `flex` + `gap-*` |
| `MGW.FE-NO-STORY` | New/changed UI component without `.stories.tsx` |
| `MGW.FE-PLACEHOLDER` | Stub route/nav for unimplemented work |

### TypeScript conventions

- Explicit props: `function Foo({ id }: FooProps)` — avoid `React.FC`
- API results: discriminated unions `{ ok: true; items: T[] } | { ok: false; message?: string }` in `*.api.ts`
- Types live in `*.types.ts` per feature — no duplicate backend shapes
- Prefer `ROUTES` helpers from `core/router/routes.ts` over string literals

### React patterns in this repo

- **Colocate state** in feature pages/hooks (`use-*.hook.ts`)
- **No derived state in `useEffect`** — compute during render
- **Fetch in `*.api.ts`** via `apiFetch` from `core/api/client.ts` — never raw `fetch` with manual token headers in features
- **Hooks naming:** `use-workspaces.hook.ts`, `use-inbox-logs.hook.ts`
- **Components:** atomic primitives in `components/ui/<name>/` with `index.ts`, test, and story

### shadcn / a11y

- `DialogTitle` required (or `sr-only`)
- `aria-invalid` on form fields with errors
- Icon buttons: `data-icon` attribute per conventions
- Semantic color tokens — no hardcoded hex in features

---

## Verification chain (mandatory after frontend changes)

```bash
cd frontend && npm run lint && npm run test
/smell develop          # fix BLOCKER/HIGH
make audit              # from repo root
```

Project gate: **zero BLOCKER/HIGH smell findings** and **`make audit` exit 0**.

---

## Test expectations

| Surface | Requirement |
| ------- | ----------- |
| `components/ui/*` | colocated `*.test.tsx` + `*.stories.tsx` |
| Feature components | colocated tests + stories |
| Feature pages/hooks | add tests when adding non-trivial logic |
| API modules | type-safe return unions; mock `apiFetch` in tests |

---

## When implementing or reviewing frontend here

1. Read `docs/frontend/architecture.md` for layer boundaries
2. Apply upstream **typescript-react-reviewer** critical checks (hooks, mutations, `any`)
3. Apply **typescript-advanced-types** + [types overlay](../../typescript-advanced-types/references/wpd-message-gateway.md) for API unions and DTO layout
4. Apply this overlay for WPD ESLint rules, smell IDs, and Storybook sync
5. Run verification chain before marking done
