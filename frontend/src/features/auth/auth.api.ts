import { apiFetch, setToken } from "@/core/api/client"

type AuthResponseBody = {
  token?: string
  message?: string
}

export async function signInWithPassword(
  email: string,
  password: string,
): Promise<{ ok: true } | { ok: false; message: string }> {
  const res = await apiFetch("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  })
  const data = (await res.json()) as AuthResponseBody
  if (!res.ok) return { ok: false, message: data.message ?? "Login failed" }
  if (data.token) setToken(data.token)
  return { ok: true }
}

export async function registerAccount(
  email: string,
  password: string,
  displayName: string,
): Promise<{ ok: true } | { ok: false; message: string }> {
  const res = await apiFetch("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password, display_name: displayName }),
  })
  const data = (await res.json()) as AuthResponseBody
  if (!res.ok) {
    return {
      ok: false,
      message: typeof data.message === "string" ? data.message : "Registration failed",
    }
  }
  if (data.token) setToken(data.token)
  return { ok: true }
}
