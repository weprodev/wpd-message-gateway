-- Default RBAC roles, permissions, and mappings for WPD Message Gateway.
-- Schema and assignment model follow wpd-gogate (global roles; workspace scope via team_id).
-- Guard name must match domain.RBACGuardName and gogate.DefaultGuardName in wire.go.

-- 1. Global workspace roles (available to every workspace via model_has_roles.team_id)
INSERT INTO roles (name, guard_name)
VALUES
    ('admin', 'msg_web'),
    ('member', 'msg_web'),
    ('viewer', 'msg_web')
ON CONFLICT (name, guard_name) DO NOTHING;

-- 2. Portal permissions
INSERT INTO permissions (name, guard_name)
VALUES
    ('workspaces.read', 'msg_web'),
    ('workspaces.write', 'msg_web'),
    ('members.read', 'msg_web'),
    ('members.write', 'msg_web'),
    ('apikeys.read', 'msg_web'),
    ('apikeys.write', 'msg_web'),
    ('logs.read', 'msg_web'),
    ('integrations.read', 'msg_web'),
    ('integrations.write', 'msg_web'),
    ('templates.read', 'msg_web'),
    ('templates.write', 'msg_web'),
    ('settings.read', 'msg_web'),
    ('settings.write', 'msg_web'),
    ('invitations.read', 'msg_web'),
    ('invitations.write', 'msg_web'),
    ('inbox.write', 'msg_web')
ON CONFLICT (name, guard_name) DO NOTHING;

-- 3. Admin — full access
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.guard_name = 'msg_web' AND p.guard_name = 'msg_web'
ON CONFLICT (permission_id, role_id) DO NOTHING;

-- 4. Member — manage workspace resources; cannot change workspace, members, or invitations
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'member'
  AND r.guard_name = 'msg_web'
  AND p.guard_name = 'msg_web'
  AND p.name IN (
    'workspaces.read',
    'members.read',
    'apikeys.read',
    'apikeys.write',
    'logs.read',
    'integrations.read',
    'integrations.write',
    'templates.read',
    'templates.write',
    'settings.read',
    'settings.write',
    'invitations.read',
    'inbox.write'
  )
ON CONFLICT (permission_id, role_id) DO NOTHING;

-- 5. Viewer — read-only
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'viewer'
  AND r.guard_name = 'msg_web'
  AND p.guard_name = 'msg_web'
  AND p.name IN (
    'workspaces.read',
    'members.read',
    'apikeys.read',
    'logs.read',
    'integrations.read',
    'templates.read',
    'settings.read',
    'invitations.read'
  )
ON CONFLICT (permission_id, role_id) DO NOTHING;
