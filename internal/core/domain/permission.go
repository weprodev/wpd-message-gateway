package domain

// Roles in the RBAC system.
const (
	RoleAdmin  = "admin"
	RoleMember = "member"
)

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
