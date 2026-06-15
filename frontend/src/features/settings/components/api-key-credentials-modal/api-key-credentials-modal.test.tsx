import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ApiKeyCredentialsModal } from "./api-key-credentials-modal"

const sampleCredentials = {
  clientId: "wk_test_client_id",
  clientSecret: "wk_test_client_secret",
  keyName: "Production",
  mode: "created" as const,
}

describe("ApiKeyCredentialsModal", () => {
  it("renders credentials and warning when open", () => {
    render(
      <ApiKeyCredentialsModal
        isOpen
        credentials={sampleCredentials}
        onConfirm={vi.fn()}
      />,
    )

    expect(screen.getByText("API key created")).toBeInTheDocument()
    expect(screen.getByDisplayValue("wk_test_client_id")).toBeInTheDocument()
    expect(screen.getByDisplayValue("wk_test_client_secret")).toBeInTheDocument()
    expect(screen.getByText(/only once/i)).toBeInTheDocument()
  })

  it("calls onConfirm when the user closes the modal", async () => {
    const user = userEvent.setup()
    const onConfirm = vi.fn()

    render(
      <ApiKeyCredentialsModal
        isOpen
        credentials={sampleCredentials}
        onConfirm={onConfirm}
      />,
    )

    await user.click(screen.getByRole("button", { name: /i've saved my credentials/i }))
    expect(onConfirm).toHaveBeenCalledTimes(1)
  })
})
