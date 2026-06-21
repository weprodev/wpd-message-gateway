import type React from "react"
import { useEffect, useState } from "react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"

import type { IntegrationViewModel } from "@/features/integrations/hooks/use-integrations.hook"
import { fetchProviderConfigFields, type ProviderConfigField } from "@/features/integrations/integrations.api"
import type { IntegrationActionResult } from "@/features/integrations/integrations.types"

interface ConnectModalProps {
  isOpen: boolean
  onClose: () => void
  workspaceId: string
  provider: IntegrationViewModel | null
  onConnect: (
    provider: IntegrationViewModel,
    config: Record<string, unknown>,
  ) => Promise<IntegrationActionResult>
}

export function ConnectModal({ isOpen, onClose, workspaceId, provider, onConnect }: ConnectModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={provider ? `Connect ${provider.name}` : undefined}>
      {isOpen && provider ? (
        <ConnectForm
          workspaceId={workspaceId}
          provider={provider}
          onConnect={onConnect}
          onClose={onClose}
        />
      ) : null}
    </Modal>
  )
}

interface ConnectFormProps {
  workspaceId: string
  provider: IntegrationViewModel
  onConnect: (
    provider: IntegrationViewModel,
    config: Record<string, unknown>,
  ) => Promise<IntegrationActionResult>
  onClose: () => void
}

function ConnectForm({ workspaceId, provider, onConnect, onClose }: ConnectFormProps) {
  const [fields, setFields] = useState<ProviderConfigField[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState<Record<string, string>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    const loadFields = async () => {
      setIsLoading(true)
      setError(null)
      try {
        const result = await fetchProviderConfigFields(workspaceId, provider.id)
        setFields(result)

        const defaults: Record<string, string> = {}
        result.forEach((field) => {
          defaults[field.key] = field.default_value || ""
        })
        setFormData(defaults)
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load configuration fields")
        setFields([])
      } finally {
        setIsLoading(false)
      }
    }

    void loadFields()
  }, [workspaceId, provider.id])

  function handleInputChange(key: string, value: string) {
    setFormData((prev) => ({ ...prev, [key]: value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setIsSubmitting(true)
    setError(null)
    try {
      const result = await onConnect(provider, formData)
      if (!result.ok) {
        setError(result.message ?? "Failed to connect provider")
      } else {
        onClose()
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to connect provider")
    } finally {
      setIsSubmitting(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center gap-3 py-8">
        <Spinner />
        <span className="text-sm text-text-secondary">Loading configuration...</span>
      </div>
    )
  }

  if (error && fields.length === 0) {
    return (
      <div className="flex flex-col gap-4 py-4">
        <p className="text-sm text-destructive">{error}</p>
        <Button onClick={onClose} variant="outline">
          Close
        </Button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-6">
      {error ? (
        <p className="rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </p>
      ) : null}

      <div className="flex max-h-[350px] flex-col gap-4 overflow-y-auto pr-1">
        {fields.map((field) => (
          <div key={field.key} className="flex flex-col gap-1.5">
            <label htmlFor={field.key} className="text-sm font-medium text-foreground">
              {field.label}
              {field.required ? <span className="ml-1 text-destructive">*</span> : null}
            </label>
            <Input
              id={field.key}
              type={field.field_type === "password" ? "password" : field.field_type === "email" ? "email" : "text"}
              required={field.required}
              value={formData[field.key] || ""}
              onChange={(e) => handleInputChange(field.key, e.target.value)}
              placeholder={field.description || `Enter ${field.label.toLowerCase()}`}
              className="bg-input"
            />
            {field.description ? (
              <span className="text-[12px] leading-normal text-text-secondary">{field.description}</span>
            ) : null}
          </div>
        ))}
      </div>

      <div className="flex items-center justify-end gap-3 border-t border-border pt-2">
        <Button type="button" variant="outline" onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={isSubmitting}
          className="bg-primary-brand hover:bg-primary-brand-hover"
        >
          {isSubmitting ? (
            <div className="flex items-center gap-2">
              <Spinner size="sm" variant="onSolid" />
              <span>Connecting...</span>
            </div>
          ) : (
            "Connect"
          )}
        </Button>
      </div>
    </form>
  )
}
