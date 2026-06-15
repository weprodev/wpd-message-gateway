import { Button } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { Spinner } from "@/components/ui/spinner"

interface ApiKeyDeleteModalProps {
  isOpen: boolean
  onCancel: () => void
  onDelete: () => void
  isDeleting?: boolean
}

export function ApiKeyDeleteModal({ isOpen, onCancel, onDelete, isDeleting = false }: ApiKeyDeleteModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onCancel} title="Delete API key" preventDismiss={isDeleting}>
      <div className="flex flex-col gap-6">
        <p className="text-sm text-text-secondary leading-relaxed">
          Are you sure you want to delete this API key?
        </p>

        <div className="flex items-center justify-end gap-3 border-t border-border pt-4">
          <Button type="button" variant="outline" onClick={onCancel} disabled={isDeleting}>
            Cancel
          </Button>
          <Button type="button" variant="destructive" onClick={onDelete} disabled={isDeleting} className="min-w-[7.25rem]">
            {isDeleting ? (
              <>
                <Spinner size="sm" variant="onSolid" />
                Deleting…
              </>
            ) : (
              "Delete"
            )}
          </Button>
        </div>
      </div>
    </Modal>
  )
}
