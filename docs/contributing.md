# Contributing Guide

## Quick Start

```bash
git clone git@github.com:weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make install
```

## Adding a New Provider

Adding a new provider **does not require modifying existing code** — just create your provider files. This follows the Open/Closed Principle (OCP).

### 1. Create the Provider

```bash
mkdir -p internal/infrastructure/provider/sendgrid
```

```go
// internal/infrastructure/provider/sendgrid/sendgrid.go
package sendgrid

import (
    "context"
    "fmt"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const ProviderName = "sendgrid"

// Config holds SendGrid configuration.
type Config struct {
    APIKey    string
    FromEmail string
    FromName  string
}

// Provider implements port.EmailSender for SendGrid.
type Provider struct {
    config Config
}

// New creates a new SendGrid provider.
func New(cfg Config) (*Provider, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("sendgrid: API key required")
    }
    return &Provider{config: cfg}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
    return ProviderName
}

// Send sends an email via SendGrid.
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
    // Your SendGrid implementation here
    return &contracts.SendResult{
        ID:      "msg-123",
        Message: "sent",
    }, nil
}
```

### 2. Register the Provider (Self-Registration)

Create a `register.go` file that uses `init()` to self-register:

```go
// internal/infrastructure/provider/sendgrid/register.go
package sendgrid

import (
    "github.com/weprodev/wpd-message-gateway/internal/core/port"
    "github.com/weprodev/wpd-message-gateway/internal/app/registry"
)

func init() {
    registry.RegisterEmailProvider("sendgrid", func(cfg registry.EmailConfig, _ port.MessageStore, _ registry.MailpitConfig) (port.EmailSender, error) {
        return New(Config{
            APIKey:    cfg.APIKey,
            FromEmail: cfg.FromEmail,
            FromName:  cfg.FromName,
        })
    })
}
```

### 3. Add Import in imports.go

Add a blank import to `internal/app/imports.go`:

```go
// internal/app/imports.go
package app

import (
    // Built-in providers
    _ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/mailgun"
    _ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
    _ "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/sendgrid"  // ← Add this
)
```

### 4. Configure in YAML

The configuration system already supports any key-value pairs. Just add to your `configs/local.yml`:

```yaml
providers:
  defaults:
    email: sendgrid
  email:
    sendgrid:
      api_key: "your-api-key"
      from_email: "noreply@example.com"
      from_name: "My App"
```

**Note:** No changes to `config.go` are needed! The `CommonConfig.Extra` map captures any additional fields.

### 5. Add Tests

```go
// internal/infrastructure/provider/sendgrid/sendgrid_test.go
package sendgrid

import "testing"

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

### 6. Run Quality Checks

```bash
make audit  # Runs fmt, lint, test, vulncheck
```

## Architecture: Why Self-Registration?

The registry pattern follows **SOLID principles**:

| Principle | How We Apply It |
|-----------|-----------------|
| **Open/Closed** | Add providers without modifying existing code |
| **Single Responsibility** | Each provider manages only its registration |
| **Dependency Inversion** | Providers depend on `port` interfaces, not concrete types |

```
┌─────────────────┐     registers      ┌──────────────┐
│ sendgrid/       │─────────────────►  │ app/         │
│  register.go    │   via init()       │  providers   │
│  sendgrid.go    │                    │  (registry)  │
└─────────────────┘                    └──────────────┘
```

## Project Structure

```
internal/infrastructure/provider/
├── mailgun/
│   ├── mailgun.go      # Implementation
│   └── register.go     # Self-registration
├── memory/
│   ├── store.go        # Shared message store
│   ├── email.go        # Memory email provider
│   ├── register.go     # Self-registration
│   └── ...
└── sendgrid/           # Your new provider
    ├── sendgrid.go     # Implementation
    ├── sendgrid_test.go
    └── register.go     # Self-registration ← KEY FILE
```

## Pull Request

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

```
feat(sendgrid): add email provider    → Minor release
fix(mailgun): handle rate limits      → Patch release
docs: update usage guide              → Patch release
feat!: change API response format     → Major release
```

See [Workflow Guide](./workflow.md) for full details.

### Checklist

- [ ] `make audit` passes
- [ ] Tests added with good coverage
- [ ] Interface check: `var _ port.EmailSender = (*Provider)(nil)`
- [ ] `register.go` with `init()` for self-registration
- [ ] Blank import added to `internal/app/imports.go`
- [ ] Commit messages follow [conventions](./workflow.md#commit-conventions)

## Related Documentation

- [Architecture](./architecture.md) — System design
- [Code Conventions](./code-conventions.md) — Coding standards
- [Workflow Guide](./workflow.md) — CI/CD, commits, and releases
