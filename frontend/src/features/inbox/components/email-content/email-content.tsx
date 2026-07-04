import type { StoredEmail } from "../../inbox.types"
import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"

interface EmailContentProps {
  message: StoredEmail
  onDelete?: (messageId: string) => void
  isDeleting?: boolean
}

function formatFullTime(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleString([], {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    })
  } catch {
    return dateStr
  }
}

export function EmailContent({ message, onDelete, isDeleting = false }: EmailContentProps) {
  const sender = message.email.from_name
    ? `${message.email.from_name} <${message.email.from}>`
    : message.email.from || "Unknown sender"

  const recipients = (message.email.to || []).join(", ")
  const formattedDate = formatFullTime(message.created_at)

  const handleDelete = () => {
    if (onDelete) {
      onDelete(message.id)
    }
  }

  const bodyContent = message.email.html || `<div style="font-family: system-ui, -apple-system, sans-serif; font-size: 14px; line-height: 1.6; color: #111827; white-space: pre-wrap; padding: 20px;">${message.email.plain_text || ""}</div>`

  return (
    <div className="bg-card flex flex-col flex-1 h-full min-w-0 relative">
      <div className="h-14 flex items-center justify-between px-4 w-full border-b border-border shrink-0">
        <div className="flex flex-col gap-0.5 min-w-0 pr-4">
          <h2 className="text-sm font-semibold text-foreground truncate">
            {message.email.subject || "(No Subject)"}
          </h2>
          <p className="text-xs text-text-secondary truncate">
            To: {recipients || "—"} • {formattedDate}
          </p>
        </div>

        <Button
          variant="outline"
          size="sm"
          onClick={handleDelete}
          disabled={isDeleting}
          className="text-destructive border-destructive/20 hover:bg-destructive/10 shrink-0 h-9"
        >
          <Icon name="delete" size="sm" data-icon="inline-start" />
          Delete
        </Button>
      </div>

      <div className="bg-muted/10 px-4 py-3 border-b border-border flex flex-col gap-1.5 shrink-0 font-mono text-[11px] leading-relaxed">
        <div className="flex items-start">
          <span className="w-16 text-text-tertiary font-semibold">From:</span>
          <span className="text-foreground break-all">{sender}</span>
        </div>
        <div className="flex items-start">
          <span className="w-16 text-text-tertiary font-semibold">To:</span>
          <span className="text-foreground break-all">{recipients}</span>
        </div>
        {message.email.cc && message.email.cc.length > 0 && (
          <div className="flex items-start">
            <span className="w-16 text-text-tertiary font-semibold">Cc:</span>
            <span className="text-foreground break-all">{message.email.cc.join(", ")}</span>
          </div>
        )}
        <div className="flex items-start">
          <span className="w-16 text-text-tertiary font-semibold">Date:</span>
          <span className="text-foreground">{formattedDate}</span>
        </div>
        <div className="flex items-start">
          <span className="w-16 text-text-tertiary font-semibold">Message ID:</span>
          <span className="text-foreground font-mono">{message.id}</span>
        </div>
      </div>

      <div className="flex-1 w-full bg-white relative min-h-0">
        <iframe
          srcDoc={bodyContent}
          title="Email body render"
          sandbox="allow-same-origin"
          className="w-full h-full border-0 absolute inset-0 bg-white"
        />
      </div>
    </div>
  )
}
