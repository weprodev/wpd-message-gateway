import { apiFetch } from "@/core/api/client"

import type { Integration, IntegrationActionResult, IntegrationChannel, IntegrationStatus } from "./integrations.types"

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
    status?: IntegrationStatus
    is_default?: boolean
  },
): Promise<{ ok: true; integration: Integration } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/integrations`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    })
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { message?: string }
      return { ok: false, message: err.message ?? "Failed to save integration" }
    }
    const integration = (await res.json()) as Integration
    return { ok: true, integration }
  } catch (err) {
    return {
      ok: false,
      message: err instanceof Error ? err.message : "Failed to save integration",
    }
  }
}

export async function deleteIntegration(
  workspaceId: string,
  integrationId: string,
): Promise<IntegrationActionResult> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/integrations/${integrationId}`, {
      method: "DELETE",
    })
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { message?: string }
      return { ok: false, message: err.message ?? "Failed to delete integration" }
    }
    return { ok: true }
  } catch (err) {
    return {
      ok: false,
      message: err instanceof Error ? err.message : "Failed to delete integration",
    }
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

export interface BackendProvider {
  id: string
  name: string
  channel_type: string
  status: string
  description: string
  icon_path: string
}

export async function fetchProviders(workspaceId: string): Promise<BackendProvider[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/providers`)
  if (!res.ok) {
    throw new Error("Failed to load providers catalog")
  }
  return (await res.json()) as BackendProvider[]
}

