# Agent prompts

Copy-paste snippets for common workflows. Full chain: **[verification.md](./verification.md)**.

## Finish implementation (mandatory)

Agents must run both — do not ask the user:

```
/smell develop
make audit
```

Follow `.claude/commands/smell.md` (architecture + Clean Code, KISS, DRY, DDD, SOLID) and fix BLOCKER/HIGH before `make audit`.

## Lint after Go changes (fast feedback)

```bash
golangci-lint run ./cmd/... ./internal/... ./pkg/...
go test -race ./cmd/... ./internal/... ./pkg/...
```

## Lint after frontend changes

```bash
cd frontend && npm run lint && npm run test
```

## Code review (PR or diff)

```
/smell develop
Follow docs/agents/review-agent.md
```

After applying review fixes, repeat `/smell develop` and `make audit`.

## Implement a feature

```
Follow docs/agents/delivery-agent.md
End with: /smell develop → fix BLOCKER/HIGH → make audit (repo root)
```

## Backend Go (golang-pro + software-architecture)

```
Read .cursor/skills/software-architecture/SKILL.md
Read .agents/skills/software-architecture/references/wpd-message-gateway.md
Read .cursor/skills/golang-pro/SKILL.md
Read .agents/skills/golang-pro/references/wpd-message-gateway.md
Follow docs/agents/delivery-agent.md
go test -race ./cmd/... ./internal/... ./pkg/...
/smell develop → make audit
```

## Frontend Portal (react reviewer + advanced types + software-architecture)

```
Read .cursor/skills/software-architecture/SKILL.md
Read .agents/skills/software-architecture/references/wpd-message-gateway.md
Read .cursor/skills/typescript-react-reviewer/SKILL.md
Read .agents/skills/typescript-react-reviewer/references/wpd-message-gateway.md
Read .cursor/skills/typescript-advanced-types/SKILL.md
Read .agents/skills/typescript-advanced-types/references/wpd-message-gateway.md
Follow docs/agents/delivery-agent.md
cd frontend && npm run lint && npm run test
/smell develop → make audit
```

## Architecture / refactor (software-architecture)

```
Read .cursor/skills/software-architecture/SKILL.md
Read .agents/skills/software-architecture/references/wpd-message-gateway.md
Read docs/backend/architecture.md and/or docs/frontend/architecture.md
Follow docs/agents/delivery-agent.md or review-agent.md
/smell develop → make audit
```
