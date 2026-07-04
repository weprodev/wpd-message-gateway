import type { Meta, StoryObj } from "@storybook/react-vite"

import type { IntegrationViewModel } from "@/features/integrations/hooks/use-integrations.hook"
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

const connectedProvider: IntegrationViewModel = {
  ...mailgunProvider,
  isConnected: true,
}

const deactivatedProvider: IntegrationViewModel = {
  ...mailgunProvider,
  isDeactivated: true,
}

const comingSoonProvider: IntegrationViewModel = {
  ...mailgunProvider,
  id: "twilio",
  name: "Twilio",
  isAvailable: false,
  isComingSoon: true,
}

const meta = {
  title: "Features/Integrations/IntegrationRow",
  component: IntegrationRow,
  parameters: { layout: "padded" },
} satisfies Meta<typeof IntegrationRow>

export default meta
type Story = StoryObj<typeof meta>

export const Connect: Story = {
  args: {
    provider: mailgunProvider,
    onConnect: () => undefined,
    onActivate: () => undefined,
    onDisconnect: () => undefined,
  },
}

export const Connected: Story = {
  args: {
    provider: connectedProvider,
    onConnect: () => undefined,
    onActivate: () => undefined,
    onDisconnect: () => undefined,
  },
}

export const Deactivated: Story = {
  args: {
    provider: deactivatedProvider,
    onConnect: () => undefined,
    onActivate: () => undefined,
    onDisconnect: () => undefined,
  },
}

export const ComingSoon: Story = {
  args: {
    provider: comingSoonProvider,
    onConnect: () => undefined,
    onActivate: () => undefined,
    onDisconnect: () => undefined,
  },
}
