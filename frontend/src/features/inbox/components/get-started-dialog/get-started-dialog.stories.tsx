import type { Meta, StoryObj } from "@storybook/react-vite"
import { MemoryRouter } from "react-router-dom"

import { GetStartedDialog } from "./get-started-dialog"

const meta = {
  title: "Features/Inbox/GetStartedDialog",
  component: GetStartedDialog,
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={["/workspaces/ws-demo/overview"]}>
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
    workspaceId: "00000000-0000-0000-0000-000000000001",
  },
}
