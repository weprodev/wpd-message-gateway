package domain

import "time"

// WorkspaceSetting is a key-value row for workspace-scoped configuration.
type WorkspaceSetting struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
