import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"
import { PinInput } from "./pin-input"

const meta = {
  title: "Components/PinInput",
  component: PinInput,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  argTypes: {
    length: { control: "number" },
  },
} satisfies Meta<typeof PinInput>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: (args) => {
    const [value, setValue] = useState("")
    return (
      <div className="w-[300px]">
        <PinInput {...args} value={value} onChange={setValue} />
        <p className="text-xs text-center text-text-secondary mt-4">Entered: {value}</p>
      </div>
    )
  },
  args: {
    length: 6,
    value: "",
    onChange: () => {},
  },
}
