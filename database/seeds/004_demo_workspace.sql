-- Demo workspace (idempotent — safe to re-run via init-db.sh, run_migrations.sh, or make seed-demo).
-- Requires 001_seed_permissions.sql and 002_seed_providers.sql first.
--
-- Portal accounts (password for all: secret)
--   demo@weprodev.com   → admin  on Demo Workspace
--   member@weprodev.com → member on Demo Workspace
--   viewer@weprodev.com → viewer on Demo Workspace
--
-- API key: demo-client-id / demo-secret
-- Workspace: demo (00000000-0000-0000-0000-000000000001)

-- Fixed demo UUIDs
--   workspace : 00000000-0000-0000-0000-000000000001
--   admin user: 00000000-0000-0000-0000-000000000010
--   member    : 00000000-0000-0000-0000-000000000011
--   viewer    : 00000000-0000-0000-0000-000000000012

-- bcrypt hash for password "secret" (cost 14)
-- Remove conflicting sign-ups so fixed UUIDs stay stable.
DELETE FROM users
WHERE email IN ('demo@weprodev.com', 'member@weprodev.com', 'viewer@weprodev.com')
  AND id NOT IN (
    '00000000-0000-0000-0000-000000000010',
    '00000000-0000-0000-0000-000000000011',
    '00000000-0000-0000-0000-000000000012'
  );

DELETE FROM workspaces
WHERE slug = 'demo'
  AND id <> '00000000-0000-0000-0000-000000000001';

INSERT INTO users (id, email, password_hash, first_name, last_name, email_verified)
VALUES
    (
        '00000000-0000-0000-0000-000000000010',
        'demo@weprodev.com',
        '$2a$14$fBL/4GbbqSCMTVOYotuUa.Qx0DwnRMpkZaOHxkD3h1X4gMESUhjD.',
        'Demo',
        'Admin',
        true
    ),
    (
        '00000000-0000-0000-0000-000000000011',
        'member@weprodev.com',
        '$2a$14$fBL/4GbbqSCMTVOYotuUa.Qx0DwnRMpkZaOHxkD3h1X4gMESUhjD.',
        'Demo',
        'Member',
        true
    ),
    (
        '00000000-0000-0000-0000-000000000012',
        'viewer@weprodev.com',
        '$2a$14$fBL/4GbbqSCMTVOYotuUa.Qx0DwnRMpkZaOHxkD3h1X4gMESUhjD.',
        'Demo',
        'Viewer',
        true
    )
ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    email_verified = EXCLUDED.email_verified;

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
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    slug = EXCLUDED.slug,
    owner_id = EXCLUDED.owner_id,
    status = EXCLUDED.status,
    is_private = EXCLUDED.is_private;

-- Workspace membership + wpd-gogate role assignments (admin, member, viewer)
INSERT INTO workspace_members (workspace_id, user_id, role_id)
SELECT '00000000-0000-0000-0000-000000000001', seed.user_id, r.id
FROM (
    VALUES
        ('00000000-0000-0000-0000-000000000010'::uuid, 'admin'),
        ('00000000-0000-0000-0000-000000000011'::uuid, 'member'),
        ('00000000-0000-0000-0000-000000000012'::uuid, 'viewer')
) AS seed(user_id, role_name)
JOIN roles r ON r.name = seed.role_name AND r.guard_name = 'msg_web'
ON CONFLICT (workspace_id, user_id) DO UPDATE SET
    role_id = EXCLUDED.role_id;

INSERT INTO model_has_roles (role_id, model_type, model_id, team_id)
SELECT r.id, 'users', seed.user_id, '00000000-0000-0000-0000-000000000001'
FROM (
    VALUES
        ('00000000-0000-0000-0000-000000000010'::uuid, 'admin'),
        ('00000000-0000-0000-0000-000000000011'::uuid, 'member'),
        ('00000000-0000-0000-0000-000000000012'::uuid, 'viewer')
) AS seed(user_id, role_name)
JOIN roles r ON r.name = seed.role_name AND r.guard_name = 'msg_web'
ON CONFLICT (role_id, model_id, model_type, team_id) DO NOTHING;

INSERT INTO workspace_channels (workspace_id, channel_type, enabled)
VALUES ('00000000-0000-0000-0000-000000000001', 'email', true)
ON CONFLICT (workspace_id, channel_type) DO UPDATE SET
    enabled = EXCLUDED.enabled;

INSERT INTO api_keys (id, workspace_id, client_id, client_secret_hash, name, is_active)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'demo-client-id',
    '$2a$14$0DAJjkBNrPj84vYy7dcYP.61t4BTzxE0ZvGbKjmuAviEPKWf9/Bee',
    'Demo API Key',
    true
)
ON CONFLICT (client_id) DO UPDATE SET
    client_secret_hash = EXCLUDED.client_secret_hash,
    name = EXCLUDED.name,
    is_active = EXCLUDED.is_active,
    workspace_id = EXCLUDED.workspace_id;

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
ON CONFLICT (workspace_id, provider_id) DO UPDATE SET
    encrypted_config = EXCLUDED.encrypted_config,
    status = EXCLUDED.status,
    is_default = EXCLUDED.is_default;
