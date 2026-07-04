package domain

import "time"

// MessageRequestLog records a gateway API call for monitoring.
type MessageRequestLog struct {
	ID             string
	WorkspaceID    string
	APIKeyID       string
	ChannelType    string
	HTTPMethod     string
	StatusCode     int
	Endpoint       string
	ProviderName   string
	RequestID      string
	DurationMs     int
	ErrorMessage   string
	InboxMessageID string

	CreatedAt time.Time

	Payload *MessageRequestPayload
}
