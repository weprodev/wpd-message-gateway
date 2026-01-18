import { EmptyState, ListItem } from '@/components/shared'
import type { StoredPush } from '@/types/messages'

interface PushListProps {
  notifications: StoredPush[]
  onDelete: (id: string) => void
}

export function PushList({ notifications, onDelete }: PushListProps) {
  if (notifications.length === 0) {
    return <EmptyState message="No push notifications yet" />
  }

  return (
    <div className="divide-y">
      {notifications.map((notif) => (
        <ListItem
          key={notif.id}
          id={notif.id}
          title={notif.push.title}
          subtitle={`${notif.push.device_tokens?.length ?? 0} device(s)`}
          timestamp={notif.created_at}
          onDelete={() => onDelete(notif.id)}
        >
          <div className="text-sm bg-muted p-2 rounded mt-2">
            {notif.push.body}
          </div>
          {hasData(notif.push.data) && (
            <pre className="text-xs bg-muted/50 p-2 rounded mt-2 overflow-auto">
              {JSON.stringify(notif.push.data, null, 2)}
            </pre>
          )}
        </ListItem>
      ))}
    </div>
  )
}

function hasData(data?: Record<string, string>): boolean {
  return Boolean(data && Object.keys(data).length > 0)
}
