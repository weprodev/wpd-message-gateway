# Portal frontend — documentation index

The **Message Gateway Portal** UI lives in [`frontend/`](../../frontend/) (Vite, React, TypeScript, Tailwind). This folder collects **repo-specific** guidance and links; implementation details and scripts are in [`frontend/README.md`](../../frontend/README.md).

## Stack

| Topic | Where to read |
| ----- | ------------- |
| **shadcn/ui** — docs & component catalogue | [ui.shadcn.com/docs](https://ui.shadcn.com/docs) |
| **shadcn/ui** — installation (Vite + React setup) | [ui.shadcn.com/docs/installation](https://ui.shadcn.com/docs/installation) |
| **Vite** | [vite.dev/guide](https://vite.dev/guide/) |
| **TypeScript** | [typescriptlang.org/docs](https://www.typescriptlang.org/docs/) |
| **React** | [react.dev](https://react.dev/) |
| **Tailwind CSS** | [tailwindcss.com/docs](https://tailwindcss.com/docs) |

## In this repository

| Document | Purpose |
| -------- | ------- |
| [conventions.md](./conventions.md) | DDD-style feature folders, KISS/DRY/SOLID, naming — aligns with this repo’s Portal |
| [frontend-engineer.md](./frontend-engineer.md) | **Role & workflow** — principal-style understand→plan→act, architecture, security, Storybook, shadcn alignment |
| [shadcn/SKILL.md](./shadcn/SKILL.md) | **shadcn/ui skill** (rules, CLI, patterns) — in-repo copy; upstream [skills.sh/shadcn/ui/shadcn](https://skills.sh/shadcn/ui/shadcn) |
| [shadcn/rules/](./shadcn/rules/) | Styling, forms, composition, icons, base vs radix |
| [`frontend/README.md`](../../frontend/README.md) | Scripts, `make` targets, `src/` layout |

## Backend / shared docs

- [Architecture](../backend/architecture.md) — system boundaries, gateway vs portal
- [Code conventions](../backend/code-conventions.md) — **Go** conventions (backend); Portal uses TypeScript rules in [conventions.md](./conventions.md)
- [Usage](../backend/usage.md) — HTTP API and SDK usage
- [Backend engineer role](../backend/backend-engineer.md) — Go layers, registry, `make audit`

## Cursor / AI

To install the **interactive** shadcn skill into Cursor (CLI-backed), use the upstream installer:

```bash
npx skills add https://github.com/shadcn/ui --skill shadcn
```

The copy under **`docs/frontend/shadcn/`** remains the **source of truth** for humans and for agents that read the repo without the skill runtime.
