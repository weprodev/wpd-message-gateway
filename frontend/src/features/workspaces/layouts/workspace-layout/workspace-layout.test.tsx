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

describe("WorkspaceLayout Component", () => {
  it("renders layout header and sidebar", () => {
    render(
      <MemoryRouter>
        <WorkspaceLayout />
      </MemoryRouter>
    )
    expect(screen.getByText("Message Gateway")).toBeInTheDocument()
    expect(screen.getByText("Demo Workspace")).toBeInTheDocument()
    expect(screen.getByText("Navigation")).toBeInTheDocument()
  })
})
