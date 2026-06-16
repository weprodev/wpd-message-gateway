import { apiFetch } from "@/core/api/client"
import { httpError, requireClientSecret } from "@/lib/errors"

import type { ApiKey, RetentionMode, WorkspaceSettings } from "./settings.types"

const DISPATCH_TO_RETENTION: Record<string, RetentionMode> = {
  memory_only: "memory",
  memory_and_provider: "both",
  provider_only: "providers",
  provider_and_database: "provider_database",
}

function mapSettingsFromApi(raw: Record<string, string>): WorkspaceSettings {
  const dataRetention =
    (raw.data_retention as RetentionMode | undefined) ??
    (raw.message_dispatch_mode ? DISPATCH_TO_RETENTION[raw.message_dispatch_mode] : undefined)

  return {
    ...raw,
    data_retention: dataRetention ?? "memory",
  }
}

async function ensureOk(response: Response, fallback: string): Promise<void> {
  if (!response.ok) {
    throw await httpError(response, fallback)
  }
}

export async function getSettings(workspaceId: string): Promise<WorkspaceSettings> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/settings`)
  if (!res.ok) {
    throw new Error("Failed to load settings")
  }
  return mapSettingsFromApi((await res.json()) as Record<string, string>)
}

export async function patchSettings(
  workspaceId: string,
  body: Record<string, string>,
): Promise<void> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/settings`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    throw new Error("Failed to save settings")
  }
}

export async function listApiKeys(workspaceId: string): Promise<ApiKey[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`)
  if (!res.ok) {
    throw new Error("Failed to load API keys")
  }
  return (await res.json()) as ApiKey[]
}

export async function createApiKey(
  workspaceId: string,
  name: string,
): Promise<ApiKey & { client_secret: string }> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name }),
  })
  await ensureOk(res, "Failed to create API key")

  const created = (await res.json()) as ApiKey & { client_secret?: string }

  return { ...created, client_secret: requireClientSecret(created) }
}

export async function deleteApiKey(workspaceId: string, keyId: string): Promise<void> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys/${keyId}`, {
    method: "DELETE",
  })
  if (!res.ok) {
    throw new Error("Failed to delete API key")
  }
}

export async function regenerateApiKey(
  workspaceId: string,
  keyId: string,
): Promise<{ client_secret: string }> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys/${keyId}/regenerate`, {
    method: "POST",
  })
  await ensureOk(res, "Failed to regenerate API key")

  const regenerated = (await res.json()) as { client_secret?: string }

  return { client_secret: requireClientSecret(regenerated) }
}
