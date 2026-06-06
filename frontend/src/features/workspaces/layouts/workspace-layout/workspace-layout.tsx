import { useEffect, useRef, useState } from "react"
import { Link, Outlet, useParams } from "react-router-dom"

import { Icon } from "@/components/ui/icon"
import { SidebarNav } from "@/components/ui/sidebar-nav"
import { cn } from "@/lib/utils"
import { ROUTES } from "@/core/router/routes"
import { useTheme } from "@/shared/context/theme-context"
import { useWorkspaces } from "../../hooks/use-workspaces.hook"
import { WORKSPACE_NAV_ITEMS, workspaceHref } from "../workspace-nav.config"

export function WorkspaceLayout() {
  const { wid = "" } = useParams<{ wid: string }>()
  const { theme, toggleTheme } = useTheme()
  const { workspaces, activeWorkspace } = useWorkspaces({ activeWorkspaceId: wid })
  const [isDropdownOpen, setIsDropdownOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsDropdownOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])


  return (
    <div className="flex min-h-screen flex-col bg-muted/30 font-sans">
      <header className="sticky top-0 z-40 w-full border-b border-border/80 bg-background shadow-xs">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <Link to={ROUTES.workspaces} className="text-base font-bold tracking-tight text-foreground">
            Message Gateway
          </Link>

          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={toggleTheme}
              className="flex items-center justify-center size-9 rounded-lg border bg-card hover:bg-muted text-muted-foreground shadow-xs transition-colors shrink-0"
              aria-label="Toggle theme"
            >
              <Icon name={theme === "light" ? "dark_mode" : "light_mode"} size="sm" />
            </button>

            <div className="relative flex items-center gap-3" ref={dropdownRef}>
              <button
                type="button"
                onClick={() => setIsDropdownOpen((open) => !open)}
                className="flex cursor-pointer select-none items-center justify-between gap-3 rounded-lg border bg-card px-4 py-2 text-sm font-semibold text-foreground shadow-xs transition-colors hover:bg-muted/40"
              >
                <span>{activeWorkspace?.name ?? "Select workspace"}</span>
                <Icon name={isDropdownOpen ? "expand_less" : "expand_more"} size="sm" className="text-muted-foreground" />
              </button>

            {isDropdownOpen && (
              <div className="absolute right-0 top-full z-50 mt-1.5 w-52 animate-in fade-in slide-in-from-top-1 rounded-lg border bg-card py-1 shadow-lg duration-100">
                <div className="border-b px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
                  Switch workspace
                </div>
                {workspaces.map((workspace) => (
                  <Link
                    key={workspace.id}
                    to={ROUTES.workspace.overview(workspace.id)}
                    onClick={() => setIsDropdownOpen(false)}
                    className={cn(
                      "flex items-center justify-between px-4 py-2 text-xs font-semibold transition-colors hover:bg-muted/60",
                      workspace.id === wid ? "font-bold text-primary" : "text-foreground",
                    )}
                  >
                    <span>{workspace.name}</span>
                    {workspace.id === wid && <Icon name="check" size="sm" className="text-primary" />}
                  </Link>
                ))}
                <div className="my-1 border-t" />
                <Link
                  to={ROUTES.workspaces}
                  onClick={() => setIsDropdownOpen(false)}
                  className="flex items-center gap-2 px-4 py-2 text-xs font-semibold text-muted-foreground transition-colors hover:bg-muted/60"
                >
                  <Icon name="grid_view" size="sm" />
                  All workspaces
                </Link>
              </div>
            )}
          </div>
          </div>
        </div>
      </header>

      <div className="mx-auto flex w-full max-w-6xl flex-1 gap-6 px-6 py-8">
        <aside className="w-56 shrink-0">
          <SidebarNav
            items={WORKSPACE_NAV_ITEMS}
            workspaceId={wid}
            buildHref={workspaceHref}
          />
        </aside>

        <main className="min-w-0 flex-1">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
