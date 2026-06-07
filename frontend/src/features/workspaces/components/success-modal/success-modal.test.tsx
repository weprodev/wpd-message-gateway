import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { SuccessModal } from "./success-modal"

describe("SuccessModal Component", () => {
  it("renders correctly with workspace name", () => {
    render(<SuccessModal isOpen={true} workspaceName="My Workspace" onContinue={vi.fn()} />)
    expect(screen.getByText("Workspace Created!")).toBeInTheDocument()
    expect(screen.getByText(/My Workspace/)).toBeInTheDocument()
  })

  it("calls onContinue when button is clicked", () => {
    const handleContinue = vi.fn()
    render(<SuccessModal isOpen={true} workspaceName="My Workspace" onContinue={handleContinue} />)
    screen.getByRole("button", { name: "Go to Dashboard" }).click()
    expect(handleContinue).toHaveBeenCalledTimes(1)
  })
})
