import { useState } from "react"
import { apiFetch } from "@/core/api/client"
import type { CreateWorkspaceData, Workspace } from "../workspace.types"

export function useCreateWorkspace() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const createWorkspace = async (data: CreateWorkspaceData): Promise<Workspace> => {
    setIsLoading(true)
    setError(null)

    try {
      const res = await apiFetch("/api/v1/workspaces", {
        method: "POST",
        body: JSON.stringify({
          name: data.name,
          unique_key: data.unique_key,
          icon_key: data.icon_key,
        }),
      })

      if (!res.ok) {
        const err = (await res.json().catch(() => ({}))) as { message?: string }
        throw new Error(err.message ?? "Failed to create workspace")
      }

      const response = (await res.json()) as Workspace
      return response
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Failed to create workspace"
      setError(msg)
      throw new Error(msg, { cause: err })
    } finally {
      setIsLoading(false)
    }
  }

  return {
    createWorkspace,
    isLoading,
    error,
  }
}
