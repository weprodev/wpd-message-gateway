export type MessageChannel = "email" | "sms" | "push" | "chat"

export type LogRow = {
  id: string
  workspace_id: string
  api_key_id?: string
  channel_type: string
  http_method: string
  status_code: number
  endpoint: string
  provider_name?: string
  request_id?: string
  duration_ms?: number
  error_message?: string
  created_at: string
  source_name?: string
}

export interface StoredEmail {
  id: string
  workspace_id?: string
  created_at: string
  email: {
    from?: string
    from_name?: string
    to: string[]
    cc?: string[]
    bcc?: string[]
    reply_to?: string
    subject: string
    html?: string
    plain_text?: string
    attachments?: unknown[]
    headers?: Record<string, string>
  }
}

export type InboxCredentials = {
  clientId: string
  clientSecret: string
}

export type EmailTemplate = {
  id: string
  workspace_id: string
  name: string
  unique_key: string
  channel_type: string
  category: string
  subject?: string
  content_html: string
  is_active: boolean
  is_default: boolean
  created_at: string
  updated_at: string
}
