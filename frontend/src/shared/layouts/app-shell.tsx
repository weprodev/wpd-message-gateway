import { Link, NavLink, Outlet } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { cn } from "@/lib/utils"

export function AppShell() {
  return (
    <div className="min-h-screen bg-muted/30">
      <header className="border-b bg-card">
        <div className="mx-auto flex max-w-5xl items-center justify-between gap-4 px-4 py-4">
          <Link to={ROUTES.root} className="font-semibold text-primary">
            Message Gateway
          </Link>
          <nav className="flex items-center gap-4 text-sm text-muted-foreground">
            <NavLink
              to={ROUTES.workspaces}
              className={({ isActive }) =>
                cn("hover:text-foreground", isActive && "font-medium text-foreground")
              }
            >
              Workspaces
            </NavLink>
          </nav>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
