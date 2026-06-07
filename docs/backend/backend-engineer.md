# Backend Engineer — Role & Workflow

The operational contract for backend engineering in the WPD Message Gateway. As a Principal Go Engineer working on this project, you are expected to follow strict Domain-Driven Design (DDD) principles, Hexagonal Architecture, and idiomatic Go best practices (KISS, DRY, SOLID).

**Agent skills (load both for backend work):**

- [golang-pro](../../.cursor/skills/golang-pro/SKILL.md) — idioms + [overlay](../../.agents/skills/golang-pro/references/wpd-message-gateway.md) — [skills.sh/jeffallan/claude-skills/golang-pro](https://www.skills.sh/jeffallan/claude-skills/golang-pro)
- [software-architecture](../../.cursor/skills/software-architecture/SKILL.md) — Clean Architecture + DDD + [overlay](../../.agents/skills/software-architecture/references/wpd-message-gateway.md) — [skills.sh/sickn33/antigravity-awesome-skills/software-architecture](https://www.skills.sh/sickn33/antigravity-awesome-skills/software-architecture)

## 1. Execution Model

Use a structured, TDD-focused workflow for all backend modifications:

1. **Understand**: Review existing architecture, domain boundaries, and dependencies.
2. **Plan**: Formulate a complete design before writing code.
3. **Test First (RED-GREEN-REFACTOR)**: Write a failing table-driven test representing the requirement, then write the minimal code to pass it, and lastly refactor for cleanliness and performance.
4. **Implement**: Execute using standard patterns (small interfaces, concrete return types, and context-aware functions).
5. **Verify**: Run [verification chain](../agents/verification.md) — lint → **`/smell develop`** → fix BLOCKER/HIGH → **`make audit`**.
6. **Communicate**: Provide clear, conventionally-formatted commit messages detailing *why* a change was made.

## 2. DDD Layers & Import Rules

The backend strictly enforces Hexagonal Architecture (Ports and Adapters) mapping to Domain-Driven Design concepts.

### Import Matrix

| Layer | Can Import | Cannot Import | Description |
|-------|------------|---------------|-------------|
| **Domain** | stdlib only | *Anything else* | Core behavior, types, consts. Zero dependencies. |
| **Port** | Domain, `pkg/contracts` | Service, Infra, Presentation | Repository + inbox interfaces; sentinel errors. |
| **Service** | Port, Domain, `pkg/contracts`, `pkg/registry` | Infra, Presentation | Business logic orchestrator; coordinates ports. |
| **Infra** | Port, Domain, `pkg/contracts` | Presentation | DB, inbox store, logger, authgate (adapters). |
| **Presentation** | Service, Port, Domain | Infra (direct SQL/provider SDK) | HTTP handlers, middleware, routers. |
| **App** (`internal/app`) | All internal layers + `pkg/*` | — | Wire, config, validation, provider blank imports. |
| **Public SDK** (`pkg/*`) | `pkg/*` only | `internal/*` | Embedded gateway, contracts, registry, providers. |

**The Iron Rule:** Outer layers depend on inner layers. Inner layers *never* depend on outer layers.

## 3. Go Best Practices & Patterns

As a Principal Go Engineer, you must adhere to the following idiomatic patterns:

### Context Management & Concurrency
- **Mandatory `context.Context`**: The first argument of every blocking, outbound, or long-running function must be `ctx context.Context`.
- **No Goroutine Leaks**: Any goroutine spawned must have a mapped lifecycle. Track them using `sync.WaitGroup` or coordinate parallel tasks via `golang.org/x/sync/errgroup`. Always use `defer cancel()` immediately after creating a context with a timeout or cancellation.

### Interface & Struct Design
- **Accept Interfaces, Return Structs**: Define interfaces closely localized to where they are consumed (e.g., in the Service layer), rather than where they are implemented. Return concrete types (`*MyStruct`) instead of interfaces from factories to keep implementation capabilities transparent.
- **Useful Zero Values**: Design structs so they are immediately usable without complex constructors when possible (e.g. `sync.Mutex`, `bytes.Buffer`).
- **PortalDeps Pattern**: For components with >4 dependencies (like Handlers and Services), use a `Deps` struct to keep function signatures clean.
  ```go
  type GatewayHandlerDeps struct {
      GatewayService port.GatewayService
      Inbox          port.Inbox
      Logger         *slog.Logger
  }
  func NewGatewayHandler(deps GatewayHandlerDeps) *GatewayHandler {}
  ```
- **Functional Options Pattern**: Use `WithOption(...)` functions to configure complex structs without polluting constructors.

### Error Handling
- **Wrap Errors (`%w`)**: Always trace errors up the stack using `fmt.Errorf("action %s: %w", context, err)`.
- **Use Sentinels / Custom Types**: Do not return raw string errors. All domain and port errors must return sentinel values (e.g., `port.ErrNotFound`) or custom types implementing `error`. 
- **Verify using `errors.Is` / `errors.As`**:
  - **✅ Correct**: `if errors.Is(err, port.ErrNotFound) { ... }`
  - **❌ Incorrect**: `if strings.Contains(err.Error(), "not found") { ... }`
- **Never Ignore Errors**: Never assign an error to the blank identifier `_` unless explicitly documented why a panic or failure is structurally impossible.

### Structured Logging
All logging must be handled by the `logger` adapter wrapper (`internal/infrastructure/logger`). Write context-aware structured logs:
- **✅ Correct**: `slog.InfoContext(ctx, "message processed", slog.String("msg_id", id))`
- **❌ Incorrect**: `log.Printf("processed: %s", msg.ID)`

## 4. Adding a Provider

1. **Implementation** — Package inside `pkg/provider/<name>/` implementing `pkg/contracts.*Sender` (see `mailgun`, `memory`).
2. **Registration** — `register.go` with `init()` calling registry methods (e.g., `registry.RegisterEmailProvider`).
3. **Wiring** — Add a blank import in [`internal/app/imports.go`](../../internal/app/imports.go).

## 5. Security & Stability

| Rule | Description |
|------|-------------|
| **Graceful Shutdown** | The server must capture `os.Signal` (`SIGINT`, `SIGTERM`) and invoke `server.Shutdown(ctx)` with a timeout to flush pending requests gracefully before dropping connections. |
| **Secrets in DB/Env** | Never commit API keys. Provider configurations reside in the DB and Environment variables. |
| **Fail Fast** | Validate all requests at the presentation boundary before invoking domain services. |
| **SQL Safety** | Always rely on query parameters (`$1, $2`) using `pgx` to prevent injection. |
| **Log Sanitization** | Never log JWTs, API keys, or full Message payload outputs. |

## 6. Testing Strategy

All new development requires testing, achieving >80% coverage where practical:

| Type | Approach |
|------|----------|
| **Unit (TDD)** | Use **Table-Driven Tests** with independent parallel sub-tests (`t.Run(..., func(t *testing.T){ t.Parallel() })`). Rely heavily on `t.Helper()` for assertions and `t.Cleanup()` for deferring teardowns. |
| **Fuzzing & Benchmarks** | Write `func FuzzX(f *testing.F)` routines for complex logic. Track optimization using `func BenchmarkX(b *testing.B)` combined with the `-benchmem` flag to identify high allocation bottlenecks. |
| **Integration** | Test implementations against real Postgres dependencies (using testcontainers if possible). |
| **Provider** | Ensure Providers respect E2E boundaries by testing against the Memory Inbox. |
| **Race Tests** | Run all tests with the `-race` flag to guarantee Concurrency Safety. |

## 7. Migrations

Use additive migrations.
- Create pairs `.up.sql` and `.down.sql`.
- One conceptual change per file pair.
- Name deliberately: `YYYYMMDDHHMMSS_add_table.up.sql`.

## 8. Common Anti-patterns

| Anti-Pattern (To Avoid) | Correct Pattern |
|----------------------------|-----------------|
| Injecting DB directly into handler | Inject Service interface into handler |
| Catching unknown errors with 500s | Map explicitly to Application domains using `port.ErrX` |
| Using `pkg/errors` or `fmt.Errorf("%v")`| Use standard `errors.New` or `errors.Is` with sentinels or `%w` |
| Bypassing the registry for providers | Blank import in `internal/app/imports.go` |
| Printing directly to standard out | Use `slog` injected context |
| `context.Context` inside a struct | Pass `ctx context.Context` explicitly as the first function parameter |
| Naked Returns | Explicitly state returned variables: `return result, err` |
| Leaking Goroutines | Ensure predictable shutdown via `context.Context` and wait groups |

## 9. Decision Heuristics

When unsure how to implement a business feature:
1. Does this belong in pure Domain logic? (Yes if it relies on no state/dependencies).
2. Is this orchestration? (Service layer).
3. Is this an external Side-effect? (Infrastructure Port).
4. Is it a presentation concern? (Handler).

## 10. Quality Gate

Before pushing any PR, run [verification.md](../agents/verification.md):

```bash
/smell develop   # fix BLOCKER/HIGH — .claude/commands/smell.md
make audit       # format, lint, test, govulncheck, builds
```

Agents must execute both themselves. Ensure exit 0 for the scope of your change.

## 11. References

- [Code Conventions](./code-conventions.md)
- [Architecture](./architecture.md)
- [Contributing](./contributing.md)
- [E2E Testing](./e2e-testing.md)
