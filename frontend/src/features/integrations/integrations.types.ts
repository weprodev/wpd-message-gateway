export type IntegrationChannel = "email" | "sms" | "push" | "chat"

export const INTEGRATION_STATUS = {
  CONNECTED: "connected",
  DISCONNECTED: "disconnected",
} as const

export type IntegrationStatus = (typeof INTEGRATION_STATUS)[keyof typeof INTEGRATION_STATUS]

export type IntegrationActionResult =
  | { ok: true }
  | { ok: false; message?: string }

export interface Integration {
  id: string
  workspace_id: string
  channel_type: IntegrationChannel
  provider_name: string
  config: Record<string, unknown>
  status: IntegrationStatus
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface IntegrationViewModel {
  id: string
  name: string
  description: string
  icon: string
  category: IntegrationChannel
  isAvailable: boolean
  isComingSoon?: boolean
  integration?: Integration
  isConnected: boolean
  isDeactivated: boolean
}
