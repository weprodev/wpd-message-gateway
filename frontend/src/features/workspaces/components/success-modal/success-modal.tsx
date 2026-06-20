import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Modal } from "@/components/ui/modal"

export type SuccessModalVariant = "created" | "joined"

interface SuccessModalProps {
  isOpen: boolean
  workspaceName: string
  variant?: SuccessModalVariant
  onClose: () => void
  onContinue: () => void
}

const copy: Record<
  SuccessModalVariant,
  { title: string; prefix: string; suffix: string }
> = {
  created: {
    title: "Workspace Created!",
    prefix: "Your workspace",
    suffix: " has been created successfully.",
  },
  joined: {
    title: "Workspace Joined!",
    prefix: "You have successfully joined",
    suffix: ".",
  },
}

export function SuccessModal({
  isOpen,
  workspaceName,
  variant = "created",
  onClose,
  onContinue,
}: SuccessModalProps) {
  const { title, prefix, suffix } = copy[variant]

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="">
      <div className="flex flex-col gap-6 items-center text-center py-4">
        <div className="bg-emerald-500/10 dark:bg-emerald-500/20 rounded-full size-16 flex items-center justify-center border border-emerald-500/30">
          <Icon name="check_circle" size="lg" className="text-emerald-600 dark:text-emerald-400" />
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-2xl font-bold text-foreground">{title}</h2>
          <p className="text-sm text-text-secondary leading-relaxed">
            {prefix}{" "}
            <span className="font-semibold text-foreground">&quot;{workspaceName}&quot;</span>
            {suffix}
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
