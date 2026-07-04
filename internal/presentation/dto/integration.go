package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

// UpsertIntegrationRequest is the JSON body for PUT /api/v1/workspaces/:wid/integrations.
type UpsertIntegrationRequest struct {
	ChannelType  string          `json:"channel_type"`
	ProviderName string          `json:"provider_name"`
	Config       json.RawMessage `json:"config"`
	Status       string          `json:"status"`
	IsDefault    bool            `json:"is_default"`
}

// ToDomain maps the request to a domain integration for upsert (passed to service → repository).
func (r UpsertIntegrationRequest) ToDomain(workspaceID string) (*domain.Integration, error) {
	cfg, err := normalizeIntegrationConfig(r.Config)
	if err != nil {
		return nil, fmt.Errorf("integration config: %w", err)
	}
	status := r.Status
	if status == "" {
		status = domain.IntegrationStatusConnected
	}
	return &domain.Integration{
		WorkspaceID:  workspaceID,
		ChannelType:  r.ChannelType,
		ProviderName: r.ProviderName,
		Config:       cfg,
		Status:       status,
		IsDefault:    r.IsDefault,
	}, nil
}

// Integration is the portal response shape for workspace integrations.
type Integration struct {
	ID           string          `json:"id"`
	WorkspaceID  string          `json:"workspace_id"`
	ChannelType  string          `json:"channel_type"`
	ProviderName string          `json:"provider_name"`
	Config       json.RawMessage `json:"config"`
	Status       string          `json:"status"`
	IsDefault    bool            `json:"is_default"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// IntegrationFromDomain maps a domain integration to its portal response shape.
func IntegrationFromDomain(intg domain.Integration) Integration {
	return Integration{
		ID:           intg.ID,
		WorkspaceID:  intg.WorkspaceID,
		ChannelType:  intg.ChannelType,
		ProviderName: intg.ProviderName,
		Config:       json.RawMessage(intg.Config),
		Status:       intg.Status,
		IsDefault:    intg.IsDefault,
		CreatedAt:    intg.CreatedAt,
		UpdatedAt:    intg.UpdatedAt,
	}
}

func normalizeIntegrationConfig(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return []byte("{}"), nil
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return json.Marshal(v)
}
