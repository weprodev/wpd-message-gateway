import { useEffect, useState } from "react"

import {
  createApiKey,
  deleteApiKey,
  getSettings,
  listApiKeys,
  patchSettings,
  regenerateApiKey,
} from "../settings.api"
import type { ApiKey, RetentionMode, WorkspaceSettings } from "../settings.types"

export function useSettings(workspaceId: string) {
  const [settings, setSettings] = useState<WorkspaceSettings>({})
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [trigger, setTrigger] = useState(0)

  const reload = () => setTrigger((value) => value + 1)

  useEffect(() => {
    let cancelled = false

    ;(async () => {
      setIsLoading(true)
      setError(null)
      try {
        const [settingsData, keysData] = await Promise.all([
          getSettings(workspaceId),
          listApiKeys(workspaceId),
        ])
        if (cancelled) return
        setSettings(settingsData)
        setApiKeys(keysData)
      } catch (err) {
        if (cancelled) return
        setError(err instanceof Error ? err.message : "Failed to load settings")
      } finally {
        if (!cancelled) setIsLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [workspaceId, trigger])

  async function saveSettings(patch: Record<string, string>) {
    await patchSettings(workspaceId, patch)
    setSettings((prev) => {
      const next = { ...prev, ...patch }
      delete next.pin_code
      if (patch.pin_code) {
        next.pin_configured = "true"
      }
      return next
    })
  }

  async function addApiKey(name: string): Promise<ApiKey & { client_secret: string }> {
    const created = await createApiKey(workspaceId, name)
    const { client_secret: clientSecret, ...apiKey } = created
    setApiKeys((prev) => [...prev, apiKey])
    return { ...apiKey, client_secret: clientSecret }
  }

  async function removeApiKey(keyId: string) {
    await deleteApiKey(workspaceId, keyId)
    setApiKeys((prev) => prev.filter((key) => key.id !== keyId))
  }

  async function rotateApiKey(keyId: string): Promise<{ client_secret: string }> {
    return regenerateApiKey(workspaceId, keyId)
  }

  const retentionMode = (settings.data_retention as RetentionMode | undefined) ?? "memory"

  return {
    settings,
    apiKeys,
    retentionMode,
    isLoading,
    error,
    reload,
    saveSettings,
    addApiKey,
    removeApiKey,
    rotateApiKey,
  }
}
