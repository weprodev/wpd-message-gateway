import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

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

describe("ApiKeyRow", () => {
  it("renders key details and action buttons", () => {
    render(<ApiKeyRow apiKey={sampleKey} onRegenerate={vi.fn()} onDelete={vi.fn()} />)

    expect(screen.getByText("Production")).toBeInTheDocument()
    expect(screen.getByText(/Last used: Never/i)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /regenerate/i })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /delete/i })).toBeInTheDocument()
  })

  it("calls regenerate and delete handlers", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn()
    const onDelete = vi.fn()

    render(<ApiKeyRow apiKey={sampleKey} onRegenerate={onRegenerate} onDelete={onDelete} />)

    await user.click(screen.getByRole("button", { name: /regenerate/i }))
    await user.click(screen.getByRole("button", { name: /delete/i }))

    expect(onRegenerate).toHaveBeenCalledWith("key-1")
    expect(onDelete).toHaveBeenCalledWith("key-1")
  })
})
