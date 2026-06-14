import { useState } from "react"

import { Button } from "@/components/ui/button"
import { DialogDescription, DialogTitle } from "@/components/ui/dialog"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"

import { IntegrationProviderIcon } from "../integration-provider-icon"
import type { IntegrationViewModel } from "../../hooks/use-integrations.hook"
import type { IntegrationActionResult } from "../../integrations.types"
import { DisconnectOptionCard } from "./disconnect-option-card"

type SubmittingAction = "deactivate" | "remove" | null

interface DisconnectModalProps {
  isOpen: boolean
  onClose: () => void
  provider: IntegrationViewModel | null
  onDeactivate: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>
  onRemove: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>
}

export function DisconnectModal({
  isOpen,
  onClose,
  provider,
  onDeactivate,
  onRemove,
}: DisconnectModalProps) {
  const [submittingAction, setSubmittingAction] = useState<SubmittingAction>(null)
  const [error, setError] = useState<string | null>(null)

  const isSubmitting = submittingAction !== null

  const runAction = async (
    action: SubmittingAction,
    handler: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>,
  ) => {
    if (!provider || !action) return

    setSubmittingAction(action)
    setError(null)
    try {
      const result = await handler(provider)
      if (!result.ok) {
        setError(result.message ?? "Failed to update provider")
      } else {
        onClose()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update provider")
    } finally {
      setSubmittingAction(null)
    }
  }

  const handleClose = () => {
    if (isSubmitting) return
    setError(null)
    onClose()
  }

  return (
    <Modal isOpen={isOpen} onClose={handleClose}>
      {isOpen && provider ? (
        <div className="flex flex-col gap-6">
          <div className="flex items-start gap-4">
            <IntegrationProviderIcon
              icon={provider.icon}
              name={provider.name}
              className="size-12 p-2.5 text-2xl"
            />
            <div className="min-w-0 flex-1">
              <DialogTitle className="text-xl font-semibold text-foreground">
                Disconnect {provider.name}
              </DialogTitle>
              <DialogDescription className="mt-1 text-sm text-text-secondary">
                Are you sure you want to disconnect this provider?
              </DialogDescription>
            </div>
          </div>

          {error ? (
            <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
              {error}
            </p>
          ) : null}

          <div className="flex flex-col gap-3">
            <DisconnectOptionCard
              variant="neutral"
              icon="pause"
              title="Deactivate Connection"
              description="Stop all message delivery immediately. Your API keys and configurations will be preserved for later use."
            />
            <DisconnectOptionCard
              variant="danger"
              icon="delete"
              title="Remove Integration"
              description="Permanently delete all API keys and configurations associated with this provider. This action cannot be undone."
            />
          </div>

          <div className="flex flex-wrap items-center justify-end gap-3 border-t border-border pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => runAction("deactivate", onDeactivate)}
              disabled={isSubmitting}
              className="border-primary-brand text-primary-brand hover:bg-primary-brand/5"
            >
              {submittingAction === "deactivate" ? (
                <span className="flex items-center gap-2">
                  <Spinner size="sm" />
                  Deactivating...
                </span>
              ) : (
                "Deactivate"
              )}
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={() => runAction("remove", onRemove)}
              disabled={isSubmitting}
            >
              {submittingAction === "remove" ? (
                <span className="flex items-center gap-2">
                  <Spinner size="sm" />
                  Removing...
                </span>
              ) : (
                "Remove Integration"
              )}
            </Button>
          </div>
        </div>
      ) : null}
    </Modal>
  )
}
