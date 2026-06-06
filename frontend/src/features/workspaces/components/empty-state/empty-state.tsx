import { Icon } from "@/components/ui/icon"
import { WorkspaceActions } from "../workspace-actions"

interface EmptyStateProps {
  onCreateWorkspace?: () => void
  onJoinWorkspace?: () => void
}

export function EmptyState({ onCreateWorkspace, onJoinWorkspace }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-8 w-full max-w-[640px] py-12 px-6">
      <div className="bg-secondary rounded-full size-24 flex items-center justify-center border border-border">
        <Icon name="folder_off" size="lg" className="text-text-tertiary" />
      </div>

      <div className="flex flex-col gap-3 items-center text-center w-full">
        <h2 className="text-2xl sm:text-3xl font-bold text-foreground">
          No Workspaces Yet
        </h2>
        <p className="text-sm sm:text-base text-text-secondary">
          Create your first workspace to get started with Message Gateway.
        </p>
        <p className="text-xs font-semibold text-text-tertiary">
          Connecting your world seamlessly!
        </p>
      </div>

      <div className="w-full max-w-sm">
        <WorkspaceActions
          variant="inline"
          onCreateWorkspace={onCreateWorkspace}
          onJoinWorkspace={onJoinWorkspace}
        />
      </div>
    </div>
  )
}
