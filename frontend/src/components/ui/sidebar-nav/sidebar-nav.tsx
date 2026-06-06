import { NavLink } from "react-router-dom"
import { cn } from "@/lib/utils"
import { Icon } from "@/components/ui/icon"

export interface SidebarNavItem<T extends string = string> {
  readonly label: string
  readonly segment: T
  readonly icon: string
}

interface SidebarNavProps<T extends string = string> extends React.HTMLAttributes<HTMLElement> {
  items: readonly SidebarNavItem<T>[]
  workspaceId: string
  buildHref: (wid: string, segment: T) => string
}

export function SidebarNav<T extends string = string>({
  items,
  workspaceId,
  buildHref,
  className,
  ...props
}: SidebarNavProps<T>) {
  return (
    <nav
      className={cn(
        "flex flex-col gap-1 rounded-xl border border-border bg-card p-4 shadow-xs",
        className
      )}
      {...props}
    >
      <span className="mb-1 px-3 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
        Channels
      </span>
      {items.map((item) => (
        <NavLink
          key={item.segment}
          to={buildHref(workspaceId, item.segment)}
          className={({ isActive }) =>
            cn(
              "flex select-none items-center gap-3 rounded-lg px-4 py-2.5 text-sm font-semibold transition-all duration-150",
              isActive
                ? "bg-primary/10 font-bold text-primary"
                : "text-muted-foreground hover:bg-muted/50 hover:text-foreground",
            )
          }
        >
          <Icon name={item.icon} size="sm" />
          {item.label}
        </NavLink>
      ))}
    </nav>
  )
}
