import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { JoinWorkspaceModal } from "./join-workspace-modal"

describe("JoinWorkspaceModal Component", () => {
  it("renders form fields correctly when open", () => {
    render(<JoinWorkspaceModal isOpen={true} onClose={vi.fn()} onSuccess={vi.fn()} />)
    expect(screen.getByLabelText("Workspace Slug")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Join Workspace" })).toBeInTheDocument()
  })
})
