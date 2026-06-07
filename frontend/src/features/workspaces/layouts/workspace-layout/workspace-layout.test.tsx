import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"
import { MemoryRouter } from "react-router-dom"

import { WorkspaceLayout } from "./workspace-layout"

vi.mock("../../hooks/use-workspaces.hook", () => ({
  useWorkspaces: () => ({
    workspaces: [{ id: "w1", name: "Demo Workspace" }],
    activeWorkspace: { id: "w1", name: "Demo Workspace" },
  }),
}))

vi.mock("@/shared/context/theme-context", () => ({
  useTheme: () => ({
    theme: "light",
    toggleTheme: vi.fn(),
  }),
}))

vi.mock("@/core/api/client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/core/api/client")>()
  return {
    ...actual,
    fetchUserProfile: vi.fn().mockResolvedValue({
      id: "u1",
      first_name: "Demo",
      last_name: "User",
      email: "demo@weprodev.com",
      email_verified: true,
      created_at: "",
      updated_at: "",
    }),
  }
})

describe("WorkspaceLayout Component", () => {
  it("renders layout header, sidebar, and profile details", async () => {
    render(
      <MemoryRouter>
        <WorkspaceLayout />
      </MemoryRouter>
    )
    expect(screen.getByText("Message Gateway")).toBeInTheDocument()
    expect(screen.getByText("Demo Workspace")).toBeInTheDocument()
    expect(screen.getByText("Navigation")).toBeInTheDocument()

    // Wait for the mock user profile to be loaded and rendered
    const userEmail = await screen.findByText("demo@weprodev.com")
    expect(userEmail).toBeInTheDocument()
    expect(screen.getByText("Demo User")).toBeInTheDocument()
  })
})

