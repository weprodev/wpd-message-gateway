import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { AppFooter } from "./app-footer"

describe("AppFooter", () => {
  it("renders attribution link", () => {
    render(<AppFooter />)

    const link = screen.getByRole("link", { name: "WeProDev" })
    expect(link).toHaveAttribute("href", "https://weprodev.com")
  })

  it("applies dashboard link styling", () => {
    render(<AppFooter variant="dashboard" />)

    expect(screen.getByRole("link", { name: "WeProDev" })).toHaveClass("text-indigo-400")
  })
})
