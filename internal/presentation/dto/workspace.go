package dto

import (
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

// CreateWorkspaceRequest is the JSON body for POST /api/v1/workspaces.
type CreateWorkspaceRequest struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	IconKey string `json:"icon_key"`
}

// JoinWorkspaceRequest is the JSON body for POST /api/v1/workspaces/join.
type JoinWorkspaceRequest struct {
	Slug string `json:"slug"`
	PIN  string `json:"pin"`
}

// PatchWorkspaceRequest is the JSON body for PATCH /api/v1/workspaces/:wid.
type PatchWorkspaceRequest struct {
	Name       *string `json:"name"`
	Visibility *string `json:"visibility"`
	IconKey    *string `json:"icon_key"`
}

// ToPatch maps the request to the service patch command (service loads domain, applies, saves via repository).
func (r PatchWorkspaceRequest) ToPatch() service.WorkspacePatch {
	return service.WorkspacePatch{
		Name:       r.Name,
		Visibility: r.Visibility,
		IconKey:    r.IconKey,
	}
}

// CreateInvitationRequest is the JSON body for POST /api/v1/workspaces/:wid/invitations.
type CreateInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ToDomain maps the request to a domain invitation for persistence.
func (r CreateInvitationRequest) ToDomain(workspaceID string, expiresAt time.Time) *domain.Invitation {
	return &domain.Invitation{
		WorkspaceID: workspaceID,
		Email:       r.Email,
		Role:        r.Role,
		ExpiresAt:   expiresAt,
		Status:      "pending",
	}
}

// CreateInvitationResponse is returned once when an invitation is created.
type CreateInvitationResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}

// CreateInvitationResponseFromDomain maps a persisted invitation and one-time token.
func CreateInvitationResponseFromDomain(inv domain.Invitation, token string) CreateInvitationResponse {
	return CreateInvitationResponse{
		ID:        inv.ID,
		Email:     inv.Email,
		Role:      inv.Role,
		ExpiresAt: inv.ExpiresAt,
		Token:     token,
	}
}
