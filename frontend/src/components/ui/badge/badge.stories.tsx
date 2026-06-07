import type { Meta, StoryObj } from "@storybook/react-vite"

import { Badge } from "./badge"

const meta = {
  title: "Components/Badge",
  component: Badge,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component: "Status and label chips; colors come from CSS variables in `src/index.css`.",
      },
    },
  },
  argTypes: {
    variant: {
      control: "select",
      options: ["default", "secondary", "destructive", "success", "warning", "outline"],
      description: "Maps to semantic color tokens.",
    },
    children: { control: "text" },
  },
} satisfies Meta<typeof Badge>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: "Default",
  },
}

export const Secondary: Story = {
  args: {
    variant: "secondary",
    children: "Secondary",
  },
}

export const Destructive: Story = {
  args: {
    variant: "destructive",
    children: "400 Bad Request",
  },
}

export const Success: Story = {
  args: {
    variant: "success",
    children: "200 OK",
  },
}

export const Warning: Story = {
  args: {
    variant: "warning",
    children: "202 Accepted",
  },
}

export const Outline: Story = {
  args: {
    variant: "outline",
    children: "Outline",
  },
}

export const AllVariantsLight: Story = {
  name: "All variants (light)",
  render: () => (
    <div className="flex max-w-2xl flex-wrap items-center gap-2 rounded-lg border border-border bg-card p-6">
      <Badge>Default</Badge>
      <Badge variant="secondary">Secondary</Badge>
      <Badge variant="destructive">Destructive</Badge>
      <Badge variant="success">Success</Badge>
      <Badge variant="warning">Warning</Badge>
      <Badge variant="outline">Outline</Badge>
    </div>
  ),
}

export const AllVariantsDark: Story = {
  name: "All variants (dark)",
  parameters: {
    docs: {
      description: {
        story: "Rendered inside a `dark` container to preview `.dark` token values.",
      },
    },
  },
  decorators: [
    (Story) => (
      <div className="dark min-h-[200px] w-full min-w-[min(100%,420px)] rounded-xl border border-border bg-background p-8 text-foreground shadow-sm">
        <Story />
      </div>
    ),
  ],
  render: () => (
    <div className="flex max-w-2xl flex-wrap items-center gap-2">
      <Badge>Default</Badge>
      <Badge variant="secondary">Secondary</Badge>
      <Badge variant="destructive">Destructive</Badge>
      <Badge variant="success">Success</Badge>
      <Badge variant="warning">Warning</Badge>
      <Badge variant="outline">Outline</Badge>
    </div>
  ),
}
