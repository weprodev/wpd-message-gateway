import { render, screen } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { Button } from "./button"

describe("Button Component", () => {
  it("renders with children correctly", () => {
    render(<Button>Click me</Button>)
    expect(screen.getByRole("button", { name: "Click me" })).toBeInTheDocument()
  })

  it("handles click events", async () => {
    const handleClick = vi.fn()
    render(<Button onClick={handleClick}>Click me</Button>)
    const button = screen.getByRole("button", { name: "Click me" })
    button.click()
    expect(handleClick).toHaveBeenCalledTimes(1)
  })

  it("applies variant and size classes", () => {
    const { container } = render(<Button variant="destructive" size="sm">Click me</Button>)
    const button = container.querySelector("button")
    expect(button).toHaveClass("bg-destructive")
    expect(button).toHaveClass("h-8")
  })
})
