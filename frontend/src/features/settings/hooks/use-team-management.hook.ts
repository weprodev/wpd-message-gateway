import { useEffect, useState } from "react"

import { fetchUserProfile } from "@/core/api/client"

import {
  createInvitation,
  listInvitations,
  listMembers,
  removeMember,
} from "../api/team.api"
import type { CreateInvitationResult, InvitableRole, WorkspaceInvitation, WorkspaceMember } from "../team.types"

export function useTeamManagement(workspaceId: string, enabled: boolean) {
  const [members, setMembers] = useState<WorkspaceMember[]>([])
  const [invitations, setInvitations] = useState<WorkspaceInvitation[]>([])
  const [currentUserId, setCurrentUserId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [busyUserId, setBusyUserId] = useState<string | null>(null)
  const [trigger, setTrigger] = useState(0)

  const reload = () => setTrigger((value) => value + 1)

  useEffect(() => {
    if (!enabled || !workspaceId) return

    let cancelled = false

    ;(async () => {
      setIsLoading(true)
      setError(null)
      try {
        const [membersData, invitationsData, profile] = await Promise.all([
          listMembers(workspaceId),
          listInvitations(workspaceId),
          fetchUserProfile(),
        ])
        if (cancelled) return
        setMembers(membersData ?? [])
        setInvitations(invitationsData ?? [])
        setCurrentUserId(profile?.id ?? null)
      } catch (err) {
        if (cancelled) return
        setError(err instanceof Error ? err.message : "Failed to load team")
      } finally {
        if (!cancelled) setIsLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [workspaceId, enabled, trigger])

  async function inviteMember(email: string, role: InvitableRole): Promise<CreateInvitationResult> {
    const created = await createInvitation(workspaceId, email, role)
    setInvitations((prev) => [
      {
        id: created.id,
        workspace_id: workspaceId,
        email: created.email,
        role: created.role,
        expires_at: created.expires_at,
        status: "pending",
        created_at: new Date().toISOString(),
      },
      ...prev,
    ])
    setError(null)
    return created
  }

  async function removeTeamMember(userId: string) {
    setBusyUserId(userId)
    try {
      await removeMember(workspaceId, userId)
      setMembers((prev) => prev.filter((member) => member.user_id !== userId))
      setError(null)
    } finally {
      setBusyUserId(null)
    }
  }

  return {
    members,
    invitations,
    currentUserId,
    isLoading,
    error,
    busyUserId,
    reload,
    inviteMember,
    removeTeamMember,
  }
}
