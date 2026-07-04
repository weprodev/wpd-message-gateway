import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { RadioOption } from "./radio-option"

describe("RadioOption", () => {
  it("renders label and description", () => {
    render(
      <RadioOption
        id="dispatch-memory"
        name="dispatch-mode"
        label="Memory"
        description="Capture messages in memory."
        checked
        onChange={vi.fn()}
      />,
    )

    expect(screen.getByText("Memory")).toBeInTheDocument()
    expect(screen.getByText("Capture messages in memory.")).toBeInTheDocument()
    expect(screen.getByRole("radio")).toBeChecked()
  })

  it("calls onChange when selected", async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()

    render(
      <RadioOption
        id="dispatch-provider"
        name="dispatch-mode"
        label="Provider"
        description="Send through provider."
        checked={false}
        onChange={onChange}
      />,
    )

    await user.click(screen.getByRole("radio"))
    expect(onChange).toHaveBeenCalledTimes(1)
  })

  it("does not call onChange when disabled", async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()

    render(
      <RadioOption
        id="dispatch-provider"
        name="dispatch-mode"
        label="Provider"
        description="Send through provider."
        checked={false}
        onChange={onChange}
        disabled
      />,
    )

    expect(screen.getByRole("radio")).toBeDisabled()
    await user.click(screen.getByRole("radio"))
    expect(onChange).not.toHaveBeenCalled()
  })
})
