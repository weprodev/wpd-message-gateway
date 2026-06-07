import type { Meta, StoryObj } from "@storybook/react-vite"
import { SuccessModal } from "./success-modal"

const meta = {
  title: "Workspaces/SuccessModal",
  component: SuccessModal,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof SuccessModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    isOpen: true,
    workspaceName: "Design Team",
    onContinue: () => {},
  },
}
