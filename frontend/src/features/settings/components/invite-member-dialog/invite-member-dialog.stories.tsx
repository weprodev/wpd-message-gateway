import type { Meta, StoryObj } from "@storybook/react-vite"

import { InviteMemberDialog } from "./invite-member-dialog"

const meta = {
  title: "Features/Settings/InviteMemberDialog",
  component: InviteMemberDialog,
  parameters: { layout: "centered" },
} satisfies Meta<typeof InviteMemberDialog>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    open: true,
    onClose: () => undefined,
    onInvite: async (email, role) => ({
      id: "inv-1",
      email,
      role,
      expires_at: "2026-01-08T00:00:00Z",
      token: "sample-token",
    }),
  },
}
