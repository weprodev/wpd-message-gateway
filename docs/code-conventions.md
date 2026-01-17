# Code Conventions

Follow [Effective Go](https://go.dev/doc/effective_go) plus these project-specific rules.

## Naming

```go
// Packages: lowercase, single word
package mailgun  // Good
package mail_gun // Bad

// Interfaces: verb-based
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

## Context

First parameter for any I/O operation:

```go
func (p *Provider) Send(ctx context.Context, email *Email) (*SendResult, error)
```

## Errors

Include context in error messages:

```go
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

## Tests

Use table-driven tests:

```go
func TestNew(t *testing.T) {
    tests := []struct {
        name    string
        cfg     config.EmailConfig
        wantErr bool
    }{
        {"valid", config.EmailConfig{APIKey: "key"}, false},
        {"missing key", config.EmailConfig{}, true},
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
test(manager): add concurrent tests
```
