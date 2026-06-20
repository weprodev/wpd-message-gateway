import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { SuccessModal } from "./success-modal"

describe("SuccessModal Component", () => {
  it("renders correctly with workspace name", () => {
    render(
      <SuccessModal
        isOpen={true}
        workspaceName="My Workspace"
        onClose={vi.fn()}
        onContinue={vi.fn()}
      />,
    )
    expect(screen.getByText("Workspace Created!")).toBeInTheDocument()
    expect(screen.getByText(/My Workspace/)).toBeInTheDocument()
  })

  it("renders joined variant copy", () => {
    render(
      <SuccessModal
        isOpen={true}
        workspaceName="Design Team"
        variant="joined"
        onClose={vi.fn()}
        onContinue={vi.fn()}
      />,
    )
    expect(screen.getByText("Workspace Joined!")).toBeInTheDocument()
    expect(screen.getByText(/Design Team/)).toBeInTheDocument()
  })

  it("calls onContinue when button is clicked", () => {
    const handleContinue = vi.fn()
    render(
      <SuccessModal
        isOpen={true}
        workspaceName="My Workspace"
        onClose={vi.fn()}
        onContinue={handleContinue}
      />,
    )
    screen.getByRole("button", { name: "Go to Dashboard" }).click()
    expect(handleContinue).toHaveBeenCalledTimes(1)
  })

  it("calls onClose when close button is clicked", () => {
    const handleClose = vi.fn()
    const handleContinue = vi.fn()
    render(
      <SuccessModal
        isOpen={true}
        workspaceName="My Workspace"
        onClose={handleClose}
        onContinue={handleContinue}
      />,
    )
    screen.getByRole("button", { name: /close/i }).click()
    expect(handleClose).toHaveBeenCalledTimes(1)
    expect(handleContinue).not.toHaveBeenCalled()
  })
})
