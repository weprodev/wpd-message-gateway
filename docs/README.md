# WPD Message Gateway Documentation

The central documentation hub for the WPD Message Gateway.

## Getting Started

- **[Architecture](./backend/architecture.md)** — Core design, layers, and the dual-mode gateway overview.
- **[Development Workflow](./workflow.md)** — CI/CD pipelines, release process, and branch strategy.
- **[Spec Kit development flow](./development-flow.md)** — Spec → plan → tasks → implement → issue update → PR.

## Backend (Go)

Guidelines, specs, and operational playbooks for the Go messaging service.

- **[Backend engineer — role & workflow](./backend/backend-engineer.md)** — Execution model, layers, registry, security, quality gate.
- **[Code Conventions](./backend/code-conventions.md)** — Go coding standards and DDD patterns.
- **[Contributing a New Provider](./backend/contributing.md)** — Step-by-step guide to adding an SMS/Push/Email provider.
- **[API Usage Guide](./backend/usage.md)** — Complete HTTP API and embedded SDK reference.
- **[Portal Inbox System](./backend/portal-inbox.md)** — Core logic behind memory routing and inbox capture.
- **[E2E Testing Gateway](./backend/e2e-testing.md)** — How to use the Gateway to intercept output in your other application's CI tests.

## Frontend (React/Portal UI)

Documentation for the Portal web UI interface.

- **[Frontend Docs Hub](./frontend/README.md)** — The entrypoint for all frontend docs.
- **[Frontend engineer — role & workflow](./frontend/frontend-engineer.md)** — Operating model, React boundaries, security, Storybook, verification.
- **[Frontend Conventions](./frontend/conventions.md)** — Architecture boundaries and file structures.
- **[shadcn/ui Skill](./frontend/shadcn/SKILL.md)** — Rules and usage guide for the `components/ui` library.

## Meta

- **[AI agents](./agents/master-agent.md)** — Playbooks in `docs/agents/`; **[verification chain](./agents/verification.md)** (lint → `/smell` → `make audit`).
- **`/smell`** — Diff review (`.claude/commands/smell.md`, base `develop`).
- **[System architecture (Excalidraw)](./assets/diagram.excalidraw)** — HTTP server, gateway, providers, portal vs memory inbox.
- **[PostgreSQL ERD (draw.io)](./assets/database-schema.drawio)** — Tables and FKs from `database/migrations/*.up.sql`.
