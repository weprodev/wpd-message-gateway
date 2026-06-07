---
name: frontend-engineer
description: Execution model for the Message Gateway Portal (Vite/React/TypeScript). Use for frontend/, Storybook, shadcn/ui, and docs/frontend — guides humans and AI agents (phases, boundaries, verification, output).
---

# Frontend engineer — role & operating model

This document is both a **team playbook** and an **operational contract** for anyone (including agents) changing the Portal under [`frontend/`](../../frontend/). It demands the mindset and restraint of a **Principal UI Engineer** who rejects generic aesthetics in favor of deliberate, production-grade output.

It complements **[conventions.md](./conventions.md)**, **[frontend/README.md](../../frontend/README.md)**, and **[shadcn/SKILL.md](./shadcn/SKILL.md)**.

**Agent skills (load both for frontend work):**

- [typescript-react-reviewer](../../.cursor/skills/typescript-react-reviewer/SKILL.md) — React 19 + component review + [overlay](../../.agents/skills/typescript-react-reviewer/references/wpd-message-gateway.md) — [skills.sh/dotneet/claude-code-marketplace/typescript-react-reviewer](https://www.skills.sh/dotneet/claude-code-marketplace/typescript-react-reviewer)
- [typescript-advanced-types](../../.cursor/skills/typescript-advanced-types/SKILL.md) — API unions, DTOs, generics + [overlay](../../.agents/skills/typescript-advanced-types/references/wpd-message-gateway.md) — [skills.sh/wshobson/agents/typescript-advanced-types](https://www.skills.sh/wshobson/agents/typescript-advanced-types)
- [software-architecture](../../.cursor/skills/software-architecture/SKILL.md) — feature-sliced boundaries + [overlay](../../.agents/skills/software-architecture/references/wpd-message-gateway.md) — [skills.sh/sickn33/antigravity-awesome-skills/software-architecture](https://www.skills.sh/sickn33/antigravity-awesome-skills/software-architecture)

---

## 1. Execution model (use for substantive work)

For non-trivial tasks (new feature slice, API wiring, component API changes, routing), follow this sequence:

### 1.1 Aesthetic & Interaction Thesis
Before writing any code, establish a BOLD point-of-view:
- **Visual thesis**: One sentence describing mood and tone. Reject generic "AI slop."
- **Content plan**: What is the singular job of this view? Keep copy extremely sparse. 
- **Interaction thesis**: 1-2 intentional UI motions (via Framer Motion) that change the feel of the page while maintaining clarity.

### 1.2 Understand & Map
- Read **`app/`**, affected **`features/`**, **`core/`**, and **`components/`** before editing.
- Map **user intent**, **data flow** (API → state → UI), **routing**, and **feature boundaries**.

### 1.3 Plan
- List **files to touch** and **why** (minimal set).
- Outline **Performance & Composition Hooks**: Memoization (`useMemo`/`React.memo`), component lifecycle, and code isolation boundaries.
- Note **Storybook**, **tests**, and **security** (tokens, user input).

### 1.4 Implement
- Execute with restraint. Default toward typography, spacing, and bold contrast.
- Follow [conventions.md](./conventions.md) and [architecture.md](./architecture.md) (barrels, `*.api.ts`, shadcn).

### 1.5 Verify (self-check)

Run [verification.md](../agents/verification.md):

- [ ] Fast lint on touched frontend paths
- [ ] **`/smell develop`** — no BLOCKER/HIGH on the diff
- [ ] **`make audit`** passes from repo root
- [ ] No **cross-feature imports** (maintain DDD isolation)
- [ ] Accessible markup: focus traps, `aria-*`, keyboard navigation
- [ ] **Storybook** updated for visual assertions

---

## 2. Architecture — layers

See **[architecture.md](./architecture.md)** for the layer diagram, ESLint boundaries, and how to add features.

| Layer | Responsibility |
| ----- | ---------------- |
| **`src/core/`** | HTTP client, React Query, **`router/`** (`ROUTES` + route tree) |
| **`src/features/<name>/`** | Pages, layouts, hooks, `*.api.ts`, `*.types.ts`, `index.ts` barrel |
| **`src/shared/`** | Generic layouts/UI — **no feature imports** |
| **`src/components/ui/`** | shadcn primitives |
| **`src/app/`** | Bootstrap helpers (e.g. home redirect) |
| **`src/lib/`** | Utilities |

Compose features in **`core/router/router.tsx`** only — not by importing feature A from feature B.

---

## 3. Aesthetic Constraints & Art Direction

Ship interfaces that feel deliberate, premium, and current (e.g., Linear-style restraint):
- **Start with composition, not components**. Prefer full-bleed layouts, generous whitespace, and sharp alignment.
- **No cards by default**. If a layout works as plain text hierarchy, remove the card border. Avoid mosaic "dashboard-card" soup.
- **Reject overused AI aesthetics**. Avoid clichéd purple gradients on white backgrounds and uninspired floating components. 
- **Sparse, utility-driven typography**. Treat the first view as a poster. If 30% of the copy can be deleted without losing meaning, delete it.
- **Motion must be intentional**. Exclude decorative motion noise. Use it for section reveals, shared layout transitions, and context persistence.

---

## 4. Advanced Frontend Patterns

- **Composition Over Inheritance**: Avoid "God components" with 20 configure props. Use Compound Components (e.g., `<Tabs><TabList><Tab/></TabList></Tabs>`).
- **Code Splitting & Lazy Loading**: Use `lazy()` and `<Suspense>` to separate heavy charts, 3D backgrounds, or complex below-the-fold interactions.
- **Virtualization**: Use tools like `@tanstack/react-virtual` for data tables or infinite feeds that exceed 50+ rows to maintain 60FPS scrolling.
- **Decoupled Data Flow**: Create pure presentation components fed by custom data-fetching hooks (e.g., `useMarketsQuery`) or generalized state wrappers (e.g., `useToggle`, `useDebounce`).

---

## 5. shadcn/ui Strict Compliance

- **[shadcn/SKILL.md](./shadcn/SKILL.md)** and **`shadcn/rules/`** are the absolute contract for styling, forms, and composition.
- **Do not guess APIs:** Use `npx shadcn@latest info`. 
- Incorporate upstream components cleanly via the CLI using `--dry-run` and `--diff` checks for local overwrites.
- **Canonical docs:** [ui.shadcn.com/docs](https://ui.shadcn.com/docs) — [Installation Guide](https://ui.shadcn.com/docs/installation)
- **Semantic tokens exclusively**: Use `bg-background` and `text-muted-foreground` in `cn()` strings over raw `bg-gray-900`.

---

## 6. Security & privacy

- **Auth:** Never log or commit tokens. Rely entirely on **`core/api/client`** logic to hydrate `Authorization` headers.
- **XSS:** Rely on React’s default escaping; avoid **`dangerouslySetInnerHTML`**.
- **Env:** Do not embed secrets in client bundles.

---

## 7. Storybook

- Stories: **`src/components/**/*.stories.tsx`**, foundations: **`src/stories/`**.
- **Update stories in the same PR** when component structural changes occur.
- Run `make audit` to implicitly ensure `npm run build-storybook` validates successfully.

---

## 8. Anti-patterns (Common Mistakes)

- **AI Aesthetic Slop**: Defaulting to Inter/Roboto, arbitrary shadow blobs, and repeating generic "SaaS card grid" impressions.
- **Cross-feature imports** to “save a line” — directly couples execution.
- **Space-y/Space-x hacks**: Using `space-x-4` instead of `flex gap-4`. Avoid at all costs.
- **God Props**: `export function Card({ title, subtitle, image, footer, actions ... })` instead of pure composition.
- **Prop-drilled callbacks** everywhere instead of Context/Reducer or centralized Zustand stores.

---

## 9. Decision heuristics (when unsure)

1. **Simpler** over clever.
2. **Explicit composition** over implicit prop configuration.
3. **Local (feature-scoped) states** over global wrappers.
4. **Temporary duplication** over **wrong abstractions**.

---

## 10. Guiding principle

> Prefer frameworks and visual rules that stay **bold**, **intentional**, and **easy to alter**. Precision always outperforms excessive decoration.

---

## References

| Doc | Link |
| --- | ---- |
| Architecture (layers, routing) | [architecture.md](./architecture.md) |
| Portal README (scripts) | [frontend/README.md](../../frontend/README.md) |
| Conventions | [conventions.md](./conventions.md) |
| Doc index | [README.md](./README.md) |
| shadcn skill | [shadcn/SKILL.md](./shadcn/SKILL.md) |
