import { useEffect, useState } from "react"
import { apiFetch } from "@/core/api/client"
import type { InboxCredentials } from "../inbox.types"

interface APIKeyItem {
  id: string
  workspace_id: string
  client_id: string
  name: string
  is_active: boolean
}

export function useInboxKey(workspaceId: string | undefined) {
  const [creds, setCreds] = useState<InboxCredentials | null>(null)
  const [isLoading, setIsLoading] = useState(() => !!workspaceId)
  const [error, setError] = useState<string | null>(null)

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

    const cacheKey = `wpd_inbox_creds_${workspaceId}`

    const bootstrap = async () => {
      await Promise.resolve() // yield to avoid synchronous state update in effect
      setIsLoading(true)
      setError(null)

      try {
        // 1. Check local storage cache
        const cached = localStorage.getItem(cacheKey)
        if (cached) {
          try {
            const parsed = JSON.parse(cached) as InboxCredentials
            if (parsed.clientId && parsed.clientSecret) {
              if (!cancelled) {
                setCreds(parsed)
                setIsLoading(false)
              }
              return
            }
          } catch {
            localStorage.removeItem(cacheKey)
          }
        }

        // 2. Fetch keys list from backend
        const listRes = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`)
        if (!listRes.ok) {
          throw new Error("Failed to retrieve workspace API keys")
        }
        const keys = (await listRes.json()) as APIKeyItem[]

        const portalKey = keys.find((k) => k.name === "Portal UI" && k.is_active)

        let finalCreds: InboxCredentials

        if (portalKey) {
          // Key exists but secret is hidden; regenerate to get new plaintext secret
          const regenRes = await apiFetch(
            `/api/v1/workspaces/${workspaceId}/api-keys/${portalKey.id}/regenerate`,
            { method: "POST" }
          )
          if (!regenRes.ok) {
            throw new Error("Failed to regenerate Portal UI API credentials")
          }
          const data = (await regenRes.json()) as { client_secret: string }
          finalCreds = {
            clientId: portalKey.client_id,
            clientSecret: data.client_secret,
          }
        } else {
          // No Portal UI key exists; create one
          const createRes = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`, {
            method: "POST",
            body: JSON.stringify({ name: "Portal UI" }),
          })
          if (!createRes.ok) {
            throw new Error("Failed to auto-provision Portal UI API credentials")
          }
          const data = (await createRes.json()) as { client_id: string; client_secret: string }
          finalCreds = {
            clientId: data.client_id,
            clientSecret: data.client_secret,
          }
        }

        // 3. Cache and update state
        localStorage.setItem(cacheKey, JSON.stringify(finalCreds))
        if (!cancelled) {
          setCreds(finalCreds)
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Authentication error")
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

  return { creds, isLoading, error }
}
