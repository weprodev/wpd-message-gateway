/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useMemo, type ReactNode } from "react"

import {
  createWorkspaceAuthorization,
  type WorkspaceAuthorization,
  type WorkspaceAuthorizationScope,
} from "./authorization"

const WorkspaceAuthorizationContext = createContext<WorkspaceAuthorization | null>(null)

type WorkspaceAuthorizationProviderProps = WorkspaceAuthorizationScope & {
  children: ReactNode
}

export function WorkspaceAuthorizationProvider({
  role,
  permissions,
  children,
}: WorkspaceAuthorizationProviderProps) {
  const value = useMemo(
    () => createWorkspaceAuthorization({ role, permissions }),
    [role, permissions],
  )

  return (
    <WorkspaceAuthorizationContext.Provider value={value}>
      {children}
    </WorkspaceAuthorizationContext.Provider>
  )
}

export function useWorkspaceAuthorization(): WorkspaceAuthorization {
  const context = useContext(WorkspaceAuthorizationContext)
  if (!context) {
    throw new Error(
      "useWorkspaceAuthorization must be used within WorkspaceAuthorizationProvider",
    )
  }
  return context
}
