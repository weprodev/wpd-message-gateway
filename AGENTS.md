# Agent instructions (`wpd-message-gateway`)

Canonical playbooks: **`docs/agents/`**. Entry points: **`.cursorrules`**, **`.github/copilot-instructions.md`**.

## Required workflow

1. Read **`docs/agents/master-agent.md`** at session start.
2. **After any code change** (implement, fix, refactor): follow **`docs/agents/verification.md`** — lint → **`/smell develop`** → fix BLOCKER/HIGH → **`make audit`**. Agents run this themselves; never ask the user.
3. **Code review:** **`/smell develop`** first, then **`docs/agents/review-agent.md`**. After applying review fixes, repeat the verification chain.
4. **Before PR or marking done:** verification chain complete; zero BLOCKER/HIGH smell findings.

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

Architecture: **`docs/backend/architecture.md`**, **`docs/frontend/architecture.md`**.
