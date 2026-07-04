import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"
import { MemoryRouter } from "react-router-dom"

import { GetStartedDialog } from "./get-started-dialog"
import { useGetStartedContext } from "../../hooks/use-get-started-context.hook"

vi.mock("../../hooks/use-get-started-context.hook", () => ({
  useGetStartedContext: vi.fn(),
}))

const mockedUseGetStartedContext = vi.mocked(useGetStartedContext)

describe("GetStartedDialog", () => {
  it("renders workspace credentials and curl examples", () => {
    mockedUseGetStartedContext.mockReturnValue({
      context: {
        workspaceId: "00000000-0000-0000-0000-000000000001",
        apiKeys: [{ client_id: "demo-client-id" }],
      },
      isLoading: false,
      error: null,
    })

    render(
      <MemoryRouter>
        <GetStartedDialog open onOpenChange={vi.fn()} workspaceId="00000000-0000-0000-0000-000000000001" />
      </MemoryRouter>,
    )

    expect(screen.getByText("Your credentials")).toBeInTheDocument()
    expect(screen.getByText("demo-client-id")).toBeInTheDocument()
    expect(screen.getAllByText(/X-Workspace-Key: 00000000-0000-0000-0000-000000000001/).length).toBeGreaterThan(0)
    expect(screen.getByText(/\/v1\/email/)).toBeInTheDocument()
  })

  it("shows loading state", () => {
    mockedUseGetStartedContext.mockReturnValue({
      context: null,
      isLoading: true,
      error: null,
    })

    render(
      <MemoryRouter>
        <GetStartedDialog open onOpenChange={vi.fn()} workspaceId="ws-1" />
      </MemoryRouter>,
    )

    expect(screen.getByText(/loading workspace credentials/i)).toBeInTheDocument()
  })
})
