import { renderHook, waitFor } from "@testing-library/react"
import { beforeEach, describe, expect, it, vi } from "vitest"

vi.mock("@/features/integrations/api/integrations.api", () => ({
  fetchProviderConfigFields: vi.fn(),
}))

import { fetchProviderConfigFields } from "@/features/integrations/api/integrations.api"

import { useProviderConfigFields } from "./use-provider-config-fields.hook"

const mockedFetchProviderConfigFields = vi.mocked(fetchProviderConfigFields)

const mockFields = [
  {
    id: "f1",
    provider_id: "mailgun",
    key: "api_key",
    label: "API Key",
    description: "Your Mailgun Private API Key",
    field_type: "password",
    required: true,
    default_value: "default-secret",
    sort_order: 1,
  },
]

describe("useProviderConfigFields", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("loads fields and initializes form data when enabled", async () => {
    mockedFetchProviderConfigFields.mockResolvedValue(mockFields)

    const { result } = renderHook(() => useProviderConfigFields("ws-1", "mailgun", true))

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(mockedFetchProviderConfigFields).toHaveBeenCalledWith("ws-1", "mailgun")
    expect(result.current.fields).toEqual(mockFields)
    expect(result.current.formData).toEqual({ api_key: "default-secret" })
    expect(result.current.error).toBeNull()
  })

  it("does not fetch when disabled or provider id is missing", () => {
    renderHook(() => useProviderConfigFields("ws-1", undefined, true))
    renderHook(() => useProviderConfigFields("ws-1", "mailgun", false))

    expect(mockedFetchProviderConfigFields).not.toHaveBeenCalled()
  })
})
