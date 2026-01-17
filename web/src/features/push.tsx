import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import type { StoredPush } from '@/types/messages'
import { formatDate } from '@/lib/date'

interface PushListProps {
  notifications: StoredPush[]
  onDelete: (id: string) => void
}

export function PushList({ notifications, onDelete }: PushListProps) {
  if (notifications.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <p className="text-sm">No push notifications yet</p>
      </div>
    )
  }

  return (
    <div className="divide-y">
      {notifications.map((notif) => (
        <div key={notif.id} className="p-3 hover:bg-accent/50 transition-colors">
          <div className="flex justify-between items-start mb-1">
            <span className="font-medium text-sm">{notif.push.title}</span>
            <span className="text-xs text-muted-foreground">
              {formatDate(notif.created_at)}
            </span>
          </div>
          <div className="text-xs text-muted-foreground mb-1">
            {notif.push.device_tokens?.length ?? 0} device(s)
          </div>
          <div className="text-sm bg-muted p-2 rounded mt-2">
            {notif.push.body}
          </div>
          {notif.push.data && Object.keys(notif.push.data).length > 0 && (
            <pre className="text-xs bg-muted/50 p-2 rounded mt-2 overflow-auto">
              {JSON.stringify(notif.push.data, null, 2)}
            </pre>
          )}
          <div className="flex justify-end mt-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => onDelete(notif.id)}
            >
              <Trash2 className="h-3 w-3" />
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
