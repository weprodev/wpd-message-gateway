import { Modal } from "@/components/ui/modal"

import { ApiKeyCredentialsView } from "../api-key-credentials-view"
import type { ApiKeyCredentials } from "../../settings.types"

interface ApiKeyCredentialsModalProps {
  isOpen: boolean
  credentials: ApiKeyCredentials | null
  onConfirm: () => void
}

export function ApiKeyCredentialsModal({ isOpen, credentials, onConfirm }: ApiKeyCredentialsModalProps) {
  if (!credentials) {
    return null
  }

  const title = credentials.mode === "created" ? "API key created" : "API key regenerated"

  return (
    <Modal isOpen={isOpen} onClose={onConfirm} title={title} preventDismiss>
      <ApiKeyCredentialsView credentials={credentials} onConfirm={onConfirm} />
    </Modal>
  )
}
