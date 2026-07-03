import { useState } from "react"
import type { StoredEmail } from "../../inbox.types"
import { Button } from "@/components/ui/button"
import { SearchInput } from "@/components/ui/search-input"
import { EmailItem } from "../email-item"

interface EmailListProps {
  messages: StoredEmail[]
  selectedMessageId: string | null
  onSelectMessage: (messageId: string) => void
  hasMore?: boolean
  isLoadingMore?: boolean
  onLoadMore?: () => void
}

export function EmailList({
  messages,
  selectedMessageId,
  onSelectMessage,
  hasMore = false,
  isLoadingMore = false,
  onLoadMore,
}: EmailListProps) {
  const [searchQuery, setSearchQuery] = useState("")

  const filtered = messages.filter((msg) => {
    const q = searchQuery.toLowerCase()
    const subject = (msg.email.subject || "").toLowerCase()
    const sender = (msg.email.from_name || msg.email.from || "").toLowerCase()
    return subject.includes(q) || sender.includes(q)
  })

  return (
    <div className="bg-card flex flex-col h-full items-start w-full md:w-[320px] border-r border-border shrink-0 min-w-0">
      <div className="h-14 flex items-center justify-between px-4 w-full border-b border-border shrink-0">
        <span className="text-sm font-semibold text-foreground">Inbox</span>
        <div className="bg-secondary px-2.5 py-0.5 rounded-full border border-border text-xs text-text-secondary font-semibold">
          {messages.length}
        </div>
      </div>

      <div className="p-2 w-full border-b border-border shrink-0 bg-muted/10 flex items-center justify-center">
        <SearchInput
          placeholder="Search inbox..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <div className="flex-1 w-full overflow-y-auto divide-y divide-border">
        {filtered.map((message) => (
          <EmailItem
            key={message.id}
            message={message}
            isSelected={selectedMessageId === message.id}
            onClick={() => onSelectMessage(message.id)}
          />
        ))}

        {filtered.length === 0 && (
          <div className="text-center text-xs text-text-tertiary py-12 px-4">
            {searchQuery ? "No messages match your search." : "Your inbox is empty."}
          </div>
        )}

        {hasMore && !searchQuery && onLoadMore && (
          <div className="p-3 border-t border-border">
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="w-full"
              disabled={isLoadingMore}
              onClick={onLoadMore}
            >
              {isLoadingMore ? "Loading..." : "Load more"}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
