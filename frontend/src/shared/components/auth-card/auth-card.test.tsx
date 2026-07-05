import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { AuthCard } from "./auth-card"

describe("AuthCard", () => {
  it("renders children inside the card container", () => {
    render(
      <AuthCard>
        <p>Portal login</p>
      </AuthCard>,
    )

    expect(screen.getByText("Portal login")).toBeInTheDocument()
  })
})
