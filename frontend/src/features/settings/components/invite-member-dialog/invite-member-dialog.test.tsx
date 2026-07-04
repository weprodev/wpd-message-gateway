import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { Role } from "@/core/auth"

import type { WorkspaceMember } from "../../team.types"
import { InviteMemberDialog } from "./invite-member-dialog"

const existingMember: WorkspaceMember = {
  workspace_id: "ws-1",
  user_id: "user-member",
  role: Role.Member,
  joined_at: "2026-01-01T00:00:00Z",
  user_email: "member@weprodev.com",
}

describe("InviteMemberDialog", () => {
  it("submits email and selected role", async () => {
    const user = userEvent.setup()
    const onInvite = vi.fn().mockResolvedValue({
      id: "inv-1",
      email: "member@example.com",
      role: Role.Viewer,
      expires_at: "2026-01-08T00:00:00Z",
      token: "token",
    })

    render(<InviteMemberDialog open onClose={vi.fn()} onInvite={onInvite} />)

    await user.type(screen.getByLabelText(/email address/i), "member@example.com")
    await user.selectOptions(screen.getByLabelText(/role/i), Role.Viewer)
    await user.click(screen.getByRole("button", { name: /send invitation/i }))

    expect(onInvite).toHaveBeenCalledWith("member@example.com", Role.Viewer)
  })

  it("shows validation error for invalid email", async () => {
    const user = userEvent.setup()

    render(<InviteMemberDialog open onClose={vi.fn()} onInvite={vi.fn()} />)

    await user.type(screen.getByLabelText(/email address/i), "not-an-email")
    await user.click(screen.getByRole("button", { name: /send invitation/i }))

    expect(screen.getByText(/valid email address/i)).toBeInTheDocument()
  })

  it("rejects invitations for existing members", async () => {
    const user = userEvent.setup()
    const onInvite = vi.fn()

    render(
      <InviteMemberDialog
        open
        onClose={vi.fn()}
        onInvite={onInvite}
        members={[existingMember]}
      />,
    )

    await user.type(screen.getByLabelText(/email address/i), "member@weprodev.com")
    await user.click(screen.getByRole("button", { name: /send invitation/i }))

    expect(screen.getByText(/already a member/i)).toBeInTheDocument()
    expect(onInvite).not.toHaveBeenCalled()
  })
})
