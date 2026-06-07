import { render, screen } from "@testing-library/react"
import { MemoryRouter } from "react-router-dom"
import { describe, expect, it } from "vitest"

import { SidebarNav, type SidebarNavSection } from "./sidebar-nav"

const mockSections: SidebarNavSection[] = [
  {
    label: "Navigation",
    items: [
      { label: "Inbox", segment: "email", icon: "mail" },
      { label: "SMS", segment: "sms", icon: "sms" },
    ],
  },
]

describe("SidebarNav Component", () => {
  it("renders the navigation items correctly", () => {
    render(
      <MemoryRouter>
        <SidebarNav
          sections={mockSections}
          workspaceId="ws-123"
          buildHref={(wid, seg) => `/workspaces/${wid}/${seg}`}
        />
      </MemoryRouter>,
    )

    expect(screen.getByText("Inbox")).toBeInTheDocument()
    expect(screen.getByText("SMS")).toBeInTheDocument()
    expect(screen.getByText("Inbox").closest("a")).toHaveAttribute(
      "href",
      "/workspaces/ws-123/email",
    )
  })
})
