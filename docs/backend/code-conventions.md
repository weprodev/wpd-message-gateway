# Code Conventions

Follow [Effective Go](https://go.dev/doc/effective_go) plus these project-specific rules.

## Project Structure

```
internal/           # Private application code (HTTP server mode)
├── app/            # Config, wire, validation, provider blank imports (imports.go)
├── core/           # Domain, ports, services
│   ├── domain/     # Domain types
│   ├── port/       # Repository + inbox interfaces, sentinel errors
│   └── service/    # GatewayService, PortalService, provider cache, etc.
├── infrastructure/ # External integrations
│   ├── authgate/   # RBAC adapter
│   ├── inbox/      # In-memory inbox store (port adapter)
│   ├── logger/     # Structured logging
│   └── repository/ # Postgres repositories
└── presentation/   # HTTP layer
    ├── dto/          # Request/response shapes + domain→JSON mappers (never in handler/)
    ├── handler/    # Request handlers (thin: bind DTO → service → DTO)
    └── middleware/ # Auth, JWT, inbox API key

pkg/                # Public packages (embedded SDK + shared contracts)
├── contracts/      # Message types + sender interfaces (EmailSender, …)
├── gateway/        # Embedded SDK entry point
├── registry/       # Provider factory registry (init() self-registration)
└── provider/       # Provider implementations (mailgun, memory, …)
```

## Naming

```go
// Packages: lowercase, single word
package mailgun  // Good
package mail_gun // Bad

// Interfaces: verb-based, in port/
type EmailSender interface { ... }  // Good
type EmailService interface { ... } // Bad

// Constants: exported, descriptive
const ProviderName = "mailgun"
```

## Interfaces

Always add compile-time check:

```go
var _ contracts.EmailSender = (*Provider)(nil)
```

Sender interfaces (`EmailSender`, `SMSSender`, …) live in **`pkg/contracts/`**. Repository and inbox interfaces live in **`internal/core/port/`**.

## Provider Registration

Use the **self-registration pattern** via `init()`:

```go
// pkg/provider/sendgrid/register.go
package sendgrid

import (
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/registry"
)

func init() {
    registry.RegisterEmailProvider("sendgrid", func(cfg registry.EmailConfig) (contracts.EmailSender, error) {
        return New(Config{
            APIKey:    cfg.APIKey,
            FromEmail: cfg.FromEmail,
        })
    })
}
```

**Why:** Follows Open/Closed Principle — add providers without modifying business logic.

**Required:** Add blank import in `internal/app/imports.go`:

```go
_ "github.com/weprodev/wpd-message-gateway/pkg/provider/sendgrid"
```

## Context

First parameter for any I/O operation:

```go
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error)
```

## Errors

Use sentinel errors from `internal/core/port/errors.go` and wrap with `fmt.Errorf`:

```go
import "github.com/weprodev/wpd-message-gateway/internal/core/port"

// Not found — callers can use errors.Is(err, port.ErrNotFound)
return nil, fmt.Errorf("workspace %s: %w", id, port.ErrNotFound)

// Conflict — duplicate key etc.
return nil, fmt.Errorf("email %s already registered: %w", email, port.ErrConflict)

// Wrap provider errors with context
return nil, fmt.Errorf("mailgun: send to %v: %w", email.To, err)

// Always wrap — never swallow
return fmt.Errorf("context: %w", err) // Good
return fmt.Errorf("error: %v", err)   // Bad — breaks errors.Is chain
```

HTTP handlers map sentinels to status codes:

```go
switch {
case errors.Is(err, port.ErrNotFound):           return c.JSON(http.StatusNotFound, ...)
case errors.Is(err, port.ErrConflict):           return c.JSON(http.StatusConflict, ...)
case errors.Is(err, port.ErrUnauthorized):       return c.JSON(http.StatusUnauthorized, ...)
case errors.Is(err, port.ErrInvalidCredentials): return c.JSON(http.StatusUnauthorized, ...)
default:                                        return c.JSON(http.StatusInternalServerError, ...)
}
```

## Validation

Validate at the start of public functions:

```go
func (p *Provider) Send(ctx context.Context, email contracts.Email) (*contracts.SendResult, error) {
    if len(email.To) == 0 {
        return nil, fmt.Errorf("%s: recipient required", ProviderName)
    }
    // ...
}
```

## Types

Use `pkg/contracts` types — they are the single source of truth:

```go
import "github.com/weprodev/wpd-message-gateway/pkg/contracts"

// Good: Use contracts types directly
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error)

// Bad: Creating duplicate types
type Email struct { ... }  // Don't duplicate contracts.Email
```

## Tests

Use table-driven tests:

```go
func TestNew(t *testing.T) {
    tests := []struct {
        name    string
        cfg     Config
        wantErr bool
    }{
        {"valid", Config{APIKey: "key"}, false},
        {"missing key", Config{}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := New(tt.cfg)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Quality Checks

Before committing:

```bash
make audit  # fmt + lint + test + vulncheck
```

## Commit Messages

```
feat(mailgun): add attachment support
fix(config): handle empty values
docs: update readme
test(service): add concurrent tests
refactor(provider): simplify factory
```

## Import Order

```go
import (
    // Standard library
    "context"
    "fmt"

    // External packages
    "github.com/mailgun/mailgun-go/v4"

    // Internal packages (local module)
    "github.com/weprodev/wpd-message-gateway/internal/core/port"
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
)
```

Use `goimports -local github.com/weprodev/wpd-message-gateway` to auto-format.

## HTTP DTOs (presentation layer)

**Never** define request/response structs in `handler/*.go`. HTTP JSON shapes belong in **`internal/presentation/dto/`** — one file per resource (e.g. `api_key.go`, `integration.go`).

| Do | Don't |
| -- | ----- |
| `dto.CreateAPIKeyRequest` in handler bind | `type createAPIKeyBody struct` in handler |
| `dto.IntegrationFromDomain(intg)` | `integrationToJSON` helper in handler |
| Exported PascalCase names (`LoginRequest`) | Unexported `*Body` suffix types in handler |
| Pass `domain.*` to services/repositories | Pass `dto.*` to services or repositories |

### Data flow (DTOs stop at the handler boundary)

Repositories and port interfaces speak **domain only** — never presentation DTOs.

```
Inbound:  JSON → handler.Bind(dto.Request) → req.ToDomain() → domain.* → Service → Port → Repository
Outbound: Repository → domain.* → Service → handler → dto.XFromDomain(domain.*) → JSON
Partial:  JSON → dto.PatchRequest → req.ToPatch() → service.*Patch → Service (load domain, apply, save)
```

| Layer | Types it uses | Maps how |
| ----- | ------------- | -------- |
| **Handler** | `dto.*` | Bind JSON; call `ToDomain` / `ToPatch`; call `XFromDomain` before `c.JSON` |
| **Service** | `domain.*`, `service.*Patch` | Business rules; orchestrates ports |
| **Port / Repository** | `domain.*` | Postgres adapters scan rows ↔ domain structs |
| **Infrastructure** | `domain.*` | No `json` tags required; no import of `presentation/dto` |

**Never send DTOs to a repository.** If a handler has a `dto.CreateTemplateRequest`, call `body.ToDomain(wid)` first, then pass the resulting `*domain.Template` to `PortalService.CreateTemplate`, which passes it to the template repository port. The repository interface looks like:

```go
Create(ctx context.Context, t *domain.Template) error  // port — domain only
```

Gateway `/v1/*` handlers may bind **`pkg/contracts`** types directly — those are the public SDK message API, not portal DTOs.

### Mapping convention (idiomatic Go)

| Direction | Style | Example |
| --------- | ----- | ------- |
| Domain → response DTO | Package func `{Type}FromDomain` | `dto.APIKeyPublicFromDomain(k)` |
| Request DTO → domain (create/upsert) | Value-receiver method `ToDomain()` | `body.ToDomain(wid)` |
| Request DTO → partial update | Value-receiver method `ToPatch()` | `body.ToPatch()` → `service.WorkspacePatch` |

Use **`{Type}FromDomain`** package functions for outbound mapping (not methods on a zero-value receiver). Use **`ToDomain()`** on request DTOs when the result is a **`domain.*`** entity for persistence. Use **`ToPatch()`** when the result is a **service patch command** (service loads domain, applies fields, saves via repository).

**Not infrastructure:** DTOs are presentation concerns (JSON field names, omitempty). `internal/infrastructure/` maps DB rows to **domain**, not HTTP.

New or changed portal endpoints: add/update types in `dto/` and colocate mapper tests (`dto/*_test.go`) when mapping logic is non-trivial.
