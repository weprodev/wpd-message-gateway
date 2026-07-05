import type { Meta, StoryObj } from "@storybook/react-vite"

import { AppFooter } from "./app-footer"

const meta = {
  title: "Shared/AppFooter",
  component: AppFooter,
  parameters: { layout: "fullscreen" },
} satisfies Meta<typeof AppFooter>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {},
}

export const Dashboard: Story = {
  args: {
    variant: "dashboard",
  },
}
