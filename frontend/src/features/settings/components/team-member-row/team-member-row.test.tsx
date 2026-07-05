import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { TeamMemberRow } from "./team-member-row"

const sampleMember = {
  workspace_id: "ws-1",
  user_id: "user-1",
  role: Role.Admin,
  joined_at: "2026-01-01T00:00:00Z",
  user_email: "admin@example.com",
  display_name: "Admin User",
}

function renderRow(
  props: React.ComponentProps<typeof TeamMemberRow>,
  permissions: string[] = [Permission.MembersWrite],
) {
  return render(
    <WorkspaceAuthorizationProvider role={Role.Admin} permissions={permissions}>
      <TeamMemberRow {...props} />
    </WorkspaceAuthorizationProvider>,
  )
}

describe("TeamMemberRow", () => {
  it("renders member details and remove action", () => {
    renderRow({
      member: sampleMember,
      isCurrentUser: false,
      onRemove: vi.fn(),
    })

    expect(screen.getByText("Admin User")).toBeInTheDocument()
    expect(screen.getByText("Admin")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /remove/i })).toBeInTheDocument()
  })

  it("hides remove action for current user", () => {
    renderRow({
      member: sampleMember,
      isCurrentUser: true,
      onRemove: vi.fn(),
    })

    expect(screen.getByText("You")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: /remove/i })).not.toBeInTheDocument()
  })

  it("calls remove handler", async () => {
    const user = userEvent.setup()
    const onRemove = vi.fn()

    renderRow({
      member: sampleMember,
      isCurrentUser: false,
      onRemove,
    })

    await user.click(screen.getByRole("button", { name: /remove/i }))
    expect(onRemove).toHaveBeenCalledWith("user-1")
  })
})
