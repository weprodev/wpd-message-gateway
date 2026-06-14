import * as Tabs from "@radix-ui/react-tabs"
import { useParams } from "react-router-dom"
import { useState } from "react"

import { PageHeader } from "@/shared/components/page-header"
import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { cn } from "@/lib/utils"

import { IntegrationRow } from "../components/integration-row"
import { ConnectModal } from "../components/connect-modal"
import { DisconnectModal } from "../components/disconnect-modal"
import {
  filterIntegrationsByTab,
  groupByCategory,
  useIntegrations,
} from "../hooks/use-integrations.hook"
import type { IntegrationChannel } from "../integrations.types"

const CATEGORY_LABELS: Record<IntegrationChannel, string> = {
  email: "Email",
  sms: "SMS",
  push: "Push",
  chat: "Chat",
}

export function IntegrationsPage() {
  const { wid = "" } = useParams<{ wid: string }>()
  const [activeTab, setActiveTab] = useState<"all" | "connected" | "available">("all")
  const [busyId, setBusyId] = useState<string | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)
  const [connectProvider, setConnectProvider] = useState<(typeof items)[number] | null>(null)
  const [disconnectProvider, setDisconnectProvider] = useState<(typeof items)[number] | null>(null)
  const { items, isLoading, error, connect, activate, deactivate, removeIntegration } = useIntegrations(wid)

  async function handleConnect(provider: (typeof items)[number], config: Record<string, unknown>) {
    await connect(provider, config)
  }

  async function runProviderAction(
    provider: (typeof items)[number],
    action: (provider: (typeof items)[number]) => Promise<void>,
    options?: { surfaceErrors?: boolean },
  ) {
    setBusyId(provider.id)
    if (options?.surfaceErrors) setActionError(null)
    try {
      await action(provider)
    } catch (err) {
      if (options?.surfaceErrors) {
        setActionError(err instanceof Error ? err.message : "Failed to update provider")
        return
      }
      throw err
    } finally {
      setBusyId(null)
    }
  }

  if (isLoading) {
    return (
      <div className="flex min-h-[320px] items-center justify-center gap-3">
        <span className="text-sm text-text-secondary">Loading integrations...</span>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="Integrations"
        description="Connect and manage your messaging service providers"
      />

      {error ? (
        <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </p>
      ) : null}

      {actionError ? (
        <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
          {actionError}
        </p>
      ) : null}

      <Tabs.Root value={activeTab} onValueChange={(v) => setActiveTab(v as typeof activeTab)}>
        <Tabs.List className="mb-6 flex gap-2 border-b border-border">
          {(["all", "connected", "available"] as const).map((tab) => (
            <Tabs.Trigger
              key={tab}
              value={tab}
              className={cn(
                "cursor-pointer border-b-2 bg-transparent px-4 py-3 text-sm font-medium capitalize transition-colors",
                "border-b-transparent text-text-secondary data-[state=active]:border-b-primary-brand data-[state=active]:text-primary-brand",
              )}
            >
              {tab}
            </Tabs.Trigger>
          ))}
        </Tabs.List>

        {(["all", "connected", "available"] as const).map((tab) => {
          const tabItems = filterIntegrationsByTab(items, tab)
          const tabGrouped = groupByCategory(tabItems)

          return (
            <Tabs.Content key={tab} value={tab} className="flex flex-col gap-8">
              {tabItems.length === 0 ? (
                <div className="flex flex-col items-center justify-center gap-6 rounded-xl border border-dashed border-border bg-card/50 py-16 px-4 text-center">
                  <div className="rounded-full bg-secondary/50 p-4 text-text-tertiary">
                    <Icon name="link_off" className="size-8" />
                  </div>
                  <div className="flex flex-col gap-2 max-w-sm">
                    <h3 className="text-lg font-semibold text-foreground">
                      {tab === "connected"
                        ? "No Connected Integrations"
                        : tab === "available"
                        ? "No Available Integrations"
                        : "No Integrations Found"}
                    </h3>
                    <p className="text-sm text-text-secondary">
                      {tab === "connected"
                        ? "You haven't connected any integrations to this workspace yet."
                        : tab === "available"
                        ? "All available integrations are already connected to your workspace."
                        : "There are no integration providers matching your filter."}
                    </p>
                  </div>
                  {tab === "connected" && (
                    <Button
                      onClick={() => setActiveTab("available")}
                      className="bg-primary-brand hover:bg-primary-brand-hover"
                    >
                      Connect a Provider
                    </Button>
                  )}
                </div>
              ) : (
                (Object.keys(tabGrouped) as IntegrationChannel[]).map((category) => {
                  const categoryItems = tabGrouped[category]
                  if (categoryItems.length === 0) return null

                  return (
                    <section key={category}>
                      <h2 className="mb-3 text-base font-semibold text-foreground">
                        {CATEGORY_LABELS[category]}
                      </h2>
                      <div className="overflow-hidden rounded-lg border border-border bg-card">
                        {categoryItems.map((provider) => (
                          <IntegrationRow
                            key={provider.id}
                            provider={provider}
                            isBusy={busyId === provider.id}
                            onConnect={setConnectProvider}
                            onActivate={(provider) =>
                              runProviderAction(provider, activate, { surfaceErrors: true })
                            }
                            onDisconnect={setDisconnectProvider}
                          />
                        ))}
                      </div>
                    </section>
                  )
                })
              )}
            </Tabs.Content>
          )
        })}
      </Tabs.Root>

      <ConnectModal
        isOpen={connectProvider !== null}
        onClose={() => setConnectProvider(null)}
        workspaceId={wid}
        provider={connectProvider}
        onConnect={handleConnect}
      />

      <DisconnectModal
        isOpen={disconnectProvider !== null}
        onClose={() => setDisconnectProvider(null)}
        provider={disconnectProvider}
        onDeactivate={(provider) => runProviderAction(provider, deactivate)}
        onRemove={(provider) => runProviderAction(provider, removeIntegration)}
      />
    </div>
  )
}
