import { useCallback, useEffect, useState } from "react"

import {
  clearInboxCredentialsCache,
  provisionInboxCredentials,
  readInboxCredentialsCache,
  saveInboxCredentialsCache,
  validateInboxCredentials,
} from "../inbox.api"
import type { InboxCredentials } from "../inbox.types"

async function resolveInboxCredentials(workspaceId: string): Promise<InboxCredentials> {
  const cached = readInboxCredentialsCache(workspaceId)
  if (cached) {
    const validation = await validateInboxCredentials(workspaceId, cached)
    if (validation.ok) {
      return cached
    }
    clearInboxCredentialsCache(workspaceId)
  }

  const creds = await provisionInboxCredentials(workspaceId)
  saveInboxCredentialsCache(workspaceId, creds)
  return creds
}

export function useInboxKey(workspaceId: string | undefined) {
  const [creds, setCreds] = useState<InboxCredentials | null>(null)
  const [isLoading, setIsLoading] = useState(() => !!workspaceId)
  const [error, setError] = useState<string | null>(null)

  const refreshCreds = useCallback(async (): Promise<InboxCredentials | null> => {
    if (!workspaceId) return null

    setIsLoading(true)
    setError(null)
    setCreds(null)
    clearInboxCredentialsCache(workspaceId)

    try {
      const nextCreds = await resolveInboxCredentials(workspaceId)
      setCreds(nextCreds)
      return nextCreds
    } catch {
      setError("Could not set up inbox access. Try again.")
      return null
    } finally {
      setIsLoading(false)
    }
  }, [workspaceId])

  useEffect(() => {
    let cancelled = false

    if (!workspaceId) {
      Promise.resolve().then(() => {
        if (!cancelled) {
          setCreds(null)
          setIsLoading(false)
          setError(null)
        }
      })
      return
    }

    const bootstrap = async () => {
      await Promise.resolve()
      setIsLoading(true)
      setError(null)

      try {
        const nextCreds = await resolveInboxCredentials(workspaceId)
        if (!cancelled) {
          setCreds(nextCreds)
        }
      } catch {
        if (!cancelled) {
          setError("Could not set up inbox access. Try again.")
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false)
        }
      }
    }

    void bootstrap()

    return () => {
      cancelled = true
    }
  }, [workspaceId])

  return { creds, isLoading, error, refreshCreds }
}
