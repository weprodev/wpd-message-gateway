<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="frontend/public/logo-dark-mode.svg">
    <img src="frontend/public/logo-light-mode.svg" alt="Message Gateway" width="400" />
  </picture>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/weprodev/wpd-message-gateway"><img src="https://pkg.go.dev/badge/github.com/weprodev/wpd-message-gateway.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/weprodev/wpd-message-gateway"><img src="https://goreportcard.com/badge/github.com/weprodev/wpd-message-gateway" alt="Go Report Card"></a>
  <a href="https://github.com/weprodev/wpd-message-gateway/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
</p>
<p align="center">
  <strong>A unified Go library and HTTP API for sending Email, SMS, Push, and Chat messages.</strong>

One interface, multiple providers. Write your messaging code once — switch between Mailgun, Twilio, Firebase, WhatsApp, and more without changing a single line of application code.

</p>

---

## Two Ways to Use

```
┌──────────────────────────────────┬──────────────────────────────────┐
│  Go Package (Embedded SDK)       │  HTTP Server (any language)      │
│                                  │                                  │
│  go get .../wpd-message-gateway  │  git clone → make start          │
│  gateway.New(config)             │  POST /v1/email (HTTP)           │
│                                  │                                  │
│  ✓ No server needed              │  ✓ Any language (Python, JS...)  │
│  ✓ No database needed            │  ✓ React UI to manage config     │
│  ✓ Config in your code           │  ✓ PostgreSQL stores everything  │
└──────────────────────────────────┴──────────────────────────────────┘
```

---

## Quick Start: Go Package

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

## Quick Start: HTTP Server

```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
cp configs/local.example.yml configs/local.yml
make start
```

1. Open **http://localhost:10104** — the Portal UI
2. Create an account (**email + password**) and sign in
3. Create a **workspace**, add an **Integration** (Mailgun, etc.), generate an **API key**
4. Send messages from any language:

```bash
curl -X POST http://localhost:10101/v1/email \
  -u "wk_client_id:your_secret" \
  -H "X-Workspace-Key: myapp" \
  -H "Content-Type: application/json" \
  -d '{"to":["user@example.com"],"subject":"Hello","html":"<h1>World</h1>"}'
```

---

## Why Message Gateway?

Building applications that send messages across multiple channels is complex. You need to integrate different APIs, handle various authentication methods, manage provider-specific quirks, and test everything without spamming real users.

**Message Gateway provides:**

- **Unified API** — Email, SMS, Push, and Chat through a single, consistent interface
- **Provider abstraction** — Switch from Mailgun to SendGrid with a config change—no code changes
- **DB-first config** — In server mode, all provider credentials live in PostgreSQL, managed via the Portal UI
- **Workspace isolation** — Multiple workspaces, each with its own providers, API keys, templates, and members
- **Memory provider** — Captures messages locally for dev and testing, no external services needed
- **E2E testing** — Assert real message payloads in CI/CD without mocking

---

## Supported Message Types

| Type      | Description                           | Providers                 |
| --------- | ------------------------------------- | ------------------------- |
| **Email** | HTML, plain text, CC/BCC, attachments | Mailgun, Memory           |
| **SMS**   | Text to mobile numbers                | Memory (Twilio planned)   |
| **Push**  | Mobile and web notifications          | Memory (Firebase planned) |
| **Chat**  | Slack, WhatsApp, Telegram             | Memory (planned)          |

---

## Portal: Configuration UI

The Portal is **always available** at `http://localhost:10104` when the server runs.

**Access**: Email + password (Portal JWT)

**Manage per workspace:**

- **Integrations** — add provider credentials (encrypted in DB, not in files)
- **API Keys** — credentials for your apps to send messages
- **Templates** — reusable HTML email templates
- **Settings** — dispatch mode (`memory_only` | `provider_only` | `memory_and_provider`)
- **Members** — invite colleagues to your workspace
- **Inbox** — view all captured messages (in `memory_only` mode)
- **Logs** — full audit trail of send requests

---

## Authentication & Gateway Modes

Message Gateway features multiple dispatch modes (`memory_only`, `provider_only`, `memory_and_provider`).

To interact with the Portal, users authenticate via **Email + Password** to receive a JWT. Client applications interacting with the **Send API** use an isolated **Workspace API Key**.

For detailed authentication flow and dispatch mode behavior, see the [Usage Guide](docs/backend/usage.md).

---

## E2E Testing

Instead of mocking messy email or SMS interactions in your test suite, use the Message Gateway's `memory_only` dispatch mode in CI/CD. It intercepts calls over HTTP, storing them in-memory, allowing you to use the Portal Inbox API to strictly assert sent payloads without sending real messages.

Read the [End-to-End Testing Guide](docs/backend/e2e-testing.md) for full docker-compose and GitHub Actions setups.

---

## Project Structure

```
wpd-message-gateway/
├── cmd/server/          # HTTP server entry point
├── configs/             # Server config (port, JWT) — NO provider credentials
├── database/
│   ├── migrations/      # SQL schema migrations
│   └── seeds/           # Optional demo data
├── internal/
│   ├── app/             # Config, wire, validation, provider blank imports
│   ├── core/            # Domain, services, ports
│   ├── infrastructure/ # Providers (mailgun, memory…) + Postgres repos + logger
│   ├── presentation/    # HTTP router, handlers, middleware
│   └── registry/        # Provider factory registry (self-registration via init)
├── pkg/                 # Public packages (contracts, gateway SDK, auth, encryption)
│   ├── contracts/       # Email, SMS, Push, Chat types
│   ├── gateway/         # Embedded Go SDK — gateway.New()
│   ├── auth/            # Password hashing (portal)
│   └── encryption/      # AES helpers
├── frontend/            # React Portal UI (Vite + TypeScript + Tailwind)
├── tests/bruno/         # HTTP API test collection
└── docs/                # Documentation
```

---

## AI-Assisted Development (GitHub Spec Kit)

We use [GitHub Spec Kit](https://github.com/github/spec-kit) **inside this repository** to enforce a strict Specification‑Driven Development workflow (spec → plan → tasks → implement → review).

The key to “best outcomes” is that Spec Kit is _not_ free-form prompting here: our Spec Kit commands are bound to **repository agents** (Principal personas) and must follow our backend DDD + frontend composition rules.

### What Spec Kit Creates (Where to Look)

Every feature gets a directory under `specs/` named after the feature branch, for example `specs/023-portal-inbox-search/`:

- `spec.md`: requirements, user stories, success criteria
- `plan.md`: technical plan + design artifacts (research/data model/contracts/quickstart)
- `tasks.md`: dependency-ordered, file-path-specific implementation tasks
- `checklists/*.md`: “unit tests for requirements writing” (not implementation tests)

### Agent Mapping (How “created Agents” are used)

- **`/speckit.specify` + `/speckit.plan`**: must follow **Master Agent** (`docs/agents/master-agent.md`)
- **`/speckit.tasks` + `/speckit.implement`**: must follow **Delivery Agent** (`docs/agents/delivery-agent.md`)
- **`/speckit.checklist` + `/speckit.analyze`**: must follow **Review Agent** (`docs/agents/review-agent.md`)

Those playbooks are the source of truth for DDD layering, frontend composition, security, and verification.

### The Happy Path (Recommended Sequence)

1. **Create a spec + feature branch**

   Run:
   - **`/speckit.specify <your feature description>`**

   This creates and checks out a feature branch and initializes `specs/<branch>/spec.md`.

2. **Clarify (optional but recommended)**

   Run:
   - **`/speckit.clarify`**

   This asks up to 5 high-impact questions and writes the answers back into `spec.md`.

3. **Plan the implementation**

   Run:
   - **`/speckit.plan`**

   This generates `plan.md` plus Phase 0/1 design artifacts (research, data model, contracts, quickstart) and updates agent context.

4. **Generate tasks**

   Run:
   - **`/speckit.tasks`**

   This generates `tasks.md` with atomic tasks, strict file paths, and dependency ordering.

5. **Sanity-check for consistency (recommended)**

   Run:
   - **`/speckit.analyze`**

   This is read-only and highlights gaps (requirements with no tasks, tasks with no requirement, constitution violations).

6. **Implement**

   Run:
   - **`/speckit.implement`**

   Complete tasks phase-by-phase; keep tasks checked off in `tasks.md` as you go.

7. **Pre-PR quality gates**

   Run:
   - **`/speckit.checklist <domain>`** (e.g. `security`, `api`, `ux`)
   - `make audit`

### Troubleshooting

- **Wrong branch / “not on a feature branch”**: Spec Kit expects to run on a branch that matches a `specs/<branch>/` directory. Re-run `/speckit.specify ...` or set `SPECIFY_FEATURE` in your shell to point at the intended feature directory.
- **“tasks.md not found”**: run `/speckit.tasks` before `/speckit.implement`.
- **CI/quality mismatch**: this repo’s hard quality gate is `make audit`; plans and tasks should always include it as a final validation step.

---

## Commands

```bash
make install    # Install Go + frontend dependencies
make start      # Gateway + Portal UI (Vite)
make ui         # Portal UI only (Vite, port 10104)
make storybook  # Storybook (port 6006)
make test       # Go tests only
make audit      # fmt+lint+test (Go+frontend), govulncheck, builds (Go+Vite+Storybook)
make build      # Go compile + frontend build:all (no tests)
make upgrade    # Upgrade dependencies
make dev        # Run via Docker Compose
```

---

## Documentation

| Document                                                     | Description                                                                 |
| ------------------------------------------------------------ | --------------------------------------------------------------------------- |
| [Docs hub](docs/README.md)                                   | Index of backend + frontend documentation                                   |
| [Usage](docs/backend/usage.md)                               | SDK and HTTP API reference, authentication, multi-language examples         |
| [Architecture](docs/backend/architecture.md)                 | System design, two modes of operation, DB schema                            |
| [Portal inbox](docs/backend/portal-inbox.md)                 | Message inbox, dispatch modes, inbox API                                    |
| [Frontend docs](docs/frontend/README.md)                     | Portal UI index — Vite, TypeScript, shadcn skill, conventions               |
| [Frontend engineer role](docs/frontend/frontend-engineer.md) | Principal-style workflow, architecture, security, Storybook, self-review    |
| [shadcn/ui skill (in-repo)](docs/frontend/shadcn/SKILL.md)   | Component rules, CLI patterns ([official docs](https://ui.shadcn.com/docs)) |
| [Backend engineer role](docs/backend/backend-engineer.md)    | Go layers, registry, security, quality gate                                 |
| [E2E Testing](docs/backend/e2e-testing.md)                   | CI/CD integration, capturing and asserting messages                         |
| [Contributing](docs/backend/contributing.md)                 | Adding new providers                                                        |
| [Code Conventions](docs/backend/code-conventions.md)         | Go coding standards                                                         |
| [Workflow](docs/workflow.md)                                 | CI/CD and release process                                                   |
| [Bruno collections](tests/bruno/)                            | HTTP API tests (`bru run`)                                                  |

---

## Licensing and Sponsorship

Released under the **[MIT License](LICENSE)**. The same grant applies to everyone who receives the code.

We want Message Gateway to stay easy to adopt for individuals and small teams, while asking larger organizations that get sustained value from it to help fund maintenance and features.

| Who                                      | What we ask                                                                  |
| ---------------------------------------- | ---------------------------------------------------------------------------- |
| **Individuals, learning, side projects** | Use freely under MIT. Sponsorship optional.                                  |
| **Small teams**                          | Use freely under MIT. Consider sponsoring if it's central to your stack.     |
| **Mid-size companies and enterprises**   | **Please sponsor** — [GitHub Sponsors](https://github.com/sponsors/weprodev) |
| **Qualifying non-profits**               | Use freely under MIT; sponsorship optional.                                  |

<p align="center">
  <a href="https://github.com/sponsors/weprodev">
    <img src="https://img.shields.io/badge/Sponsor-❤️-ea4aaa?style=for-the-badge" alt="Sponsor on GitHub">
  </a>
</p>

---

## Contributing

1. **Report bugs** — Open an issue
2. **Suggest features** — Ideas welcome
3. **Pull requests** — Code and docs
4. **Add providers** — Expand integrations
5. **Sponsor** — Support ongoing work

See [docs/backend/contributing.md](docs/backend/contributing.md). Run `make audit` before opening a PR.

---

<p align="center">
  Built with ❤️ by <a href="https://github.com/weprodev">WeProDev</a>
</p>
