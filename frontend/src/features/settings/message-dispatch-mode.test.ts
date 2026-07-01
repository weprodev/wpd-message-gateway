import { describe, expect, it } from "vitest"

import {
  parseMessageDispatchConfig,
  toMessageDispatchApiValue,
} from "./message-dispatch-mode"

describe("parseMessageDispatchConfig", () => {
  it("maps canonical API values", () => {
    expect(parseMessageDispatchConfig("memory_only")).toEqual({ mode: "memory", storeInDb: false })
    expect(parseMessageDispatchConfig("memory_and_database")).toEqual({ mode: "memory", storeInDb: true })
    expect(parseMessageDispatchConfig("provider_only")).toEqual({ mode: "provider", storeInDb: false })
    expect(parseMessageDispatchConfig("provider_and_database")).toEqual({
      mode: "provider",
      storeInDb: true,
    })
  })

  it("defaults to memory without persistence", () => {
    expect(parseMessageDispatchConfig()).toEqual({ mode: "memory", storeInDb: false })
    expect(parseMessageDispatchConfig("unknown")).toEqual({ mode: "memory", storeInDb: false })
  })
})

describe("toMessageDispatchApiValue", () => {
  it("writes canonical gateway mode values", () => {
    expect(toMessageDispatchApiValue({ mode: "memory", storeInDb: false })).toBe("memory_only")
    expect(toMessageDispatchApiValue({ mode: "memory", storeInDb: true })).toBe("memory_and_database")
    expect(toMessageDispatchApiValue({ mode: "provider", storeInDb: false })).toBe("provider_only")
    expect(toMessageDispatchApiValue({ mode: "provider", storeInDb: true })).toBe("provider_and_database")
  })

  it("round-trips canonical values", () => {
    const cases = [
      "memory_only",
      "memory_and_database",
      "provider_only",
      "provider_and_database",
    ] as const

    for (const apiValue of cases) {
      expect(toMessageDispatchApiValue(parseMessageDispatchConfig(apiValue))).toBe(apiValue)
    }
  })
})
