package domain

import "time"

// WorkspaceMember links a user to a workspace with a role.
type WorkspaceMember struct {
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
	UserEmail   string    `json:"user_email,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
}
