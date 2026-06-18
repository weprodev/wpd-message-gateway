import type React from "react"
import * as Tabs from "@radix-ui/react-tabs"
import { useParams } from "react-router-dom"
import { useEffect, useState } from "react"

import { PageHeader } from "@/shared/components/page-header"
import { Button } from "@/components/ui/button"
import { DialogDescription, DialogTitle } from "@/components/ui/dialog"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"
import { cn } from "@/lib/utils"

import { IntegrationProviderIcon } from "../components/integration-provider-icon"
import { IntegrationRow } from "../components/integration-row"
import {
  filterIntegrationsByTab,
  groupByCategory,
  useIntegrations,
  type IntegrationViewModel,
} from "../hooks/use-integrations.hook"
import { fetchProviderConfigFields, type ProviderConfigField } from "../integrations.api"
import type { IntegrationChannel, IntegrationActionResult } from "../integrations.types"

const CATEGORY_LABELS: Record<IntegrationChannel, string> = {
  email: "Email",
  sms: "SMS",
  push: "Push",
  chat: "Chat",
}

type DisconnectSubmittingAction = "deactivate" | "remove" | null

export function IntegrationsPage() {
  const { wid = "" } = useParams<{ wid: string }>()
  const [activeTab, setActiveTab] = useState<"all" | "connected" | "available">("all")
  const [busyId, setBusyId] = useState<string | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)
  const [connectProvider, setConnectProvider] = useState<IntegrationViewModel | null>(null)
  const [disconnectProvider, setDisconnectProvider] = useState<IntegrationViewModel | null>(null)
  const [connectFields, setConnectFields] = useState<ProviderConfigField[]>([])
  const [isConnectFieldsLoading, setIsConnectFieldsLoading] = useState(false)
  const [connectError, setConnectError] = useState<string | null>(null)
  const [connectFormData, setConnectFormData] = useState<Record<string, string>>({})
  const [isConnectSubmitting, setIsConnectSubmitting] = useState(false)
  const [disconnectSubmittingAction, setDisconnectSubmittingAction] = useState<DisconnectSubmittingAction>(null)
  const [disconnectError, setDisconnectError] = useState<string | null>(null)
  const { items, isLoading, error, connect, activate, deactivate, removeIntegration } = useIntegrations(wid)

  const isDisconnectSubmitting = disconnectSubmittingAction !== null

  useEffect(() => {
    if (!connectProvider) return

    const loadFields = async () => {
      setIsConnectFieldsLoading(true)
      setConnectError(null)
      try {
        const result = await fetchProviderConfigFields(wid, connectProvider.id)
        setConnectFields(result)

        const defaults: Record<string, string> = {}
        result.forEach((field) => {
          defaults[field.key] = field.default_value || ""
        })
        setConnectFormData(defaults)
      } catch (err) {
        setConnectError(err instanceof Error ? err.message : "Failed to load configuration fields")
        setConnectFields([])
      } finally {
        setIsConnectFieldsLoading(false)
      }
    }

    void loadFields()
  }, [connectProvider, wid])

  function closeConnectDialog() {
    setConnectProvider(null)
    setConnectFields([])
    setConnectFormData({})
    setConnectError(null)
    setIsConnectSubmitting(false)
  }

  function closeDisconnectDialog() {
    if (isDisconnectSubmitting) return
    setDisconnectProvider(null)
    setDisconnectError(null)
    setDisconnectSubmittingAction(null)
  }

  function handleConnectInputChange(key: string, value: string) {
    setConnectFormData((prev) => ({ ...prev, [key]: value }))
  }

  async function handleConnectSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!connectProvider) return

    setIsConnectSubmitting(true)
    setConnectError(null)
    try {
      const result = await connect(connectProvider, connectFormData)
      if (!result.ok) {
        setConnectError(result.message ?? "Failed to connect provider")
      } else {
        closeConnectDialog()
      }
    } catch (err) {
      setConnectError(err instanceof Error ? err.message : "Failed to connect provider")
    } finally {
      setIsConnectSubmitting(false)
    }
  }

  async function runDisconnectAction(
    action: DisconnectSubmittingAction,
    handler: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>,
  ) {
    if (!disconnectProvider || !action) return

    setDisconnectSubmittingAction(action)
    setDisconnectError(null)
    try {
      const result = await handler(disconnectProvider)
      if (!result.ok) {
        setDisconnectError(result.message ?? "Failed to update provider")
      } else {
        setDisconnectProvider(null)
        setDisconnectError(null)
      }
    } catch (err) {
      setDisconnectError(err instanceof Error ? err.message : "Failed to update provider")
    } finally {
      setDisconnectSubmittingAction(null)
    }
  }

  async function runProviderAction(
    provider: IntegrationViewModel,
    action: (provider: IntegrationViewModel) => Promise<IntegrationActionResult>,
    options?: { surfaceErrors?: boolean },
  ): Promise<IntegrationActionResult> {
    setBusyId(provider.id)
    if (options?.surfaceErrors) setActionError(null)
    try {
      const result = await action(provider)
      if (!result.ok && options?.surfaceErrors) {
        setActionError(result.message ?? "Failed to update provider")
      }
      return result
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to update provider"
      if (options?.surfaceErrors) {
        setActionError(message)
        return { ok: false, message }
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

      <Modal
        isOpen={connectProvider !== null}
        onClose={closeConnectDialog}
        title={connectProvider ? `Connect ${connectProvider.name}` : undefined}
      >
        {connectProvider && isConnectFieldsLoading ? (
          <div className="flex items-center justify-center gap-3 py-8">
            <Spinner />
            <span className="text-sm text-text-secondary">Loading configuration...</span>
          </div>
        ) : null}

        {connectProvider && connectError && connectFields.length === 0 && !isConnectFieldsLoading ? (
          <div className="flex flex-col gap-4 py-4">
            <p className="text-sm text-destructive">{connectError}</p>
            <Button onClick={closeConnectDialog} variant="outline">
              Close
            </Button>
          </div>
        ) : null}

        {connectProvider && connectFields.length > 0 && !isConnectFieldsLoading ? (
          <form onSubmit={handleConnectSubmit} className="flex flex-col gap-6">
            {connectError ? (
              <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
                {connectError}
              </p>
            ) : null}

            <div className="flex max-h-[350px] flex-col gap-4 overflow-y-auto pr-1">
              {connectFields.map((field) => (
                <div key={field.key} className="flex flex-col gap-1.5">
                  <label htmlFor={field.key} className="text-sm font-medium text-foreground">
                    {field.label}
                    {field.required ? <span className="ml-1 text-destructive">*</span> : null}
                  </label>
                  <Input
                    id={field.key}
                    type={field.field_type === "password" ? "password" : field.field_type === "email" ? "email" : "text"}
                    required={field.required}
                    value={connectFormData[field.key] || ""}
                    onChange={(e) => handleConnectInputChange(field.key, e.target.value)}
                    placeholder={field.description || `Enter ${field.label.toLowerCase()}`}
                    className="bg-input"
                  />
                  {field.description ? (
                    <span className="text-[12px] leading-normal text-text-secondary">
                      {field.description}
                    </span>
                  ) : null}
                </div>
              ))}
            </div>

            <div className="flex items-center justify-end gap-3 border-t border-border pt-2">
              <Button type="button" variant="outline" onClick={closeConnectDialog} disabled={isConnectSubmitting}>
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={isConnectSubmitting}
                className="bg-primary-brand hover:bg-primary-brand-hover"
              >
                {isConnectSubmitting ? (
                  <div className="flex items-center gap-2">
                    <Spinner size="sm" />
                    <span>Connecting...</span>
                  </div>
                ) : (
                  "Connect"
                )}
              </Button>
            </div>
          </form>
        ) : null}
      </Modal>

      <Modal
        isOpen={disconnectProvider !== null}
        onClose={closeDisconnectDialog}
        preventClose={isDisconnectSubmitting}
        header={
          disconnectProvider ? (
            <div className="flex items-start gap-4">
              <IntegrationProviderIcon
                icon={disconnectProvider.icon}
                name={disconnectProvider.name}
                className="size-12 p-2.5 text-2xl"
              />
              <div className="min-w-0 flex-1">
                <DialogTitle className="text-xl font-semibold text-foreground">
                  Disconnect {disconnectProvider.name}
                </DialogTitle>
                <DialogDescription className="mt-1 text-sm text-text-secondary">
                  Are you sure you want to disconnect this provider?
                </DialogDescription>
              </div>
            </div>
          ) : undefined
        }
      >
        {disconnectProvider ? (
          <div className="flex flex-col gap-6">
            {disconnectError ? (
              <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
                {disconnectError}
              </p>
            ) : null}

            <div className="flex flex-col gap-3">
              <div className="flex items-start gap-3 rounded-lg border border-border bg-muted/40 p-4">
                <div className="flex size-10 shrink-0 items-center justify-center rounded-full border border-border bg-card">
                  <Icon name="pause" size="sm" className="text-text-secondary" />
                </div>
                <div className="min-w-0 flex-1">
                  <h3 className="text-sm font-semibold text-foreground">Deactivate Connection</h3>
                  <p className="mt-1 text-[13px] leading-normal text-text-secondary">
                    Stop all message delivery immediately. Your API keys and configurations will be preserved for later use.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3 rounded-lg border border-destructive/20 bg-destructive/5 p-4">
                <div className="flex size-10 shrink-0 items-center justify-center rounded-full border border-destructive/20 bg-card">
                  <Icon name="delete" size="sm" className="text-destructive" />
                </div>
                <div className="min-w-0 flex-1">
                  <h3 className="text-sm font-semibold text-foreground">Remove Integration</h3>
                  <p className="mt-1 text-[13px] leading-normal text-text-secondary">
                    Permanently delete all API keys and configurations associated with this provider. This action cannot be undone.
                  </p>
                </div>
              </div>
            </div>

            <div className="flex flex-wrap items-center justify-end gap-3 border-t border-border pt-4">
              <Button
                type="button"
                variant="outline"
                onClick={closeDisconnectDialog}
                disabled={isDisconnectSubmitting}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => runDisconnectAction("deactivate", (provider) => runProviderAction(provider, deactivate))}
                disabled={isDisconnectSubmitting}
                className="border-primary-brand text-primary-brand hover:bg-primary-brand/5"
              >
                {disconnectSubmittingAction === "deactivate" ? (
                  <span className="flex items-center gap-2">
                    <Spinner size="sm" />
                    Deactivating...
                  </span>
                ) : (
                  "Deactivate"
                )}
              </Button>
              <Button
                type="button"
                variant="destructive"
                onClick={() => runDisconnectAction("remove", (provider) => runProviderAction(provider, removeIntegration))}
                disabled={isDisconnectSubmitting}
              >
                {disconnectSubmittingAction === "remove" ? (
                  <span className="flex items-center gap-2">
                    <Spinner size="sm" />
                    Removing...
                  </span>
                ) : (
                  "Remove Integration"
                )}
              </Button>
            </div>
          </div>
        ) : null}
      </Modal>
    </div>
  )
}
