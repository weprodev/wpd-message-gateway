import { useEffect, useState } from "react"

import { deleteIntegration, listIntegrations, upsertIntegration, fetchProviders } from "../integrations.api"
import type { BackendProvider } from "../integrations.api"
import type { Integration, IntegrationActionResult, IntegrationChannel } from "../integrations.types"

export interface IntegrationViewModel {
  id: string
  name: string
  description: string
  icon: string
  category: IntegrationChannel
  isAvailable: boolean
  isComingSoon?: boolean
  integration?: Integration
  isConnected: boolean
  isDeactivated: boolean
}

function mergeCatalogWithIntegrations(providers: BackendProvider[], integrations: Integration[]): IntegrationViewModel[] {
  return providers.map((provider) => {
    const integration = integrations.find(
      (item) =>
        item.provider_name.toLowerCase() === provider.name.toLowerCase(),
    )
    return {
      id: provider.name.toLowerCase(),
      name: provider.name,
      description: provider.description,
      icon: provider.icon_path,
      category: provider.channel_type as IntegrationChannel,
      isAvailable: provider.status === "active",
      isComingSoon: provider.status === "not_supported",
      integration,
      isConnected: integration?.status === "connected",
      isDeactivated: integration?.status === "disconnected",
    }
  })
}

export function useIntegrations(workspaceId: string) {
  const [items, setItems] = useState<IntegrationViewModel[]>([])
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
        const [providers, integrations] = await Promise.all([
          fetchProviders(workspaceId),
          listIntegrations(workspaceId),
        ])
        if (cancelled) return
        setItems(mergeCatalogWithIntegrations(providers, integrations))
      } catch (err) {
        if (cancelled) return
        setError(err instanceof Error ? err.message : "Failed to load integrations")
      } finally {
        if (!cancelled) setIsLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [workspaceId, trigger])

  async function connect(
    provider: IntegrationViewModel,
    config: Record<string, unknown> = {},
  ): Promise<IntegrationActionResult> {
    const result = await upsertIntegration(workspaceId, {
      channel_type: provider.category,
      provider_name: provider.id,
      status: "connected",
      config: config,
    })
    if (!result.ok) return result
    reload()
    return { ok: true }
  }

  async function activate(provider: IntegrationViewModel): Promise<IntegrationActionResult> {
    if (!provider.integration) {
      return { ok: false, message: `No integration found for ${provider.name}` }
    }

    const result = await upsertIntegration(workspaceId, {
      channel_type: provider.category,
      provider_name: provider.id,
      status: "connected",
      config: provider.integration.config,
      is_default: provider.integration.is_default,
    })
    if (!result.ok) return result
    reload()
    return { ok: true }
  }

  async function deactivate(provider: IntegrationViewModel): Promise<IntegrationActionResult> {
    if (!provider.integration) {
      return { ok: false, message: `No integration found for ${provider.name}` }
    }

    const result = await upsertIntegration(workspaceId, {
      channel_type: provider.category,
      provider_name: provider.id,
      status: "disconnected",
      config: provider.integration.config,
      is_default: provider.integration.is_default,
    })
    if (!result.ok) return result
    reload()
    return { ok: true }
  }

  async function removeIntegration(provider: IntegrationViewModel): Promise<IntegrationActionResult> {
    if (!provider.integration) {
      return { ok: false, message: `No integration found for ${provider.name}` }
    }
    if (!provider.integration.id) {
      return { ok: false, message: `Cannot remove ${provider.name}: integration id is missing` }
    }

    const result = await deleteIntegration(workspaceId, provider.integration.id)
    if (!result.ok) return result
    reload()
    return { ok: true }
  }

  return { items, isLoading, error, reload, connect, activate, deactivate, removeIntegration }
}

export function filterIntegrationsByTab(
  items: IntegrationViewModel[],
  tab: "all" | "connected" | "available",
): IntegrationViewModel[] {
  if (tab === "connected") return items.filter((item) => item.isConnected)
  if (tab === "available") return items.filter((item) => item.isAvailable && !item.integration)
  return items
}

export function groupByCategory(
  items: IntegrationViewModel[],
): Record<IntegrationChannel, IntegrationViewModel[]> {
  return {
    email: items.filter((item) => item.category === "email"),
    sms: items.filter((item) => item.category === "sms"),
    push: items.filter((item) => item.category === "push"),
    chat: items.filter((item) => item.category === "chat"),
  }
}
