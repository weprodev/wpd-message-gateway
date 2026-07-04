import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ApiKeyCreateDialog } from "./api-key-create-dialog"

describe("ApiKeyCreateDialog", () => {
  it("submits trimmed API key name", async () => {
    const user = userEvent.setup()
    const onCreate = vi.fn().mockResolvedValue(undefined)

    render(<ApiKeyCreateDialog open onClose={vi.fn()} onCreate={onCreate} />)

    await user.type(screen.getByLabelText(/api key name/i), "  Production  ")
    await user.click(screen.getByRole("button", { name: /generate key/i }))

    expect(onCreate).toHaveBeenCalledWith("Production")
  })

  it("shows validation error when name is empty", async () => {
    const user = userEvent.setup()

    render(<ApiKeyCreateDialog open onClose={vi.fn()} onCreate={vi.fn()} />)

    await user.click(screen.getByRole("button", { name: /generate key/i }))

    expect(screen.getByText(/api key name is required/i)).toBeInTheDocument()
  })
})
