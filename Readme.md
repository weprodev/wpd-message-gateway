# WPD Message Gateway

[![Go Reference](https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg)](https://pkg.go.dev/github.com/weprodev/wpd-message-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway)](https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway)

A unified Go package for sending messages through multiple providers. One API, any provider.

## Installation

```bash
go get github.com/weprodev/wpd-message-gateway
```

## Quick Start

**1. Configure environment:**
```bash
export MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun
export MESSAGE_MAILGUN_API_KEY=your-api-key
export MESSAGE_MAILGUN_DOMAIN=mg.yourdomain.com
export MESSAGE_MAILGUN_FROM_EMAIL=noreply@yourdomain.com
```

**2. Send an email:**
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
        HTML:    "<h1>Hello from Go!</h1>",
    })
}
```

## Providers

| Type | Provider | Status |
|------|----------|--------|
| 📧 Email | Mailgun | ✅ Ready |
| 📧 Email | SendGrid | 📋 Planned |
| 📧 Email | AWS SES | 📋 Planned |
| 📱 SMS | Twilio | 📋 Planned |
| 🔔 Push | Firebase | 📋 Planned |
| 💬 Chat | WhatsApp | 📋 Planned |

## Documentation

- **[Usage Guide](docs/usage.md)** — Installation, configuration, examples
- **[Architecture](docs/architecture.md)** — Design patterns and principles
- **[Contributing](docs/contributing.md)** — How to add new providers
- **[Code Conventions](docs/code-conventions.md)** — Style guide

## Development

```bash
make test       # Run tests
make lint       # Run linter
make sandbox    # Interactive testing CLI
```

## License

MIT
