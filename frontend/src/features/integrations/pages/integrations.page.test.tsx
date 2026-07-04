import { render, screen } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { useIntegrations } from "../hooks/use-integrations.hook"
import { IntegrationsPage } from "./integrations.page"

vi.mock("../hooks/use-integrations.hook", () => ({
  useIntegrations: vi.fn(),
  filterIntegrationsByTab: vi.fn((items: unknown[]) => items),
  groupByCategory: vi.fn(() => ({})),
}))

const mockedUseIntegrations = vi.mocked(useIntegrations)

function renderPage() {
  return render(
    <MemoryRouter initialEntries={["/workspaces/ws-1/integrations"]}>
      <Routes>
        <Route path="/workspaces/:wid/integrations" element={<IntegrationsPage />} />
      </Routes>
    </MemoryRouter>,
  )
}

describe("IntegrationsPage", () => {
  it("shows loading state", () => {
    mockedUseIntegrations.mockReturnValue({
      items: [],
      isLoading: true,
      error: null,
      connect: vi.fn(),
      activate: vi.fn(),
      deactivate: vi.fn(),
      removeIntegration: vi.fn(),
    })

    renderPage()

    expect(screen.getByText(/loading integrations/i)).toBeInTheDocument()
  })

  it("shows error banner", () => {
    mockedUseIntegrations.mockReturnValue({
      items: [],
      isLoading: false,
      error: "Failed to load providers",
      connect: vi.fn(),
      activate: vi.fn(),
      deactivate: vi.fn(),
      removeIntegration: vi.fn(),
    })

    renderPage()

    expect(screen.getByText("Failed to load providers")).toBeInTheDocument()
    expect(screen.getByText("Integrations")).toBeInTheDocument()
  })
})
