import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { Badge } from "./badge"

describe("Badge Component", () => {
  it("renders the badge children correctly", () => {
    render(<Badge>Test Badge</Badge>)
    expect(screen.getByText("Test Badge")).toBeInTheDocument()
  })

  it("applies the variant classes correctly", () => {
    const { container } = render(<Badge variant="destructive">Destructive</Badge>)
    expect(container.firstChild).toHaveClass("bg-destructive/15")
  })
})
