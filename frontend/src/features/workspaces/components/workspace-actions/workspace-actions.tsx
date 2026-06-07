import { Icon } from "@/components/ui/icon"

interface WorkspaceActionsProps {
  onCreateWorkspace?: () => void
  onJoinWorkspace?: () => void
  variant?: "card" | "inline"
}

export function WorkspaceActions({
  onCreateWorkspace,
  onJoinWorkspace,
  variant = "card",
}: WorkspaceActionsProps) {
  const containerClass = variant === "card"
    ? "bg-card border border-border rounded-2xl p-4 w-full max-w-[560px] flex flex-col sm:flex-row gap-4 shadow-sm"
    : "flex flex-col sm:flex-row gap-4 w-full"

  return (
    <div className={containerClass}>
      <button
        type="button"
        onClick={onCreateWorkspace}
        className="flex-1 bg-primary-brand hover:bg-primary-brand-hover text-white px-5 h-12 rounded-lg font-semibold text-sm flex items-center justify-center gap-2 transition-all active:scale-[0.98]"
      >
        <Icon name="add" size="sm" />
        Create Workspace
      </button>

      <button
        type="button"
        onClick={onJoinWorkspace}
        className="flex-1 bg-secondary hover:bg-muted text-foreground border border-border px-5 h-12 rounded-lg font-semibold text-sm flex items-center justify-center gap-2 transition-all active:scale-[0.98]"
      >
        <Icon name="login" size="sm" className="text-text-secondary" />
        Join Workspace
      </button>
    </div>
  )
}
