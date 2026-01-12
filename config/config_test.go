package config

import (
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	// Use t.Setenv for automatic cleanup
	t.Setenv("MESSAGE_DEFAULT_EMAIL_PROVIDER", "mailgun")
	t.Setenv("MESSAGE_MAILGUN_API_KEY", "test-api-key")
	t.Setenv("MESSAGE_MAILGUN_DOMAIN", "test.mailgun.org")
	t.Setenv("MESSAGE_MAILGUN_FROM_EMAIL", "test@test.com")
	t.Setenv("MESSAGE_MAILGUN_FROM_NAME", "Test App")
	t.Setenv("MESSAGE_MAILGUN_BASE_URL", "https://api.eu.mailgun.net")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error: %v", err)
	}

	// Verify default provider
	if cfg.DefaultEmailProvider != "mailgun" {
		t.Errorf("DefaultEmailProvider = %s, want mailgun", cfg.DefaultEmailProvider)
	}

	// Verify provider was auto-discovered and loaded
	mailgunCfg, ok := cfg.EmailProviders["mailgun"]
	if !ok {
		t.Fatal("Expected mailgun provider to be loaded")
	}

	// Verify all fields loaded correctly
	if mailgunCfg.APIKey != "test-api-key" {
		t.Errorf("APIKey = %s, want test-api-key", mailgunCfg.APIKey)
	}
	if mailgunCfg.Domain != "test.mailgun.org" {
		t.Errorf("Domain = %s, want test.mailgun.org", mailgunCfg.Domain)
	}
	if mailgunCfg.FromEmail != "test@test.com" {
		t.Errorf("FromEmail = %s, want test@test.com", mailgunCfg.FromEmail)
	}
	if mailgunCfg.BaseURL != "https://api.eu.mailgun.net" {
		t.Errorf("BaseURL = %s, want https://api.eu.mailgun.net", mailgunCfg.BaseURL)
	}
}

func TestRegisterProvider(t *testing.T) {
	// Register a new custom provider
	RegisterProvider("customtest", ProviderTypeEmail)

	// Use t.Setenv for automatic cleanup
	t.Setenv("MESSAGE_CUSTOMTEST_API_KEY", "custom-key")

	// Load config
	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error: %v", err)
	}

	// Verify custom provider was loaded
	provider, ok := cfg.EmailProviders["customtest"]
	if !ok {
		t.Fatal("Expected custom provider to be loaded")
	}
	if provider.APIKey != "custom-key" {
		t.Errorf("APIKey = %s, want custom-key", provider.APIKey)
	}
}
