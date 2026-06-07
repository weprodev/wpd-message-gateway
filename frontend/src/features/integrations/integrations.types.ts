export type IntegrationChannel = "email" | "sms" | "push" | "chat"

export interface Integration {
  id: string
  workspace_id: string
  channel_type: IntegrationChannel
  provider_name: string
  config: Record<string, unknown>
  status: string
  is_default: boolean
  created_at: string
  updated_at: string
}
