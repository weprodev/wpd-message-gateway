import { render, screen, fireEvent } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { Input } from "./input"

describe("Input Component", () => {
  it("renders correctly", () => {
    render(<Input placeholder="Test Input" />)
    expect(screen.getByPlaceholderText("Test Input")).toBeInTheDocument()
  })

  it("handles changes", () => {
    const handleChange = vi.fn()
    render(<Input placeholder="Test Input" onChange={handleChange} />)
    const input = screen.getByPlaceholderText("Test Input")
    fireEvent.change(input, { target: { value: "hello" } })
    expect(handleChange).toHaveBeenCalled()
  })

  it("can be disabled", () => {
    render(<Input placeholder="Test Input" disabled />)
    expect(screen.getByPlaceholderText("Test Input")).toBeDisabled()
  })
})
