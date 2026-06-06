import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"

import { Button } from "@/components/ui/button"
import { SendTestModal } from "./send-test-modal"

const meta = {
  title: "Features/Inbox/SendTestModal",
  component: SendTestModal,
  tags: ["autodocs"],
  parameters: { layout: "centered" },
} satisfies Meta<typeof SendTestModal>

export default meta

type Story = StoryObj<typeof meta>

function SendTestModalDemo() {
  const [open, setOpen] = useState(true)
  return (
    <>
      <Button type="button" onClick={() => setOpen(true)}>
        Open send test
      </Button>
      <SendTestModal
        workspaceId="00000000-0000-0000-0000-000000000001"
        open={open}
        onOpenChange={setOpen}
        onSent={() => undefined}
        initialChannel="email"
      />
    </>
  )
}

export const EmailChannel: Story = {
  render: () => <SendTestModalDemo />,
  args: {
    workspaceId: "00000000-0000-0000-0000-000000000001",
    open: true,
    onOpenChange: () => undefined,
    onSent: () => undefined,
    initialChannel: "email",
  },
}
