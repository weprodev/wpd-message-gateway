import { useState } from "react"
import { useNavigate } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { SearchInput } from "@/components/ui/search-input"
import { Spinner } from "@/components/ui/spinner"
import { WorkspaceCard } from "../components/workspace-card"
import { WorkspaceActions } from "../components/workspace-actions"
import { EmptyState } from "../components/empty-state"
import { CreateWorkspaceModal } from "../components/create-workspace-modal"
import { SuccessModal, type SuccessModalVariant } from "../components/success-modal"
import { JoinWorkspaceModal } from "../components/join-workspace-modal"
import { useWorkspaces } from "../hooks/use-workspaces.hook"
import { fetchWorkspaces } from "../workspaces.api"
import type { Workspace } from "../workspace.types"

export function WorkspacesPage() {
  const navigate = useNavigate()
  const { workspaces, isLoading, error, reload } = useWorkspaces()
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState<string | null>(null)

  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [successModalOpen, setSuccessModalOpen] = useState(false)
  const [successVariant, setSuccessVariant] = useState<SuccessModalVariant>("created")
  const [joinModalOpen, setJoinModalOpen] = useState(false)
  const [newWorkspace, setNewWorkspace] = useState<{ id: string; name: string } | null>(null)

  const filteredWorkspaces = workspaces.filter(
    (workspace) =>
      workspace.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      workspace.slug.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  const handleSelectWorkspace = (workspace: Workspace) => {
    setSelectedWorkspaceId(workspace.id)
    navigate(ROUTES.workspace.overview(workspace.id))
  }

  const handleCreateWorkspace = () => {
    setCreateModalOpen(true)
  }

  const handleJoinWorkspace = () => {
    setJoinModalOpen(true)
  }

  const handleWorkspaceCreated = (workspaceId: string, workspaceName: string) => {
    setNewWorkspace({ id: workspaceId, name: workspaceName })
    setSuccessVariant("created")
    setSuccessModalOpen(true)
    reload()
  }

  const handleWorkspaceJoined = async (slug: string) => {
    const normalizedSlug = slug.trim().toLowerCase()
    const result = await fetchWorkspaces()
    const joined =
      result.ok
        ? result.workspaces.find((workspace) => workspace.slug.toLowerCase() === normalizedSlug)
        : undefined

    setNewWorkspace({
      id: joined?.id ?? "",
      name: joined?.name ?? slug.trim(),
    })
    setSuccessVariant("joined")
    setSuccessModalOpen(true)
    reload()
  }

  const handleCloseSuccessModal = () => {
    setSuccessModalOpen(false)
  }

  const handleContinueToDashboard = () => {
    if (newWorkspace?.id) {
      navigate(ROUTES.workspace.overview(newWorkspace.id))
    }
    setSuccessModalOpen(false)
  }

  const successModal = (
    <SuccessModal
      isOpen={successModalOpen}
      workspaceName={newWorkspace?.name ?? ""}
      variant={successVariant}
      onClose={handleCloseSuccessModal}
      onContinue={handleContinueToDashboard}
    />
  )

  if (isLoading) {
    return (
      <div className="flex min-h-[400px] flex-col items-center justify-center gap-3 py-16">
        <Spinner size="lg" />
        <span className="text-sm text-text-secondary">Loading workspaces...</span>
      </div>
    )
  }

  if (workspaces.length === 0) {
    return (
      <div className="flex flex-1 items-center justify-center px-6 py-16">
        <EmptyState
          onCreateWorkspace={handleCreateWorkspace}
          onJoinWorkspace={handleJoinWorkspace}
        />

        <CreateWorkspaceModal
          isOpen={createModalOpen}
          onClose={() => setCreateModalOpen(false)}
          onSuccess={handleWorkspaceCreated}
        />

        <JoinWorkspaceModal
          isOpen={joinModalOpen}
          onClose={() => setJoinModalOpen(false)}
          onSuccess={handleWorkspaceJoined}
        />

        {successModal}
      </div>
    )
  }

  return (
    <div className="flex w-full max-w-[720px] flex-col items-center gap-12 px-6 py-16">
      <div className="flex w-full max-w-[640px] flex-col items-center gap-2 text-center">
        <h1 className="text-4xl font-bold leading-tight text-foreground">Your Workspaces</h1>
        <p className="text-base leading-normal text-text-secondary">Select a workspace to continue</p>
        <p className="text-xs leading-4 text-text-tertiary">Connecting your world seamlessly!</p>
      </div>

      {error ? (
        <p className="w-full rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-center text-sm font-medium text-destructive">
          {error}
        </p>
      ) : null}

      <WorkspaceActions
        variant="card"
        onCreateWorkspace={handleCreateWorkspace}
        onJoinWorkspace={handleJoinWorkspace}
      />

      <SearchInput
        placeholder="Search workspaces..."
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        className="w-full max-w-[560px]"
      />

      <div className="flex w-full flex-col gap-3">
        {filteredWorkspaces.map((workspace) => (
          <WorkspaceCard
            key={workspace.id}
            workspace={workspace}
            isSelected={selectedWorkspaceId === workspace.id}
            onSelect={handleSelectWorkspace}
          />
        ))}
        {filteredWorkspaces.length === 0 ? (
          <p className="py-6 text-center text-sm text-text-tertiary">
            No workspaces matched your search.
          </p>
        ) : null}
      </div>

      <CreateWorkspaceModal
        isOpen={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSuccess={handleWorkspaceCreated}
      />

      {successModal}

      <JoinWorkspaceModal
        isOpen={joinModalOpen}
        onClose={() => setJoinModalOpen(false)}
        onSuccess={handleWorkspaceJoined}
      />
    </div>
  )
}
