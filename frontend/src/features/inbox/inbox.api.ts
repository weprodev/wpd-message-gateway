import { apiFetch, getToken } from "@/core/api/client"
import type { EmailTemplate, InboxCredentials, LogRow, StoredEmail } from "./inbox.types"

const PORTAL_UI_KEY_NAME = "Portal UI"

export type InboxRequestFailure = {
  ok: false
  status: number
  unauthorized: boolean
  message: string
}

function inboxAuthHeaders(creds: InboxCredentials): HeadersInit {
  return {
    "X-Api-Client-Id": creds.clientId,
    "X-Api-Client-Secret": creds.clientSecret,
  }
}

/** Safe, user-facing inbox errors — never forward raw API bodies or network details. */
function inboxErrorForStatus(status: number): Pick<InboxRequestFailure, "message" | "unauthorized"> {
  if (status === 401) {
    return {
      message: "Inbox credentials are invalid or expired. Try again to reconnect.",
      unauthorized: true,
    }
  }
  if (status === 403) {
    return {
      message: "You do not have access to this workspace inbox.",
      unauthorized: false,
    }
  }
  if (status >= 500) {
    return {
      message: "The inbox service is temporarily unavailable. Try again later.",
      unauthorized: false,
    }
  }
  return {
    message: "Failed to fetch email inbox messages.",
    unauthorized: false,
  }
}

function inboxNetworkFailureMessage(action: "load" | "delete"): string {
  if (action === "delete") {
    return "Could not delete the email. Check your connection and try again."
  }
  return "Could not reach the inbox service. Check your connection and try again."
}

export function inboxCredentialsCacheKey(workspaceId: string): string {
  return `wpd_inbox_creds_${workspaceId}`
}

export function readInboxCredentialsCache(workspaceId: string): InboxCredentials | null {
  if (typeof localStorage === "undefined") return null
  const cached = localStorage.getItem(inboxCredentialsCacheKey(workspaceId))
  if (!cached) return null
  try {
    const parsed = JSON.parse(cached) as InboxCredentials
    if (parsed.clientId && parsed.clientSecret) {
      return parsed
    }
  } catch {
    clearInboxCredentialsCache(workspaceId)
  }
  return null
}

export function saveInboxCredentialsCache(workspaceId: string, creds: InboxCredentials): void {
  if (typeof localStorage === "undefined") return
  localStorage.setItem(inboxCredentialsCacheKey(workspaceId), JSON.stringify(creds))
}

export function clearInboxCredentialsCache(workspaceId: string): void {
  if (typeof localStorage === "undefined") return
  localStorage.removeItem(inboxCredentialsCacheKey(workspaceId))
}

interface WorkspaceAPIKeyItem {
  id: string
  client_id: string
  name: string
  is_active: boolean
}

export async function provisionInboxCredentials(workspaceId: string): Promise<InboxCredentials> {
  const listRes = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`)
  if (!listRes.ok) {
    throw new Error("INBOX_CREDENTIALS_UNAVAILABLE")
  }

  const keys = (await listRes.json()) as WorkspaceAPIKeyItem[]
  const portalKey = keys.find((key) => key.name === PORTAL_UI_KEY_NAME && key.is_active)

  if (portalKey) {
    const regenRes = await apiFetch(
      `/api/v1/workspaces/${workspaceId}/api-keys/${portalKey.id}/regenerate`,
      { method: "POST" }
    )
    if (!regenRes.ok) {
      throw new Error("INBOX_CREDENTIALS_UNAVAILABLE")
    }
    const data = (await regenRes.json()) as { client_secret: string }
    return {
      clientId: portalKey.client_id,
      clientSecret: data.client_secret,
    }
  }

  const createRes = await apiFetch(`/api/v1/workspaces/${workspaceId}/api-keys`, {
    method: "POST",
    body: JSON.stringify({ name: PORTAL_UI_KEY_NAME }),
  })
  if (!createRes.ok) {
    throw new Error("INBOX_CREDENTIALS_UNAVAILABLE")
  }
  const data = (await createRes.json()) as { client_id: string; client_secret: string }
  return {
    clientId: data.client_id,
    clientSecret: data.client_secret,
  }
}

export async function validateInboxCredentials(
  workspaceId: string,
  creds: InboxCredentials
): Promise<{ ok: true } | { ok: false; unauthorized: boolean }> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/stats`, {
      headers: inboxAuthHeaders(creds),
    })
    if (res.ok) {
      return { ok: true }
    }
    return { ok: false, unauthorized: res.status === 401 }
  } catch {
    return { ok: false, unauthorized: false }
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
): Promise<{ ok: true; items: StoredEmail[] } | InboxRequestFailure> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/emails`, {
      headers: inboxAuthHeaders(creds),
    })
    if (!res.ok) {
      return { ok: false, status: res.status, ...inboxErrorForStatus(res.status) }
    }
    const data = (await res.json()) as StoredEmail[]
    const sorted = [...data].sort(
      (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )
    return { ok: true, items: sorted }
  } catch {
    return {
      ok: false,
      status: 0,
      unauthorized: false,
      message: inboxNetworkFailureMessage("load"),
    }
  }
}

export async function deleteInboxEmail(
  workspaceId: string,
  messageId: string,
  creds: InboxCredentials
): Promise<{ ok: true } | InboxRequestFailure> {
  try {
    const res = await apiFetch(`/api/v1/workspaces/${workspaceId}/inbox/emails/${messageId}`, {
      method: "DELETE",
      headers: inboxAuthHeaders(creds),
    })
    if (!res.ok) {
      const failure = inboxErrorForStatus(res.status)
      return {
        ok: false,
        status: res.status,
        message: res.status === 401 ? failure.message : "Could not delete the email. Try again.",
        unauthorized: failure.unauthorized,
      }
    }
    return { ok: true }
  } catch {
    return {
      ok: false,
      status: 0,
      unauthorized: false,
      message: inboxNetworkFailureMessage("delete"),
    }
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
