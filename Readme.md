<p align="center">
  <img src="assets/logo.png" alt="Message Gateway Logo" width="400" />
</p>

<h1 align="center">Message Gateway</h1>

<p align="center">
  <strong>A unified Go library and HTTP API for sending Email, SMS, Push, and Chat messages.</strong>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/weprodev/wpd-message-gateway"><img src="https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway"><img src="https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway" alt="Go Report Card"></a>
  <a href="https://github.com/weprodev/wpd-message-gateway/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
</p>

<p align="center">
  One interface, multiple providers. Write your messaging code once — switch between Mailgun, Twilio, Firebase, WhatsApp, and more without changing a single line of application code.
</p>

---

## Why Message Gateway?

Building applications that send messages across multiple channels is complex. You need to integrate different APIs, handle various authentication methods, manage provider-specific quirks, and test everything without spamming real users.

**Message Gateway solves this by providing:**

- **Unified API** — Send emails, SMS, push notifications, and chat messages through a single, consistent interface
- **Provider Abstraction** — Switch from Mailgun to SendGrid, or Twilio to Vonage, with a config change—no code modifications
- **DevBox** — Built-in web UI to preview all messages during development
- **E2E Testing** — Memory provider captures messages for automated testing without external dependencies

## Quick Start

### As a Go Package

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

### As an HTTP Server

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
cp configs/local.example.yml configs/local.yml
make start
```

Open **http://localhost:10104** to see the DevBox UI with all captured messages.

## Supported Message Types

| Type | Description | Providers |
|------|-------------|-----------|
| **Email** | Send emails with HTML, plain text, and attachments | Mailgun, Memory |
| **SMS** | Send text messages to mobile phones | Memory (Twilio planned) |
| **Push** | Send notifications to mobile and web apps | Memory (Firebase planned) |
| **Chat** | Send messages to chat platforms | Memory (WhatsApp, Slack planned) |

## Configuration

Configure providers in `configs/local.yml`:

```yaml
providers:
  defaults:
    email: mailgun
    sms: memory
    push: memory
    chat: memory
  
  email:
    mailgun:
      api_key: "your-api-key"
      domain: "mg.yourdomain.com"
```

Or use environment variables:

```bash
MESSAGE_MAILGUN_API_KEY=key-xxxxx
MESSAGE_MAILGUN_DOMAIN=mg.example.com
```

## Development with DevBox

During development, use the **memory** provider to capture all messages locally without sending them to real recipients:

```yaml
providers:
  defaults:
    email: memory
    sms: memory
    push: memory
    chat: memory
```

The DevBox UI at **http://localhost:10104** displays all captured messages organized by type.

### Optional: Mailpit Integration

For realistic HTML email preview with proper rendering:

```bash
make mailpit          # Start Mailpit
# Enable in configs/local.yml: mailpit.enabled: true
make start            # Start the gateway
```

View emails at:
- **DevBox**: http://localhost:10104 (all message types)
- **Mailpit**: http://localhost:10103 (rich email preview)

## E2E Testing

Capture and verify messages your application sends during automated tests—no mocks required.

```yaml
# docker-compose.yml
services:
  gateway:
    image: ghcr.io/weprodev/wpd-message-gateway:latest
    ports:
      - 10101:10101
```

```bash
# Verify the welcome email was sent correctly
curl -s http://localhost:10101/api/v1/emails | \
  jq -e '.emails[0].email.subject == "Welcome!"'
```

See the [E2E Testing Guide](docs/e2e-testing.md) for complete examples.

## Commands

```bash
make install    # Install dependencies
make start      # Start development server with DevBox
make test       # Run tests
make audit      # Format, lint, test, and security check
make build      # Build all packages
make upgrade    # Upgrade all dependencies
```

## Documentation

| Document | Description |
|----------|-------------|
| [Usage Guide](docs/usage.md) | Installation, configuration, and examples |
| [E2E Testing](docs/e2e-testing.md) | Automated testing patterns |
| [Architecture](docs/architecture.md) | System design and principles |
| [DevBox](docs/devbox.md) | Development inbox documentation |
| [Contributing](docs/contributing.md) | How to add providers |
| [Code Conventions](docs/code-conventions.md) | Coding standards |

## Project Structure

```
wpd-message-gateway/
├── cmd/server/          # HTTP server entry point
├── configs/             # YAML configuration files
├── internal/            # Private application code
│   ├── app/             # Configuration and wiring
│   ├── core/            # Business logic and interfaces
│   └── infrastructure/  # Provider implementations
├── pkg/                 # Public packages
│   ├── contracts/       # Message types (Email, SMS, Push, Chat)
│   └── gateway/         # Embedded SDK
├── web/                 # DevBox React UI
└── docs/                # Documentation
```

## Contributing

We welcome contributions! Whether you're fixing bugs, adding features, or improving documentation, your help makes Message Gateway better for everyone.

**Ways to contribute:**

1. **Report bugs** — Open an issue describing the problem
2. **Suggest features** — Share your ideas for improvements
3. **Submit pull requests** — Code contributions are always welcome
4. **Add providers** — Help expand support for more messaging services
5. **Improve docs** — Help others understand and use the project
6. **Become a sponsor** — Support ongoing development and maintenance

See [CONTRIBUTING.md](docs/contributing.md) for detailed guidelines.

### Development Setup

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make install
make start
```

Run `make audit` before submitting pull requests to ensure code quality.

## Become a Sponsor

Message Gateway is open source and free to use. If it helps your team ship faster, consider supporting its continued development.

**Your sponsorship helps:**

- Maintain and improve the codebase
- Add support for more providers 
- Build better documentation and examples
- Keep the project actively maintained

<p align="center">
  <a href="https://github.com/sponsors/weprodev">
    <img src="https://img.shields.io/badge/Sponsor-❤️-ea4aaa?style=for-the-badge" alt="Sponsor">
  </a>
</p>

**Sponsor tiers:**

- **Individual** — Support open source development
- **Startup** — Priority support and feature requests
- **Enterprise** — Custom integrations and dedicated support

[Become a sponsor →](https://github.com/sponsors/weprodev)

## License

[MIT](LICENSE) — Free for personal and commercial use.

---

<p align="center">
  Built with ❤️ by <a href="https://github.com/weprodev">WeProDev</a>
</p>
