package app

import (
	"os"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	// Backup env vars
	origJWT := os.Getenv("MESSAGE_JWT_SECRET")
	origEncKey := os.Getenv("MESSAGE_CONFIG_ENCRYPTION_KEY")
	origDBURL := os.Getenv("DATABASE_URL")
	defer func() {
		os.Setenv("MESSAGE_JWT_SECRET", origJWT)
		os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", origEncKey)
		os.Setenv("DATABASE_URL", origDBURL)
	}()

	tests := []struct {
		name       string
		setup      func()
		config     *Config
		wantErr    bool
		errMessage string
	}{
		{
			name: "valid configuration",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "0123456789abcdef0123456789abcdef")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
				os.Setenv("DATABASE_URL", "postgres://localhost/test")
			},
			config: &Config{
				Providers: ProviderConfig{
					Defaults: ProviderDefaults{
						Email: "memory",
						SMS:   "memory",
						Push:  "memory",
						Chat:  "memory",
					},
				},
				Portal: PortalConfig{
					BaseURL: "https://message-gateway.weprodev.com",
				},
			},
			wantErr: false,
		},
		{
			name: "missing portal base url - log warning only",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "0123456789abcdef0123456789abcdef")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
				os.Setenv("DATABASE_URL", "postgres://localhost/test")
			},
			config: &Config{
				Providers: ProviderConfig{
					Defaults: ProviderDefaults{
						Email: "memory",
						SMS:   "memory",
						Push:  "memory",
						Chat:  "memory",
					},
				},
				Portal: PortalConfig{
					BaseURL: "", // empty base URL
				},
			},
			wantErr: false,
		},
		{
			name: "jwt secret too short",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "short")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
				os.Setenv("DATABASE_URL", "postgres://localhost/test")
			},
			config:     &Config{},
			wantErr:    true,
			errMessage: "MESSAGE_JWT_SECRET (or portal.jwt_secret) must be at least 32 characters",
		},
		{
			name: "jwt secret uses default placeholder",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "REPLACE_JWT_SECRET")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
				os.Setenv("DATABASE_URL", "postgres://localhost/test")
			},
			config:     &Config{},
			wantErr:    true,
			errMessage: "MESSAGE_JWT_SECRET must not use the default placeholder 'REPLACE_JWT_SECRET'",
		},
		{
			name: "encryption key too short",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "0123456789abcdef0123456789abcdef")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "short-key")
				os.Setenv("DATABASE_URL", "postgres://localhost/test")
			},
			config:     &Config{},
			wantErr:    true,
			errMessage: "MESSAGE_CONFIG_ENCRYPTION_KEY must be at least 32 characters",
		},
		{
			name: "invalid database config",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "0123456789abcdef0123456789abcdef")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
				os.Unsetenv("DATABASE_URL")
				os.Unsetenv("DB_HOST")
				os.Unsetenv("DB_CONNECTION_NAME")
			},
			config:     &Config{},
			wantErr:    true,
			errMessage: "invalid database configuration",
		},
		{
			name: "no default providers configured",
			setup: func() {
				os.Setenv("MESSAGE_JWT_SECRET", "0123456789abcdef0123456789abcdef")
				os.Setenv("MESSAGE_CONFIG_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
				os.Setenv("DATABASE_URL", "postgres://localhost/test")
			},
			config: &Config{
				Providers: ProviderConfig{
					Defaults: ProviderDefaults{},
				},
			},
			wantErr:    true,
			errMessage: "no default providers configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				if tt.errMessage != "" && !contains(err.Error(), tt.errMessage) {
					t.Errorf("expected error containing %q, got %q", tt.errMessage, err.Error())
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || checkSubstr(s, substr)))
}

func checkSubstr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
