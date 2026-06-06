import { useState } from "react"
import { useNavigate } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { SearchInput } from "@/components/ui/search-input"
import { Spinner } from "@/components/ui/spinner"
import { WorkspaceCard } from "../components/workspace-card"
import { WorkspaceActions } from "../components/workspace-actions"
import { EmptyState } from "../components/empty-state"
import { CreateWorkspaceModal } from "../components/create-workspace-modal"
import { SuccessModal } from "../components/success-modal"
import { JoinWorkspaceModal } from "../components/join-workspace-modal"
import { useWorkspaces } from "../hooks/use-workspaces.hook"
import type { Workspace } from "../workspace.types"

export function WorkspacesPage() {
  const navigate = useNavigate()
  const { workspaces, isLoading, error, reload } = useWorkspaces()
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedWorkspaceId, setSelectedWorkspaceId] = useState<string | null>(null)
  
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [successModalOpen, setSuccessModalOpen] = useState(false)
  const [joinModalOpen, setJoinModalOpen] = useState(false)
  const [newWorkspace, setNewWorkspace] = useState<{ id: string; name: string } | null>(null)

  const filteredWorkspaces = workspaces.filter((workspace) =>
    workspace.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    workspace.unique_key.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const handleSelectWorkspace = (workspace: Workspace) => {
    setSelectedWorkspaceId(workspace.id)
    // Directly navigate to the workspace overview
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
    setSuccessModalOpen(true)
    reload()
  }

  const handleWorkspaceJoined = () => {
    reload()
  }

  const handleContinueToDashboard = () => {
    if (newWorkspace) {
      navigate(ROUTES.workspace.overview(newWorkspace.id))
    }
    setSuccessModalOpen(false)
  }

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px] gap-3">
        <Spinner size="lg" />
        <span className="text-sm text-text-secondary">Loading workspaces...</span>
      </div>
    )
  }

  if (workspaces.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center py-12">
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
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center gap-10 py-6 max-w-2xl mx-auto w-full">
      <div className="flex flex-col gap-2 items-center text-center">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">
          Your Workspaces
        </h1>
        <p className="text-sm text-text-secondary">
          Select a workspace to view logs and manage integrations.
        </p>
      </div>

      {error ? (
        <p className="text-sm text-destructive bg-destructive/10 border border-destructive/20 rounded-lg p-3 w-full text-center font-medium">
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
      />

      <div className="flex flex-col gap-3 w-full">
        {filteredWorkspaces.map((workspace) => (
          <WorkspaceCard
            key={workspace.id}
            workspace={workspace}
            isSelected={selectedWorkspaceId === workspace.id}
            onSelect={handleSelectWorkspace}
          />
        ))}
        {filteredWorkspaces.length === 0 && (
          <p className="text-sm text-text-tertiary text-center py-6">
            No workspaces matched your search.
          </p>
        )}
      </div>

      <CreateWorkspaceModal
        isOpen={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSuccess={handleWorkspaceCreated}
      />

      <SuccessModal
        isOpen={successModalOpen}
        workspaceName={newWorkspace?.name || ""}
        onContinue={handleContinueToDashboard}
      />

      <JoinWorkspaceModal
        isOpen={joinModalOpen}
        onClose={() => setJoinModalOpen(false)}
        onSuccess={handleWorkspaceJoined}
      />
    </div>
  )
}
