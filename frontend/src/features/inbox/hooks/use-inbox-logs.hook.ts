import { useCallback, useEffect, useState } from "react"

import { fetchLogs } from "../inbox.api"
import type { LogRow } from "../inbox.types"

export function useInboxLogs(workspaceId: string | undefined, channel?: string) {
  const [logs, setLogs] = useState<LogRow[]>([])
  const [total, setTotal] = useState(0)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const applyResponse = useCallback(
    (res: Awaited<ReturnType<typeof fetchLogs>>) => {
      setIsLoading(false)
      if (!res.ok) {
        setError(res.message ?? "Failed to load request logs")
        return
      }
      setError(null)
      setLogs(res.items)
      setTotal(res.total)
    },
    [],
  )

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      if (!workspaceId) return
      setIsLoading(true)
      setError(null)
      const res = await fetchLogs(workspaceId, { channel, limit: 50 })
      if (cancelled) return
      applyResponse(res)
    })()
    return () => {
      cancelled = true
    }
  }, [workspaceId, channel, applyResponse])

  const reload = useCallback(async () => {
    if (!workspaceId) return
    setIsLoading(true)
    setError(null)
    const res = await fetchLogs(workspaceId, { channel, limit: 50 })
    applyResponse(res)
  }, [workspaceId, channel, applyResponse])

  return { logs, total, isLoading, error, reload }
}
