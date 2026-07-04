# WPD Message Gateway SDK (`pkg/`)

The `pkg/` tree is the **embedded Go SDK** — fully independent from `internal/`, PostgreSQL, and the HTTP portal. Third-party apps import it to send email, SMS, push, and chat through provider plugins without running the gateway server.

**Iron rule:** `pkg/*` must **never** import `internal/*`. Enforced by `pkg/import_boundary_test.go`.

## Packages

| Package | Role |
| ------- | ---- |
| `contracts/` | Channel message types and sender interfaces (`EmailSender`, `SMSSender`, …) |
| `registry/` | Provider factory registry (register + resolve by name) |
| `provider/*` | Concrete providers (`mailgun`, `memory`, …) |
| `gateway/` | High-level `gateway.New(Config)` facade for in-process use |

## Mode 1 — Embedded SDK (no server, no DB)

```go
import (
    "context"

    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/gateway"
    _ "github.com/weprodev/wpd-message-gateway/pkg/provider/memory" // registers "memory"
)

func main() {
    gw, err := gateway.New(gateway.Config{
        DefaultEmailProvider: "memory",
        EmailProviders: map[string]gateway.EmailConfig{
            "memory": {},
        },
    })
    if err != nil {
        panic(err)
    }

    result, err := gw.SendEmail(context.Background(), &contracts.Email{
        To:      []string{"user@example.com"},
        Subject: "Hello",
        HTML:    "<p>Hi</p>",
    })
    if err != nil {
        panic(err)
    }
    _ = result.ID
}
```

Import only the providers you need via blank import, or use `pkg/gateway` which bundles common providers through `gateway/imports.go`.

## Mode 2 — Dynamic resolution via registry

For config loaded from your own store (not this repo's PostgreSQL):

```go
import (
    "encoding/json"

    "github.com/weprodev/wpd-message-gateway/pkg/registry"
    _ "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
)

func main() {
    raw := map[string]any{
        "api_key":    "key-xxx",
        "domain":     "mg.example.com",
        "from_email": "noreply@example.com",
    }
    b, _ := json.Marshal(raw)

    factory, err := registry.GetEmailFactory("mailgun")
    if err != nil {
        panic(err)
    }

    var cfg registry.EmailConfig
    if err := json.Unmarshal(b, &cfg); err != nil {
        panic(err)
    }

    sender, err := factory(cfg)
    if err != nil {
        panic(err)
    }
    _ = sender
}
```

## Direct provider usage

```go
import (
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/provider/mailgun"
)

sender, err := mailgun.New(mailgun.Config{
    Domain:    "mg.example.com",
    APIKey:    "key-xxx",
    FromEmail: "noreply@example.com",
})
```

## Adding a provider

1. Create `pkg/provider/<name>/`
2. Implement the relevant `contracts.*Sender` interface
3. Register in `init()` via `registry.RegisterEmailProvider` (or SMS/push/chat)
4. Add tests colocated in the provider package

## Relationship to HTTP server

The HTTP server (`cmd/server`, `internal/`) uses the **same** `pkg/contracts` and `pkg/registry` but loads integration config from PostgreSQL. Portal UI and `/v1/*` routes are **not** part of `pkg/`.
