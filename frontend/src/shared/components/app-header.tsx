import { MessageGatewayLogo } from "./message-gateway-logo"
import { ThemeToggle } from "./theme-toggle"

function BrandText() {
  return (
    <div className="flex flex-col gap-1">
      <p className="text-xl font-semibold leading-normal text-primary-brand">Message Gateway</p>
      <p className="text-xs leading-4 text-text-secondary">Connecting your world seamlessly!</p>
    </div>
  )
}

export function AppHeader() {
  return (
    <header className="w-full border-b border-divider bg-surface">
      <div className="flex items-center justify-between px-12 py-6">
        <div className="flex items-center gap-3">
          <MessageGatewayLogo />
          <BrandText />
        </div>
        <ThemeToggle />
      </div>
    </header>
  )
}
