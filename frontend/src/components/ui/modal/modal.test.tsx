import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { Modal } from "./modal"

describe("Modal Component", () => {
  it("does not render when isOpen is false", () => {
    render(<Modal isOpen={false} onClose={vi.fn()}>Content</Modal>)
    expect(screen.queryByText("Content")).not.toBeInTheDocument()
  })

  it("renders when isOpen is true", () => {
    render(<Modal isOpen={true} onClose={vi.fn()} title="Modal Title">Content</Modal>)
    expect(screen.getByText("Modal Title")).toBeInTheDocument()
    expect(screen.getByText("Content")).toBeInTheDocument()
  })

  it("calls onClose when close button is clicked", () => {
    const handleClose = vi.fn()
    render(<Modal isOpen={true} onClose={handleClose} title="Modal Title">Content</Modal>)
    const closeButton = screen.getByRole("button", { name: /close/i })
    closeButton.click()
    expect(handleClose).toHaveBeenCalledTimes(1)
  })
})
