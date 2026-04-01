import { AUTH_PATHS } from "@/features/auth/paths"
import { EMAIL_PATHS } from "@/features/email/paths"
import { WORKSPACES_PATHS } from "@/features/workspaces/paths"

export const ROUTES = {
  root: "/",
  ...AUTH_PATHS,
  workspaces: WORKSPACES_PATHS.list,
  email: EMAIL_PATHS.overview,
} as const
