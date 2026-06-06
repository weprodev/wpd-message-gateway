<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="frontend/public/logo-dark-mode.svg">
    <img src="frontend/public/logo-light-mode.svg" alt="Message Gateway" width="400" />
  </picture>
</p>

<p align="center">
  <!-- Build & CI -->
  <a href="https://github.com/weprodev/wpd-message-gateway/actions/workflows/ci.yml"><img src="https://github.com/weprodev/wpd-message-gateway/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
  <!-- Quality & Go -->
  <a href="https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway"><img src="https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway" alt="Go Report Card"></a>
  <a href="https://pkg.go.dev/github.com/weprodev/wpd-message-gateway"><img src="https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg" alt="Go Reference"></a>
  <!-- License -->
  <a href="https://github.com/weprodev/wpd-message-gateway/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
</p>

<p align="center">
  <strong>A unified messaging solution for Email, SMS, Push, and Chat.</strong>
</p>

This project offers two distinct ways to integrate messaging into your architecture:
1. **Embedded Go SDK:** A lightweight library for your Go applications. Send messages natively without running a separate server or database.
2. **Standalone HTTP Gateway:** A fully-featured service deployable via Docker. Provides a unified REST API for any programming language, backed by PostgreSQL and a React UI.

Write your messaging code once — switch between Mailgun, Twilio, Firebase, WhatsApp, and more without changing a single line of application code.

---

## Key Features

- **Unified Interface:** Send Email, SMS, Push, and Chat messages through one consistent API.
- **Provider Abstraction:** Change a configuration value to switch providers, not your code.
- **DB-First Configuration:** Provider credentials and workspace settings in PostgreSQL (Portal REST API; UI pages partial).
- **Workspace Isolation:** Multi-tenant support for separate providers, API keys, and templates per workspace.
- **Developer-Friendly:** Includes a local Memory provider to capture messages locally and assert real payloads without mocking.

---

## Two Ways to Use

```mermaid
graph TD
    classDef main fill:#f9f9f9,stroke:#333,stroke-width:2px;
    classDef feature fill:#e1f5fe,stroke:#0288d1,stroke-width:1px;
    classDef embedded fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef standalone fill:#fff3e0,stroke:#e65100,stroke-width:2px;

    %% SDK Usage Path
    subgraph SDK ["Go Package (Embedded SDK)"]
        direction TB
        A["💻 go get wpd-message-gateway"] --> B["📦 gateway.New(config)"]
        B -.-> C["✓ No server needed"]:::feature
        B -.-> D["✓ No database needed"]:::feature
        B -.-> E["✓ Config in your code"]:::feature
    end
    class SDK embedded;

    %% Server Usage Path
    subgraph Server ["HTTP Server (Standalone Gateway)"]
        direction TB
        F["🐳 make dev (Docker)"] --> G["🌐 POST /v1/email (HTTP)"]
        G -.-> H["✓ Any language (Python, JS...)"]:::feature
        G -.-> I["✓ React Portal UI included"]:::feature
        G -.-> J["✓ PostgreSQL stores secrets"]:::feature
    end
    class Server standalone;
```

---

## Quick Start

### Option A: Embedded Go SDK

Install the SDK:

```bash
go get github.com/weprodev/wpd-message-gateway
```

Send an email in 10 lines:

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

No server. No database. Just `go get` and send.

---

### Option B: Standalone HTTP Gateway

The easiest way to get the standalone server running is with Docker Compose.

**Prerequisites:** Docker and Docker Compose

```bash
# 1. Clone the repository
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway

# 2. Start the gateway, PostgreSQL, and Portal UI
make dev
```

Once running:

1. Open **http://localhost:10104** — the Portal UI
2. Create an account and sign in (or use the pre-seeded demo user: `demo@example.com` / `password`)
3. Select the pre-seeded workspace or bootstrap a workspace and API key using the REST API (the Portal UI currently covers authentication, workspaces, message logs, and send tests; other provisioning is REST-only). See [E2E Bootstrap Guide](docs/backend/e2e-testing.md#2-bootstrap-create-workspace-and-api-key) for details.
4. Send a message from any language:

```bash
curl -X POST http://localhost:10101/v1/email \
  -u "wk_client_id:your_secret" \
  -H "X-Workspace-Key: myapp" \
  -H "Content-Type: application/json" \
  -d '{"to":["user@example.com"],"subject":"Hello","html":"<h1>World</h1>"}'
```

<details>
<summary><strong>Alternative: Manual Setup (No Docker)</strong></summary>

If you prefer to run components directly on your machine:
**Prerequisites:** Go 1.22+, PostgreSQL, Node.js 20+

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
cp configs/local.example.yml configs/local.yml   # edit DB credentials here

make install
make start
```
</details>

---

## Why Message Gateway?

| Problem | Solution |
|---------|----------|
| Each provider has a different API | **Unified interface** — Email, SMS, Push, Chat through one consistent API |
| Switching providers means rewriting code | **Provider abstraction** — change a config value, not your code |
| Credentials scattered in env vars | **DB-first config** — all secrets in PostgreSQL, managed via Portal UI |
| Hard to test messaging in CI | **Memory provider** — captures messages locally, assert real payloads without mocking |
| Multi-tenant apps need isolation | **Workspace isolation** — separate providers, API keys, templates per workspace |

---

## Supported Providers

| Type | Providers | Status |
|------|----------|--------|
| **Email** | Mailgun, Memory | ✅ Available |
| **SMS** | Memory | ✅ Memory · 🔜 Twilio |
| **Push** | Memory | ✅ Memory · 🔜 Firebase |
| **Chat** | Memory | ✅ Memory · 🔜 Slack, WhatsApp, Telegram |

> **Adding a provider?** See [Contributing: Adding Providers](docs/backend/contributing.md).

---

## Portal UI

Available at **http://localhost:10104** when the server runs.

| Feature | Description |
|---------|-------------|
| **Sign in / Register** | Portal account (JWT for `/api/v1/*`) |
| **Workspaces** | List and open a workspace |
| **Message logs** | Audit trail of gateway send requests |
| **Send test** | Test send per channel from the dashboard |

Memory **inbox** capture, integrations, API keys, templates, members, and settings are **REST/Bruno only** today — see [Usage guide](docs/backend/usage.md) and [Portal inbox](docs/backend/portal-inbox.md).

---

## Configuration

The gateway enforces a strict separation between **server infrastructure** and **messaging credentials**:

1. **Server Configuration (`configs/local.yml`):**
   This file handles the operational infrastructure. Copy `configs/local.example.yml` to `configs/local.yml`.
   - `DISPATCH_MODE`: Defines if messages are sent for real (`provider_only`), captured (`memory_only`), or both.
   - Database credentials (`DB_HOST`, `DB_USER`, etc.) for PostgreSQL.
   - `JWT_SECRET` for Portal UI session authentication.
   *(Note: Never place provider credentials like Mailgun API keys in this file.)*

2. **Provider credentials (server mode):**
   Stored AES-encrypted in PostgreSQL. Configure via the **Portal REST API** (Integrations endpoints) — there is no Portal UI page for this yet. See [Usage guide](docs/backend/usage.md) and [E2E bootstrap](docs/backend/e2e-testing.md).

## E2E Testing

Use `memory_only` dispatch mode in CI/CD to capture messages without sending real emails or SMS. Then assert the payloads via the Portal Inbox API.

No mocking. Real HTTP calls. Full payload assertions.

```yaml
# docker-compose.test.yml — drop in with your CI
services:
  message-gateway:
    image: weprodev/wpd-message-gateway
    environment:
      DISPATCH_MODE: memory_only
```

📖 Full setup guide: [E2E Testing with Docker & GitHub Actions](docs/backend/e2e-testing.md)

---

## Development Environment & Workflow

There are two primary ways to set up the development environment:

1. **Local Setup:** Using `make install` and `make start` with local dependencies (Go, Node, PostgreSQL).
2. **Docker Setup:** Using `make dev` as the recommended one-command method to run everything via Docker Compose.

| Command | Description |
|---------|------------|
| `make install` | Install Go + frontend dependencies |
| `make start` | Start Gateway API + Portal UI |
| `make dev` | Run via Docker Compose (includes PostgreSQL) |
| `make test` | Run Go tests |
| `make audit` | Full quality gate: fmt + lint + test (Go + frontend) + govulncheck + build |
| `make build` | Compile Go binary + build frontend (no tests) |
| `make ui` | Portal UI only (Vite dev server, port 10104) |
| `make storybook` | Component library (Storybook, port 6006) |
| `make upgrade` | Upgrade all dependencies |

This project uses a spec-driven, AI-assisted development workflow. See the [Development Flow](docs/development-flow.md) documentation for details on Spec Kit commands and Agent Mapping.

---

## Project Structure

```
wpd-message-gateway/
├── cmd/server/          # HTTP server entry point
├── configs/             # Server config (ports, JWT) — NOT provider credentials
├── database/
│   ├── migrations/      # SQL schema migrations
│   └── seeds/           # Optional demo data
├── internal/
│   ├── app/             # Config, wire, validation, provider registration
│   ├── core/            # Domain models, services, ports (interfaces)
│   ├── infrastructure/  # Providers (mailgun, memory…) + Postgres repos + logger
│   ├── presentation/    # HTTP router, handlers, middleware
│   └── registry/        # Provider factory registry (self-registration via init)
├── pkg/                 # Public Go packages
│   ├── contracts/       # Email, SMS, Push, Chat message types
│   └── gateway/         # Embedded SDK — gateway.New()
├── frontend/            # React Portal UI (Vite + TypeScript + Tailwind)
├── tests/bruno/         # HTTP API test collection (Bruno)
├── specs/               # Feature specifications (Spec Kit)
└── docs/                # Documentation
```

---

## Documentation

| Document | Description |
|----------|------------|
| [Docs Hub](docs/README.md) | Index of all backend + frontend documentation |
| [Usage Guide](docs/backend/usage.md) | SDK and HTTP API reference, authentication, multi-language examples |
| [Architecture](docs/backend/architecture.md) | System design, two modes of operation, DB schema |
| [Portal Inbox](docs/backend/portal-inbox.md) | Message inbox, dispatch modes, inbox API |
| [E2E Testing](docs/backend/e2e-testing.md) | CI/CD integration, capturing and asserting messages |
| [Contributing](docs/backend/contributing.md) | Adding new providers, code quality gates |
| [Code Conventions](docs/backend/code-conventions.md) | Go coding standards |
| [Frontend Docs](docs/frontend/README.md) | Portal UI — Vite, TypeScript, shadcn/ui, conventions |
| [Frontend Engineer](docs/frontend/frontend-engineer.md) | Principal-style workflow, architecture, Storybook |
| [Backend Engineer](docs/backend/backend-engineer.md) | Go layers, registry, security, quality gate |
| [Workflow](docs/workflow.md) | CI/CD and release process |
| [Development Flow](docs/development-flow.md) | Spec Kit workflow and AI-Assisted Development |
| [Bruno Collections](tests/bruno/) | HTTP API tests |

---

## Contributing

1. **Report bugs** — [Open an issue](https://github.com/weprodev/wpd-message-gateway/issues)
2. **Suggest features** — Ideas and discussions welcome
3. **Pull requests** — Code, docs, and tests
4. **Add providers** — See [Contributing Guide](docs/backend/contributing.md)
5. **Sponsor** — Support ongoing development

> Run **`/smell develop`** then **`make audit`** before opening a PR — see [docs/agents/verification.md](docs/agents/verification.md).

---

## Licensing & Sponsorship

Released under the **[MIT License](LICENSE)**.

| Who | What we ask |
|-----|------------|
| Individuals, learning, side projects | Use freely. Sponsorship optional. |
| Small teams | Use freely. Consider sponsoring if it's central to your stack. |
| Mid-size companies & enterprises | **Please [sponsor](https://github.com/sponsors/weprodev)** |
| Qualifying non-profits | Use freely. Sponsorship optional. |

<p align="center">
  <a href="https://github.com/sponsors/weprodev">
    <img src="https://img.shields.io/badge/Sponsor-❤️-ea4aaa?style=for-the-badge" alt="Sponsor on GitHub">
  </a>
</p>

---

<p align="center">
  Built with ❤️ by <a href="https://github.com/weprodev">WeProDev</a>
</p>
