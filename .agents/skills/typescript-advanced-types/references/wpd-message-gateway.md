# WPD Message Gateway — Portal TypeScript overlay

Project-specific type rules for **typescript-advanced-types** in `frontend/`. Complements **typescript-react-reviewer** — use both for frontend work.

**Canonical docs:** `docs/frontend/architecture.md`, `docs/frontend/conventions.md`, `docs/frontend/frontend-engineer.md`

---

## Compiler baseline

`frontend/tsconfig.app.json` enforces:

- `strict: true`
- `noUnusedLocals`, `noUnusedParameters`
- `verbatimModuleSyntax`, `erasableSyntaxOnly`
- `noFallthroughCasesInSwitch`

Do not weaken these. Prefer fixing types over `@ts-expect-error`.

---

## File layout (feature-sliced)

| File pattern | Purpose |
| ------------ | ------- |
| `features/<name>/<name>.types.ts` | Domain DTOs, unions, enums — **no fetch logic** |
| `features/<name>/<name>.api.ts` | Typed API calls returning discriminated unions |
| `features/<name>/use-*.hook.ts` | Narrow `ok` unions; expose hook-friendly state |
| `core/router/routes.ts` | `ROUTES` const object — single source for path literals |

**Rule:** Backend JSON shapes live once in `*.types.ts`. API modules map responses; hooks/pages consume narrowed types.

---

## API result pattern (required)

All `*.api.ts` functions return **discriminated unions** on `ok`:

```typescript
type ListResult =
  | { ok: true; items: Workspace[] }
  | { ok: false; status: number; message?: string }
```

**Narrowing in consumers:**

```typescript
const result = await fetchWorkspaces()
if (!result.ok) {
  setError(result.message ?? "Failed")
  return
}
setWorkspaces(result.items) // items is Workspace[]
```

| Smell ID | Issue |
| -------- | ----- |
| `MGW.FE-API-THROW` | Throwing from `*.api.ts` instead of `{ ok: false }` |
| `MGW.FE-API-ANY` | `any` on response bodies — use `unknown` + narrow |
| `MGW.FE-UNION-SKIP` | Accessing success fields without `ok` check |

---

## Preferred type patterns in this repo

### String literal unions (channels, statuses)

```typescript
export type MessageChannel = "email" | "sms" | "push" | "chat"
```

Use unions over `enum` unless interop requires it. Exhaustive `switch` with `noFallthroughCasesInSwitch`.

### Props interfaces

```typescript
type InboxTableProps = {
  rows: LogRow[]
  onDelete: (id: string) => void
}

function InboxTable({ rows, onDelete }: InboxTableProps) { ... }
```

Avoid `React.FC`. Prefer explicit props + `children?: React.ReactNode` when needed.

### Unknown over any

For opaque backend fields (e.g. `attachments?: unknown[]`), keep `unknown` until a parser/validator exists. Do not widen to `any`.

### Generics — when to use

| Use generics | Skip generics |
| ------------ | ------------- |
| Reusable `components/ui/*` (DataTable columns, Select options) | One-off feature pages |
| Shared helpers in `lib/` | API modules (concrete DTOs) |
| Type-safe event maps in shared hooks | Simple CRUD with fixed shapes |

**KISS:** Portal features are bounded contexts — prefer concrete types in `*.types.ts` over clever conditional/mapped types unless reuse is proven.

---

## Integration with React reviewer skill

| Concern | Skill |
| ------- | ----- |
| Hooks, effects, components, Storybook | **typescript-react-reviewer** |
| API unions, generics, narrowing, DTO layout | **typescript-advanced-types** (this overlay) |

Load **both** for frontend implementation and review.

---

## Verification

```bash
cd frontend && npm run lint && npm run test
/smell develop
make audit
```

Type-only changes must still pass `tsc` (via `npm run lint` / build) and existing tests.
