import { createBrowserRouter, Navigate } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { ProtectedRoute } from "@/core/router/protected-route"
import { LoginPage, RegisterPage } from "@/features/auth"
import { IntegrationsPage } from "@/features/integrations"
import { OverviewPage, EmailInboxPage, EmailTemplatesPage } from "@/features/inbox"
import { SettingsPage } from "@/features/settings"
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
      {
        element: <ProtectedRoute />,
        children: [
          { path: "workspaces", element: <WorkspacesPage /> },
        ],
      },
    ],
  },
  {
    path: ROUTES.workspace.pattern,
    element: <ProtectedRoute />,
    children: [
      {
        element: <WorkspaceLayout />,
        children: [
          { index: true, element: <Navigate to="overview" replace /> },
          { path: "overview", element: <OverviewPage /> },
          { path: "email", element: <EmailInboxPage /> },
          { path: "email/templates", element: <EmailTemplatesPage /> },
          { path: "sms", element: <OverviewPage channel="sms" /> },
          { path: "push", element: <OverviewPage channel="push" /> },
          { path: "chat", element: <OverviewPage channel="chat" /> },
          { path: "integrations", element: <IntegrationsPage /> },
          { path: "settings", element: <SettingsPage /> },
        ],
      },
    ],
  },
])
