import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { INTEGRATION_STATUS } from "@/features/integrations/integrations.types"
import { DisconnectModal } from "./disconnect-modal"

const connectedMailgun = {
  id: "mailgun",
  name: "Mailgun",
  description: "Reliable transactional email sending service.",
  icon: "📧",
  category: "email" as const,
  isAvailable: true,
  isConnected: true,
  isDeactivated: false,
  integration: {
    id: "intg-1",
    workspace_id: "demo-wid",
    channel_type: "email" as const,
    provider_name: "mailgun",
    config: { api_key: "hidden" },
    status: INTEGRATION_STATUS.CONNECTED,
    is_default: true,
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
  },
}

describe("DisconnectModal", () => {
  it("renders disconnect actions for connected provider", () => {
    render(
      <DisconnectModal
        isOpen
        provider={connectedMailgun}
        onClose={vi.fn()}
        onDeactivate={vi.fn().mockResolvedValue({ ok: true })}
        onRemove={vi.fn().mockResolvedValue({ ok: true })}
      />,
    )

    expect(screen.getByText(/disconnect mailgun/i)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /deactivate/i })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /remove integration/i })).toBeInTheDocument()
  })

  it("calls deactivate handler", async () => {
    const user = userEvent.setup()
    const onDeactivate = vi.fn().mockResolvedValue({ ok: true })
    const onClose = vi.fn()

    render(
      <DisconnectModal
        isOpen
        provider={connectedMailgun}
        onClose={onClose}
        onDeactivate={onDeactivate}
        onRemove={vi.fn().mockResolvedValue({ ok: true })}
      />,
    )

    await user.click(screen.getByRole("button", { name: /^deactivate$/i }))

    expect(onDeactivate).toHaveBeenCalledWith(connectedMailgun)
    expect(onClose).toHaveBeenCalledTimes(1)
  })
})
