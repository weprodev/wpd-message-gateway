import * as Tabs from "@radix-ui/react-tabs"
import { useParams } from "react-router-dom"
import { useState } from "react"

import { PageHeader } from "@/shared/components/page-header"
import { cn } from "@/lib/utils"

import { IntegrationRow } from "../components/integration-row"
import { ConnectModal } from "../components/connect-modal"
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
  const [connectProvider, setConnectProvider] = useState<(typeof items)[number] | null>(null)
  const { items, isLoading, error, connect, disconnect } = useIntegrations(wid)

  const filtered = filterIntegrationsByTab(items, activeTab)
  const grouped = groupByCategory(filtered)

  async function handleConnect(provider: (typeof items)[number], config: Record<string, unknown>) {
    await connect(provider, config)
  }

  async function handleDisconnect(provider: (typeof items)[number]) {
    setBusyId(provider.id)
    try {
      await disconnect(provider)
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

        {(["all", "connected", "available"] as const).map((tab) => (
          <Tabs.Content key={tab} value={tab} className="flex flex-col gap-8">
            {(Object.keys(grouped) as IntegrationChannel[]).map((category) => {
              const categoryItems = grouped[category]
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
                        onDisconnect={handleDisconnect}
                      />
                    ))}
                  </div>
                </section>
              )
            })}
          </Tabs.Content>
        ))}
      </Tabs.Root>

      <ConnectModal
        isOpen={connectProvider !== null}
        onClose={() => setConnectProvider(null)}
        workspaceId={wid}
        provider={connectProvider}
        onConnect={handleConnect}
      />
    </div>
  )
}
