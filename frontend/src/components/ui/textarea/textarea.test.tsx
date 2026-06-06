import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { Textarea } from "./textarea"

describe("Textarea Component", () => {
  it("renders correctly with placeholder", () => {
    render(<Textarea placeholder="Type something here" />)
    expect(screen.getByPlaceholderText("Type something here")).toBeInTheDocument()
  })

  it("can be disabled", () => {
    render(<Textarea disabled placeholder="Disabled textarea" />)
    expect(screen.getByPlaceholderText("Disabled textarea")).toBeDisabled()
  })
})
