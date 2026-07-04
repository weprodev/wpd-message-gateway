import type { Meta, StoryObj } from "@storybook/react-vite"

import { MessageGatewayLogo } from "./message-gateway-logo"

const meta = {
  title: "Shared/MessageGatewayLogo",
  component: MessageGatewayLogo,
  parameters: { layout: "centered" },
} satisfies Meta<typeof MessageGatewayLogo>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {},
}
