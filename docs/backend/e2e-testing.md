# E2E Testing Guide

Use the Message Gateway to capture and verify all messages your app sends during CI/CD tests — without mocking, without real sends, without external dependencies.

## Why Use This?

| Traditional Mocking | With Message Gateway |
|---------------------|----------------------|
| Mock email service at code level | Real HTTP calls to real endpoint |
| Only tests mock was called | Tests actual request/response payload |
| Can't introspect message content easily | Query exact subject, body, recipients via API |
| Different mocks for email/SMS/push | One service captures all channels |
| No capture of messages from external libs | Captures everything your app sends |

---

## Architecture

```
Your App                 Gateway (memory dispatch)        Test
   │                          │                        │
   │  POST /v1/email ─────────▶│ captured in RAM        │
   │  POST /v1/email ─────────▶│ captured in RAM        │
   │                          │                        │
   │                          │◀─── GET /inbox/emails ─│
   │                          │─── [email1, email2] ───▶│
   │                          │                        │ assert
```

Messages are **never sent** to real providers. Everything stays in the gateway's in-process memory until your test clears it.

---

## Security: What's Required to Read the Inbox

The inbox is protected. Every request to `/api/v1/workspaces/:wid/inbox/*` must provide:

1. **Portal JWT** — `Authorization: Bearer <token>` from `POST /api/v1/auth/login` (or `/api/v1/auth/register`)
2. **Workspace membership** — user must be a member of the workspace
3. **Workspace API key** — `X-Api-Client-Id` + `X-Api-Client-Secret` matching the workspace

The **send** endpoint (`POST /v1/*`) uses API key + `X-Workspace-Key` (workspace UUID preferred; slug also accepted).  
The **inbox** endpoint uses the workspace **UUID** in the URL path.

---

## Quick Setup (GitHub Actions)

### 1. Add the Gateway Service

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
        env:
          DATABASE_URL: postgres://gateway:gateway@postgres:5432/gateway
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: gateway
          POSTGRES_PASSWORD: gateway
          POSTGRES_DB: gateway
```

### 2. Bootstrap: Create Workspace and API Key

Before bootstrapping, ensure that the base roles and permissions seed file `database/seeds/001_seed_permissions.sql` has been executed against the database. Workspace creation programmatically maps the creator to the `admin` role and registers permissions in the `gogate` cache; this process fails if the roles have not been seeded.

After applying the base seed, execute the bootstrap sequence:

```bash
BASE=http://localhost:10101/api/v1
UK="ci-${{ github.run_id }}"   # unique workspace key per run

# 1) Register (idempotent in CI) then login
curl -fsS -X POST "$BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"ci-$UK@e2e.test\",\"password\":\"ci-secret-123\",\"first_name\":\"CI\",\"last_name\":\"Runner\"}" >/dev/null 2>&1 || true

PORTAL_JWT=$(curl -sS -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"ci-$UK@e2e.test\",\"password\":\"ci-secret-123\"}" \
  | jq -r .token)

# 2) Create workspace (creator becomes admin member automatically)
WS_JSON=$(curl -sS -X POST "$BASE/workspaces" \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"CI\",\"unique_key\":\"$UK\"}")
WORKSPACE_ID=$(echo "$WS_JSON" | jq -r .id)

# 3) Create API key (secret shown once — save it)
KEY_JSON=$(curl -sS -X POST "$BASE/workspaces/$WORKSPACE_ID/api-keys" \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "Content-Type: application/json" \
  -d '{"name":"ci-key"}')
API_CLIENT_ID=$(echo "$KEY_JSON"   | jq -r .client_id)
API_CLIENT_SECRET=$(echo "$KEY_JSON" | jq -r .client_secret)

# 4) Dispatch defaults to memory — nothing to configure
echo "PORTAL_JWT=$PORTAL_JWT"           >> $GITHUB_ENV
echo "WORKSPACE_ID=$WORKSPACE_ID"      >> $GITHUB_ENV
echo "WORKSPACE_UK=$UK"                >> $GITHUB_ENV
echo "API_CLIENT_ID=$API_CLIENT_ID"    >> $GITHUB_ENV
echo "API_CLIENT_SECRET=$API_CLIENT_SECRET" >> $GITHUB_ENV
```

### 3. Configure Your App

Point your app's message sending to the gateway:

```bash
EMAIL_GATEWAY_URL=http://localhost:10101
EMAIL_GATEWAY_WORKSPACE=ci-${{ github.run_id }}   # = WORKSPACE_UK
EMAIL_GATEWAY_CLIENT_ID=${{ env.API_CLIENT_ID }}
EMAIL_GATEWAY_CLIENT_SECRET=${{ env.API_CLIENT_SECRET }}
```

Your app sends to:
```
POST http://localhost:10101/v1/email
Authorization: Basic <base64(client_id:secret)>
X-Workspace-Key: ci-<run_id>
```

### 4. Query and Assert Captured Messages

```bash
BASE=http://localhost:10101/api/v1

# List emails
curl -sS \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "X-Api-Client-Id: $API_CLIENT_ID" \
  -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
  "$BASE/workspaces/$WORKSPACE_ID/inbox/emails"
```

Clear between test cases:
```bash
curl -sS -X DELETE \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "X-Api-Client-Id: $API_CLIENT_ID" \
  -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
  "$BASE/workspaces/$WORKSPACE_ID/inbox/messages"
```

---

## Complete GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: gateway
          POSTGRES_PASSWORD: gateway
          POSTGRES_DB: gateway
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

      gateway:
        image: ghcr.io/weprodev/wpd-message-gateway:latest
        ports:
          - 10101:10101
        env:
          DATABASE_URL: postgres://gateway:gateway@postgres:5432/gateway
          MESSAGE_JWT_SECRET: ci-test-jwt-secret-32-chars-minimum
        depends_on:
          postgres:
            condition: service_healthy

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Bootstrap gateway workspace
        run: |
          BASE=http://localhost:10101/api/v1
          UK="ci-${{ github.run_id }}"

          # Register then login to get JWT
          curl -fsS -X POST "$BASE/auth/register" \
            -H "Content-Type: application/json" \
            -d "{\"email\":\"ci-$UK@e2e.test\",\"password\":\"ci-secret-123\",\"first_name\":\"CI\",\"last_name\":\"Runner\"}" >/dev/null 2>&1 || true
            
          TOKEN=$(curl -sS -X POST "$BASE/auth/login" \
            -H "Content-Type: application/json" \
            -d "{\"email\":\"ci-$UK@e2e.test\",\"password\":\"ci-secret-123\"}" | jq -r .token)

          WS=$(curl -sS -X POST "$BASE/workspaces" \
            -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
            -d "{\"name\":\"CI\",\"unique_key\":\"$UK\"}")
          WID=$(echo "$WS" | jq -r .id)

          KEY=$(curl -sS -X POST "$BASE/workspaces/$WID/api-keys" \
            -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
            -d '{"name":"ci"}')

          echo "PORTAL_JWT=$TOKEN"                                 >> $GITHUB_ENV
          echo "WORKSPACE_ID=$WID"                                >> $GITHUB_ENV
          echo "WORKSPACE_UK=$UK"                                 >> $GITHUB_ENV
          echo "API_CLIENT_ID=$(echo "$KEY" | jq -r .client_id)" >> $GITHUB_ENV
          echo "API_CLIENT_SECRET=$(echo "$KEY" | jq -r .client_secret)" >> $GITHUB_ENV

      - run: npm ci

      - name: Start app
        run: |
          EMAIL_GATEWAY_URL=http://localhost:10101 \
          EMAIL_GATEWAY_WORKSPACE=${{ env.WORKSPACE_UK }} \
          EMAIL_GATEWAY_CLIENT_ID=${{ env.API_CLIENT_ID }} \
          EMAIL_GATEWAY_CLIENT_SECRET=${{ env.API_CLIENT_SECRET }} \
          npm start &
          sleep 5

      - name: Test signup flow — verify welcome email
        env:
          BASE: http://localhost:10101/api/v1
        run: |
          # Clear inbox
          curl -sS -X DELETE \
            -H "Authorization: Bearer $PORTAL_JWT" \
            -H "X-Api-Client-Id: $API_CLIENT_ID" \
            -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
            "$BASE/workspaces/$WORKSPACE_ID/inbox/messages"

          # Trigger your app's signup
          curl -sS -X POST http://localhost:3000/api/signup \
            -H "Content-Type: application/json" \
            -d '{"email": "user@example.com", "name": "John"}'

          sleep 2

          # Assert email was sent with correct content
          curl -sS \
            -H "Authorization: Bearer $PORTAL_JWT" \
            -H "X-Api-Client-Id: $API_CLIENT_ID" \
            -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
            "$BASE/workspaces/$WORKSPACE_ID/inbox/emails" \
          | jq -e '
            length >= 1 and
            .[0].email.to[0] == "user@example.com" and
            .[0].email.subject == "Welcome!"
          ' || exit 1

          echo "✅ Welcome email verified"
```

---

## Go SDK: Testing with Memory Provider

If your app uses the gateway Go SDK, switch to `memory` provider in tests — no server needed:

```go
// main.go / real app
gw, _ := gateway.New(gateway.Config{
    DefaultEmailProvider: os.Getenv("EMAIL_PROVIDER"), // "mailgun" in prod
    EmailProviders: map[string]gateway.EmailConfig{
        "mailgun": { /* ... */ },
    },
})

// xxx_test.go — override with memory provider
func TestSignupEmail(t *testing.T) {
    gw, _ := gateway.New(gateway.Config{
        DefaultEmailProvider: "memory",
    })
    // Use gw — messages captured in RAM
    // Assert via the inbox HTTP API (if server running)
    // OR use the memory store directly in unit tests
}
```

---

## API Reference

### Send (Gateway API)

```bash
curl -X POST http://localhost:10101/v1/email \
  -u "client_id:client_secret" \
  -H "X-Workspace-Key: <workspace-uuid>" \
  -H "Content-Type: application/json" \
  -d '{"to":["user@example.com"],"subject":"Hello","html":"<p>World</p>"}'
```

### Inbox Queries

All paths: `/api/v1/workspaces/{uuid}/inbox/...`  
Required headers: `Authorization: Bearer <portal-jwt>`, `X-Api-Client-Id`, `X-Api-Client-Secret`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `.../inbox/stats` | GET | Message counts per channel |
| `.../inbox/emails` | GET | List captured emails |
| `.../inbox/emails/{id}` | GET | Get one email |
| `.../inbox/emails/{id}` | DELETE | Delete one email |
| `.../inbox/sms` | GET | List captured SMS |
| `.../inbox/push` | GET | List captured push notifications |
| `.../inbox/chat` | GET | List captured chat messages |
| `.../inbox/messages` | DELETE | Clear **all** captured messages |
| `.../inbox/events` | GET | SSE stream for real-time updates |

### Response Format

List endpoints return a JSON **array** of stored messages (not wrapped in an object):

```json
[
  {
    "id": "abc123",
    "created_at": "2024-01-15T10:30:00Z",
    "workspace_id": "uuid-of-workspace",
    "email": {
      "to": ["user@example.com"],
      "subject": "Welcome!",
      "html": "<h1>Hello John</h1>",
      "plain_text": "Hello John"
    }
  }
]
```

---

## Common Testing Patterns

### Assert Email Count

```bash
COUNT=$(curl -sS \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "X-Api-Client-Id: $API_CLIENT_ID" \
  -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
  "$BASE/workspaces/$WORKSPACE_ID/inbox/emails" | jq 'length')
[ "$COUNT" -eq 1 ] || { echo "Expected 1 email, got $COUNT"; exit 1; }
```

### Assert Email Content

```bash
curl -sS \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "X-Api-Client-Id: $API_CLIENT_ID" \
  -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
  "$BASE/workspaces/$WORKSPACE_ID/inbox/emails" | jq -e '
  .[0].email.subject == "Welcome!" and
  .[0].email.to[0] == "user@example.com"
'
```

### Wait for Message (polling)

```bash
for i in $(seq 1 10); do
  COUNT=$(curl -sS -H "Authorization: Bearer $PORTAL_JWT" \
    -H "X-Api-Client-Id: $API_CLIENT_ID" -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
    "$BASE/workspaces/$WORKSPACE_ID/inbox/emails" | jq 'length')
  [ "$COUNT" -ge 1 ] && break
  sleep 1
done
```

### Real-Time Updates (SSE)

```bash
curl -N \
  -H "Authorization: Bearer $PORTAL_JWT" \
  -H "X-Api-Client-Id: $API_CLIENT_ID" \
  -H "X-Api-Client-Secret: $API_CLIENT_SECRET" \
  "$BASE/workspaces/$WORKSPACE_ID/inbox/events"
```

---

## Local Development

```bash
# Run gateway locally
docker run -p 10101:10101 \
  -e DATABASE_URL=postgres://... \
  ghcr.io/weprodev/wpd-message-gateway:latest
```

Or from source:
```bash
git clone https://github.com/weprodev/wpd-message-gateway.git
cd wpd-message-gateway
make start  # starts server + Portal UI at http://localhost:10104
```

---

## Related Documentation

- [Usage](./usage.md) — Full API reference
- [Architecture](./architecture.md) — System design
- [Portal inbox](./portal-inbox.md) — Inbox API details
