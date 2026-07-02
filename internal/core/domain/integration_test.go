package domain

import (
	"errors"
	"testing"
)

func TestValidateProviderIntegration(t *testing.T) {
	valid := Integration{
		ID:           "int-1",
		ProviderName: "mailgun",
		Config:       []byte(`{"api_key":"key"}`),
		Status:       IntegrationStatusConnected,
	}

	tests := []struct {
		name    string
		intg    Integration
		wantErr bool
	}{
		{"valid integration", valid, false},
		{"empty integration", Integration{}, true},
		{"memory provider", Integration{ProviderName: ProviderNameMemory, Config: []byte(`{}`), Status: IntegrationStatusConnected}, true},
		{"disconnected", Integration{ProviderName: "mailgun", Config: []byte(`{}`), Status: IntegrationStatusDisconnected}, true},
		{"empty config", Integration{ProviderName: "mailgun", Config: []byte(`   `), Status: IntegrationStatusConnected}, true},
		{"invalid json", Integration{ProviderName: "mailgun", Config: []byte(`{`), Status: IntegrationStatusConnected}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProviderIntegration(tt.intg)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if !errors.Is(err, ErrProviderNotReady) {
					t.Fatalf("expected ErrProviderNotReady, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
