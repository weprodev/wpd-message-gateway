-- Demo workspace, API key, email integration, template (aligned with latest schema)

INSERT INTO users (id, email, password, first_name, last_name, status)
VALUES (
    '00000000-0000-0000-0000-000000000010',
    'demo@weprodev.com',
    '$2a$14$dummy.hash.for.demo.only',
    'Demo',
    'User',
    'active'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO roles (id, name, display_name)
VALUES
    ('00000000-0000-0000-0000-000000000020', 'admin', 'Admin'),
    ('00000000-0000-0000-0000-000000000021', 'member', 'Member')
ON CONFLICT (id) DO NOTHING;

INSERT INTO workspaces (id, name, slug, owner_id, status, is_private, hashed_pin_code, icon_key)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Demo Workspace',
    'demo',
    '00000000-0000-0000-0000-000000000010',
    'active',
    true,
    NULL,
    NULL
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO workspace_members (workspace_id, user_id, role_id)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000020'
)
ON CONFLICT (workspace_id, user_id) DO NOTHING;

INSERT INTO workspace_channels (workspace_id, channel_type, enabled)
VALUES ('00000000-0000-0000-0000-000000000001', 'email', true)
ON CONFLICT (workspace_id, channel_type) DO NOTHING;

INSERT INTO api_keys (id, workspace_id, client_id, client_secret_hash, name, is_active)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'demo-client-id',
    '$2a$14$0DAJjkBNrPj84vYy7dcYP.61t4BTzxE0ZvGbKjmuAviEPKWf9/Bee',
    'Demo API Key',
    true
)
ON CONFLICT (client_id) DO NOTHING;

INSERT INTO providers (id, name, channel_type, status)
VALUES (
    '00000000-0000-0000-0000-000000000030',
    'memory',
    'email',
    'active'
)
ON CONFLICT (id) DO NOTHING;

-- Encrypted JSON "{}" for memory provider (AES-GCM, dev only)
INSERT INTO integrations (id, workspace_id, channel_type, provider_id, encrypted_config, status, is_default)
VALUES (
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000001',
    'email',
    '00000000-0000-0000-0000-000000000030',
    E'\\x720231d45885b8d808a3e2930264f2b60ca166004abf449b6be50328217e',
    'connected',
    true
)
ON CONFLICT (workspace_id, provider_id) DO NOTHING;

INSERT INTO templates (id, workspace_id, name, unique_key, channel_type, subject, content, category, is_active, is_default)
VALUES (
    '00000000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0000-000000000001',
    'Welcome Email',
    'welcome_email',
    'email',
    'Welcome',
    '<h1>Welcome to WeProDev!</h1><p>Your demo workspace is active.</p>',
    'transactional',
    true,
    true
)
ON CONFLICT (workspace_id, unique_key) DO NOTHING;
