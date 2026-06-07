export type IntegrationChannel = "email" | "sms" | "push" | "chat"

export interface Integration {
  id: string
  workspace_id: string
  channel_type: IntegrationChannel
  provider_name: string
  config: Record<string, unknown>
  status: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface ProviderCatalogItem {
  id: string
  name: string
  description: string
  icon: string
  category: IntegrationChannel
  isAvailable: boolean
  isComingSoon?: boolean
}

export const PROVIDER_CATALOG: readonly ProviderCatalogItem[] = [
  {
    id: "mailgun",
    name: "Mailgun",
    description: "Send transactional and marketing emails at scale",
    icon: "📧",
    category: "email",
    isAvailable: true,
  },
  {
    id: "sendgrid",
    name: "SendGrid",
    description: "Email delivery and marketing platform",
    icon: "✉️",
    category: "email",
    isAvailable: false,
    isComingSoon: true,
  },
  {
    id: "twilio",
    name: "Twilio",
    description: "SMS and voice communication APIs",
    icon: "💬",
    category: "sms",
    isAvailable: false,
    isComingSoon: true,
  },
  {
    id: "firebase",
    name: "Firebase Cloud Messaging",
    description: "Cross-platform push notification service",
    icon: "🔔",
    category: "push",
    isAvailable: false,
    isComingSoon: true,
  },
  {
    id: "telegram",
    name: "Telegram",
    description: "Send messages via Telegram Bot API",
    icon: "✈️",
    category: "chat",
    isAvailable: true,
  },
  {
    id: "slack",
    name: "Slack",
    description: "Team messaging and notifications",
    icon: "💼",
    category: "chat",
    isAvailable: false,
    isComingSoon: true,
  },
] as const
