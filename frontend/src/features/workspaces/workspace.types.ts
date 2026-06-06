export interface Workspace {
  id: string
  name: string
  unique_key: string
  status: string
  admin_email?: string
  visibility?: "public" | "private"
  icon_key?: string
  created_at?: string
  updated_at?: string
}

export interface CreateWorkspaceData {
  name: string
  unique_key: string
  icon_key: string
}

export interface WorkspaceCreatedResponse {
  id: string
  name: string
  unique_key: string
  icon_key?: string
  created_at: string
}
