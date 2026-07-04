import { useEffect, useState } from "react"

import { fetchProviderConfigFields, type ProviderConfigField } from "@/features/integrations/integrations.api"

export function useProviderConfigFields(workspaceId: string, providerId: string | undefined, enabled: boolean) {
  const [fields, setFields] = useState<ProviderConfigField[]>([])
  const [formData, setFormData] = useState<Record<string, string>>({})
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!enabled || !providerId) {
      return
    }

    let cancelled = false

    ;(async () => {
      setIsLoading(true)
      setError(null)
      try {
        const result = await fetchProviderConfigFields(workspaceId, providerId)
        if (cancelled) return

        setFields(result)
        const defaults: Record<string, string> = {}
        result.forEach((field) => {
          defaults[field.key] = field.default_value || ""
        })
        setFormData(defaults)
      } catch (err) {
        if (cancelled) return
        setError(err instanceof Error ? err.message : "Failed to load configuration fields")
        setFields([])
      } finally {
        if (!cancelled) setIsLoading(false)
      }
    })()

    return () => {
      cancelled = true
    }
  }, [workspaceId, providerId, enabled])

  function updateField(key: string, value: string) {
    setFormData((prev) => ({ ...prev, [key]: value }))
  }

  return {
    fields,
    formData,
    isLoading,
    error,
    updateField,
    setError,
  }
}
