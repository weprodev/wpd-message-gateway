package dto

import (
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

// MessageRequestLog is the portal list shape for gateway request logs.
type MessageRequestLog struct {
	ID             string    `json:"id"`
	WorkspaceID    string    `json:"workspace_id"`
	APIKeyID       string    `json:"api_key_id,omitempty"`
	ChannelType    string    `json:"channel_type"`
	HTTPMethod     string    `json:"http_method"`
	StatusCode     int       `json:"status_code"`
	Endpoint       string    `json:"endpoint"`
	ProviderName   string    `json:"provider_name,omitempty"`
	RequestID      string    `json:"request_id,omitempty"`
	DurationMs     int       `json:"duration_ms,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	InboxMessageID string    `json:"inbox_message_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	SourceName     string    `json:"source_name"`
	ClientID       string    `json:"client_id,omitempty"`
}

// MessageRequestLogFromDomain maps a domain log row (with API key source metadata) to the portal DTO.
func MessageRequestLogFromDomain(l domain.MessageRequestLogWithSource) MessageRequestLog {
	return MessageRequestLog{
		ID:             l.ID,
		WorkspaceID:    l.WorkspaceID,
		APIKeyID:       l.APIKeyID,
		ChannelType:    l.ChannelType,
		HTTPMethod:     l.HTTPMethod,
		StatusCode:     l.StatusCode,
		Endpoint:       l.Endpoint,
		ProviderName:   l.ProviderName,
		RequestID:      l.RequestID,
		DurationMs:     l.DurationMs,
		ErrorMessage:   l.ErrorMessage,
		InboxMessageID: l.InboxMessageID,
		CreatedAt:      l.CreatedAt,
		SourceName:     l.SourceName,
		ClientID:       l.ClientID,
	}
}

// MessageRequestLogList is the paginated logs API response.
type MessageRequestLogList struct {
	Items []MessageRequestLog `json:"items"`
	Total int                 `json:"total"`
}

// MessageRequestLogListFromDomain maps domain log rows and total count to the paginated response.
func MessageRequestLogListFromDomain(rows []domain.MessageRequestLogWithSource, total int) MessageRequestLogList {
	items := make([]MessageRequestLog, 0, len(rows))
	for _, r := range rows {
		items = append(items, MessageRequestLogFromDomain(r))
	}
	return MessageRequestLogList{Items: items, Total: total}
}
