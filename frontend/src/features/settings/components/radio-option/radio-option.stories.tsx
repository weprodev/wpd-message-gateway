import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"

import { RadioOption } from "./radio-option"

const meta = {
  title: "Features/Settings/RadioOption",
  component: RadioOption,
  parameters: { layout: "padded" },
} satisfies Meta<typeof RadioOption>

export default meta
type Story = StoryObj<typeof meta>

function RadioOptionDemo() {
  const [checked, setChecked] = useState(true)
  return (
    <RadioOption
      id="dispatch-memory"
      name="dispatch-mode"
      label="Memory"
      description="Capture messages in memory for development and testing."
      checked={checked}
      onChange={() => setChecked((value) => !value)}
    />
  )
}

export const Selected: Story = {
  render: () => <RadioOptionDemo />,
  args: {
    id: "dispatch-memory",
    name: "dispatch-mode",
    label: "Memory",
    description: "Capture messages in memory for development and testing.",
    checked: true,
    onChange: () => undefined,
  },
}

export const Unselected: Story = {
  args: {
    id: "dispatch-provider",
    name: "dispatch-mode",
    label: "Provider",
    description: "Send messages through the connected channel integration.",
    checked: false,
    onChange: () => undefined,
  },
}
