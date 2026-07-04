import type { Meta, StoryObj } from "@storybook/react-vite"

import { ThemeProvider } from "@/shared/context/theme-context"
import { ThemeToggle } from "./theme-toggle"

const meta = {
  title: "Shared/ThemeToggle",
  component: ThemeToggle,
  parameters: { layout: "centered" },
  decorators: [
    (Story) => (
      <ThemeProvider>
        <Story />
      </ThemeProvider>
    ),
  ],
} satisfies Meta<typeof ThemeToggle>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {},
}
