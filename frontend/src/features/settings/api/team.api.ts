import { apiFetch } from "@/core/api/client"

import type {
  CreateInvitationResult,
  InvitableRole,
  WorkspaceInvitation,
  WorkspaceMember,
} from "../team.types"

async function readApiErrorMessage(res: Response, fallback: string): Promise<string> {
  try {
    const body = (await res.json()) as { message?: string }
    const message = body.message?.trim()
    return message || fallback
  } catch {
    return fallback
  }
}

function ensureArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : []
}

export async function listMembers(workspaceId: string): Promise<WorkspaceMember[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/members`)
  if (!res.ok) {
    throw new Error(await readApiErrorMessage(res, "Failed to load team members"))
  }
  return ensureArray(await res.json())
}

export async function removeMember(workspaceId: string, userId: string): Promise<void> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/members/${userId}`, {
    method: "DELETE",
  })
  if (!res.ok) {
    throw new Error(await readApiErrorMessage(res, "Failed to remove team member"))
  }
}

export async function listInvitations(workspaceId: string): Promise<WorkspaceInvitation[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/invitations`)
  if (!res.ok) {
    throw new Error(await readApiErrorMessage(res, "Failed to load pending invitations"))
  }
  return ensureArray(await res.json())
}

export async function createInvitation(
  workspaceId: string,
  email: string,
  role: InvitableRole,
): Promise<CreateInvitationResult> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/invitations`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, role }),
  })
  if (!res.ok) {
    throw new Error(await readApiErrorMessage(res, "Failed to send invitation"))
  }
  return (await res.json()) as CreateInvitationResult
}
