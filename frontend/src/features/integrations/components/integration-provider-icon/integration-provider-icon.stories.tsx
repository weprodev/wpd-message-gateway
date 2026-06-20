import type { Meta, StoryObj } from "@storybook/react-vite"

import { IntegrationProviderIcon } from "./integration-provider-icon"

const meta = {
  title: "Features/Integrations/IntegrationProviderIcon",
  component: IntegrationProviderIcon,
  tags: ["autodocs"],
  parameters: { layout: "centered" },
} satisfies Meta<typeof IntegrationProviderIcon>

export default meta

type Story = StoryObj<typeof meta>

export const EmojiIcon: Story = {
  args: {
    icon: "📧",
    name: "Mailgun",
    className: "size-12 p-2.5 text-2xl",
  },
}

export const ImageIcon: Story = {
  args: {
    icon: "/providers/mailgun.svg",
    name: "Mailgun",
    className: "size-12 p-2.5",
  },
}
