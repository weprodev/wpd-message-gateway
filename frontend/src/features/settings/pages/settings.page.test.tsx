import { render, screen } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { useSettings } from "../hooks/use-settings.hook"
import { SettingsPage } from "./settings.page"

vi.mock("../hooks/use-settings.hook", () => ({
  useSettings: vi.fn(),
}))

const mockedUseSettings = vi.mocked(useSettings)

const sampleApiKey = {
  id: "key-1",
  workspace_id: "w1",
  client_id: "demo-client-id-abcdefghijklmnop",
  name: "Production",
  is_active: true,
  created_at: "2026-01-01T00:00:00Z",
  last_used_at: null,
}

function renderSettingsPage(tab = "developer", permissions: string[] = [Permission.APIKeysWrite]) {
  return render(
    <MemoryRouter initialEntries={[`/workspaces/w1/settings?tab=${tab}`]}>
      <WorkspaceAuthorizationProvider role={Role.Admin} permissions={permissions}>
        <Routes>
          <Route path="/workspaces/:wid/settings" element={<SettingsPage />} />
        </Routes>
      </WorkspaceAuthorizationProvider>
    </MemoryRouter>,
  )
}

describe("SettingsPage", () => {
  it("shows loading state while settings are loading", () => {
    mockedUseSettings.mockReturnValue({
      settings: {},
      apiKeys: [],
      messageDispatchConfig: { mode: "memory", storeMessageContent: false },
      isLoading: true,
      error: null,
      reload: vi.fn(),
      saveSettings: vi.fn(),
      addApiKey: vi.fn(),
      removeApiKey: vi.fn(),
      rotateApiKey: vi.fn(),
    })

    renderSettingsPage()

    expect(screen.getByText(/loading settings/i)).toBeInTheDocument()
  })

  it("renders API keys tab content from useSettings", () => {
    mockedUseSettings.mockReturnValue({
      settings: {},
      apiKeys: [sampleApiKey],
      messageDispatchConfig: { mode: "memory", storeMessageContent: false },
      isLoading: false,
      error: null,
      reload: vi.fn(),
      saveSettings: vi.fn(),
      addApiKey: vi.fn(),
      removeApiKey: vi.fn(),
      rotateApiKey: vi.fn(),
    })

    renderSettingsPage("developer")

    expect(screen.getByRole("heading", { name: "API Keys" })).toBeInTheDocument()
    expect(screen.getByText("Production")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /generate key/i })).toBeInTheDocument()
  })

  it("hides generate key action for read-only users", () => {
    mockedUseSettings.mockReturnValue({
      settings: {},
      apiKeys: [sampleApiKey],
      messageDispatchConfig: { mode: "memory", storeMessageContent: false },
      isLoading: false,
      error: null,
      reload: vi.fn(),
      saveSettings: vi.fn(),
      addApiKey: vi.fn(),
      removeApiKey: vi.fn(),
      rotateApiKey: vi.fn(),
    })

    renderSettingsPage("developer", [Permission.APIKeysRead])

    expect(screen.queryByRole("button", { name: /generate key/i })).not.toBeInTheDocument()
  })
})
