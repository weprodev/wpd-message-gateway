import { useEffect, useState } from "react"

import { deleteIntegration, listIntegrations, upsertIntegration, fetchProviders } from "@/features/integrations/api/integrations.api"
import type { BackendProvider } from "@/features/integrations/api/integrations.api"
import {
  INTEGRATION_STATUS,
  type Integration,
  type IntegrationActionResult,
  type IntegrationChannel,
  type IntegrationStatus,
  type IntegrationViewModel,
} from "@/features/integrations/integrations.types"

export type { IntegrationViewModel } from "@/features/integrations/integrations.types"

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
      isConnected: integration?.status === INTEGRATION_STATUS.CONNECTED,
      isDeactivated: integration?.status === INTEGRATION_STATUS.DISCONNECTED,
    }
  })
}

function integrationNotFound(provider: IntegrationViewModel): IntegrationActionResult {
  return { ok: false, message: `No integration found for ${provider.name}` }
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

  async function saveIntegrationStatus(
    provider: IntegrationViewModel,
    status: IntegrationStatus,
    config: Record<string, unknown>,
  ): Promise<IntegrationActionResult> {
    const result = await upsertIntegration(workspaceId, {
      channel_type: provider.category,
      provider_name: provider.id,
      status,
      config,
      ...(provider.integration ? { is_default: provider.integration.is_default } : {}),
    })
    if (!result.ok) return result
    reload()
    return { ok: true }
  }

  async function connect(
    provider: IntegrationViewModel,
    config: Record<string, unknown> = {},
  ): Promise<IntegrationActionResult> {
    return saveIntegrationStatus(provider, INTEGRATION_STATUS.CONNECTED, config)
  }

  async function activate(provider: IntegrationViewModel): Promise<IntegrationActionResult> {
    if (!provider.integration) return integrationNotFound(provider)
    return saveIntegrationStatus(provider, INTEGRATION_STATUS.CONNECTED, provider.integration.config)
  }

  async function deactivate(provider: IntegrationViewModel): Promise<IntegrationActionResult> {
    if (!provider.integration) return integrationNotFound(provider)
    return saveIntegrationStatus(provider, INTEGRATION_STATUS.DISCONNECTED, provider.integration.config)
  }

  async function removeIntegration(provider: IntegrationViewModel): Promise<IntegrationActionResult> {
    if (!provider.integration) return integrationNotFound(provider)
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
