import { ROUTES } from "@/core/router/routes"

export type WorkspaceNavSegment = Exclude<keyof typeof ROUTES.workspace, "pattern">

export function workspaceHref(workspaceId: string, segment: WorkspaceNavSegment) {
  return ROUTES.workspace[segment](workspaceId)
}

export const WORKSPACE_NAV_ITEMS = [
  { segment: "overview" as const, label: "Overview", icon: "grid_view" },
  { segment: "email" as const, label: "Email", icon: "mail" },
  { segment: "sms" as const, label: "SMS", icon: "forum" },
  { segment: "push" as const, label: "Push", icon: "notifications" },
  { segment: "chat" as const, label: "Chat", icon: "chat" },
] as const
