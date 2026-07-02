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
│  No DB required                  │  React Portal UI (partial)       │
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

#### Option A: Docker Compose (Recommended)
This spins up the entire stack including a PostgreSQL database (pre-loaded with migrations and permission/demo seeds), a hot-reloading Go backend (via `air`), and the Portal UI dev server:

```bash
# Start all services
make dev

# Stop all services
make dev-down

# Reset Database: Since DB init scripts run only on first boot, run this to reset data/schema:
docker compose down -v
```

Open **http://localhost:10104** — the Portal UI.

#### Option B: Native Setup (No Docker)
If you prefer to run components directly on your host machine:

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
cp configs/local.example.yml configs/local.yml
make start
```

For Option B, you must manually run PostgreSQL locally, apply schema migrations, and apply the SQL seeds (see below).

Open **http://localhost:10104** — the Portal UI.

### First time setup (Portal UI)

The Portal UI currently supports:

1. **Register** and **sign in** (email + password)
2. **Workspaces** — list workspaces you belong to and open one
3. **Integrations** — connect, activate, deactivate, or remove messaging providers
4. **Message logs** — audit trail of gateway send requests (Overview + channel tabs)
5. **Send test** — send a test message through the gateway for the workspace

For local dev with PostgreSQL, run migrations then apply seeds in order (see `database/init-db.sh`): `001_seed_permissions.sql`, `002_seed_providers.sql`, `003_seed_mailgun_config.sql`, and `004_demo_workspace.sql` (optional demo data: `demo@weprodev.com` / `secret`).

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

### Portal REST API (JWT auth)

Routes registered in `internal/presentation/router.go`. **Portal UI pages exist only for auth, workspace list, message logs, and send test** — other management routes are REST-only until UI is built.

#### Auth

| Method | Endpoint | Portal UI | Description |
|--------|----------|:---------:|-------------|
| POST | `/api/v1/auth/register` | ✓ | Create portal user (`first_name`, `last_name`, `email`, `password`) |
| POST | `/api/v1/auth/login` | ✓ | Login (returns JWT) |
| GET | `/api/v1/auth/verify-email` | | Email verification link handler |
| GET | `/api/v1/auth/me` | | Current user (JWT) |

#### Workspaces & messaging (Portal UI)

| Method | Endpoint | Portal UI | Description |
|--------|----------|:---------:|-------------|
| GET | `/api/v1/workspaces` | ✓ | List workspaces for the logged-in user |
| GET | `/api/v1/workspaces/:wid/logs` | ✓ | Message request logs (optionally filter by `channel`) |
| POST | `/api/v1/workspaces/:wid/send-test/:channel` | ✓ | Send test via gateway (`email` / `sms` / `push` / `chat`) |

#### Workspace provisioning & management (REST only — no Portal UI)

Used by Bruno, curl, and CI bootstrap. Not documented endpoint-by-endpoint here; see `internal/presentation/router.go` and [E2E testing](./e2e-testing.md) for create workspace, API keys, settings, integrations, templates, and members.

### Workspace access (portal JWT routes & RBAC)

Workspace-scoped routes require a valid portal JWT and are governed by a Role-Based Access Control (RBAC) middleware backed by `wpd-gogate`.

Roles and permissions are defined in [permission.go](../../internal/core/domain/permission.go):
- **admin**: Full read/write access. The creator of a workspace is automatically assigned this role (as `admin` in both members repository and `gogate`).
- **member**: Read-only access to workspaces, members, API keys, logs, integrations, templates, settings, and invitations, plus the ability to send test messages (`send.test`).

To bootstrap roles and permissions, make sure you run the database seeds:
1. Run `database/seeds/001_seed_permissions.sql` to populate roles (`admin`, `member`), default permissions, and role-to-permission mappings.
2. Run `database/seeds/002_seed_providers.sql` to seed the provider catalog (memory, mailgun, etc.).
3. Run `database/seeds/003_seed_mailgun_config.sql` to seed Mailgun Portal config field metadata.
4. Run `database/seeds/004_demo_workspace.sql` (optional) — demo user (`demo@weprodev.com` / `secret`), workspace, admin role, memory integration, API key (`demo-client-id` / `demo-secret`).

### Dual Authorization Models

The HTTP server uses two distinct authorization approaches by design:
- **Portal Management API (`/api/v1/workspaces/:wid/...`)**: Restricts access using JWT tokens combined with fine-grained `wpd-gogate` permission middleware (`RequirePermission`).
- **Inbox SSE API (`/api/v1/workspaces/:wid/inbox/...`)**: Restricts access using JWT tokens combined with a direct database workspace membership check (`RequireWorkspaceMember`) and the workspace API key (`RequireWorkspaceAPIKey`). This decoupled model allows client-side SDKs, SSE event streaming, and automation runners to interact with the simulated inbox without requiring full portal RBAC configuration.


### Inbox API (captured messages — REST / Bruno, not Portal UI)

The Portal UI shows **request logs** (`/logs`), not the memory **inbox** (`/inbox/*`). Inbox capture, SSE, and internal ingest are documented in [Portal inbox](./portal-inbox.md).

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
> They are configured via the Portal REST API and stored encrypted in PostgreSQL.

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
```

---

## Message Dispatch & Content Storage

Control how messages are routed and whether bodies are captured in the portal inbox.

Defaults: `message_dispatch_mode=memory`, `store_message_content=false`.

Invalid values for dispatch settings return **400 Bad Request** on PATCH.

Set via Portal **Settings → Message Dispatch** or REST:

```bash
PATCH /api/v1/workspaces/:wid/settings
Authorization: Bearer <portal-jwt>
Content-Type: application/json

{ "message_dispatch_mode": "provider", "store_message_content": "true" }
```

Portal UI: select **Memory** or **Provider**, then check **Store message content in inbox** on or off.

See [Portal inbox](./portal-inbox.md) for routing behavior.

Every send also creates a `message_request_logs` row for operational tracing (correlation via `request_id`).

---

## Troubleshooting

### Mailgun: 403 Forbidden

Using a Mailgun Sandbox Domain:
1. Log into Mailgun → **Sending → Domains → [Sandbox]**
2. Add recipient email to **Authorized Recipients**
3. Recipient must click the verification link

### Provider Not Sending

Check workspace settings via `GET /api/v1/workspaces/:wid/settings`. If `message_dispatch_mode` is `memory`, messages are captured locally — not sent to the provider.

### API Key Rejected

- Ensure `X-Workspace-Key` matches the workspace's `unique_key` (slug, not UUID)
- Verify the API key is active (`GET /api/v1/workspaces/:wid/api-keys` with portal JWT)
- Check `client_id` and `client_secret` — secret is shown once at creation

---

## Related Documentation

- [Architecture](./architecture.md) — System design
- [Portal inbox](./portal-inbox.md) — Memory capture and inbox REST API
- [E2E Testing](./e2e-testing.md) — Automated testing patterns
- [Contributing](./contributing.md) — Adding providers
