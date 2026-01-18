import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { EmptyState, ListItem } from '@/components/shared'
import type { StoredChat } from '@/types/messages'
import { formatFullDate } from '@/lib/date'
import { PREVIEW_TEXT_LENGTH } from '@/lib/constants'

interface ChatListProps {
  messages: StoredChat[]
  selected: StoredChat | null
  onSelect: (chat: StoredChat) => void
  onDelete: (id: string) => void
}

export function ChatList({ messages, selected, onSelect, onDelete }: ChatListProps) {
  if (messages.length === 0) {
    return <EmptyState message="No chat messages yet" />
  }

  return (
    <div className="divide-y">
      {messages.map((msg) => (
        <ListItem
          key={msg.id}
          id={msg.id}
          title={msg.chat.to?.join(', ') || 'Unknown'}
          subtitle={msg.chat.from ? `From: ${msg.chat.from}` : undefined}
          preview={msg.chat.message?.slice(0, PREVIEW_TEXT_LENGTH)}
          timestamp={msg.created_at}
          isSelected={selected?.id === msg.id}
          isSelectable
          onSelect={() => onSelect(msg)}
          onDelete={() => onDelete(msg.id)}
        />
      ))}
    </div>
  )
}

interface ChatDetailProps {
  chat: StoredChat
}

export function ChatDetail({ chat }: ChatDetailProps) {
  const { to, from, message, template_id, template_params, media_url, media_type, buttons, metadata } = chat.chat

  return (
    <>
      <div className="p-4 border-b">
        <h2 className="text-lg font-semibold">Chat Message</h2>
        <div className="text-sm text-muted-foreground mt-1 space-y-0.5">
          <div>To: {to?.join(', ')}</div>
          {from && <div>From: {from}</div>}
          <div>{formatFullDate(chat.created_at)}</div>
        </div>
      </div>
      <ScrollArea className="flex-1 p-4">
        <div className="bg-muted p-4 rounded-lg">{message}</div>

        {template_id && (
          <Section title="Template">
            <div className="text-sm bg-muted/50 p-2 rounded">
              <div>ID: {template_id}</div>
              {hasItems(template_params) && (
                <div>Params: {template_params?.join(', ')}</div>
              )}
            </div>
          </Section>
        )}

        {media_url && (
          <Section title="Media">
            <a
              href={media_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary underline text-sm"
            >
              {media_type || 'View Media'}
            </a>
          </Section>
        )}

        {hasItems(buttons) && (
          <Section title="Buttons">
            <div className="flex flex-wrap gap-2">
              {buttons?.map((btn, i) => (
                <Badge key={i} variant="outline">{btn.text}</Badge>
              ))}
            </div>
          </Section>
        )}

        {hasData(metadata) && (
          <Section title="Metadata">
            <pre className="text-xs bg-muted/50 p-2 rounded overflow-auto">
              {JSON.stringify(metadata, null, 2)}
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

function hasItems<T>(arr?: T[]): boolean {
  return Boolean(arr && arr.length > 0)
}

function hasData(data?: Record<string, string>): boolean {
  return Boolean(data && Object.keys(data).length > 0)
}
