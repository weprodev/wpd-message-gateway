import type { Meta, StoryObj } from "@storybook/react-vite"

import { ApiKeySecretDialog } from "./api-key-secret-dialog"

const meta = {
  title: "Features/Settings/ApiKeySecretDialog",
  component: ApiKeySecretDialog,
} satisfies Meta<typeof ApiKeySecretDialog>

export default meta
type Story = StoryObj<typeof meta>

export const Open: Story = {
  args: {
    open: true,
    clientId: "demo-client-id",
    clientSecret: "ZVtUvv8aDOfM0R2krW3qKpn_EQ_FX-KXFa0lRZVfvpU",
    onClose: () => undefined,
  },
}
