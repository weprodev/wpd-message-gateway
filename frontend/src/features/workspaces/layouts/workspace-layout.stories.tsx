import type { Meta, StoryObj } from "@storybook/react-vite"
import { MemoryRouter, Route, Routes } from "react-router-dom"

import { WorkspaceLayout } from "./workspace-layout"

const meta = {
  title: "Features/Workspaces/WorkspaceLayout",
  component: WorkspaceLayout,
  tags: ["autodocs"],
  parameters: { layout: "fullscreen" },
  decorators: [
    (Story) => (
      <MemoryRouter initialEntries={["/workspaces/demo-wid/overview"]}>
        <Routes>
          <Route path="/workspaces/:wid" element={<Story />}>
            <Route path="overview" element={<p className="text-sm text-muted-foreground">Page content</p>} />
          </Route>
        </Routes>
      </MemoryRouter>
    ),
  ],
} satisfies Meta<typeof WorkspaceLayout>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
