import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { EmailList } from "./email-list"

const mockEmailList = [
  {
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
      plain_text: "Let's get started with your new communications setup.",
    },
  },
]

describe("EmailList Component", () => {
  it("renders correctly with search bar and item counts", () => {
    render(<EmailList messages={mockEmailList} selectedMessageId="e1" onSelectMessage={vi.fn()} />)
    expect(screen.getByText("Inbox")).toBeInTheDocument()
    expect(screen.getByText("1")).toBeInTheDocument()
    expect(screen.getByPlaceholderText("Search inbox...")).toBeInTheDocument()
    expect(screen.getByText("recipient@example.com")).toBeInTheDocument()
  })
})
