import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { CreateWorkspaceModal } from "./create-workspace-modal"

vi.mock("../../hooks/use-create-workspace.hook", () => ({
  useCreateWorkspace: () => ({
    createWorkspace: vi.fn(),
    isLoading: false,
    error: null,
  }),
}))

describe("CreateWorkspaceModal Component", () => {
  it("renders form fields correctly when open", () => {
    render(<CreateWorkspaceModal isOpen={true} onClose={vi.fn()} onSuccess={vi.fn()} />)
    expect(screen.getByLabelText("Workspace Name")).toBeInTheDocument()
    expect(screen.getByLabelText("Workspace Slug")).toBeInTheDocument()
  })
})
