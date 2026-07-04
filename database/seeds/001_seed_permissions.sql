-- Default RBAC Permissions & Role Bindings for WPD Message Gateway

-- 1. Insert default roles
INSERT INTO roles (id, name, display_name, guard_name)
VALUES
    ('00000000-0000-0000-0000-000000000020', 'admin', 'Admin', 'msg_web'),
    ('00000000-0000-0000-0000-000000000021', 'member', 'Member', 'msg_web')
ON CONFLICT (name, guard_name) DO NOTHING;

-- 2. Insert permissions
INSERT INTO permissions (id, name, guard_name)
VALUES
    (gen_random_uuid(), 'workspaces.read', 'msg_web'),
    (gen_random_uuid(), 'workspaces.write', 'msg_web'),
    (gen_random_uuid(), 'members.read', 'msg_web'),
    (gen_random_uuid(), 'members.write', 'msg_web'),
    (gen_random_uuid(), 'apikeys.read', 'msg_web'),
    (gen_random_uuid(), 'apikeys.write', 'msg_web'),
    (gen_random_uuid(), 'logs.read', 'msg_web'),
    (gen_random_uuid(), 'integrations.read', 'msg_web'),
    (gen_random_uuid(), 'integrations.write', 'msg_web'),
    (gen_random_uuid(), 'templates.read', 'msg_web'),
    (gen_random_uuid(), 'templates.write', 'msg_web'),
    (gen_random_uuid(), 'settings.read', 'msg_web'),
    (gen_random_uuid(), 'settings.write', 'msg_web'),
    (gen_random_uuid(), 'invitations.read', 'msg_web'),
    (gen_random_uuid(), 'invitations.write', 'msg_web'),
    (gen_random_uuid(), 'inbox.write', 'msg_web')
ON CONFLICT (name, guard_name) DO NOTHING;

-- 3. Bind permissions to 'admin' role
INSERT INTO role_has_permissions (permission_id, role_id)
SELECT p.id, r.id
FROM permissions p
CROSS JOIN roles r
WHERE r.name = 'admin'
ON CONFLICT (permission_id, role_id) DO NOTHING;

-- 4. Bind permissions to 'member' role (read-only + inbox.write)
INSERT INTO role_has_permissions (permission_id, role_id)
SELECT p.id, r.id
FROM permissions p
CROSS JOIN roles r
WHERE r.name = 'member'
  AND p.name IN (
    'workspaces.read',
    'members.read',
    'apikeys.read',
    'logs.read',
    'integrations.read',
    'templates.read',
    'settings.read',
    'invitations.read',
    'inbox.write'
  )
ON CONFLICT (permission_id, role_id) DO NOTHING;
