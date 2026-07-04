import { useCallback, useEffect, useState } from "react"

import { fetchGetStartedContext, type GetStartedContext } from "../api/get-started.api"

export function useGetStartedContext(workspaceId: string | undefined, enabled: boolean) {
  const [context, setContext] = useState<GetStartedContext | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const applyResult = useCallback(
    (result: Awaited<ReturnType<typeof fetchGetStartedContext>>) => {
      setIsLoading(false)
      if (!result.ok) {
        setContext(null)
        setError(result.message ?? "Failed to load workspace credentials")
        return
      }
      setError(null)
      setContext(result.context)
    },
    [],
  )

  useEffect(() => {
    if (!enabled || !workspaceId) {
      return
    }

    let cancelled = false
    ;(async () => {
      setIsLoading(true)
      setError(null)
      const result = await fetchGetStartedContext(workspaceId)
      if (cancelled) return
      applyResult(result)
    })()

    return () => {
      cancelled = true
    }
  }, [workspaceId, enabled, applyResult])

  return { context, isLoading, error }
}
