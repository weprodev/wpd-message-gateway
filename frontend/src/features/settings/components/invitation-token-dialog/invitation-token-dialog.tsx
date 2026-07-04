import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"

import { roleLabel } from "../../team.utils"

interface InvitationTokenDialogProps {
  open: boolean
  email: string
  role: string
  token: string
  onClose: () => void
}

export function InvitationTokenDialog({
  open,
  email,
  role,
  token,
  onClose,
}: InvitationTokenDialogProps) {
  const [copied, setCopied] = useState(false)

  async function copyToken() {
    try {
      await navigator.clipboard.writeText(token)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 2000)
    } catch {
      setCopied(false)
    }
  }

  return (
    <Modal isOpen={open} onClose={onClose} title="Invitation created">
      <div className="flex flex-col gap-4">
        <p className="text-sm text-text-secondary">
          Share this one-time token with <span className="font-medium text-foreground">{email}</span> as{" "}
          <span className="font-medium text-foreground">{roleLabel(role)}</span>. It will not be shown again.
        </p>

        <div className="flex flex-col gap-2">
          <label htmlFor="invitation-token" className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            Invitation token
          </label>
          <div className="flex gap-2">
            <Input
              id="invitation-token"
              readOnly
              value={token}
              className="font-mono text-sm"
              onFocus={(event) => event.target.select()}
            />
            <Button type="button" variant="outline" size="sm" className="shrink-0" onClick={() => void copyToken()}>
              <Icon name="content_copy" size="sm" data-icon="inline-start" />
              {copied ? "Copied" : "Copy"}
            </Button>
          </div>
        </div>

        <Button type="button" className="w-full" onClick={onClose}>
          Done
        </Button>
      </div>
    </Modal>
  )
}
