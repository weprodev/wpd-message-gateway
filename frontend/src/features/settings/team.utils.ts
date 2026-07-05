import { Role, WorkspaceRoles, type WorkspaceRole } from "@/core/auth"

const ROLE_LABELS: Record<WorkspaceRole, string> = {
  [Role.Admin]: "Admin",
  [Role.Member]: "Member",
  [Role.Viewer]: "Viewer",
}

export const INVITABLE_ROLES: readonly WorkspaceRole[] = WorkspaceRoles

export function roleLabel(role: string): string {
  if (role in ROLE_LABELS) {
    return ROLE_LABELS[role as WorkspaceRole]
  }
  return role
}

export function formatMemberName(member: { display_name?: string; user_email?: string }): string {
  const name = member.display_name?.trim()
  if (name) return name
  return member.user_email ?? "Unknown member"
}

export function formatJoinedDate(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString()
}

export function formatExpiryDate(value: string): string {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
}

export function normalizeEmail(email: string): string {
  return email.trim().toLowerCase()
}

export function isExistingMemberEmail(
  members: Array<{ user_email?: string }>,
  email: string,
): boolean {
  const normalized = normalizeEmail(email)
  return members.some(
    (member) => member.user_email && normalizeEmail(member.user_email) === normalized,
  )
}

export function hasPendingInvitationEmail(
  invitations: Array<{ email: string }>,
  email: string,
): boolean {
  const normalized = normalizeEmail(email)
  return invitations.some((invitation) => normalizeEmail(invitation.email) === normalized)
}
