# DevBox - Development Message Inbox

A web UI to view all messages sent during development. No real emails, SMS, or push notifications are sent — everything is captured locally.

## Why Use DevBox?

- **See all messages** — Email, SMS, Push, and Chat in one place
- **No real sends** — Messages stay on your computer
- **Real-time updates** — New messages appear instantly via SSE
- **Test easily** — Verify your app sends the right messages

## Quick Start

### 1. Configure Memory Provider

```yaml
# configs/local.yml
providers:
  defaults:
    email: memory
    sms: memory
    push: memory
    chat: memory
```

### 2. Start Everything

```bash
make start    # Starts Gateway server + DevBox UI
```

### 3. Open Browser

Go to http://localhost:10104

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

When providers are set to `memory`, messages are stored in RAM instead of being sent. The DevBox UI fetches these messages via REST API and receives real-time updates via Server-Sent Events (SSE).

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
| `/api/v1/stats` | GET | Message counts by type |
| `/api/v1/emails` | GET | List all emails |
| `/api/v1/emails/{id}` | GET | Get single email |
| `/api/v1/emails/{id}` | DELETE | Delete an email |
| `/api/v1/sms` | GET | List all SMS |
| `/api/v1/sms/{id}` | DELETE | Delete an SMS |
| `/api/v1/push` | GET | List all push notifications |
| `/api/v1/push/{id}` | DELETE | Delete a push notification |
| `/api/v1/chat` | GET | List all chat messages |
| `/api/v1/chat/{id}` | DELETE | Delete a chat message |
| `/api/v1/messages` | DELETE | Clear all messages |
| `/api/v1/events` | GET | Real-time updates (SSE) |

## E2E Testing Example

```go
// In your test
resp, _ := http.Get("http://localhost:10101/api/v1/emails")

var response struct {
    Emails []struct {
        ID    string `json:"id"`
        Email struct {
            Subject string   `json:"subject"`
            To      []string `json:"to"`
        } `json:"email"`
    } `json:"emails"`
}
json.NewDecoder(resp.Body).Decode(&response)

// Assert the email was sent
assert.Equal(t, 1, len(response.Emails))
assert.Equal(t, "Welcome!", response.Emails[0].Email.Subject)
```

## Mailpit Integration (Optional)

For realistic email preview with HTML rendering, you can optionally forward emails to Mailpit:

```bash
# 1. Start Mailpit
make mailpit

# 2. Enable in configs/local.yml:
mailpit:
  enabled: true

# 3. Start server
make start
```

With Mailpit enabled, emails are:
- Stored in DevBox (viewable at http://localhost:10104)
- Forwarded to Mailpit (viewable at http://localhost:10103)

This is useful when you need to preview HTML email templates with proper rendering.

## Configuration

```yaml
# configs/local.yml

# DevBox UI settings
devbox:
  enabled: true
  port: 10104    # DevBox UI port

# Memory provider for all message types
providers:
  defaults:
    email: memory
    sms: memory
    push: memory
    chat: memory

# Optional: Forward emails to Mailpit
mailpit:
  enabled: false  # Set to true when running Mailpit
```

## Tech Stack

- **Backend**: Go (memory provider + REST API + SSE)
- **Frontend**: React 19, TypeScript, Tailwind CSS, shadcn/ui
- **Build**: Vite

## Related

- [Usage Guide](./usage.md) — How to use the message gateway
- [Architecture](./architecture.md) — System design
