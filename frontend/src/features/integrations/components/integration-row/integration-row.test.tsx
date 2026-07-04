import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import type { IntegrationViewModel } from "@/features/integrations/integrations.types"
import { IntegrationRow } from "./integration-row"

const mailgunProvider: IntegrationViewModel = {
  id: "mailgun",
  name: "Mailgun",
  description: "Reliable transactional email sending service.",
  icon: "📧",
  category: "email",
  isAvailable: true,
  isConnected: false,
  isDeactivated: false,
}

describe("IntegrationRow", () => {
  it("shows connect action for available providers", async () => {
    const user = userEvent.setup()
    const onConnect = vi.fn()

    render(
      <IntegrationRow
        provider={mailgunProvider}
        onConnect={onConnect}
        onActivate={vi.fn()}
        onDisconnect={vi.fn()}
      />,
    )

    await user.click(screen.getByRole("button", { name: /connect/i }))
    expect(onConnect).toHaveBeenCalledWith(mailgunProvider)
  })

  it("shows disconnect action for connected providers", async () => {
    const user = userEvent.setup()
    const onDisconnect = vi.fn()
    const connected = { ...mailgunProvider, isConnected: true }

    render(
      <IntegrationRow
        provider={connected}
        onConnect={vi.fn()}
        onActivate={vi.fn()}
        onDisconnect={onDisconnect}
      />,
    )

    expect(screen.getByText("Connected")).toBeInTheDocument()
    await user.click(screen.getByRole("button", { name: /disconnect/i }))
    expect(onDisconnect).toHaveBeenCalledWith(connected)
  })

  it("shows coming soon badge for unavailable providers", () => {
    render(
      <IntegrationRow
        provider={{ ...mailgunProvider, isAvailable: false, isComingSoon: true }}
        onConnect={vi.fn()}
        onActivate={vi.fn()}
        onDisconnect={vi.fn()}
      />,
    )

    expect(screen.getByText("Coming soon")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: /connect/i })).not.toBeInTheDocument()
  })
})
