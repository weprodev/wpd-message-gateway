import { render } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { MessageGatewayLogo } from "./message-gateway-logo"

describe("MessageGatewayLogo", () => {
  it("renders decorative logo markup", () => {
    const { container } = render(<MessageGatewayLogo />)
    expect(container.querySelector('[aria-hidden="true"]')).toBeInTheDocument()
  })
})
