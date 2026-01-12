package manager

import (
	"testing"

	"github.com/weprodev/wpd-message-gateway/config"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "valid mailgun config",
			config: &config.Config{
				EmailProviders: map[string]config.EmailConfig{
					"mailgun": {
						CommonConfig: config.CommonConfig{
							APIKey: "key",
						},
						Domain: "domain",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid mailgun config",
			config: &config.Config{
				EmailProviders: map[string]config.EmailConfig{
					"mailgun": {}, // Missing required APIKey and Domain
				},
			},
			wantErr: true,
		},
		{
			name: "empty config",
			config: &config.Config{
				EmailProviders: make(map[string]config.EmailConfig),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr, err := New(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("New() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
				return
			}
			if mgr == nil {
				t.Error("New() returned nil manager")
			}
		})
	}
}
