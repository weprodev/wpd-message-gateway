import { describe, expect, it } from "vitest"

import { buildGatewayCurlExamples } from "./get-started-dialog.utils"

describe("buildGatewayCurlExamples", () => {
  it("embeds workspace id and client id in curl commands", () => {
    const examples = buildGatewayCurlExamples(
      "00000000-0000-0000-0000-000000000001",
      "demo-client-id",
    )

    expect(examples.email).toContain("X-Workspace-Key: 00000000-0000-0000-0000-000000000001")
    expect(examples.email).toContain('-u "demo-client-id:YOUR_CLIENT_SECRET"')
    expect(examples.sms).toContain("X-Workspace-Key: 00000000-0000-0000-0000-000000000001")
    expect(examples.sms).toContain('-u "demo-client-id:YOUR_CLIENT_SECRET"')
  })
})
