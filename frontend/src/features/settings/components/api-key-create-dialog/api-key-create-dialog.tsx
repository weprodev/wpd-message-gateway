import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"

interface ApiKeyCreateDialogProps {
  open: boolean
  onClose: () => void
  onCreate: (name: string) => Promise<void>
}

export function ApiKeyCreateDialog({ open, onClose, onCreate }: ApiKeyCreateDialogProps) {
  const [name, setName] = useState("")
  const [validationError, setValidationError] = useState<string | null>(null)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  function resetForm() {
    setName("")
    setValidationError(null)
    setSubmitError(null)
    setIsSubmitting(false)
  }

  function handleClose() {
    if (isSubmitting) return
    resetForm()
    onClose()
  }

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) {
      setValidationError("API key name is required")
      return
    }

    setValidationError(null)
    setSubmitError(null)
    setIsSubmitting(true)
    try {
      await onCreate(trimmed)
      resetForm()
      onClose()
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Failed to create API key")
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal isOpen={open} onClose={handleClose} title="Generate API key" preventClose={isSubmitting}>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <p className="text-sm text-text-secondary">
          Give this key a name so you can identify it later (for example, &quot;Production&quot; or &quot;CI&quot;).
        </p>

        <div className="flex flex-col gap-2">
          <label htmlFor="api-key-name" className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            API key name
          </label>
          <Input
            id="api-key-name"
            type="text"
            value={name}
            onChange={(event) => {
              setName(event.target.value)
              setValidationError(null)
              setSubmitError(null)
            }}
            placeholder="e.g. Production"
            autoFocus
            disabled={isSubmitting}
          />
          {validationError ? (
            <p className="text-xs font-medium text-destructive">{validationError}</p>
          ) : null}
        </div>

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
            {isSubmitting ? "Generating…" : "Generate key"}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
