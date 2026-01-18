# System Design & Architecture

This document describes the structure and design principles of the **WPD Message Gateway**.

## Architecture Overview

The gateway follows **Clean Architecture** principles with clear separation between layers:

```text
┌─────────────────────────────────────────────────────────────────┐
│                        External World                           │
│  (HTTP Clients, Go Applications using pkg/gateway)              │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Presentation Layer                           │
│                  (internal/presentation/)                       │
│  ┌─────────────┐   ┌─────────────┐  ┌─────────────┐             │
│  │   Router    │   │  Gateway    │  │   DevBox    │             │
│  │             │──▶│  Handler    │  │   Handler   │             │
│  └─────────────┘   └─────┬───────┘  └─────┬───────┘             │
└──────────────────────────┼────────────────┼─────────────────────┘
                           │                │
                           ▼                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Core Layer                                 │
│                    (internal/core/)                             │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                 GatewayService                          │    │
│  │  - SendEmail()   - SendSMS()                            │    │
│  │  - SendPush()    - SendChat()                           │    │
│  └────────────────────────┬────────────────────────────────┘    │
│                           │                                     │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Registry                             │    │
│  │  (Thread-safe provider management)                      │    │
│  └────────────────────────┬────────────────────────────────┘    │
│                           │                                     │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                  Ports (Interfaces)                     │    │
│  │  EmailSender | SMSSender | PushSender | ChatSender      │    │
│  └─────────────────────────────────────────────────────────┘    │
└──────────────────────────┬──────────────────────────────────────┘
                           │ (implements)
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                           │
│               (internal/infrastructure/)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │   Memory    │  │   Mailgun   │  │  SendGrid   │              │
│  │  Provider   │  │  Provider   │  │  Provider   │  ...         │
│  └──────┬──────┘  └─────────────┘  └─────────────┘              │
│         │                                                       │
│         │ (if mailpit.enabled)                                  │
│         ▼                                                       │
│  ┌─────────────┐                                                │
│  │  Mailpit    │                                                │
│  │  Forwarder  │                                                │
│  └─────────────┘                                                │
└─────────────────────────────────────────────────────────────────┘
```

## Request Flow

### Production Flow (Real Provider)

```text
[ Your App ]
     │
     │ SendEmail(email)
     ▼
[ GatewayService ]
     │
     │ providers.defaults.email = mailgun
     ▼
[ Mailgun Provider ]
     │
     │ Convert to API Request
     ▼
( Mailgun API → Email Delivered )
```

### Development Flow (Memory Provider)

```text
[ Your App ]
     │
     │ SendEmail(email)
     ▼
[ GatewayService ]
     │
     │ providers.defaults.email = memory
     ▼
[ Memory Provider ]
     │
     ├──────────────────┬─────────────────────┐
     │                  │                     │
     ▼                  ▼                     ▼
[ RAM Storage ]   [ SMTP Forwarder ]    [ SSE Events ]
     │            (if mailpit.enabled)        │
     ▼                  │                     ▼
[ DevBox API ]          ▼             [ Real-time UI ]
GET /api/v1/emails   [ Mailpit ]      (via EventSource)
                   http://localhost:10103
```

### Configuration Options

```text
┌──────────────────────┬──────────────┬─────────────┐
│ Config                │ DevBox UI    │ Mailpit     │
├──────────────────────┼──────────────┼─────────────┤
│ memory only          │ ✅ Yes       │ ❌ No        │
│ memory + mailpit     │ ✅ Yes       │ ✅ Yes       │
│ mailgun (production) │ ❌ No        │ ❌ No        │
└──────────────────────┴──────────────┴─────────────┘
```

## Directory Structure

```text
wpd-message-gateway/
├── cmd/
│   └── server/              # HTTP server entry point (lean)
│       └── main.go
│
├── configs/                 # YAML configuration files
│   ├── local.yml            # Local development config
│   └── local.example.yml    # Example configuration
│
├── internal/                # Private application code
│   ├── app/                 # Application bootstrap
│   │   ├── config.go        # Configuration structs & loading
│   │   ├── providers.go     # Provider factory (uses registry)
│   │   ├── validation.go    # Configuration validation
│   │   └── wire.go          # Dependency injection
│   │
│   │
│   │   └── registry/        # Provider registration (sub-package)
│   │       └── registry.go  # RegisterEmailProvider, etc.
│   │
│   ├── core/                # Business logic (domain)
│   │   ├── port/            # Interface definitions
│   │   │   ├── email.go     # EmailSender interface
│   │   │   ├── sms.go       # SMSSender interface
│   │   │   ├── push.go      # PushSender interface
│   │   │   └── chat.go      # ChatSender interface
│   │   └── service/
│   │       ├── gateway_service.go  # Core business logic
│   │       └── registry.go         # Provider registry
│   │
│   ├── infrastructure/      # External integrations
│   │   └── provider/        # Provider implementations
│   │       ├── mailgun/     # Mailgun email provider
│   │       └── memory/      # In-memory provider (DevBox)
│   │           ├── store.go # Thread-safe message store
│   │           ├── email.go # Email provider
│   │           ├── sms.go   # SMS provider
│   │           ├── push.go  # Push provider
│   │           └── chat.go  # Chat provider
│   │
│   └── presentation/        # HTTP layer
│       ├── router.go        # Route definitions
│       └── handler/
│           ├── gateway_handler.go  # /v1/* endpoints
│           └── devbox_handler.go   # /api/v1/* endpoints
│
├── pkg/                     # Public packages
│   ├── contracts/           # Message types (single source of truth)
│   │   ├── email.go
│   │   ├── sms.go
│   │   ├── push.go
│   │   ├── chat.go
│   │   └── message.go       # SendResult, Attachment
│   ├── errors/              # Structured error types
│   └── gateway/             # Embedded SDK for Go applications
│       └── gateway.go
│
├── web/                     # DevBox React UI
└── tests/bruno/             # API test collection
```

## Core Concepts

### 1. GatewayService (The Orchestrator)

The **GatewayService** is the central business logic layer. It abstracts away provider management and routing.

- **Location**: `internal/core/service/gateway_service.go`
- **Responsibility**: Route messages to the correct provider based on configuration
- **Benefit**: Your application only interacts with a single service interface

### 2. Ports (The Interfaces)

Ports define **capabilities** — the "What" (Send Email), not the "How" (using Mailgun API).

- **Location**: `internal/core/port/`
- **Interfaces**: `EmailSender`, `SMSSender`, `PushSender`, `ChatSender`
- **Benefit**: Providers are interchangeable — any implementation that satisfies the interface works

### 3. Providers (The Adapters)

Providers are concrete implementations of Ports. They translate generic requests into vendor-specific API calls.

- **Location**: `internal/infrastructure/provider/`
- **Responsibility**: Handle vendor-specific logic, rate limiting, retries
- **Benefit**: Vendor logic is isolated and testable

### 4. Contracts (The Public API)

Contracts define the message types used across the system — the single source of truth.

- **Location**: `pkg/contracts/`
- **Types**: `Email`, `SMS`, `PushNotification`, `ChatMessage`, `SendResult`
- **Benefit**: Consistent types for both internal code and external consumers

### 5. Registry (Thread-Safe Provider Management)

The Registry manages provider instances with thread-safe access.

- **Location**: `internal/core/service/registry.go`
- **Responsibility**: Store and retrieve providers by name
- **Benefit**: Safe concurrent access in HTTP server context

### 6. Provider Self-Registration Pattern

Providers register themselves via Go's `init()` mechanism, following the **Open/Closed Principle**:

```text
┌──────────────────────────────────────────────────────────────────────┐
│  Adding a new provider requires NO modifications to existing code     │
└──────────────────────────────────────────────────────────────────────┘

┌─────────────────────┐         init()           ┌────────────────────────┐
│  Provider Package   │─────────────────────────▶│  Provider Registry     │
│  (register.go)      │   RegisterEmailProvider  │ (internal/app/registry)│
│                     │                          │                        │
│  sendgrid/          │                          │  emailFactories[       │
│   ├── sendgrid.go   │                          │    "sendgrid"          │
│   └── register.go ◄─┤                          │  ] = factory           │
└─────────────────────┘                          └────────────────────────┘
```

- **Location**: `internal/app/registry/` (sub-package), `provider/*/register.go` (registrations)
- **Mechanism**: Each provider has a `register.go` with `init()` that calls `registry.RegisterEmailProvider()`
- **Benefit**: Add providers by creating files, not modifying config or factory code

## Design Principles

### Clean Architecture

- **Dependency Rule**: Dependencies point inward. Infrastructure depends on Core, never the reverse.
- **Testability**: Core logic can be tested without HTTP or external APIs.

### SOLID Principles

| Principle | Application |
|-----------|-------------|
| **Single Responsibility** | Each provider handles one vendor. GatewayService handles routing. |
| **Open/Closed** | Add new providers without modifying existing code. |
| **Liskov Substitution** | Any `EmailSender` implementation works interchangeably. |
| **Interface Segregation** | Separate interfaces for Email, SMS, Push, Chat. |
| **Dependency Inversion** | Core depends on abstractions (ports), not implementations. |

### Other Principles

- **KISS**: Minimal API surface — `Send(ctx, message)` is all you need
- **DRY**: Types defined once in `pkg/contracts/`, reused everywhere
- **12-Factor App**: Configuration via YAML files with environment variable overrides

## Public SDK (`pkg/gateway`)

For Go applications that want to use the gateway as a library:

```go
import (
    "github.com/weprodev/wpd-message-gateway/pkg/contracts"
    "github.com/weprodev/wpd-message-gateway/pkg/gateway"
)

// Simple: memory provider (no config needed)
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "memory",
})

// Production: with provider config
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "mailgun",
    EmailProviders: map[string]gateway.EmailConfig{
        "mailgun": {
            CommonConfig: gateway.CommonConfig{APIKey: "key-xxx"},
            Domain:       "mg.example.com",
            FromEmail:    "noreply@example.com",
        },
    },
})

gw.SendEmail(ctx, &contracts.Email{...})
```

The SDK:
- Uses `registry` types (re-exported as `gateway.EmailConfig`, etc.)
- Leverages the same provider self-registration pattern as the server
- Provides a clean, minimal public API

## Related Documentation

- [Usage Guide](./usage.md) — How to use the package
- [Contributing](./contributing.md) — Guide for adding new providers
- [DevBox](./devbox.md) — Development UI for testing
- [Workflow](./workflow.md) — CI/CD and release process
