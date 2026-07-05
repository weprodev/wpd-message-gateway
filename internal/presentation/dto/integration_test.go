package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/weprodev/wpd-message-gateway/internal/core/domain"
)

func TestIntegrationFromDomain(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	intg := domain.Integration{
		ID:           "intg-1",
		WorkspaceID:  "ws-1",
		ChannelType:  "email",
		ProviderName: "mailgun",
		Config:       []byte(`{"api_key":"secret"}`),
		Status:       domain.IntegrationStatusConnected,
		IsDefault:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	got := IntegrationFromDomain(intg)
	if got.ID != intg.ID || got.WorkspaceID != intg.WorkspaceID {
		t.Fatalf("unexpected ids: %+v", got)
	}
	if string(got.Config) != string(intg.Config) {
		t.Fatalf("config = %s, want %s", got.Config, intg.Config)
	}

	var cfg map[string]string
	if err := json.Unmarshal(got.Config, &cfg); err != nil {
		t.Fatalf("config not valid JSON: %v", err)
	}
	if cfg["api_key"] != "secret" {
		t.Errorf("api_key = %q, want secret", cfg["api_key"])
	}
}

func TestUpsertIntegrationRequestToDomain(t *testing.T) {
	req := UpsertIntegrationRequest{
		ChannelType:  "email",
		ProviderName: "mailgun",
		Config:       json.RawMessage(`{"api_key":"secret"}`),
		IsDefault:    true,
	}

	got, err := req.ToDomain("ws-1")
	if err != nil {
		t.Fatalf("ToDomain() error = %v", err)
	}
	if got.WorkspaceID != "ws-1" || got.ProviderName != "mailgun" {
		t.Fatalf("unexpected domain mapping: %+v", got)
	}
	if got.Status != domain.IntegrationStatusConnected {
		t.Errorf("status = %q, want connected", got.Status)
	}
}
