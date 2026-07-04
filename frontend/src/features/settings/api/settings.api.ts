import { fetchWorkspaceApiKeys } from "@/core/api/workspace-api-keys"
import { apiFetch } from "@/core/api/client"

import { parseWorkspaceSettings } from "../settings.schema"
import type { ApiKey, WorkspaceSettings, WorkspaceSettingsPatch } from "../settings.types"

async function readApiErrorMessage(res: Response, fallback: string): Promise<string> {
  try {
    const body = (await res.json()) as { message?: string }
    const message = body.message?.trim()
    return message || fallback
  } catch {
    return fallback
  }
}

export async function getSettings(workspaceId: string): Promise<WorkspaceSettings> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/settings`)
  if (!res.ok) {
    throw new Error("Failed to load settings")
  }
  return parseWorkspaceSettings(await res.json())
}

export async function patchSettings(workspaceId: string, body: WorkspaceSettingsPatch): Promise<void> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/settings`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    throw new Error(await readApiErrorMessage(res, "Failed to save settings"))
  }
}

export async function listApiKeys(workspaceId: string): Promise<ApiKey[]> {
  return fetchWorkspaceApiKeys(workspaceId)
}

export async function createApiKey(
  workspaceId: string,
  name: string,
): Promise<ApiKey & { client_secret?: string }> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name }),
  })
  if (!res.ok) {
    throw new Error("Failed to create API key")
  }
  return (await res.json()) as ApiKey & { client_secret?: string }
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
  if (!res.ok) {
    throw new Error("Failed to regenerate API key")
  }
  return (await res.json()) as { client_secret: string }
}
