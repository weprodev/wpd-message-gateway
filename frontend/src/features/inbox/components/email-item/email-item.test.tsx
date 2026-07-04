import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { EmailItem } from "./email-item"

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
    subject: "Test Subject",
    plain_text: "This is the body snippet of the test email.",
  },
}

describe("EmailItem Component", () => {
  it("renders correctly with email properties", () => {
    render(<EmailItem message={mockEmail} isSelected={false} onClick={vi.fn()} />)
    expect(screen.getByText("recipient@example.com")).toBeInTheDocument()
    expect(screen.getByText("Test Subject")).toBeInTheDocument()
    expect(screen.getByText("This is the body snippet of the test email.")).toBeInTheDocument()
  })

  it("triggers onClick when clicked", () => {
    const handleClick = vi.fn()
    render(<EmailItem message={mockEmail} isSelected={false} onClick={handleClick} />)
    screen.getByRole("button").click()
    expect(handleClick).toHaveBeenCalledTimes(1)
  })
})
