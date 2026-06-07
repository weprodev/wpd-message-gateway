import { ROUTES } from "@/core/router/routes"

export type WorkspaceNavSegment = Exclude<keyof typeof ROUTES.workspace, "pattern">

export function workspaceHref(workspaceId: string, segment: WorkspaceNavSegment) {
  return ROUTES.workspace[segment](workspaceId)
}

export interface WorkspaceNavItem {
  readonly segment: WorkspaceNavSegment
  readonly label: string
  readonly icon: string
  readonly disabled?: boolean
}

export interface WorkspaceNavSection {
  readonly label?: string
  readonly items: readonly WorkspaceNavItem[]
}

export const WORKSPACE_NAV_SECTIONS: readonly WorkspaceNavSection[] = [
  {
    label: "Navigation",
    items: [
      { segment: "overview", label: "Overview", icon: "dashboard" },
      { segment: "email", label: "Email", icon: "mail" },
      { segment: "sms", label: "SMS", icon: "forum" },
      { segment: "push", label: "Push", icon: "notifications" },
      { segment: "chat", label: "Chat", icon: "chat" },
    ],
  },
  {
    items: [
      { segment: "integrations", label: "Integrations", icon: "extension" },
      { segment: "settings", label: "Settings", icon: "settings" },
    ],
  },
] as const
