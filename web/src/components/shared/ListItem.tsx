import { cn } from '@/lib/utils'
import { formatDate } from '@/lib/date'
import { DeleteButton } from './DeleteButton'

interface ListItemProps {
  id: string
  title: string
  subtitle?: string
  preview?: string
  timestamp: string
  isSelected?: boolean
  isSelectable?: boolean
  onSelect?: () => void
  onDelete: () => void
  children?: React.ReactNode
}

export function ListItem({
  title,
  subtitle,
  preview,
  timestamp,
  isSelected = false,
  isSelectable = false,
  onSelect,
  onDelete,
  children,
}: ListItemProps) {
  return (
    <div
      onClick={isSelectable ? onSelect : undefined}
      className={cn(
        'p-3 hover:bg-accent/50 transition-colors',
        isSelectable && 'cursor-pointer',
        isSelected && 'bg-accent'
      )}
    >
      <div className="flex justify-between items-start mb-1">
        <span className="font-medium text-sm truncate flex-1">{title}</span>
        <span className="text-xs text-muted-foreground ml-2">
          {formatDate(timestamp)}
        </span>
      </div>

      {subtitle && (
        <div className="text-xs text-muted-foreground mb-1">{subtitle}</div>
      )}

      {preview && (
        <div className="text-xs text-muted-foreground truncate mt-1">
          {preview}
        </div>
      )}

      {children}

      <DeleteButton onClick={onDelete} stopPropagation={isSelectable} />
    </div>
  )
}
