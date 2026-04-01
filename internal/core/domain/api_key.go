package domain

import "time"

// APIKey is a client credential scoped to a workspace.
type APIKey struct {
	ID               string `json:"id"`
	WorkspaceID      string `json:"workspace_id"`
	ClientID         string `json:"client_id"`
	ClientSecretHash string `json:"-"`
	// Name identifies the consumer (product or service) for logs and UI.
	Name       string     `json:"name"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}
