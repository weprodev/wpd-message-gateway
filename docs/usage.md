# Usage Guide

This guide covers how to use Go Message Gateway in your projects.

## Installation

```bash
go get github.com/weprodev/wpd-message-gateway
```

## Quick Start

### 1. Configure Environment

```bash
# Required: Set default provider and credentials
export MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun
export MESSAGE_MAILGUN_API_KEY=your-api-key
export MESSAGE_MAILGUN_DOMAIN=mg.yourdomain.com
export MESSAGE_MAILGUN_FROM_EMAIL=noreply@yourdomain.com
export MESSAGE_MAILGUN_FROM_NAME=YourApp

# Optional: EU region
export MESSAGE_MAILGUN_BASE_URL=https://api.eu.mailgun.net
```

### 2. Send Your First Email

```go
package main

import (
    "context"
    "log"

    "github.com/weprodev/wpd-message-gateway/config"
    "github.com/weprodev/wpd-message-gateway/contracts"
    "github.com/weprodev/wpd-message-gateway/manager"
)

func main() {
    cfg, _ := config.LoadFromEnv()
    mgr, _ := manager.New(cfg)

    result, err := mgr.SendEmail(context.Background(), &contracts.Email{
        To:      []string{"user@example.com"},
        Subject: "Welcome!",
        HTML:    "<h1>Hello!</h1>",
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Sent! ID: %s", result.ID)
}
```

## Configuration Methods

### Environment Variables (Recommended)

```bash
MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun
MESSAGE_DEFAULT_SMS_PROVIDER=twilio
MESSAGE_DEFAULT_CHAT_PROVIDER=whatsapp

MESSAGE_{PROVIDER}_API_KEY=xxx
MESSAGE_{PROVIDER}_DOMAIN=xxx
MESSAGE_{PROVIDER}_FROM_EMAIL=xxx
MESSAGE_{PROVIDER}_FROM_NAME=xxx
MESSAGE_{PROVIDER}_BASE_URL=xxx  # optional
```

### Programmatic Configuration

```go
cfg := config.NewConfig()
cfg.DefaultEmailProvider = "mailgun"
cfg.AddProvider("mailgun", config.ProviderConfig{
    APIKey:    os.Getenv("MAILGUN_API_KEY"),
    Domain:    "mg.yourdomain.com",
    FromEmail: "noreply@yourdomain.com",
    FromName:  "My App",
})
```

## Sending Messages

### Email

```go
result, err := mgr.SendEmail(ctx, &contracts.Email{
    To:          []string{"user@example.com"},
    CC:          []string{"cc@example.com"},
    BCC:         []string{"bcc@example.com"},
    Subject:     "Subject",
    HTML:        "<h1>HTML Body</h1>",
    PlainText:   "Plain text fallback",
    ReplyTo:     "reply@example.com",
    Attachments: []contracts.Attachment{
        {Filename: "doc.pdf", Data: pdfBytes},
    },
})
```

### SMS (Coming Soon)

```go
result, err := mgr.SendSMS(ctx, &contracts.SMS{
    To:      []string{"+1234567890"},
    Message: "Your verification code is 123456",
})
```

### Push Notification (Coming Soon)

```go
result, err := mgr.SendPush(ctx, &contracts.PushNotification{
    DeviceTokens: []string{"token1", "token2"},
    Title:        "New Message",
    Body:         "You have a new message",
    Data:         map[string]string{"action": "open_chat"},
})
```

### Chat / Social Media (Coming Soon)

```go
result, err := mgr.SendChat(ctx, &contracts.ChatMessage{
    To:      []string{"+1234567890"},  // WhatsApp
    Message: "Hello from Go Message Gateway!",
})
```

## Using Specific Providers

```go
// Use a specific provider instead of default
result, err := mgr.SendEmailWith(ctx, "sendgrid", email)

// Get provider instance directly
mailgun, err := mgr.EmailProvider("mailgun")
result, err := mailgun.Send(ctx, email)

// List available providers
providers := mgr.AvailableEmailProviders()
```

## Custom Providers

Implement the contract interface:

```go
type MyEmailProvider struct{}

func (p *MyEmailProvider) Send(ctx context.Context, email *contracts.Email) (*contracts.SendResult, error) {
    // Your implementation
    return &contracts.SendResult{ID: "123", Message: "sent"}, nil
}

func (p *MyEmailProvider) Name() string {
    return "my-provider"
}

// Register it
mgr.RegisterEmailProvider("my-provider", &MyEmailProvider{})
```

## Error Handling

```go
result, err := mgr.SendEmail(ctx, email)
if err != nil {
    var providerErr *errors.ProviderError
    if errors.As(err, &providerErr) {
        log.Printf("Provider %s failed: %s (code: %d)", 
            providerErr.Provider, providerErr.Message, providerErr.StatusCode)
    }
    return err
}
```

## Troubleshooting

### Mailgun: 403 Forbidden / "Free accounts are for test purposes only"
If you receive this error: 'Domain ... is not allowed to send', it means you are using a Sandbox Domain.
**Solution**:
1.  Log in to Mailgun.
2.  Go to **Sending > Domains > [Your Sandbox Domain]**.
3.  Add the recipient's email to **Authorized Recipients**.
4.  The recipient must click the verification link sent by Mailgun.

### Config Issues
Ensure your `.env` keys match the provider requirements (e.g., `MESSAGE_MAILGUN_API_KEY`). Use `make sandbox` to verify configuration interactively.

## Related Documentation

- [Architecture](./architecture.md) - How the package is designed
- [Code Conventions](./code-conventions.md) - Coding standards
