import type { Meta, StoryObj } from "@storybook/react-vite"

import { Role } from "@/core/auth"

import { InvitationTokenDialog } from "./invitation-token-dialog"

const meta = {
  title: "Features/Settings/InvitationTokenDialog",
  component: InvitationTokenDialog,
  parameters: { layout: "centered" },
} satisfies Meta<typeof InvitationTokenDialog>

export default meta
type Story = StoryObj<typeof meta>

export const Open: Story = {
  args: {
    open: true,
    email: "invitee@example.com",
    role: Role.Member,
    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.demo-invitation-token",
    onClose: () => undefined,
  },
}

export const ViewerRole: Story = {
  args: {
    open: true,
    email: "viewer@weprodev.com",
    role: Role.Viewer,
    token: "one-time-viewer-invite-token",
    onClose: () => undefined,
  },
}
