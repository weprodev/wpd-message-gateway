import { render, screen, waitFor } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { fetchEmailTemplates } from "../api/inbox.api"
import { EmailTemplatesPage } from "./email-templates.page"

vi.mock("../api/inbox.api", () => ({
  fetchEmailTemplates: vi.fn(),
  createEmailTemplate: vi.fn(),
  deleteEmailTemplate: vi.fn(),
}))

const mockedFetchTemplates = vi.mocked(fetchEmailTemplates)

function renderPage(permissions: string[] = [Permission.TemplatesRead, Permission.TemplatesWrite]) {
  return render(
    <MemoryRouter initialEntries={["/workspaces/ws-1/inbox/templates"]}>
      <WorkspaceAuthorizationProvider role={Role.Viewer} permissions={permissions}>
        <Routes>
          <Route path="/workspaces/:wid/inbox/templates" element={<EmailTemplatesPage />} />
        </Routes>
      </WorkspaceAuthorizationProvider>
    </MemoryRouter>,
  )
}

describe("EmailTemplatesPage", () => {
  it("shows loading state", () => {
    mockedFetchTemplates.mockReturnValue(new Promise(() => {}))

    renderPage()

    expect(screen.getByText(/loading assets/i)).toBeInTheDocument()
  })

  it("renders template table rows", async () => {
    mockedFetchTemplates.mockResolvedValue({
      ok: true,
      items: [
        {
          id: "tpl-1",
          workspace_id: "ws-1",
          name: "Welcome",
          unique_key: "welcome",
          channel_type: "email",
          category: "transactional",
          subject: "Welcome!",
          content_html: "<p>Hi</p>",
          is_active: true,
          is_default: false,
          created_at: "2026-01-01T00:00:00Z",
          updated_at: "2026-01-01T00:00:00Z",
        },
      ],
    })

    renderPage()

    await waitFor(() => {
      expect(screen.getByText("Welcome")).toBeInTheDocument()
    })
    expect(screen.getByText("welcome")).toBeInTheDocument()
  })

  it("hides template write actions for read-only users", async () => {
    mockedFetchTemplates.mockResolvedValue({
      ok: true,
      items: [
        {
          id: "tpl-1",
          workspace_id: "ws-1",
          name: "Welcome",
          unique_key: "welcome",
          channel_type: "email",
          category: "transactional",
          subject: "Welcome!",
          content_html: "<p>Hi</p>",
          is_active: true,
          is_default: false,
          created_at: "2026-01-01T00:00:00Z",
          updated_at: "2026-01-01T00:00:00Z",
        },
      ],
    })

    renderPage([Permission.TemplatesRead])

    await waitFor(() => {
      expect(screen.getByText("Welcome")).toBeInTheDocument()
    })
    expect(screen.queryByRole("button", { name: /create template/i })).not.toBeInTheDocument()
    expect(screen.queryByTitle("Delete Template")).not.toBeInTheDocument()
  })
})
