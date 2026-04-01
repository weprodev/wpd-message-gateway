package domain

import "time"

// MessageRequestLog records a gateway API call for monitoring.
type MessageRequestLog struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id"`
	APIKeyID     string    `json:"api_key_id,omitempty"`
	ChannelType  string    `json:"channel_type"`
	HTTPMethod   string    `json:"http_method"`
	StatusCode   int       `json:"status_code"`
	Endpoint     string    `json:"endpoint"`
	ProviderName string    `json:"provider_name,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	DurationMs   int       `json:"duration_ms,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
