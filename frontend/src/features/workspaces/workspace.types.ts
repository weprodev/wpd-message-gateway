export interface Workspace {
  id: string
  name: string
  slug: string
  status: string
  admin_email?: string
  visibility?: "public" | "private"
  icon_key?: string
  created_at?: string
  updated_at?: string
  role?: string
  permissions?: string[]
}

export interface CreateWorkspaceData {
  name: string
  slug: string
  icon_key: string
  pin: string
}

export interface WorkspaceCreatedResponse {
  id: string
  name: string
  slug: string
  icon_key?: string
  created_at: string
}
