import type { Meta, StoryObj } from "@storybook/react-vite"
import { MemoryRouter } from "react-router-dom"

import { GetStartedDialog } from "./get-started-dialog"

const meta = {
  title: "Features/Inbox/GetStartedDialog",
  component: GetStartedDialog,
  decorators: [
    (Story) => (
      <MemoryRouter>
        <Story />
      </MemoryRouter>
    ),
  ],
} satisfies Meta<typeof GetStartedDialog>

export default meta
type Story = StoryObj<typeof meta>

export const Open: Story = {
  args: {
    open: true,
    onOpenChange: () => undefined,
    workspaceId: "ws-demo",
  },
}
