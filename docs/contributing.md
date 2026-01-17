# Contributing Guide

## Quick Start

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make install
make test
```

## Adding a New Provider

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

    "github.com/weprodev/wpd-message-gateway/internal/core/port"
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const ProviderName = "sendgrid"

// Compile-time interface check
var _ port.EmailSender = (*Provider)(nil)

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

### 2. Add Configuration

Update `internal/app/config.go`:

```go
type EmailConfig struct {
    // ... existing fields
    SendGrid SendGridConfig `yaml:"sendgrid"`
}

type SendGridConfig struct {
    APIKey    string `yaml:"api_key"`
    FromEmail string `yaml:"from_email"`
    FromName  string `yaml:"from_name"`
}
```

### 3. Register in Provider Factory

Update `internal/app/providers.go`:

```go
import "github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/sendgrid"

func (f *ProviderFactory) createEmailProvider(name string) (port.EmailSender, error) {
    switch name {
    // ... existing cases
    case "sendgrid":
        cfg := f.config.Providers.Email.SendGrid
        return sendgrid.New(sendgrid.Config{
            APIKey:    cfg.APIKey,
            FromEmail: cfg.FromEmail,
            FromName:  cfg.FromName,
        })
    default:
        return nil, fmt.Errorf("unknown email provider: %s", name)
    }
}
```

### 4. Add Tests

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

### 5. Run Quality Checks

```bash
make audit  # Runs fmt, lint, test, vulncheck
```

## Project Structure

When adding a new provider, follow this structure:

```
internal/infrastructure/provider/
├── mailgun/
│   └── mailgun.go      # Mailgun implementation
├── memory/
│   ├── store.go        # Shared message store
│   ├── email.go        # Memory email provider
│   └── ...
└── sendgrid/           # Your new provider
    ├── sendgrid.go
    └── sendgrid_test.go
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
- [ ] Configuration added to `internal/app/config.go`
- [ ] Provider registered in `internal/app/providers.go`
- [ ] Commit messages follow [conventions](./workflow.md#commit-conventions)

## Related Documentation

- [Architecture](./architecture.md) — System design
- [Code Conventions](./code-conventions.md) — Coding standards
- [Workflow Guide](./workflow.md) — CI/CD, commits, and releases
