import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"

interface ApiKeySecretDialogProps {
  open: boolean
  clientId?: string
  clientSecret: string
  onClose: () => void
}

export function ApiKeySecretDialog({
  open,
  clientId,
  clientSecret,
  onClose,
}: ApiKeySecretDialogProps) {
  const [copiedField, setCopiedField] = useState<"clientId" | "clientSecret" | null>(null)

  async function copyValue(value: string, field: "clientId" | "clientSecret") {
    try {
      await navigator.clipboard.writeText(value)
      setCopiedField(field)
      window.setTimeout(() => setCopiedField((current) => (current === field ? null : current)), 2000)
    } catch {
      setCopiedField(null)
    }
  }

  return (
    <Modal isOpen={open} onClose={onClose} title="Save your API credentials">
      <div className="flex flex-col gap-4">
        <p className="text-sm text-text-secondary">
          Copy these values now. The client secret is only shown once and cannot be retrieved later.
        </p>

        {clientId ? (
          <div className="flex flex-col gap-2">
            <label htmlFor="api-key-client-id" className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
              Client ID
            </label>
            <div className="flex gap-2">
              <Input
                id="api-key-client-id"
                readOnly
                value={clientId}
                className="font-mono text-sm"
                onFocus={(event) => event.target.select()}
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="shrink-0"
                onClick={() => void copyValue(clientId, "clientId")}
              >
                <Icon name="content_copy" size="sm" data-icon="inline-start" />
                {copiedField === "clientId" ? "Copied" : "Copy"}
              </Button>
            </div>
          </div>
        ) : null}

        <div className="flex flex-col gap-2">
          <label htmlFor="api-key-client-secret" className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            Client secret
          </label>
          <div className="flex gap-2">
            <Input
              id="api-key-client-secret"
              readOnly
              value={clientSecret}
              className="font-mono text-sm"
              onFocus={(event) => event.target.select()}
            />
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="shrink-0"
              onClick={() => void copyValue(clientSecret, "clientSecret")}
            >
              <Icon name="content_copy" size="sm" data-icon="inline-start" />
              {copiedField === "clientSecret" ? "Copied" : "Copy"}
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
