import { Navigate } from "react-router-dom"

import { getToken } from "@/core/api/client"
import { ROUTES } from "@/app/paths"

export function HomeRedirect() {
  return <Navigate to={getToken() ? ROUTES.workspaces : ROUTES.login} replace />
}
