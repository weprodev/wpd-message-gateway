import type { Meta, StoryObj } from "@storybook/react-vite"
import { MemoryRouter, Route, Routes } from "react-router-dom"

import { OverviewPage } from "./overview.page"

const meta = {
  title: "Features/Inbox/OverviewPage",
  component: OverviewPage,
  tags: ["autodocs"],
  parameters: { layout: "padded" },
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={["/workspaces/demo-wid/overview"]}>
        <Routes>
          <Route path="/workspaces/:wid/overview" element={<Story />} />
        </Routes>
      </MemoryRouter>
    ),
  ],
} satisfies Meta<typeof OverviewPage>

export default meta

type Story = StoryObj<typeof meta>

export const AllChannels: Story = {}

export const EmailChannel: Story = {
  args: { channel: "email" },
}
