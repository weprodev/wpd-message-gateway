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
import { ROUTES } from "@/core/router/routes"

interface GetStartedDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  workspaceId?: string
}

export function GetStartedDialog({ open, onOpenChange, workspaceId }: GetStartedDialogProps) {
  const settingsPath = workspaceId ? ROUTES.workspace.settings(workspaceId) : ROUTES.workspaces

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[85vh] max-w-lg gap-0 overflow-hidden p-0">
        <DialogHeader className="border-b bg-muted/20 px-5 py-4 text-left">
          <DialogTitle>Get started</DialogTitle>
          <DialogDescription>
            Send messages through the same gateway endpoints your integrations use in production.
          </DialogDescription>
        </DialogHeader>
        <div className="flex flex-col gap-4 overflow-y-auto p-6 text-sm text-foreground">
          <p className="text-xs leading-relaxed text-text-secondary">
            Create an API key in{" "}
            <Link to={settingsPath} className="font-semibold text-primary-brand hover:underline">
              Settings
            </Link>
            , then authenticate each request with Basic auth and your workspace slug header.
          </p>
          <div className="flex flex-col gap-3">
            <div>
              <h4 className="mb-1 text-xs font-bold uppercase text-text-secondary">Email</h4>
              <pre className="overflow-x-auto rounded-md border bg-muted/80 p-2.5 font-mono text-xs whitespace-pre-wrap">{`curl -X POST http://localhost:10101/v1/email \\
  -u "client_id:client_secret" \\
  -H "X-Workspace-Key: your-workspace-slug" \\
  -H "Content-Type: application/json" \\
  -d '{"to":["user@example.com"],"subject":"Test","html":"<p>Hi</p>"}'`}</pre>
            </div>
            <div>
              <h4 className="mb-1 text-xs font-bold uppercase text-text-secondary">SMS</h4>
              <pre className="overflow-x-auto rounded-md border bg-muted/80 p-2.5 font-mono text-xs whitespace-pre-wrap">{`curl -X POST http://localhost:10101/v1/sms \\
  -u "client_id:client_secret" \\
  -H "X-Workspace-Key: your-workspace-slug" \\
  -H "Content-Type: application/json" \\
  -d '{"to":["+1234567890"],"message":"Hello"}'`}</pre>
            </div>
          </div>
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
