import type { Meta, StoryObj } from "@storybook/react-vite"

import { AuthCard } from "./auth-card"

const meta = {
  title: "Shared/AuthCard",
  component: AuthCard,
  parameters: { layout: "centered" },
} satisfies Meta<typeof AuthCard>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: (
      <>
        <h2 className="text-lg font-semibold text-foreground">Sign in</h2>
        <p className="mt-2 text-sm text-text-secondary">Use your portal credentials.</p>
      </>
    ),
  },
}
