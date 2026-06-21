import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"

import { Button } from "@/components/ui/button"
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

const meta = {
  title: "Features/Integrations/DisconnectModal",
  component: DisconnectModal,
  tags: ["autodocs"],
  parameters: { layout: "centered" },
} satisfies Meta<typeof DisconnectModal>

export default meta

type Story = StoryObj<typeof meta>

function DisconnectModalDemo() {
  const [open, setOpen] = useState(true)
  return (
    <>
      <Button type="button" onClick={() => setOpen(true)}>
        Open Disconnect Modal
      </Button>
      <DisconnectModal
        isOpen={open}
        provider={connectedMailgun}
        onClose={() => setOpen(false)}
        onDeactivate={async () => {
          await new Promise((resolve) => setTimeout(resolve, 500))
          return { ok: true }
        }}
        onRemove={async () => {
          await new Promise((resolve) => setTimeout(resolve, 500))
          return { ok: true }
        }}
      />
    </>
  )
}

export const MailgunDisconnect: Story = {
  render: () => <DisconnectModalDemo />,
  args: {
    isOpen: true,
    provider: connectedMailgun,
    onClose: () => undefined,
    onDeactivate: async () => ({ ok: true }),
    onRemove: async () => ({ ok: true }),
  },
}
