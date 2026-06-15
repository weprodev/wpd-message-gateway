import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ApiKeyDeleteModal } from "./api-key-delete-modal"

describe("ApiKeyDeleteModal", () => {
  it("renders confirmation message and actions when open", () => {
    render(
      <ApiKeyDeleteModal
        isOpen
        onCancel={vi.fn()}
        onDelete={vi.fn()}
      />,
    )

    expect(screen.getByText("Are you sure you want to delete this API key?")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument()
  })

  it("calls the appropriate handler for cancel and delete", async () => {
    const user = userEvent.setup()
    const onCancel = vi.fn()
    const onDelete = vi.fn()

    render(
      <ApiKeyDeleteModal
        isOpen
        onCancel={onCancel}
        onDelete={onDelete}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Cancel" }))
    expect(onCancel).toHaveBeenCalledTimes(1)
    expect(onDelete).not.toHaveBeenCalled()

    await user.click(screen.getByRole("button", { name: "Delete" }))
    expect(onDelete).toHaveBeenCalledTimes(1)
  })
})
