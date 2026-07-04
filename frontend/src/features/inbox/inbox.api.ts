import { apiFetch, getToken } from "@/core/api/client"
import type { EmailTemplate, InboxEmailPage, LogRow, StoredEmail } from "./inbox.types"
import { parseLogsResponse } from "./logs.schema"

async function parseApiError(res: Response, fallback: string): Promise<string> {
  try {
    const err = (await res.json()) as { message?: string; error?: string }
    return err.message ?? err.error ?? fallback
  } catch {
    return fallback
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
      return { ok: false, message: await parseApiError(res, "Failed to fetch logs") }
    }
    const parsed = parseLogsResponse(await res.json())
    if (!parsed.ok) {
      return { ok: false, message: parsed.message }
    }
    return { ok: true, items: parsed.items, total: parsed.total }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to fetch logs" }
  }
}

export async function fetchInboxEmails(
  workspaceId: string,
  params: { limit?: number; cursor?: string } = {}
): Promise<{ ok: true; page: InboxEmailPage } | { ok: false; message?: string }> {
  try {
    const query = new URLSearchParams()
    if (params.limit !== undefined) {
      query.set("limit", String(params.limit))
    }
    if (params.cursor) {
      query.set("cursor", params.cursor)
    }
    const qs = query.toString()
    const path = `/api/v1/workspaces/${workspaceId}/inbox/emails${qs ? `?${qs}` : ""}`
    const res = await apiFetch(path)
    if (!res.ok) {
      return { ok: false, message: await parseApiError(res, "Failed to fetch email inbox messages") }
    }
    const raw = await res.json()
    const page: InboxEmailPage = Array.isArray(raw)
      ? { items: raw, has_more: false }
      : {
          items: raw.items ?? [],
          next_cursor: raw.next_cursor,
          has_more: raw.has_more ?? false,
        }
    return { ok: true, page }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to load emails" }
  }
}

export async function fetchInboxEmailById(
  workspaceId: string,
  messageId: string
): Promise<{ ok: true; item: StoredEmail } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/emails/${messageId}`)
    if (!res.ok) {
      return { ok: false, message: "Failed to fetch email" }
    }
    const item = (await res.json()) as StoredEmail
    return { ok: true, item }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to fetch email" }
  }
}

export async function deleteInboxEmail(
  workspaceId: string,
  messageId: string
): Promise<{ ok: true } | { ok: false; message?: string }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/emails/${messageId}`, {
      method: "DELETE",
    })
    if (!res.ok) {
      return { ok: false, message: "Failed to delete email message" }
    }
    return { ok: true }
  } catch (err) {
    return { ok: false, message: err instanceof Error ? err.message : "Failed to delete email" }
  }
}

export function buildInboxEventsUrl(workspaceId: string): string {
  const token = getToken()
  const params = new URLSearchParams({
    access_token: token ?? "",
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
