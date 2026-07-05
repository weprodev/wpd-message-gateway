export type WorkspaceAuthorizationScope = {
  role?: string
  permissions?: string[]
}

export type WorkspaceAuthorization = {
  role: string
  permissions: readonly string[]
  /** Current user holds one of the given workspace roles. */
  hasRole: (...roles: string[]) => boolean
  /** Current user has the permission (resolved from API / wpd-gogate). */
  can: (permission: string) => boolean
  canAny: (...permissions: string[]) => boolean
  canAll: (...permissions: string[]) => boolean
}

export function createWorkspaceAuthorization(
  scope: WorkspaceAuthorizationScope,
): WorkspaceAuthorization {
  const role = scope.role ?? ""
  const permissions = scope.permissions ?? []
  const permissionSet = new Set(permissions)

  return {
    role,
    permissions,
    hasRole: (...roles: string[]) => role !== "" && roles.includes(role),
    can: (permission: string) => permissionSet.has(permission),
    canAny: (...requested) => requested.some((permission) => permissionSet.has(permission)),
    canAll: (...requested) => requested.every((permission) => permissionSet.has(permission)),
  }
}
