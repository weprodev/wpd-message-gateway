# Contributing Guide

## Quick Start

### Option A: Docker Compose (Recommended)
This spins up the entire stack including a PostgreSQL database (preloaded with migrations and permission/demo seeds), a hot-reloading Go backend, and the Portal UI dev server:

```bash
git clone git@github.com:weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make dev
```

### Option B: Native Setup (No Docker)
If you prefer to run components directly on your host machine:

```bash
git clone git@github.com:weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make install
make start
```

---

## Adding a New Provider

Adding a provider **does not require modifying any existing code**. This follows the Open/Closed Principle: open for extension, closed for modification.

The system uses a **self-registration pattern** via Go's `init()` mechanism. When you create a new provider package with a `register.go` file and add a blank import in `internal/app/imports.go`, the provider becomes available in both:

- The **embedded SDK** (`gateway.New(...)` — standalone Go usage)
- The **HTTP server** (via UI-managed integrations in the database)

### Step 1: Create the Provider Package

```bash
mkdir -p pkg/provider/sendgrid
```

```go
// pkg/provider/sendgrid/sendgrid.go
package sendgrid

import (
    "context"
    "fmt"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

const ProviderName = "sendgrid"

// Provider implements contracts.EmailSender for SendGrid.
type Provider struct {
    apiKey    string
    fromEmail string
    fromName  string
}

// compile-time interface check
var _ contracts.EmailSender = (*Provider)(nil)

func New(cfg Config) (*Provider, error) {
    if cfg.APIKey == "" {
        return nil, fmt.Errorf("sendgrid: API key required")
    }
    return &Provider{
        apiKey:    cfg.APIKey,
        fromEmail: cfg.FromEmail,
        fromName:  cfg.FromName,
    }, nil
}

func (p *Provider) Send(ctx context.Context, email contracts.Email) (*contracts.SendResult, error) {
    // Call SendGrid HTTP API here
    return &contracts.SendResult{
        ID:      "sg-msg-id",
        Message: "queued",
    }, nil
}

// Config holds SendGrid credentials.
type Config struct {
    APIKey    string
    FromEmail string
    FromName  string
}
```

### Step 2: Self-Register via `init()`

```go
// pkg/provider/sendgrid/register.go
package sendgrid

import (
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/registry"
)

func init() {
    registry.RegisterEmailProvider("sendgrid", func(cfg registry.EmailConfig) (contracts.EmailSender, error) {
        return New(Config{
            APIKey:    cfg.APIKey,
            FromEmail: cfg.FromEmail,
            FromName:  cfg.FromName,
        })
    })
}
```

### Step 3: Add a Blank Import

```go
// internal/app/imports.go
package app

import (
    // Provider self-registration — order does not matter
    _ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
    _ "github.com/weprodev/wpd-message-gateway/pkg/provider/memory"
    _ "github.com/weprodev/wpd-message-gateway/pkg/provider/sendgrid" // ← add this
)
```

That's all that's needed for both SDK and server modes.

### Step 4: Use in SDK Mode (Go package)

```go
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "sendgrid",
    EmailProviders: map[string]gateway.EmailConfig{
        "sendgrid": {
            CommonConfig: gateway.CommonConfig{APIKey: "SG.xxx"},
            FromEmail:    "noreply@example.com",
            FromName:     "MyApp",
        },
    },
})
gw.SendEmail(ctx, email)
```

### Step 5: Use in server mode (REST API)

1. Start the server and obtain a portal JWT (`POST /api/v1/auth/login`)
2. `POST /api/v1/workspaces/:wid/integrations` with channel **email**, provider **sendgrid**, and JSON config (`api_key`, `from_email`, `from_name`) — or use the Portal **Integrations** page
3. The `GatewayService` reads this from DB and instantiates your provider via the registry factory

> Provider config in server mode is stored **AES-encrypted** in the `integrations` table — configure via Portal UI or REST, not YAML files.

### Step 6: Add Tests

```go
// pkg/provider/sendgrid/sendgrid_test.go
package sendgrid

import (
    "context"
    "testing"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestNew(t *testing.T) {
    tests := []struct {
        name    string
        cfg     Config
        wantErr bool
    }{
        {"valid", Config{APIKey: "key", FromEmail: "a@b.com"}, false},
        {"missing key", Config{}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := New(tt.cfg)
            if (err != nil) != tt.wantErr {
                t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestSend(t *testing.T) {
    p, err := New(Config{APIKey: "key", FromEmail: "noreply@example.com"})
    if err != nil {
        t.Fatal(err)
    }
    result, err := p.Send(context.Background(), &contracts.Email{
        To:      []string{"user@example.com"},
        Subject: "Test",
        HTML:    "<p>Hello</p>",
    })
    if err != nil {
        t.Errorf("Send() error = %v", err)
    }
    if result == nil || result.ID == "" {
        t.Error("expected non-empty result ID")
    }
}
```

### Step 7: Quality Check

```bash
make audit   # format, lint, test, security scan
```

---

## Provider Directory Layout

```
pkg/provider/
├── mailgun/
│   ├── mailgun.go      # Provider implementation
│   └── register.go     # init() self-registration
│
├── memory/             # Built-in dev/test capture provider
│   ├── store.go        # Shared in-process message store
│   ├── email.go        # Memory email provider
│   ├── sms.go          # Memory SMS provider
│   ├── push.go         # Memory push provider
│   ├── chat.go         # Memory chat provider
│   └── register.go     # init() self-registration
│
└── sendgrid/           # Your new provider
    ├── sendgrid.go
    ├── sendgrid_test.go
    └── register.go     ← KEY FILE: always required
```

---

## Self-Registration Pattern

```text
┌─────────────────────┐   init()   ┌────────────────────────────┐
│  Provider Package   │──────────▶ │  Provider Registry         │
│  (register.go)      │            │  (pkg/registry)            │
│                     │            │                            │
│  emailFactories[    │            │  emailFactories[           │
│    "sendgrid"       │            │    "sendgrid"              │
│  ] = factory func   │            │  ] = factory func          │
└─────────────────────┘            └────────────────────────────┘
                                           │
                              ┌────────────┴────────────┐
                              │                         │
                       SDK gateway.New()         Server GatewayService
                       (static config)           (reads from DB integrations)
```

---

## Adding Other Channel Types (SMS, Push, Chat)

Same pattern — use the corresponding registry function and port interface:

| Channel | Registry function | Contract interface |
|---------|------------------|----------------|
| Email | `registry.RegisterEmailProvider(name, factory)` | `contracts.EmailSender` |
| SMS | `registry.RegisterSMSProvider(name, factory)` | `contracts.SMSSender` |
| Push | `registry.RegisterPushProvider(name, factory)` | `contracts.PushSender` |
| Chat | `registry.RegisterChatProvider(name, factory)` | `contracts.ChatSender` |

---

## Pull Request Checklist

- [ ] **`/smell develop`** — no BLOCKER/HIGH ([verification.md](../agents/verification.md))
- [ ] **`make audit`** passes (format, lint, test, security)
- [ ] Tests added with meaningful coverage
- [ ] Interface check: `var _ contracts.EmailSender = (*Provider)(nil)`
- [ ] `register.go` with `init()` for self-registration
- [ ] Blank import added to `internal/app/imports.go`
- [ ] No provider credentials in YAML — use environment variables (embedded SDK) or Portal REST API (server mode)
- [ ] Commit messages follow [Conventional Commits](../workflow.md#commit-conventions)

---

## Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

| Message | Version bump |
|---------|-------------|
| `feat(sendgrid): add email provider` | Minor |
| `fix(mailgun): handle rate limit 429` | Patch |
| `docs: update contributing guide` | Patch |
| `feat!: change auth to email+PIN` | **Major** |

---

## Related Documentation

- [Architecture](./architecture.md) — System design and layers
- [Usage](./usage.md) — API reference
- [Code Conventions](./code-conventions.md) — Coding standards
- [Workflow Guide](../workflow.md) — CI/CD, commits, releases
