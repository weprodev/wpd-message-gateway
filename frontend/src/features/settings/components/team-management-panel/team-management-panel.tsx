import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Spinner } from "@/components/ui/spinner"
import { Can, Permission } from "@/core/auth"

import { useTeamManagement } from "../../hooks/use-team-management.hook"
import { InvitationRow } from "../invitation-row"
import { InvitationTokenDialog } from "../invitation-token-dialog"
import { InviteMemberDialog } from "../invite-member-dialog"
import { TeamMemberRow } from "../team-member-row"

interface TeamManagementPanelProps {
  workspaceId: string
  enabled: boolean
}

export function TeamManagementPanel({ workspaceId, enabled }: TeamManagementPanelProps) {
  const {
    members,
    invitations,
    currentUserId,
    isLoading,
    error,
    busyUserId,
    inviteMember,
    removeTeamMember,
  } = useTeamManagement(workspaceId, enabled)

  const [inviteDialogOpen, setInviteDialogOpen] = useState(false)
  const [invitationToken, setInvitationToken] = useState<{
    email: string
    role: string
    token: string
  } | null>(null)

  async function handleInvite(email: string, role: Parameters<typeof inviteMember>[1]) {
    const created = await inviteMember(email, role)
    setInvitationToken({
      email: created.email,
      role: created.role,
      token: created.token,
    })
    return created
  }

  async function handleRemoveMember(userId: string) {
    if (!window.confirm("Remove this member from the workspace?")) return
    await removeTeamMember(userId)
  }

  if (isLoading) {
    return (
      <div className="flex min-h-[240px] items-center justify-center gap-3">
        <Spinner size="lg" />
        <span className="text-sm text-text-secondary">Loading team…</span>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-8">
      {error ? (
        <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </p>
      ) : null}

      <section className="flex flex-col gap-4">
        <div className="flex items-center justify-between gap-4">
          <div>
            <h2 className="text-base font-semibold text-foreground">Members</h2>
            <p className="text-sm text-text-secondary">
              People with access to this workspace and their assigned roles.
            </p>
          </div>
          <Can permission={Permission.InvitationsWrite}>
            <Button type="button" onClick={() => setInviteDialogOpen(true)}>
              <Icon name="person_add" size="sm" />
              Invite member
            </Button>
          </Can>
        </div>

        <div className="overflow-hidden rounded-lg border border-border bg-card">
          {(members ?? []).length === 0 ? (
            <p className="px-4 py-8 text-center text-sm text-text-secondary">
              No members yet.
            </p>
          ) : (
            (members ?? []).map((member) => (
              <TeamMemberRow
                key={member.user_id}
                member={member}
                isCurrentUser={member.user_id === currentUserId}
                isBusy={busyUserId === member.user_id}
                onRemove={handleRemoveMember}
              />
            ))
          )}
        </div>
      </section>

      <Can permission={Permission.InvitationsRead}>
        <section className="flex flex-col gap-4">
          <div>
            <h2 className="text-base font-semibold text-foreground">Pending invitations</h2>
            <p className="text-sm text-text-secondary">
              Invitations waiting to be accepted.
            </p>
          </div>

          <div className="overflow-hidden rounded-lg border border-border bg-card">
            {(invitations ?? []).length === 0 ? (
              <p className="px-4 py-8 text-center text-sm text-text-secondary">
                No pending invitations.
              </p>
            ) : (
              (invitations ?? []).map((invitation) => (
                <InvitationRow key={invitation.id} invitation={invitation} />
              ))
            )}
          </div>
        </section>
      </Can>

      <InviteMemberDialog
        key={inviteDialogOpen ? "open" : "closed"}
        open={inviteDialogOpen}
        onClose={() => setInviteDialogOpen(false)}
        onInvite={handleInvite}
        members={members ?? []}
        invitations={invitations ?? []}
      />

      <InvitationTokenDialog
        open={invitationToken !== null}
        email={invitationToken?.email ?? ""}
        role={invitationToken?.role ?? ""}
        token={invitationToken?.token ?? ""}
        onClose={() => setInvitationToken(null)}
      />
    </div>
  )
}
