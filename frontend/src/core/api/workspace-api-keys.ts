import { apiFetch } from "@/core/api/client"

export interface WorkspaceApiKey {
  id: string
  workspace_id: string
  client_id: string
  name: string
  is_active: boolean
  last_used_at?: string | null
  created_at: string
  expires_at?: string | null
}

export async function fetchWorkspaceApiKeys(workspaceId: string): Promise<WorkspaceApiKey[]> {
  const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`)
  if (!res.ok) {
    throw new Error("Failed to load API keys")
  }
  return (await res.json()) as WorkspaceApiKey[]
}

export function activeWorkspaceApiKeys(keys: WorkspaceApiKey[]): WorkspaceApiKey[] {
  return keys.filter((key) => key.is_active)
}
