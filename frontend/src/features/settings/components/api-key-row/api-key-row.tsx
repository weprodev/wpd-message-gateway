import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Permission, useWorkspaceAuthorization } from "@/core/auth"

import type { ApiKey } from "../../settings.types"

interface ApiKeyRowProps {
  apiKey: ApiKey
  onRegenerate: (id: string) => void
  onDelete: (id: string) => void
  isBusy?: boolean
}

function formatLastUsed(value?: string | null): string {
  if (!value) return "Never"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString()
}

export function ApiKeyRow({ apiKey, onRegenerate, onDelete, isBusy }: ApiKeyRowProps) {
  const { can } = useWorkspaceAuthorization()
  const canManageKeys = can(Permission.APIKeysWrite)

  return (
    <div className="flex items-center justify-between gap-4 border-b border-border px-4 py-4 last:border-b-0">
      <div className="min-w-0 flex-1">
        <p className="text-sm font-medium text-foreground">{apiKey.name}</p>
        <p className="mt-1 font-mono text-xs text-text-secondary">
          {apiKey.client_id.replace(/^(.{8}).*(.{8})$/, "............$2")}
        </p>
        <p className="mt-1 text-xs text-text-tertiary">
          Last used: {formatLastUsed(apiKey.last_used_at)}
        </p>
      </div>
      <div className="flex shrink-0 items-center gap-2">
        {canManageKeys ? (
          <>
            <Button
              type="button"
              variant="outline"
              size="sm"
              disabled={isBusy}
              onClick={() => onRegenerate(apiKey.id)}
            >
              <Icon name="refresh" size="sm" />
              Regenerate
            </Button>
            <Button
              type="button"
              variant="outline"
              size="sm"
              disabled={isBusy}
              onClick={() => onDelete(apiKey.id)}
            >
              <Icon name="delete" size="sm" />
              Delete
            </Button>
          </>
        ) : null}
      </div>
    </div>
  )
}
