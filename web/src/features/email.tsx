import { ScrollArea } from '@/components/ui/scroll-area'
import { EmptyState, ListItem } from '@/components/shared'
import type { StoredEmail } from '@/types/messages'
import { PREVIEW_TEXT_LENGTH } from '@/lib/constants'
import { formatFullDate } from '@/lib/date'

interface EmailListProps {
  emails: StoredEmail[]
  selected: StoredEmail | null
  onSelect: (email: StoredEmail) => void
  onDelete: (id: string) => void
}

export function EmailList({ emails, selected, onSelect, onDelete }: EmailListProps) {
  if (emails.length === 0) {
    return <EmptyState message="No emails yet" />
  }

  return (
    <div className="divide-y">
      {emails.map((email) => (
        <ListItem
          key={email.id}
          id={email.id}
          title={email.email.to?.join(', ') || 'Unknown'}
          timestamp={email.created_at}
          isSelected={selected?.id === email.id}
          isSelectable
          onSelect={() => onSelect(email)}
          onDelete={() => onDelete(email.id)}
        >
          <div className="text-sm font-medium truncate">
            {email.email.subject || '(No Subject)'}
          </div>
          <div className="text-xs text-muted-foreground truncate mt-1">
            {email.email.plain_text?.slice(0, PREVIEW_TEXT_LENGTH) || 'HTML content...'}
          </div>
        </ListItem>
      ))}
    </div>
  )
}

interface EmailDetailProps {
  email: StoredEmail
}

export function EmailDetail({ email }: EmailDetailProps) {
  const { from, from_name, to, subject, html, plain_text } = email.email

  return (
    <>
      <div className="p-4 border-b">
        <h2 className="text-lg font-semibold">{subject || '(No Subject)'}</h2>
        <div className="text-sm text-muted-foreground mt-1 space-y-0.5">
          <div>To: {to?.join(', ')}</div>
          {from && (
            <div>
              From: {from_name ? `${from_name} <${from}>` : from}
            </div>
          )}
          <div>{formatFullDate(email.created_at)}</div>
        </div>
      </div>
      <ScrollArea className="flex-1 p-4">
        {html ? (
          <div
            className="prose prose-sm dark:prose-invert max-w-none bg-white text-black p-4 rounded border"
            dangerouslySetInnerHTML={{ __html: html }}
          />
        ) : (
          <pre className="whitespace-pre-wrap text-sm font-mono">
            {plain_text}
          </pre>
        )}
      </ScrollArea>
    </>
  )
}
