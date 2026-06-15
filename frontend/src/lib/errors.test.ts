import { describe, expect, it } from "vitest"

import { httpError, toUserMessage } from "./errors"

describe("httpError", () => {
  it("returns the API message for client errors", async () => {
    const response = new Response(JSON.stringify({ message: "name required" }), { status: 400 })

    await expect(httpError(response, "Failed to create API key")).resolves.toEqual(
      new Error("name required"),
    )
  })

  it("returns the fallback for client errors without a message", async () => {
    const response = new Response("{}", { status: 404 })

    await expect(httpError(response, "Failed to create API key")).resolves.toEqual(
      new Error("Failed to create API key"),
    )
  })

  it("returns the fallback for server errors even when the body includes a message", async () => {
    const response = new Response(
      JSON.stringify({ message: "pq: relation \"api_keys\" does not exist" }),
      { status: 500 },
    )

    await expect(httpError(response, "Failed to create API key")).resolves.toEqual(
      new Error("Failed to create API key"),
    )
  })

  it("returns the fallback for non-JSON server error bodies", async () => {
    const response = new Response("Internal Server Error", { status: 502 })

    await expect(httpError(response, "Failed to regenerate API key")).resolves.toEqual(
      new Error("Failed to regenerate API key"),
    )
  })
})

describe("toUserMessage", () => {
  it("returns the error message when the value is an Error", () => {
    expect(toUserMessage(new Error("network error"), "fallback")).toBe("network error")
  })

  it("returns the fallback for unknown values", () => {
    expect(toUserMessage("boom", "fallback")).toBe("fallback")
  })
})

describe("requireClientSecret", () => {
  it("returns the secret when present", async () => {
    const { requireClientSecret } = await import("./errors")
    expect(requireClientSecret({ client_secret: "wk_secret" })).toBe("wk_secret")
  })

  it("throws when the secret is missing or blank", async () => {
    const { requireClientSecret, MISSING_CLIENT_SECRET_MESSAGE } = await import("./errors")

    expect(() => requireClientSecret({})).toThrow(MISSING_CLIENT_SECRET_MESSAGE)
    expect(() => requireClientSecret({ client_secret: "   " })).toThrow(MISSING_CLIENT_SECRET_MESSAGE)
  })
})
