import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { PageHeader } from "./page-header"

describe("PageHeader", () => {
  it("renders title and optional description", () => {
    render(<PageHeader title="Settings" description="Manage workspace settings." />)

    expect(screen.getByRole("heading", { name: "Settings" })).toBeInTheDocument()
    expect(screen.getByText("Manage workspace settings.")).toBeInTheDocument()
  })
})
