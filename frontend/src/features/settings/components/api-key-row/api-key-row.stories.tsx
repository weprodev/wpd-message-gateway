import type { Meta, StoryObj } from "@storybook/react-vite"

import { ApiKeyRow } from "./api-key-row"

const meta = {
  title: "Features/Settings/ApiKeyRow",
  component: ApiKeyRow,
  parameters: { layout: "padded" },
} satisfies Meta<typeof ApiKeyRow>

export default meta
type Story = StoryObj<typeof meta>

const sampleKey = {
  id: "key-1",
  workspace_id: "00000000-0000-0000-0000-000000000001",
  client_id: "demo-client-id-abcdefghijklmnop",
  name: "Production",
  is_active: true,
  created_at: "2026-01-01T00:00:00Z",
  last_used_at: null,
}

export const Default: Story = {
  args: {
    apiKey: sampleKey,
    onRegenerate: () => undefined,
    onDelete: () => undefined,
  },
}

export const Busy: Story = {
  args: {
    apiKey: sampleKey,
    isBusy: true,
    onRegenerate: () => undefined,
    onDelete: () => undefined,
  },
}
