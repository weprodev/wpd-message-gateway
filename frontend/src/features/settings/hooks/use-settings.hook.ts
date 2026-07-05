import { useEffect, useState } from "react"

import {
  createApiKey,
  deleteApiKey,
  getSettings,
  listApiKeys,
  patchSettings,
  regenerateApiKey,
} from "../api/settings.api"
import { parseMessageDispatchConfig } from "../message-dispatch-mode"
import { parseWorkspaceSettings } from "../settings.schema"
import type { ApiKey, WorkspaceSettings, WorkspaceSettingsPatch } from "../settings.types"

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

  async function saveSettings(patch: WorkspaceSettingsPatch) {
    await patchSettings(workspaceId, patch)
    setSettings((prev) => parseWorkspaceSettings({ ...prev, ...patch }))
    setError(null)
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

  const messageDispatchConfig = parseMessageDispatchConfig(settings.message_dispatch_mode, settings.store_message_content)

  return {
    settings,
    apiKeys,
    messageDispatchConfig,
    isLoading,
    error,
    reload,
    saveSettings,
    addApiKey,
    removeApiKey,
    rotateApiKey,
  }
}
