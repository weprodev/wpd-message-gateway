import { useEffect, useState, useCallback } from "react"
import { useNavigate, useParams, useSearchParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Spinner } from "@/components/ui/spinner"
import { ROUTES } from "@/core/router/routes"
import {
  buildInboxEventsUrl,
  deleteInboxEmail,
  fetchInboxEmailById,
  fetchInboxEmails,
} from "../inbox.api"
import { EmailList } from "../components/email-list"
import { EmailContent } from "../components/email-content"
import type { StoredEmail } from "../inbox.types"

const PAGE_SIZE = 50

export function EmailInboxPage() {
  const navigate = useNavigate()
  const { wid } = useParams<{ wid: string }>()
  const [searchParams, setSearchParams] = useSearchParams()
  const deepLinkMessageId = searchParams.get("message")

  const [messages, setMessages] = useState<StoredEmail[]>([])
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [nextCursor, setNextCursor] = useState<string | undefined>()
  const [hasMore, setHasMore] = useState(false)

  const fetchEmails = useCallback(async () => {
    if (!wid) return
    setError(null)
    setIsLoading(true)

    const result = await fetchInboxEmails(wid, { limit: PAGE_SIZE })
    if (!result.ok) {
      setError(result.message ?? "Failed to load emails")
      setIsLoading(false)
      return
    }

    setMessages(result.page.items ?? [])
    setNextCursor(result.page.next_cursor)
    setHasMore(result.page.has_more)
    if ((result.page.items ?? []).length > 0) {
      setSelectedMessageId((prev) => prev ?? (result.page.items ?? [])[0].id)
    }
    setIsLoading(false)
  }, [wid])

  const loadMoreEmails = useCallback(async () => {
    if (!wid || !nextCursor || isLoadingMore) return
    setIsLoadingMore(true)

    const result = await fetchInboxEmails(wid, { limit: PAGE_SIZE, cursor: nextCursor })
    if (!result.ok) {
      setError(result.message ?? "Failed to load more emails")
      setIsLoadingMore(false)
      return
    }

    setMessages((prev) => [...(prev ?? []), ...(result.page.items ?? [])])
    setNextCursor(result.page.next_cursor)
    setHasMore(result.page.has_more)
    setIsLoadingMore(false)
  }, [wid, nextCursor, isLoadingMore])

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      if (!wid) {
        setIsLoading(false)
        return
      }
      setError(null)
      setIsLoading(true)

      const result = await fetchInboxEmails(wid, { limit: PAGE_SIZE })
      if (cancelled) return

      if (!result.ok) {
        setError(result.message ?? "Failed to load emails")
        setIsLoading(false)
        return
      }

      setMessages(result.page.items ?? [])
      setNextCursor(result.page.next_cursor)
      setHasMore(result.page.has_more)
      if ((result.page.items ?? []).length > 0) {
        setSelectedMessageId((prev) => prev ?? (result.page.items ?? [])[0].id)
      }
      setIsLoading(false)
    })()

    return () => {
      cancelled = true
    }
  }, [wid])

  useEffect(() => {
    if (!wid || !deepLinkMessageId) return

    let cancelled = false
    ;(async () => {
      const result = await fetchInboxEmailById(wid, deepLinkMessageId)
      if (cancelled) return

      if (result.ok) {
        setMessages((prev) => {
          if (prev.some((m) => m.id === result.item.id)) return prev
          return [result.item, ...prev]
        })
        setSelectedMessageId(deepLinkMessageId)
      }

      setSearchParams((params) => {
        const next = new URLSearchParams(params)
        next.delete("message")
        return next
      }, { replace: true })
    })()

    return () => {
      cancelled = true
    }
  }, [wid, deepLinkMessageId, setSearchParams])

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
              setMessages((prev) => {
                if (prev.some((m) => m.id === result.item.id)) return prev
                return [result.item, ...prev]
              })
              setSelectedMessageId((current) => current ?? id)
            })
            return
          }

          if (parsed.type === "email_deleted") {
            const id = typeof parsed.data === "string" ? parsed.data : undefined
            if (!id) return
            setMessages((prev) => {
              const remaining = prev.filter((m) => m.id !== id)
              setSelectedMessageId((current) => {
                if (current === id) {
                  return remaining.length > 0 ? remaining[0].id : null
                }
                return current
              })
              return remaining
            })
            return
          }

          if (parsed.type === "messages_cleared") {
            setMessages([])
            setSelectedMessageId(null)
            setNextCursor(undefined)
            setHasMore(false)
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
  }, [wid])

  const handleDeleteMessage = async (messageId: string) => {
    if (!wid) return
    setIsDeleting(true)
    const result = await deleteInboxEmail(wid, messageId)
    if (!result.ok) {
      alert(result.message ?? "Failed to delete email")
      setIsDeleting(false)
      return
    }

    setMessages((prev) => {
      const remaining = prev.filter((m) => m.id !== messageId)
      setSelectedMessageId((current) => {
        if (current === messageId) {
          return remaining.length > 0 ? remaining[0].id : null
        }
        return current
      })
      return remaining
    })
    setIsDeleting(false)
  }

  const handleManageTemplates = () => {
    if (wid) {
      navigate(ROUTES.workspace.emailTemplates(wid))
    }
  }

  const activeMessage = (messages ?? []).find((m) => m.id === selectedMessageId) ?? null

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
          <Button type="button" variant="outline" size="sm" onClick={() => void fetchEmails()} disabled={isLoading} className="h-10">
            <Icon name="refresh" size="sm" data-icon="inline-start" className={isLoading ? "animate-spin" : undefined} />
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
          <Button variant="outline" size="sm" onClick={() => void fetchEmails()}>
            Try again
          </Button>
        </div>
      ) : isLoading ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center gap-3">
          <Spinner size="lg" />
          <span className="text-sm text-text-secondary">Loading simulated inbox...</span>
        </div>
      ) : (messages ?? []).length === 0 ? (
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
            messages={messages ?? []}
            selectedMessageId={selectedMessageId}
            onSelectMessage={setSelectedMessageId}
            hasMore={hasMore}
            isLoadingMore={isLoadingMore}
            onLoadMore={() => void loadMoreEmails()}
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
