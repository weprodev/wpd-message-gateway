import type { Meta, StoryObj } from "@storybook/react-vite"

import { ApiKeyCreateDialog } from "./api-key-create-dialog"

const meta = {
  title: "Features/Settings/ApiKeyCreateDialog",
  component: ApiKeyCreateDialog,
} satisfies Meta<typeof ApiKeyCreateDialog>

export default meta
type Story = StoryObj<typeof meta>

export const Open: Story = {
  args: {
    open: true,
    onClose: () => undefined,
    onCreate: async () => {
      await new Promise((resolve) => window.setTimeout(resolve, 500))
    },
  },
}
