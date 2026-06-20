import { useEffect, useState, useCallback, useRef } from "react"
import { useNavigate, useParams } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Spinner } from "@/components/ui/spinner"
import { ROUTES } from "@/core/router/routes"
import { buildInboxEventsUrl, deleteInboxEmail, fetchInboxEmails } from "../inbox.api"
import { useInboxKey } from "../hooks/use-inbox-key.hook"
import { EmailList } from "../components/email-list"
import { EmailContent } from "../components/email-content"
import type { InboxCredentials, StoredEmail } from "../inbox.types"

export function EmailInboxPage() {
  const navigate = useNavigate()
  const { wid } = useParams<{ wid: string }>()
  const { creds, isLoading: credsLoading, error: credsError, refreshCreds } = useInboxKey(wid)

  const [messages, setMessages] = useState<StoredEmail[]>([])
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isDeleting, setIsDeleting] = useState(false)
  const [isRetrying, setIsRetrying] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const authRetryRef = useRef(false)

  const fetchEmails = useCallback(
    async (activeCreds: InboxCredentials) => {
      if (!wid) return

      setError(null)
      setIsLoading(true)

      const result = await fetchInboxEmails(wid, activeCreds)
      if (!result.ok) {
        if (result.unauthorized && !authRetryRef.current) {
          authRetryRef.current = true
          const refreshed = await refreshCreds()
          if (refreshed) {
            const retryResult = await fetchInboxEmails(wid, refreshed)
            if (retryResult.ok) {
              setMessages(retryResult.items)
              if (retryResult.items.length > 0) {
                setSelectedMessageId((prev) => prev ?? retryResult.items[0].id)
              }
              setIsLoading(false)
              return
            }
            setError(retryResult.message)
            setIsLoading(false)
            return
          }
        }

        setError(result.message)
        setIsLoading(false)
        return
      }

      authRetryRef.current = false
      setMessages(result.items)
      if (result.items.length > 0) {
        setSelectedMessageId((prev) => prev ?? result.items[0].id)
      }
      setIsLoading(false)
    },
    [wid, refreshCreds]
  )

  useEffect(() => {
    authRetryRef.current = false
    if (creds) {
      Promise.resolve().then(() => {
        void fetchEmails(creds)
      })
    }
  }, [creds, fetchEmails])

  useEffect(() => {
    if (!wid || !creds) return

    const sseUrl = buildInboxEventsUrl(wid, creds)
    let eventSource: EventSource | null = null
    try {
      eventSource = new EventSource(sseUrl)

      eventSource.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data) as { type: string }
          if (
            parsed.type === "email_received" ||
            parsed.type === "email_deleted" ||
            parsed.type === "messages_cleared"
          ) {
            void fetchEmails(creds)
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
  }, [wid, creds, fetchEmails])

  const handleRetry = async () => {
    if (!wid) return

    setIsRetrying(true)
    setError(null)
    authRetryRef.current = false

    try {
      const nextCreds = (await refreshCreds()) ?? creds
      if (nextCreds) {
        await fetchEmails(nextCreds)
      }
    } finally {
      setIsRetrying(false)
    }
  }

  const handleDeleteMessage = async (messageId: string) => {
    if (!wid || !creds) return
    setIsDeleting(true)
    const result = await deleteInboxEmail(wid, messageId, creds)
    if (!result.ok) {
      alert(result.message)
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

  const activeMessage = messages.find((m) => m.id === selectedMessageId) || null
  const isBootstrapLoading = credsLoading || isRetrying || (isLoading && !!creds)
  const displayError = credsError || error

  return (
    <div className="flex flex-col gap-6 w-full h-[calc(100vh-160px)] min-h-[500px]">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between shrink-0">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            Email Inbox
          </h1>
          <p className="text-sm text-text-secondary">
            Read and test outbound emails captured by the memory provider.
          </p>
        </div>

        <div className="flex items-center gap-2.5 shrink-0">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => creds && void fetchEmails(creds)}
            disabled={isBootstrapLoading || !creds}
            className="h-10"
          >
            <Icon name="refresh" size="sm" data-icon="inline-start" className={isBootstrapLoading ? "animate-spin" : undefined} />
            Refresh
          </Button>
          <Button type="button" size="sm" onClick={handleManageTemplates} className="h-10">
            <Icon name="dashboard_customize" size="sm" data-icon="inline-start" />
            Manage Templates
          </Button>
        </div>
      </div>

      {displayError ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center p-6 text-center">
          <span className="text-sm font-semibold text-foreground">Could not connect to Simulated Inbox</span>
          <p className="text-xs text-text-secondary max-w-sm mt-1 mb-4">
            {displayError}
          </p>
          <Button variant="outline" size="sm" onClick={() => void handleRetry()} disabled={isBootstrapLoading}>
            {isBootstrapLoading ? "Reconnecting…" : "Try again"}
          </Button>
        </div>
      ) : isBootstrapLoading ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center gap-3">
          <Spinner size="lg" />
          <span className="text-sm text-text-secondary">Connecting to secure simulated inbox...</span>
        </div>
      ) : messages.length === 0 ? (
        <div className="flex-1 bg-card border border-border rounded-2xl flex flex-col items-center justify-center p-8 text-center">
          <div className="mb-4 flex size-16 items-center justify-center rounded-full border border-border bg-muted/40 text-text-tertiary">
            <Icon name="mail" size="lg" />
          </div>
          <h3 className="mb-1 text-lg font-bold text-foreground">Inbox is empty</h3>
          <p className="max-w-xs text-sm text-text-secondary">
            Send an email to `/v1/email` using this workspace API key to see it captured here in real-time.
          </p>
        </div>
      ) : (
        <div className="flex-1 bg-card border border-border rounded-2xl overflow-hidden flex flex-col md:flex-row shadow-sm min-h-0">
          <EmailList
            messages={messages}
            selectedMessageId={selectedMessageId}
            onSelectMessage={setSelectedMessageId}
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
