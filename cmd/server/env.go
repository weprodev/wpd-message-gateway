package main

import (
	"fmt"
	"strings"

	"github.com/weprodev/wpd-message-gateway/config"
)

// validateConfig checks that all required default providers are configured.
// The gateway requires a default provider for each message type to function.
func validateConfig(cfg *config.Config) error {
	var missing []string

	if cfg.DefaultEmailProvider() == "" {
		missing = append(missing, "MESSAGE_DEFAULT_EMAIL_PROVIDER")
	} else if !config.IsKnownProvider(cfg.DefaultEmailProvider()) {
		missing = append(missing, fmt.Sprintf("MESSAGE_DEFAULT_EMAIL_PROVIDER (unknown provider: %s)", cfg.DefaultEmailProvider()))
	}

	if cfg.DefaultSMSProvider() == "" {
		missing = append(missing, "MESSAGE_DEFAULT_SMS_PROVIDER")
	} else if !config.IsKnownProvider(cfg.DefaultSMSProvider()) {
		missing = append(missing, fmt.Sprintf("MESSAGE_DEFAULT_SMS_PROVIDER (unknown provider: %s)", cfg.DefaultSMSProvider()))
	}

	if cfg.DefaultPushProvider() == "" {
		missing = append(missing, "MESSAGE_DEFAULT_PUSH_PROVIDER")
	} else if !config.IsKnownProvider(cfg.DefaultPushProvider()) {
		missing = append(missing, fmt.Sprintf("MESSAGE_DEFAULT_PUSH_PROVIDER (unknown provider: %s)", cfg.DefaultPushProvider()))
	}

	if cfg.DefaultChatProvider() == "" {
		missing = append(missing, "MESSAGE_DEFAULT_CHAT_PROVIDER")
	} else if !config.IsKnownProvider(cfg.DefaultChatProvider()) {
		missing = append(missing, fmt.Sprintf("MESSAGE_DEFAULT_CHAT_PROVIDER (unknown provider: %s)", cfg.DefaultChatProvider()))
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"missing or invalid required configuration:\n"+
				"  %s\n\n"+
				"Each message type requires a valid default provider (e.g. 'memory', 'mailgun', 'twilio').\n"+
				"Please configure these in configs/local.yml or via environment variables.",
			strings.Join(missing, "\n  "),
		)
	}

	return nil
}
