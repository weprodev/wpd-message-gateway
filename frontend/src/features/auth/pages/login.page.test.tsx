import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { MemoryRouter } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { signInWithPassword } from "../api/auth.api"
import { LoginPage } from "./login.page"

vi.mock("../api/auth.api", () => ({
  signInWithPassword: vi.fn(),
}))

const mockedSignIn = vi.mocked(signInWithPassword)

describe("LoginPage", () => {
  it("renders sign-in form", () => {
    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    )

    expect(screen.getByRole("heading", { name: /sign in/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument()
  })

  it("shows API error on failed login", async () => {
    mockedSignIn.mockResolvedValue({ ok: false, message: "Invalid credentials" })

    render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>,
    )

    const user = userEvent.setup()
    await user.type(screen.getByLabelText(/email address/i), "user@example.com")
    await user.type(screen.getByLabelText(/^password$/i), "wrong")
    await user.click(screen.getByRole("button", { name: /^login$/i }))

    expect(await screen.findByText("Invalid credentials")).toBeInTheDocument()
  })
})
