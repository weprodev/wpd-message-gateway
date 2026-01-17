# Code Conventions

Follow [Effective Go](https://go.dev/doc/effective_go) plus these project-specific rules.

## Project Structure

```
internal/           # Private application code
├── app/            # Configuration, wiring, validation
├── core/           # Business logic (domain)
│   ├── port/       # Interface definitions
│   └── service/    # Business services
├── infrastructure/ # External integrations
│   └── provider/   # Provider implementations
└── presentation/   # HTTP layer
    └── handler/    # Request handlers

pkg/                # Public packages
├── contracts/      # Message types (single source of truth)
├── errors/         # Structured error types
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

## Context

First parameter for any I/O operation:

```go
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error)
```

## Errors

Use structured errors from `pkg/errors`:

```go
import pkgerrors "github.com/weprodev/wpd-message-gateway/pkg/errors"

// Provider errors
return nil, pkgerrors.NewProviderError("mailgun", "failed to send", 500, err)

// Config errors
return nil, pkgerrors.NewConfigError("mailgun", "api_key", "required")

// Simple errors with context
return fmt.Errorf("mailgun: failed to send: %w", err)  // Good
return fmt.Errorf("error: %v", err)                    // Bad
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
