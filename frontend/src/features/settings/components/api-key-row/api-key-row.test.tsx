import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { ApiKeyRow } from "./api-key-row"

const sampleKey = {
  id: "key-1",
  workspace_id: "ws-1",
  client_id: "demo-client-id-abcdefghijklmnop",
  name: "Production",
  is_active: true,
  created_at: "2026-01-01T00:00:00Z",
  last_used_at: null,
}

function renderRow(
  props: React.ComponentProps<typeof ApiKeyRow>,
  permissions: string[] = [Permission.APIKeysWrite],
) {
  return render(
    <WorkspaceAuthorizationProvider role={Role.Member} permissions={permissions}>
      <ApiKeyRow {...props} />
    </WorkspaceAuthorizationProvider>,
  )
}

describe("ApiKeyRow", () => {
  it("renders key details and action buttons", () => {
    renderRow({ apiKey: sampleKey, onRegenerate: vi.fn(), onDelete: vi.fn() })

    expect(screen.getByText("Production")).toBeInTheDocument()
    expect(screen.getByText(/Last used: Never/i)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /regenerate/i })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /delete/i })).toBeInTheDocument()
  })

  it("hides action buttons for read-only users", () => {
    renderRow(
      { apiKey: sampleKey, onRegenerate: vi.fn(), onDelete: vi.fn() },
      [Permission.APIKeysRead],
    )

    expect(screen.queryByRole("button", { name: /regenerate/i })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: /delete/i })).not.toBeInTheDocument()
  })

  it("calls regenerate and delete handlers", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn()
    const onDelete = vi.fn()

    renderRow({ apiKey: sampleKey, onRegenerate, onDelete })

    await user.click(screen.getByRole("button", { name: /regenerate/i }))
    await user.click(screen.getByRole("button", { name: /delete/i }))

    expect(onRegenerate).toHaveBeenCalledWith("key-1")
    expect(onDelete).toHaveBeenCalledWith("key-1")
  })
})
