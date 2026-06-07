import { Navigate, Outlet } from "react-router-dom"
import { getToken } from "@/core/api/client"
import { ROUTES } from "@/core/router/routes"

export function ProtectedRoute() {
  const token = getToken()

  if (!token) {
    return <Navigate to={ROUTES.login} replace />
  }

  return <Outlet />
}
