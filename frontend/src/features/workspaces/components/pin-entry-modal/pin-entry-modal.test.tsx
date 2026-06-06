import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { PinEntryModal } from "./pin-entry-modal"

describe("PinEntryModal Component", () => {
  it("renders correctly with workspace name", () => {
    render(
      <PinEntryModal
        isOpen={true}
        workspaceName="Secret Project"
        onClose={vi.fn()}
        onSubmit={vi.fn()}
      />
    )
    expect(screen.getByText(/Secret Project/)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Unlock Workspace" })).toBeInTheDocument()
  })

  it("calls onClose when Cancel button is clicked", () => {
    const handleClose = vi.fn()
    render(
      <PinEntryModal
        isOpen={true}
        workspaceName="Secret Project"
        onClose={handleClose}
        onSubmit={vi.fn()}
      />
    )
    screen.getByRole("button", { name: "Cancel" }).click()
    expect(handleClose).toHaveBeenCalledTimes(1)
  })
})
