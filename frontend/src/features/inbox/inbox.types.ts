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
