import type { ReactNode } from "react"
import * as DialogPrimitive from "@radix-ui/react-dialog"
import {
  Dialog,
  DialogContent,
  DialogTitle,
} from "@/components/ui/dialog"

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  children: ReactNode
  title?: string
  header?: ReactNode
  description?: ReactNode
  /** When true, blocks overlay, escape, and close-button dismissal until the caller closes the modal. */
  preventDismiss?: boolean
  /** Alias for preventDismiss — kept for callers that predated the rename. */
  preventClose?: boolean
}

export function Modal({
  isOpen,
  onClose,
  children,
  title,
  header,
  description,
  preventDismiss = false,
  preventClose = false,
}: ModalProps) {
  const locked = preventDismiss || preventClose

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(open) => {
        if (!open && !locked) onClose()
      }}
    >
      <DialogContent
        className="bg-card rounded-[16px] p-6 sm:p-[32px] w-full max-w-[480px] flex flex-col gap-6 border border-border shadow-lg animate-in zoom-in-95 duration-200"
        showCloseButton={!locked}
        onInteractOutside={locked ? (event) => event.preventDefault() : undefined}
        onEscapeKeyDown={locked ? (event) => event.preventDefault() : undefined}
      >
        {header}
        {!header && title ? (
          <div className="flex w-full flex-col gap-2">
            <DialogTitle className="text-xl font-semibold text-foreground">
              {title}
            </DialogTitle>
            {description ? (
              <div className="text-sm leading-relaxed text-text-secondary">{description}</div>
            ) : null}
          </div>
        ) : null}
        {!header && !title ? (
          <DialogPrimitive.Title className="sr-only">
            Modal Dialog
          </DialogPrimitive.Title>
        ) : null}
        <div className="flex flex-col gap-4">
          {children}
        </div>
      </DialogContent>
    </Dialog>
  )
}
