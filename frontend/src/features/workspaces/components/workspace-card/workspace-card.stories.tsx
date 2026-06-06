import type { Meta, StoryObj } from "@storybook/react-vite"
import { WorkspaceCard } from "./workspace-card"

const mockWorkspace = {
  id: "w1",
  name: "Production Workspace",
  unique_key: "production_workspace",
  icon_key: "shield",
  visibility: "private" as const,
  status: "active",
  created_at: "2026-06-06T12:00:00Z",
  updated_at: "2026-06-06T12:00:00Z",
}

const meta = {
  title: "Workspaces/WorkspaceCard",
  component: WorkspaceCard,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof WorkspaceCard>

export default meta
type Story = StoryObj<typeof meta>

export const Unselected: Story = {
  args: {
    workspace: mockWorkspace,
    isSelected: false,
    onSelect: () => {},
  },
}

export const Selected: Story = {
  args: {
    workspace: mockWorkspace,
    isSelected: true,
    onSelect: () => {},
  },
}
