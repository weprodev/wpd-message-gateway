import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, expect, it, vi } from "vitest"

import { MISSING_CLIENT_SECRET_MESSAGE } from "@/lib/errors"

import { ApiKeyRegenerateModal } from "./api-key-regenerate-modal"

const sampleCredentials = {
  clientId: "wk_test_client_id",
  clientSecret: "wk_test_client_secret",
  keyName: "Production",
  mode: "regenerated" as const,
}

describe("ApiKeyRegenerateModal", () => {
  it("renders confirmation message and actions when open", () => {
    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={vi.fn()}
      />,
    )

    expect(screen.getByText("Are you sure you want to regenerate this API key?")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Cancel" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Regenerate" })).toBeInTheDocument()
  })

  it("shows credentials in the same modal after regenerate succeeds", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn().mockResolvedValue(sampleCredentials)

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={onRegenerate}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Regenerate" }))

    expect(await screen.findByText("API key regenerated")).toBeInTheDocument()
    expect(screen.getByDisplayValue("wk_test_client_id")).toBeInTheDocument()
    expect(screen.queryByText("Are you sure you want to regenerate this API key?")).not.toBeInTheDocument()
  })

  it("calls onClose when cancel is clicked", async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={onClose}
        onRegenerate={vi.fn()}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Cancel" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it("blocks dismissal while regenerate is in progress", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn(() => new Promise<typeof sampleCredentials>(() => {}))

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={onRegenerate}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Regenerate" }))

    expect(await screen.findByText("Regenerating…")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: /close/i })).not.toBeInTheDocument()
  })

  it("starts on the confirm step when reopened, even if a prior regenerate finishes late", async () => {
    const user = userEvent.setup()
    let resolveRegenerate: ((value: typeof sampleCredentials) => void) | undefined
    const onRegenerate = vi.fn(
      () =>
        new Promise<typeof sampleCredentials>((resolve) => {
          resolveRegenerate = resolve
        }),
    )
    const onClose = vi.fn()

    const { rerender } = render(
      <ApiKeyRegenerateModal isOpen onClose={onClose} onRegenerate={onRegenerate} />,
    )

    await user.click(screen.getByRole("button", { name: "Regenerate" }))

    rerender(<ApiKeyRegenerateModal isOpen={false} onClose={onClose} onRegenerate={onRegenerate} />)
    rerender(<ApiKeyRegenerateModal isOpen onClose={onClose} onRegenerate={onRegenerate} />)

    expect(screen.getByText("Are you sure you want to regenerate this API key?")).toBeInTheDocument()
    expect(screen.queryByDisplayValue("wk_test_client_secret")).not.toBeInTheDocument()

    resolveRegenerate?.(sampleCredentials)
    await vi.waitFor(() => expect(onRegenerate).toHaveBeenCalled())

    expect(screen.getByText("Are you sure you want to regenerate this API key?")).toBeInTheDocument()
    expect(screen.queryByDisplayValue("wk_test_client_secret")).not.toBeInTheDocument()
  })

  it("shows an error when regenerate fails", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn().mockRejectedValue(new Error("network error"))

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={onRegenerate}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Regenerate" }))

    expect(await screen.findByText("network error")).toBeInTheDocument()
    expect(screen.getByText("Are you sure you want to regenerate this API key?")).toBeInTheDocument()
  })

  it("shows an error when the server omits the client secret", async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn().mockRejectedValue(new Error(MISSING_CLIENT_SECRET_MESSAGE))

    render(
      <ApiKeyRegenerateModal
        isOpen
        onClose={vi.fn()}
        onRegenerate={onRegenerate}
      />,
    )

    await user.click(screen.getByRole("button", { name: "Regenerate" }))

    expect(await screen.findByText(MISSING_CLIENT_SECRET_MESSAGE)).toBeInTheDocument()
    expect(screen.getByText("Are you sure you want to regenerate this API key?")).toBeInTheDocument()
  })
})
