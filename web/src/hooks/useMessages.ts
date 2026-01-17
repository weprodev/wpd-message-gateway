import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useEffect, useCallback } from 'react'
import type { StoredEmail, StoredSMS, StoredPush, StoredChat, Stats } from '@/types/messages'

const API_BASE = '/api/v1'

// --- Fetch functions ---

async function fetchStats(): Promise<Stats> {
    const res = await fetch(`${API_BASE}/stats`)
    if (!res.ok) throw new Error('Failed to fetch stats')
    return res.json()
}

async function fetchEmails(): Promise<StoredEmail[]> {
    const res = await fetch(`${API_BASE}/emails`)
    if (!res.ok) throw new Error('Failed to fetch emails')
    return res.json()
}

async function fetchSMS(): Promise<StoredSMS[]> {
    const res = await fetch(`${API_BASE}/sms`)
    if (!res.ok) throw new Error('Failed to fetch SMS')
    return res.json()
}

async function fetchPush(): Promise<StoredPush[]> {
    const res = await fetch(`${API_BASE}/push`)
    if (!res.ok) throw new Error('Failed to fetch push notifications')
    return res.json()
}

async function fetchChat(): Promise<StoredChat[]> {
    const res = await fetch(`${API_BASE}/chat`)
    if (!res.ok) throw new Error('Failed to fetch chat messages')
    return res.json()
}

async function deleteMessage(type: string, id: string): Promise<void> {
    const res = await fetch(`${API_BASE}/${type}/${id}`, { method: 'DELETE' })
    if (!res.ok) throw new Error(`Failed to delete ${type}`)
}

async function clearAllMessages(): Promise<void> {
    const res = await fetch(`${API_BASE}/messages`, { method: 'DELETE' })
    if (!res.ok) throw new Error('Failed to clear messages')
}

// --- Hooks ---

export function useStats() {
    return useQuery({
        queryKey: ['stats'],
        queryFn: fetchStats,
        refetchInterval: 5000,
    })
}

export function useEmails() {
    return useQuery({
        queryKey: ['emails'],
        queryFn: fetchEmails,
        refetchInterval: 5000,
    })
}

export function useSMS() {
    return useQuery({
        queryKey: ['sms'],
        queryFn: fetchSMS,
        refetchInterval: 5000,
    })
}

export function usePush() {
    return useQuery({
        queryKey: ['push'],
        queryFn: fetchPush,
        refetchInterval: 5000,
    })
}

export function useChat() {
    return useQuery({
        queryKey: ['chat'],
        queryFn: fetchChat,
        refetchInterval: 5000,
    })
}

export function useDeleteMessage() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({ type, id }: { type: string; id: string }) => deleteMessage(type, id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['emails'] })
            queryClient.invalidateQueries({ queryKey: ['sms'] })
            queryClient.invalidateQueries({ queryKey: ['push'] })
            queryClient.invalidateQueries({ queryKey: ['chat'] })
            queryClient.invalidateQueries({ queryKey: ['stats'] })
        },
    })
}

export function useClearAll() {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: clearAllMessages,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['emails'] })
            queryClient.invalidateQueries({ queryKey: ['sms'] })
            queryClient.invalidateQueries({ queryKey: ['push'] })
            queryClient.invalidateQueries({ queryKey: ['chat'] })
            queryClient.invalidateQueries({ queryKey: ['stats'] })
        },
    })
}

// --- SSE Hook ---

export function useSSE() {
    const queryClient = useQueryClient()

    const handleEvent = useCallback((data: { type: string }) => {
        // Invalidate queries based on event type
        if (data.type.includes('email')) {
            queryClient.invalidateQueries({ queryKey: ['emails'] })
        } else if (data.type.includes('sms')) {
            queryClient.invalidateQueries({ queryKey: ['sms'] })
        } else if (data.type.includes('push')) {
            queryClient.invalidateQueries({ queryKey: ['push'] })
        } else if (data.type.includes('chat')) {
            queryClient.invalidateQueries({ queryKey: ['chat'] })
        } else if (data.type === 'messages_cleared') {
            queryClient.invalidateQueries({ queryKey: ['emails'] })
            queryClient.invalidateQueries({ queryKey: ['sms'] })
            queryClient.invalidateQueries({ queryKey: ['push'] })
            queryClient.invalidateQueries({ queryKey: ['chat'] })
        }
        queryClient.invalidateQueries({ queryKey: ['stats'] })
    }, [queryClient])

    useEffect(() => {
        const eventSource = new EventSource(`${API_BASE}/events`)

        eventSource.addEventListener('message', (event) => {
            try {
                const data = JSON.parse(event.data)
                handleEvent(data)
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

        return () => {
            eventSource.close()
        }
    }, [handleEvent])
}
