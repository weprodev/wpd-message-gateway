import type { Meta, StoryObj } from "@storybook/react-vite"

import { Role } from "@/core/auth"

import { InvitationRow } from "./invitation-row"

const meta = {
  title: "Features/Settings/InvitationRow",
  component: InvitationRow,
  parameters: { layout: "padded" },
  decorators: [
    (Story) => (
      <div className="max-w-2xl rounded-lg border border-border bg-card">
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof InvitationRow>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    invitation: {
      id: "inv-1",
      workspace_id: "ws-1",
      email: "invitee@example.com",
      role: Role.Viewer,
      expires_at: "2026-01-08T12:00:00Z",
      status: "pending",
      created_at: "2026-01-01T00:00:00Z",
    },
  },
}
