# Portal frontend — documentation index

Implementation lives in [`frontend/`](../../frontend/). Start with **[architecture.md](./architecture.md)** for folder layout and dependency rules.

## In this folder

| Document | Purpose |
| -------- | ------- |
| [architecture.md](./architecture.md) | Layers, routing, feature isolation |
| [conventions.md](./conventions.md) | Naming, KISS/DRY/SOLID, shadcn |
| [frontend-engineer.md](./frontend-engineer.md) | Workflow, aesthetics, self-review |
| [shadcn/SKILL.md](./shadcn/SKILL.md) | shadcn/ui contract and CLI |

## Stack (external)

[Vite](https://vite.dev/guide/) · [React](https://react.dev/) · [TypeScript](https://www.typescriptlang.org/docs/) · [Tailwind](https://tailwindcss.com/docs) · [shadcn/ui](https://ui.shadcn.com/docs)

## Backend

[Architecture](../backend/architecture.md) · [Usage](../backend/usage.md) · [Code conventions](../backend/code-conventions.md)

## Cursor / shadcn skill install

```bash
npx skills add https://github.com/shadcn/ui --skill shadcn
```

The in-repo copy under `shadcn/` remains the source of truth for this project.
