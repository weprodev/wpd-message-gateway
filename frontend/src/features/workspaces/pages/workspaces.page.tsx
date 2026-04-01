import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"

import { getToken } from "@/core/api/client"
import { ROUTES } from "@/app/paths"
import { Button } from "@/components/ui/button"
import { fetchWorkspaces } from "@/features/workspaces/workspaces.api"
import type { Workspace } from "@/features/workspaces/workspace.types"

export function WorkspacesPage() {
  const navigate = useNavigate()
  const [list, setList] = useState<Workspace[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!getToken()) {
      navigate(ROUTES.login, { replace: true })
      return
    }
    let cancelled = false
    ;(async () => {
      const result = await fetchWorkspaces()
      if (cancelled) return
      if (!result.ok) {
        if (result.status === 401) {
          navigate(ROUTES.login, { replace: true })
          return
        }
        setError(result.message ?? "Failed to load workspaces")
        return
      }
      setList(result.workspaces)
    })()
    return () => {
      cancelled = true
    }
  }, [navigate])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Your workspaces</h1>
        <p className="text-muted-foreground">Workspaces you belong to (from the portal API).</p>
      </div>
      <div className="flex flex-wrap gap-2">
        <Button type="button" variant="outline" disabled title="Use POST /api/v1/workspaces from API or Bruno">
          Create workspace
        </Button>
        <Button type="button" variant="outline" disabled title="Use POST /api/v1/workspaces/join">
          Join with PIN
        </Button>
      </div>
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
      {list === null && !error ? (
        <p className="text-sm text-muted-foreground">Loading…</p>
      ) : null}
      {list && list.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          No workspaces yet. Create one via the API or seed data, then add your user as a member.
        </p>
      ) : null}
      {list && list.length > 0 ? (
        <ul className="divide-y rounded-md border bg-card">
          {list.map((w) => (
            <li key={w.id} className="flex flex-col gap-1 px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <div className="font-medium">{w.name}</div>
                <div className="text-sm text-muted-foreground">
                  <code className="rounded bg-muted px-1">{w.unique_key}</code> · {w.status}
                </div>
              </div>
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  )
}
