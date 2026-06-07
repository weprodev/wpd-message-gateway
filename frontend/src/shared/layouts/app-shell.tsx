import { Outlet } from "react-router-dom"

import { AppFooter } from "@/shared/components/app-footer"
import { AppHeader } from "@/shared/components/app-header"

export function AppShell() {
  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <AppHeader />
      <main className="flex w-full flex-1 flex-col items-center">
        <Outlet />
      </main>
      <AppFooter />
    </div>
  )
}
