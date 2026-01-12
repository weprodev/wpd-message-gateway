# Contributing Guide

Thank you for your interest in contributing to Go Message Gateway!

## Getting Started

```bash
# Clone the repository
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway

# Install dependencies and tools
make setup

# Run tests
make test

# Run all quality checks
make audit
```

## Adding a New Provider

### 1. Create Provider Directory

```bash
mkdir -p providers/{type}/{provider}
# Example: providers/email/sendgrid
```

### 2. Implement the Contract

```go
// providers/email/sendgrid/sendgrid.go
package sendgrid

import (
    "context"
    "github.com/weprodev/wpd-message-gateway/config"
    "github.com/weprodev/wpd-message-gateway/contracts"
    msgerrors "github.com/weprodev/wpd-message-gateway/errors"
)

const ProviderName = "sendgrid"

// Compile-time interface check
var _ contracts.EmailSender = (*Provider)(nil)

type Provider struct {
    config config.ProviderConfig
    // ... provider-specific fields
}

func New(cfg config.ProviderConfig) (*Provider, error) {
    if cfg.APIKey == "" {
        return nil, msgerrors.NewConfigError(ProviderName, "APIKey", "required")
    }
    return &Provider{config: cfg}, nil
}

func (p *Provider) Name() string {
    return ProviderName
}

func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
    // Validate input
    if len(email.To) == 0 {
        return nil, msgerrors.NewProviderError(ProviderName, "recipient required", 400, nil)
    }
    
    // Implementation here...
    
    return &contracts.SendResult{
        ID:         "message-id",
        StatusCode: 200,
        Message:    "Email sent successfully",
    }, nil
}
```

### 3. Add Tests

```go
// providers/email/sendgrid/sendgrid_test.go
package sendgrid

import (
    "context"
    "testing"
    "github.com/weprodev/wpd-message-gateway/config"
    "github.com/weprodev/wpd-message-gateway/contracts"
)

func TestNew(t *testing.T) {
    tests := []struct {
        name    string
        cfg     config.ProviderConfig
        wantErr bool
    }{
        {"valid config", config.ProviderConfig{APIKey: "key"}, false},
        {"missing API key", config.ProviderConfig{}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := New(tt.cfg)
            if (err != nil) != tt.wantErr {
                t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestProvider_Send_Validation(t *testing.T) {
    p, _ := New(config.ProviderConfig{APIKey: "key"})
    
    _, err := p.Send(context.Background(), &contracts.Email{})
    if err == nil {
        t.Error("expected error for empty recipients")
    }
}
```

### 4. Register in Manager

Add to `manager/manager.go`:

```go
import "github.com/weprodev/wpd-message-gateway/providers/email/sendgrid"

func (m *Manager) initializeProviders() error {
    for name, cfg := range m.config.Providers {
        switch name {
        case mailgun.ProviderName:
            // existing...
        case sendgrid.ProviderName:  // Add this
            provider, err := sendgrid.New(cfg)
            if err != nil {
                return msgerrors.NewProviderError(name, "failed to initialize", 0, err)
            }
            m.emailProviders[name] = provider
        }
    }
    return nil
}
```

### 5. Add Example

```go
// examples/email/sendgrid/main.go
package main

// Example code...
```

## Pull Request Process

### PR Title Format

```
feat(provider): add SendGrid email provider
fix(mailgun): handle rate limit errors
docs: update architecture diagram
test: add integration tests for Twilio
```

### PR Checklist

- [ ] Code follows [code conventions](./code-conventions.md)
- [ ] Interface compliance check: `var _ contracts.X = (*Provider)(nil)`
- [ ] Unit tests with >80% coverage
- [ ] `make audit` passes (fmt, lint, test, vulncheck)
- [ ] Documentation updated if needed
- [ ] Example added for new provider

### Review Process

1. Create feature branch from `main`
2. Make changes with atomic commits
3. Run `make audit` locally
4. Open PR with description
5. Address review feedback
6. Squash merge when approved

## Testing Guidelines

### Unit Tests

```bash
make test          # All tests
make test-cover    # With coverage
```

### Mock Testing

Use mock implementations for testing:

```go
type mockEmailSender struct {
    sendCalled bool
    lastEmail  *contracts.Email
}

func (m *mockEmailSender) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
    m.sendCalled = true
    m.lastEmail = email
    return &contracts.SendResult{ID: "mock-id"}, nil
}

func (m *mockEmailSender) Name() string { return "mock" }
```

### Integration Tests

Tag integration tests to run separately:

```go
//go:build integration

func TestSendGrid_Integration(t *testing.T) {
    // Real API calls...
}
```

```bash
go test -tags=integration ./...
```

## Related Documentation

- [Architecture](./architecture.md) - Understand the design
- [Code Conventions](./code-conventions.md) - Coding standards
