# Bruno API Tests

API tests for WPD Message Gateway using [Bruno](https://www.usebruno.com/).

## Setup

1. Install Bruno: https://www.usebruno.com/downloads
2. Open Bruno and import this folder (`tests/bruno`)

## Environments

| Environment | URL | Use Case |
|-------------|-----|----------|
| **local** | `http://localhost:10101` | Local development |
| **memory** | `http://localhost:10101` | Memory provider testing |

## Running Tests

### Start the Server

```bash
make start
```

### Run in Bruno

1. Open Bruno
2. Select environment: `local` or `memory`
3. Run individual requests or entire folders

### Run from CLI

```bash
# Install Bruno CLI
npm install -g @usebruno/cli

# Run all tests
cd tests/bruno
bru run --env local
```

## Test Structure

```
tests/bruno/
├── environments/
│   ├── local.bru        # Local dev environment
│   └── memory.bru       # Memory provider environment
├── Gateway/             # Gateway API (/v1/*)
│   ├── Send Email.bru
│   ├── Get Inbox.bru
│   └── Clear Inbox.bru
└── DevBox/              # DevBox API (/api/v1/*)
    ├── Get Stats.bru
    ├── Clear All Messages.bru
    ├── Email/
    │   ├── Ingest Email.bru
    │   ├── List Emails.bru
    │   ├── Get Email by ID.bru
    │   └── Delete Email.bru
    ├── SMS/
    │   ├── Ingest SMS.bru
    │   ├── List SMS.bru
    │   ├── Get SMS by ID.bru
    │   └── Delete SMS.bru
    ├── Push/
    │   ├── Ingest Push.bru
    │   ├── List Push.bru
    │   ├── Get Push by ID.bru
    │   └── Delete Push.bru
    └── Chat/
        ├── Ingest Chat.bru
        ├── List Chat.bru
        ├── Get Chat by ID.bru
        └── Delete Chat.bru
```

## Typical Test Flow

1. **Clear All Messages** - Start fresh
2. **Ingest Email/SMS/Push/Chat** - Create test messages
3. **List Messages** - Verify they appear
4. **Get by ID** - Verify individual retrieval
5. **Delete** - Verify deletion works
6. **Get Stats** - Verify counts are correct
