import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { WorkspaceCard } from "./workspace-card"

const mockWorkspace = {
  id: "w1",
  name: "Production Workspace",
  unique_key: "production_workspace",
  icon_key: "shield",
  visibility: "private" as const,
  status: "active",
  created_at: "2026-06-06T12:00:00Z",
  updated_at: "2026-06-06T12:00:00Z",
}

describe("WorkspaceCard Component", () => {
  it("renders correctly with workspace details", () => {
    render(<WorkspaceCard workspace={mockWorkspace} isSelected={false} onSelect={vi.fn()} />)
    expect(screen.getByText("Production Workspace")).toBeInTheDocument()
    expect(screen.getByText("production_workspace")).toBeInTheDocument()
    expect(screen.getByText("private")).toBeInTheDocument()
  })

  it("triggers onSelect callback on click", () => {
    const handleSelect = vi.fn()
    render(<WorkspaceCard workspace={mockWorkspace} isSelected={false} onSelect={handleSelect} />)
    screen.getByRole("button").click()
    expect(handleSelect).toHaveBeenCalledWith(mockWorkspace)
  })
})
