import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"

import { Button } from "@/components/ui/button"
import { DisconnectModal } from "./disconnect-modal"

const mailgunProvider = {
  id: "mailgun",
  name: "Mailgun",
  description: "Reliable transactional email sending service.",
  icon: "📧",
  category: "email" as const,
  isAvailable: true,
  isConnected: true,
  isDeactivated: false,
  integration: {
    id: "int-mailgun-1",
    workspace_id: "demo-wid",
    channel_type: "email" as const,
    provider_name: "mailgun",
    status: "connected",
    is_default: true,
    config: { api_key: "key-123", domain: "mg.example.com" },
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
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
        onClose={() => setOpen(false)}
        provider={mailgunProvider}
        onDeactivate={async () => {
          await new Promise((resolve) => setTimeout(resolve, 1000))
          return { ok: true as const }
        }}
        onRemove={async () => {
          await new Promise((resolve) => setTimeout(resolve, 1000))
          return { ok: true as const }
        }}
      />
    </>
  )
}

export const MailgunDisconnect: Story = {
  render: () => <DisconnectModalDemo />,
  args: {
    isOpen: true,
    onClose: () => undefined,
    provider: mailgunProvider,
    onDeactivate: async () => ({ ok: true as const }),
    onRemove: async () => ({ ok: true as const }),
  },
}
