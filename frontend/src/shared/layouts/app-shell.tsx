import { useState } from "react"
import { Link, NavLink, Outlet, useNavigate } from "react-router-dom"

import { getToken, setToken } from "@/core/api/client"
import { ROUTES } from "@/app/paths"
import { cn } from "@/lib/utils"

export function AppShell() {
  const navigate = useNavigate()
  const [authed, setAuthed] = useState(() => !!getToken())

  function logout() {
    setToken(null)
    setAuthed(false)
    navigate(ROUTES.login, { replace: true })
  }

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
            <NavLink
              to={ROUTES.email}
              className={({ isActive }) =>
                cn("hover:text-foreground", isActive && "font-medium text-foreground")
              }
            >
              Email
            </NavLink>
            {authed ? (
              <button
                type="button"
                onClick={logout}
                className="text-foreground hover:underline"
              >
                Sign out
              </button>
            ) : (
              <Link to={ROUTES.login} className="hover:text-foreground">
                Sign in
              </Link>
            )}
          </nav>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
