import type { Meta, StoryObj } from "@storybook/react-vite"
import { EmailContent } from "./email-content"

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
    subject: "Welcome to Message Gateway!",
    html: "<h1>Welcome!</h1><p>Let's get started with your new communications setup.</p>",
  },
}

const meta = {
  title: "Features/Inbox/EmailContent",
  component: EmailContent,
  tags: ["autodocs"],
} satisfies Meta<typeof EmailContent>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    message: mockEmail,
    onDelete: () => {},
    isDeleting: false,
  },
}
