import { useState } from "react"
import { Modal } from "@/components/ui/modal"
import { PinInput } from "@/components/ui/pin-input"
import { Button } from "@/components/ui/button"

interface PinEntryModalProps {
  isOpen: boolean
  workspaceName: string
  onClose: () => void
  onSubmit: (pin: string) => void
  error?: string
  isLoading?: boolean
}

export function PinEntryModal({
  isOpen,
  workspaceName,
  onClose,
  onSubmit,
  error: externalError,
  isLoading = false,
}: PinEntryModalProps) {
  const [pin, setPin] = useState("")
  const [internalError, setInternalError] = useState("")

  const error = externalError || internalError

  const handlePinChange = (value: string) => {
    setPin(value)
    setInternalError("")
  };

  const handleSubmit = () => {
    if (pin.length !== 6) {
      setInternalError("Please enter a complete 6-digit PIN")
      return
    }
    onSubmit(pin)
  }

  const handleCancel = () => {
    setPin("")
    setInternalError("")
    onClose()
  }

  const handleComplete = (value: string) => {
    if (value.length === 6) {
      onSubmit(value)
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={handleCancel} title="Enter Workspace PIN">
      <div className="flex flex-col gap-4">
        <p className="text-sm text-text-secondary">
          The &apos;{workspaceName}&apos; workspace is private. Please enter your 6-digit security code.
        </p>

        <div className="flex flex-col gap-2 py-4">
          <PinInput
            length={6}
            value={pin}
            onChange={handlePinChange}
            onComplete={handleComplete}
          />
          {error && (
            <p className="text-xs text-destructive text-center font-medium mt-1">{error}</p>
          )}
        </div>

        <div className="flex flex-col gap-3 pt-2">
          <Button onClick={handleSubmit} disabled={pin.length !== 6 || isLoading}>
            {isLoading ? "Unlocking..." : "Unlock Workspace"}
          </Button>
          <button
            type="button"
            onClick={handleCancel}
            className="w-full bg-secondary hover:bg-muted text-text-secondary border border-border h-10 rounded-lg font-semibold text-sm transition-colors active:scale-[0.98]"
          >
            Cancel
          </button>
        </div>
      </div>
    </Modal>
  )
}
