import type { ReactNode } from "react"

import { useWorkspaceAuthorization } from "./workspace-authorization-context"

type CanBaseProps = {
  children: ReactNode
  fallback?: ReactNode
}

type CanByPermission = CanBaseProps & {
  permission: string
  anyPermission?: never
  role?: never
  anyRole?: never
}

type CanByAnyPermission = CanBaseProps & {
  anyPermission: readonly string[]
  permission?: never
  role?: never
  anyRole?: never
}

type CanByRole = CanBaseProps & {
  role: string
  anyRole?: never
  permission?: never
  anyPermission?: never
}

type CanByAnyRole = CanBaseProps & {
  anyRole: readonly string[]
  role?: never
  permission?: never
  anyPermission?: never
}

export type CanProps = CanByPermission | CanByAnyPermission | CanByRole | CanByAnyRole

export function Can(props: CanProps) {
  const auth = useWorkspaceAuthorization()
  const { children, fallback = null } = props

  let allowed = false

  if ("permission" in props && props.permission) {
    allowed = auth.can(props.permission)
  } else if ("anyPermission" in props && props.anyPermission) {
    allowed = auth.canAny(...props.anyPermission)
  } else if ("role" in props && props.role) {
    allowed = auth.hasRole(props.role)
  } else if ("anyRole" in props && props.anyRole) {
    allowed = auth.hasRole(...props.anyRole)
  }

  return allowed ? children : fallback
}
