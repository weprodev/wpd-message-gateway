import { apiFetch, setToken } from "@/core/api/client"

type AuthResponseBody = {
  token?: string
  message?: string
}

function splitDisplayName(displayName: string): { first_name: string; last_name: string } {
  const trimmed = displayName.trim()
  const spaceIdx = trimmed.indexOf(" ")
  if (spaceIdx === -1) {
    return { first_name: trimmed, last_name: "" }
  }
  return {
    first_name: trimmed.slice(0, spaceIdx),
    last_name: trimmed.slice(spaceIdx + 1).trim(),
  }
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
  const { first_name, last_name } = splitDisplayName(displayName)
  const res = await apiFetch("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password, first_name, last_name }),
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

