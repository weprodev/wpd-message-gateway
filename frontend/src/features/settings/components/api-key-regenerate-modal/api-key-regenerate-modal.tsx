import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"

import { ApiKeyCredentialsView } from "../api-key-credentials-view"
import type { ApiKeyCredentials } from "../../settings.types"

interface ApiKeyRegenerateModalProps {
  isOpen: boolean
  onClose: () => void
  onRegenerate: () => Promise<ApiKeyCredentials | null>
}

type RegenerateStep = "confirm" | "credentials"

export function ApiKeyRegenerateModal({ isOpen, onClose, onRegenerate }: ApiKeyRegenerateModalProps) {
  const [step, setStep] = useState<RegenerateStep>("confirm")
  const [isRegenerating, setIsRegenerating] = useState(false)
  const [credentials, setCredentials] = useState<ApiKeyCredentials | null>(null)

  function resetAndClose() {
    setStep("confirm")
    setCredentials(null)
    onClose()
  }

  async function handleRegenerate() {
    setIsRegenerating(true)
    try {
      const result = await onRegenerate()
      if (result) {
        setCredentials(result)
        setStep("credentials")
      }
    } finally {
      setIsRegenerating(false)
    }
  }

  const title = step === "confirm" ? "Regenerate API key" : "API key regenerated"

  return (
    <Modal
      isOpen={isOpen}
      onClose={resetAndClose}
      title={title}
      preventDismiss={step === "credentials"}
    >
      {step === "confirm" ? (
        <div className="flex flex-col gap-6">
          <p className="text-sm text-text-secondary leading-relaxed">
            Are you sure you want to regenerate this API key?
          </p>

          <div className="flex items-center justify-end gap-3 border-t border-border pt-4">
            <Button type="button" variant="outline" onClick={resetAndClose} disabled={isRegenerating}>
              Cancel
            </Button>
            <Button type="button" onClick={handleRegenerate} disabled={isRegenerating} className="min-w-[9.75rem]">
              {isRegenerating ? (
                <>
                  <Spinner size="sm" variant="onSolid" />
                  Regenerating…
                </>
              ) : (
                "Regenerate"
              )}
            </Button>
          </div>
        </div>
      ) : credentials ? (
        <ApiKeyCredentialsView credentials={credentials} onConfirm={resetAndClose} />
      ) : null}
    </Modal>
  )
}
