package registry

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

type mockEmailSender struct{}

func (m *mockEmailSender) Send(ctx context.Context, e contracts.Email) (*contracts.SendResult, error) {
	return nil, nil
}
func (m *mockEmailSender) Name() string { return "mock" }

func TestRegistry(t *testing.T) {
	// Register a mock
	RegisterEmailProvider("mock", func(cfg EmailConfig) (contracts.EmailSender, error) {
		return &mockEmailSender{}, nil
	})

	// Test success
	factory, err := GetEmailFactory("mock")
	if err != nil {
		t.Errorf("GetEmailFactory failed: %v", err)
	}
	if factory == nil {
		t.Errorf("expected non-nil factory")
	}

	// Test missing
	_, err = GetEmailFactory("missing")
	if err == nil {
		t.Errorf("expected error for missing provider, got nil")
	}
}

func TestEmailConfig_UnmarshalPortalIntegrationJSON(t *testing.T) {
	raw := []byte(`{
		"api_key": "key-test",
		"domain": "mg.example.com",
		"base_url": "https://api.mailgun.net/v3",
		"from_email": "noreply@mg.example.com",
		"from_name": "Support"
	}`)

	var cfg EmailConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cfg.APIKey != "key-test" {
		t.Fatalf("APIKey: got %q", cfg.APIKey)
	}
	if cfg.Domain != "mg.example.com" {
		t.Fatalf("Domain: got %q", cfg.Domain)
	}
	if cfg.FromEmail != "noreply@mg.example.com" {
		t.Fatalf("FromEmail: got %q", cfg.FromEmail)
	}
}
