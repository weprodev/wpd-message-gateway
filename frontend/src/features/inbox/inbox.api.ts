import { apiFetch } from "@/core/api/client"
import type { LogRow } from "./inbox.types"

export async function fetchLogs(
  workspaceId: string,
  params: { channel?: string; limit?: number; offset?: number } = {}
): Promise<{ ok: true; items: LogRow[]; total: number } | { ok: false; message?: string }> {
  const query = new URLSearchParams()
  if (params.channel) {
    query.set("channel", params.channel)
  }
  if (params.limit !== undefined) {
    query.set("limit", String(params.limit))
  }
  if (params.offset !== undefined) {
    query.set("offset", String(params.offset))
  }

  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/logs?${query.toString()}`)
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { message?: string }
      return { ok: false, message: err.message ?? "Failed to fetch logs" }
    }
    const data = (await res.json()) as { items: LogRow[]; total: number }
    return { ok: true, items: data.items || [], total: data.total || 0 }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to fetch logs" }
  }
}

export async function sendTestRequest(
  workspaceId: string,
  channel: "email" | "sms" | "push" | "chat",
  payload: Record<string, unknown>
): Promise<{ ok: true; id: string } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/send-test/${channel}`, {
      method: "POST",
      body: JSON.stringify(payload),
    })
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { message?: string }
      return { ok: false, message: err.message ?? `Failed to send test ${channel}` }
    }
    const data = (await res.json()) as { id: string }
    return { ok: true, id: data.id }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : `Failed to send test ${channel}` }
  }
}
