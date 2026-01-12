# Code Conventions

This document defines the coding standards for Go Message Gateway.

## Go Style

Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).

## Naming

### Packages

```go
// Good: lowercase, single word
package mailgun
package config

// Bad: underscores, mixed case
package mail_gun
package MailGun
```

### Interfaces

```go
// Good: verb-based, single method = -er suffix
type EmailSender interface {
    Send(ctx context.Context, email *Email) (*SendResult, error)
    Name() string
}

// Bad: noun-based
type EmailService interface {}
```

### Constants

```go
// Good: exported, descriptive
const ProviderName = "mailgun"
const DefaultTimeout = 30 * time.Second

// Bad: unexported for public API
const providerName = "mailgun"
```

### Errors

```go
// Good: start with package/context
return fmt.Errorf("mailgun: failed to send: %w", err)

// Bad: generic
return fmt.Errorf("error: %v", err)
```

## Interface Compliance

Always add compile-time checks:

```go
var _ contracts.EmailSender = (*Provider)(nil)
```

## Context

Always accept `context.Context` as first parameter:

```go
// Good
func (p *Provider) Send(ctx context.Context, email *Email) (*SendResult, error)

// Bad
func (p *Provider) Send(email *Email) (*SendResult, error)
```

## Error Handling

### Use Custom Error Types

```go
// Good: rich error with context
return msgerrors.NewProviderError(ProviderName, "rate limited", 429, err)

// Bad: generic error
return errors.New("failed")
```

### Wrap Errors

```go
// Good: preserve error chain
return fmt.Errorf("mailgun: send failed: %w", err)

// Bad: lose original error
return fmt.Errorf("send failed: %v", err)
```

## Validation

Validate at provider boundary:

```go
func (p *Provider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
    // Validate first
    if len(email.To) == 0 {
        return nil, msgerrors.NewProviderError(ProviderName, "recipient required", 400, nil)
    }
    
    // Then process...
}
```

## File Organization

### Provider Structure

```
providers/email/mailgun/
├── mailgun.go       # Main implementation
├── mailgun_test.go  # Unit tests
└── doc.go           # Package documentation (optional)
```

### Test Files

- Same package for unit tests: `mailgun_test.go`
- Integration tests: use build tag `//go:build integration`

## Testing

### Table-Driven Tests

```go
func TestNew(t *testing.T) {
    tests := []struct {
        name    string
        cfg     config.ProviderConfig
        wantErr bool
    }{
        {"valid", config.ProviderConfig{APIKey: "key", Domain: "d"}, false},
        {"missing key", config.ProviderConfig{Domain: "d"}, true},
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

### Mock Interfaces

```go
type mockEmailSender struct {
    contracts.EmailSender
    sendFunc func(ctx context.Context, e *contracts.Email) (*contracts.SendResult, error)
}

func (m *mockEmailSender) Send(ctx context.Context, e *contracts.Email) (*contracts.SendResult, error) {
    return m.sendFunc(ctx, e)
}
```

## Documentation

### Package Comments

```go
// Package mailgun implements the EmailSender interface using the Mailgun API.
//
// Configuration requires MESSAGE_MAILGUN_API_KEY and MESSAGE_MAILGUN_DOMAIN
// environment variables.
package mailgun
```

### Function Comments

```go
// New creates a new Mailgun provider from the given configuration.
// Returns an error if APIKey or Domain is empty.
func New(cfg config.ProviderConfig) (*Provider, error)
```

## Commit Messages

```
feat(mailgun): add attachment support
fix(config): handle empty provider names
docs: add architecture diagram
test(manager): add concurrent access tests
refactor(errors): simplify error types
```

## Quality Checks

Before submitting:

```bash
make fmt        # Format code
make lint       # Run linter
make test       # Run tests
make vulncheck  # Check vulnerabilities
make audit      # All of the above
```
