import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { WorkspaceActions } from "./workspace-actions"

describe("WorkspaceActions Component", () => {
  it("renders buttons correctly", () => {
    render(<WorkspaceActions />)
    expect(screen.getByRole("button", { name: "Create Workspace" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Join Workspace" })).toBeInTheDocument()
  })

  it("calls handlers on click", () => {
    const handleCreate = vi.fn()
    const handleJoin = vi.fn()
    render(<WorkspaceActions onCreateWorkspace={handleCreate} onJoinWorkspace={handleJoin} />)
    
    screen.getByRole("button", { name: "Create Workspace" }).click()
    screen.getByRole("button", { name: "Join Workspace" }).click()

    expect(handleCreate).toHaveBeenCalledTimes(1)
    expect(handleJoin).toHaveBeenCalledTimes(1)
  })
})
