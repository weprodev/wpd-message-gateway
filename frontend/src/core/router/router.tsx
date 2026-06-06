import { createBrowserRouter, Navigate } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { LoginPage, RegisterPage } from "@/features/auth"
import { OverviewPage } from "@/features/inbox"
import { WorkspaceLayout, WorkspacesPage } from "@/features/workspaces"
import { AppShell } from "@/shared/layouts/app-shell"

export const router = createBrowserRouter([
  {
    path: ROUTES.root,
    element: <AppShell />,
    children: [
      { index: true, element: <Navigate to={ROUTES.workspaces} replace /> },
      { path: "login", element: <LoginPage /> },
      { path: "register", element: <RegisterPage /> },
      { path: "workspaces", element: <WorkspacesPage /> },
    ],
  },
  {
    path: ROUTES.workspace.pattern,
    element: <WorkspaceLayout />,
    children: [
      { index: true, element: <Navigate to="overview" replace /> },
      { path: "overview", element: <OverviewPage /> },
      { path: "email", element: <OverviewPage channel="email" /> },
      { path: "sms", element: <OverviewPage channel="sms" /> },
      { path: "push", element: <OverviewPage channel="push" /> },
      { path: "chat", element: <OverviewPage channel="chat" /> },
    ],
  },
])
