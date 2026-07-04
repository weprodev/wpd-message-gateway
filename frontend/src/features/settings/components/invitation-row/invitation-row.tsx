import { Badge } from "@/components/ui/badge"

import { formatExpiryDate, roleLabel } from "../../team.utils"
import type { WorkspaceInvitation } from "../../team.types"

interface InvitationRowProps {
  invitation: WorkspaceInvitation
}

export function InvitationRow({ invitation }: InvitationRowProps) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-border px-4 py-4 last:border-b-0">
      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-2">
          <p className="truncate text-sm font-medium text-foreground">{invitation.email}</p>
          <Badge variant="outline">{roleLabel(invitation.role)}</Badge>
          <Badge variant="warning">Pending</Badge>
        </div>
        <p className="mt-1 text-xs text-text-tertiary">
          Expires {formatExpiryDate(invitation.expires_at)}
        </p>
      </div>
    </div>
  )
}
