/** Mirrors internal/core/domain/permission.go — keep in sync with backend RBAC seeds. */

export const Role = {
  Admin: "admin",
  Member: "member",
  Viewer: "viewer",
} as const

export type WorkspaceRole = (typeof Role)[keyof typeof Role]

export const WorkspaceRoles: readonly WorkspaceRole[] = [
  Role.Admin,
  Role.Member,
  Role.Viewer,
]

export const Permission = {
  WorkspacesRead: "workspaces.read",
  WorkspacesWrite: "workspaces.write",
  MembersRead: "members.read",
  MembersWrite: "members.write",
  APIKeysRead: "apikeys.read",
  APIKeysWrite: "apikeys.write",
  LogsRead: "logs.read",
  IntegrationsRead: "integrations.read",
  IntegrationsWrite: "integrations.write",
  TemplatesRead: "templates.read",
  TemplatesWrite: "templates.write",
  SettingsRead: "settings.read",
  SettingsWrite: "settings.write",
  InvitationsRead: "invitations.read",
  InvitationsWrite: "invitations.write",
  InboxWrite: "inbox.write",
} as const

export type PermissionName = (typeof Permission)[keyof typeof Permission]

export const AllPermissions: readonly PermissionName[] = Object.values(Permission)

const memberPermissions = new Set<PermissionName>([
  Permission.WorkspacesRead,
  Permission.MembersRead,
  Permission.APIKeysRead,
  Permission.APIKeysWrite,
  Permission.LogsRead,
  Permission.IntegrationsRead,
  Permission.IntegrationsWrite,
  Permission.TemplatesRead,
  Permission.TemplatesWrite,
  Permission.SettingsRead,
  Permission.SettingsWrite,
  Permission.InvitationsRead,
  Permission.InboxWrite,
])

const viewerPermissions = new Set<PermissionName>([
  Permission.WorkspacesRead,
  Permission.MembersRead,
  Permission.APIKeysRead,
  Permission.LogsRead,
  Permission.IntegrationsRead,
  Permission.TemplatesRead,
  Permission.SettingsRead,
  Permission.InvitationsRead,
])

/** Static role → permission matrix (wpd-gogate HasRolePermission equivalent). */
export function hasRolePermission(role: string, permission: string): boolean {
  if (role === Role.Admin) {
    return (AllPermissions as readonly string[]).includes(permission)
  }
  if (role === Role.Member) {
    return memberPermissions.has(permission as PermissionName)
  }
  if (role === Role.Viewer) {
    return viewerPermissions.has(permission as PermissionName)
  }
  return false
}

export function isWorkspaceRole(role: string): role is WorkspaceRole {
  return (WorkspaceRoles as readonly string[]).includes(role)
}
