import type { Meta, StoryObj } from "@storybook/react-vite"

import { Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { TeamMemberRow } from "./team-member-row"

const meta = {
  title: "Features/Settings/TeamMemberRow",
  component: TeamMemberRow,
  parameters: { layout: "padded" },
  decorators: [
    (Story) => (
      <WorkspaceAuthorizationProvider role={Role.Admin} permissions={["members.write"]}>
        <div className="max-w-2xl rounded-lg border border-border bg-card">
          <Story />
        </div>
      </WorkspaceAuthorizationProvider>
    ),
  ],
} satisfies Meta<typeof TeamMemberRow>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    member: {
      workspace_id: "ws-1",
      user_id: "user-1",
      role: Role.Member,
      joined_at: "2026-01-01T00:00:00Z",
      user_email: "member@example.com",
      display_name: "Member User",
    },
    isCurrentUser: false,
    onRemove: () => undefined,
  },
}

export const CurrentUser: Story = {
  args: {
    member: {
      workspace_id: "ws-1",
      user_id: "user-1",
      role: Role.Admin,
      joined_at: "2026-01-01T00:00:00Z",
      user_email: "admin@example.com",
      display_name: "Admin User",
    },
    isCurrentUser: true,
    onRemove: () => undefined,
  },
}
