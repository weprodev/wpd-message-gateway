import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"

import { IntegrationProviderIcon } from "@/features/integrations/components/integration-provider-icon"
import type { IntegrationViewModel } from "@/features/integrations/hooks/use-integrations.hook"
import type { IntegrationActionResult } from "@/features/integrations/integrations.types"

type DisconnectAction = "deactivate" | "remove"

interface DisconnectModalProps {
  isOpen: boolean
  provider: IntegrationViewModel | null
  onClose: () => void
  onDeactivate: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>
  onRemove: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>
}

export function DisconnectModal({
  isOpen,
  provider,
  onClose,
  onDeactivate,
  onRemove,
}: DisconnectModalProps) {
  const [submittingAction, setSubmittingAction] = useState<DisconnectAction | null>(null)
  const [error, setError] = useState<string | null>(null)

  const isSubmitting = submittingAction !== null

  function handleClose() {
    if (isSubmitting) return
    setError(null)
    setSubmittingAction(null)
    onClose()
  }

  async function runAction(action: DisconnectAction) {
    if (!provider) return

    setSubmittingAction(action)
    setError(null)
    try {
      const result = action === "deactivate" ? await onDeactivate(provider) : await onRemove(provider)
      if (!result.ok) {
        setError(result.message ?? "Failed to update provider")
      } else {
        setError(null)
        onClose()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update provider")
    } finally {
      setSubmittingAction(null)
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      preventDismiss={isSubmitting}
      title={
        provider ? (
          <div className="flex items-start gap-4">
            <IntegrationProviderIcon
              icon={provider.icon}
              name={provider.name}
              className="size-12 p-2.5 text-2xl"
            />
            <span className="min-w-0 flex-1 pt-1">Disconnect {provider.name}</span>
          </div>
        ) : undefined
      }
      description="Are you sure you want to disconnect this provider?"
    >
      {provider ? (
        <div className="flex flex-col gap-6">
          {error ? (
            <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
              {error}
            </p>
          ) : null}

          <div className="flex flex-col gap-3">
            <div className="flex items-start gap-3 rounded-lg border border-border bg-muted/40 p-4">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-full border border-border bg-card">
                <Icon name="pause" size="sm" className="text-text-secondary" />
              </div>
              <div className="min-w-0 flex-1">
                <h3 className="text-sm font-semibold text-foreground">Deactivate Connection</h3>
                <p className="mt-1 text-[13px] leading-normal text-text-secondary">
                  Stop all message delivery immediately. Your API keys and configurations will be preserved for later use.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3 rounded-lg border border-destructive/20 bg-destructive/5 p-4">
              <div className="flex size-10 shrink-0 items-center justify-center rounded-full border border-destructive/20 bg-card">
                <Icon name="delete" size="sm" className="text-destructive" />
              </div>
              <div className="min-w-0 flex-1">
                <h3 className="text-sm font-semibold text-foreground">Remove Integration</h3>
                <p className="mt-1 text-[13px] leading-normal text-text-secondary">
                  Permanently delete all API keys and configurations associated with this provider. This action cannot be undone.
                </p>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap items-center justify-end gap-3 border-t border-border pt-4">
            <Button type="button" variant="outline" onClick={handleClose} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => runAction("deactivate")}
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
              onClick={() => runAction("remove")}
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
