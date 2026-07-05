import { apiFetch } from "@/core/api/client"

import type { Workspace } from "../workspace.types"

export async function fetchWorkspaces(): Promise<
  { ok: true; workspaces: Workspace[] } | { ok: false; status: number; message?: string }
> {
  const res = await apiFetch("/api/v1/workspaces")
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { message?: string }
    return { ok: false, status: res.status, message: err.message }
  }
  const workspaces = (await res.json()) as Workspace[] | null
  return { ok: true, workspaces: workspaces ?? [] }
}
