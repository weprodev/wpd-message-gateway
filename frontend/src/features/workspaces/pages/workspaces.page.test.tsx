import { render, screen } from "@testing-library/react"
import { MemoryRouter } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { useWorkspaces } from "../hooks/use-workspaces.hook"
import { WorkspacesPage } from "./workspaces.page"

vi.mock("../hooks/use-workspaces.hook", () => ({
  useWorkspaces: vi.fn(),
}))

const mockedUseWorkspaces = vi.mocked(useWorkspaces)

describe("WorkspacesPage", () => {
  it("shows loading spinner while workspaces load", () => {
    mockedUseWorkspaces.mockReturnValue({
      workspaces: [],
      isLoading: true,
      error: null,
      reload: vi.fn(),
    })

    render(
      <MemoryRouter>
        <WorkspacesPage />
      </MemoryRouter>,
    )

    expect(screen.getByText(/loading workspaces/i)).toBeInTheDocument()
  })

  it("renders workspace cards", () => {
    mockedUseWorkspaces.mockReturnValue({
      workspaces: [
        {
          id: "ws-1",
          name: "Demo Workspace",
          slug: "demo",
          role: "admin",
          status: "active",
          visibility: "public",
          created_at: "2026-01-01T00:00:00Z",
        },
      ],
      isLoading: false,
      error: null,
      reload: vi.fn(),
    })

    render(
      <MemoryRouter>
        <WorkspacesPage />
      </MemoryRouter>,
    )

    expect(screen.getByText("Demo Workspace")).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/search workspaces/i)).toBeInTheDocument()
  })
})
