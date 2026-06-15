import { useEffect, useRef, useState, type FormEvent } from "react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"

import { toUserMessage } from "@/lib/errors"

import { ApiKeyCredentialsView } from "../api-key-credentials-view"
import type { ApiKeyCredentials } from "../../settings.types"

const CREATE_FAILED_MESSAGE = "Could not create the API key. Please try again."

interface ApiKeyCreateModalProps {
  isOpen: boolean
  onClose: () => void
  onCreate: (name: string) => Promise<ApiKeyCredentials>
}

type CreateStep = "form" | "credentials"

export function ApiKeyCreateModal({ isOpen, onClose, onCreate }: ApiKeyCreateModalProps) {
  const [step, setStep] = useState<CreateStep>("form")
  const [name, setName] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [isCreating, setIsCreating] = useState(false)
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
      setStep("form")
      setName("")
      setError(null)
      setIsCreating(false)
      setCredentials(null)
    }
  }

  function resetForm() {
    setStep("form")
    setName("")
    setError(null)
    setIsCreating(false)
    setCredentials(null)
  }

  function resetAndClose() {
    resetForm()
    onClose()
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) {
      setError("API key name is required")
      return
    }
    const session = sessionRef.current
    setError(null)
    setIsCreating(true)
    try {
      const result = await onCreate(trimmed)
      if (session !== sessionRef.current) return
      setCredentials(result)
      setStep("credentials")
    } catch (err) {
      if (session !== sessionRef.current) return
      setError(toUserMessage(err, CREATE_FAILED_MESSAGE))
    } finally {
      if (session === sessionRef.current) {
        setIsCreating(false)
      }
    }
  }

  const title = step === "form" ? "Create API key" : "API key created"

  return (
    <Modal
      isOpen={isOpen}
      onClose={resetAndClose}
      title={title}
      preventDismiss={step === "credentials" || isCreating}
    >
      {step === "form" ? (
        <form onSubmit={handleSubmit} className="flex flex-col gap-6">
          <div className="flex flex-col gap-1.5">
            <label htmlFor="api-key-name" className="text-xs font-semibold uppercase text-text-secondary">
              Name
            </label>
            <Input
              id="api-key-name"
              type="text"
              value={name}
              onChange={(event) => {
                setName(event.target.value)
                if (error) setError(null)
              }}
              placeholder="e.g. Production"
              disabled={isCreating}
              autoFocus
            />
            {error ? <p className="text-sm text-destructive">{error}</p> : null}
          </div>

          <div className="flex items-center justify-end gap-3 border-t border-border pt-4">
            <Button type="button" variant="outline" onClick={resetAndClose} disabled={isCreating}>
              Cancel
            </Button>
            <Button type="submit" disabled={isCreating} className="min-w-[7.25rem]">
              {isCreating ? (
                <>
                  <Spinner size="sm" variant="onSolid" />
                  Creating…
                </>
              ) : (
                "Create"
              )}
            </Button>
          </div>
        </form>
      ) : credentials ? (
        <ApiKeyCredentialsView credentials={credentials} onConfirm={resetAndClose} />
      ) : null}
    </Modal>
  )
}
