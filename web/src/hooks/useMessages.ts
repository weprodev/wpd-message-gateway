import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useEffect, useCallback } from 'react'
import type { StoredEmail, StoredSMS, StoredPush, StoredChat, Stats } from '@/types/messages'
import { API_BASE, REFETCH_INTERVAL_MS } from '@/lib/constants'

// --- Generic Fetcher ---

async function fetchResource<T>(endpoint: string, errorMessage: string): Promise<T> {
  const response = await fetch(`${API_BASE}/${endpoint}`)
  if (!response.ok) throw new Error(errorMessage)
  return response.json()
}

// --- Query Keys ---

const QUERY_KEYS = {
  stats: ['stats'],
  emails: ['emails'],
  sms: ['sms'],
  push: ['push'],
  chat: ['chat'],
} as const

const ALL_MESSAGE_KEYS = [
  QUERY_KEYS.emails,
  QUERY_KEYS.sms,
  QUERY_KEYS.push,
  QUERY_KEYS.chat,
  QUERY_KEYS.stats,
]

// --- Hooks ---

export function useStats() {
  return useQuery({
    queryKey: QUERY_KEYS.stats,
    queryFn: () => fetchResource<Stats>('stats', 'Failed to fetch stats'),
    refetchInterval: REFETCH_INTERVAL_MS,
  })
}

export function useEmails() {
  return useQuery({
    queryKey: QUERY_KEYS.emails,
    queryFn: () => fetchResource<StoredEmail[]>('emails', 'Failed to fetch emails'),
    refetchInterval: REFETCH_INTERVAL_MS,
  })
}

export function useSMS() {
  return useQuery({
    queryKey: QUERY_KEYS.sms,
    queryFn: () => fetchResource<StoredSMS[]>('sms', 'Failed to fetch SMS'),
    refetchInterval: REFETCH_INTERVAL_MS,
  })
}

export function usePush() {
  return useQuery({
    queryKey: QUERY_KEYS.push,
    queryFn: () => fetchResource<StoredPush[]>('push', 'Failed to fetch push notifications'),
    refetchInterval: REFETCH_INTERVAL_MS,
  })
}

export function useChat() {
  return useQuery({
    queryKey: QUERY_KEYS.chat,
    queryFn: () => fetchResource<StoredChat[]>('chat', 'Failed to fetch chat messages'),
    refetchInterval: REFETCH_INTERVAL_MS,
  })
}

export function useDeleteMessage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ type, id }: { type: string; id: string }) => {
      const response = await fetch(`${API_BASE}/${type}/${id}`, { method: 'DELETE' })
      if (!response.ok) throw new Error(`Failed to delete ${type}`)
    },
    onSuccess: () => invalidateAllQueries(queryClient),
  })
}

export function useClearAll() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async () => {
      const response = await fetch(`${API_BASE}/messages`, { method: 'DELETE' })
      if (!response.ok) throw new Error('Failed to clear messages')
    },
    onSuccess: () => invalidateAllQueries(queryClient),
  })
}

// --- SSE Hook ---

export function useSSE() {
  const queryClient = useQueryClient()

  const handleEvent = useCallback((data: { type: string }) => {
    const typeToKeyMap: Record<string, readonly string[]> = {
      email: QUERY_KEYS.emails,
      sms: QUERY_KEYS.sms,
      push: QUERY_KEYS.push,
      chat: QUERY_KEYS.chat,
    }

    if (data.type === 'messages_cleared') {
      invalidateAllQueries(queryClient)
      return
    }

    for (const [keyword, queryKey] of Object.entries(typeToKeyMap)) {
      if (data.type.includes(keyword)) {
        queryClient.invalidateQueries({ queryKey })
        break
      }
    }

    queryClient.invalidateQueries({ queryKey: QUERY_KEYS.stats })
  }, [queryClient])

  useEffect(() => {
    const eventSource = new EventSource(`${API_BASE}/events`)

    eventSource.addEventListener('message', (event) => {
      try {
        handleEvent(JSON.parse(event.data))
      } catch {
        // Ignore parse errors
      }
    })

    eventSource.addEventListener('connected', () => {
      console.log('SSE connected')
    })

    eventSource.onerror = () => {
      console.log('SSE error, will reconnect...')
    }

    return () => eventSource.close()
  }, [handleEvent])
}

// --- Helpers ---

function invalidateAllQueries(queryClient: ReturnType<typeof useQueryClient>) {
  ALL_MESSAGE_KEYS.forEach(key => {
    queryClient.invalidateQueries({ queryKey: key })
  })
}
