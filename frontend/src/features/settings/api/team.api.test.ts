import { describe, expect, it, vi } from "vitest"

import { listInvitations, listMembers } from "./team.api"

vi.mock("@/core/api/client", () => ({
  apiFetch: vi.fn(),
}))

import { apiFetch } from "@/core/api/client"

const mockedApiFetch = vi.mocked(apiFetch)

describe("team.api", () => {
  it("listMembers returns empty array when API responds with null", async () => {
    mockedApiFetch.mockResolvedValue({
      ok: true,
      json: async () => null,
    } as Response)

    await expect(listMembers("ws-1")).resolves.toEqual([])
  })

  it("listInvitations returns empty array when API responds with null", async () => {
    mockedApiFetch.mockResolvedValue({
      ok: true,
      json: async () => null,
    } as Response)

    await expect(listInvitations("ws-1")).resolves.toEqual([])
  })
})
