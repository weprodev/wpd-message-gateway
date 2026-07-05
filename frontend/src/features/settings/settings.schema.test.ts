import { describe, expect, it } from "vitest"

import { parseWorkspaceSettings } from "./settings.schema"

describe("parseWorkspaceSettings", () => {
  it("parses valid workspace settings", () => {
    const result = parseWorkspaceSettings({
      owner_email: "owner@example.com",
      pin_code: "123456",
      message_dispatch_mode: "provider",
      store_message_content: "true",
    })

    expect(result).toEqual({
      owner_email: "owner@example.com",
      pin_code: "123456",
      message_dispatch_mode: "provider",
      store_message_content: "true",
    })
  })

  it("accepts an empty object", () => {
    expect(parseWorkspaceSettings({})).toEqual({})
  })

  it("throws when settings contain invalid values", () => {
    expect(() =>
      parseWorkspaceSettings({ message_dispatch_mode: "invalid-mode" }),
    ).toThrow("Invalid settings response from server")
  })
})
