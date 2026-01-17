package app

import (
	"fmt"
	"strings"
)

// ValidateConfig validates required configuration.
func ValidateConfig(cfg *Config) error {
	missingProviders := []string{}

	// At least one default provider should be configured
	if cfg.DefaultEmailProvider() == "" {
		missingProviders = append(missingProviders, "EMAIL")
	} else if !IsKnownProvider(cfg.DefaultEmailProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_EMAIL_PROVIDER (unknown provider: %s)",
			cfg.DefaultEmailProvider(),
		)
	}

	if cfg.DefaultSMSProvider() == "" {
		missingProviders = append(missingProviders, "SMS")
	} else if !IsKnownProvider(cfg.DefaultSMSProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_SMS_PROVIDER (unknown provider: %s)",
			cfg.DefaultSMSProvider(),
		)
	}

	if cfg.DefaultPushProvider() == "" {
		missingProviders = append(missingProviders, "PUSH")
	} else if !IsKnownProvider(cfg.DefaultPushProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_PUSH_PROVIDER (unknown provider: %s)",
			cfg.DefaultPushProvider(),
		)
	}

	if cfg.DefaultChatProvider() == "" {
		missingProviders = append(missingProviders, "CHAT")
	} else if !IsKnownProvider(cfg.DefaultChatProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_CHAT_PROVIDER (unknown provider: %s)",
			cfg.DefaultChatProvider(),
		)
	}

	// If ALL providers are missing, that's an error
	if len(missingProviders) == 4 {
		return fmt.Errorf(
			"no default providers configured. Please set at least one in configs/local.yml:\n" +
				"  providers:\n" +
				"    defaults:\n" +
				"      email: memory\n" +
				"      sms: memory\n" +
				"      push: memory\n" +
				"      chat: memory",
		)
	}

	// Log warnings for missing providers
	if len(missingProviders) > 0 {
		fmt.Printf("Note: No default provider configured for: %s\n", strings.Join(missingProviders, ", "))
	}

	return nil
}
