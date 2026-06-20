# System Design & Architecture

This document describes the structure and design principles of the **WPD Message Gateway** following the current refactored architecture.

---

## Two Modes of Operation

The gateway is designed to serve two completely different use cases from a single codebase:

```
┌──────────────────────────────────────────────────────────────────────┐
│  MODE 1 — Go Package (Embedded SDK)                                  │
│                                                                      │
│  go get github.com/weprodev/wpd-message-gateway                      │
│                                                                      │
│  ┌────────────────┐     ┌─────────────────────┐                      │
│  │  Your Go App   │────▶│  pkg/gateway.New()  │                      │
│  └────────────────┘     │  (no server, no DB) │                      │
│                         └──────┬──────────────┘                      │
│                                 │ uses registry                      │
│                                 ▼                                    │
│                   ┌─────────────────────────┐                        │
│                   │  Provider (Mailgun, etc)│                        │
│                   └─────────────────────────┘                        │
│                                                                      │
│  Config lives in your code. No infrastructure required.              │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│  MODE 2 — HTTP Server (Full Stack)                                   │
│                                                                      │
│  git clone → make start                                              │
│                                                                      │
│  ┌───────────┐  ┌──────────────┐  ┌─────────────────┐                │
│  │ React UI  │  │  REST API    │  │  Any Language   │                │
│  │ (Portal)  │  │  /api/v1/*   │  │  HTTP client    │                │
│  └─────┬─────┘  └──────┬───────┘  └────────┬────────┘                │
│        │               │                   │ /v1/email               │
│        └───────────────┴───────────────────┘                         │
│                        │                                             │
│                        ▼                                             │
│             ┌────────────────────┐                                   │
│             │   Go HTTP Server   │                                   │
│             │  (Echo framework)  │                                   │
│             └──────────┬─────────┘                                   │
│                        │                                             │
│             ┌──────────▼──────────┐                                  │
│             │  PostgreSQL DB      │ ← single source of truth         │
│             │  workspaces         │   for ALL configuration          │
│             │  integrations       │   (provider credentials,         │
│             │  api_keys           │    workspace settings,           │
│             │  templates          │    members, logs, etc.)          │
│             │  workspace_settings │                                  │
│             └─────────────────────┘                                  │
│                                                                      │
│  Config lives in PostgreSQL — configured via the React UI.           │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Mode 1: Embedded Go SDK

For Go applications that want to send messages **without running any infrastructure**.

```go
import (
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/gateway"
)

// Pass provider config directly — no DB, no server
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

gw.SendEmail(ctx, &contracts.Email{
    To:      []string{"user@example.com"},
    Subject: "Hello",
    HTML:    "<h1>Hi!</h1>",
})
```

**Key characteristics:**

- No PostgreSQL required
- No HTTP server required
- Provider credentials live in your application config (env vars, secrets manager, etc.)
- Uses the same provider registry and contracts as the server mode
- Memory provider available for testing (no external dependencies)

---

## Mode 2: HTTP Server Architecture

### Full System Architecture

```text
┌─────────────────────────────────────────────────────────────────┐
│                        External World                           │
│  ┌────────────────────┐      ┌──────────────────────────────┐   │
│  │   React Portal UI  │      │  HTTP Clients (any language) │   │
│  │  (config, auth,     │      │  POST /v1/email             │   │
│  │   logs, templates) │      │  POST /v1/sms  etc.          │   │
│  └────────────────────┘      └──────────────────────────────┘   │
└──────────────┬──────────────────────────┬───────────────────────┘
               │ /api/v1/*                │ /v1/*
               │ (Portal JWT)             │ (Workspace JWT or API key)
               ▼                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Presentation Layer                           │
│                  (internal/presentation/)                       │
│  ┌─────────────┐   ┌─────────────────┐  ┌─────────────────┐     │
│  │   Router    │   │  PortalHandler  │  │  GatewayHandler │     │
│  │             │   │  (auth, WS mgmt)│  │  (send email,   │     │
│  └─────────────┘   │  InboxHandler   │  │   sms, push,    │     │
│                    │  (captured msgs)│  │   chat)         │     │
│                    └────────┬────────┘  └────────┬────────┘     │
└─────────────────────────────┼────────────────────┼──────────────┘
                              │                    │
                              ▼                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Core Layer                                 │
│                    (internal/core/)                             │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │   PortalService                                          │   │
│  │   Email+password auth · Workspaces · Members · API keys  │   │
│  │   Integrations · Templates · Settings · Logs             │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │   GatewayService                                         │   │
│  │   SendEmail() · SendSMS() · SendPush() · SendChat()      │   │
│  │   Reads dispatch_mode from workspace_settings            │   │
│  │   Reads provider credentials from integrations table     │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │   Ports (Interfaces)                                     │   │
│  │   EmailSender · SMSSender · PushSender · ChatSender      │   │
│  │   Repository interfaces for all DB entities              │   │
│  └──────────────────────────────────────────────────────────┘   │
└──────────────────────────────┬──────────────────────────────────┘
                               │ implements
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                           │
│               (internal/infrastructure/)                        │
│  ┌─────────────────────┐   ┌───────────────────────────────┐    │
│  │  Inbox UI Capture   │   │  Repository Implementations   │    │
│  │  ─────────────────  │   │  ──────────────────────────── │    │
│  │  inbox/             │   │  postgres/                    │    │
│  │  (portal store)     │   │    WorkspaceRepository        │    │
│  │                     │   │    APIKeyRepository           │    │
│  └──────────┬──────────┘   │    IntegrationRepository      │    │
│             │              │    ... (all entities)         │    │
│    ┌────────▼──────┐       └───────────────────────────────┘    │
│    │  Memory Store │                                            │
│    │  (in-process  │                                            │
│    │   inbox for   │                                            │
│    │   dev/testing)│                                            │
│    └───────────────┘                                            │
└─────────────────────────────────────────────────────────────────┘
                               │
                               ▼
                  ┌────────────────────────┐
                  │   PostgreSQL Database  │
                  │   (all persistent data)│
                  └────────────────────────┘
```

---

## Authentication & Authorization Model

### Portal Access (UI, Management API & RBAC)

Portal accounts use **email + password** (passwords are stored hashed). All management requests to `/api/v1/workspaces/:wid/...` require a valid JWT token (`Authorization: Bearer <portal-jwt>`) and are authorized via Role-Based Access Control (RBAC) powered by the `wpd-gogate` library.

The core service layer decouples from `wpd-gogate` by depending on the `port.AuthorizationGate` interface. The implementation adapter lives in `internal/infrastructure/authgate/gate_adapter.go`.

- **admin**: Full read/write access to settings, integrations, templates, API keys, and member removal. Assigned automatically to the creator of a workspace.
- **member**: Read-only access to workspaces, settings, templates, API keys, plus permission to send test messages (`send.test`).

Public workspaces (`is_private = false`) are dynamically accessible to any authenticated user as a `"viewer"` with read-only permissions (i.e. `*.read` operations are bypassed and approved automatically without requiring explicit membership or Casbin checks). All write operations on public workspaces remain strictly restricted to workspace admins.

### Send API Auth (Machine-to-Machine)

Sending messages via `/v1/*` requires **workspace API credentials** and the workspace unique key:

```text
POST /v1/email
Authorization: Bearer <workspace-jwt-token>
    OR
X-Api-Client-Id: wk_abc123
X-Api-Client-Secret: <secret>
X-Workspace-Key: <workspace-unique-key>
```

### Dual Authorization Design

The application enforces two distinct authorization streams to balance security and usability:
1. **Portal REST API**: Restricts access using JWT tokens combined with fine-grained `wpd-gogate` permission middleware (`RequirePermission`).
2. **Inbox SSE API (`/api/v1/workspaces/:wid/inbox/...`)**: Restricts access using JWT tokens combined with a direct database workspace membership check (`RequireWorkspaceMember`) and the workspace API key (`RequireWorkspaceAPIKey`). This design allows client-side SDKs, SSE event streaming, and automation runners to interact with the simulated inbox without requiring full portal RBAC configuration.


---

## PostgreSQL — Single Source of Truth

When running as an HTTP server, **PostgreSQL holds all configuration**. There are no YAML files for provider credentials.

### Database Schema (authoritative references)

To avoid docs drifting from the real schema, treat these as the sources of truth:

- **Migration SQL**: `database/migrations/*.up.sql`
- **ERD**: `docs/assets/database-schema.drawio`

This doc intentionally does **not** duplicate full column lists.

### Configuring providers (server mode)

Provider credentials are stored AES-encrypted in the `integrations` table. Configure via the **Integrations** page in the Portal UI or `POST /api/v1/workspaces/:wid/integrations`.

### Portal UI coverage

The React Portal (`frontend/`) currently implements: auth, workspace list, **integrations**, message **logs**, and **send test**. API keys, templates, members, settings, and memory inbox browsing are REST/Bruno only.

---

## Message Dispatch Modes

Each workspace independently controls how outbound messages are handled:

| Mode                  | Behavior                                                                            |
| --------------------- | ----------------------------------------------------------------------------------- |
| `memory_only`         | Captured in-process RAM only. No external provider called. Default for development. |
| `provider_only`       | Sent through the connected integration only. No in-memory copy.                     |
| `memory_and_provider` | Stored in memory AND sent through the integration.                                  |

Configured via `PATCH /api/v1/workspaces/:wid/settings` (REST — no Portal UI page yet):

```
PATCH /api/v1/workspaces/:wid/settings
{ "message_dispatch_mode": "provider_only" }
```

---

## Request Flow: Sending a Message

### Server Mode — `POST /v1/email`

```text
HTTP Client (your app)
    │
    │ POST /v1/email
    │ Authorization: Bearer <workspace-jwt>
    │ X-Workspace-Key: myapp
    │ { "to": [...], "subject": "...", "html": "..." }
    ▼
APIKeyAuthMiddleware
    │ Validates workspace-jwt OR client_id+secret
    │ Resolves workspace UUID from X-Workspace-Key
    │ Sets workspace_id in request context
    ▼
GatewayHandler.HandleSendEmail()
    │ Reads workspace_id from context
    ▼
GatewayService.SendEmail(ctx, workspaceID, email)
    │ Reads workspace_settings.message_dispatch_mode from DB
    │
    ├── memory_only  → Memory Provider (in-process RAM)
    │                        │
    │                        ▼
    │                  Portal Inbox (SSE + REST)
    │
    ├── provider_only → reads integrations table for workspace+channel
    │                  Decrypts AES config
    │                  Instantiates provider via registry
    │                  → sends to Mailgun/etc
    │
    └── memory_and_provider → both paths above
```

### SDK Mode — `gateway.New(config).SendEmail(...)`

```text
Your Go App
    │
    │ gw.SendEmail(ctx, email)
    ▼
pkg/gateway.Gateway
    │ No DB lookup — config was passed to New()
    ▼
Provider Registry
    │ Finds factory for DefaultEmailProvider name
    ▼
Provider (Mailgun/Memory/etc)
    │ Sends message
    ▼
contracts.SendResult
```

---

## Workspace Isolation

Every entity in the system is scoped to a workspace:

```text
Workspace "myapp"
├── Members (users with roles: admin, member)
├── API Keys (for machine-to-machine sends)
├── Integrations (provider credentials per channel)
│   ├── email: mailgun { api_key: "...", domain: "..." }
│   └── sms: (not configured)
├── Settings { message_dispatch_mode: "provider_only" }
├── Templates (email HTML templates)
└── Message Logs (request audit trail)
```

A user can be a member of **multiple workspaces**. Each workspace has **its own** provider config, API keys, members, messages, etc.

---

## Directory Structure

```text
wpd-message-gateway/
├── cmd/
│   └── server/              # HTTP server entry point
│       └── main.go
│
├── configs/                 # Minimal server config (NO provider credentials)
│   ├── local.yml            # Port, JWT secret
│   └── local.example.yml
│
├── database/
│   ├── migrations/          # SQL migrations (schema evolution)
│   └── seeds/               # Optional demo data
│
├── internal/                # Private application code
│   ├── app/                 # Bootstrap and wiring
│   │   ├── config.go        # Server config (port, JWT)
│   │   ├── wire.go          # Dependency injection — wires all layers
│   │   ├── validation.go    # Config validation
│   │   └── imports.go       # Blank imports to trigger provider init()
│   │
│   ├── core/                # Business logic (domain-pure)
│   │   ├── domain/          # Types: Workspace, User, APIKey, Integration...
│   │   ├── port/            # Interfaces: EmailSender, repositories...
│   │   ├── service/
│   │   │   ├── gateway_service.go   # Message dispatch (reads DB integrations)
│   │   │   └── portal_service.go   # Auth, workspaces, API keys, etc.
│   │   └── authjwt/         # JWT sign/parse
│   │
│   ├── infrastructure/      # DB + provider + auth adapters
│   │   ├── authgate/        # wpd-gogate RBAC adapter implementation
│   │   ├── inbox/           # In-process UI capture store
│   │   ├── logger/          # Application-wide context-aware logger
│   │   └── repository/
│   │       └── postgres/    # All PostgreSQL repository implementations
│   │
│   ├── presentation/        # HTTP layer
│   │   ├── router.go        # All route definitions
│   │   ├── handler/
│   │   │   ├── gateway_handler.go       # POST /v1/* (send messages)
│   │   │   ├── portal_handler.go        # /api/v1/* (auth, workspace CRUD)
│   │   │   ├── portal_inbox_handler.go  # /api/v1/workspaces/:wid/inbox/*
│   │   │   └── send_helper.go           # Shared message dispatcher & logger
│   │   └── middleware/
│   │       ├── auth.go              # API key auth for /v1/*
│   │       ├── portal_jwt.go        # Portal JWT for /api/v1/*
│   │       ├── portal_inbox_auth.go # Inbox-specific auth (JWT + member + key)
│   │       └── rbac.go              # wpd-gogate permission validator
│   │
│
├── pkg/                     # Public packages (imported by external Go apps)
│   ├── contracts/           # Message types: Email, SMS, PushNotification...
│   ├── auth/                # Hash utilities (shared between SDK and server)
│   ├── encryption/          # AES encryption (for DB-stored provider config)
│   ├── provider/            # Provider Adapters (mailgun, memory, etc)
│   ├── registry/            # Provider factory registry
│   └── gateway/             # Embedded SDK — gateway.New() for pure Go usage
│       ├── gateway.go       # Gateway struct, SendEmail/SMS/Push/Chat
│       ├── config.go        # Config struct + New() constructor
│       └── errors.go
│
├── frontend/                # React Portal UI (Vite + TypeScript + Tailwind)
│   └── src/
│       ├── features/        # auth, workspaces, inbox (message logs)
│       └── components/      # design system components
│
└── docs/                    # Documentation
```

---

## Core Concepts

### 1. Workspaces — Tenant Boundary

A **workspace** is the primary isolation unit. All resources (API keys, integrations, templates, members, logs) are scoped to a workspace. Your applications connect to a specific workspace. You can have multiple workspaces for different products, environments, or teams.

### 2. Integrations — DB-Stored Provider Config (Server Mode)

In server mode, provider credentials (Mailgun API keys, etc.) are stored **encrypted** in the `integrations` table. Configure them via the Portal **Integrations** page or the REST API.

### 3. Ports (Interfaces)

Ports define **capabilities** — the "What" (Send Email), not the "How" (using Mailgun API).

- **Sender interfaces**: `pkg/contracts/` — `EmailSender`, `SMSSender`, `PushSender`, `ChatSender`
- **Server ports**: `internal/core/port/` — repository contracts, `InboxReader`, `InboxWriter`
- **Message payloads** (`Email`, `SMS`, `SendResult`, …): defined **only** in `pkg/contracts/` — ports and services use those types; they are not duplicated under `internal/core/domain/`.
- **Benefit**: Providers are interchangeable — any implementation that satisfies the interface works

### 4. Provider Self-Registration

Providers register via Go's `init()` mechanism (Open/Closed Principle):

```text
Provider Package          Provider Registry
(register.go)   ──init()──▶  (pkg/registry)
                 RegisterEmailProvider("mailgun", factory)

Adding a new provider requires NO changes to existing code.
Only create new files in pkg/provider/<name>/
```

### 5. Memory Provider & Dispatch Modes

The **memory provider** captures messages in process RAM. It's always available — no external service needed. Combined with `dispatch_mode`, you can:

- Use `memory_only` for local development (captured messages via inbox REST API / Bruno)
- Use `provider_only` in production (messages go to real provider)
- Use `memory_and_provider` to keep a local copy AND send to the real provider

### 6. Portal (always on)

The React Portal runs at `portal.ui_port` (default **10104**) when the server starts.

**Portal UI today:** sign in, list workspaces, manage **integrations**, view message logs, send test messages.

**REST only (no UI pages yet):** workspace create, API keys, templates, members, settings, memory inbox browser.

---

## Request Tracing & Correlation

To enable robust trace tracking across all layers, the gateway employs an end-to-end correlation ID pipeline:

1. **Correlation Generation**: Echo's `RequestID` middleware generates a unique request UUID for every incoming HTTP request, placing it in the `X-Request-ID` header.
2. **Context Propagation**: A custom correlation middleware extracts this request ID and injects it into the standard Go `context.Context` (accessible via `logger.GetRequestID(ctx)`).
3. **Structured slog Hooking**: An application-wide custom `ContextHandler` intercepts all standard `slog.InfoContext` / `slog.ErrorContext` calls. It transparently extracts the `request_id`, `workspace_id`, `api_key_id`, `channel`, and `provider` attributes from the context and automatically appends them to printed log records without requiring manual logging boilerplate.
4. **Audit Trail Tracking**: The presentation-layer `SendHelper` automatically populates the `request_id` column on the `message_request_logs` PostgreSQL table, linking individual requests directly to database audit trails.

---

## Design Principles

### Clean Architecture

- **Dependency Rule**: Dependencies point inward. Infrastructure depends on Core, never the reverse.
- **Testability**: Core logic can be tested without HTTP, without PostgreSQL, without external APIs.
- **Ports & Adapters**: The Core defines interfaces (ports); Infrastructure provides implementations (adapters).

### SOLID Principles

| Principle                 | Application                                                       |
| ------------------------- | ----------------------------------------------------------------- |
| **Single Responsibility** | Each provider handles one vendor. GatewayService handles routing. |
| **Open/Closed**           | Add new providers without modifying existing code.                |
| **Liskov Substitution**   | Any `EmailSender` implementation works interchangeably.           |
| **Interface Segregation** | Separate interfaces for Email, SMS, Push, Chat.                   |
| **Dependency Inversion**  | Core depends on abstractions (ports), not implementations.        |

### Other Principles

- **KISS**: Minimal API surface — `Send(ctx, message)` is all you need
- **DRY**: Types defined once in `pkg/contracts/`, reused everywhere
- **DB-First Config**: In server mode, PostgreSQL is the single source of truth — no credentials in files

---

## Related Documentation

- [Usage Guide](./usage.md) — SDK and HTTP API usage
- [Contributing](./contributing.md) — Adding providers
- [Portal inbox](./portal-inbox.md) — Memory capture and inbox REST API
- [E2E Testing](./e2e-testing.md) — Automated testing patterns
- [Workflow](../workflow.md) — CI/CD and release process
