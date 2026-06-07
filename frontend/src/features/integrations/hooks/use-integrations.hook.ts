import { useEffect, useState } from "react"

import { deleteIntegration, listIntegrations, upsertIntegration } from "../integrations.api"
import type { Integration, IntegrationChannel, ProviderCatalogItem } from "../integrations.types"
import { PROVIDER_CATALOG } from "../integrations.types"

export interface IntegrationViewModel extends ProviderCatalogItem {
  integration?: Integration
  isConnected: boolean
}

function mergeCatalogWithIntegrations(integrations: Integration[]): IntegrationViewModel[] {
  return PROVIDER_CATALOG.map((provider) => {
    const integration = integrations.find(
      (item) =>
        item.provider_name.toLowerCase() === provider.id ||
        item.provider_name.toLowerCase() === provider.name.toLowerCase(),
    )
    return {
      ...provider,
      integration,
      isConnected: integration?.status === "connected",
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
        const integrations = await listIntegrations(workspaceId)
        if (cancelled) return
        setItems(mergeCatalogWithIntegrations(integrations))
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

  async function connect(provider: IntegrationViewModel) {
    await upsertIntegration(workspaceId, {
      channel_type: provider.category,
      provider_name: provider.id,
      status: "connected",
      config: {},
    })
    reload()
  }

  async function disconnect(provider: IntegrationViewModel) {
    if (provider.integration?.id) {
      await deleteIntegration(workspaceId, provider.integration.id)
    }
    reload()
  }

  return { items, isLoading, error, reload, connect, disconnect }
}

export function filterIntegrationsByTab(
  items: IntegrationViewModel[],
  tab: "all" | "connected" | "available",
): IntegrationViewModel[] {
  if (tab === "connected") return items.filter((item) => item.isConnected)
  if (tab === "available") return items.filter((item) => item.isAvailable && !item.isConnected)
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
