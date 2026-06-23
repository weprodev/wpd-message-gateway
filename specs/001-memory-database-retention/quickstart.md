# Quickstart: Data Retention Modes & API Key Modals

## Prerequisites

- Gateway server running with PostgreSQL (`make run` or project quickstart)
- Portal frontend dev server (`cd frontend && npm run dev`)
- Bruno CLI for API verification (`tests/bruno/`)

## 1. Verify retention modes (Bruno / curl)

1. Set retention to **Memory only**:
   ```bash
   curl -X PATCH "$BASE/api/v1/workspaces/$WID/settings" \
     -H "Authorization: Bearer $JWT" \
     -H "Content-Type: application/json" \
     -d '{"data_retention":"memory"}'
   ```
2. Send a test message (any channel). **Expect**: no new row in `message_request_logs` for that request.
3. Set `data_retention` to `memory_database`, send again. **Expect**: request log row inserted.
4. Set `data_retention` to `provider` (requires active integration). **Expect**: provider dispatch, no request log.
5. Set `data_retention` to `provider_database`, send again. **Expect**: provider dispatch + request log row.

### SQL check

```sql
SELECT id, workspace_id, channel_type, status_code, created_at
FROM message_request_logs
WHERE workspace_id = '<wid>'
ORDER BY created_at DESC
LIMIT 5;
```

## 2. Legacy alias normalization

1. Manually insert legacy setting:
   ```sql
   INSERT INTO workspace_settings (workspace_id, key, value)
   VALUES ('<wid>', 'data_retention', 'both')
   ON CONFLICT (workspace_id, key) DO UPDATE SET value = EXCLUDED.value;
   ```
2. Open Portal → Settings → Data Retention. **Expect**: **Memory + Database** selected.
3. Save without changes. **Expect**: stored value becomes `memory_database`.

## 3. API key modals (manual UI)

1. Navigate to **Settings → Developer**.
2. **Create**: Click **Generate key** → enter name "Production" → **Generate Key**.
   - Credentials modal shows client ID + secret.
   - Copy buttons swap to checkmark with transition.
3. **Regenerate**: Click **Regenerate** on a row → **Cancel** (no change) → repeat → **Generate Key** → credentials modal.
4. **Delete**: Click **Delete** → **Cancel** (key remains) → **Delete** (key removed from list).

## 4. Automated tests

```bash
# Domain mapping tests
go test -race ./internal/core/domain/... -run DispatchMode

# Gateway dispatch + log gating
go test -race ./internal/core/service/... ./internal/presentation/handler/...

# Frontend
cd frontend && npm run lint && npm run test
```

## 5. Definition of done

- [ ] Four retention radios in portal including **Provider + Database**
- [ ] Request logs only for `memory_database` and `provider_database`
- [ ] No browser `prompt`/`confirm` on API key actions
- [ ] `make audit` passes

See [contracts/retention-modes.md](./contracts/retention-modes.md) and [data-model.md](./data-model.md) for canonical values.
