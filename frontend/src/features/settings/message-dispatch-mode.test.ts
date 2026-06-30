import { describe, expect, it } from "vitest"

import {
  parseMessageDispatchConfig,
  toMessageDispatchApiValue,
} from "./message-dispatch-mode"

describe("parseMessageDispatchConfig", () => {
  it("maps canonical API values", () => {
    expect(parseMessageDispatchConfig("memory_only")).toEqual({ mode: "memory", storeInDb: false })
    expect(parseMessageDispatchConfig("memory_and_provider")).toEqual({ mode: "memory", storeInDb: true })
    expect(parseMessageDispatchConfig("provider_only")).toEqual({ mode: "provider", storeInDb: false })
    expect(parseMessageDispatchConfig("provider_and_database")).toEqual({
      mode: "provider",
      storeInDb: true,
    })
  })

  it("maps legacy aliases", () => {
    expect(parseMessageDispatchConfig("memory")).toEqual({ mode: "memory", storeInDb: false })
    expect(parseMessageDispatchConfig("both")).toEqual({ mode: "memory", storeInDb: true })
    expect(parseMessageDispatchConfig("memory_database")).toEqual({ mode: "memory", storeInDb: true })
    expect(parseMessageDispatchConfig("providers")).toEqual({ mode: "provider", storeInDb: false })
    expect(parseMessageDispatchConfig("provider_database")).toEqual({ mode: "provider", storeInDb: true })
  })

  it("defaults to memory without persistence", () => {
    expect(parseMessageDispatchConfig()).toEqual({ mode: "memory", storeInDb: false })
    expect(parseMessageDispatchConfig("unknown")).toEqual({ mode: "memory", storeInDb: false })
  })
})

describe("toMessageDispatchApiValue", () => {
  it("round-trips canonical values", () => {
    const cases = [
      "memory_only",
      "memory_and_provider",
      "provider_only",
      "provider_and_database",
    ] as const

    for (const apiValue of cases) {
      expect(toMessageDispatchApiValue(parseMessageDispatchConfig(apiValue))).toBe(apiValue)
    }
  })
})
