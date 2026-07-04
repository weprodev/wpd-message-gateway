import type { WorkspaceRole } from "@/core/auth"

export interface WorkspaceMember {
  workspace_id: string
  user_id: string
  role: string
  joined_at: string
  user_email?: string
  display_name?: string
}

export interface WorkspaceInvitation {
  id: string
  workspace_id: string
  email: string
  role: string
  expires_at: string
  status: string
  created_at: string
}

export interface CreateInvitationResult {
  id: string
  email: string
  role: string
  expires_at: string
  token: string
}

export type InvitableRole = WorkspaceRole
