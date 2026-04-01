import { createBrowserRouter } from "react-router-dom"

import { HomeRedirect } from "@/app/home-redirect"
import { ROUTES } from "@/app/paths"
import { LoginPage } from "@/features/auth/pages/login.page"
import { RegisterPage } from "@/features/auth/pages/register.page"
import { EmailOverviewPage } from "@/features/email/pages/email-overview.page"
import { WorkspacesPage } from "@/features/workspaces/pages/workspaces.page"
import { AppShell } from "@/shared/layouts/app-shell"

export const router = createBrowserRouter([
  {
    path: ROUTES.root,
    element: <AppShell />,
    children: [
      { index: true, element: <HomeRedirect /> },
      { path: "login", element: <LoginPage /> },
      { path: "register", element: <RegisterPage /> },
      { path: "workspaces", element: <WorkspacesPage /> },
      { path: "email", element: <EmailOverviewPage /> },
    ],
  },
])
