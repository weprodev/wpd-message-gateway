import { describe, expect, it, vi } from "vitest"

vi.mock("@/core/api/client", () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from "@/core/api/client"

import {
  activeWorkspaceApiKeys,
  fetchWorkspaceApiKeys,
  type WorkspaceApiKey,
} from "./workspace-api-keys"

const mockedApiFetch = vi.mocked(apiFetch)

const sampleKeys: WorkspaceApiKey[] = [
  {
    id: "key-1",
    workspace_id: "ws-1",
    client_id: "client-active",
    name: "Active key",
    is_active: true,
    created_at: "2026-01-01T00:00:00Z",
  },
  {
    id: "key-2",
    workspace_id: "ws-1",
    client_id: "client-inactive",
    name: "Inactive key",
    is_active: false,
    created_at: "2026-01-02T00:00:00Z",
  },
]

describe("activeWorkspaceApiKeys", () => {
  it("filters out inactive keys", () => {
    const result = activeWorkspaceApiKeys(sampleKeys)

    expect(result).toHaveLength(1)
    expect(result[0]?.name).toBe("Active key")
    expect(result.every((key) => key.is_active)).toBe(true)
  })
})

describe("fetchWorkspaceApiKeys", () => {
  it("returns keys from the API", async () => {
    mockedApiFetch.mockResolvedValue({
      ok: true,
      json: async () => sampleKeys,
    } as Response)

    const result = await fetchWorkspaceApiKeys("ws-1")

    expect(mockedApiFetch).toHaveBeenCalledWith("/api/v1/workspaces/ws-1/api-keys")
    expect(result).toEqual(sampleKeys)
  })

  it("throws when the API response is not ok", async () => {
    mockedApiFetch.mockResolvedValue({ ok: false } as Response)

    await expect(fetchWorkspaceApiKeys("ws-1")).rejects.toThrow("Failed to load API keys")
  })
})
