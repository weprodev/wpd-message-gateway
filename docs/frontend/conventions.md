# Frontend conventions (Portal)

These rules apply to code under [`frontend/src/`](../../frontend/src/). They complement **[shadcn/ui](./shadcn/SKILL.md)**, **[frontend engineer — role & workflow](./frontend-engineer.md)**, and [backend code conventions](../backend/code-conventions.md) (Go).

## 1. Principles & React Patterns

- **Composition Over Inheritance** — Build interfaces via composite components `<Card><CardHeader>Title</CardHeader></Card>` rather than giant prop-driven monolithic components like `<Card title="Title" />`.
- **KISS** — Prefer small modules, clear flow, minimal abstraction until duplication hurts. Use custom hooks (`useDebounce`, `useToggle`) to keep components clean.
- **DRY** — Shared UI primitives in `components/ui`; shared utilities in `lib/`; avoid duplicating API or types across features.
- **DDD (bounded contexts)** — Each **`features/<name>/`** slice owns its **paths**, **pages**, **API calls** (`*.api.ts`), and **types** (`*.types.ts`). Compose via **`app/`** (router, paths) and **`shared/`**.
- **SOLID** — Single-purpose components and hooks. Rely on `ErrorBoundary` for fallbacks rather than polluting components with constant `if (err) return <Error />`.
- **Performance Budget** — Default to `useMemo` for heavy sorting/filtering, and virtualization (`@tanstack/react-virtual`) if lists exceed ~50 rows. 

## 2. Naming & Structure

| Kind | Convention | Example |
| ---- | ---------- | ------- |
| Routes (per feature) | `paths.ts` exporting path constants | `features/auth/paths.ts` |
| Pages | `*.page.tsx` | `login.page.tsx` |
| API layer | `*.api.ts` | `auth.api.ts` |
| Types | `*.types.ts` | `workspace.types.ts` |
| UI primitives | kebab-case files in `components/ui` | `button.tsx` |
| React components | PascalCase | `WorkspacesPage` |

Use the **`@/`** alias for imports from `src/`.

## 3. Strict shadcn/ui & Composition Rules

We adopt **shadcn** as the **copy-paste source** into `src/components/ui`. See [`docs/frontend/shadcn/SKILL.md`](./shadcn/SKILL.md) for full context.
Canonical references: [ui.shadcn.com/docs](https://ui.shadcn.com/docs) · [Installation](https://ui.shadcn.com/docs/installation)

### Styling Rules
- **Semantic tokens exclusively:** Use `bg-background` and `text-muted-foreground` internally (never `bg-gray-900`).
- **Use `gap-*`**: *Never* use `space-y-*` or `space-x-*`. Always use `flex flex-col gap-*`.
- **Equal Dimensions**: Use `size-10` rather than `w-10 h-10`.
- **Truncation**: Use `truncate`, never `overflow-hidden text-ellipsis whitespace-nowrap`.
- **Icons**: Icons inside Buttons MUST use `data-icon="inline-start"` instead of ad-hoc padding/scaling classes. 

### Forms & Structure
- **Forms Layout**: Forms heavily use `FieldGroup` + `Field`. **Never** use raw `div` with `space-y-*` for form constraints.
- **Validation Linking**: Always rely on `data-invalid` on the `Field` container and `aria-invalid` on the control element to ensure accessibility logic triggers correctly. 
- **Empty States**: Use `Empty` or `Alert` generic shadcn wrappers. Do not build custom empty-state spans.
- **Dialogs & Overlays**: Every `Dialog`, `Sheet`, and `Drawer` MUST have a title constraint (`DialogTitle`, `SheetTitle`) or `className="sr-only"` for A11y.

## 4. Aesthetic & Execution Directives

- **Reject Genericism**: Ban components and typography that replicate generic dashboard UI. Default toward **bold structural constraints** — limit colors to a rigid core palette per workspace.
- **Micro-interactions via CSS**: Rely on pure CSS for most hover states. When orchestrating massive reveals or sticky transitions, use `framer-motion` strategically (limit heavy JS-driven animations).
- **No Manual Dark Mode Strings**: Do not inject manual `dark:text-white`. Modify the CSS variables inside `src/index.css` directly so tokens absorb the theming gracefully.

## 5. Quality Gate

From repo root: **`make audit`** runs format, ESLint, test (Vitest), and **build verification** for Go and the frontend (including Storybook static build). See [`frontend/README.md`](../../frontend/README.md).

Before pushing a branch or handing off your implementation, a 0-exit code from `make audit` guarantees absolute compliance with all TypeScript strictness, linting, and structural barriers.
