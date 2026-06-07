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
  {
    id: "f3",
    provider_id: "mailgun",
    key: "base_url",
    label: "Base URL",
    description: "Mailgun API Base URL",
    field_type: "text",
    required: false,
    default_value: "https://api.mailgun.net/v3",
    sort_order: 3,
  },
  {
    id: "f4",
    provider_id: "mailgun",
    key: "from_email",
    label: "From Email",
    description: "Default sender email address",
    field_type: "email",
    required: true,
    default_value: "",
    sort_order: 4,
  },
  {
    id: "f5",
    provider_id: "mailgun",
    key: "from_name",
    label: "From Name",
    description: "Default sender name",
    field_type: "text",
    required: false,
    default_value: "",
    sort_order: 5,
  },
]

const meta = {
  title: "Features/Integrations/ConnectModal",
  component: ConnectModal,
  tags: ["autodocs"],
  parameters: { layout: "centered" },
  decorators: [
    (Story) => {
      // Mock window.fetch for loading the provider configuration fields
      const originalFetch = window.fetch
      window.fetch = async (input, init) => {
        const urlStr = typeof input === "string" ? input : input instanceof URL ? input.toString() : (input as Request).url
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
        provider={{
          id: "mailgun",
          name: "Mailgun",
          description: "Reliable transactional email sending service.",
          icon: "📧",
          category: "email",
          isAvailable: true,
          isConnected: false,
        }}
        onConnect={async (provider, config) => {
          console.log("Connect called with:", provider, config)
          await new Promise((resolve) => setTimeout(resolve, 1000))
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
    provider: {
      id: "mailgun",
      name: "Mailgun",
      description: "Reliable transactional email sending service.",
      icon: "📧",
      category: "email",
      isAvailable: true,
      isConnected: false,
    },
    onConnect: async () => undefined,
  },
}
