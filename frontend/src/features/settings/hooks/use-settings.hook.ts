import { useEffect, useState } from "react"

import {
  createApiKey,
  deleteApiKey,
  getSettings,
  listApiKeys,
  patchSettings,
  regenerateApiKey,
} from "../settings.api"
import type { ApiKey, MessageDispatchMode, WorkspaceSettings } from "../settings.types"

const DISPATCH_MODES: MessageDispatchMode[] = [
  "memory_only",
  "memory_and_provider",
  "provider_only",
  "provider_and_database",
]

function normalizeMessageDispatchMode(raw: string | undefined): MessageDispatchMode {
  if (!raw) {
    return "memory_only"
  }

  if (DISPATCH_MODES.includes(raw as MessageDispatchMode)) {
    return raw as MessageDispatchMode
  }

  switch (raw) {
    case "memory":
      return "memory_only"
    case "both":
    case "memory_database":
      return "memory_and_provider"
    case "providers":
    case "provider":
      return "provider_only"
    case "provider_database":
      return "provider_and_database"
    default:
      return "memory_only"
  }
}

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
    setSettings((prev) => ({ ...prev, ...patch }))
  }

  async function addApiKey(name: string) {
    const created = await createApiKey(workspaceId, name)
    setApiKeys((prev) => [...prev, created])
    return created
  }

  async function removeApiKey(keyId: string) {
    await deleteApiKey(workspaceId, keyId)
    setApiKeys((prev) => prev.filter((key) => key.id !== keyId))
  }

  async function rotateApiKey(keyId: string) {
    return regenerateApiKey(workspaceId, keyId)
  }

  const messageDispatchMode = normalizeMessageDispatchMode(settings.message_dispatch_mode)

  return {
    settings,
    apiKeys,
    messageDispatchMode,
    isLoading,
    error,
    reload,
    saveSettings,
    addApiKey,
    removeApiKey,
    rotateApiKey,
  }
}
