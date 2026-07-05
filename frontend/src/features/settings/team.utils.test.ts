import { describe, expect, it } from "vitest"

import { Role } from "@/core/auth"

import {
  hasPendingInvitationEmail,
  isExistingMemberEmail,
  roleLabel,
} from "./team.utils"

describe("roleLabel", () => {
  it("returns human-readable labels for workspace roles", () => {
    expect(roleLabel(Role.Admin)).toBe("Admin")
    expect(roleLabel(Role.Member)).toBe("Member")
    expect(roleLabel(Role.Viewer)).toBe("Viewer")
  })

  it("falls back to raw role name", () => {
    expect(roleLabel("custom")).toBe("custom")
  })
})

describe("invitation email guards", () => {
  it("detects existing workspace members by email", () => {
    expect(
      isExistingMemberEmail([{ user_email: "Member@WeProDev.com" }], "member@weprodev.com"),
    ).toBe(true)
    expect(isExistingMemberEmail([{ user_email: "other@example.com" }], "member@weprodev.com")).toBe(
      false,
    )
  })

  it("detects duplicate pending invitations by email", () => {
    expect(hasPendingInvitationEmail([{ email: "Pending@Example.com" }], "pending@example.com")).toBe(
      true,
    )
    expect(hasPendingInvitationEmail([], "pending@example.com")).toBe(false)
  })
})
