import { render, screen } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { useInboxLogs } from "../hooks/use-inbox-logs.hook"
import { OverviewPage } from "./overview.page"

vi.mock("../hooks/use-inbox-logs.hook", () => ({
  useInboxLogs: vi.fn(),
}))

const mockedUseInboxLogs = vi.mocked(useInboxLogs)

function renderPage() {
  return render(
    <MemoryRouter initialEntries={["/workspaces/ws-1/overview"]}>
      <Routes>
        <Route path="/workspaces/:wid/overview" element={<OverviewPage />} />
      </Routes>
    </MemoryRouter>,
  )
}

describe("OverviewPage", () => {
  it("shows loading state", () => {
    mockedUseInboxLogs.mockReturnValue({
      logs: [],
      total: 0,
      isLoading: true,
      error: null,
      reload: vi.fn(),
    })

    renderPage()

    expect(screen.getByText(/loading request logs/i)).toBeInTheDocument()
  })

  it("renders log rows", () => {
    mockedUseInboxLogs.mockReturnValue({
      logs: [
        {
          id: "log-1",
          workspace_id: "ws-1",
          channel_type: "email",
          http_method: "POST",
          status_code: 200,
          endpoint: "/v1/email",
          provider_name: "memory",
          created_at: "2026-07-04T10:00:00Z",
          source_name: "Demo Key",
          client_id: "client-1",
        },
      ],
      total: 1,
      isLoading: false,
      error: null,
      reload: vi.fn(),
    })

    renderPage()

    expect(screen.getByText("/v1/email")).toBeInTheDocument()
    expect(screen.getByText("200")).toBeInTheDocument()
  })
})
