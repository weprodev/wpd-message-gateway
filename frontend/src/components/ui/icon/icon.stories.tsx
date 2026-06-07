import type { Meta, StoryObj } from "@storybook/react-vite"
import { Icon } from "./icon"

const meta = {
  title: "Components/Icon",
  component: Icon,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  argTypes: {
    name: { control: "text", description: "Material Symbol name (e.g. settings, mail)" },
    size: { control: "select", options: ["sm", "md", "lg"] },
  },
} satisfies Meta<typeof Icon>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    name: "settings",
    size: "md",
  },
}

export const Small: Story = {
  args: {
    name: "mail",
    size: "sm",
  },
}

export const Large: Story = {
  args: {
    name: "favorite",
    size: "lg",
  },
}
