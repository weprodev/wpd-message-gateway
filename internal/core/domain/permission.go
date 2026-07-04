package domain

// Roles in the RBAC system (global catalog; scoped per workspace via wpd-gogate team_id).
const (
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleViewer = "viewer"
)

// WorkspaceRoles lists assignable workspace roles seeded for every workspace.
var WorkspaceRoles = []string{RoleAdmin, RoleMember, RoleViewer}

// RBACGuardName is the guard_name used for portal roles and permissions in this project.
const RBACGuardName = "msg_web"

// Permissions in the RBAC system.
const (
	// Workspaces
	PermissionWorkspacesRead  = "workspaces.read"
	PermissionWorkspacesWrite = "workspaces.write"

	// Workspace Members
	PermissionMembersRead  = "members.read"
	PermissionMembersWrite = "members.write"

	// API Keys
	PermissionAPIKeysRead  = "apikeys.read"
	PermissionAPIKeysWrite = "apikeys.write"

	// Logs
	PermissionLogsRead = "logs.read"

	// Integrations
	PermissionIntegrationsRead  = "integrations.read"
	PermissionIntegrationsWrite = "integrations.write"

	// Templates
	PermissionTemplatesRead  = "templates.read"
	PermissionTemplatesWrite = "templates.write"

	// Settings
	PermissionSettingsRead  = "settings.read"
	PermissionSettingsWrite = "settings.write"

	// Invitations
	PermissionInvitationsRead  = "invitations.read"
	PermissionInvitationsWrite = "invitations.write"

	// Inbox
	PermissionInboxWrite = "inbox.write"
)

// ReadPermissions are granted to viewer role and public workspace guests.
var ReadPermissions = []string{
	PermissionWorkspacesRead,
	PermissionMembersRead,
	PermissionAPIKeysRead,
	PermissionLogsRead,
	PermissionIntegrationsRead,
	PermissionTemplatesRead,
	PermissionSettingsRead,
	PermissionInvitationsRead,
}

// IsWorkspaceRole reports whether name is one of the seeded workspace roles.
func IsWorkspaceRole(name string) bool {
	for _, role := range WorkspaceRoles {
		if role == name {
			return true
		}
	}
	return false
}
