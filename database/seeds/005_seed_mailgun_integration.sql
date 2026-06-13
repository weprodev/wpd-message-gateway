-- Mailgun integration stub for the demo workspace (disconnected, empty config).

INSERT INTO integrations (id, workspace_id, channel_type, provider_id, encrypted_config, status, is_default)
VALUES (
    '00000000-0000-0000-0000-000000000005',
    '00000000-0000-0000-0000-000000000001',
    'email',
    '00000000-0000-0000-0000-000000000031',
    E'\\x1732ee13b2f560c579c73d29ed361c5ccd9eeb785369f5e934ddafdee7af',
    'disconnected',
    false
)
ON CONFLICT (workspace_id, provider_id) DO NOTHING;
