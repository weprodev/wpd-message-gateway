# Usage Guide

## Installation

### As a Go Package

```bash
go get github.com/weprodev/wpd-message-gateway
```

### As an HTTP Server

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make install
```

## Quick Start

### Option 1: Go Package (Embedded SDK)

```go
package main

import (
    "context"
    "log"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/gateway"
)

func main() {
    // Create gateway with configuration
    gw, err := gateway.New(gateway.Config{
        DefaultEmailProvider: "memory",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Send email
    result, err := gw.SendEmail(context.Background(), &contracts.Email{
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

### Option 2: HTTP Server

```bash
# 1. Configure
cp configs/local.example.yml configs/local.yml
# Edit configs/local.yml with your provider settings

# 2. Start server
make start

# 3. Send via HTTP
curl -X POST http://localhost:10101/v1/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "Hello",
    "html": "<h1>World</h1>"
  }'
```

## Configuration

### YAML Configuration

Create `configs/local.yml`:

```yaml
providers:
  defaults:
    email: mailgun
    sms: twilio
    push: firebase
    chat: slack

  email:
    mailgun:
      api_key: "key-xxxxxxxx"
      domain: "mg.yourdomain.com"
      from_email: "noreply@yourdomain.com"
      from_name: "YourApp"
```

### Environment Variable Overrides

Environment variables override YAML values (useful for secrets):

```bash
MESSAGE_MAILGUN_API_KEY=key-xxxxx
MESSAGE_MAILGUN_DOMAIN=mg.example.com
MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun
```

### SDK Configuration

```go
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "mailgun",
    Mailgun: gateway.MailgunConfig{
        APIKey:    "key-xxxxxxxx",
        Domain:    "mg.yourdomain.com",
        FromEmail: "noreply@yourdomain.com",
        FromName:  "YourApp",
    },
})
```

## Sending Messages

### Email

```go
result, err := gw.SendEmail(ctx, &contracts.Email{
    To:        []string{"user@example.com"},
    CC:        []string{"cc@example.com"},
    BCC:       []string{"bcc@example.com"},
    Subject:   "Subject",
    HTML:      "<h1>HTML Body</h1>",
    PlainText: "Plain text fallback",
})
```

### SMS

```go
result, err := gw.SendSMS(ctx, &contracts.SMS{
    To:      []string{"+1234567890"},
    Message: "Your verification code is 123456",
})
```

### Push Notification

```go
result, err := gw.SendPush(ctx, &contracts.PushNotification{
    DeviceTokens: []string{"device-token-1"},
    Title:        "New Message",
    Body:         "You have a new message",
    Data:         map[string]string{"action": "open_chat"},
})
```

### Chat (Slack, WhatsApp, Telegram)

```go
result, err := gw.SendChat(ctx, &contracts.ChatMessage{
    To:      []string{"+1234567890"},
    Message: "Hello from the gateway!",
})
```

### Using a Specific Provider

Override the default provider for a single message:

```go
// Send via SendGrid instead of the default
result, err := gw.SendEmailWith(ctx, "sendgrid", &contracts.Email{...})

// Send via Vonage instead of default SMS provider
result, err := gw.SendSMSWith(ctx, "vonage", &contracts.SMS{...})
```

## Development Mode

For local development and testing, use the **memory** provider:

```yaml
# configs/local.yml
providers:
  defaults:
    email: memory
    sms: memory
    push: memory
    chat: memory
```

Messages are stored in RAM. View them in the DevBox UI:

```bash
make start    # Starts Gateway + DevBox UI
```

Open http://localhost:10104 to see all intercepted messages.

→ See [DevBox](./devbox.md) for more details.

## HTTP API Reference

### Gateway Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/email` | Send email |
| POST | `/v1/sms` | Send SMS |
| POST | `/v1/push` | Send push notification |
| POST | `/v1/chat` | Send chat message |

### DevBox Endpoints (Development Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/stats` | Message counts |
| GET | `/api/v1/emails` | List all emails |
| GET | `/api/v1/sms` | List all SMS |
| GET | `/api/v1/push` | List all push notifications |
| GET | `/api/v1/chat` | List all chat messages |
| DELETE | `/api/v1/messages` | Clear all messages |
| GET | `/api/v1/events` | Real-time updates (SSE) |

## Troubleshooting

### Mailgun: 403 Forbidden

If using a Mailgun Sandbox Domain:
1. Log in to Mailgun
2. Go to **Sending > Domains > [Your Sandbox Domain]**
3. Add recipient email to **Authorized Recipients**
4. Recipient must click verification link

### Provider Not Found

Ensure your configuration has the provider defined:

```yaml
providers:
  defaults:
    email: mailgun   # Must match a key under providers.email
  email:
    mailgun:         # This key must exist
      api_key: "..."
```

## Related Documentation

- [Architecture](./architecture.md) — System design
- [DevBox](./devbox.md) — Development UI
- [Contributing](./contributing.md) — Adding providers
