# WPD Message Gateway

[![Go Reference](https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg)](https://pkg.go.dev/github.com/weprodev/wpd-message-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway)](https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway)

**A unified Go library and HTTP API for sending Email, SMS, Push, and Chat messages.**

One interface, multiple providers. Write your messaging code once — switch between Mailgun, Twilio, Firebase, WhatsApp, and more without changing a single line of application code.

## Why Use This?

- **🔌 One API, Many Providers** — Send emails via Mailgun today, switch to SendGrid tomorrow. No code changes.
- **📦 DevBox Included** — Built-in web UI to preview emails, SMS, push notifications, and chat messages during development.
- **🧪 E2E Testing Ready** — Memory provider stores messages in RAM. Query them via REST API for automated testing.
- **🚀 Go Library + HTTP Server** — Use as a Go package (`import`) or deploy as a standalone microservice.

## How It Works

```
┌─────────────────┐
│   Your App      │
└────────┬────────┘
         │ POST /v1/email
         ▼
┌─────────────────┐
│  Gateway Service│
│ (Routes by      │
│  provider name) │
└────────┬────────┘
         │
         │ providers.defaults.email = ?
         │
    ┌────┴────────────────────────────┐
    │                                 │
    ▼                                 ▼
┌─────────────────┐       ┌─────────────────┐
│ "memory"        │       │ "mailgun"       │
│                 │       │ "sendgrid"      │
│ ┌─────────────┐ │       │ etc.            │
│ │ DevBox UI   │ │       │                 │
│ │ (RAM store) │ │       │  Real Provider  │
│ └─────────────┘ │       │  (API calls)    │
│        +        │       │                 │
│ ┌─────────────┐ │       └─────────────────┘
│ │ Mailpit     │ │
│ │ (optional)  │ │ ← Only if mailpit.enabled: true
│ └─────────────┘ │
└─────────────────┘
```

## Message Types

| Type | What it does | Example providers |
|------|--------------|-------------------|
| 📧 **Email** | Send emails with HTML, attachments | Mailgun, SendGrid, AWS SES |
| 📱 **SMS** | Send text messages to phones | Twilio, Vonage |
| 🔔 **Push** | Send notifications to apps | Firebase, OneSignal |
| 💬 **Chat** | Send messages on chat platforms | WhatsApp, Telegram, Slack |

## Quick Start

### Option 1: Use as a Go Package

```bash
go get github.com/weprodev/wpd-message-gateway
```

```go
package main

import (
    "context"
    "log"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/gateway"
)

func main() {
    gw, _ := gateway.New(gateway.Config{
        DefaultEmailProvider: "memory",
    })

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

### Option 2: Run as HTTP Server

```bash
# 1. Clone and configure
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
cp configs/local.example.yml configs/local.yml

# 2. Start everything (Gateway + DevBox UI)
make start
```

Open http://localhost:10104 to see all intercepted messages in the DevBox UI.

→ See [Usage Guide](docs/usage.md) for more examples.

## Configuration

Configure providers in `configs/local.yml`:

```yaml
providers:
  defaults:
    email: mailgun   # or: memory, sendgrid, ses
    sms: memory      # or: twilio, vonage
    push: memory     # or: firebase, onesignal
    chat: memory     # or: slack, telegram, whatsapp
  
  email:
    mailgun:
      api_key: "your-api-key"
      domain: "mg.yourdomain.com"
```

Or use environment variables for secrets:

```bash
MESSAGE_MAILGUN_API_KEY=key-xxxxx
MESSAGE_MAILGUN_DOMAIN=mg.example.com
```

## Development Mode (DevBox)

During development, use the **memory** provider to capture all messages locally:

```yaml
# configs/local.yml
providers:
  defaults:
    email: memory
    sms: memory
    push: memory
    chat: memory
```

→ See [DevBox Documentation](docs/devbox.md) for details.

### Mailpit Integration (Optional)

For realistic email preview with HTML rendering:

```bash
# 1. Start Mailpit
make mailpit

# 2. Enable in configs/local.yml:
mailpit:
  enabled: true

# 3. Start server
make start

# View emails:
#   - DevBox UI: http://localhost:10104 (all message types)
#   - Mailpit:   http://localhost:10103 (email preview)
```

## E2E Testing in CI

Use the gateway to **capture and verify** all messages your app sends during tests. No mocks needed.

**Benefits:**
- ✅ Test real HTTP calls, not mocks
- ✅ Verify exact message content (subject, body, recipients)
- ✅ Test all channels: Email + SMS + Push + Chat
- ✅ Zero external dependencies (no Mailgun/Twilio accounts needed)

```yaml
services:
  gateway:
    image: ghcr.io/weprodev/wpd-message-gateway:latest
    ports:
      - 10101:10101

steps:
  - run: npm test  # Your app sends to http://localhost:10101
  
  - name: Verify welcome email
    run: |
      curl -s http://localhost:10101/api/v1/emails | \
        jq -e '.emails[0].email.subject == "Welcome!"'
```

→ See [E2E Testing Guide](docs/e2e-testing.md) for complete examples.

## Provider Status

| Type | Provider | Status |
|------|----------|--------|
| 📧 Email | Mailgun | ✅ Ready |
| 📧 Email | Memory (DevBox) | ✅ Ready |
| 📧 Email | SendGrid | 📋 Planned |
| 📱 SMS | Memory (DevBox) | ✅ Ready |
| 📱 SMS | Twilio | 📋 Planned |
| 🔔 Push | Memory (DevBox) | ✅ Ready |
| 🔔 Push | Firebase | 📋 Planned |
| 💬 Chat | Memory (DevBox) | ✅ Ready |
| 💬 Chat | WhatsApp | 📋 Planned |

## Commands

```bash
make install    # Install all dependencies
make start      # Start development (Gateway + DevBox UI)
make test       # Run tests
make audit      # Full check: format + lint + test + security
make build      # Build all packages

# Docker
make dev        # Start Gateway via Docker
make dev-down   # Stop Docker

# Optional (email preview)
make mailpit    # Start Mailpit for HTML email preview
```

## Project Structure

```
wpd-message-gateway/
├── cmd/server/          # HTTP server entry point
├── configs/             # YAML configuration files
├── internal/            # Private application code
│   ├── app/             # Configuration, wiring, validation
│   ├── core/            # Business logic
│   │   ├── port/        # Interface definitions (contracts)
│   │   └── service/     # Gateway service, registry
│   ├── infrastructure/  # External integrations
│   │   └── provider/    # Provider implementations
│   │       ├── mailgun/ # Mailgun email provider
│   │       └── memory/  # In-memory provider (DevBox)
│   └── presentation/    # HTTP layer
│       ├── handler/     # Request handlers
│       └── router.go    # Route definitions
├── pkg/                 # Public packages for external use
│   ├── contracts/       # Message types (Email, SMS, Push, Chat)
│   ├── errors/          # Error types
│   └── gateway/         # Embedded SDK
├── web/                 # DevBox React UI
└── tests/bruno/         # API test collection
```

## Documentation

| Document | Description |
|----------|-------------|
| [Usage Guide](docs/usage.md) | Install, configure, and send messages |
| [E2E Testing](docs/e2e-testing.md) | Test your app's messages in CI |
| [Architecture](docs/architecture.md) | System design and principles |
| [DevBox](docs/devbox.md) | Development inbox UI |
| [Contributing](docs/contributing.md) | Add new providers |
| [Workflow](docs/workflow.md) | CI/CD and releases |
| [Code Conventions](docs/code-conventions.md) | Coding standards |

## License

[MIT](LICENSE)
