import { useState } from "react"
import { Modal } from "@/components/ui/modal"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { PinInput } from "@/components/ui/pin-input"
import { apiFetch } from "@/core/api/client"

interface JoinWorkspaceModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: (workspaceName: string) => void
}

export function JoinWorkspaceModal({ isOpen, onClose, onSuccess }: JoinWorkspaceModalProps) {
  const [slug, setSlug] = useState("")
  const [pin, setPin] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  const validate = () => {
    const errors: Record<string, string> = {}
    if (!slug.trim()) {
      errors.slug = "Workspace slug is required"
    }
    if (!pin.trim()) {
      errors.pin = "Security PIN is required"
    } else if (pin.length !== 6) {
      errors.pin = "PIN must be exactly 6 digits"
    }
    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    setIsLoading(true)
    setError(null)

    try {
      const res = await apiFetch("/api/v1/workspaces/join", {
        method: "POST",
        body: JSON.stringify({
          slug: slug.trim().toLowerCase(),
          pin: pin.trim(),
        }),
      })

      if (!res.ok) {
        const err = (await res.json().catch(() => ({}))) as { message?: string }
        throw new Error(err.message ?? "Failed to join workspace. Please check the slug and PIN.")
      }

      onSuccess(slug.trim())
      handleClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to join workspace")
    } finally {
      setIsLoading(false)
    }
  }

  const handleClose = () => {
    setSlug("")
    setPin("")
    setError(null)
    setValidationErrors({})
    onClose()
  }

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Join Workspace">
      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <p className="text-sm text-text-secondary">
          Enter the slug and security PIN of the workspace you want to join.
        </p>

        <div className="flex flex-col gap-1.5">
          <label htmlFor="join-key" className="text-xs font-semibold text-text-secondary uppercase">
            Workspace Slug
          </label>
          <Input
            id="join-key"
            type="text"
            value={slug}
            onChange={(e) => {
              setSlug(e.target.value)
              setValidationErrors((prev) => ({ ...prev, slug: "" }))
            }}
            placeholder="e.g. marketing-team"
            className="w-full bg-secondary border-border h-11 px-4 font-mono"
          />
          {validationErrors.slug && (
            <p className="text-xs text-destructive font-medium mt-0.5">{validationErrors.slug}</p>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <span className="text-xs font-semibold text-text-secondary uppercase">
            Security PIN
          </span>
          <div className="py-2">
            <PinInput
              length={6}
              value={pin}
              onChange={(value) => {
                setPin(value)
                setValidationErrors((prev) => ({ ...prev, pin: "" }))
              }}
            />
          </div>
          {validationErrors.pin && (
            <p className="text-xs text-destructive font-medium text-center mt-0.5">{validationErrors.pin}</p>
          )}
        </div>

        {error && (
          <p className="text-xs text-destructive text-center font-medium mt-2">{error}</p>
        )}

        <div className="flex flex-col gap-3 pt-2">
          <Button type="submit" disabled={isLoading}>
            {isLoading ? "Joining..." : "Join Workspace"}
          </Button>
          <button
            type="button"
            onClick={handleClose}
            className="w-full bg-secondary hover:bg-muted text-text-secondary border border-border h-10 rounded-lg font-semibold text-sm transition-all active:scale-[0.98]"
          >
            Cancel
          </button>
        </div>
      </form>
    </Modal>
  )
}
