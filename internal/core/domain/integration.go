package domain

import "time"

// Integration lifecycle statuses (integrations.status column).
const (
	IntegrationStatusConnected    = "connected"
	IntegrationStatusDisconnected = "disconnected"
)

// Integration holds provider credentials for a workspace channel.
type Integration struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id"`
	ChannelType  string    `json:"channel_type"`
	ProviderName string    `json:"provider_name"`
	Config       []byte    `json:"-"` // decrypted JSON
	Status       string    `json:"status"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
