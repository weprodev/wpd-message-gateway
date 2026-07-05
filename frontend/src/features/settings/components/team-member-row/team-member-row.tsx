import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Permission, useWorkspaceAuthorization } from "@/core/auth"

import { formatJoinedDate, formatMemberName, roleLabel } from "../../team.utils"
import type { WorkspaceMember } from "../../team.types"

interface TeamMemberRowProps {
  member: WorkspaceMember
  isCurrentUser: boolean
  isBusy?: boolean
  onRemove: (userId: string) => void
}

export function TeamMemberRow({ member, isCurrentUser, isBusy, onRemove }: TeamMemberRowProps) {
  const { can } = useWorkspaceAuthorization()
  const canRemove = can(Permission.MembersWrite) && !isCurrentUser

  return (
    <div className="flex items-center justify-between gap-4 border-b border-border px-4 py-4 last:border-b-0">
      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-2">
          <p className="text-sm font-medium text-foreground">{formatMemberName(member)}</p>
          <Badge variant="secondary">{roleLabel(member.role)}</Badge>
          {isCurrentUser ? (
            <span className="text-xs text-text-tertiary">You</span>
          ) : null}
        </div>
        {member.user_email ? (
          <p className="mt-1 truncate text-xs text-text-secondary">{member.user_email}</p>
        ) : null}
        <p className="mt-1 text-xs text-text-tertiary">
          Joined {formatJoinedDate(member.joined_at)}
        </p>
      </div>

      {canRemove ? (
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={isBusy}
          onClick={() => onRemove(member.user_id)}
        >
          <Icon name="person_remove" size="sm" />
          Remove
        </Button>
      ) : null}
    </div>
  )
}
