import { Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface DeleteButtonProps {
  onClick: () => void
  stopPropagation?: boolean
}

export function DeleteButton({ onClick, stopPropagation = false }: DeleteButtonProps) {
  const handleClick = (e: React.MouseEvent) => {
    if (stopPropagation) e.stopPropagation()
    onClick()
  }

  return (
    <div className="flex justify-end mt-2">
      <Button
        variant="ghost"
        size="icon"
        className="h-6 w-6"
        onClick={handleClick}
      >
        <Trash2 className="h-3 w-3" />
      </Button>
    </div>
  )
}
