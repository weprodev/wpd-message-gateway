import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

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

describe("EmailContent Component", () => {
  it("renders correctly with metadata headers", () => {
    render(<EmailContent message={mockEmail} onDelete={vi.fn()} />)
    expect(screen.getByText("Welcome to Message Gateway!")).toBeInTheDocument()
    expect(screen.getAllByText(/John Doe/).length).toBeGreaterThan(0)
    expect(screen.getByText("e1")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Delete" })).toBeInTheDocument()
  })
})
