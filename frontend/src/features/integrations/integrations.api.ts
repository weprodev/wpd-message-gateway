import { apiFetch } from "@/core/api/client"

import type { Integration, IntegrationChannel } from "./integrations.types"

export async function listIntegrations(workspaceId: string): Promise<Integration[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/integrations`)
  if (!res.ok) {
    throw new Error("Failed to load integrations")
  }
  return (await res.json()) as Integration[]
}

export async function upsertIntegration(
  workspaceId: string,
  body: {
    channel_type: IntegrationChannel
    provider_name: string
    config?: Record<string, unknown>
    status?: string
    is_default?: boolean
  },
): Promise<Integration> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/integrations`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    throw new Error("Failed to save integration")
  }
  return (await res.json()) as Integration
}

export async function deleteIntegration(workspaceId: string, integrationId: string): Promise<void> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/integrations/${integrationId}`, {
    method: "DELETE",
  })
  if (!res.ok) {
    throw new Error("Failed to delete integration")
  }
}

export interface ProviderConfigField {
  id: string
  provider_id: string
  key: string
  label: string
  description: string
  field_type: string
  required: boolean
  default_value: string
  options?: unknown
  sort_order: number
}

export async function fetchProviderConfigFields(
  workspaceId: string,
  providerName: string,
): Promise<ProviderConfigField[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/providers/${providerName}/config`)
  if (!res.ok) {
    throw new Error("Failed to load provider configuration fields")
  }
  return (await res.json()) as ProviderConfigField[]
}

