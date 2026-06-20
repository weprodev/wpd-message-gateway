# Agent instructions (`wpd-message-gateway`)

Canonical playbooks: **`docs/agents/`**. Entry points: **`CLAUDE.md`** (Claude Code), **`.cursorrules`**, **`.github/copilot-instructions.md`**.

## Required workflow

1. Read **`docs/agents/master-agent.md`** at session start.
2. **After any code change** (implement, fix, refactor): follow **`docs/agents/verification.md`** — lint → **`/smell develop`** → fix BLOCKER/HIGH → **`make audit`**. Agents run this themselves; never ask the user.
3. **Code review:** **`/smell develop`** first, then **`docs/agents/review-agent.md`**. After applying review fixes, repeat the verification chain.
4. **Before PR or marking done:** verification chain complete; zero BLOCKER/HIGH smell findings; diff respects KISS, DRY, DDD, and SOLID (see `/smell` principles pass).

## Agent commands

| Resource | Purpose |
| -------- | ------- |
| **`docs/agents/verification.md`** | Mandatory lint → smell → audit chain |
| `/smell [branch]` | Pre-PR diff review (`.claude/commands/smell.md`) |
| `make audit` | Full quality gate after smell |
| `docs/agents/master-agent.md` | Route implement vs review |
| `docs/agents/delivery-agent.md` | Feature implementation |
| `docs/agents/review-agent.md` | PR / diff review (after smell) |
| `docs/agents/prompts.md` | Copy-paste prompts |
| **golang-pro** (`.cursor/skills/golang-pro/`) | Go idioms + WPD layer rules for `cmd/`, `internal/`, `pkg/` |
| **typescript-react-reviewer** (`.cursor/skills/typescript-react-reviewer/`) | React 19 + TS review for `frontend/` |
| **typescript-advanced-types** (`.cursor/skills/typescript-advanced-types/`) | Advanced TS types for `frontend/` (complement to react reviewer) |
| **software-architecture** (`.cursor/skills/software-architecture/`) | Clean Architecture + DDD for cross-cutting design and refactors |

Architecture: **`docs/backend/architecture.md`**, **`docs/frontend/architecture.md`**.

**Skills:**
- `npx skills add https://github.com/jeffallan/claude-skills --skill golang-pro` → `.agents/skills/golang-pro/` · entry `.cursor/skills/golang-pro/`
- `npx skills add https://github.com/dotneet/claude-code-marketplace --skill typescript-react-reviewer` → `.agents/skills/typescript-react-reviewer/` · entry `.cursor/skills/typescript-react-reviewer/`
- `npx skills add https://github.com/wshobson/agents --skill typescript-advanced-types` → `.agents/skills/typescript-advanced-types/` · entry `.cursor/skills/typescript-advanced-types/`
- `npx skills add https://github.com/sickn33/antigravity-awesome-skills --skill software-architecture` → `.agents/skills/software-architecture/` · entry `.cursor/skills/software-architecture/`
