import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { Permission, Role } from "./permissions"
import { Can } from "./can"
import { WorkspaceAuthorizationProvider } from "./workspace-authorization-context"

function renderWithAuth(
  ui: React.ReactNode,
  scope: { role?: string; permissions?: string[] },
) {
  return render(
    <WorkspaceAuthorizationProvider {...scope}>{ui}</WorkspaceAuthorizationProvider>,
  )
}

describe("Can", () => {
  it("renders children when permission is granted", () => {
    renderWithAuth(
      <Can permission={Permission.APIKeysWrite}>
        <button type="button">Generate key</button>
      </Can>,
      { role: Role.Member, permissions: [Permission.APIKeysWrite] },
    )

    expect(screen.getByRole("button", { name: "Generate key" })).toBeInTheDocument()
  })

  it("renders fallback when permission is missing", () => {
    renderWithAuth(
      <Can permission={Permission.MembersWrite} fallback={<span>Read only</span>}>
        <button type="button">Remove member</button>
      </Can>,
      { role: Role.Viewer, permissions: [Permission.MembersRead] },
    )

    expect(screen.queryByRole("button", { name: "Remove member" })).not.toBeInTheDocument()
    expect(screen.getByText("Read only")).toBeInTheDocument()
  })

  it("renders children when role matches", () => {
    renderWithAuth(
      <Can role={Role.Admin}>
        <p>Admin panel</p>
      </Can>,
      { role: Role.Admin, permissions: [] },
    )

    expect(screen.getByText("Admin panel")).toBeInTheDocument()
  })
})
