import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"
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

function renderRow(
  props: React.ComponentProps<typeof IntegrationRow>,
  permissions: string[] = [Permission.IntegrationsWrite],
) {
  return render(
    <WorkspaceAuthorizationProvider role={Role.Member} permissions={permissions}>
      <IntegrationRow {...props} />
    </WorkspaceAuthorizationProvider>,
  )
}

describe("IntegrationRow", () => {
  it("shows connect action for available providers", async () => {
    const user = userEvent.setup()
    const onConnect = vi.fn()

    renderRow({
      provider: mailgunProvider,
      onConnect,
      onActivate: vi.fn(),
      onDisconnect: vi.fn(),
    })

    await user.click(screen.getByRole("button", { name: /connect/i }))
    expect(onConnect).toHaveBeenCalledWith(mailgunProvider)
  })

  it("shows disconnect action for connected providers", async () => {
    const user = userEvent.setup()
    const onDisconnect = vi.fn()
    const connected = { ...mailgunProvider, isConnected: true }

    renderRow({
      provider: connected,
      onConnect: vi.fn(),
      onActivate: vi.fn(),
      onDisconnect,
    })

    expect(screen.getByText("Connected")).toBeInTheDocument()
    await user.click(screen.getByRole("button", { name: /disconnect/i }))
    expect(onDisconnect).toHaveBeenCalledWith(connected)
  })

  it("hides management actions for read-only users", () => {
    renderRow(
      {
        provider: mailgunProvider,
        onConnect: vi.fn(),
        onActivate: vi.fn(),
        onDisconnect: vi.fn(),
      },
      [Permission.IntegrationsRead],
    )

    expect(screen.getByText("Not connected")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: /connect/i })).not.toBeInTheDocument()
  })

  it("shows coming soon badge for unavailable providers", () => {
    renderRow({
      provider: { ...mailgunProvider, isAvailable: false, isComingSoon: true },
      onConnect: vi.fn(),
      onActivate: vi.fn(),
      onDisconnect: vi.fn(),
    })

    expect(screen.getByText("Coming soon")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: /connect/i })).not.toBeInTheDocument()
  })
})
