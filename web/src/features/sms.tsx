import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import type { StoredSMS } from '@/types/messages'
import { formatDate } from '@/lib/date'

interface SMSListProps {
  messages: StoredSMS[]
  onDelete: (id: string) => void
}

export function SMSList({ messages, onDelete }: SMSListProps) {
  if (messages.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <p className="text-sm">No SMS messages yet</p>
      </div>
    )
  }

  return (
    <div className="divide-y">
      {messages.map((msg) => (
        <div key={msg.id} className="p-3 hover:bg-accent/50 transition-colors">
          <div className="flex justify-between items-start mb-1">
            <span className="font-medium text-sm">
              {msg.sms.to?.join(', ')}
            </span>
            <span className="text-xs text-muted-foreground">
              {formatDate(msg.created_at)}
            </span>
          </div>
          {msg.sms.from && (
            <div className="text-xs text-muted-foreground mb-1">
              From: {msg.sms.from}
            </div>
          )}
          <div className="text-sm bg-muted p-2 rounded mt-2">
            {msg.sms.message}
          </div>
          <div className="flex justify-end mt-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => onDelete(msg.id)}
            >
              <Trash2 className="h-3 w-3" />
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
