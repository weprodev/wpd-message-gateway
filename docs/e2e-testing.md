# E2E Testing Guide

Use the message gateway to capture and verify all messages your app sends during CI tests.

## Why Use This?

| Traditional Mocking | With Message Gateway |
|---------------------|----------------------|
| Mock email service at code level | Real HTTP calls to real server |
| Only tests mock was called | Tests actual request/response |
| Can't verify message content easily | Query exact message content via API |
| Need different mocks for email/SMS/push | One service captures all channels |

## Quick Setup

### 1. Add Gateway Service

```yaml
# .github/workflows/test.yml
jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      gateway:
        image: ghcr.io/weprodev/wpd-message-gateway:latest
        ports:
          - 10101:10101
```

### 2. Configure Your App

Point your app's message sending to the gateway:

```bash
# Environment variable
EMAIL_API_URL=http://localhost:10101
```

### 3. Query Captured Messages

```bash
# Get all emails
curl http://localhost:10101/api/v1/emails

# Get all SMS
curl http://localhost:10101/api/v1/sms

# Get message counts
curl http://localhost:10101/api/v1/stats

# Clear messages between tests
curl -X DELETE http://localhost:10101/api/v1/messages
```

---

## Complete Example: Testing User Signup

Your app sends a welcome email when users sign up. Here's how to test it:

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      gateway:
        image: ghcr.io/weprodev/wpd-message-gateway:latest
        ports:
          - 10101:10101
    
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - run: npm ci
      
      - name: Start app
        run: |
          EMAIL_API_URL=http://localhost:10101 npm start &
          sleep 5
      
      - name: Test signup flow
        run: |
          # Clear messages
          curl -X DELETE http://localhost:10101/api/v1/messages
          
          # Trigger signup
          curl -X POST http://localhost:3000/api/signup \
            -H "Content-Type: application/json" \
            -d '{"email": "user@example.com", "name": "John"}'
          
          sleep 2
          
          # Verify email
          EMAILS=$(curl -s http://localhost:10101/api/v1/emails)
          
          echo "$EMAILS" | jq -e '
            .emails | length == 1 and
            .[0].email.to[0] == "user@example.com" and
            .[0].email.subject == "Welcome!"
          ' || exit 1
          
          echo "✅ Welcome email verified"
```

---

## API Reference

### Send Messages

```bash
# Email
curl -X POST http://localhost:10101/v1/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["user@example.com"],
    "subject": "Hello",
    "html": "<p>World</p>"
  }'

# SMS
curl -X POST http://localhost:10101/v1/sms \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["+1234567890"],
    "message": "Your code is 123456"
  }'
```

### Query Captured Messages

| Endpoint | Description |
|----------|-------------|
| `GET /api/v1/emails` | List all captured emails |
| `GET /api/v1/sms` | List all captured SMS |
| `GET /api/v1/push` | List all captured push notifications |
| `GET /api/v1/chat` | List all captured chat messages |
| `GET /api/v1/stats` | Get message counts |
| `DELETE /api/v1/messages` | Clear all messages |

### Response Format

```json
{
  "emails": [
    {
      "id": "abc123",
      "created_at": "2024-01-15T10:30:00Z",
      "email": {
        "to": ["user@example.com"],
        "subject": "Welcome!",
        "html": "<h1>Hello John</h1>",
        "plain_text": "Hello John"
      }
    }
  ]
}
```

---

## Testing Patterns

### Assert Email Count

```bash
COUNT=$(curl -s http://localhost:10101/api/v1/emails | jq '.emails | length')
[ "$COUNT" -eq 1 ] || exit 1
```

### Assert Email Content

```bash
curl -s http://localhost:10101/api/v1/emails | jq -e '
  .emails[0].email.subject == "Welcome!" and
  .emails[0].email.to[0] == "user@example.com"
'
```

### Assert Email Contains Text

```bash
curl -s http://localhost:10101/api/v1/emails | jq -e '
  .emails[0].email.html | contains("verification code")
'
```

### Clear Between Tests

```bash
curl -X DELETE http://localhost:10101/api/v1/messages
```

---

## Go SDK Users

If your app uses the gateway Go SDK, just change the provider in tests:

```go
// Production
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "mailgun",
    EmailProviders: map[string]gateway.EmailConfig{
        "mailgun": {
            CommonConfig: gateway.CommonConfig{APIKey: os.Getenv("MAILGUN_API_KEY")},
            Domain:       "mg.example.com",
        },
    },
})

// Tests - use memory provider (no config needed)
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: "memory",
})

// Messages are stored in memory, query via DevBox API
```

---

## Local Development

Run the gateway locally for development:

```bash
docker run -p 10101:10101 ghcr.io/weprodev/wpd-message-gateway:latest
```

Or with docker-compose:

```yaml
services:
  gateway:
    image: ghcr.io/weprodev/wpd-message-gateway:latest
    ports:
      - "10101:10101"
```
