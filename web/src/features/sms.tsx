import { EmptyState, ListItem } from '@/components/shared'
import type { StoredSMS } from '@/types/messages'

interface SMSListProps {
  messages: StoredSMS[]
  onDelete: (id: string) => void
}

export function SMSList({ messages, onDelete }: SMSListProps) {
  if (messages.length === 0) {
    return <EmptyState message="No SMS messages yet" />
  }

  return (
    <div className="divide-y">
      {messages.map((msg) => (
        <ListItem
          key={msg.id}
          id={msg.id}
          title={msg.sms.to?.join(', ') || 'Unknown'}
          subtitle={msg.sms.from ? `From: ${msg.sms.from}` : undefined}
          timestamp={msg.created_at}
          onDelete={() => onDelete(msg.id)}
        >
          <div className="text-sm bg-muted p-2 rounded mt-2">
            {msg.sms.message}
          </div>
        </ListItem>
      ))}
    </div>
  )
}
