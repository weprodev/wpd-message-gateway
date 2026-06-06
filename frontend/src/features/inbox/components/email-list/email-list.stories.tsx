import type { Meta, StoryObj } from "@storybook/react-vite"
import { EmailList } from "./email-list"

const mockEmailList = [
  {
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
      plain_text: "Let's get started with your new communications setup.",
    },
  },
  {
    id: "e2",
    workspace_id: "w1",
    channel: "email",
    status: "delivered",
    created_at: "2026-06-06T11:30:00Z",
    email: {
      from: "alerts@example.com",
      from_name: "Alerts System",
      to: ["recipient@example.com"],
      subject: "Security Warning: New Login",
      plain_text: "We noticed a login to your account from a new IP address.",
    },
  },
]

const meta = {
  title: "Features/Inbox/EmailList",
  component: EmailList,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
} satisfies Meta<typeof EmailList>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    messages: mockEmailList,
    selectedMessageId: "e1",
    onSelectMessage: () => {},
  },
}
