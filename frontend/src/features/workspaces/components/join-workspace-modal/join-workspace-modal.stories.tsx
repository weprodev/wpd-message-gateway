import type { Meta, StoryObj } from "@storybook/react-vite"
import { JoinWorkspaceModal } from "./join-workspace-modal"

const meta = {
  title: "Workspaces/JoinWorkspaceModal",
  component: JoinWorkspaceModal,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof JoinWorkspaceModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    isOpen: true,
    onClose: () => {},
    onSuccess: () => {},
  },
}
