import { useCallback, useEffect, useState } from "react"

import { fetchInboxEmailById, fetchInboxEmails } from "../api/inbox.api"
import type { StoredEmail } from "../inbox.types"

const PAGE_SIZE = 50

export function useInboxEmails(workspaceId: string | undefined) {
  const [messages, setMessages] = useState<StoredEmail[]>([])
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isLoadingMore, setIsLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [nextCursor, setNextCursor] = useState<string | undefined>()
  const [hasMore, setHasMore] = useState(false)

  const upsertMessage = useCallback((message: StoredEmail) => {
    setMessages((prev) => {
      const index = prev.findIndex((m) => m.id === message.id)
      if (index === -1) return [message, ...prev]
      const next = [...prev]
      next[index] = message
      return next
    })
  }, [])

  const applyPage = useCallback((result: Awaited<ReturnType<typeof fetchInboxEmails>>) => {
    setIsLoading(false)
    if (!result.ok) {
      setError(result.message ?? "Failed to load emails")
      return
    }

    setError(null)
    setMessages(result.page.items ?? [])
    setNextCursor(result.page.next_cursor)
    setHasMore(result.page.has_more)
    if ((result.page.items ?? []).length > 0) {
      setSelectedMessageId((prev) => prev ?? (result.page.items ?? [])[0].id)
    }
  }, [])

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      if (!workspaceId) return
      setIsLoading(true)
      setError(null)
      const result = await fetchInboxEmails(workspaceId, { limit: PAGE_SIZE })
      if (cancelled) return
      applyPage(result)
    })()

    return () => {
      cancelled = true
    }
  }, [workspaceId, applyPage])

  useEffect(() => {
    if (!workspaceId || !selectedMessageId) return

    let cancelled = false
    void fetchInboxEmailById(workspaceId, selectedMessageId).then((result) => {
      if (cancelled || !result.ok) return
      upsertMessage(result.item)
    })

    return () => {
      cancelled = true
    }
  }, [workspaceId, selectedMessageId, upsertMessage])

  const reload = useCallback(async () => {
    if (!workspaceId) return
    setIsLoading(true)
    setError(null)
    const result = await fetchInboxEmails(workspaceId, { limit: PAGE_SIZE })
    applyPage(result)
  }, [workspaceId, applyPage])

  const loadMore = useCallback(async () => {
    if (!workspaceId || !nextCursor || isLoadingMore) return
    setIsLoadingMore(true)

    const result = await fetchInboxEmails(workspaceId, { limit: PAGE_SIZE, cursor: nextCursor })
    setIsLoadingMore(false)
    if (!result.ok) {
      setError(result.message ?? "Failed to load more emails")
      return
    }

    setMessages((prev) => [...prev, ...(result.page.items ?? [])])
    setNextCursor(result.page.next_cursor)
    setHasMore(result.page.has_more)
  }, [workspaceId, nextCursor, isLoadingMore])

  const removeMessage = useCallback((messageId: string) => {
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
  }, [])

  const clearMessages = useCallback(() => {
    setMessages([])
    setSelectedMessageId(null)
    setNextCursor(undefined)
    setHasMore(false)
  }, [])

  return {
    messages,
    selectedMessageId,
    setSelectedMessageId,
    isLoading,
    isLoadingMore,
    error,
    nextCursor,
    hasMore,
    reload,
    loadMore,
    upsertMessage,
    removeMessage,
    clearMessages,
  }
}
