package mailgun

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  config.EmailConfig
		wantErr bool
	}{
		{
			name:    "empty config",
			config:  config.EmailConfig{},
			wantErr: true,
		},
		{
			name: "config with mailgun missing required fields",
			config: config.EmailConfig{
				CommonConfig: config.CommonConfig{
					APIKey: "key-123",
				},
			},
			wantErr: true,
		},
		{
			name: "config with valid mailgun",
			config: config.EmailConfig{
				CommonConfig: config.CommonConfig{
					APIKey: "key-123",
				},
				Domain: "example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("New() expected error, got nil")
					return
				}
				// No specific error message check in new tests
				return
			}
			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
				return
			}
			if provider == nil {
				t.Error("New() returned nil provider")
			}
		})
	}
}

func TestProvider_buildFromAddress(t *testing.T) {
	tests := []struct {
		name       string
		config     config.EmailConfig
		email      *contracts.Email
		wantResult string
	}{
		{
			name: "use default from config",
			config: config.EmailConfig{
				CommonConfig: config.CommonConfig{APIKey: "k"},
				Domain:       "d",
				FromEmail:    "default@example.com",
				FromName:     "Default",
			},
			email:      &contracts.Email{To: []string{"User"}},
			wantResult: "Default <default@example.com>",
		},
		{
			name: "override from address",
			config: config.EmailConfig{
				CommonConfig: config.CommonConfig{APIKey: "k"},
				Domain:       "d",
				FromEmail:    "default@example.com",
			},
			email: &contracts.Email{
				To:   []string{"User"},
				From: "custom@example.com",
			},
			wantResult: "custom@example.com",
		},
		{
			name: "only override address",
			config: config.EmailConfig{
				CommonConfig: config.CommonConfig{APIKey: "k"},
				Domain:       "d",
				FromEmail:    "default@example.com",
				FromName:     "Default",
			},
			email: &contracts.Email{
				To:   []string{"User"},
				From: "custom@example.com",
			},
			wantResult: "Default <custom@example.com>",
		},
		{
			name: "only override name",
			config: config.EmailConfig{
				CommonConfig: config.CommonConfig{APIKey: "k"},
				Domain:       "d",
				FromEmail:    "default@example.com",
				FromName:     "Default",
			},
			email: &contracts.Email{
				To:       []string{"User"},
				FromName: "Custom Name",
			},
			wantResult: "Custom Name <default@example.com>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.config)
			if err != nil {
				t.Fatalf("Failed to create provider for test %s: %v", tt.name, err)
			}
			got := provider.buildFromAddress(tt.email)
			if got != tt.wantResult {
				t.Errorf("buildFromAddress() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestProvider_Send_Validation(t *testing.T) {
	provider, _ := New(config.EmailConfig{
		CommonConfig: config.CommonConfig{APIKey: "key"},
		Domain:       "domain",
		FromEmail:    "from@test.com",
	})

	tests := []struct {
		name    string
		email   *contracts.Email
		wantErr bool
	}{
		{
			name:    "no recipients",
			email:   &contracts.Email{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Send(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
