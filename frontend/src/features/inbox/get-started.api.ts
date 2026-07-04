import { activeWorkspaceApiKeys, fetchWorkspaceApiKeys } from "@/core/api/workspace-api-keys"
import type { WorkspaceApiKey } from "@/core/api/workspace-api-keys"

export interface GetStartedContext {
  workspaceId: string
  apiKeys: Array<Pick<WorkspaceApiKey, "client_id" | "name" | "is_active">>
}

export async function fetchGetStartedContext(
  workspaceId: string,
): Promise<{ ok: true; context: GetStartedContext } | { ok: false; message?: string }> {
  try {
    const apiKeys = activeWorkspaceApiKeys(await fetchWorkspaceApiKeys(workspaceId))
    return {
      ok: true,
      context: {
        workspaceId,
        apiKeys,
      },
    }
  } catch (err) {
    return {
      ok: false,
      message: err instanceof Error ? err.message : "Failed to load credentials",
    }
  }
}
