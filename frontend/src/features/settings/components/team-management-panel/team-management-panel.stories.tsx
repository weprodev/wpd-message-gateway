import type { Meta, StoryObj } from "@storybook/react-vite"
import { useEffect, type ReactNode } from "react"

import { AllPermissions, Role, WorkspaceAuthorizationProvider } from "@/core/auth"

import type { WorkspaceInvitation, WorkspaceMember } from "../../team.types"
import { TeamManagementPanel } from "./team-management-panel"

const workspaceId = "00000000-0000-0000-0000-000000000001"
const currentUserId = "00000000-0000-0000-0000-000000000010"

const mockMembers: WorkspaceMember[] = [
  {
    workspace_id: workspaceId,
    user_id: currentUserId,
    role: Role.Admin,
    joined_at: "2026-01-01T00:00:00Z",
    user_email: "demo@weprodev.com",
    display_name: "Demo Admin",
  },
  {
    workspace_id: workspaceId,
    user_id: "00000000-0000-0000-0000-000000000011",
    role: Role.Member,
    joined_at: "2026-01-02T00:00:00Z",
    user_email: "member@weprodev.com",
    display_name: "Demo Member",
  },
  {
    workspace_id: workspaceId,
    user_id: "00000000-0000-0000-0000-000000000012",
    role: Role.Viewer,
    joined_at: "2026-01-03T00:00:00Z",
    user_email: "viewer@weprodev.com",
    display_name: "Demo Viewer",
  },
]

const mockInvitations: WorkspaceInvitation[] = [
  {
    id: "inv-1",
    workspace_id: workspaceId,
    email: "pending@example.com",
    role: Role.Viewer,
    expires_at: "2026-02-01T12:00:00Z",
    status: "pending",
    created_at: "2026-01-05T00:00:00Z",
  },
]

type TeamFetchMock = {
  members?: WorkspaceMember[]
  invitations?: WorkspaceInvitation[]
  pending?: boolean
  errorMessage?: string
}

function resolveUrl(input: RequestInfo | URL): string {
  if (typeof input === "string") return input
  if (input instanceof URL) return input.toString()
  return input.url
}

function jsonResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response
}

function installTeamFetchMock(mock: TeamFetchMock) {
  const originalFetch = window.fetch

  window.fetch = async (input, init) => {
    const url = resolveUrl(input)

    if (mock.pending) {
      if (
        url.includes("/members") ||
        url.includes("/invitations") ||
        url.includes("/auth/me")
      ) {
        return new Promise(() => undefined)
      }
    }

    if (mock.errorMessage) {
      if (url.includes("/members") || url.includes("/invitations") || url.includes("/auth/me")) {
        return jsonResponse({ message: mock.errorMessage }, 500)
      }
    }

    if (url.includes("/members")) {
      return jsonResponse(mock.members ?? [])
    }

    if (url.includes("/invitations")) {
      return jsonResponse(mock.invitations ?? [])
    }

    if (url.includes("/auth/me")) {
      return jsonResponse({
        id: currentUserId,
        email: "demo@weprodev.com",
        first_name: "Demo",
        last_name: "Admin",
      })
    }

    return originalFetch(input, init)
  }

  return () => {
    window.fetch = originalFetch
  }
}

function withFetchMock(mock: TeamFetchMock) {
  return (Story: () => ReactNode) => {
    function FetchMockWrapper() {
      useEffect(() => installTeamFetchMock(mock), [])
      return <Story />
    }
    return <FetchMockWrapper />
  }
}

const meta = {
  title: "Features/Settings/TeamManagementPanel",
  component: TeamManagementPanel,
  parameters: { layout: "padded" },
  decorators: [
    (Story) => (
      <WorkspaceAuthorizationProvider role={Role.Admin} permissions={[...AllPermissions]}>
        <div className="max-w-3xl">
          <Story />
        </div>
      </WorkspaceAuthorizationProvider>
    ),
  ],
  args: {
    workspaceId,
    enabled: true,
  },
} satisfies Meta<typeof TeamManagementPanel>

export default meta
type Story = StoryObj<typeof meta>

export const Loading: Story = {
  decorators: [withFetchMock({ pending: true })],
}

export const EmptyTeam: Story = {
  decorators: [withFetchMock({ members: [], invitations: [] })],
}

export const PopulatedTeam: Story = {
  decorators: [
    withFetchMock({ members: mockMembers, invitations: mockInvitations }),
  ],
}

export const LoadError: Story = {
  decorators: [withFetchMock({ errorMessage: "Failed to load team" })],
}
