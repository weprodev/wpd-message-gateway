import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { EmptyState } from "./empty-state"

describe("EmptyState Component", () => {
  it("renders correctly with message content", () => {
    render(<EmptyState />)
    expect(screen.getByText("No Workspaces Yet")).toBeInTheDocument()
    expect(screen.getByText("Create your first workspace to get started with Message Gateway.")).toBeInTheDocument()
  })
})
