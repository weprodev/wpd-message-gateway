import { Badge } from "@/components/ui/badge"
import { Icon } from "@/components/ui/icon"
import { resolveWorkspaceIconName } from "../../workspace-icons"
import type { Workspace } from "../../workspace.types"

interface WorkspaceCardProps {
  workspace: Workspace
  isSelected: boolean
  onSelect: (workspace: Workspace) => void
}

export function WorkspaceCard({ workspace, isSelected, onSelect }: WorkspaceCardProps) {
  const borderClass = isSelected
    ? "border-2 border-primary-brand shadow-sm"
    : "border border-border hover:border-primary-brand"

  return (
    <div
      role="button"
      tabIndex={0}
      className={`bg-card ${borderClass} rounded-[12px] p-6 flex items-center gap-6 w-full cursor-pointer transition-all duration-200 select-none`}
      onClick={() => onSelect(workspace)}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault()
          onSelect(workspace)
        }
      }}
    >
      <div className="bg-secondary border border-border rounded-full size-12 flex items-center justify-center shrink-0">
        <Icon name={resolveWorkspaceIconName(workspace.icon_key)} size="sm" className="text-text-secondary" />
      </div>

      <div className="flex-1 flex flex-col gap-1 min-w-0">
        <p className="text-base font-semibold leading-none text-foreground truncate">
          {workspace.name}
        </p>
        <span className="text-xs text-text-secondary font-mono truncate">
          {workspace.unique_key}
        </span>
      </div>

      <div className="flex items-center gap-3 shrink-0">
        <Badge variant={workspace.visibility === "private" ? "secondary" : "outline"} className="capitalize">
          {workspace.visibility || "private"}
        </Badge>

        <div className="bg-secondary rounded-lg size-10 flex items-center justify-center border border-border">
          <Icon
            name="chevron_right"
            size="sm"
            className={isSelected ? "text-primary-brand translate-x-0.5" : "text-text-secondary"}
          />
        </div>
      </div>
    </div>
  )
}
