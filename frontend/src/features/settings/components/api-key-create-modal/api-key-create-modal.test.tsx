import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ApiKeyCreateModal } from "./api-key-create-modal"

const sampleCredentials = {
  clientId: "wk_test_client_id",
  clientSecret: "wk_test_client_secret",
  keyName: "Production",
  mode: "created" as const,
}

describe("ApiKeyCreateModal", () => {
  it("renders name field and actions when open", () => {
    render(
      <ApiKeyCreateModal
        isOpen
        onClose={vi.fn()}
        onCreate={vi.fn()}
      />,
    )

    expect(screen.getByLabelText("Name")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Create" })).toBeInTheDocument()
  })

  it("requires a name before creating", async () => {
    const user = userEvent.setup()
    const onCreate = vi.fn()

    render(
      <ApiKeyCreateModal
        isOpen
        onClose={vi.fn()}
        onCreate={onCreate}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Create" }))
    expect(screen.getByText("API key name is required")).toBeInTheDocument()
    expect(onCreate).not.toHaveBeenCalled()
  })

  it("shows credentials in the same modal after create succeeds", async () => {
    const user = userEvent.setup()
    const onCreate = vi.fn().mockResolvedValue(sampleCredentials)

    render(
      <ApiKeyCreateModal
        isOpen
        onClose={vi.fn()}
        onCreate={onCreate}
      />,
    )

    await user.type(screen.getByLabelText("Name"), "Production")
    await user.click(screen.getByRole("button", { name: "Create" }))

    expect(await screen.findByText("API key created")).toBeInTheDocument()
    expect(screen.getByDisplayValue("wk_test_client_id")).toBeInTheDocument()
    expect(screen.getByDisplayValue("wk_test_client_secret")).toBeInTheDocument()
    expect(screen.queryByLabelText("Name")).not.toBeInTheDocument()
  })

  it("calls onClose when cancel is clicked", async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()

    render(
      <ApiKeyCreateModal
        isOpen
        onClose={onClose}
        onCreate={vi.fn()}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Cancel" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })
})
