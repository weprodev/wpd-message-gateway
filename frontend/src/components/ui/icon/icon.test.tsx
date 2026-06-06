import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { Icon } from "./icon"

describe("Icon Component", () => {
  it("renders the icon with correct name as content", () => {
    render(<Icon name="settings" aria-label="Settings Icon" />)
    const iconSpan = screen.getByLabelText("Settings Icon")
    expect(iconSpan).toBeInTheDocument()
    expect(iconSpan).toHaveTextContent("settings")
    expect(iconSpan).toHaveClass("material-symbols-outlined")
  })

  it("applies the size classes correctly", () => {
    const { container } = render(<Icon name="mail" size="sm" />)
    expect(container.firstChild).toHaveClass("text-[18px]")
  })
})
