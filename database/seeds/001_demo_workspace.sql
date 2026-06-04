-- Demo data for push/telegram (adapted to OTP-branch schema)

-- ── Channels: enable push for existing workspace ────────────────────────
INSERT INTO workspace_channels (workspace_id, channel_type, enabled)
VALUES ('00000000-0000-0000-0000-000000000001', 'push', true)
ON CONFLICT (workspace_id, channel_type) DO NOTHING;

-- ── Providers: add push providers to global catalog ─────────────────────
INSERT INTO providers (id, name, channel_type, status)
VALUES
    ('00000000-0000-0000-0000-000000000031', 'memory', 'push', 'active'),
    ('00000000-0000-0000-0000-000000000032', 'telegram', 'push', 'active')
ON CONFLICT (id) DO NOTHING;

-- ── Integrations: push/telegram with encrypted bot token ────────────────
-- Encrypted JSON: {"api_key":"REPLACE_WITH_YOUR_TELEGRAM_BOT_API_KEY"}
INSERT INTO integrations (id, workspace_id, channel_type, provider_id, provider_name, encrypted_config, status, is_default)
VALUES (
    '00000000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0000-000000000001',
    'push',
    '00000000-0000-0000-0000-000000000032',
    'telegram',
    E'\\x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
    'disconnected',
    true
)
ON CONFLICT (workspace_id, channel_type, provider_name) DO NOTHING;

-- ── Workspace settings: route push through provider ─────────────────────
INSERT INTO workspace_settings (workspace_id, key, value)
VALUES ('00000000-0000-0000-0000-000000000001', 'message_dispatch_mode', 'provider_only')
ON CONFLICT (workspace_id, key) DO UPDATE SET value = EXCLUDED.value;
