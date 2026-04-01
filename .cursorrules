# WPD Message Gateway: AI Agent Rules

This project enforces a rigorous engineering standard via specific `.md` definitions that govern coding, layering, UI styling, and orchestration.

## Primary Directive (The Master Agent)
Whenever taking action in this repository, you must initially adopt the mindset and protocols outlined in:
`docs/agents/master-agent.md`

### 1. Plan-and-Execute Protocol
You are forbidden from utilizing unbounded "ReAct" loops or infinite context-seeking that overwrites your token limits. You must:
1. Formulate a short execution plan.
2. Provide the plan to the user.
3. Execute the minimal required change in complete adherence to DDD Layering (Go) and Composition (React).
4. Verify via `make audit` or explicit test execution.

### 2. Specialist Delegation
You must delegate and refer to specific playbooks based on intent:
- For implementing new code, refactoring, or bug fixing: Follow the strict safety, security, and layer isolation rules mapped in `docs/agents/delivery-agent.md`
- For analyzing PRs, diffs, or code structure: Follow the vulnerability, optimization, and edge-case testing limits in `docs/agents/review-agent.md`

### 3. Absolute Constraints
- **Security**: Never expose or guess passwords, database connection strings, or JWTs. Explicitly load from ENV or configuration.
- **Frontend Constraints**: Use `shadcn/ui` composition limits. `gap-*` instead of `space-x`. Do not produce "generic SaaS AI slop". Favor stark, beautiful typography with zero cards by default.
- **Backend Constraints**: Domain layers are zero-dependency. Handlers must never touch infrastructure directly. Concurrency requires Context propagation. Unhandled errors must be tagged with `//nolint:errcheck` ONLY if explicitly justified.

### 4. Continuous Documentation Sync
You MUST ensure documentation is always kept perfectly up to date with any codebase modifications:
- **UI Components**: If you modify or add a frontend UI component, you must guarantee `Storybook` stories are updated.
- **Database Schema**: If you add/modify a database table or column, you must update the database architecture diagrams/markdown.
- **Architecture Layers**: If you create a new layer or directory, ensure it does not introduce an anti-pattern, explicitly verify it aligns with `docs/backend/architecture.md`, and update the documentation to reflect the new boundaries.

By executing tasks here, you agree to follow the "Antigravity Protocol" and Safety guidelines mapped inside `docs/agents/master-agent.md`.
