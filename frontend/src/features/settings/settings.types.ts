export interface WorkspaceSettings {
  owner_email?: string
  pin_code?: string
  message_dispatch_mode?: MessageDispatchMode
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

/** Gateway dispatch mode values stored in workspace_settings.message_dispatch_mode */
export type MessageDispatchMode =
  | "memory_only"
  | "memory_and_provider"
  | "provider_only"
  | "provider_and_database"
