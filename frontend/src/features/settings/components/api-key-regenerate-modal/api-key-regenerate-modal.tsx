import { useEffect, useRef, useState } from "react"

import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"
import { toUserMessage } from "@/lib/errors"

import { ApiKeyCredentialsView } from "../api-key-credentials-view"
import type { ApiKeyCredentials } from "../../settings.types"

const REGENERATE_FAILED_MESSAGE = "Could not regenerate the API key. Please try again."

interface ApiKeyRegenerateModalProps {
  isOpen: boolean
  onClose: () => void
  onRegenerate: () => Promise<ApiKeyCredentials>
}

type RegenerateStep = "confirm" | "credentials"

export function ApiKeyRegenerateModal({ isOpen, onClose, onRegenerate }: ApiKeyRegenerateModalProps) {
  const [step, setStep] = useState<RegenerateStep>("confirm")
  const [isRegenerating, setIsRegenerating] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [credentials, setCredentials] = useState<ApiKeyCredentials | null>(null)
  const [prevIsOpen, setPrevIsOpen] = useState(isOpen)
  const [session, setSession] = useState(0)
  const sessionRef = useRef(session)

  useEffect(() => {
    sessionRef.current = session
  }, [session])

  if (isOpen !== prevIsOpen) {
    setPrevIsOpen(isOpen)
    setSession((current) => current + 1)
    if (isOpen) {
      setStep("confirm")
      setError(null)
      setIsRegenerating(false)
      setCredentials(null)
    }
  }

  function resetForm() {
    setStep("confirm")
    setError(null)
    setIsRegenerating(false)
    setCredentials(null)
  }

  function resetAndClose() {
    resetForm()
    onClose()
  }

  async function handleRegenerate() {
    const activeSession = sessionRef.current
    setError(null)
    setIsRegenerating(true)
    try {
      const result = await onRegenerate()
      if (activeSession !== sessionRef.current) return
      setCredentials(result)
      setStep("credentials")
    } catch (err) {
      if (activeSession !== sessionRef.current) return
      setError(toUserMessage(err, REGENERATE_FAILED_MESSAGE))
    } finally {
      if (activeSession === sessionRef.current) {
        setIsRegenerating(false)
      }
    }
  }

  const title = step === "confirm" ? "Regenerate API key" : "API key regenerated"

  return (
    <Modal
      isOpen={isOpen}
      onClose={resetAndClose}
      title={title}
      preventDismiss={step === "credentials" || isRegenerating}
    >
      {step === "confirm" ? (
        <div className="flex flex-col gap-6">
          <p className="text-sm text-text-secondary leading-relaxed">
            Are you sure you want to regenerate this API key?
          </p>
          {error ? <p className="text-sm text-destructive">{error}</p> : null}

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
