import type { Meta, StoryObj } from "@storybook/react-vite"
import { CreateWorkspaceModal } from "./create-workspace-modal"

const meta = {
  title: "Workspaces/CreateWorkspaceModal",
  component: CreateWorkspaceModal,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof CreateWorkspaceModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    isOpen: true,
    onClose: () => {},
    onSuccess: () => {},
  },
}
