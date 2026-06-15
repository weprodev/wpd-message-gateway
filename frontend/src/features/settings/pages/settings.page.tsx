import * as Tabs from "@radix-ui/react-tabs"
import { useState } from "react"
import { useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"
import { PageHeader } from "@/shared/components/page-header"
import { cn } from "@/lib/utils"

import { ApiKeyRow } from "../components/api-key-row"
import { ApiKeyCreateModal } from "../components/api-key-create-modal"
import { ApiKeyDeleteModal } from "../components/api-key-delete-modal"
import { ApiKeyRegenerateModal } from "../components/api-key-regenerate-modal"
import { RadioOption } from "../components/radio-option"
import { useSettings } from "../hooks/use-settings.hook"
import type { ApiKeyCredentials, RetentionMode, SettingsTab, WorkspaceSettings } from "../settings.types"

interface GeneralSettingsPanelProps {
  settings: WorkspaceSettings
  onSave: (patch: Record<string, string>) => Promise<void>
}

function GeneralSettingsPanel({ settings, onSave }: GeneralSettingsPanelProps) {
  const [ownerEmail, setOwnerEmail] = useState(settings.owner_email ?? "")
  const [pinCode, setPinCode] = useState(settings.pin_code ?? "")
  const [showPin, setShowPin] = useState(false)
  const [isSaving, setIsSaving] = useState(false)

  async function handleSave() {
    setIsSaving(true)
    try {
      await onSave({ owner_email: ownerEmail, pin_code: pinCode })
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <div className="flex max-w-xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <label htmlFor="owner-email" className="text-sm font-medium text-text-secondary">
          Owner email
        </label>
        <Input
          id="owner-email"
          type="email"
          value={ownerEmail}
          onChange={(ev) => setOwnerEmail(ev.target.value)}
          placeholder="owner@company.com"
        />
      </div>

      <div className="flex flex-col gap-2">
        <label htmlFor="pin-code" className="text-sm font-medium text-text-secondary">
          Workspace PIN
        </label>
        <div className="relative">
          <Input
            id="pin-code"
            type={showPin ? "text" : "password"}
            value={pinCode}
            onChange={(ev) => setPinCode(ev.target.value)}
            placeholder="••••••"
            className="pr-10"
          />
          <button
            type="button"
            onClick={() => setShowPin((v) => !v)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-text-placeholder"
            aria-label={showPin ? "Hide PIN" : "Show PIN"}
          >
            <Icon name={showPin ? "visibility_off" : "visibility"} size="sm" />
          </button>
        </div>
      </div>

      <Button type="button" onClick={handleSave} disabled={isSaving} className="w-fit">
        {isSaving ? "Saving…" : "Save changes"}
      </Button>
    </div>
  )
}

interface RetentionSettingsPanelProps {
  initialMode: RetentionMode
  onSave: (patch: Record<string, string>) => Promise<void>
}

function RetentionSettingsPanel({ initialMode, onSave }: RetentionSettingsPanelProps) {
  const [retentionMode, setRetentionMode] = useState<RetentionMode>(initialMode)
  const [isSaving, setIsSaving] = useState(false)

  async function handleSave() {
    setIsSaving(true)
    try {
      await onSave({ data_retention: retentionMode })
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <div className="flex max-w-xl flex-col gap-4">
      <RadioOption
        id="retention-memory"
        name="retention"
        label="Memory only"
        description="Store messages in memory for testing — no persistence."
        checked={retentionMode === "memory"}
        onChange={() => setRetentionMode("memory")}
      />
      <RadioOption
        id="retention-both"
        name="retention"
        label="Memory + Database"
        description="Persist messages in the portal inbox and database."
        checked={retentionMode === "both"}
        onChange={() => setRetentionMode("both")}
      />
      <RadioOption
        id="retention-providers"
        name="retention"
        label="Providers only"
        description="Send through providers without storing message content."
        checked={retentionMode === "providers"}
        onChange={() => setRetentionMode("providers")}
      />

      <Button type="button" onClick={handleSave} disabled={isSaving} className="mt-2 w-fit">
        {isSaving ? "Saving…" : "Save retention policy"}
      </Button>
    </div>
  )
}

export function SettingsPage() {
  const { wid = "" } = useParams<{ wid: string }>()
  const [activeTab, setActiveTab] = useState<SettingsTab>("general")
  const [busyKeyId, setBusyKeyId] = useState<string | null>(null)
  const [pendingDeleteKeyId, setPendingDeleteKeyId] = useState<string | null>(null)
  const [pendingRegenerateKeyId, setPendingRegenerateKeyId] = useState<string | null>(null)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)

  const { settings, apiKeys, retentionMode, isLoading, error, saveSettings, addApiKey, removeApiKey, rotateApiKey } =
    useSettings(wid)

  async function handleCreateApiKey(name: string): Promise<ApiKeyCredentials | null> {
    const created = await addApiKey(name)
    if (!created.client_secret) return null
    return {
      clientId: created.client_id,
      clientSecret: created.client_secret,
      keyName: created.name,
      mode: "created",
    }
  }

  function handleRequestRegenerate(keyId: string) {
    setPendingRegenerateKeyId(keyId)
  }

  async function handleRegenerateApiKey(): Promise<ApiKeyCredentials | null> {
    if (!pendingRegenerateKeyId) return null
    const key = apiKeys.find((item) => item.id === pendingRegenerateKeyId)
    const { client_secret: clientSecret } = await rotateApiKey(pendingRegenerateKeyId)
    if (!key || !clientSecret) return null
    return {
      clientId: key.client_id,
      clientSecret,
      keyName: key.name,
      mode: "regenerated",
    }
  }

  function handleRequestDelete(keyId: string) {
    setPendingDeleteKeyId(keyId)
  }

  function handleCancelDelete() {
    setPendingDeleteKeyId(null)
  }

  async function handleConfirmDelete() {
    if (!pendingDeleteKeyId) return
    setBusyKeyId(pendingDeleteKeyId)
    try {
      await removeApiKey(pendingDeleteKeyId)
      setPendingDeleteKeyId(null)
    } finally {
      setBusyKeyId(null)
    }
  }

  if (isLoading) {
    return (
      <div className="flex min-h-[320px] items-center justify-center gap-3">
        <Spinner size="lg" />
        <span className="text-sm text-text-secondary">Loading settings...</span>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="Settings"
        description="Manage your workspace settings and preferences."
      />

      {error ? (
        <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </p>
      ) : null}

      <Tabs.Root value={activeTab} onValueChange={(v) => setActiveTab(v as SettingsTab)}>
        <Tabs.List className="mb-6 flex flex-wrap gap-2 border-b border-border pb-2">
          {(
            [
              ["general", "General"],
              ["developer", "Developer"],
              ["team", "Team Management"],
              ["retention", "Data Retention"],
            ] as const
          ).map(([value, label]) => (
            <Tabs.Trigger
              key={value}
              value={value}
              className={cn(
                "cursor-pointer rounded-lg border px-4 py-2 text-sm font-medium transition-colors",
                "border-border bg-transparent text-text-secondary",
                "data-[state=active]:border-primary-brand data-[state=active]:bg-card data-[state=active]:text-primary-brand",
              )}
            >
              {label}
            </Tabs.Trigger>
          ))}
        </Tabs.List>

        <Tabs.Content value="general">
          <GeneralSettingsPanel key={wid} settings={settings} onSave={saveSettings} />
        </Tabs.Content>

        <Tabs.Content value="developer" className="flex flex-col gap-4">
          <div className="flex items-center justify-between gap-4">
            <div>
              <h2 className="text-base font-semibold text-foreground">API Keys</h2>
              <p className="text-sm text-text-secondary">
                Manage keys used to authenticate gateway requests.
              </p>
            </div>
            <Button type="button" onClick={() => setIsCreateModalOpen(true)}>
              <Icon name="add" size="sm" />
              Create key
            </Button>
          </div>

          <div className="overflow-hidden rounded-lg border border-border bg-card">
            {apiKeys.length === 0 ? (
              <p className="px-4 py-8 text-center text-sm text-text-secondary">
                No API keys yet. Create one to send test requests.
              </p>
            ) : (
              apiKeys.map((key) => (
                <ApiKeyRow
                  key={key.id}
                  apiKey={key}
                  isBusy={busyKeyId === key.id}
                  onRegenerate={handleRequestRegenerate}
                  onDelete={handleRequestDelete}
                />
              ))
            )}
          </div>
        </Tabs.Content>

        <Tabs.Content value="team">
          <div className="rounded-lg border border-border bg-card p-8 text-center">
            <Icon name="groups" size="lg" className="mx-auto text-text-tertiary" />
            <h2 className="mt-4 text-base font-semibold text-foreground">Team Management</h2>
            <p className="mt-2 text-sm text-text-secondary">Coming soon — invite teammates and manage roles.</p>
          </div>
        </Tabs.Content>

        <Tabs.Content value="retention">
          <RetentionSettingsPanel key={`${wid}-${retentionMode}`} initialMode={retentionMode} onSave={saveSettings} />
        </Tabs.Content>
      </Tabs.Root>

      <ApiKeyCreateModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCreate={handleCreateApiKey}
      />

      <ApiKeyDeleteModal
        isOpen={pendingDeleteKeyId !== null}
        onCancel={handleCancelDelete}
        onDelete={handleConfirmDelete}
        isDeleting={pendingDeleteKeyId !== null && busyKeyId === pendingDeleteKeyId}
      />

      <ApiKeyRegenerateModal
        isOpen={pendingRegenerateKeyId !== null}
        onClose={() => setPendingRegenerateKeyId(null)}
        onRegenerate={handleRegenerateApiKey}
      />

      <section className="mt-4 rounded-lg border border-destructive/30 bg-destructive/5 p-6">
        <h2 className="text-base font-semibold text-destructive">Danger zone</h2>
        <p className="mt-1 text-sm text-text-secondary">
          Permanently delete this workspace and all associated data.
        </p>
        <Button type="button" variant="destructive" className="mt-4" disabled>
          Delete workspace
        </Button>
      </section>
    </div>
  )
}
