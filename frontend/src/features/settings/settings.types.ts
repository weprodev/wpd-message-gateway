export interface WorkspaceSettings {
  owner_email?: string
  pin_code?: string
  message_dispatch_mode?: MessageDispatchApiValue
  [key: string]: string | undefined
}

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

export type SettingsTab = "general" | "developer" | "team" | "retention"

export type MessageDispatchMode = "memory" | "provider"

export interface MessageDispatchConfig {
  mode: MessageDispatchMode
  storeInDb: boolean
}

/** Gateway dispatch mode values stored in workspace_settings.message_dispatch_mode */
export type MessageDispatchApiValue =
  | "memory_only"
  | "memory_and_database"
  | "provider_only"
  | "provider_and_database"
