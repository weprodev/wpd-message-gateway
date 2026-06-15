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
  /** When true, blocks overlay, escape, and close-button dismissal until the caller closes the modal. */
  preventDismiss?: boolean
}

export function Modal({ isOpen, onClose, children, title, preventDismiss = false }: ModalProps) {
  return (
    <Dialog
      open={isOpen}
      onOpenChange={(open) => {
        if (!open && !preventDismiss) onClose()
      }}
    >
      <DialogContent
        className="bg-card rounded-[16px] p-6 sm:p-[32px] w-full max-w-[480px] flex flex-col gap-6 border border-border shadow-lg animate-in zoom-in-95 duration-200"
        showCloseButton={!preventDismiss}
        onInteractOutside={preventDismiss ? (event) => event.preventDefault() : undefined}
        onEscapeKeyDown={preventDismiss ? (event) => event.preventDefault() : undefined}
      >
        {title && (
          <div className="flex items-center justify-between w-full">
            <DialogTitle className="text-xl font-semibold text-foreground">
              {title}
            </DialogTitle>
          </div>
        )}
        {!title && (
          <DialogPrimitive.Title className="sr-only">
            Modal Dialog
          </DialogPrimitive.Title>
        )}
        <div className="flex flex-col gap-4">
          {children}
        </div>
      </DialogContent>
    </Dialog>
  )
}
