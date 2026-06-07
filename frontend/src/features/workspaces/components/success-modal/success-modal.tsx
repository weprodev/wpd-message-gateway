import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Modal } from "@/components/ui/modal"

interface SuccessModalProps {
  isOpen: boolean
  workspaceName: string
  onContinue: () => void
}

export function SuccessModal({ isOpen, workspaceName, onContinue }: SuccessModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onContinue} title="">
      <div className="flex flex-col gap-6 items-center text-center py-4">
        <div className="bg-emerald-500/10 dark:bg-emerald-500/20 rounded-full size-16 flex items-center justify-center border border-emerald-500/30">
          <Icon name="check_circle" size="lg" className="text-emerald-600 dark:text-emerald-400" />
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-2xl font-bold text-foreground">
            Workspace Created!
          </h2>
          <p className="text-sm text-text-secondary leading-relaxed">
            Your workspace <span className="font-semibold text-foreground">&quot;{workspaceName}&quot;</span> has been created successfully.
          </p>
        </div>

        <p className="text-xs text-text-tertiary">
          You can now start managing your messages and communications.
        </p>

        <div className="w-full pt-2">
          <Button onClick={onContinue} className="w-full">
            Go to Dashboard
          </Button>
        </div>
      </div>
    </Modal>
  )
}
