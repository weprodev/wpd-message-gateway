export interface WorkspaceSettings {
  owner_email?: string
  pin_code?: string
  data_retention?: "memory" | "both" | "providers"
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

export type ApiKeyCredentialsMode = "created" | "regenerated"

export interface ApiKeyCredentials {
  clientId: string
  clientSecret: string
  keyName: string
  mode: ApiKeyCredentialsMode
}

export type SettingsTab = "general" | "developer" | "team" | "retention"
export type RetentionMode = "memory" | "both" | "providers"
