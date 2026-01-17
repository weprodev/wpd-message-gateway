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
mkdir -p providers/email/sendgrid
```

```go
// providers/email/sendgrid/sendgrid.go
package sendgrid

import (
    "context"
    "github.com/weprodev/wpd-message-gateway/config"
    "github.com/weprodev/wpd-message-gateway/contracts"
)

const ProviderName = "sendgrid"

// Compile-time interface check
var _ contracts.EmailSender = (*Provider)(nil)

type Provider struct {
    apiKey string
}

func New(cfg config.EmailConfig) (*Provider, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("sendgrid: API key required")
    }
    return &Provider{apiKey: cfg.APIKey}, nil
}

func (p *Provider) Name() string {
    return ProviderName
}

func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
    // Your implementation here
    return &contracts.SendResult{ID: "msg-123"}, nil
}
```

### 2. Add Tests

```go
// providers/email/sendgrid/sendgrid_test.go
func TestNew(t *testing.T) {
    _, err := New(config.EmailConfig{})
    if err == nil {
        t.Error("expected error for missing API key")
    }
}
```

### 3. Register in Factory

Add your provider to `manager/factory.go`:

```go
case "sendgrid":
    return sendgrid.New(cfg)
```

### 4. Run Quality Checks

```bash
make audit  # Runs fmt, lint, test, vulncheck
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

See [Workflow Guide](./workflow.md) for full details on commit conventions.

### Checklist

- [ ] `make audit` passes
- [ ] Tests added
- [ ] Interface check: `var _ contracts.X = (*Provider)(nil)`
- [ ] Commit messages follow [conventions](./workflow.md#commit-conventions)

## Related Documentation

- [Code Conventions](./code-conventions.md) - Coding standards
- [Workflow Guide](./workflow.md) - CI/CD, commits, and releases
