import type { Meta, StoryObj } from "@storybook/react-vite"
import { EmailItem } from "./email-item"

const mockEmail = {
  id: "e1",
  workspace_id: "w1",
  channel: "email",
  status: "delivered",
  created_at: "2026-06-06T12:00:00Z",
  email: {
    from: "sender@example.com",
    from_name: "John Doe",
    to: ["recipient@example.com"],
    subject: "Test Subject",
    plain_text: "This is the body snippet of the test email.",
  },
}

const meta = {
  title: "Features/Inbox/EmailItem",
  component: EmailItem,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof EmailItem>

export default meta
type Story = StoryObj<typeof meta>

export const Unselected: Story = {
  args: {
    message: mockEmail,
    isSelected: false,
    onClick: () => {},
  },
}

export const Selected: Story = {
  args: {
    message: mockEmail,
    isSelected: true,
    onClick: () => {},
  },
}
