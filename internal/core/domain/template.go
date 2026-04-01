package domain

import "time"

// Template is a reusable message body for a workspace channel.
type Template struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Name        string    `json:"name"`
	UniqueKey   string    `json:"unique_key"`
	ChannelType string    `json:"channel_type"`
	Category    string    `json:"category,omitempty"`
	Subject     string    `json:"subject,omitempty"`
	ContentHTML string    `json:"content_html"`
	IsActive    bool      `json:"is_active"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
