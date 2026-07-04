package mailgun

import (
	"strings"
	"testing"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

func TestNew_validationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name:    "missing api key",
			cfg:     Config{Domain: "mg.example.com"},
			wantErr: "API key is required",
		},
		{
			name:    "missing domain",
			cfg:     Config{APIKey: "key-123"},
			wantErr: "domain is required",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := New(tt.cfg)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected %q in error, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNew_success(t *testing.T) {
	t.Parallel()

	p, err := New(Config{APIKey: "key-123", Domain: "mg.example.com", FromEmail: "noreply@mg.example.com"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if p.fromAddress != "noreply@mg.example.com" {
		t.Fatalf("expected from address on provider, got %q", p.fromAddress)
	}
}

func TestProvider_buildFromAddress(t *testing.T) {
	t.Parallel()

	defaultProvider := &Provider{
		fromAddress: "noreply@mg.example.com",
		fromName:    "Acme",
	}

	tests := []struct {
		name     string
		provider *Provider
		email    contracts.Email
		want     string
	}{
		{
			name:     "email overrides defaults",
			provider: defaultProvider,
			email:    contracts.Email{From: "custom@example.com", FromName: "Custom"},
			want:     "Custom <custom@example.com>",
		},
		{
			name:     "uses provider defaults",
			provider: defaultProvider,
			email:    contracts.Email{},
			want:     "Acme <noreply@mg.example.com>",
		},
		{
			name:     "email from with default name",
			provider: defaultProvider,
			email:    contracts.Email{From: "only@example.com"},
			want:     "Acme <only@example.com>",
		},
		{
			name:     "address only when no name",
			provider: &Provider{fromAddress: "noreply@mg.example.com"},
			email:    contracts.Email{From: "plain@example.com"},
			want:     "plain@example.com",
		},
		{
			name:     "empty when no address available",
			provider: &Provider{fromName: "Name Only"},
			email:    contracts.Email{},
			want:     "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.provider.buildFromAddress(tt.email)
			if got != tt.want {
				t.Fatalf("buildFromAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProvider_Send_rejectsMissingFrom(t *testing.T) {
	t.Parallel()

	p := &Provider{}
	_, err := p.Send(t.Context(), contracts.Email{To: []string{"a@b.com"}})
	if err == nil || err.Error() != "no from address specified" {
		t.Fatalf("expected missing from error, got %v", err)
	}
}
