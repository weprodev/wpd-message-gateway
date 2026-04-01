package domain

import "time"

// Invitation is a pending email invite to a workspace.
type Invitation struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	TokenHash   string    `json:"-"`
	ExpiresAt   time.Time `json:"expires_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
