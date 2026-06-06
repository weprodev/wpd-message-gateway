# Agent prompts

Copy-paste snippets for common workflows. Full chain: **[verification.md](./verification.md)**.

## Finish implementation (mandatory)

Agents must run both — do not ask the user:

```
/smell develop
make audit
```

Follow `.claude/commands/smell.md` and fix BLOCKER/HIGH before `make audit`.

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
