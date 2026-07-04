import { useEffect, useState } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Spinner } from "@/components/ui/spinner"
import { ROUTES } from "@/core/router/routes"
import {
  buildInboxEventsUrl,
  deleteInboxEmail,
  fetchInboxEmailById,
} from "../inbox.api"
import { EmailList } from "../components/email-list"
import { EmailContent } from "../components/email-content"
import { useInboxEmails } from "../hooks/use-inbox-emails.hook"

export function EmailInboxPage() {
  const navigate = useNavigate()
  const { wid } = useParams<{ wid: string }>()
  const [searchParams, setSearchParams] = useSearchParams()
  const deepLinkMessageId = searchParams.get("message")
  const [isDeleting, setIsDeleting] = useState(false)

  const {
    messages,
    selectedMessageId,
    setSelectedMessageId,
    isLoading,
    isLoadingMore,
    error,
    hasMore,
    reload,
    loadMore,
    upsertMessage,
    removeMessage,
    clearMessages,
  } = useInboxEmails(wid)

  useEffect(() => {
    if (!deepLinkMessageId) return

    setSelectedMessageId(deepLinkMessageId)
    setSearchParams((params) => {
      const next = new URLSearchParams(params)
      next.delete("message")
      return next
    }, { replace: true })
  }, [deepLinkMessageId, setSearchParams, setSelectedMessageId])

  useEffect(() => {
    if (!wid) return

    const sseUrl = buildInboxEventsUrl(wid)
    let eventSource: EventSource | null = null
    try {
      eventSource = new EventSource(sseUrl)

      eventSource.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data) as {
            type: string
            data?: { id?: string } | string | null
          }

          if (parsed.type === "email_received") {
            const id =
              typeof parsed.data === "object" && parsed.data !== null
                ? parsed.data.id
                : undefined
            if (!id) return
            void fetchInboxEmailById(wid, id).then((result) => {
              if (!result.ok) return
              upsertMessage(result.item)
              setSelectedMessageId((current) => current ?? id)
            })
            return
          }

          if (parsed.type === "email_deleted") {
            const id = typeof parsed.data === "string" ? parsed.data : undefined
            if (!id) return
            removeMessage(id)
            return
          }

          if (parsed.type === "messages_cleared") {
            clearMessages()
          }
        } catch (err) {
          console.error("Failed to parse SSE event data:", err)
        }
      }
    } catch (err) {
      console.error("SSE connection failed:", err)
    }

    return () => {
      eventSource?.close()
    }
  }, [wid, upsertMessage, removeMessage, clearMessages, setSelectedMessageId])

  const handleDeleteMessage = async (messageId: string) => {
    if (!wid) return
    setIsDeleting(true)
    const result = await deleteInboxEmail(wid, messageId)
    if (!result.ok) {
      alert(result.message ?? "Failed to delete email")
      setIsDeleting(false)
      return
    }

    removeMessage(messageId)
    setIsDeleting(false)
  }

  const handleManageTemplates = () => {
    if (wid) {
      navigate(ROUTES.workspace.emailTemplates(wid))
    }
  }

  const activeMessage = messages.find((m) => m.id === selectedMessageId) ?? null
  const showLoading = Boolean(wid) && isLoading

  return (
    <div className="flex flex-col gap-6 w-full h-[calc(100vh-160px)] min-h-[500px]">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between shrink-0">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Email Inbox
          </h1>
          <p className="text-sm text-text-secondary">
            Outbound emails captured by the gateway for review and testing.
          </p>
        </div>

        <div className="flex items-center gap-2.5 shrink-0">
          <Button type="button" variant="outline" size="sm" onClick={() => void reload()} disabled={showLoading} className="h-10">
            <Icon name="refresh" size="sm" data-icon="inline-start" className={showLoading ? "animate-spin" : undefined} />
            Refresh
          </Button>
          <Button type="button" size="sm" onClick={handleManageTemplates} className="h-10">
            <Icon name="dashboard_customize" size="sm" data-icon="inline-start" />
            Manage Templates
          </Button>
        </div>
      </div>

      {error ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center p-6 text-center">
          <span className="text-sm font-semibold text-foreground">Could not connect to Simulated Inbox</span>
          <p className="text-xs text-text-secondary max-w-sm mt-1 mb-4">
            {error}
          </p>
          <Button variant="outline" size="sm" onClick={() => void reload()}>
            Try again
          </Button>
        </div>
      ) : showLoading ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center gap-3">
          <Spinner size="lg" />
          <span className="text-sm text-text-secondary">Loading simulated inbox...</span>
        </div>
      ) : messages.length === 0 ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center p-8 text-center">
          <div className="mb-4 flex size-16 items-center justify-center rounded-full border border-border bg-muted/40 text-text-tertiary">
            <Icon name="mail" size="lg" />
          </div>
          <h3 className="mb-1 text-lg font-bold text-foreground">Inbox is empty</h3>
          <p className="max-w-xs text-sm text-text-secondary">
            POST to `/v1/email` with your API key to see captured messages here in real-time.
          </p>
        </div>
      ) : (
        <div className="flex-1 bg-card border border-border rounded-2xl overflow-hidden flex flex-col md:flex-row shadow-sm min-h-0">
          <EmailList
            messages={messages}
            selectedMessageId={selectedMessageId}
            onSelectMessage={setSelectedMessageId}
            hasMore={hasMore}
            isLoadingMore={isLoadingMore}
            onLoadMore={() => void loadMore()}
          />
          {activeMessage ? (
            <EmailContent
              message={activeMessage}
              onDelete={handleDeleteMessage}
              isDeleting={isDeleting}
            />
          ) : (
            <div className="flex-1 flex items-center justify-center text-sm text-text-tertiary bg-muted/5">
              Select an email to read its contents.
            </div>
          )}
        </div>
      )}
    </div>
  )
}
