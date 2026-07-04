import { apiFetch } from "@/core/api/client"

export interface GetStartedApiKey {
  client_id: string
  name: string
  is_active: boolean
}

export interface GetStartedContext {
  workspaceId: string
  apiKeys: GetStartedApiKey[]
}

export async function fetchGetStartedContext(
  workspaceId: string,
): Promise<{ ok: true; context: GetStartedContext } | { ok: false; message?: string }> {
  try {
    const keysRes = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`)
    if (!keysRes.ok) {
      return { ok: false, message: "Failed to load API keys" }
    }

    const apiKeys = (await keysRes.json()) as GetStartedApiKey[]
    return {
      ok: true,
      context: {
        workspaceId,
        apiKeys: apiKeys.filter((key) => key.is_active),
      },
    }
  } catch (err) {
    return {
      ok: false,
      message: err instanceof Error ? err.message : "Failed to load credentials",
    }
  }
}
