import { describe, expect, it } from "vitest"

import {
  dispatchConfigsEqual,
  normalizeDispatchMode,
  normalizeStoreMessageContent,
  parseMessageDispatchConfig,
  toStoreMessageContentSetting,
} from "./message-dispatch-mode"

describe("parseMessageDispatchConfig", () => {
  it("maps valid combinations", () => {
    expect(parseMessageDispatchConfig("memory", "false")).toEqual({ mode: "memory", storeMessageContent: false })
    expect(parseMessageDispatchConfig("memory", "true")).toEqual({ mode: "memory", storeMessageContent: true })
    expect(parseMessageDispatchConfig("provider", "false")).toEqual({ mode: "provider", storeMessageContent: false })
    expect(parseMessageDispatchConfig("provider", "true")).toEqual({ mode: "provider", storeMessageContent: true })
  })

  it("defaults to memory without inbox capture on missing input", () => {
    expect(parseMessageDispatchConfig()).toEqual({ mode: "memory", storeMessageContent: false })
    expect(parseMessageDispatchConfig("unknown", "invalid")).toEqual({ mode: "memory", storeMessageContent: false })
  })
})

describe("normalizeDispatchMode", () => {
  it("is case-insensitive for provider", () => {
    expect(normalizeDispatchMode("PROVIDER")).toBe("provider")
  })
})

describe("normalizeStoreMessageContent", () => {
  it("treats unknown strings as false", () => {
    expect(normalizeStoreMessageContent("yes")).toBe(false)
  })
})

describe("toStoreMessageContentSetting", () => {
  it("maps boolean to API string values", () => {
    expect(toStoreMessageContentSetting(true)).toBe("true")
    expect(toStoreMessageContentSetting(false)).toBe("false")
  })
})

describe("dispatchConfigsEqual", () => {
  it("compares mode and store flag", () => {
    expect(dispatchConfigsEqual({ mode: "memory", storeMessageContent: false }, { mode: "memory", storeMessageContent: false })).toBe(true)
    expect(dispatchConfigsEqual({ mode: "provider", storeMessageContent: true }, { mode: "provider", storeMessageContent: false })).toBe(false)
  })
})
