import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { IntegrationProviderIcon } from "./integration-provider-icon"

describe("IntegrationProviderIcon", () => {
  it("renders emoji icon", () => {
    render(<IntegrationProviderIcon icon="📧" name="Mailgun" className="size-12" />)
    expect(screen.getByText("📧")).toBeInTheDocument()
  })

  it("renders image icon with alt text", () => {
    render(<IntegrationProviderIcon icon="/providers/mailgun.svg" name="Mailgun" className="size-12" />)
    expect(screen.getByRole("img", { name: "Mailgun logo" })).toHaveAttribute("src", "/providers/mailgun.svg")
  })
})
