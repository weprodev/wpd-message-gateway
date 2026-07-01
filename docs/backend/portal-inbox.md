# Portal — Message Inbox

The Portal is the web UI and REST surface for the WPD Message Gateway when running in server mode.

---

## What exists today

### Portal UI (http://localhost:10104)

| Feature | Description |
|---------|-------------|
| **Authentication** | Register and sign in (email + password → JWT) |
| **Workspaces** | List workspaces you belong to; open a workspace dashboard |
| **Integrations** | Connect, activate, deactivate, or remove messaging providers |
| **Message logs** | Outbound send **request audit trail** (not memory inbox capture) |
| **Send test** | Send a test message per channel from the dashboard |

No Portal pages yet for: workspace create, API keys, templates, members, settings, or memory inbox browsing.

### REST API (server)

| Feature | Portal UI | Description |
|---------|:---------:|-------------|
| **Inbox capture** | | Messages stored in RAM (`memory_only` / `memory_and_database`) — [Inbox API](#inbox-api-reference) |
| **Internal ingest** | | Automation writes to inbox (requires Portal auth) |
| **Workspace provisioning** | | Create workspace, API keys, settings — curl/Bruno/CI only; see [E2E bootstrap](./e2e-testing.md) |

Provider credentials are stored in PostgreSQL (encrypted at rest) and managed via the Portal **Integrations** page. Dispatch mode is REST-only (`PATCH /api/v1/workspaces/:wid/settings`).

---

## Getting Started

```bash
make start       # Starts Go server + React Portal UI
```

Open **http://localhost:10104**

Register (email + password), then sign in.

---

## Message Inbox (Memory Capture)

When a workspace's **dispatch mode** is `memory_only` or `memory_and_database`, outbound messages are captured in-process RAM and displayed in the Portal inbox.

```
Your App
   │
   │  POST /v1/email  (workspace: memory_only)
   ▼
Memory Provider
   └──────────────────▶ Portal Inbox (REST + SSE)
                          └── "inbox" tab shows email
```

### Message Types in the Inbox

| Channel | List view | Detail view |
|---------|-----------|-------------|
| **Email** | Subject, recipient, preview | Full HTML render |
| **SMS** | Full message inline | — |
| **Push** | Title, body, data | — |
| **Chat** | Message content | — |

---

## Dispatch Modes

Each workspace controls how its outbound messages are handled:

| `message_dispatch_mode` | Behavior |
|------------------------|----------|
| `memory_only` | Captured in RAM only, **no** external provider called. **Default.** |
| `memory_and_database` | Captured in RAM only; request logs have `retained = true`. |
| `provider_only` | Sent through the connected integration, **no** memory copy. |
| `provider_and_database` | Same dispatch as provider only; request logs have `retained = true`. |

Set via Portal **Settings → Data Retention** or REST:

```bash
PATCH /api/v1/workspaces/:wid/settings
Content-Type: application/json
Authorization: Bearer <portal-jwt>

{ "message_dispatch_mode": "provider_only" }
```

---

## Authentication (Inbox & API)

The inbox and all Portal API endpoints are **not** public. Every request must be authenticated.

### Portal UI & Management API (`/api/v1/*`)

Authenticated via **Portal JWT**:

```bash
# 1. Register (first time) OR login
curl -X POST http://localhost:10101/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "your-password"}'
# → { "token": "<jwt>", "user": {...} }

# 2. Use JWT for all /api/v1/* requests
curl -H "Authorization: Bearer <jwt>" \
  http://localhost:10101/api/v1/workspaces
```

### Inbox Endpoints (`/api/v1/workspaces/:wid/inbox/*`)

Three things are required:

1. **Portal JWT** — `Authorization: Bearer <token>` (identifies the user)
2. **Workspace membership** — user must be in `workspace_members` for `{wid}`
3. **Workspace API key** — `X-Api-Client-Id` + `X-Api-Client-Secret` for that workspace

```bash
curl -H "Authorization: Bearer $PORTAL_JWT" \
     -H "X-Api-Client-Id: $CLIENT_ID" \
     -H "X-Api-Client-Secret: $SECRET" \
     http://localhost:10101/api/v1/workspaces/$WORKSPACE_ID/inbox/emails
```

> **Note:** For GET requests only, JWT may alternatively be passed as `?access_token=<jwt>` (useful for SSE connections from browser EventSource). API key may also be passed as `?client_id=&client_secret=` — prefer headers over HTTPS.

### Send API (`/v1/*`)

Uses API key credentials + `X-Workspace-Key` header. See [Usage Guide](./usage.md).

---

## Inbox API Reference

All paths are under `/api/v1/workspaces/{workspaceId}/inbox/...`  
Required headers: `Authorization: Bearer <portal-jwt>`, `X-Api-Client-Id`, `X-Api-Client-Secret`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `.../inbox/stats` | Message counts for this workspace |
| GET | `.../inbox/emails` | List captured emails |
| GET | `.../inbox/emails/{id}` | Get one email |
| DELETE | `.../inbox/emails/{id}` | Delete one email |
| GET | `.../inbox/sms` | List captured SMS |
| GET | `.../inbox/sms/{id}` | Get one SMS |
| DELETE | `.../inbox/sms/{id}` | Delete one SMS |
| GET | `.../inbox/push` | List captured push notifications |
| GET | `.../inbox/push/{id}` | Get one push notification |
| DELETE | `.../inbox/push/{id}` | Delete one push notification |
| GET | `.../inbox/chat` | List captured chat messages |
| GET | `.../inbox/chat/{id}` | Get one chat message |
| DELETE | `.../inbox/chat/{id}` | Delete one chat message |
| DELETE | `.../inbox/messages` | Clear **all** captured messages for this workspace |
| GET | `.../inbox/events` | SSE stream for real-time updates |

---

## Internal Ingest (Automation)

For writing directly into the inbox via external automation (e.g. from CI or background jobs):

```
POST /api/v1/workspaces/{workspaceId}/internal/email
POST /api/v1/workspaces/{workspaceId}/internal/sms
POST /api/v1/workspaces/{workspaceId}/internal/push
POST /api/v1/workspaces/{workspaceId}/internal/chat
```

These endpoints are protected by Portal JWT auth and workspace membership. They do not require the workspace API key headers used by inbox read routes.

---

## Server Configuration

```yaml
# configs/local.yml

server:
  port: 10101        # Send API + Portal REST API

portal:
  jwt_secret: "your-secret-min-32-chars"
  jwt_ttl_hours: 72
  ui_port: 10104     # React dev server (Portal UI)
```

> **Provider credentials** (Mailgun API keys, etc.) are **not** in this file.  
> They are stored AES-encrypted in PostgreSQL and configured via the **Portal REST API** (no Portal UI page yet).

---

## Tech Stack

- **Backend**: Go — Echo framework, memory provider, PostgreSQL repositories
- **Frontend**: React, TypeScript, Tailwind CSS, shadcn/ui, Vite
- **Realtime**: Server-Sent Events (SSE) for live inbox updates

---

## Related

- [Usage Guide](./usage.md) — Full API reference
- [Architecture](./architecture.md) — System design
- [E2E Testing](./e2e-testing.md) — Using the inbox in CI/CD
