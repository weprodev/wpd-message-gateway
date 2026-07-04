import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { MemoryRouter } from "react-router-dom"
import { describe, expect, it, vi } from "vitest"

import { registerAccount } from "../auth.api"
import { RegisterPage } from "./register.page"

vi.mock("../auth.api", () => ({
  registerAccount: vi.fn(),
}))

const mockedRegister = vi.mocked(registerAccount)

describe("RegisterPage", () => {
  it("renders sign-up form", () => {
    render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>,
    )

    expect(screen.getByRole("heading", { name: /sign up/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/full name/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument()
  })

  it("shows API error on failed registration", async () => {
    mockedRegister.mockResolvedValue({ ok: false, message: "Email already taken" })

    render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>,
    )

    const user = userEvent.setup()
    await user.type(screen.getByLabelText(/full name/i), "Demo User")
    await user.type(screen.getByLabelText(/email address/i), "user@example.com")
    await user.type(screen.getByLabelText(/^password$/i), "secret123456")
    await user.click(screen.getByRole("button", { name: /sign up/i }))

    expect(await screen.findByText("Email already taken")).toBeInTheDocument()
  })
})
