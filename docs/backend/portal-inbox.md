# Portal — Message Inbox

The Portal is the **primary interface** for the WPD Message Gateway. It is always enabled when running the server. You use it to configure workspaces, manage providers, view captured messages, manage API keys, and more.

---

## What the Portal Does

| Feature | Description |
|---------|-------------|
| **Authentication** | Email + password (Portal JWT) |
| **Workspaces** | Multi-tenant isolation — each workspace has its own config |
| **Integrations** | Configure email/SMS/push/chat providers (stored encrypted in DB) |
| **API Keys** | Create and manage machine-to-machine credentials |
| **Templates** | HTML email templates per workspace |
| **Message Inbox** | View messages captured in memory (`memory_only` mode) |
| **Message Logs** | Audit trail of all outbound send API requests |
| **Settings** | Per-workspace dispatch mode and other configuration |

---

## Getting Started

```bash
make start       # Starts Go server + React Portal UI
```

Open **http://localhost:10104**

Register (email + password), then sign in.

---

## Message Inbox (Memory Capture)

When a workspace's **dispatch mode** is `memory_only` or `memory_and_provider`, outbound messages are captured in-process RAM and displayed in the Portal inbox.

```
Your App
   │
   │  POST /v1/email  (workspace: memory_only)
   ▼
Memory Provider
   ├──────────────────▶ Portal Inbox (REST + SSE)
   │                      └── "inbox" tab shows email
   │
   └──────────────────▶ Mailpit (if enabled)
                          └── HTML preview
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
| `provider_only` | Sent through the connected integration, **no** memory copy. |
| `memory_and_provider` | Stored in memory **and** sent through the integration. |

Set via Portal → **Settings → General**, or via API:

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

For writing into the inbox without Portal JWT (e.g. from CI automation or background jobs):

```
POST /api/v1/workspaces/{workspaceId}/internal/email
POST /api/v1/workspaces/{workspaceId}/internal/sms
POST /api/v1/workspaces/{workspaceId}/internal/push
POST /api/v1/workspaces/{workspaceId}/internal/chat
```

Protected by `X-Internal-Secret` header when `MESSAGE_INTERNAL_INGEST_SECRET` env var is set.  
In local dev (no env var), the endpoint is open.

---

## Optional: Mailpit

For HTML email preview with rich rendering, forward captured emails to Mailpit:

```bash
make mailpit     # Starts Mailpit at http://localhost:10103
```

Enable forwarding in `configs/local.yml`:

```yaml
mailpit:
  enabled: true
```

With Mailpit enabled: emails are stored in the Portal inbox **and** forwarded to Mailpit for rich preview.

| | Portal Inbox | Mailpit |
|--|--|--|
| All channels | ✅ | ❌ email only |
| REST API | ✅ | Limited |
| SSE real-time | ✅ | ✅ |
| HTML render | Basic | ✅ Rich |
| Mobile preview | ❌ | ✅ |

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

# Optional:
# mailpit:
#   enabled: true
```

> **Provider credentials** (Mailgun API keys, etc.) are **not** in this file.  
> They are configured in the Portal UI and stored AES-encrypted in PostgreSQL.

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
