import type { Meta, StoryObj } from "@storybook/react-vite"

import { PageHeader } from "./page-header"

const meta = {
  title: "Shared/PageHeader",
  component: PageHeader,
  parameters: { layout: "padded" },
} satisfies Meta<typeof PageHeader>

export default meta
type Story = StoryObj<typeof meta>

export const WithDescription: Story = {
  args: {
    title: "Settings",
    description: "Manage your workspace settings and preferences.",
  },
}

export const TitleOnly: Story = {
  args: {
    title: "Overview",
  },
}
