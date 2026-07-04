import { describe, expect, it } from "vitest"

import { Permission, Role, hasRolePermission } from "./permissions"
import { createWorkspaceAuthorization } from "./authorization"

describe("createWorkspaceAuthorization", () => {
  it("checks role and permissions for the active workspace", () => {
    const auth = createWorkspaceAuthorization({
      role: Role.Member,
      permissions: [Permission.IntegrationsRead, Permission.IntegrationsWrite],
    })

    expect(auth.hasRole(Role.Member)).toBe(true)
    expect(auth.hasRole(Role.Admin)).toBe(false)
    expect(auth.can(Permission.IntegrationsWrite)).toBe(true)
    expect(auth.can(Permission.MembersWrite)).toBe(false)
    expect(auth.canAny(Permission.MembersWrite, Permission.IntegrationsRead)).toBe(true)
    expect(auth.canAll(Permission.IntegrationsRead, Permission.IntegrationsWrite)).toBe(true)
    expect(auth.canAll(Permission.IntegrationsRead, Permission.MembersWrite)).toBe(false)
  })

  it("denies everything when scope is empty", () => {
    const auth = createWorkspaceAuthorization({})

    expect(auth.hasRole(Role.Admin)).toBe(false)
    expect(auth.can(Permission.SettingsRead)).toBe(false)
  })
})

describe("hasRolePermission", () => {
  it("matches seeded role matrices", () => {
    expect(hasRolePermission(Role.Admin, Permission.MembersWrite)).toBe(true)
    expect(hasRolePermission(Role.Member, Permission.MembersWrite)).toBe(false)
    expect(hasRolePermission(Role.Member, Permission.IntegrationsWrite)).toBe(true)
    expect(hasRolePermission(Role.Viewer, Permission.SettingsRead)).toBe(true)
    expect(hasRolePermission(Role.Viewer, Permission.SettingsWrite)).toBe(false)
  })
})
