# Agent verification chain

Mandatory sequence for **every agent** after code changes — implementation, bugfix, refactor, or post-review fixes.

Agents **must run these themselves**. Do not ask the user to run them. Do not mark work complete until all steps pass.

**Command reference:** `.claude/commands/smell.md` · **`/smell [branch]`** (default base: `develop`)

---

## When to run

| Phase | Steps |
| ----- | ----- |
| **After implementation / fix** | Fast lint → `/smell` → fix BLOCKER/HIGH → `make audit` |
| **Code review (PR or diff)** | `/smell` **first** → [review-agent](./review-agent.md) checklists |
| **Review → fix cycle** | Apply fixes → repeat full chain below |
| **Before PR or “done”** | Full chain; zero BLOCKER/HIGH smell findings |

---

## Chain (in order)

### 1. Fast lint (touched paths)

Go:

```bash
golangci-lint run ./cmd/... ./internal/... ./pkg/...
go test -race ./cmd/... ./internal/... ./pkg/...
```

Frontend:

```bash
cd frontend && npm run lint && npm run test
```

Fix compile/lint failures before smell — smell is not a substitute for lint.

### 2. `/smell develop`

Follow `.claude/commands/smell.md` on the git diff vs `develop`.

- Produce the **Smell Report**
- **Fix all BLOCKER and HIGH** catalog findings (`MGW.*`, `GO.*`, `CC.*`, `TS.*`)
- Re-run `/smell` after fixes until clean or only MEDIUM/LOW/NIT remain

### 3. `make audit`

From repository root:

```bash
make audit
```

Exit code **0** required. Fix failures and repeat from step 2 if the diff changed materially.

### 4. Report to user

State explicitly:

- Smell summary (findings count / clean)
- `make audit` result
- Any MEDIUM items deferred (with reason)

---

## Specialist routing

| Agent | Smell role |
| ----- | ---------- |
| [master-agent](./master-agent.md) | Enforces this chain on all delegated work |
| [delivery-agent](./delivery-agent.md) | Runs chain after §4.3 Implement |
| [review-agent](./review-agent.md) | Runs `/smell` before §2 checklists |

See [prompts.md](./prompts.md) for copy-paste snippets.
