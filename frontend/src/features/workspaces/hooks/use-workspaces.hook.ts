import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { fetchWorkspaces } from "../workspaces.api"
import type { Workspace } from "../workspace.types"

type UseWorkspacesOptions = {
  /** When set, redirects if the active workspace is missing from the list. */
  activeWorkspaceId?: string
}

export function useWorkspaces(options: UseWorkspacesOptions = {}) {
  const { activeWorkspaceId } = options
  const navigate = useNavigate()
  const [workspaces, setWorkspaces] = useState<Workspace[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [trigger, setTrigger] = useState(0)

  const reload = () => setTrigger((t) => t + 1)

  useEffect(() => {
    let cancelled = false

    ;(async () => {
      setIsLoading(true)
      setError(null)
      const result = await fetchWorkspaces()
      if (cancelled) return

      setIsLoading(false)
      if (!result.ok) {
        setError(result.message ?? "Failed to load workspaces")
        return
      }

      setWorkspaces(result.workspaces)

      if (!activeWorkspaceId) return

      const activeExists = result.workspaces.some((w) => w.id === activeWorkspaceId)
      if (!activeExists && result.workspaces.length > 0) {
        navigate(ROUTES.workspace.overview(result.workspaces[0].id), { replace: true })
      } else if (result.workspaces.length === 0) {
        navigate(ROUTES.workspaces, { replace: true })
      }
    })()

    return () => {
      cancelled = true
    }
  }, [activeWorkspaceId, navigate, trigger])

  const activeWorkspace =
    (workspaces || []).find((workspace) => workspace.id === activeWorkspaceId) ?? null

  return { workspaces: workspaces || [], activeWorkspace, isLoading, error, reload }
}
