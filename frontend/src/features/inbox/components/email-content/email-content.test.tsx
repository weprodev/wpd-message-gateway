import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { EmailContent } from "./email-content"

const mockEmail = {
  id: "e1",
  workspace_id: "w1",
  channel: "email",
  status: "delivered",
  created_at: "2026-06-06T12:00:00Z",
  email: {
    from: "sender@example.com",
    from_name: "John Doe",
    to: ["recipient@example.com"],
    subject: "Welcome to Message Gateway!",
    html: "<h1>Welcome!</h1><p>Let's get started with your new communications setup.</p>",
  },
}

function renderContent(
  permissions: string[] = [Permission.InboxWrite],
) {
  return render(
    <WorkspaceAuthorizationProvider role={Role.Member} permissions={permissions}>
      <EmailContent message={mockEmail} onDelete={vi.fn()} />
    </WorkspaceAuthorizationProvider>,
  )
}

describe("EmailContent Component", () => {
  it("renders correctly with metadata headers", () => {
    renderContent()
    expect(screen.getByText("Welcome to Message Gateway!")).toBeInTheDocument()
    expect(screen.getAllByText(/John Doe/).length).toBeGreaterThan(0)
    expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument()
  })

  it("hides delete action for read-only users", () => {
    renderContent([Permission.LogsRead])

    expect(screen.queryByRole("button", { name: "Delete" })).not.toBeInTheDocument()
  })
})
