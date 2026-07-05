package dto

import (
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

// RegisterRequest is the JSON body for POST /api/v1/auth/register.
type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// LoginRequest is the JSON body for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User is the portal-safe account shape (no password hash).
type User struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UserFromDomain maps a domain user to the portal response shape.
func UserFromDomain(u *domain.User) User {
	if u == nil {
		return User{}
	}
	return User{
		ID:            u.ID.String(),
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// UserWorkspace is a workspace as returned on login and /auth/me (no PIN hash).
type UserWorkspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	AdminEmail  string    `json:"admin_email"`
	Status      string    `json:"status"`
	Visibility  string    `json:"visibility"`
	IconKey     string    `json:"icon_key,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Role        string    `json:"role,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

// UserWorkspaceFromDomain maps a domain workspace (with access metadata) to the portal shape.
func UserWorkspaceFromDomain(w domain.Workspace) UserWorkspace {
	return UserWorkspace{
		ID:          w.ID,
		Name:        w.Name,
		Slug:        w.Slug,
		AdminEmail:  w.AdminEmail,
		Status:      w.Status,
		Visibility:  w.Visibility,
		IconKey:     w.IconKey,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		Role:        w.Role,
		Permissions: w.Permissions,
	}
}

func userWorkspacesFromDomain(workspaces []domain.Workspace) []UserWorkspace {
	out := make([]UserWorkspace, 0, len(workspaces))
	for _, w := range workspaces {
		out = append(out, UserWorkspaceFromDomain(w))
	}
	return out
}

// LoginResponse is returned after successful authentication.
type LoginResponse struct {
	Token      string          `json:"token"`
	User       User            `json:"user"`
	Workspaces []UserWorkspace `json:"workspaces"`
}

// LoginResponseFromDomain maps auth results to the login response.
func LoginResponseFromDomain(token string, user *domain.User, workspaces []domain.Workspace) LoginResponse {
	return LoginResponse{
		Token:      token,
		User:       UserFromDomain(user),
		Workspaces: userWorkspacesFromDomain(workspaces),
	}
}

// UserProfileResponse is returned by GET /api/v1/auth/me.
// User fields are embedded so JSON matches the flat portal contract (id, first_name, …, workspaces).
type UserProfileResponse struct {
	User
	Workspaces []UserWorkspace `json:"workspaces"`
}

// UserProfileResponseFromDomain maps the current user and workspaces to the profile response.
func UserProfileResponseFromDomain(user *domain.User, workspaces []domain.Workspace) UserProfileResponse {
	return UserProfileResponse{
		User:       UserFromDomain(user),
		Workspaces: userWorkspacesFromDomain(workspaces),
	}
}
