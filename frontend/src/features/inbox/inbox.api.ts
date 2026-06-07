import { apiFetch, getToken } from "@/core/api/client"
import type { EmailTemplate, InboxCredentials, LogRow, StoredEmail } from "./inbox.types"

function inboxAuthHeaders(creds: InboxCredentials): HeadersInit {
  return {
    "X-Api-Client-Id": creds.clientId,
    "X-Api-Client-Secret": creds.clientSecret,
  }
}

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

export async function fetchInboxEmails(
  workspaceId: string,
  creds: InboxCredentials
): Promise<{ ok: true; items: StoredEmail[] } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/emails`, {
      headers: inboxAuthHeaders(creds),
    })
    if (!res.ok) {
      return { ok: false, message: "Failed to fetch email inbox messages" }
    }
    const data = (await res.json()) as StoredEmail[]
    const sorted = [...data].sort(
      (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )
    return { ok: true, items: sorted }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to load emails" }
  }
}

export async function deleteInboxEmail(
  workspaceId: string,
  messageId: string,
  creds: InboxCredentials
): Promise<{ ok: true } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/emails/${messageId}`, {
      method: "DELETE",
      headers: inboxAuthHeaders(creds),
    })
    if (!res.ok) {
      return { ok: false, message: "Failed to delete email message" }
    }
    return { ok: true }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to delete email" }
  }
}

export function buildInboxEventsUrl(
  workspaceId: string,
  creds: InboxCredentials
): string {
  const token = getToken()
  const params = new URLSearchParams({
    access_token: token ?? "",
    client_id: creds.clientId,
    client_secret: creds.clientSecret,
  })
  return `/api/v1/workspaces/${workspaceId}/inbox/events?${params.toString()}`
}

export async function fetchEmailTemplates(
  workspaceId: string
): Promise<{ ok: true; items: EmailTemplate[] } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/templates`)
    if (!res.ok) {
      return { ok: false, message: "Failed to load templates" }
    }
    const data = (await res.json()) as EmailTemplate[]
    return { ok: true, items: data ?? [] }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to load templates" }
  }
}

export async function deleteEmailTemplate(
  workspaceId: string,
  templateId: string
): Promise<{ ok: true } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/templates/${templateId}`, {
      method: "DELETE",
    })
    if (!res.ok) {
      return { ok: false, message: "Failed to delete template" }
    }
    return { ok: true }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to delete template" }
  }
}

export async function createEmailTemplate(
  workspaceId: string,
  payload: {
    name: string
    unique_key: string
    category: string
    subject: string
    content_html: string
  }
): Promise<{ ok: true; item: EmailTemplate } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/templates`, {
      method: "POST",
      body: JSON.stringify({
        name: payload.name,
        unique_key: payload.unique_key,
        channel_type: "email",
        category: payload.category,
        subject: payload.subject,
        content_html: payload.content_html,
        is_active: true,
        is_default: false,
      }),
    })
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { message?: string }
      return { ok: false, message: err.message ?? "Failed to save template" }
    }
    const item = (await res.json()) as EmailTemplate
    return { ok: true, item }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to create template" }
  }
}
