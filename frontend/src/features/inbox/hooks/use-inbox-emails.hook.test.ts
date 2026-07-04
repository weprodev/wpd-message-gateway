import { act, renderHook, waitFor } from "@testing-library/react"
import { describe, expect, it, vi, beforeEach } from "vitest"

import { fetchInboxEmailById, fetchInboxEmails } from "../inbox.api"
import { useInboxEmails } from "./use-inbox-emails.hook"

vi.mock("../inbox.api", () => ({
  fetchInboxEmails: vi.fn(),
  fetchInboxEmailById: vi.fn(),
}))

const mockedFetchInboxEmails = vi.mocked(fetchInboxEmails)
const mockedFetchInboxEmailById = vi.mocked(fetchInboxEmailById)

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

describe("useInboxEmails", () => {
  beforeEach(() => {
    mockedFetchInboxEmailById.mockResolvedValue({ ok: false })
  })

  it("loads emails for a workspace", async () => {
    mockedFetchInboxEmails.mockResolvedValue({
      ok: true,
      page: { items: [sampleEmail], has_more: false, next_cursor: undefined },
    })

    const { result } = renderHook(() => useInboxEmails("ws-1"))

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(mockedFetchInboxEmails).toHaveBeenCalledWith("ws-1", { limit: 50 })
    expect(result.current.messages).toHaveLength(1)
    expect(result.current.selectedMessageId).toBe("email-1")
    expect(result.current.error).toBeNull()
  })

  it("surfaces API errors", async () => {
    mockedFetchInboxEmails.mockResolvedValue({
      ok: false,
      message: "inbox unavailable",
    })

    const { result } = renderHook(() => useInboxEmails("ws-1"))

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.error).toBe("inbox unavailable")
    expect(result.current.messages).toHaveLength(0)
  })

  it("appends more emails when loadMore is called", async () => {
    mockedFetchInboxEmails
      .mockResolvedValueOnce({
        ok: true,
        page: {
          items: [sampleEmail],
          has_more: true,
          next_cursor: "cursor-1",
        },
      })
      .mockResolvedValueOnce({
        ok: true,
        page: {
          items: [{
            ...sampleEmail,
            id: "email-2",
            email: { ...sampleEmail.email, subject: "Second" },
          }],
          has_more: false,
          next_cursor: undefined,
        },
      })

    const { result } = renderHook(() => useInboxEmails("ws-1"))

    await waitFor(() => {
      expect(result.current.hasMore).toBe(true)
    })

    await act(async () => {
      await result.current.loadMore()
    })

    expect(result.current.messages).toHaveLength(2)
    expect(result.current.hasMore).toBe(false)
  })
})
