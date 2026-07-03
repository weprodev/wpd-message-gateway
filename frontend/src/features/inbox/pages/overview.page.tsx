import { useState } from "react"
import { useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { SearchInput } from "@/components/ui/search-input"
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
import type { MessageChannel, LogRow } from "../inbox.types"
import { PageHeader } from "@/shared/components/page-header"

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

function getChannelIconName(channelType: string) {
  switch (channelType.toLowerCase()) {
    case "email":
      return "mail"
    case "sms":
      return "chat"
    case "push":
      return "notifications"
    case "chat":
      return "forum"
    default:
      return "mail"
  }
}

function getStatusStyle(code: number) {
  if (code >= 200 && code < 300) {
    return "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-semibold"
  }
  if (code >= 400) {
    return "bg-rose-500/10 text-rose-600 dark:text-rose-400 font-semibold"
  }
  return "bg-amber-500/10 text-amber-600 dark:text-amber-400 font-semibold"
}

function groupLogsByDate(logs: LogRow[]) {
  const groups: { date: string; logs: LogRow[] }[] = []
  logs.forEach((log) => {
    const d = new Date(log.created_at)
    const dateStr = d.toLocaleDateString([], {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    })
    let group = groups.find((g) => g.date === dateStr)
    if (!group) {
      group = { date: dateStr, logs: [] }
      groups.push(group)
    }
    group.logs.push(log)
  })
  return groups
}

export function OverviewPage({ channel }: OverviewPageProps) {
  const { wid } = useParams<{ wid: string }>()
  const { logs, total, isLoading, error, reload } = useInboxLogs(wid, channel)
  const [isTestModalOpen, setIsTestModalOpen] = useState(false)
  const [isDocsOpen, setIsDocsOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")

  const filteredLogs = logs.filter((log) => {
    const q = searchQuery.toLowerCase()
    return (
      (log.endpoint ?? "").toLowerCase().includes(q) ||
      (log.http_method ?? "").toLowerCase().includes(q) ||
      String(log.status_code ?? "").includes(q) ||
      (log.source_name ?? "").toLowerCase().includes(q)
    )
  })

  const dateGroups = groupLogsByDate(filteredLogs)

  return (
    <div className="flex flex-col gap-6 w-full">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <PageHeader
          title={channel ? `${channel} logs` : "Message logs"}
          description="Monitor API requests across communication channels."
        />

        <div className="flex items-center gap-3 shrink-0">
          <SearchInput
            placeholder="Search logs..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-64 shadow-xs"
          />

          <Button type="button" variant="outline" size="sm" onClick={() => void reload()} disabled={isLoading} className="h-10">
            <Icon name="refresh" size="sm" data-icon="inline-start" className={isLoading ? "animate-spin" : undefined} />
            Refresh
          </Button>

          <Button type="button" size="sm" onClick={() => setIsTestModalOpen(true)} className="h-10">
            <Icon name="send" size="sm" data-icon="inline-start" />
            Send test request
          </Button>
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-border bg-card shadow-card">
        <div className="grid grid-cols-12 bg-muted/30 px-6 py-3 border-b border-border text-xs font-semibold uppercase tracking-wider text-text-secondary">
          <div className="col-span-2">Time</div>
          <div className="col-span-2">Channel</div>
          <div className="col-span-3">Source</div>
          <div className="col-span-1">Method</div>
          <div className="col-span-1">Status</div>
          <div className="col-span-3">Endpoint</div>
        </div>

        {isLoading && logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-3 py-20 text-sm text-text-secondary">
            <span className="size-6 animate-spin rounded-full border-2 border-primary-brand border-t-transparent" />
            <span>Loading request logs...</span>
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center gap-2 px-4 py-20 text-center">
            <span className="text-sm font-semibold text-foreground">Failed to load logs</span>
            <span className="max-w-sm text-xs text-text-tertiary">{error}</span>
            <Button variant="outline" size="sm" onClick={() => void reload()} className="mt-2">
              Try again
            </Button>
          </div>
        ) : logs.length === 0 ? (
          <div className="flex flex-col items-center justify-center px-4 py-24 text-center">
            <div className="mb-4 flex size-16 items-center justify-center rounded-full border border-border bg-muted/40 text-text-tertiary">
              <Icon name="mail" size="lg" />
            </div>
            <h3 className="mb-1 text-lg font-bold text-foreground">No requests yet</h3>
            <p className="mb-6 max-w-xs text-sm text-text-secondary">
              Send your first test request to see gateway logs here.
            </p>
            <div className="flex gap-3 justify-center">
              <Button type="button" onClick={() => setIsTestModalOpen(true)}>
                <Icon name="send" size="sm" data-icon="inline-start" />
                Send test request
              </Button>
              <Button type="button" variant="outline" onClick={() => setIsDocsOpen(true)}>
                View API Docs
              </Button>
            </div>
          </div>
        ) : filteredLogs.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-sm text-text-tertiary">
            No logs matched your search filters.
          </div>
        ) : (
          <div className="divide-y divide-border">
            {dateGroups.map((group) => (
              <div key={group.date} className="flex flex-col">
                <div className="bg-muted/10 px-6 py-1.5 border-b border-border text-xs font-semibold text-text-secondary">
                  {group.date}
                </div>

                <div className="divide-y divide-border">
                  {group.logs.map((log) => (
                    <div
                      key={log.id}
                      className="grid grid-cols-12 items-center px-6 py-3.5 text-sm transition-colors hover:bg-muted/10"
                    >
                      <div className="col-span-2 font-medium text-text-secondary">
                        {formatTime(log.created_at)}
                      </div>

                      <div className="col-span-2 flex items-center gap-2">
                        <div className="bg-secondary rounded-md size-7 flex items-center justify-center border border-border shrink-0">
                          <Icon name={getChannelIconName(log.channel_type)} size="sm" className="text-text-secondary" />
                        </div>
                        <span className="font-medium capitalize text-foreground">{log.channel_type}</span>
                      </div>

                      <div className="col-span-3 text-text-secondary truncate pr-4 font-medium" title={log.source_name || "API Request"}>
                        {log.source_name || "API Request"}
                      </div>

                      <div className="col-span-1">
                        <span className="font-mono text-xs font-bold text-foreground">
                          {log.http_method || "POST"}
                        </span>
                      </div>

                      <div className="col-span-1">
                        <span className={`rounded-full px-2.5 py-0.5 text-xs font-mono ${getStatusStyle(log.status_code)}`}>
                          {log.status_code}
                        </span>
                      </div>

                      <div className="col-span-3 truncate font-mono text-xs text-text-tertiary pr-2" title={log.endpoint}>
                        {log.endpoint || `/v1/${log.channel_type}`}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {logs.length > 0 && (
        <div className="flex items-center justify-between px-1 text-xs text-text-secondary">
          <span>
            Showing latest {filteredLogs.length} logs for {channel ? `${channel} channel` : "all channels"}
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
            <p className="text-xs leading-relaxed text-text-secondary">
              Attach workspace API credentials to each request. See backend usage docs for full reference.
            </p>
            <div className="flex flex-col gap-3">
              <div>
                <h4 className="mb-1 text-xs font-bold uppercase text-text-secondary">Email</h4>
                <pre className="rounded-md border bg-muted/80 p-2.5 font-mono text-xs">{`POST /v1/email`}</pre>
              </div>
              <div>
                <h4 className="mb-1 text-xs font-bold uppercase text-text-secondary">SMS</h4>
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
