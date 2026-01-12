# System Design & Lifecycle

This document illustrates the structure of the **WPD Message Gateway** using text-based diagrams.

## Architecture Blueprint

This diagram shows how your application connects to the Gateway, how the Gateway manages configurations, and how different Providers implement the core Contracts.

```text
  +---------------------------+
  |  ThirdPartyApp (Your App) |
  +---------------------------+
  | + LoadConfig()            |
  | + InitManager()           |      1. Creates
  | + SendEmail()             |----------------------------------.
  +---------------------------+                                  |
               |                                                 |
               | 2. Initializes                                  v
               |                                       +-------------------+
               v                                       |      Config       |
  +---------------------------+                        +-------------------+
  |          Manager          | 3. Uses                | + DefaultEmail    |
  +---------------------------+<-----------------------| + Providers Map   |
  | - config: Config          |                        | + LoadFromEnv()   |
  | - emailProviders: Map     |                        +-------------------+
  |                           |
  | + New(cfg)                |
  | + SendEmail(ctx, email)   |
  +---------------------------+
               |
               | 4. Manages
               v
  +---------------------------+
  | <<Interface>> EmailSender |
  +---------------------------+
  | + Send(ctx, email)        |
  | + Name() string           |
  +---------------------------+
               ^
              / \  (Implements)
             /   \
            /_____\
               |
      +--------+--------+
      |                 |
+------------+   +-------------+
|   Mailgun  |   |  SendGrid   |
+------------+   +-------------+
| - apiKey   |   | - apiKey    |
| - domain   |   |             |
|            |   |             |
| + Send()   |   | + Send()    |
+------------+   +-------------+
```

## Lifecycle Explanation

1.  **Configuration (The Setup)**
    *   **Box:** `ThirdPartyApp` creates `Config`.
    *   **Action:** Your app loads environment variables (API keys, defaults) into the `Config` struct.

2.  **Initialization (The Wiring)**
    *   **Box:** `ThirdPartyApp` initializes `Manager`.
    *   **Connection:** `Manager` reads the `Config` and instantiates the specific Providers (like `MailgunProvider`) that you defined.

3.  **Execution (The Call)**
    *   **Flow:** `ThirdPartyApp` calls `Manager.SendEmail()`.
    *   **Routing:** `Manager` looks up the correct `EmailSender` (the interface).

4.  **Implementation (The Triangle)**
    *   **Polymorphism:** The `Triangle` symbol (`/_\`) represents inheritance or implementation. `Mailgun` and `SendGrid` **implement** the `EmailSender` interface.
    *   **Result:** The Manager accepts any provider that fits the "shape" of the interface, without knowing the specific details of the provider.

## Request Flow

```text
[ Your App ]
     |
     | SendEmail(email)
     v
[ Manager ]
     |
     | 1. Find Provider ("mailgun")
     v
[ Mailgun Provider ]
     |
     | 2. Convert to API Request
     v
( Internet / External API )
```

## Core Concepts

The gateway is built around four main pillars that separate concerns and ensure extensibility.

### 1. Manager ( The Gateway )
The **Manager** acts as the central entry point for your application. It abstracts away the complexity of managing multiple providers.
*   **Role**: Orchestration and Dispatch.
*   **Responsibility**: Loads configuration, initializes providers, and routes "Send" requests to the correct active provider.
*   **Benefit**: Your application only talks to `Manager`, never directly to Mailgun or Twilio.

### 2. Contracts ( The Interfaces )
Contracts are Go interfaces that define **capabilities**. They represent the "What" (Send Email), not the "How" (using Mailgun API).
*   **Role**: Definition and Abstraction.
*   **Responsibility**: Define strict signatures for `EmailSender`, `SMSSender`, etc.
*   **Benefit**: Allows providers to be swapped easily. If a new provider satisfies the interface, it fits.

### 3. Providers ( The Adapters )
Providers are the concrete implementations of Contracts. They act as adapters between our internal domain and external APIs.
*   **Role**: Implementation.
*   **Responsibility**: Translate the generic request (e.g., `contracts.Email`) into the specific API call required by the vendor (e.g., Mailgun JSON payload).
*   **Benefit**: Vendor-specific logic is isolated. Rate limiting, retries, and HTTP calls happen here.

### 4. Config ( The Configuration )
Configuration acts as the blueprint for the Manager.
*   **Role**: Setup.
*   **Responsibility**: Reads from environment variables (12-factor app) or programmatic input to tell the Manager *which* providers to enable and *what* credentials to use.

---

## Directory Structure

The project structure reflects these concepts:

```text
wpd-message-gateway/
├── config/              # Configuration logic
├── contracts/           # The Interface definitions (Ports)
├── manager/             # The Central Gateway (Orchestrator)
├── providers/           # Concrete Implementations (Adapters)
│   ├── email/           # Email-specific adapters
│   ├── sms/             # SMS-specific adapters
│   └── ...
└── errors/              # Domain-specific error types
```

---

## Design Principles

We adhere strictly to the following engineering principles:

### SOLID Principles
*   **Single Responsibility Principle (SRP)**: A provider implementation only handles sending for its specific service. The Manager only handles routing.
*   **Open/Closed Principle (OCP)**: You can add a new provider (e.g., "Postmark") by creating a new struct in `providers/email/postmark` without modifying a single line of the Manager's core logic.
*   **Interface Segregation (ISP)**: We have separate interfaces for `EmailSender`, `SMSSender`, etc., rather than one giant "MessageSender" interface.
*   **Dependency Inversion (DIP)**: The Manager depends on the abstract `contracts.EmailSender`, not on `mailgun.Provider`.

### Other Principles
*   **KISS (Keep It Simple, Stupid)**: The API surface is minimal. `Send(ctx, message)` is all you need.
*   **DRY (Don't Repeat Yourself)**: Common types like `Attachment` and `SendResult` are defined once in `contracts/message.go` and reused across all provider types.

---

## Related Documentation

- **[System Design & Lifecycle](./system-design.md)** - Visual diagrams of the architecture.
- **[Usage Guide](./usage.md)** - How to use the package in your code.
- **[Contributing](./contributing.md)** - Guide for adding new providers.
1