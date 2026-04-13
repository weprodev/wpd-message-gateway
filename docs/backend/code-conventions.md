# Code Conventions

Follow [Effective Go](https://go.dev/doc/effective_go) plus these project-specific rules.

## Project Structure

```
internal/           # Private application code
├── app/            # Config, wire, validation, provider blank imports (imports.go)
├── core/           # Domain, ports, services
│   ├── domain/     # Domain types
│   ├── port/       # Interface definitions + sentinel errors
│   └── service/    # GatewayService, PortalService, provider cache, etc.
├── infrastructure/ # External integrations
│   ├── logger/     # Structured logging
│   ├── provider/   # Email/SMS/Push/Chat provider implementations
│   └── repository/ # Postgres repositories
├── presentation/   # HTTP layer
│   ├── handler/    # Request handlers
│   └── middleware/ # Auth, JWT, inbox API key
└── registry/       # Provider factory registry (init() self-registration)

pkg/                # Public packages
├── contracts/      # Message types (single source of truth)
├── auth/           # Hashing helpers (portal)
├── encryption/     # AES helpers
└── gateway/        # Embedded SDK
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
var _ port.EmailSender = (*Provider)(nil)
```

Define interfaces in `internal/core/port/`, not alongside implementations.

## Provider Registration

Use the **self-registration pattern** via `init()`:

```go
// internal/infrastructure/provider/sendgrid/register.go
package sendgrid

import (
    "github.com/weprodev/wpd-message-gateway/internal/core/port"
    "github.com/weprodev/wpd-message-gateway/internal/registry"
)

func init() {
    registry.RegisterEmailProvider("sendgrid", func(cfg registry.EmailConfig, _ registry.MailpitConfig) (port.EmailSender, error) {
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
_ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/sendgrid"
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
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
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
