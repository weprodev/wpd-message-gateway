import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { SendTestModal } from "./send-test-modal"

vi.mock("../../inbox.api", () => ({
  sendTestRequest: vi.fn(),
}))

describe("SendTestModal Component", () => {
  it("renders correctly when open", () => {
    render(
      <SendTestModal
        workspaceId="w1"
        open={true}
        onOpenChange={vi.fn()}
        onSent={vi.fn()}
      />
    )
    expect(screen.getByText("Send Test Request")).toBeInTheDocument()
    expect(screen.getByLabelText("To (comma-separated)")).toBeInTheDocument()
  })
})
