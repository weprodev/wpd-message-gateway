import type { Meta, StoryObj } from "@storybook/react-vite"
import { PinEntryModal } from "./pin-entry-modal"

const meta = {
  title: "Workspaces/PinEntryModal",
  component: PinEntryModal,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof PinEntryModal>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    isOpen: true,
    workspaceName: "Engineering Team",
    onClose: () => {},
    onSubmit: () => {},
    error: "",
    isLoading: false,
  },
}

export const Loading: Story = {
  args: {
    isOpen: true,
    workspaceName: "Engineering Team",
    onClose: () => {},
    onSubmit: () => {},
    error: "",
    isLoading: true,
  },
}

export const WithError: Story = {
  args: {
    isOpen: true,
    workspaceName: "Engineering Team",
    onClose: () => {},
    onSubmit: () => {},
    error: "Incorrect 6-digit passcode. Please try again.",
    isLoading: false,
  },
}
