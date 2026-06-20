import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"

import { Button } from "@/components/ui/button"
import { ConnectModal } from "./connect-modal"

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
  {
    id: "f2",
    provider_id: "mailgun",
    key: "domain",
    label: "Domain",
    description: "Your Mailgun Domain",
    field_type: "text",
    required: true,
    default_value: "",
    sort_order: 2,
  },
]

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

const meta = {
  title: "Features/Integrations/ConnectModal",
  component: ConnectModal,
  tags: ["autodocs"],
  parameters: { layout: "centered" },
  decorators: [
    (Story) => {
      const originalFetch = window.fetch
      window.fetch = async (input, init) => {
        const urlStr =
          typeof input === "string"
            ? input
            : input instanceof URL
              ? input.toString()
              : (input as Request).url
        if (urlStr.includes("/providers/mailgun/config")) {
          return {
            ok: true,
            status: 200,
            json: async () => mockFields,
          } as Response
        }
        return originalFetch(input, init)
      }
      return <Story />
    },
  ],
} satisfies Meta<typeof ConnectModal>

export default meta

type Story = StoryObj<typeof meta>

function ConnectModalDemo() {
  const [open, setOpen] = useState(true)
  return (
    <>
      <Button type="button" onClick={() => setOpen(true)}>
        Open Connect Modal
      </Button>
      <ConnectModal
        isOpen={open}
        onClose={() => setOpen(false)}
        workspaceId="demo-wid"
        provider={mailgunProvider}
        onConnect={async () => {
          await new Promise((resolve) => setTimeout(resolve, 500))
          return { ok: true }
        }}
      />
    </>
  )
}

export const MailgunConnection: Story = {
  render: () => <ConnectModalDemo />,
  args: {
    isOpen: true,
    onClose: () => undefined,
    workspaceId: "demo-wid",
    provider: mailgunProvider,
    onConnect: async () => ({ ok: true }),
  },
}
