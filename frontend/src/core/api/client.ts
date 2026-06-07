const TOKEN_KEY = "portal_token"

export function getToken(): string | null {
  if (typeof localStorage === "undefined") return null
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string | null): void {
  if (typeof localStorage === "undefined") return
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

export async function apiFetch(path: string, init?: RequestInit): Promise<Response> {
  const headers = new Headers(init?.headers)
  if (!headers.has("Content-Type") && init?.body) {
    headers.set("Content-Type", "application/json")
  }
  const token = getToken()
  if (token) headers.set("Authorization", `Bearer ${token}`)
  return fetch(path, { ...init, headers })
}

export interface UserWorkspace {
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

export interface UserProfile {
  id: string
  first_name: string
  last_name: string
  email: string
  email_verified: boolean
  created_at: string
  updated_at: string
  workspaces?: UserWorkspace[]
}

export async function fetchUserProfile(): Promise<UserProfile | null> {
  const res = await apiFetch("/api/v1/auth/me")
  if (!res.ok) return null
  return res.json() as Promise<UserProfile>
}

