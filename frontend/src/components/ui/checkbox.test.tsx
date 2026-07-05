import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { Checkbox } from "./checkbox"

describe("Checkbox", () => {
  it("renders and toggles checked state", async () => {
    const user = userEvent.setup()
    const onCheckedChange = vi.fn()

    render(<Checkbox aria-label="Accept terms" onCheckedChange={onCheckedChange} />)

    await user.click(screen.getByRole("checkbox", { name: /accept terms/i }))
    expect(onCheckedChange).toHaveBeenCalledWith(true)
  })
})
