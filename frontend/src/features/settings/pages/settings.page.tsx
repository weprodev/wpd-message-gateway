import * as Tabs from "@radix-ui/react-tabs"
import { useState } from "react"
import { useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"
import { PageHeader } from "@/shared/components/page-header"
import { cn } from "@/lib/utils"

import { ApiKeyRow } from "../components/api-key-row"
import { RadioOption } from "../components/radio-option"
import { useSettings } from "../hooks/use-settings.hook"
import { dispatchConfigsEqual, toStoreMessageContentSetting } from "../message-dispatch-mode"
import type {
  MessageDispatchConfig,
  MessageDispatchMode,
  SettingsTab,
  WorkspaceSettings,
  WorkspaceSettingsPatch,
} from "../settings.types"

interface GeneralSettingsPanelProps {
  settings: WorkspaceSettings
  onSave: (patch: WorkspaceSettingsPatch) => Promise<void>
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

interface DispatchSettingsPanelProps {
  initialConfig: MessageDispatchConfig
  onSave: (patch: WorkspaceSettingsPatch) => Promise<void>
}

function DispatchSettingsPanel({ initialConfig, onSave }: DispatchSettingsPanelProps) {
  const [mode, setMode] = useState<MessageDispatchMode>(initialConfig.mode)
  const [storeMessageContent, setStoreMessageContent] = useState(initialConfig.storeMessageContent)
  const [isSaving, setIsSaving] = useState(false)

  const currentConfig: MessageDispatchConfig = { mode, storeMessageContent }
  const isDirty = !dispatchConfigsEqual(currentConfig, initialConfig)

  async function handleSave() {
    setIsSaving(true)
    try {
      await onSave({
        message_dispatch_mode: mode,
        store_message_content: toStoreMessageContentSetting(storeMessageContent),
      })
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <div className="flex max-w-xl flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h2 className="text-sm font-medium text-foreground">Dispatch mode</h2>
        <p className="text-sm text-text-secondary">Choose where outbound messages are routed.</p>
        <div className="mt-1 flex flex-col gap-3">
          <RadioOption
            id="dispatch-memory"
            name="dispatch-mode"
            label="Memory"
            description="Capture messages in memory for development and testing."
            checked={mode === "memory"}
            onChange={() => setMode("memory")}
          />
          <RadioOption
            id="dispatch-provider"
            name="dispatch-mode"
            label="Provider"
            description="Send messages through the connected channel integration."
            checked={mode === "provider"}
            onChange={() => setMode("provider")}
          />
        </div>
      </div>

      <div className="flex items-center justify-between gap-4 rounded-lg border border-border p-4">
        <label htmlFor="store-message-content" className={cn("flex flex-1 cursor-pointer flex-col gap-1", mode === "memory" && "opacity-50 cursor-not-allowed")}>
          <span className="text-sm font-medium text-foreground">Store message content in inbox</span>
          <span className="text-sm text-text-secondary">
            {mode === "memory" 
              ? "Message bodies are always captured in Memory mode." 
              : "When enabled, message bodies are captured in the portal inbox for review."}
          </span>
        </label>
        <Checkbox
          id="store-message-content"
          checked={mode === "memory" ? true : storeMessageContent}
          onCheckedChange={(checked) => setStoreMessageContent(checked === true)}
          disabled={mode === "memory"}
        />
      </div>

      <Button type="button" onClick={handleSave} disabled={isSaving || !isDirty} className="w-fit">
        {isSaving ? "Saving…" : "Save dispatch settings"}
      </Button>
    </div>
  )
}

export function SettingsPage() {
  const { wid = "" } = useParams<{ wid: string }>()
  const [activeTab, setActiveTab] = useState<SettingsTab>("general")
  const [busyKeyId, setBusyKeyId] = useState<string | null>(null)

  const { settings, apiKeys, messageDispatchConfig, isLoading, error, saveSettings, addApiKey, removeApiKey, rotateApiKey } =
    useSettings(wid)

  async function handleCreateApiKey() {
    const name = window.prompt("API key name")
    if (!name?.trim()) return
    await addApiKey(name.trim())
  }

  async function handleRegenerate(keyId: string) {
    setBusyKeyId(keyId)
    try {
      await rotateApiKey(keyId)
    } finally {
      setBusyKeyId(null)
    }
  }

  async function handleDeleteKey(keyId: string) {
    if (!window.confirm("Delete this API key?")) return
    setBusyKeyId(keyId)
    try {
      await removeApiKey(keyId)
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
              ["dispatch", "Message Dispatch"],
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
            <Button type="button" onClick={handleCreateApiKey}>
              <Icon name="add" size="sm" />
              Generate key
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
                  onRegenerate={handleRegenerate}
                  onDelete={handleDeleteKey}
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

        <Tabs.Content value="dispatch">
          <DispatchSettingsPanel
            key={`${wid}-${messageDispatchConfig.mode}-${messageDispatchConfig.storeMessageContent}`}
            initialConfig={messageDispatchConfig}
            onSave={saveSettings}
          />
        </Tabs.Content>
      </Tabs.Root>

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
