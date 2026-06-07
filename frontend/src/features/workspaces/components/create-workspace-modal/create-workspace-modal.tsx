import { useState } from "react"
import { Modal } from "@/components/ui/modal"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Icon } from "@/components/ui/icon"
import { useCreateWorkspace } from "../../hooks/use-create-workspace.hook"
import { WORKSPACE_ICON_OPTIONS } from "../../workspace-icons"

interface CreateWorkspaceModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess: (workspaceId: string, workspaceName: string) => void
}
export function CreateWorkspaceModal({ isOpen, onClose, onSuccess }: CreateWorkspaceModalProps) {
  const [name, setName] = useState("")
  const [slug, setSlug] = useState("")
  const [iconKey, setIconKey] = useState("package")
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  const { createWorkspace, isLoading, error: apiError } = useCreateWorkspace()

  const handleNameChange = (val: string) => {
    setName(val)
    setValidationErrors((prev) => ({ ...prev, name: "" }))
    const generatedSlug = val
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/(^-|-$)/g, "")
    setSlug(generatedSlug)
  }

  const validate = () => {
    const errors: Record<string, string> = {}
    if (!name.trim()) {
      errors.name = "Workspace name is required"
    }
    if (!slug.trim()) {
      errors.slug = "Workspace slug is required"
    } else if (!/^[a-z0-9-]+$/.test(slug)) {
      errors.slug = "Workspace slug must only contain lowercase letters, numbers, and dashes"
    }
    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    try {
      const workspace = await createWorkspace({
        name: name.trim(),
        slug: slug.trim(),
        icon_key: iconKey,
      })
      onSuccess(workspace.id, workspace.name)
      handleClose()
    } catch (err) {
      console.error(err)
    }
  }

  const handleClose = () => {
    setName("")
    setSlug("")
    setIconKey("package")
    setValidationErrors({})
    onClose()
  }

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Create New Workspace">
      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <p className="text-sm text-text-secondary">
          Set up a new workspace for your team and communications.
        </p>

        <div className="flex flex-col gap-1.5">
          <label htmlFor="ws-name" className="text-xs font-semibold text-text-secondary uppercase">
            Workspace Name
          </label>
          <Input
            id="ws-name"
            type="text"
            value={name}
            onChange={(e) => handleNameChange(e.target.value)}
            placeholder="e.g. Engineering Team"
            className="w-full bg-secondary border-border h-11 px-4"
          />
          {validationErrors.name && (
            <p className="text-xs text-destructive font-medium mt-0.5">{validationErrors.name}</p>
          )}
        </div>

        <div className="flex flex-col gap-1.5">
          <label htmlFor="ws-slug" className="text-xs font-semibold text-text-secondary uppercase">
            Workspace Slug
          </label>
          <Input
            id="ws-slug"
            type="text"
            value={slug}
            onChange={(e) => setSlug(e.target.value)}
            placeholder="e.g. engineering-team"
            className="w-full bg-secondary border-border h-11 px-4 font-mono"
          />
          {validationErrors.slug && (
            <p className="text-xs text-destructive font-medium mt-0.5">{validationErrors.slug}</p>
          )}
        </div>

        <div className="flex flex-col gap-2">
          <span className="text-xs font-semibold text-text-secondary uppercase">
            Workspace Icon
          </span>
          <div className="grid grid-cols-4 gap-2">
            {WORKSPACE_ICON_OPTIONS.map((item) => {
              const isSelected = iconKey === item.key
              return (
                <button
                  key={item.key}
                  type="button"
                  onClick={() => setIconKey(item.key)}
                  className={`flex flex-col items-center justify-center p-3 rounded-lg border text-xs font-medium transition-all ${
                    isSelected
                      ? "border-primary-brand bg-primary-brand/5 text-primary-brand"
                      : "border-border bg-card text-text-secondary hover:border-text-secondary"
                  }`}
                >
                  <Icon name={item.iconName} size="sm" className="mb-1 shrink-0" />
                  {item.label}
                </button>
              )
            })}
          </div>
        </div>

        {apiError && (
          <p className="text-xs text-destructive text-center font-medium mt-2">{apiError}</p>
        )}

        <div className="flex flex-col gap-3 pt-2">
          <Button type="submit" disabled={isLoading}>
            {isLoading ? "Creating..." : "Create Workspace"}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
