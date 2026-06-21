import { Button } from "@/components/ui/button"
import { Spinner } from "@/components/ui/spinner"

type ConfirmVariant = "default" | "destructive" | "outline" | "secondary" | "ghost" | "link" | "brand"

interface ModalActionsProps {
  cancelLabel?: string
  confirmLabel: string
  onCancel: () => void
  onConfirm?: () => void
  confirmVariant?: ConfirmVariant
  confirmType?: "button" | "submit"
  confirmForm?: string
  isLoading?: boolean
  loadingLabel?: string
  cancelDisabled?: boolean
  confirmDisabled?: boolean
  confirmClassName?: string
}

export function ModalActions({
  cancelLabel = "Cancel",
  confirmLabel,
  onCancel,
  onConfirm,
  confirmVariant = "default",
  confirmType = "button",
  confirmForm,
  isLoading = false,
  loadingLabel,
  cancelDisabled = false,
  confirmDisabled = false,
  confirmClassName,
}: ModalActionsProps) {
  const busy = isLoading || cancelDisabled

  return (
    <div className="flex items-center justify-end gap-3 border-t border-border pt-4">
      <Button type="button" variant="outline" onClick={onCancel} disabled={busy}>
        {cancelLabel}
      </Button>
      <Button
        type={confirmType}
        form={confirmForm}
        variant={confirmVariant}
        onClick={confirmType === "button" ? onConfirm : undefined}
        disabled={busy || confirmDisabled || isLoading}
        className={confirmClassName}
      >
        {isLoading ? (
          <>
            <Spinner size="sm" variant="onSolid" />
            {loadingLabel ?? confirmLabel}
          </>
        ) : (
          confirmLabel
        )}
      </Button>
    </div>
  )
}
