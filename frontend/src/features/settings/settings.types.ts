export type MessageDispatchMode = "memory" | "provider"

export type StoreMessageContentSetting = "true" | "false"

export interface WorkspaceSettings {
  owner_email?: string
  pin_code?: string
  message_dispatch_mode?: MessageDispatchMode
  store_message_content?: StoreMessageContentSetting
}

export type WorkspaceSettingsPatch = Partial<Record<keyof WorkspaceSettings, string>>

export interface ApiKey {
  id: string
  workspace_id: string
  client_id: string
  name: string
  is_active: boolean
  last_used_at?: string | null
  created_at: string
  expires_at?: string | null
}

export type SettingsTab = "general" | "developer" | "team" | "dispatch"

export interface MessageDispatchConfig {
  mode: MessageDispatchMode
  storeMessageContent: boolean
}
