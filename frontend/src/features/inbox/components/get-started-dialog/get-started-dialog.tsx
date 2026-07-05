import { Link } from "react-router-dom"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Spinner } from "@/components/ui/spinner"
import { ROUTES } from "@/core/router/routes"
import { buildGatewayCurlExamples } from "../../get-started-dialog.utils"
import { useGetStartedContext } from "../../hooks/use-get-started-context.hook"

interface GetStartedDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  workspaceId?: string
}

function developerSettingsPath(workspaceId: string) {
  return `${ROUTES.workspace.settings(workspaceId)}?tab=developer`
}

export function GetStartedDialog({ open, onOpenChange, workspaceId }: GetStartedDialogProps) {
  const { context, isLoading, error } = useGetStartedContext(workspaceId, open)
  const developerPath = workspaceId ? developerSettingsPath(workspaceId) : ROUTES.workspaces

  const primaryKey = context?.apiKeys[0]
  const examples =
    context && primaryKey
      ? buildGatewayCurlExamples(context.workspaceId, primaryKey.client_id)
      : null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] max-w-lg gap-0 overflow-hidden p-0">
        <DialogHeader className="border-b bg-muted/20 px-5 py-4 text-left">
          <DialogTitle>Get started</DialogTitle>
          <DialogDescription>
            Send test messages with your workspace API key and ID.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-4 overflow-y-auto p-6 text-sm text-foreground">
          {isLoading ? (
            <div className="flex items-center gap-2 text-sm text-text-secondary">
              <Spinner size="sm" />
              Loading workspace credentials…
            </div>
          ) : null}

          {error ? (
            <p className="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-xs text-destructive">
              {error}
            </p>
          ) : null}

          {context ? (
            <div className="rounded-lg border border-border bg-muted/30 p-4">
              <h3 className="text-xs font-bold uppercase tracking-wide text-text-secondary">
                Your credentials
              </h3>
              <dl className="mt-3 flex flex-col gap-2 text-xs">
                <div className="flex flex-col gap-0.5 sm:flex-row sm:items-baseline sm:gap-2">
                  <dt className="shrink-0 font-medium text-text-secondary">Workspace ID</dt>
                  <dd className="font-mono text-foreground break-all">{context.workspaceId}</dd>
                </div>
                <div className="flex flex-col gap-0.5 sm:flex-row sm:items-baseline sm:gap-2">
                  <dt className="shrink-0 font-medium text-text-secondary">Header</dt>
                  <dd className="font-mono text-foreground break-all">
                    X-Workspace-Key: {context.workspaceId}
                  </dd>
                </div>
                {primaryKey ? (
                  <div className="flex flex-col gap-0.5 sm:flex-row sm:items-baseline sm:gap-2">
                    <dt className="shrink-0 font-medium text-text-secondary">Client ID</dt>
                    <dd className="font-mono text-foreground">{primaryKey.client_id}</dd>
                  </div>
                ) : (
                  <p className="text-text-secondary">
                    No API keys yet. Create one in{" "}
                    <Link to={developerPath} className="font-semibold text-primary-brand hover:underline">
                      Settings → Developer
                    </Link>
                    .
                  </p>
                )}
                <div className="flex flex-col gap-0.5">
                  <dt className="font-medium text-text-secondary">Client secret</dt>
                  <dd className="text-text-secondary">
                    Shown once when you generate or regenerate a key in{" "}
                    <Link to={developerPath} className="font-semibold text-primary-brand hover:underline">
                      Settings → Developer
                    </Link>
                    . Replace <span className="font-mono text-foreground">YOUR_CLIENT_SECRET</span> in the
                    examples below.
                  </dd>
                </div>
              </dl>
            </div>
          ) : null}

          {examples ? (
            <div className="flex flex-col gap-3">
              <div>
                <h4 className="mb-1 text-xs font-bold uppercase text-text-secondary">Email</h4>
                <pre className="overflow-x-auto rounded-md border bg-muted/80 p-2.5 font-mono text-xs whitespace-pre-wrap">
                  {examples.email}
                </pre>
              </div>
              <div>
                <h4 className="mb-1 text-xs font-bold uppercase text-text-secondary">SMS</h4>
                <pre className="overflow-x-auto rounded-md border bg-muted/80 p-2.5 font-mono text-xs whitespace-pre-wrap">
                  {examples.sms}
                </pre>
              </div>
            </div>
          ) : null}

          {!isLoading && !context && !error ? (
            <p className="text-xs text-text-secondary">Open this dialog from a workspace to load credentials.</p>
          ) : null}
        </div>
        <DialogFooter className="border-t bg-muted/20 px-5 py-4">
          <Button type="button" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
