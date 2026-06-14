import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"

import { IntegrationProviderIcon } from "../integration-provider-icon"
import type { IntegrationViewModel } from "../../hooks/use-integrations.hook"

interface IntegrationRowProps {
  provider: IntegrationViewModel
  onConnect: (provider: IntegrationViewModel) => void
  onActivate: (provider: IntegrationViewModel) => void
  onDisconnect: (provider: IntegrationViewModel) => void
  isBusy?: boolean
}

export function IntegrationRow({ provider, onConnect, onActivate, onDisconnect, isBusy }: IntegrationRowProps) {
  const isDisabled = !provider.isAvailable || isBusy

  return (
    <div
      className={cn(
        "flex items-center justify-between gap-4 border-b border-border px-4 py-4 last:border-b-0",
        !provider.isAvailable && "opacity-50",
      )}
    >
      <div className="flex flex-1 items-center gap-4">
        <IntegrationProviderIcon
          icon={provider.icon}
          name={provider.name}
          className="size-12 p-2.5 text-2xl"
        />
        <div className="min-w-0 flex-1">
          <p className="text-sm font-medium text-foreground">{provider.name}</p>
          <p className="mt-1 text-[13px] text-text-secondary">{provider.description}</p>
        </div>
      </div>

      <div className="flex shrink-0 items-center gap-3">
        {provider.isComingSoon ? (
          <Badge variant="secondary">Coming soon</Badge>
        ) : provider.isConnected ? (
          <>
            <Badge className="border-success/30 bg-success-bg text-success">Connected</Badge>
            <Button
              type="button"
              variant="outline"
              size="sm"
              disabled={isDisabled}
              onClick={() => onDisconnect(provider)}
            >
              Disconnect
            </Button>
          </>
        ) : provider.isDeactivated ? (
          <>
            <Badge variant="secondary">Deactivated</Badge>
            <Button
              type="button"
              size="sm"
              disabled={isDisabled}
              className="bg-primary-brand hover:bg-primary-brand-hover"
              onClick={() => onActivate(provider)}
            >
              Activate
            </Button>
          </>
        ) : (
          <Button
            type="button"
            size="sm"
            disabled={isDisabled}
            className="bg-primary-brand hover:bg-primary-brand-hover"
            onClick={() => onConnect(provider)}
          >
            Connect
          </Button>
        )}
      </div>
    </div>
  )
}
