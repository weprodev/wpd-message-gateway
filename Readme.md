# WPD Message Gateway

[![Go Reference](https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg)](https://pkg.go.dev/github.com/weprodev/wpd-message-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway)](https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway)

**A unified Go library and HTTP API for sending Email, SMS, Push, and Chat messages.**

One interface, multiple providers. Write your messaging code once — switch between Mailgun, Twilio, Firebase, WhatsApp, and more without changing a single line of application code.

## Why Use This?

- **🔌 One API, Many Providers** — Send emails via Mailgun today, switch to SendGrid tomorrow. No code changes.
- **📦 DevBox Included** — Built-in web UI to preview emails, SMS, push notifications, and chat messages during development. No real messages sent.
- **🧪 E2E Testing Ready** — Memory provider stores messages in RAM. Query them via REST API for automated testing.
- **🚀 Go Library + HTTP Server** — Use as a Go package (`import`) or deploy as a standalone microservice for any language.

## What is this?

Think of it as a **universal adapter for messaging**. Instead of learning how Mailgun, Twilio, Firebase, and WhatsApp each work differently, you use **one simple interface**:

## How It Works

```
┌─────────────────┐
│   Your App      │
└────────┬────────┘
         │ POST /v1/email
         ▼
┌─────────────────┐
│    Manager      │
│ (Routes by      │
│  provider name) │
└────────┬────────┘
         │
         │ MESSAGE_DEFAULT_EMAIL_PROVIDER = ?
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
│ │ (optional)  │ │ ← Only if MAILPIT_ENABLED=true
│ └─────────────┘ │
└─────────────────┘
```

**Development** (`MESSAGE_DEFAULT_EMAIL_PROVIDER=memory`):
- Emails stored in RAM → View in DevBox UI
- Optionally forward to Mailpit (`MAILPIT_ENABLED=true`)

**Production** (`MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun`):
- Emails sent via real provider API
- Nothing in DevBox

```go
// Send an email - same code works with any email provider
mgr.SendEmail(ctx, &contracts.Email{
    To:      []string{"user@example.com"},
    Subject: "Hello!",
    HTML:    "<h1>Welcome!</h1>",
})
```

## Message Types

| Type | What it does | Example providers |
|------|--------------|-------------------|
| 📧 **Email** | Send emails with HTML, attachments | Mailgun, SendGrid, AWS SES |
| 📱 **SMS** | Send text messages to phones | Twilio, Vonage |
| 🔔 **Push** | Send notifications to apps | Firebase, OneSignal |
| 💬 **Chat** | Send messages on chat platforms | WhatsApp, Telegram |

## Quick Start

### 1. Install

```bash
go get github.com/weprodev/wpd-message-gateway
```

### 2. Configure

Configure your providers in `configs/local.yml`:

```yaml
# configs/local.yml
providers:
  defaults:
    email: mailgun
  email:
    mailgun:
      api_key: "your-api-key"
      domain: "mg.yourdomain.com"
```

Or use environment variables for secrets:

### 3. Send

```go
package main

import (
    "context"
    "github.com/weprodev/wpd-message-gateway/config"
    "github.com/weprodev/wpd-message-gateway/contracts"
    "github.com/weprodev/wpd-message-gateway/manager"
)

func main() {
    cfg, _ := config.LoadFromEnv()
    mgr, _ := manager.New(cfg)

    mgr.SendEmail(context.Background(), &contracts.Email{
        To:      []string{"user@example.com"},
        Subject: "Welcome!",
        HTML:    "<h1>Hello!</h1>",
    })
}
```

That's it! See [Usage Guide](docs/usage.md) for more examples.

## Development Mode (DevBox)

During development, you don't want to send real messages. The **DevBox** catches all messages and shows them in a web UI:

```bash
# 1. Copy config example
cp configs/local.example.yml configs/local.yml

# 2. Start everything (Gateway + DevBox UI)
make start
```

Open http://localhost:10104 to see all intercepted messages.

→ See [DevBox Documentation](docs/devbox.md) for details.

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
make clean      # Clean artifacts

# Docker
make dev        # Start Gateway via Docker
make dev-down   # Stop Docker

# Optional (SMTP provider testing only)
make mailpit    # Start Mailpit
```

### When do I need Mailpit?

**Most developers don't need it.** The DevBox UI shows all messages stored in memory.

Use Mailpit when you want **realistic email preview** (HTML rendering, attachments):

```bash
# 1. Start Mailpit
make mailpit

# 2. Set in configs/local.yml:
providers:
  defaults:
    email: memory
mailpit:
  enabled: true

# 3. Send emails → View in BOTH:
#    - DevBox UI: http://localhost:10104 (all message types)
#    - Mailpit:   http://localhost:10103 (email preview)
```

## Documentation

| Document | Description |
|----------|-------------|
| [Usage Guide](docs/usage.md) | How to install, configure, and send messages |
| [Architecture](docs/architecture.md) | How the package is designed |
| [DevBox](docs/devbox.md) | Development inbox for testing |
| [Contributing](docs/contributing.md) | How to add new providers |
| [Workflow](docs/workflow.md) | CI/CD, commit conventions, and releases |
| [Code Conventions](docs/code-conventions.md) | Coding style guide |

## Project Structure

```
wpd-message-gateway/
├── config/         # Configuration loading
├── contracts/      # Message types (Email, SMS, Push, Chat)
├── manager/        # Main API you use
├── providers/      # Provider implementations
│   ├── email/      # Mailgun, SendGrid, etc.
│   ├── sms/        # Twilio, etc.
│   ├── push/       # Firebase, etc.
│   └── chat/       # WhatsApp, Telegram, etc.
├── internal/       # Internal packages (DevBox API)
└── web/            # DevBox React UI
```

## License

[MIT](LICENSE)
