package dto

import (
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

// APIKeyPublic is the portal-safe API key representation (no secret hash).
type APIKeyPublic struct {
	ID          string     `json:"id"`
	WorkspaceID string     `json:"workspace_id"`
	ClientID    string     `json:"client_id"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// APIKeyPublicFromDomain maps a domain API key to its portal response shape.
func APIKeyPublicFromDomain(k domain.APIKey) APIKeyPublic {
	return APIKeyPublic{
		ID:          k.ID,
		WorkspaceID: k.WorkspaceID,
		ClientID:    k.ClientID,
		Name:        k.Name,
		IsActive:    k.IsActive,
		LastUsedAt:  k.LastUsedAt,
		CreatedAt:   k.CreatedAt,
		ExpiresAt:   k.ExpiresAt,
	}
}

// CreateAPIKeyRequest is the JSON body for POST /api/v1/workspaces/:wid/api-keys.
type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

// CreateAPIKeyResponse is returned once when a new API key secret is issued.
type CreateAPIKeyResponse struct {
	APIKeyPublic
	ClientSecret string `json:"client_secret"`
}

// CreateAPIKeyResponseFromDomain maps a newly created domain API key and one-time secret.
func CreateAPIKeyResponseFromDomain(k domain.APIKey, secret string) CreateAPIKeyResponse {
	return CreateAPIKeyResponse{
		APIKeyPublic: APIKeyPublicFromDomain(k),
		ClientSecret: secret,
	}
}

// RegenerateAPIKeyResponse is returned when an API key secret is re-issued.
type RegenerateAPIKeyResponse struct {
	ClientSecret string `json:"client_secret"`
}

// RegenerateAPIKeyResponseFromDomain maps a re-issued secret to the portal response.
func RegenerateAPIKeyResponseFromDomain(secret string) RegenerateAPIKeyResponse {
	return RegenerateAPIKeyResponse{ClientSecret: secret}
}
