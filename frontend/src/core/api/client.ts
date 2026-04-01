const TOKEN_KEY = "portal_token"
const WS_API_PREFIX = "portal_ws_api_"

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string | null): void {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

export type WorkspaceApiCredentials = {
  clientId: string
  clientSecret: string
}

export function setWorkspaceApiKey(workspaceId: string, cred: WorkspaceApiCredentials | null): void {
  const k = WS_API_PREFIX + workspaceId
  if (cred) localStorage.setItem(k, JSON.stringify(cred))
  else localStorage.removeItem(k)
}

export function getWorkspaceApiKey(workspaceId: string): WorkspaceApiCredentials | null {
  const raw = localStorage.getItem(WS_API_PREFIX + workspaceId)
  if (!raw) return null
  try {
    return JSON.parse(raw) as WorkspaceApiCredentials
  } catch {
    return null
  }
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

export async function apiFetchWorkspace(workspaceId: string, path: string, init?: RequestInit): Promise<Response> {
  const headers = new Headers(init?.headers)
  if (!headers.has("Content-Type") && init?.body) {
    headers.set("Content-Type", "application/json")
  }
  const token = getToken()
  if (token) headers.set("Authorization", `Bearer ${token}`)
  const cred = getWorkspaceApiKey(workspaceId)
  if (cred) {
    headers.set("X-Api-Client-Id", cred.clientId)
    headers.set("X-Api-Client-Secret", cred.clientSecret)
  }
  const base = `/api/v1/workspaces/${workspaceId}/inbox`
  const url = path.startsWith("http") ? path : `${base}${path.startsWith("/") ? path : "/" + path}`
  return fetch(url, { ...init, headers })
}
