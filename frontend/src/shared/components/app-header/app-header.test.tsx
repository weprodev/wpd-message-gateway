import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"
import { MemoryRouter } from "react-router-dom"

import { AppHeader } from "./app-header"

const { navigate, setToken } = vi.hoisted(() => ({
  navigate: vi.fn(),
  setToken: vi.fn(),
}))

vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual<typeof import("react-router-dom")>("react-router-dom")
  return {
    ...actual,
    useNavigate: () => navigate,
  }
})

vi.mock("@/core/api/client", () => ({
  getToken: vi.fn(() => "portal-jwt"),
  setToken,
}))

vi.mock("@/shared/context/theme-context", () => ({
  useTheme: () => ({
    theme: "light",
    toggleTheme: vi.fn(),
  }),
}))

describe("AppHeader", () => {
  it("renders brand and sign out when authenticated", async () => {
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <AppHeader />
      </MemoryRouter>,
    )

    expect(screen.getByText("Message Gateway")).toBeInTheDocument()
    await user.click(screen.getByRole("button", { name: /sign out/i }))

    expect(setToken).toHaveBeenCalledWith(null)
    expect(navigate).toHaveBeenCalled()
  })
})
