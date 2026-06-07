import type { Meta, StoryObj } from "@storybook/react-vite"
import { WorkspaceActions } from "./workspace-actions"

const meta = {
  title: "Workspaces/WorkspaceActions",
  component: WorkspaceActions,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof WorkspaceActions>

export default meta
type Story = StoryObj<typeof meta>

export const CardVariant: Story = {
  args: {
    variant: "card",
  },
}

export const InlineVariant: Story = {
  args: {
    variant: "inline",
  },
}
