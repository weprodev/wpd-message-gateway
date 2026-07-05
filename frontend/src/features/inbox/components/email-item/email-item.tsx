import type { StoredEmail } from "../../inbox.types"

interface EmailItemProps {
  message: StoredEmail
  isSelected: boolean
  onClick: () => void
}

function formatRecipients(to: string[] | undefined) {
  if (!to || to.length === 0) return "Unknown recipient"
  return to.join(", ")
}

function getRecipientInitial(recipient: string) {
  if (!recipient) return "?"
  const cleaned = recipient.trim()
  return cleaned[0]?.toUpperCase() ?? "?"
}

function formatTime(dateStr: string) {
  try {
    const d = new Date(dateStr)
    const now = new Date()
    if (d.toDateString() === now.toDateString()) {
      return d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
    }
    return d.toLocaleDateString([], { month: "short", day: "numeric" })
  } catch {
    return dateStr
  }
}

export function EmailItem({ message, isSelected, onClick }: EmailItemProps) {
  const recipientLabel = formatRecipients(message.email.to)
  const initial = getRecipientInitial(message.email.to?.[0] ?? recipientLabel)
  const timestamp = formatTime(message.created_at)
  
  const snippet = message.email.plain_text || message.email.html || message.email.subject || ""
  const cleanSnippet = snippet.replace(/<[^>]*>/g, "").slice(0, 60)

  const rowBg = isSelected ? "bg-primary-brand/10 border-l-4 border-primary-brand" : "bg-card hover:bg-muted/30"
  const paddingClass = isSelected ? "pl-3 pr-4" : "px-4"

  return (
    <div
      role="button"
      tabIndex={0}
      className={`${rowBg} ${paddingClass} py-3.5 flex flex-col gap-1.5 cursor-pointer border-b border-border transition-colors outline-none`}
      onClick={onClick}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault()
          onClick()
        }
      }}
    >
      <div className="flex items-center justify-between gap-2">
        <div className="flex items-center gap-2 min-w-0">
          <div className={`size-6 rounded-full flex items-center justify-center text-xs font-semibold shrink-0 ${
            isSelected ? "bg-primary-brand text-white" : "bg-secondary text-text-secondary"
          }`}>
            {initial}
          </div>
          <span className="text-sm font-semibold text-foreground truncate">
            {recipientLabel}
          </span>
        </div>
        <span className="text-xs text-text-tertiary shrink-0 font-medium">
          {timestamp}
        </span>
      </div>

      <div className="flex flex-col gap-0.5 min-w-0">
        <span className="text-xs font-semibold text-foreground truncate">
          {message.email.subject || "(No Subject)"}
        </span>
        <p className="text-xs text-text-secondary truncate">
          {cleanSnippet}
        </p>
      </div>
    </div>
  )
}
