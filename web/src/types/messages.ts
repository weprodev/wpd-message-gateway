// TypeScript types for DevBox API messages

export interface StoredEmail {
    id: string
    created_at: string
    email: {
        from: string
        from_name: string
        to: string[]
        cc?: string[]
        bcc?: string[]
        reply_to?: string
        subject: string
        html: string
        plain_text?: string
        attachments?: Attachment[]
        headers?: Record<string, string>
    }
}

export interface Attachment {
    filename: string
    content_type: string
    data: string
}

export interface StoredSMS {
    id: string
    created_at: string
    sms: {
        from: string
        to: string[]
        message: string
    }
}

export interface StoredPush {
    id: string
    created_at: string
    push: {
        device_tokens: string[]
        title: string
        body: string
        data?: Record<string, string>
        badge?: number
        sound?: string
    }
}

export interface StoredChat {
    id: string
    created_at: string
    chat: {
        from: string
        to: string[]
        message: string
        template_id?: string
        template_params?: string[]
        media_url?: string
        media_type?: string
        buttons?: ChatButton[]
        reply_to_id?: string
        metadata?: Record<string, string>
    }
}

export interface ChatButton {
    id: string
    text: string
    url?: string
    phone?: string
}

export interface Stats {
    emails: number
    sms: number
    push: number
    chat: number
    total: number
}

export type MessageType = 'email' | 'sms' | 'push' | 'chat'
