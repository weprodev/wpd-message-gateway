import * as Tabs from "@radix-ui/react-tabs"
import { useState, type FormEvent } from "react"
import { useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { Modal, ModalActions, useModalSession } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"
import { PageHeader } from "@/shared/components/page-header"
import { cn } from "@/lib/utils"
import { toUserMessage } from "@/lib/errors"

import { ApiKeyCredentialsView } from "../components/api-key-credentials-view"
import { ApiKeyRow } from "../components/api-key-row"
import { RadioOption } from "../components/radio-option"
import { useSettings } from "../hooks/use-settings.hook"
import type { ApiKeyCredentials, RetentionMode, SettingsTab, WorkspaceSettings } from "../settings.types"

const CREATE_FAILED_MESSAGE = "Could not create the API key. Please try again."
const REGENERATE_FAILED_MESSAGE = "Could not regenerate the API key. Please try again."

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
        id="retention-memory-database"
        name="retention"
        label="Memory + Database"
        description="Persist messages in the portal inbox and database."
        checked={retentionMode === "memory_database"}
        onChange={() => setRetentionMode("memory_database")}
      />
      <RadioOption
        id="retention-providers"
        name="retention"
        label="Providers only"
        description="Send through providers without storing message content."
        checked={retentionMode === "providers"}
        onChange={() => setRetentionMode("providers")}
      />
      <RadioOption
        id="retention-provider-database"
        name="retention"
        label="Provider + Database"
        description="Send through providers and persist the full payload in PostgreSQL."
        checked={retentionMode === "provider_database"}
        onChange={() => setRetentionMode("provider_database")}
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

  const [createStep, setCreateStep] = useState<"form" | "credentials">("form")
  const [createName, setCreateName] = useState("")
  const [createError, setCreateError] = useState<string | null>(null)
  const [isCreating, setIsCreating] = useState(false)
  const [createdCredentials, setCreatedCredentials] = useState<ApiKeyCredentials | null>(null)

  const [regenerateStep, setRegenerateStep] = useState<"confirm" | "credentials">("confirm")
  const [regenerateError, setRegenerateError] = useState<string | null>(null)
  const [isRegenerating, setIsRegenerating] = useState(false)
  const [regeneratedCredentials, setRegeneratedCredentials] = useState<ApiKeyCredentials | null>(null)

  const { settings, apiKeys, retentionMode, isLoading, error, saveSettings, addApiKey, removeApiKey, rotateApiKey } =
    useSettings(wid)

  function resetCreateModal() {
    setCreateStep("form")
    setCreateName("")
    setCreateError(null)
    setIsCreating(false)
    setCreatedCredentials(null)
  }

  function resetRegenerateModal() {
    setRegenerateStep("confirm")
    setRegenerateError(null)
    setIsRegenerating(false)
    setRegeneratedCredentials(null)
  }

  const createSessionRef = useModalSession(isCreateModalOpen, resetCreateModal)
  const regenerateSessionRef = useModalSession(pendingRegenerateKeyId !== null, resetRegenerateModal)

  function closeCreateModal() {
    resetCreateModal()
    setIsCreateModalOpen(false)
  }

  async function handleCreateSubmit(event: FormEvent) {
    event.preventDefault()
    const trimmed = createName.trim()
    if (!trimmed) {
      setCreateError("API key name is required")
      return
    }

    const session = createSessionRef.current
    setCreateError(null)
    setIsCreating(true)
    try {
      const created = await addApiKey(trimmed)
      if (session !== createSessionRef.current) return
      setCreatedCredentials({
        clientId: created.client_id,
        clientSecret: created.client_secret,
        keyName: created.name,
        mode: "created",
      })
      setCreateStep("credentials")
    } catch (err) {
      if (session !== createSessionRef.current) return
      setCreateError(toUserMessage(err, CREATE_FAILED_MESSAGE))
    } finally {
      if (session === createSessionRef.current) {
        setIsCreating(false)
      }
    }
  }

  function handleRequestRegenerate(keyId: string) {
    setPendingRegenerateKeyId(keyId)
  }

  async function handleRegenerateConfirm() {
    if (!pendingRegenerateKeyId) {
      setRegenerateError("No API key selected for regeneration.")
      return
    }

    const key = apiKeys.find((item) => item.id === pendingRegenerateKeyId)
    if (!key) {
      setRegenerateError("This API key could not be found. Refresh the page and try again.")
      return
    }

    const session = regenerateSessionRef.current
    setRegenerateError(null)
    setIsRegenerating(true)
    try {
      const { client_secret: clientSecret } = await rotateApiKey(pendingRegenerateKeyId)
      if (session !== regenerateSessionRef.current) return
      setRegeneratedCredentials({
        clientId: key.client_id,
        clientSecret,
        keyName: key.name,
        mode: "regenerated",
      })
      setRegenerateStep("credentials")
    } catch (err) {
      if (session !== regenerateSessionRef.current) return
      setRegenerateError(toUserMessage(err, REGENERATE_FAILED_MESSAGE))
    } finally {
      if (session === regenerateSessionRef.current) {
        setIsRegenerating(false)
      }
    }
  }

  function closeRegenerateModal() {
    resetRegenerateModal()
    setPendingRegenerateKeyId(null)
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

      <Modal
        isOpen={isCreateModalOpen}
        onClose={closeCreateModal}
        title={createStep === "form" ? "Create API key" : "API key created"}
        preventDismiss={createStep === "credentials" || isCreating}
      >
        {createStep === "form" ? (
          <form onSubmit={handleCreateSubmit} className="flex flex-col gap-6">
            <div className="flex flex-col gap-1.5">
              <label htmlFor="api-key-name" className="text-xs font-semibold uppercase text-text-secondary">
                Name
              </label>
              <Input
                id="api-key-name"
                type="text"
                value={createName}
                onChange={(event) => {
                  setCreateName(event.target.value)
                  if (createError) setCreateError(null)
                }}
                placeholder="e.g. Production"
                disabled={isCreating}
                autoFocus
              />
              {createError ? <p className="text-sm text-destructive">{createError}</p> : null}
            </div>

            <ModalActions
              confirmLabel="Create"
              confirmType="submit"
              onCancel={closeCreateModal}
              isLoading={isCreating}
              loadingLabel="Creating…"
              cancelDisabled={isCreating}
              confirmClassName="min-w-[7.25rem]"
            />
          </form>
        ) : createdCredentials ? (
          <ApiKeyCredentialsView credentials={createdCredentials} onConfirm={closeCreateModal} />
        ) : null}
      </Modal>

      <Modal
        isOpen={pendingDeleteKeyId !== null}
        onClose={handleCancelDelete}
        title="Delete API key"
        description="Are you sure you want to delete this API key?"
        preventDismiss={pendingDeleteKeyId !== null && busyKeyId === pendingDeleteKeyId}
      >
        <ModalActions
          confirmLabel="Delete"
          confirmVariant="destructive"
          onCancel={handleCancelDelete}
          onConfirm={handleConfirmDelete}
          isLoading={pendingDeleteKeyId !== null && busyKeyId === pendingDeleteKeyId}
          loadingLabel="Deleting…"
          confirmClassName="min-w-[7.25rem]"
        />
      </Modal>

      <Modal
        isOpen={pendingRegenerateKeyId !== null}
        onClose={closeRegenerateModal}
        title={regenerateStep === "confirm" ? "Regenerate API key" : "API key regenerated"}
        description={
          regenerateStep === "confirm" ? "Are you sure you want to regenerate this API key?" : undefined
        }
        preventDismiss={regenerateStep === "credentials" || isRegenerating}
      >
        {regenerateStep === "confirm" ? (
          <div className="flex flex-col gap-6">
            {regenerateError ? <p className="text-sm text-destructive">{regenerateError}</p> : null}
            <ModalActions
              confirmLabel="Regenerate"
              onCancel={closeRegenerateModal}
              onConfirm={handleRegenerateConfirm}
              isLoading={isRegenerating}
              loadingLabel="Regenerating…"
              confirmClassName="min-w-[9.75rem]"
            />
          </div>
        ) : regeneratedCredentials ? (
          <ApiKeyCredentialsView credentials={regeneratedCredentials} onConfirm={closeRegenerateModal} />
        ) : null}
      </Modal>

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
