import { useState } from "react"
import { useNavigate, useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { SearchInput } from "@/components/ui/search-input"
import { ROUTES } from "@/core/router/routes"
import { GetStartedDialog } from "../components/get-started-dialog"
import { useInboxLogs } from "../hooks/use-inbox-logs.hook"
import type { MessageChannel, LogRow } from "../inbox.types"
import { hasInboxContent, isEmailInboxLink } from "../log-display.utils"
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
  const navigate = useNavigate()
  const { wid } = useParams<{ wid: string }>()
  const { logs, total, isLoading, error, reload } = useInboxLogs(wid, channel)
  const [isGetStartedOpen, setIsGetStartedOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")

  const filteredLogs = logs.filter((log) => {
    const q = searchQuery.toLowerCase()
    return (
      (log.endpoint ?? "").toLowerCase().includes(q) ||
      (log.http_method ?? "").toLowerCase().includes(q) ||
      String(log.status_code ?? "").includes(q) ||
      (log.source_name ?? "").toLowerCase().includes(q) ||
      (log.provider_name ?? "").toLowerCase().includes(q)
    )
  })

  const dateGroups = groupLogsByDate(filteredLogs)

  const handleLogClick = (log: LogRow) => {
    if (!wid || !isEmailInboxLink(log) || !log.inbox_message_id) return
    navigate(ROUTES.workspace.emailWithMessage(wid, log.inbox_message_id))
  }

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

          <Button type="button" size="sm" onClick={() => setIsGetStartedOpen(true)} className="h-10">
            <Icon name="code" size="sm" data-icon="inline-start" />
            Get started
          </Button>
        </div>
      </div>

      <div className="overflow-hidden rounded-2xl border border-border bg-card shadow-card">
        <div className="grid grid-cols-12 bg-muted/30 px-6 py-3 border-b border-border text-xs font-semibold uppercase tracking-wider text-text-secondary">
          <div className="col-span-2">Time</div>
          <div className="col-span-2">Channel</div>
          <div className="col-span-2">Source</div>
          <div className="col-span-1">Method</div>
          <div className="col-span-1">Status</div>
          <div className="col-span-1">Provider</div>
          <div className="col-span-1">Retained</div>
          <div className="col-span-2">Endpoint</div>
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
              Send your first request via the gateway API to see logs here.
            </p>
            <Button type="button" onClick={() => setIsGetStartedOpen(true)}>
              <Icon name="code" size="sm" data-icon="inline-start" />
              Get started
            </Button>
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
                  {group.logs.map((log) => {
                    const clickable = isEmailInboxLink(log)
                    return (
                      <div
                        key={log.id}
                        role={clickable ? "button" : undefined}
                        tabIndex={clickable ? 0 : undefined}
                        onClick={() => handleLogClick(log)}
                        onKeyDown={(e) => {
                          if (!clickable) return
                          if (e.key === "Enter" || e.key === " ") {
                            e.preventDefault()
                            handleLogClick(log)
                          }
                        }}
                        className={`grid grid-cols-12 items-center px-6 py-3.5 text-sm transition-colors hover:bg-muted/10 ${
                          clickable ? "cursor-pointer hover:bg-primary-brand/5" : ""
                        }`}
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

                        <div className="col-span-2 text-text-secondary truncate pr-4 font-medium" title={log.source_name || "API Request"}>
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

                        <div className="col-span-1 truncate text-xs font-medium text-foreground" title={log.provider_name}>
                          {log.provider_name || "—"}
                        </div>

                        <div className="col-span-1">
                          {hasInboxContent(log) ? (
                            <span className="rounded-full bg-emerald-500/10 px-2 py-0.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400">
                              Yes
                            </span>
                          ) : (
                            <span className="text-xs text-text-tertiary">No</span>
                          )}
                        </div>

                        <div className="col-span-2 truncate font-mono text-xs text-text-tertiary pr-2" title={log.endpoint}>
                          {log.endpoint || `/v1/${log.channel_type}`}
                        </div>
                      </div>
                    )
                  })}
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

      <GetStartedDialog open={isGetStartedOpen} onOpenChange={setIsGetStartedOpen} workspaceId={wid} />
    </div>
  )
}
