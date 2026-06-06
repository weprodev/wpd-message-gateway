import { Link } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { Icon } from "@/components/ui/icon"
import { useWorkspaces } from "../hooks/use-workspaces.hook"

export function WorkspacesPage() {
  const { workspaces, isLoading, error } = useWorkspaces()

  if (isLoading) {
    return <p className="text-sm text-muted-foreground">Loading…</p>
  }

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Your workspaces</h1>
        <p className="text-muted-foreground">Workspaces you belong to.</p>
      </div>
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      {workspaces.length === 0 ? (
        <p className="text-sm text-muted-foreground">No workspaces yet.</p>
      ) : (
        <ul className="divide-y overflow-hidden rounded-md border bg-card">
          {workspaces.map((workspace) => (
            <li key={workspace.id} className="transition-colors hover:bg-muted/30">
              <Link
                to={ROUTES.workspace.overview(workspace.id)}
                className="flex w-full flex-col gap-1 px-4 py-4 sm:flex-row sm:items-center sm:justify-between"
              >
                <div>
                  <div className="font-semibold text-foreground">{workspace.name}</div>
                  <div className="mt-0.5 text-sm text-muted-foreground">
                    <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">
                      {workspace.unique_key}
                    </code>{" "}
                    · {workspace.status}
                  </div>
                </div>
                <div className="flex items-center gap-1 text-xs font-semibold text-primary">
                  Open
                  <Icon name="arrow_forward" size="sm" />
                </div>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
