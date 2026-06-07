import { NavLink } from "react-router-dom"

import { cn } from "@/lib/utils"
import { Icon } from "@/components/ui/icon"

export interface SidebarNavItem<T extends string = string> {
  readonly label: string
  readonly segment: T
  readonly icon: string
  readonly disabled?: boolean
}

export interface SidebarNavSection<T extends string = string> {
  readonly label?: string
  readonly items: readonly SidebarNavItem<T>[]
}

interface SidebarNavProps<T extends string = string> extends React.HTMLAttributes<HTMLElement> {
  sections: readonly SidebarNavSection<T>[]
  workspaceId: string
  buildHref: (wid: string, segment: T) => string
}

export function SidebarNav<T extends string = string>({
  sections,
  workspaceId,
  buildHref,
  className,
  ...props
}: SidebarNavProps<T>) {
  return (
    <nav
      className={cn(
        "flex w-[260px] flex-col gap-2 rounded-2xl border border-border bg-card p-4",
        className,
      )}
      {...props}
    >
      {sections.map((section, sectionIndex) => (
        <div key={section.label ?? sectionIndex} className="flex flex-col gap-2">
          {section.label ? (
            <span className="text-xs leading-4 text-text-tertiary">{section.label}</span>
          ) : sectionIndex > 0 ? (
            <div className="h-4" aria-hidden="true" />
          ) : null}
          {section.items.map((item) =>
            item.disabled ? (
              <span
                key={item.segment}
                className="flex h-11 w-full cursor-not-allowed items-center gap-3 rounded-[10px] px-3 opacity-50"
              >
                <Icon name={item.icon} size="sm" className="text-text-secondary" />
                <span className="text-sm font-medium text-text-secondary">{item.label}</span>
              </span>
            ) : (
              <NavLink
                key={item.segment}
                to={buildHref(workspaceId, item.segment)}
                className={({ isActive }) =>
                  cn(
                    "flex h-11 w-full items-center gap-3 rounded-[10px] px-3 text-sm font-medium transition-colors",
                    isActive
                      ? "bg-input text-primary-brand [&_.material-symbols-outlined]:text-indigo-500"
                      : "text-text-secondary hover:bg-input [&_.material-symbols-outlined]:text-text-secondary",
                  )
                }
              >
                <Icon name={item.icon} size="sm" />
                {item.label}
              </NavLink>
            ),
          )}
        </div>
      ))}
    </nav>
  )
}
