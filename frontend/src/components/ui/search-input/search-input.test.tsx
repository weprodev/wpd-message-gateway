import { render, screen, fireEvent } from "@testing-library/react"
import { describe, expect, it, vi } from "vitest"

import { SearchInput } from "./search-input"

describe("SearchInput Component", () => {
  it("renders correctly with search icon and placeholder", () => {
    render(<SearchInput placeholder="Search logs..." />)
    expect(screen.getByPlaceholderText("Search logs...")).toBeInTheDocument()
  })

  it("handles user typing input", () => {
    const handleChange = vi.fn()
    render(<SearchInput onChange={handleChange} />)
    const input = screen.getByPlaceholderText("Search...")
    fireEvent.change(input, { target: { value: "query" } })
    expect(handleChange).toHaveBeenCalled()
  })
})
