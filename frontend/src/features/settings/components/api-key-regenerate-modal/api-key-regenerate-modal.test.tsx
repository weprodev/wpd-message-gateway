import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ApiKeyRegenerateModal } from "./api-key-regenerate-modal"

const sampleCredentials = {
  clientId: "wk_test_client_id",
  clientSecret: "wk_test_client_secret",
  keyName: "Production",
  mode: "regenerated" as const,
}

describe("ApiKeyRegenerateModal", () => {
  it("renders confirmation message and actions when open", () => {
    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={vi.fn()}
      />,
    )

    expect(screen.getByText("Are you sure you want to regenerate this API key?")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Regenerate" })).toBeInTheDocument()
  })

  it("shows credentials in the same modal after regenerate succeeds", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn().mockResolvedValue(sampleCredentials)

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={onRegenerate}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Regenerate" }))

    expect(await screen.findByText("API key regenerated")).toBeInTheDocument()
    expect(screen.getByDisplayValue("wk_test_client_id")).toBeInTheDocument()
    expect(screen.queryByText("Are you sure you want to regenerate this API key?")).not.toBeInTheDocument()
  })

  it("calls onClose when cancel is clicked", async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={onClose}
        onRegenerate={vi.fn()}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Cancel" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })
})
