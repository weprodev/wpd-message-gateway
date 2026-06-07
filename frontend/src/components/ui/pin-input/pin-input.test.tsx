import { render, fireEvent } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { PinInput } from "./pin-input"

describe("PinInput Component", () => {
  it("renders the correct number of input boxes", () => {
    const { container } = render(<PinInput value="" onChange={vi.fn()} length={6} />)
    const inputs = container.querySelectorAll("input")
    expect(inputs).toHaveLength(6)
  })

  it("calls onChange when a digit is entered", () => {
    const handleChange = vi.fn()
    const { container } = render(<PinInput value="" onChange={handleChange} length={6} />)
    const inputs = container.querySelectorAll("input")
    
    const input = inputs[0] as HTMLInputElement
    fireEvent.change(input, { target: { value: "1" } })

    expect(handleChange).toHaveBeenCalledWith("1")
  })
})
