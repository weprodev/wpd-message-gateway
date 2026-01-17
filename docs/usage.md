# Usage Guide

## Installation

```bash
go get github.com/weprodev/wpd-message-gateway
```

## Quick Start

### 1. Configure Environment

```bash
# Set default providers
export MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun
export MESSAGE_DEFAULT_SMS_PROVIDER=twilio
export MESSAGE_DEFAULT_CHAT_PROVIDER=whatsapp
export MESSAGE_DEFAULT_PUSH_PROVIDER=firebase

# Provider credentials
export MESSAGE_MAILGUN_API_KEY=your-api-key
export MESSAGE_MAILGUN_DOMAIN=mg.yourdomain.com
export MESSAGE_MAILGUN_FROM_EMAIL=noreply@yourdomain.com
export MESSAGE_MAILGUN_FROM_NAME=YourApp
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

## Sending Messages

### Email

```go
result, err := mgr.SendEmail(ctx, &contracts.Email{
    To:      []string{"user@example.com"},
    Subject: "Subject",
    HTML:    "<h1>HTML Body</h1>",
})
```

### SMS

```go
result, err := mgr.SendSMS(ctx, &contracts.SMS{
    To:      []string{"+1234567890"},
    Message: "Your code is 123456",
})
```

### Push Notification

```go
result, err := mgr.SendPush(ctx, &contracts.PushNotification{
    DeviceTokens: []string{"token1"},
    Title:        "New Message",
    Body:         "You have a new message",
})
```

### Chat (WhatsApp, Telegram)

```go
result, err := mgr.SendChat(ctx, &contracts.ChatMessage{
    To:      []string{"+1234567890"},
    Message: "Hello!",
})
```

## Development Mode

For local development and testing, use the **memory** provider:

```bash
export MESSAGE_DEFAULT_EMAIL_PROVIDER=memory
export MESSAGE_DEFAULT_SMS_PROVIDER=memory
```

Messages are stored in RAM. Use the DevBox UI to view them:

```bash
make server    # Start gateway with DevBox
make web-dev   # Start DevBox UI (http://localhost:5173)
```

See [DevBox](./devbox.md) for more details.

## Troubleshooting

### Mailgun: 403 Forbidden

If using a Mailgun Sandbox Domain:
1. Log in to Mailgun
2. Go to **Sending > Domains > [Your Sandbox Domain]**
3. Add recipient email to **Authorized Recipients**
4. Recipient must click verification link

## Related Documentation

- [Architecture](./architecture.md) - System design
- [DevBox](./devbox.md) - Development UI
