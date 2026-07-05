import { render, screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ConnectModal } from "./connect-modal"

const mailgunProvider = {
  id: "mailgun",
  name: "Mailgun",
  description: "Reliable transactional email sending service.",
  icon: "📧",
  category: "email" as const,
  isAvailable: true,
  isConnected: false,
  isDeactivated: false,
}

const mockFields = [
  {
    id: "f1",
    provider_id: "mailgun",
    key: "api_key",
    label: "API Key",
    description: "Your Mailgun Private API Key",
    field_type: "password",
    required: true,
    default_value: "",
    sort_order: 1,
  },
]

describe("ConnectModal", () => {
  it("loads config fields and submits connection", async () => {
    const user = userEvent.setup()
    const onConnect = vi.fn().mockResolvedValue({ ok: true })
    const onClose = vi.fn()

    vi.spyOn(globalThis, "fetch").mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => mockFields,
    } as Response)

    render(
      <ConnectModal
        isOpen
        onClose={onClose}
        workspaceId="demo-wid"
        provider={mailgunProvider}
        onConnect={onConnect}
      />,
    )

    expect(screen.getByText(/loading configuration/i)).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByLabelText(/api key/i)).toBeInTheDocument()
    })

    await user.type(screen.getByLabelText(/api key/i), "secret-key")
    await user.click(screen.getByRole("button", { name: /^connect$/i }))

    await waitFor(() => {
      expect(onConnect).toHaveBeenCalledWith(mailgunProvider, { api_key: "secret-key" })
      expect(onClose).toHaveBeenCalledTimes(1)
    })
  })
})
