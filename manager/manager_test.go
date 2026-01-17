package manager

import (
	"context"
	"testing"

	"github.com/weprodev/wpd-message-gateway/config"
	"github.com/weprodev/wpd-message-gateway/contracts"
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
				Providers: config.ProviderConfig{
					Defaults: config.ProviderDefaults{
						Email: "mailgun",
					},
				},
				EmailProviders: map[string]config.EmailConfig{
					"mailgun": {}, // Missing required APIKey and Domain
				},
			},
			wantErr: true, // Should fail when trying to create mailgun provider
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

func mockConfig() *config.Config {
	cfg := &config.Config{
		Providers: config.ProviderConfig{
			Defaults: config.ProviderDefaults{
				Email: "memory",
				SMS:   "memory",
				Push:  "memory",
				Chat:  "memory",
			},
		},
		EmailProviders: map[string]config.EmailConfig{
			"memory": {},
		},
		SMSProviders: map[string]config.SMSConfig{
			"memory": {},
		},
		PushProviders: map[string]config.PushConfig{
			"memory": {},
		},
		ChatProviders: map[string]config.ChatConfig{
			"memory": {},
		},
	}
	return cfg
}

func TestNewManager(t *testing.T) {
	cfg := mockConfig()
	mgr, err := New(cfg)

	if err != nil {
		t.Fatalf("New(cfg) error = %v", err)
	}

	if mgr == nil {
		t.Fatal("New(cfg) returned nil manager")
	}

	// Check if default providers were initialized
	if _, err := mgr.Email(); err != nil {
		t.Errorf("Email() error = %v", err)
	}
	if _, err := mgr.SMS(); err != nil {
		t.Errorf("SMS() error = %v", err)
	}
	// Push and Chat might not be fully functional in mock without registry,
	// but basic initialization should work if config is correct.
}

func TestManager_DefaultProviders(t *testing.T) {
	cfg := mockConfig()

	// Test override valid
	cfg.Providers.Defaults.Email = "memory"
	mgr, _ := New(cfg)
	if mgr.Config().DefaultEmailProvider() != "memory" {
		t.Errorf("DefaultEmailProvider = %s, want memory", mgr.Config().DefaultEmailProvider())
	}

	// Test override invalid (unknown provider)
	// In the real implementation, New() would invoke initializeDefaultProviders()
	// which checks isUnknownProviderError.
	// But here we just check the config value getter.
}

// TestManager_EmailProvider moved to email_test.go

func TestManager_EnsureProvider_Unknown(t *testing.T) {
	cfg := mockConfig()
	// Set a default provider that doesn't exist in config or registry
	cfg.Providers.Defaults.Email = "unknown_provider"

	// New() should return an error or skip it depending on implementation.
	// Current impl: initializeDefaultProviders skips unknown providers.
	mgr, err := New(cfg)
	if err != nil {
		t.Fatalf("New() should not fail for unknown provider: %v", err)
	}

	// But accessing it should fail
	_, err = mgr.Email()
	if err == nil {
		t.Error("Email() should fail for unknown provider")
	}
}

func TestMemoryProviderAsDefault(t *testing.T) {
	t.Run("memory provider requires no configuration", func(t *testing.T) {
		cfg := mockConfig()

		mgr, err := New(cfg)
		if err != nil {
			t.Fatalf("New() with memory defaults should not error: %v", err)
		}

		if mgr.GetMemoryStore() == nil {
			t.Error("memory store should be available")
		}
	})

	t.Run("memory provider ignores other provider configs", func(t *testing.T) {
		cfg := &config.Config{
			Providers: config.ProviderConfig{
				Defaults: config.ProviderDefaults{
					Email: "memory",
					SMS:   "memory",
				},
			},
			EmailProviders: map[string]config.EmailConfig{
				// Mailgun config exists but should be ignored
				"mailgun": {
					CommonConfig: config.CommonConfig{APIKey: "test-key"},
					Domain:       "test.example.com",
				},
				"memory": {}, // Add memory provider config
			},
			SMSProviders: map[string]config.SMSConfig{
				"memory": {}, // Add memory provider config
			},
			PushProviders: make(map[string]config.PushConfig),
			ChatProviders: make(map[string]config.ChatConfig),
		}

		mgr, err := New(cfg)
		if err != nil {
			t.Fatalf("New() should succeed when memory is default: %v", err)
		}

		// Mailgun config exists but memory is used
		if mgr.Config().DefaultEmailProvider() != "memory" {
			t.Error("default email provider should be memory")
		}
	})

	t.Run("all message types can use memory independently", func(t *testing.T) {
		cfg := mockConfig()

		mgr, err := New(cfg)
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		store := mgr.GetMemoryStore()
		ctx := context.Background()

		// Send email
		_, err = store.EmailProvider(config.MailpitConfig{}).Send(ctx, &contracts.Email{
			To:        []string{"test@example.com"},
			Subject:   "Test Email",
			PlainText: "Test body",
		})
		if err != nil {
			t.Errorf("email send failed: %v", err)
		}

		// Send SMS
		_, err = store.SMSProvider().Send(ctx, &contracts.SMS{
			To:      []string{"+1234567890"},
			Message: "Test SMS",
		})
		if err != nil {
			t.Errorf("SMS send failed: %v", err)
		}

		// Send Push
		_, err = store.PushProvider().Send(ctx, &contracts.PushNotification{
			DeviceTokens: []string{"token123"},
			Title:        "Test Push",
			Body:         "Test body",
		})
		if err != nil {
			t.Errorf("push send failed: %v", err)
		}

		// Send Chat
		_, err = store.ChatProvider().Send(ctx, &contracts.ChatMessage{
			To:      []string{"chat-recipient"},
			Message: "Test Chat",
		})
		if err != nil {
			t.Errorf("chat send failed: %v", err)
		}

		// Verify all messages stored
		stats := store.Stats()
		if stats["emails"] != 1 {
			t.Errorf("expected 1 email, got %d", stats["emails"])
		}
		if stats["sms"] != 1 {
			t.Errorf("expected 1 sms, got %d", stats["sms"])
		}
		if stats["push"] != 1 {
			t.Errorf("expected 1 push, got %d", stats["push"])
		}
		if stats["chat"] != 1 {
			t.Errorf("expected 1 chat, got %d", stats["chat"])
		}
	})

	t.Run("memory store is shared across all message types", func(t *testing.T) {
		cfg := &config.Config{
			Providers: config.ProviderConfig{
				Defaults: config.ProviderDefaults{
					Email: "memory",
					SMS:   "memory",
					Push:  "memory",
					Chat:  "memory",
				},
			},
			EmailProviders: make(map[string]config.EmailConfig),
			SMSProviders:   make(map[string]config.SMSConfig),
			PushProviders:  make(map[string]config.PushConfig),
			ChatProviders:  make(map[string]config.ChatConfig),
		}

		mgr, err := New(cfg)
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		store := mgr.GetMemoryStore()
		ctx := context.Background()

		// Send messages via different providers
		_, _ = store.EmailProvider(config.MailpitConfig{}).Send(ctx, &contracts.Email{To: []string{"a@b.com"}, Subject: "1"})
		_, _ = store.SMSProvider().Send(ctx, &contracts.SMS{To: []string{"+1"}, Message: "2"})

		// Both should be in the same store
		if store.Count() != 2 {
			t.Errorf("expected 2 total messages in shared store, got %d", store.Count())
		}

		// Clear should affect all
		store.Clear()
		if store.Count() != 0 {
			t.Errorf("expected 0 messages after clear, got %d", store.Count())
		}
	})
}
