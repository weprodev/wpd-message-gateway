import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import type { StoredEmail } from '@/types/messages'
import { cn } from '@/lib/utils'
import { formatDate, formatFullDate } from '@/lib/date'

interface EmailListProps {
  emails: StoredEmail[]
  selected: StoredEmail | null
  onSelect: (email: StoredEmail) => void
  onDelete: (id: string) => void
}

export function EmailList({ emails, selected, onSelect, onDelete }: EmailListProps) {
  if (emails.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <p className="text-sm">No emails yet</p>
      </div>
    )
  }

  return (
    <div className="divide-y">
      {emails.map((email) => (
        <div
          key={email.id}
          onClick={() => onSelect(email)}
          className={cn(
            "p-3 cursor-pointer hover:bg-accent/50 transition-colors",
            selected?.id === email.id && "bg-accent"
          )}
        >
          <div className="flex justify-between items-start mb-1">
            <span className="font-medium text-sm truncate flex-1">
              {email.email.to?.join(', ') || 'Unknown'}
            </span>
            <span className="text-xs text-muted-foreground ml-2">
              {formatDate(email.created_at)}
            </span>
          </div>
          <div className="text-sm font-medium truncate">
            {email.email.subject || '(No Subject)'}
          </div>
          <div className="text-xs text-muted-foreground truncate mt-1">
            {email.email.plain_text?.slice(0, 100) || 'HTML content...'}
          </div>
          <div className="flex justify-end mt-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={(e) => {
                e.stopPropagation()
                onDelete(email.id)
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

interface EmailDetailProps {
  email: StoredEmail
}

export function EmailDetail({ email }: EmailDetailProps) {
  return (
    <>
      <div className="p-4 border-b">
        <h2 className="text-lg font-semibold">
          {email.email.subject || '(No Subject)'}
        </h2>
        <div className="text-sm text-muted-foreground mt-1 space-y-0.5">
          <div>To: {email.email.to?.join(', ')}</div>
          {email.email.from && (
            <div>
              From: {email.email.from_name 
                ? `${email.email.from_name} <${email.email.from}>` 
                : email.email.from}
            </div>
          )}
          <div>{formatFullDate(email.created_at)}</div>
        </div>
      </div>
      <ScrollArea className="flex-1 p-4">
        {email.email.html ? (
          <div
            className="prose prose-sm dark:prose-invert max-w-none bg-white text-black p-4 rounded border"
            dangerouslySetInnerHTML={{ __html: email.email.html }}
          />
        ) : (
          <pre className="whitespace-pre-wrap text-sm font-mono">
            {email.email.plain_text}
          </pre>
        )}
      </ScrollArea>
    </>
  )
}
