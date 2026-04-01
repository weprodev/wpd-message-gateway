# Usage Guide

## Two Ways to Use Message Gateway

```
┌──────────────────────────────────┬──────────────────────────────────┐
│  Embedded Go SDK                 │  HTTP Server (any language)      │
│  (no server needed)              │  (full infrastructure)           │
├──────────────────────────────────┼──────────────────────────────────┤
│  go get .../wpd-message-gateway  │  git clone → make start          │
│  gateway.New(config)             │  POST /v1/email (HTTP)           │
│  Config lives in your code       │  Config lives in PostgreSQL      │
│  No DB required                  │  React Portal UI to manage       │
└──────────────────────────────────┴──────────────────────────────────┘
```

---

## Option 1: Go Package (Embedded SDK)

Use this when you have a **Go application** and want to send messages directly — no server setup required.

### Install

```bash
go get github.com/weprodev/wpd-message-gateway
```

### Basic Usage

```go
package main

import (
    "context"
    "log"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/gateway"
)

func main() {
    gw, err := gateway.New(gateway.Config{
        DefaultEmailProvider: "mailgun",
        EmailProviders: map[string]gateway.EmailConfig{
            "mailgun": {
                CommonConfig: gateway.CommonConfig{APIKey: "key-xxx"},
                Domain:       "mg.example.com",
                FromEmail:    "noreply@example.com",
                FromName:     "MyApp",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

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

### Using Memory Provider (Dev / Tests)

```go
// No external service — messages are captured in RAM
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "memory",
})
gw.SendEmail(ctx, email) // captured locally, not sent
```

### Provider Configuration

Pass credentials directly to `gateway.New()`. Use environment variables or your secrets manager — never hard-code credentials.

```go
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: os.Getenv("EMAIL_PROVIDER"), // "mailgun", "memory", etc.
    EmailProviders: map[string]gateway.EmailConfig{
        "mailgun": {
            CommonConfig: gateway.CommonConfig{
                APIKey: os.Getenv("MAILGUN_API_KEY"),
            },
            Domain:    os.Getenv("MAILGUN_DOMAIN"),
            FromEmail: os.Getenv("MAILGUN_FROM_EMAIL"),
            FromName:  os.Getenv("MAILGUN_FROM_NAME"),
        },
    },
})
```

---

## Option 2: HTTP Server

Use this when you want to:
- Send messages from **any language** (Python, PHP, Ruby, JS, etc.)
- Manage provider credentials via a **UI** (not in code)
- Have a **team** managing multiple workspaces
- Use the built-in **message inbox** (dev/testing)

### Setup

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
cp configs/local.example.yml configs/local.yml
make start
```

Open **http://localhost:10104** — the Portal UI.

### First Time Setup (Portal)

1. Go to **http://localhost:10104**
2. Register (email + password), then sign in
3. **Create a workspace** (e.g. "myapp")
4. Go to **Integrations** → add your email provider (Mailgun, etc.)
5. Go to **API Keys** → create an API key (copy the secret — shown once)
6. Go to **Settings** → set dispatch mode (`provider_only` for real sends)

### Sending via HTTP

```bash
curl -X POST http://localhost:10101/v1/email \
  -u "wk_abc123:your-secret" \
  -H "X-Workspace-Key: myapp" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "Hello",
    "html": "<h1>World</h1>"
  }'
```

Or with Bearer token:

```bash
curl -X POST http://localhost:10101/v1/email \
  -H "Authorization: Bearer wk_abc123:your-secret" \
  -H "X-Workspace-Key: myapp" \
  -H "Content-Type: application/json" \
  -d '{"to": ["user@example.com"], "subject": "Hello", "html": "<h1>World</h1>"}'
```

### Sending from Python

```python
import requests

response = requests.post(
    "http://localhost:10101/v1/email",
    auth=("wk_abc123", "your-secret"),
    headers={"X-Workspace-Key": "myapp"},
    json={
        "to": ["user@example.com"],
        "subject": "Hello from Python",
        "html": "<h1>Hello!</h1>"
    }
)
print(response.json())
```

### Sending from Node.js / TypeScript

```typescript
const response = await fetch("http://localhost:10101/v1/email", {
  method: "POST",
  headers: {
    "Authorization": "Basic " + btoa("wk_abc123:your-secret"),
    "X-Workspace-Key": "myapp",
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    to: ["user@example.com"],
    subject: "Hello from JS",
    html: "<h1>Hello!</h1>",
  }),
});
const result = await response.json();
```

---

## Authentication

### Portal (UI & Management API `/api/v1/*`)

Authentication is **email + password** — the Portal returns a JWT used for all `/api/v1/*` requests.

```bash
# Login (or register first, then login)
curl -X POST http://localhost:10101/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "your-password"}'

# Response:
# { "token": "<jwt>", "user": { "id": "...", "email": "..." } }
```

Use the JWT for all `/api/v1/*` requests:
```bash
curl -H "Authorization: Bearer <jwt>" http://localhost:10101/api/v1/workspaces
```

### Send API (`/v1/*`)

Use your **workspace API key** credentials:
```bash
# Basic auth
curl -u "wk_abc123:secret" -H "X-Workspace-Key: myapp" ...

# Bearer (colon-separated)
curl -H "Authorization: Bearer wk_abc123:secret" -H "X-Workspace-Key: myapp" ...
```

---

## Message Types

### Email

```go
// Go SDK
result, err := gw.SendEmail(ctx, &contracts.Email{
    To:          []string{"user@example.com"},
    CC:          []string{"cc@example.com"},         // optional
    BCC:         []string{"bcc@example.com"},        // optional
    Subject:     "Hello",
    HTML:        "<h1>HTML body</h1>",
    PlainText:   "Plain text fallback",              // optional
    ReplyTo:     "reply@example.com",                // optional
})
```

```bash
# HTTP API
curl -X POST /v1/email -d '{
  "to": ["user@example.com"],
  "cc": ["cc@example.com"],
  "subject": "Hello",
  "html": "<h1>HTML body</h1>",
  "plain_text": "Plain text fallback"
}'
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
    To:      []string{"#channel"},
    Message: "Hello from the gateway!",
})
```

---

## HTTP API Reference

### Send Endpoints (API Key auth required)

| Method | Endpoint | Auth |
|--------|----------|------|
| POST | `/v1/email` | API key + `X-Workspace-Key` |
| POST | `/v1/sms` | API key + `X-Workspace-Key` |
| POST | `/v1/push` | API key + `X-Workspace-Key` |
| POST | `/v1/chat` | API key + `X-Workspace-Key` |

### Portal Endpoints (JWT auth required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Create portal user (returns JWT) |
| POST | `/api/v1/auth/login` | Login portal user (returns JWT) |
| GET | `/api/v1/auth/me` | Current user |
| GET | `/api/v1/workspaces` | List my workspaces |
| POST | `/api/v1/workspaces` | Create workspace |
| POST | `/api/v1/workspaces/join` | Join workspace with PIN |
| GET | `/api/v1/workspaces/:wid` | Get workspace |
| PATCH | `/api/v1/workspaces/:wid` | Update workspace |
| GET | `/api/v1/workspaces/:wid/members` | List members |
| DELETE | `/api/v1/workspaces/:wid/members/:uid` | Remove member |
| GET | `/api/v1/workspaces/:wid/api-keys` | List API keys |
| POST | `/api/v1/workspaces/:wid/api-keys` | Create API key |
| DELETE | `/api/v1/workspaces/:wid/api-keys/:kid` | Delete API key |
| POST | `/api/v1/workspaces/:wid/api-keys/:kid/regenerate` | Regenerate secret |
| GET | `/api/v1/workspaces/:wid/integrations` | List integrations |
| POST | `/api/v1/workspaces/:wid/integrations` | Create/update integration |
| DELETE | `/api/v1/workspaces/:wid/integrations/:iid` | Delete integration |
| GET | `/api/v1/workspaces/:wid/templates` | List templates |
| POST | `/api/v1/workspaces/:wid/templates` | Create template |
| PATCH | `/api/v1/workspaces/:wid/templates/:tid` | Update template |
| DELETE | `/api/v1/workspaces/:wid/templates/:tid` | Delete template |
| GET | `/api/v1/workspaces/:wid/settings` | Get settings |
| PATCH | `/api/v1/workspaces/:wid/settings` | Update settings |
| GET | `/api/v1/workspaces/:wid/logs` | Message logs |

### Portal Inbox Endpoints (JWT + member + API key required)

See [Portal inbox](./portal-inbox.md) for full reference. Base: `/api/v1/workspaces/:wid/inbox/`

---

## Configuration

### Server (`configs/local.yml`)

```yaml
environment: local

server:
  port: 10101

portal:
  jwt_secret: "your-long-secret-here-minimum-32-chars"
  jwt_ttl_hours: 72
  ui_port: 10104
```

> **Note:** Provider credentials (Mailgun API keys, etc.) are **NOT** in YAML files.
> They are configured in the Portal UI and stored encrypted in PostgreSQL.

### Environment Variables (Server)

```bash
# PostgreSQL connection
DATABASE_URL=postgres://user:pass@localhost:5432/gateway?sslmode=disable
# or individual components:
DB_HOST=localhost
DB_PORT=5432
DB_USER=gateway
DB_PASSWORD=secret
DB_NAME=gateway

# Portal JWT (override yaml)
MESSAGE_JWT_SECRET=your-jwt-secret

# AES encryption key for provider credentials (32 bytes)
MESSAGE_CONFIG_ENCRYPTION_KEY=your-32-byte-key-here

# Internal ingest secret (optional)
MESSAGE_INTERNAL_INGEST_SECRET=your-ingest-secret
```

---

## Message Dispatch Modes

Control how messages are handled per workspace (set in Portal → Settings):

| Mode | Behavior | Use Case |
|------|----------|----------|
| `memory_only` | In-process RAM only, **not sent** | Development, testing |
| `provider_only` | Sent to real provider, no local copy | Production |
| `memory_and_provider` | Both — local copy + real send | Staging, debugging |

---

## Troubleshooting

### Mailgun: 403 Forbidden

Using a Mailgun Sandbox Domain:
1. Log into Mailgun → **Sending → Domains → [Sandbox]**
2. Add recipient email to **Authorized Recipients**
3. Recipient must click the verification link

### Provider Not Sending

Check workspace dispatch mode in Portal → Settings. If it's `memory_only`, messages are captured locally — not sent to the provider.

### API Key Rejected

- Ensure `X-Workspace-Key` matches the workspace's `unique_key` (slug, not UUID)
- Verify the API key is active (Portal → API Keys)
- Check `client_id` and `client_secret` — secret is shown once at creation

---

## Related Documentation

- [Architecture](./architecture.md) — System design
- [Portal inbox](./portal-inbox.md) — Memory capture and inbox UI
- [E2E Testing](./e2e-testing.md) — Automated testing patterns
- [Contributing](./contributing.md) — Adding providers
