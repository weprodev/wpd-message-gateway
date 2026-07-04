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

// ReadPermissions are granted to the viewer workspace role (full read-only member).
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

// WritePermissions require admin or member role (never viewer).
var WritePermissions = []string{
	PermissionWorkspacesWrite,
	PermissionMembersWrite,
	PermissionAPIKeysWrite,
	PermissionIntegrationsWrite,
	PermissionTemplatesWrite,
	PermissionSettingsWrite,
	PermissionInvitationsWrite,
	PermissionInboxWrite,
}

// AllPermissions is the full portal RBAC catalog (read + write).
var AllPermissions = append(append([]string{}, ReadPermissions...), WritePermissions...)

// IsWritePermission reports whether name is a mutating portal permission.
func IsWritePermission(name string) bool {
	for _, permission := range WritePermissions {
		if permission == name {
			return true
		}
	}
	return false
}

// IsReadPermission reports whether name is a read-only portal permission for workspace members.
func IsReadPermission(name string) bool {
	for _, permission := range ReadPermissions {
		if permission == name {
			return true
		}
	}
	return false
}

// PublicGuestPermissions are the only permissions non-members receive on public workspaces.
// Narrower than viewer role: guests must not read members, API keys, logs, or settings.
var PublicGuestPermissions = []string{
	PermissionWorkspacesRead,
	PermissionTemplatesRead,
}

// IsPublicGuestPermission reports whether permission is allowed for public workspace guests.
func IsPublicGuestPermission(name string) bool {
	for _, permission := range PublicGuestPermissions {
		if permission == name {
			return true
		}
	}
	return false
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
