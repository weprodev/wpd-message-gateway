import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { StoredChat } from '@/types/messages'
import { cn } from '@/lib/utils'
import { formatDate, formatFullDate } from '@/lib/date'

interface ChatListProps {
  messages: StoredChat[]
  selected: StoredChat | null
  onSelect: (chat: StoredChat) => void
  onDelete: (id: string) => void
}

export function ChatList({ messages, selected, onSelect, onDelete }: ChatListProps) {
  if (messages.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <p className="text-sm">No chat messages yet</p>
      </div>
    )
  }

  return (
    <div className="divide-y">
      {messages.map((msg) => (
        <div
          key={msg.id}
          onClick={() => onSelect(msg)}
          className={cn(
            "p-3 cursor-pointer hover:bg-accent/50 transition-colors",
            selected?.id === msg.id && "bg-accent"
          )}
        >
          <div className="flex justify-between items-start mb-1">
            <span className="font-medium text-sm truncate flex-1">
              {msg.chat.to?.join(', ')}
            </span>
            <span className="text-xs text-muted-foreground ml-2">
              {formatDate(msg.created_at)}
            </span>
          </div>
          {msg.chat.from && (
            <div className="text-xs text-muted-foreground">
              From: {msg.chat.from}
            </div>
          )}
          <div className="text-xs text-muted-foreground truncate mt-1">
            {msg.chat.message?.slice(0, 80)}
          </div>
          <div className="flex justify-end mt-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={(e) => {
                e.stopPropagation()
                onDelete(msg.id)
              }}
            >
              <Trash2 className="h-3 w-3" />
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}

interface ChatDetailProps {
  chat: StoredChat
}

export function ChatDetail({ chat }: ChatDetailProps) {
  return (
    <>
      <div className="p-4 border-b">
        <h2 className="text-lg font-semibold">Chat Message</h2>
        <div className="text-sm text-muted-foreground mt-1 space-y-0.5">
          <div>To: {chat.chat.to?.join(', ')}</div>
          {chat.chat.from && <div>From: {chat.chat.from}</div>}
          <div>{formatFullDate(chat.created_at)}</div>
        </div>
      </div>
      <ScrollArea className="flex-1 p-4">
        <div className="bg-muted p-4 rounded-lg">
          {chat.chat.message}
        </div>

        {chat.chat.template_id && (
          <Section title="Template">
            <div className="text-sm bg-muted/50 p-2 rounded">
              <div>ID: {chat.chat.template_id}</div>
              {(chat.chat.template_params?.length ?? 0) > 0 && (
                <div>Params: {chat.chat.template_params?.join(', ')}</div>
              )}
            </div>
          </Section>
        )}

        {chat.chat.media_url && (
          <Section title="Media">
            <a
              href={chat.chat.media_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary underline text-sm"
            >
              {chat.chat.media_type || 'View Media'}
            </a>
          </Section>
        )}

        {(chat.chat.buttons?.length ?? 0) > 0 && (
          <Section title="Buttons">
            <div className="flex flex-wrap gap-2">
              {chat.chat.buttons?.map((btn, i) => (
                <Badge key={i} variant="outline">{btn.text}</Badge>
              ))}
            </div>
          </Section>
        )}

        {chat.chat.metadata && Object.keys(chat.chat.metadata).length > 0 && (
          <Section title="Metadata">
            <pre className="text-xs bg-muted/50 p-2 rounded overflow-auto">
              {JSON.stringify(chat.chat.metadata, null, 2)}
            </pre>
          </Section>
        )}
      </ScrollArea>
    </>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="mt-4">
      <h3 className="text-sm font-medium mb-2">{title}</h3>
      {children}
    </div>
  )
}
