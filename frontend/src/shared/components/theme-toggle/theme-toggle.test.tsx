import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { ThemeToggle } from "./theme-toggle"

const { toggleTheme } = vi.hoisted(() => ({
  toggleTheme: vi.fn(),
}))

vi.mock("@/shared/context/theme-context", () => ({
  useTheme: () => ({
    theme: "light",
    toggleTheme,
  }),
}))

describe("ThemeToggle", () => {
  it("calls toggleTheme when clicked", async () => {
    const user = userEvent.setup()

    render(<ThemeToggle />)

    await user.click(screen.getByRole("button", { name: /toggle theme/i }))
    expect(toggleTheme).toHaveBeenCalledTimes(1)
  })
})
