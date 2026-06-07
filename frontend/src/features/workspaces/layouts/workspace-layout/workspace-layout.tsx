import { useEffect, useRef, useState } from "react"
import { Link, Outlet, useParams } from "react-router-dom"

import { Icon } from "@/components/ui/icon"
import { SidebarNav } from "@/components/ui/sidebar-nav"
import { cn } from "@/lib/utils"
import { ROUTES } from "@/core/router/routes"
import { AppFooter } from "@/shared/components/app-footer"
import { MessageGatewayLogo } from "@/shared/components/message-gateway-logo"
import { ThemeToggle } from "@/shared/components/theme-toggle"
import { useWorkspaces } from "../../hooks/use-workspaces.hook"
import { WORKSPACE_NAV_SECTIONS, workspaceHref } from "../workspace-nav.config"

export function WorkspaceLayout() {
  const { wid = "" } = useParams<{ wid: string }>()
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
    <div className="flex min-h-screen flex-col bg-background font-sans">
      <header className="sticky top-0 z-40 h-20 w-full border-b border-border bg-card">
        <div className="flex h-full items-center justify-between px-12 py-4">
          <div className="flex items-center gap-3">
            <MessageGatewayLogo />
            <div className="flex flex-col gap-1">
              <p className="text-xl font-semibold leading-normal text-primary-brand">
                Message Gateway
              </p>
              <p className="text-xs leading-4 text-text-secondary">
                Connecting your world seamlessly!
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <ThemeToggle />

            <div className="relative" ref={dropdownRef}>
              <button
                type="button"
                onClick={() => setIsDropdownOpen((open) => !open)}
                className="flex cursor-pointer select-none items-center gap-2 rounded-lg border border-border bg-input px-3 py-2 text-sm font-medium text-text-secondary transition-colors hover:bg-secondary-button-hover"
              >
                <span>{activeWorkspace?.name ?? "Select workspace"}</span>
                <Icon
                  name={isDropdownOpen ? "expand_less" : "expand_more"}
                  size="sm"
                  className="text-text-secondary"
                />
              </button>

              {isDropdownOpen ? (
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
                        workspace.id === wid ? "font-bold text-primary-brand" : "text-foreground",
                      )}
                    >
                      <span>{workspace.name}</span>
                      {workspace.id === wid ? (
                        <Icon name="check" size="sm" className="text-primary-brand" />
                      ) : null}
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
              ) : null}
            </div>
          </div>
        </div>
      </header>

      <main className="flex min-h-0 flex-1 gap-6 px-12 py-8">
        <aside className="w-[260px] shrink-0 self-start sticky top-[112px]">
          <SidebarNav
            sections={WORKSPACE_NAV_SECTIONS}
            workspaceId={wid}
            buildHref={workspaceHref}
          />
        </aside>

        <div className="flex min-w-0 flex-1 flex-col">
          <Outlet />
        </div>
      </main>

      <AppFooter variant="dashboard" />
    </div>
  )
}
