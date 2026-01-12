package config

// ChatConfig specific configuration for Chat providers
type ChatConfig struct {
	CommonConfig
	FromPhone  string
	WebhookURL string // e.g. Slack/Teams Webhook
}
