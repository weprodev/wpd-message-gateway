<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="frontend/public/logo-dark-mode.svg">
    <img src="frontend/public/logo-light-mode.svg" alt="Message Gateway" width="400" />
  </picture>
</p>

<p align="center">
  <a href="https://github.com/weprodev/wpd-message-gateway/actions/workflows/ci.yml"><img src="https://github.com/weprodev/wpd-message-gateway/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://pkg.go.dev/github.com/weprodev/wpd-message-gateway"><img src="https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg" alt="Go Reference"></a>
  <a href="https://github.com/weprodev/wpd-message-gateway/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
</p>

Unified messaging for **Email, SMS, Push, and Chat** — one API, swappable providers.

Use it as an **embedded Go SDK** (`pkg/gateway`, no server) or as a **standalone HTTP gateway** with PostgreSQL and a React portal. Details: [docs/](docs/README.md).

## Two ways to use

```mermaid
graph TD
    classDef main fill:#f9f9f9,stroke:#333,stroke-width:2px;
    classDef feature fill:#e1f5fe,stroke:#0288d1,stroke-width:1px;
    classDef embedded fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;
    classDef standalone fill:#fff3e0,stroke:#e65100,stroke-width:2px;

    subgraph SDK ["Go Package (Embedded SDK)"]
        direction TB
        A["go get wpd-message-gateway"] --> B["gateway.New(config)"]
        B -.-> C["No server"]:::feature
        B -.-> D["No database"]:::feature
        B -.-> E["Config in code"]:::feature
    end
    class SDK embedded;

    subgraph Server ["HTTP Server (Standalone Gateway)"]
        direction TB
        F["make dev (Docker)"] --> G["POST /v1/email"]
        G -.-> H["Any language"]:::feature
        G -.-> I["React Portal UI"]:::feature
        G -.-> J["PostgreSQL"]:::feature
    end
    class Server standalone;
```

## Architecture

```mermaid
graph LR
    classDef client fill:#e1f5fe,stroke:#0288d1,stroke-width:1px;
    classDef main fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef db fill:#e8f5e9,stroke:#2e7d32,stroke-width:1px;
    classDef ext fill:#f5f5f5,stroke:#9e9e9e,stroke-width:1px;

    subgraph ClientLayer ["Clients"]
        UI["React Portal UI"]:::client
        AppClient["API clients"]:::client
    end

    subgraph ServerLayer ["Go HTTP Gateway"]
        API["REST API"]:::main
        Svc["Gateway & Portal services"]:::main
        Reg["Provider registry"]:::main
    end

    subgraph DBLayer ["Database"]
        DB[(PostgreSQL)]:::db
    end

    subgraph ExternalServices ["Providers"]
        EmailProvider["Email"]:::ext
        SMSProvider["SMS"]:::ext
        ChatProvider["Chat"]:::ext
    end

    UI -->|"/api/v1/* JWT"| API
    AppClient -->|"/v1/* API key"| API
    API --> Svc
    Svc --> DB
    Svc --> Reg
    Reg --> EmailProvider
    Reg --> SMSProvider
    Reg --> ChatProvider
```

---

## Quick start (Docker)

**Requires:** Docker

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make dev
```

| Service | URL |
|---------|-----|
| Portal UI | http://localhost:10104 |
| Gateway API | http://localhost:10101 |

### Demo account

Seeded on first DB init. Reset anytime: `make seed-demo`.

| | |
|--|--|
| Portal | `demo@weprodev.com` / `secret` |
| API key | `demo-client-id` / `demo-secret` |
| Workspace ID | `00000000-0000-0000-0000-000000000001` |

```bash
curl -X POST http://localhost:10101/v1/email \
  -u "demo-client-id:demo-secret" \
  -H "X-Workspace-Key: 00000000-0000-0000-0000-000000000001" \
  -H "Content-Type: application/json" \
  -d '{"to":["user@example.com"],"subject":"Hello","html":"<h1>World</h1>"}'
```

---

## Go SDK

```bash
go get github.com/weprodev/wpd-message-gateway
```

```go
gw, _ := gateway.New(gateway.Config{ /* providers in code */ })
gw.SendEmail(ctx, &contracts.Email{To: []string{"a@b.com"}, Subject: "Hi", HTML: "<p>Hi</p>"})
```

No database, no HTTP server. Full examples: [usage guide](docs/backend/usage.md).

---

## Local development (no Docker)

Go 1.26+, PostgreSQL 17+, Node 20+.

```bash
cp configs/local.example.yml configs/local.yml   # set DB + secrets
make install && make start
```

Migrations and seeds: [usage guide](docs/backend/usage.md).

| Command | |
|---------|--|
| `make dev` / `make dev-down` | Docker up / down |
| `make seed-demo` | Re-apply demo seeds |
| `make audit` | Lint, test, build (pre-PR gate) |

Run `make help` for the full list.

---

## Documentation

Everything else lives under **[docs/](docs/README.md)** — API reference, architecture, E2E testing, contributing.

---

## Contributing

Issues and PRs welcome. Before opening a PR: `make audit`. See [contributing](docs/backend/contributing.md).

MIT License — [LICENSE](LICENSE).

<p align="center">
  <a href="https://github.com/sponsors/weprodev">Sponsor</a> ·
  <a href="https://github.com/weprodev">WeProDev</a>
</p>
