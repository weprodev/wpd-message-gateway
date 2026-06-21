import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { Spinner } from "./spinner"

describe("Spinner Component", () => {
  it("renders with correct accessibility roles", () => {
    render(<Spinner />)
    expect(screen.getByRole("status")).toBeInTheDocument()
    expect(screen.getByLabelText("loading")).toBeInTheDocument()
  })

  it("applies the sizing classes correctly", () => {
    const { container: containerSm } = render(<Spinner size="sm" />)
    expect(containerSm.firstChild).toHaveClass("size-4")

    const { container: containerMd } = render(<Spinner size="md" />)
    expect(containerMd.firstChild).toHaveClass("size-6")

    const { container: containerLg } = render(<Spinner size="lg" />)
    expect(containerLg.firstChild).toHaveClass("size-8")
  })

  it("uses a contrasting border on solid backgrounds", () => {
    const { container } = render(<Spinner variant="onSolid" />)
    expect(container.firstChild).toHaveClass("border-current")
  })
})
