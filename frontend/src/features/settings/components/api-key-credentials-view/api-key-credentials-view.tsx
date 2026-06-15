import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"

import type { ApiKeyCredentials } from "../../settings.types"

async function copyToClipboard(value: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(value)
    return true
  } catch {
    return false
  }
}

interface CredentialFieldProps {
  id: string
  label: string
  value: string
}

function CredentialField({ id, label, value }: CredentialFieldProps) {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    const ok = await copyToClipboard(value)
    if (!ok) return
    setCopied(true)
    window.setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={id} className="text-xs font-semibold uppercase text-text-secondary">
        {label}
      </label>
      <div className="flex gap-2">
        <Input id={id} readOnly value={value} className="font-mono text-xs bg-input" />
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="shrink-0"
          onClick={handleCopy}
          aria-label={`Copy ${label}`}
        >
          <Icon name={copied ? "check" : "content_copy"} size="sm" />
          {copied ? "Copied" : "Copy"}
        </Button>
      </div>
    </div>
  )
}

interface ApiKeyCredentialsViewProps {
  credentials: ApiKeyCredentials
  onConfirm: () => void
}

export function ApiKeyCredentialsView({ credentials, onConfirm }: ApiKeyCredentialsViewProps) {
  return (
    <div className="flex flex-col gap-5">
      <div
        role="alert"
        className="flex gap-3 rounded-lg border border-amber-500/30 bg-amber-500/10 p-4 text-sm text-foreground"
      >
        <Icon name="warning" size="sm" className="mt-0.5 shrink-0 text-amber-600 dark:text-amber-400" />
        <p className="leading-relaxed">
          This secret is shown <span className="font-semibold">only once</span>. Copy both values and store them
          securely before closing — you will not be able to view the secret again.
        </p>
      </div>

      <p className="text-sm text-text-secondary">
        Key: <span className="font-medium text-foreground">{credentials.keyName}</span>
      </p>

      <div className="flex flex-col gap-4">
        <CredentialField id="api-client-id" label="Client ID" value={credentials.clientId} />
        <CredentialField id="api-client-secret" label="Client secret" value={credentials.clientSecret} />
      </div>

      <Button type="button" className="w-full" onClick={onConfirm}>
        I&apos;ve saved my credentials — close
      </Button>
    </div>
  )
}
