import { useNavigate } from "react-router-dom"
import { getToken, setToken } from "@/core/api/client"
import { ROUTES } from "@/core/router/routes"
import { Icon } from "@/components/ui/icon"
import { MessageGatewayLogo } from "../message-gateway-logo"
import { ThemeToggle } from "../theme-toggle"

function BrandText() {
  return (
    <div className="flex flex-col gap-1">
      <p className="text-xl font-semibold leading-normal text-primary-brand">Message Gateway</p>
      <p className="text-xs leading-4 text-text-secondary">Connecting your world seamlessly!</p>
    </div>
  )
}

export function AppHeader() {
  const navigate = useNavigate()
  const hasToken = Boolean(getToken())

  const handleSignOut = () => {
    setToken(null)
    navigate(ROUTES.login, { replace: true })
  }

  return (
    <header className="w-full border-b border-divider bg-surface">
      <div className="flex items-center justify-between px-12 py-6">
        <div className="flex items-center gap-3">
          <MessageGatewayLogo />
          <BrandText />
        </div>
        <div className="flex items-center gap-3">
          <ThemeToggle />
          {hasToken ? (
            <button
              type="button"
              onClick={handleSignOut}
              className="flex cursor-pointer items-center gap-2 rounded-lg border border-border bg-input px-3 py-1.5 text-xs font-semibold text-text-secondary transition-colors hover:bg-secondary-button-hover"
              aria-label="Sign out"
            >
              <Icon name="logout" size="sm" />
              <span>Sign out</span>
            </button>
          ) : null}
        </div>
      </div>
    </header>
  )
}

