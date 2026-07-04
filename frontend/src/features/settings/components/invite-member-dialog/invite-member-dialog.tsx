import { useState } from "react"

import { Role, type WorkspaceRole } from "@/core/auth"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"
import { cn } from "@/lib/utils"

import { INVITABLE_ROLES, hasPendingInvitationEmail, isExistingMemberEmail, roleLabel } from "../../team.utils"
import type { CreateInvitationResult, InvitableRole, WorkspaceInvitation, WorkspaceMember } from "../../team.types"

interface InviteMemberDialogProps {
  open: boolean
  onClose: () => void
  onInvite: (email: string, role: InvitableRole) => Promise<CreateInvitationResult>
  members?: WorkspaceMember[]
  invitations?: WorkspaceInvitation[]
}

export function InviteMemberDialog({
  open,
  onClose,
  onInvite,
  members = [],
  invitations = [],
}: InviteMemberDialogProps) {
  const [email, setEmail] = useState("")
  const [role, setRole] = useState<InvitableRole>(Role.Member)
  const [validationError, setValidationError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  function resetForm() {
    setEmail("")
    setRole(Role.Member)
    setValidationError(null)
    setSubmitError(null)
    setIsSubmitting(false)
  }

  function handleClose() {
    if (isSubmitting) return
    resetForm()
    onClose()
  }

  function isValidEmail(value: string): boolean {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)
  }

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault()
    const trimmedEmail = email.trim().toLowerCase()
    if (!trimmedEmail) {
      setValidationError("Email is required")
      return
    }
    if (!isValidEmail(trimmedEmail)) {
      setValidationError("Enter a valid email address")
      return
    }
    if (isExistingMemberEmail(members, trimmedEmail)) {
      setValidationError("This person is already a member of the workspace")
      return
    }
    if (hasPendingInvitationEmail(invitations, trimmedEmail)) {
      setValidationError("A pending invitation already exists for this email")
      return
    }

    setValidationError(null)
    setSubmitError(null)
    setIsSubmitting(true)
    try {
      await onInvite(trimmedEmail, role)
      resetForm()
      onClose()
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Failed to send invitation")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal isOpen={open} onClose={handleClose} title="Invite teammate" preventClose={isSubmitting}>
      <form noValidate onSubmit={handleSubmit} className="flex flex-col gap-4">
        <p className="text-sm text-text-secondary">
          Send an invitation email with a role. The invitee uses the one-time token to join this workspace.
        </p>

        <div className="flex flex-col gap-2">
          <label htmlFor="invite-email" className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            Email address
          </label>
          <Input
            id="invite-email"
            type="email"
            value={email}
            onChange={(event) => {
              setEmail(event.target.value)
              setValidationError(null)
              setSubmitError(null)
            }}
            placeholder="teammate@company.com"
            autoFocus
            disabled={isSubmitting}
          />
        </div>

        <div className="flex flex-col gap-2">
          <label htmlFor="invite-role" className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            Role
          </label>
          <select
            id="invite-role"
            value={role}
            onChange={(event) => setRole(event.target.value as WorkspaceRole)}
            disabled={isSubmitting}
            className={cn(
              "flex h-10 w-full rounded-md border border-border bg-input px-3 py-2 text-sm text-foreground",
              "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
            )}
          >
            {INVITABLE_ROLES.map((option) => (
              <option key={option} value={option}>
                {roleLabel(option)}
              </option>
            ))}
          </select>
        </div>

        {validationError ? (
          <p className="text-xs font-medium text-destructive">{validationError}</p>
        ) : null}

        {submitError ? (
          <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
            {submitError}
          </p>
        ) : null}

        <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <Button type="button" variant="outline" onClick={handleClose} disabled={isSubmitting}>
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Sending…" : "Send invitation"}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
