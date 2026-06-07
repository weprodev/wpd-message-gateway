# WPD Message Gateway — Go project overlay

Project-specific rules for **golang-pro** in this repository. These **override** generic microservice/gRPC guidance where they conflict.

**Canonical project docs:** `docs/backend/architecture.md`, `docs/backend/code-conventions.md`, `docs/backend/backend-engineer.md`, `docs/agents/verification.md`

---

## Dual-mode gateway

| Mode | Package | Requires DB/server |
| ---- | ------- | ------------------ |
| Embedded SDK | `pkg/gateway`, `pkg/contracts`, `pkg/registry`, `pkg/provider/*` | No |
| HTTP server | `cmd/server`, `internal/*` | Yes (PostgreSQL) |

**Iron rule:** `pkg/*` must **never** import `internal/*`. External Go consumers use only `pkg/`.

---

## Layer map (Hexagonal / DDD)

```
Router → Middleware → Handler → Service → Port ← Repository / Provider
```

| Layer | Path | May import | Must not |
| ----- | ---- | ---------- | -------- |
| Domain | `internal/core/domain/` | stdlib only | Echo, HTTP, infra |
| Port | `internal/core/port/` | domain, `pkg/contracts` | implementations |
| Service | `internal/core/service/` | port, domain, `pkg/contracts`, `pkg/registry` | SQL, HTTP, Echo |
| Infrastructure | `internal/infrastructure/` | port, domain | presentation |
| Presentation | `internal/presentation/` | service, port, domain | direct SQL |
| Public SDK | `pkg/*` | `pkg/*` only | `internal/*` |

Sender interfaces (`EmailSender`, `SMSSender`, …) live in **`pkg/contracts/`**, not `internal/core/port/`.

---

## Provider registration

1. Implement under **`pkg/provider/<name>/`**
2. Register via `init()` in `register.go` → `pkg/registry`
3. Blank-import in **`internal/app/imports.go`** (server mode only)

Do **not** add providers under `internal/infrastructure/provider/` (deprecated).

---

## Idiomatic patterns (this repo)

### Context & logging

- First param: `ctx context.Context` on all blocking/outbound calls
- Use `slog.InfoContext` / `slog.ErrorContext` — global default is wired in `internal/infrastructure/logger`
- Enrich ctx: `logger.WithRequestID`, `WithWorkspace`, `WithChannel`, `WithProvider`
- Never log JWTs, API keys, or provider secrets

### Errors

- Wrap with `fmt.Errorf("…: %w", err)`
- Domain/port sentinels: `port.ErrNotFound`, `port.ErrConflict`, etc.
- Handlers: generic messages on 500; full error in slog only (`send_helper.go` pattern)

### Concurrency

- SSE subscribers in `portal_inbox_handler.go`: respect `r.Context().Done()`
- Provider cache (`provider_cache.go`): RWMutex; factory outside write lock
- No unbounded goroutines; use buffered channels or `errgroup` for pipelines

### Value vs pointer semantics

- **Messages** (`contracts.Email`, `SMS`, …): pass **by value** into `Send` and `InboxWriter` — required inputs, no nil.
- **Results** (`*contracts.SendResult`): keep pointer when metadata is stamped after send.
- **Inbox storage** (`port.Stored*`): embed message **values**; slices are `[]StoredEmail` not `[]*StoredEmail`; lookups return `(StoredEmail, bool)`.
- **Services/handlers**: decode JSON into stack values (`var req contracts.Email`), pass by value to services.

### Testing

- Table-driven tests with `t.Run` subtests
- Run race detector: `go test -race ./cmd/... ./internal/... ./pkg/...`
- Prefer fakes implementing `port.*` interfaces over mocking frameworks

---

## Verification chain (mandatory after Go changes)

Agents run this themselves — do not ask the user:

```bash
golangci-lint run ./cmd/... ./internal/... ./pkg/...
go test -race ./cmd/... ./internal/... ./pkg/...
/smell develop          # fix BLOCKER/HIGH
make audit              # from repo root
```

Project gate supersedes generic "80%+ coverage" targets: **zero BLOCKER/HIGH smell findings** and **`make audit` exit 0**.

---

## Layout reference

```
wpd-message-gateway/
├── cmd/server/           # HTTP entry
├── internal/
│   ├── app/              # wire, imports, config
│   ├── core/             # domain, port, service
│   ├── infrastructure/   # postgres repos, logger, inbox store
│   └── presentation/     # router, handlers, middleware
├── pkg/
│   ├── contracts/        # message types + sender interfaces
│   ├── gateway/          # embedded SDK
│   ├── registry/         # provider factory registry
│   └── provider/         # mailgun, memory, …
└── database/migrations/
```

---

## Architecture complement

Load **software-architecture** + [architecture overlay](../../software-architecture/references/wpd-message-gateway.md) for new modules, refactors, and layer moves.

## When implementing Go here

1. Read `docs/backend/architecture.md` for the change surface (SDK vs server)
2. Place code in the correct layer; add port interface before infra adapter
3. Implement with context, wrapped errors, structured logs
4. Add table-driven tests in the same package
5. Run verification chain before marking done
