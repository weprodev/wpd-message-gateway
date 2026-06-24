# Quickstart: Data Retention Modes

## Prerequisites

- Gateway server running with PostgreSQL (`make dev` or equivalent)
- Workspace with API key
- Bruno or curl for send requests

## 1. Memory only (default)

1. Set retention to `memory` in Settings → Data Retention (or leave unset).
2. Send a test email via API.
3. **Expect**: Message in portal inbox (memory); `SELECT COUNT(*) FROM message_request_logs WHERE workspace_id = '<wid>'` unchanged.

## 2. Memory + Database

1. Set retention to `memory_database`; save.
2. Send a successful test message.
3. **Expect**: Inbox capture + new row in `message_request_logs`.
4. Send with invalid body (400).
5. **Expect**: No new `message_request_logs` row.

## 3. Provider only

1. Connect an integration; set retention to `provider`.
2. Send successfully.
3. **Expect**: Provider dispatch; no `message_request_logs` row.

## 4. Provider + Database

1. Set retention to `provider_database`; save.
2. Send successfully.
3. **Expect**: Same dispatch as provider only + new `message_request_logs` row with provider metadata.
4. Force dispatch failure (e.g. invalid provider config).
5. **Expect**: Error response; no new log row.

## 5. Legacy alias normalization

1. Manually set `workspace_settings.value = 'both'` for `data_retention`.
2. Reload Settings page.
3. **Expect**: UI shows **Memory + Database** selected.

## Verification commands

```bash
go test -race ./internal/core/domain/... ./internal/core/service/... ./internal/presentation/handler/...
cd frontend && npm run lint && npm run test
make audit
```

See [data-model.md](./data-model.md) and [contracts/retention-modes.md](./contracts/retention-modes.md) for enum mappings.
