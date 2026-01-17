# DevBox - Development Message Inbox

A web UI to view all messages sent during development. No real emails, SMS, or push notifications are sent — everything is captured locally.

## Why Use DevBox?

- **See all messages** — Email, SMS, Push, and Chat in one place
- **No real sends** — Messages stay on your computer
- **Real-time updates** — New messages appear instantly
- **Test easily** — Verify your app sends the right messages

## Quick Start

### 1. Set Memory Provider

```bash
# Add to your .env or export
MESSAGE_DEFAULT_EMAIL_PROVIDER=memory
MESSAGE_DEFAULT_SMS_PROVIDER=memory
MESSAGE_DEFAULT_PUSH_PROVIDER=memory
MESSAGE_DEFAULT_CHAT_PROVIDER=memory
```

### 2. Start the UI

```bash
make web-install   # First time only
make web-dev       # Start dev server
```

### 3. Open Browser

Go to http://localhost:5173

## How It Works

```
Your App                    DevBox
   │                          │
   │  SendEmail(...)          │
   ▼                          │
Memory Provider ──────────► Web UI shows the email
   │                          │
   │  SendSMS(...)            │
   ▼                          │
Memory Provider ──────────► Web UI shows the SMS
```

When you set providers to `memory`, messages are stored in memory instead of being sent. The DevBox UI fetches these messages via REST API.

## Features

| Message Type | List View | Detail View |
|--------------|-----------|-------------|
| **Email** | Subject, recipient, preview | Full HTML template |
| **SMS** | Full message inline | — |
| **Push** | Title, body, data | — |
| **Chat** | Message preview | Template, media, buttons |

## API Endpoints

The DevBox exposes a REST API for programmatic access (useful for E2E tests):

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/stats` | GET | Message counts |
| `/api/v1/emails` | GET | List all emails |
| `/api/v1/emails/{id}` | DELETE | Delete an email |
| `/api/v1/sms` | GET | List all SMS |
| `/api/v1/push` | GET | List all push notifications |
| `/api/v1/chat` | GET | List all chat messages |
| `/api/v1/messages` | DELETE | Clear all messages |
| `/api/v1/events` | GET | Real-time updates (SSE) |

## E2E Testing Example

```go
// In your test
resp, _ := http.Get("http://localhost:8080/api/v1/emails")
var emails []memory.StoredEmail
json.NewDecoder(resp.Body).Decode(&emails)

// Assert the email was sent
assert.Equal(t, 1, len(emails))
assert.Equal(t, "Welcome!", emails[0].Email.Subject)
```

## Tech Stack

- **Backend**: Go (memory provider + REST API)
- **Frontend**: React 19, TypeScript, Tailwind CSS, shadcn/ui
- **Build**: Vite

## Related

- [Usage Guide](./usage.md) — How to use the message gateway
- [Architecture](./architecture.md) — System design
