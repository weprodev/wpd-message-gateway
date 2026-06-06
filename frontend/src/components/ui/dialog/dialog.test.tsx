import { render, screen, fireEvent } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { Dialog, DialogContent, DialogDescription, DialogTitle, DialogTrigger } from "./dialog"

describe("Dialog Component", () => {
  it("renders trigger and displays content when clicked", async () => {
    render(
      <Dialog>
        <DialogTrigger>Open</DialogTrigger>
        <DialogContent>
          <DialogTitle>Dialog Title</DialogTitle>
          <DialogDescription>Dialog Description</DialogDescription>
        </DialogContent>
      </Dialog>
    )

    const trigger = screen.getByText("Open")
    expect(trigger).toBeInTheDocument()
    expect(screen.queryByText("Dialog Title")).not.toBeInTheDocument()

    fireEvent.click(trigger)

    expect(await screen.findByText("Dialog Title")).toBeInTheDocument()
    expect(screen.getByText("Dialog Description")).toBeInTheDocument()
  })
})
