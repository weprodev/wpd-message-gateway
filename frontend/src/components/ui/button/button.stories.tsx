import type { Meta, StoryObj } from "@storybook/react-vite"
import { Icon } from "../icon"

import { Button } from "./button"

const meta = {
  title: "Components/Button",
  component: Button,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component: "Primary action control; supports `Icon` as leading, trailing, or icon-only content.",
      },
    },
  },
  argTypes: {
    variant: {
      control: "select",
      options: ["default", "destructive", "outline", "secondary", "ghost", "link"],
      description: "Maps to semantic Tailwind colors backed by CSS variables.",
    },
    size: {
      control: "select",
      options: ["default", "sm", "lg", "icon"],
    },
    disabled: { control: "boolean" },
    asChild: { table: { disable: true } },
    children: { control: "text" },
  },
} satisfies Meta<typeof Button>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: "Button",
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
    children: "Destructive",
  },
}

export const Outline: Story = {
  args: {
    variant: "outline",
    children: "Outline",
  },
}

export const Ghost: Story = {
  args: {
    variant: "ghost",
    children: "Ghost",
  },
}

export const Link: Story = {
  args: {
    variant: "link",
    children: "Link",
  },
}

export const Small: Story = {
  args: {
    size: "sm",
    children: "Small",
  },
}

export const Large: Story = {
  args: {
    size: "lg",
    children: "Large",
  },
}

export const IconLeft: Story = {
  name: "Icon (Left)",
  render: (args) => (
    <Button {...args}>
      <Icon name="mail" size="sm" /> 
      Send Email
    </Button>
  ),
  args: {
    variant: "default",
  },
}

export const IconRight: Story = {
  name: "Icon (Right)",
  render: (args) => (
    <Button {...args}>
      Process Payment
      <Icon name="arrow_forward" size="sm" /> 
    </Button>
  ),
  args: {
    variant: "default",
  },
}

export const IconOnly: Story = {
  name: "Icon (Only)",
  render: (args) => (
    <Button {...args} aria-label="Settings">
      <Icon name="settings" size="md" />
    </Button>
  ),
  args: {
    size: "icon",
    variant: "outline",
  },
}

export const Disabled: Story = {
  args: {
    disabled: true,
    children: "Disabled",
  },
}

export const AllVariants: Story = {
  name: "All variants (light)",
  render: () => (
    <div className="flex max-w-2xl flex-wrap items-center gap-3 rounded-lg border border-border bg-card p-6">
      <Button type="button">Default</Button>
      <Button type="button" variant="secondary">
        Secondary
      </Button>
      <Button type="button" variant="destructive">
        Destructive
      </Button>
      <Button type="button" variant="outline">
        Outline
      </Button>
      <Button type="button" variant="ghost">
        Ghost
      </Button>
      <Button type="button" variant="link">
        Link
      </Button>
    </div>
  ),
}

export const DarkMode: Story = {
  name: "All variants (dark)",
  parameters: {
    docs: {
      description: {
        story: "Same variants; the decorator adds a parent with the dark theme class so `.dark` tokens from index.css apply.",
      },
    },
  },
  decorators: [
    (Story) => (
      <div className="dark min-h-[280px] w-full min-w-[min(100%,420px)] rounded-xl border border-border bg-background p-8 text-foreground shadow-sm">
        <Story />
      </div>
    ),
  ],
  render: () => (
    <div className="flex max-w-2xl flex-wrap items-center gap-3">
      <Button type="button">Default</Button>
      <Button type="button" variant="secondary">
        Secondary
      </Button>
      <Button type="button" variant="destructive">
        Destructive
      </Button>
      <Button type="button" variant="outline">
        Outline
      </Button>
      <Button type="button" variant="ghost">
        Ghost
      </Button>
      <Button type="button" variant="link">
        Link
      </Button>
      <Button type="button" variant="default">
        <Icon name="mail" size="sm" /> 
        Send Email
      </Button>
      <Button type="button" variant="outline" size="icon">
        <Icon name="settings" size="md" /> 
      </Button>
    </div>
  ),
}

export const IconLeftDark: Story = {
  name: "Icon (Left) - Dark",
  decorators: [
    (Story) => (
      <div className="dark rounded-xl border border-border bg-background p-8 text-foreground shadow-sm">
        <Story />
      </div>
    ),
  ],
  render: (args) => (
    <Button {...args}>
      <Icon name="mail" size="sm" /> 
      Send Email
    </Button>
  ),
  args: {
    variant: "default",
  },
}

export const IconRightDark: Story = {
  name: "Icon (Right) - Dark",
  decorators: [
    (Story) => (
      <div className="dark rounded-xl border border-border bg-background p-8 text-foreground shadow-sm">
        <Story />
      </div>
    ),
  ],
  render: (args) => (
    <Button {...args}>
      Process Payment
      <Icon name="arrow_forward" size="sm" /> 
    </Button>
  ),
  args: {
    variant: "default",
  },
}

export const IconOnlyDark: Story = {
  name: "Icon (Only) - Dark",
  decorators: [
    (Story) => (
      <div className="dark rounded-xl border border-border bg-background p-8 text-foreground shadow-sm">
        <Story />
      </div>
    ),
  ],
  render: (args) => (
    <Button {...args} aria-label="Settings">
      <Icon name="settings" size="md" />
    </Button>
  ),
  args: {
    size: "icon",
    variant: "outline",
  },
}
