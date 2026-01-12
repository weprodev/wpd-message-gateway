# Go Message Gateway

[![Go Reference](https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg)](https://pkg.go.dev/github.com/weprodev/wpd-message-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway)](https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway)

A unified Go gateway for sending messages through various providers.

## Features

- **Provider Agnostic** - Single API for Email, SMS, Push, and Chat
- **Easy Switching** - Change providers via configuration
- **Extensible** - Add custom providers easily
- **Type Safe** - Leverages Go interfaces

## Supported Providers

### 📧 Email
| Provider | Status |
|----------|--------|
| Mailgun | ✅ Ready |
| SendGrid | 📋 Planned |

### 📱 SMS
| Provider | Status |
|----------|--------|
| CM.com | 📋 Planned |
| Twilio | 📋 Planned |


### 🔔 Push
| Provider | Status |
|----------|--------|
| Firebase | 📋 Planned |

### 💬 Chat
| Provider | Status |
|----------|--------|
| WhatsApp | 📋 Planned |
| Telegram | 📋 Planned |

## Quick Start

```bash
go get github.com/weprodev/wpd-message-gateway
```

```bash
export MESSAGE_DEFAULT_EMAIL_PROVIDER=mailgun
export MESSAGE_MAILGUN_API_KEY=your-key
export MESSAGE_MAILGUN_DOMAIN=mg.yourdomain.com
```

```go
cfg, _ := config.LoadFromEnv()
mgr, _ := manager.New(cfg)

mgr.SendEmail(ctx, &contracts.Email{
    To:      []string{"user@example.com"},
    Subject: "Hello!",
    HTML:    "<h1>Welcome!</h1>",
})
```

## Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/architecture.md) | Design, diagrams, and principles |
| [Usage Guide](docs/usage.md) | Installation, configuration, examples |
| [Contributing](docs/contributing.md) | How to add providers, PR process |
| [Code Conventions](docs/code-conventions.md) | Go style, testing, commits |

## Development

```bash
make setup    # Install tools
make test     # Run tests
make audit    # Full quality check (fmt, lint, test, vuln)
```

## Support
- ⭐️ **Star** this repository if you find it useful!
- 🤝 **Contribute** by submitting a Pull Request.
- 💖 **Sponsor** us to support development.

## License

MIT License
