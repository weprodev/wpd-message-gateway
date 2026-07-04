import { describe, expect, it } from "vitest"

import type { LogRow } from "./inbox.types"
import { hasInboxContent, isEmailInboxLink } from "./log-display.utils"

const baseLog: LogRow = {
  id: "log-1",
  workspace_id: "ws-1",
  channel_type: "email",
  http_method: "POST",
  status_code: 200,
  endpoint: "/v1/email",
  created_at: "2026-01-01T00:00:00Z",
}

describe("hasInboxContent", () => {
  it("returns true when inbox_message_id is present", () => {
    expect(hasInboxContent({ ...baseLog, inbox_message_id: "msg-1" })).toBe(true)
  })

  it("returns false when inbox_message_id is missing", () => {
    expect(hasInboxContent(baseLog)).toBe(false)
  })
})

describe("isEmailInboxLink", () => {
  it("returns true for email logs with inbox content", () => {
    expect(isEmailInboxLink({ ...baseLog, inbox_message_id: "msg-1" })).toBe(true)
  })

  it("returns false for non-email channels", () => {
    expect(isEmailInboxLink({ ...baseLog, channel_type: "sms", inbox_message_id: "msg-1" })).toBe(
      false,
    )
  })

  it("returns false for email logs without inbox content", () => {
    expect(isEmailInboxLink(baseLog)).toBe(false)
  })

  it("matches email channel type case-insensitively", () => {
    expect(isEmailInboxLink({ ...baseLog, channel_type: "EMAIL", inbox_message_id: "msg-1" })).toBe(
      true,
    )
  })
})
