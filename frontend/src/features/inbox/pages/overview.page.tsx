import { useState } from "react"
import { useParams } from "react-router-dom"

import { Icon } from "@/components/ui/icon"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { SendTestModal } from "../components/send-test-modal"
import { useInboxLogs } from "../hooks/use-inbox-logs.hook"
import type { MessageChannel } from "../inbox.types"

interface OverviewPageProps {
  channel?: MessageChannel
}

function formatTime(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    })
  } catch {
    return dateStr
  }
}

function statusBadgeClass(code: number) {
  if (code >= 200 && code < 300) {
    return "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-semibold"
  }
  if (code >= 400) {
    return "bg-rose-500/10 text-rose-600 dark:text-rose-400 font-semibold"
  }
  return "bg-amber-500/10 text-amber-600 dark:text-amber-400 font-semibold"
}

function channelIcon(channelType: string) {
  switch (channelType.toLowerCase()) {
    case "email":
      return "mail"
    case "sms":
      return "forum"
    case "push":
      return "notifications"
    case "chat":
      return "chat"
    default:
      return "sms"
  }
}

function exampleEndpoint(channel?: MessageChannel) {
  return channel ? `POST /v1/${channel}` : "POST /v1/email"
}

export function OverviewPage({ channel }: OverviewPageProps) {
  const { wid } = useParams<{ wid: string }>()
  const { logs, total, isLoading, error, reload } = useInboxLogs(wid, channel)
  const [isTestModalOpen, setIsTestModalOpen] = useState(false)
  const [isDocsOpen, setIsDocsOpen] = useState(false)

  return (
    <div className="flex w-full flex-col gap-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">Message Logs</h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            Monitor gateway API requests across channels.
          </p>
        </div>
        {logs.length > 0 && (
          <div className="flex shrink-0 items-center gap-2">
            <Button type="button" variant="outline" size="sm" onClick={() => void reload()} disabled={isLoading}>
              <Icon
                name="refresh"
                size="sm"
                className={isLoading ? "animate-spin text-muted-foreground" : "text-muted-foreground"}
              />
              Refresh
            </Button>
            <Button type="button" size="sm" onClick={() => setIsTestModalOpen(true)}>
              <Icon name="send" size="sm" data-icon="inline-start" />
              Send test request
            </Button>
          </div>
        )}
      </div>

      <div className="w-full overflow-hidden rounded-xl border bg-card shadow-xs">
        <div className="grid grid-cols-5 border-b bg-muted/20 px-6 py-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          <div>Time</div>
          <div>Channel</div>
          <div>Method</div>
          <div>Status</div>
          <div>Endpoint</div>
        </div>

        {isLoading && logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 py-20 text-sm text-muted-foreground">
            <span className="size-6 animate-spin rounded-full border-2 border-primary border-t-transparent" />
            <span>Loading request logs...</span>
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center gap-2 px-4 py-20 text-center">
            <Icon name="error_outline" size="lg" className="text-destructive" />
            <span className="text-sm font-semibold text-foreground">Failed to load logs</span>
            <span className="max-w-sm text-xs text-muted-foreground">{error}</span>
            <Button variant="outline" size="sm" onClick={() => void reload()} className="mt-2">
              Try again
            </Button>
          </div>
        ) : logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center px-4 py-24 text-center animate-in fade-in duration-300">
            <div className="mb-4 flex size-16 items-center justify-center rounded-full border bg-muted/60 text-muted-foreground">
              <Icon name="mail_outline" size="lg" className="text-2xl text-muted-foreground" />
            </div>
            <h3 className="mb-1 text-lg font-bold text-foreground">No requests yet</h3>
            <p className="mb-6 max-w-xs text-sm text-muted-foreground">
              Send your first test request to see gateway logs here.
            </p>
            <div className="mb-8 flex w-full max-w-xs flex-col items-center gap-3 sm:max-w-none sm:flex-row sm:justify-center">
              <Button type="button" onClick={() => setIsTestModalOpen(true)} className="h-10 w-full sm:w-auto">
                <Icon name="send" size="sm" data-icon="inline-start" />
                Send test request
              </Button>
              <Button type="button" variant="outline" onClick={() => setIsDocsOpen(true)} className="h-10 w-full sm:w-auto">
                <Icon name="menu_book" size="sm" data-icon="inline-start" />
                View API Docs
              </Button>
            </div>
            <div className="flex flex-col items-center gap-2">
              <span className="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
                Example endpoint
              </span>
              <code className="rounded-md border bg-muted/60 px-3 py-1.5 font-mono text-xs font-semibold text-foreground">
                {exampleEndpoint(channel)}
              </code>
            </div>
          </div>
        ) : (
          <div className="divide-y">
            {logs.map((log) => (
              <div
                key={log.id}
                className="grid grid-cols-5 items-center px-6 py-4 text-sm transition-colors hover:bg-muted/10"
              >
                <div className="font-medium text-muted-foreground">{formatTime(log.created_at)}</div>
                <div className="flex items-center gap-2">
                  <Icon name={channelIcon(log.channel_type)} size="sm" className="text-muted-foreground/80" />
                  <span className="font-semibold capitalize text-foreground">{log.channel_type}</span>
                </div>
                <div>
                  <span className="rounded border bg-muted/50 px-1.5 py-0.5 font-mono text-[11px] font-bold uppercase text-muted-foreground">
                    {log.http_method || "POST"}
                  </span>
                </div>
                <div>
                  <span className={`rounded-full px-2 py-0.5 text-xs ${statusBadgeClass(log.status_code)}`}>
                    {log.status_code}
                  </span>
                </div>
                <div className="truncate font-mono text-xs text-muted-foreground" title={log.endpoint}>
                  {log.endpoint || `/v1/${log.channel_type}`}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {logs.length > 0 && (
        <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
          <span>
            Showing latest {logs.length} logs for {channel ? `${channel} channel` : "all channels"}
          </span>
          <span>Total tracked: {total}</span>
        </div>
      )}

      {wid && (
        <SendTestModal
          workspaceId={wid}
          open={isTestModalOpen}
          onOpenChange={setIsTestModalOpen}
          onSent={() => void reload()}
          initialChannel={channel}
        />
      )}

      <Dialog open={isDocsOpen} onOpenChange={setIsDocsOpen}>
        <DialogContent className="max-h-[85vh] max-w-lg gap-0 overflow-hidden p-0">
          <DialogHeader className="border-b bg-muted/20 px-5 py-4 text-left">
            <DialogTitle>API Reference</DialogTitle>
            <DialogDescription>Integrate Message Gateway into your codebase.</DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-4 overflow-y-auto p-6 text-sm text-foreground">
            <p className="text-xs leading-relaxed text-muted-foreground">
              Attach workspace API credentials to each request. See backend usage docs for full reference.
            </p>
            <div className="flex flex-col gap-3">
              <div>
                <h4 className="mb-1 text-xs font-bold uppercase text-muted-foreground">Email</h4>
                <pre className="rounded-md border bg-muted/80 p-2.5 font-mono text-xs">{`POST /v1/email`}</pre>
              </div>
              <div>
                <h4 className="mb-1 text-xs font-bold uppercase text-muted-foreground">SMS</h4>
                <pre className="rounded-md border bg-muted/80 p-2.5 font-mono text-xs">{`POST /v1/sms`}</pre>
              </div>
            </div>
          </div>
          <DialogFooter className="border-t bg-muted/20 px-5 py-4">
            <Button type="button" onClick={() => setIsDocsOpen(false)}>
              Close
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
