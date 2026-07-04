import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ApiKeySecretDialog } from "./api-key-secret-dialog"

describe("ApiKeySecretDialog", () => {
  it("renders client credentials and closes on done", async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()

    render(
      <ApiKeySecretDialog
        open
        clientId="demo-client-id"
        clientSecret="secret-value"
        onClose={onClose}
      />,
    )

    expect(screen.getByDisplayValue("demo-client-id")).toBeInTheDocument()
    expect(screen.getByDisplayValue("secret-value")).toBeInTheDocument()

    await user.click(screen.getByRole("button", { name: /done/i }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it("copies client secret to clipboard", async () => {
    const user = userEvent.setup()
    const writeText = vi.fn().mockResolvedValue(undefined)
    vi.spyOn(navigator.clipboard, "writeText").mockImplementation(writeText)

    render(
      <ApiKeySecretDialog
        open
        clientId="demo-client-id"
        clientSecret="secret-value"
        onClose={vi.fn()}
      />,
    )

    const copyButtons = screen.getAllByRole("button", { name: /copy/i })
    await user.click(copyButtons[1])

    expect(writeText).toHaveBeenCalledWith("secret-value")
    expect(screen.getByRole("button", { name: /copied/i })).toBeInTheDocument()
  })
})
