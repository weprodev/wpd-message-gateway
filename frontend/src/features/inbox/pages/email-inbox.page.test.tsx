import { render, screen } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { Permission, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import { useInboxEmails } from "../hooks/use-inbox-emails.hook"
import { EmailInboxPage } from "./email-inbox.page"

vi.mock("../hooks/use-inbox-emails.hook", () => ({
  useInboxEmails: vi.fn(),
}))

const mockedUseInboxEmails = vi.mocked(useInboxEmails)

const sampleEmail = {
  id: "email-1",
  workspace_id: "ws-1",
  channel: "email",
  status: "delivered",
  created_at: "2026-01-01T00:00:00Z",
  email: {
    to: ["a@b.com"],
    subject: "Hello",
    html: "<p>hi</p>",
  },
}

function renderPage(permissions: string[] = [Permission.LogsRead, Permission.InboxWrite]) {
  return render(
    <MemoryRouter initialEntries={["/workspaces/ws-1/inbox/email"]}>
      <WorkspaceAuthorizationProvider role={Role.Viewer} permissions={permissions}>
        <Routes>
          <Route path="/workspaces/:wid/inbox/email" element={<EmailInboxPage />} />
        </Routes>
      </WorkspaceAuthorizationProvider>
    </MemoryRouter>,
  )
}

describe("EmailInboxPage", () => {
  it("shows loading state while emails load", () => {
    mockedUseInboxEmails.mockReturnValue({
      messages: [],
      selectedMessageId: null,
      setSelectedMessageId: vi.fn(),
      isLoading: true,
      isLoadingMore: false,
      error: null,
      nextCursor: undefined,
      hasMore: false,
      reload: vi.fn(),
      loadMore: vi.fn(),
      upsertMessage: vi.fn(),
      removeMessage: vi.fn(),
      clearMessages: vi.fn(),
    })

    renderPage()

    expect(screen.getByText(/loading simulated inbox/i)).toBeInTheDocument()
  })

  it("renders email list when messages are available", () => {
    mockedUseInboxEmails.mockReturnValue({
      messages: [sampleEmail],
      selectedMessageId: "email-1",
      setSelectedMessageId: vi.fn(),
      isLoading: false,
      isLoadingMore: false,
      error: null,
      nextCursor: undefined,
      hasMore: false,
      reload: vi.fn(),
      loadMore: vi.fn(),
      upsertMessage: vi.fn(),
      removeMessage: vi.fn(),
      clearMessages: vi.fn(),
    })

    renderPage()

    expect(screen.getAllByText("Hello").length).toBeGreaterThan(0)
  })

  it("hides delete action for read-only users", () => {
    mockedUseInboxEmails.mockReturnValue({
      messages: [sampleEmail],
      selectedMessageId: "email-1",
      setSelectedMessageId: vi.fn(),
      isLoading: false,
      isLoadingMore: false,
      error: null,
      nextCursor: undefined,
      hasMore: false,
      reload: vi.fn(),
      loadMore: vi.fn(),
      upsertMessage: vi.fn(),
      removeMessage: vi.fn(),
      clearMessages: vi.fn(),
    })

    renderPage([Permission.LogsRead])

    expect(screen.queryByRole("button", { name: "Delete" })).not.toBeInTheDocument()
  })

  it("shows error message when hook reports failure", () => {
    mockedUseInboxEmails.mockReturnValue({
      messages: [],
      selectedMessageId: null,
      setSelectedMessageId: vi.fn(),
      isLoading: false,
      isLoadingMore: false,
      error: "Failed to load emails",
      nextCursor: undefined,
      hasMore: false,
      reload: vi.fn(),
      loadMore: vi.fn(),
      upsertMessage: vi.fn(),
      removeMessage: vi.fn(),
      clearMessages: vi.fn(),
    })

    renderPage()

    expect(screen.getByText("Failed to load emails")).toBeInTheDocument()
  })
})
