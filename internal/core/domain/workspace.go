package domain

import "time"

// Workspace is a tenant boundary for messaging configuration and logs.
type Workspace struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	UniqueKey  string    `json:"unique_key"`
	AdminEmail string    `json:"admin_email"`
	Status     string    `json:"status"`
	Visibility string    `json:"visibility"`
	HashedPin  string    `json:"-"`
	IconKey    string    `json:"icon_key,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// User access metadata
	Role        string   `json:"role,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}
