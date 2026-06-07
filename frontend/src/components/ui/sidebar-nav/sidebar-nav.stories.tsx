import type { Meta, StoryObj } from "@storybook/react-vite"
import { MemoryRouter } from "react-router-dom"

import { SidebarNav, type SidebarNavSection } from "./sidebar-nav"

const mockSections: SidebarNavSection[] = [
  {
    label: "Navigation",
    items: [
      { label: "Inbox", segment: "email", icon: "mail" },
      { label: "SMS", segment: "sms", icon: "sms" },
      { label: "Push", segment: "push", icon: "notifications" },
    ],
  },
  {
    items: [
      { label: "Integrations", segment: "integrations", icon: "extension" },
      { label: "Settings", segment: "settings", icon: "settings" },
    ],
  },
]

const meta = {
  title: "Components/SidebarNav",
  component: SidebarNav,
  tags: ["autodocs"],
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={["/workspaces/ws-123/email"]}>
        <div className="w-64">
          <Story />
        </div>
      </MemoryRouter>
    ),
  ],
  parameters: {
    docs: {
      description: {
        component: "Navigation panel displaying application features/channels.",
      },
    },
  },
} satisfies Meta<typeof SidebarNav>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    sections: mockSections,
    workspaceId: "ws-123",
    buildHref: (wid, seg) => `/workspaces/${wid}/${seg}`,
  },
}
