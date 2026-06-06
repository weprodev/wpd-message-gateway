package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	pkgconfig "github.com/weprodev/go-pkg/config"

	"github.com/weprodev/wpd-message-gateway/internal/registry"
)

// ValidateConfig validates required configuration.
func ValidateConfig(cfg *Config) error {
	// 1. Validate JWT Secret
	jwtSecret := cfg.Portal.JWTSecret
	if v := os.Getenv("MESSAGE_JWT_SECRET"); v != "" {
		jwtSecret = v
	}
	if jwtSecret == "REPLACE_JWT_SECRET" {
		return fmt.Errorf("MESSAGE_JWT_SECRET must not use the default placeholder 'REPLACE_JWT_SECRET'")
	}
	if len(jwtSecret) < 32 {
		return fmt.Errorf("MESSAGE_JWT_SECRET (or portal.jwt_secret) must be at least 32 characters, got %d", len(jwtSecret))
	}

	// 2. Validate Config Encryption Key
	encKey := os.Getenv("MESSAGE_CONFIG_ENCRYPTION_KEY")
	if len(encKey) < 32 {
		return fmt.Errorf("MESSAGE_CONFIG_ENCRYPTION_KEY must be at least 32 characters, got %d", len(encKey))
	}

	// 3. Validate Database connection fields
	dbConfig := pkgconfig.ApplyDatabaseOverrides(pkgconfig.DatabaseConfig{})
	dbErrs, dbWarns := pkgconfig.ValidateDatabaseConfig(dbConfig)
	for _, w := range dbWarns {
		slog.Warn("database configuration warning", "warning", w)
	}
	if len(dbErrs) > 0 {
		return fmt.Errorf("invalid database configuration: %s", strings.Join(dbErrs, "; "))
	}

	// 4. Log warning if Portal Base URL is not configured
	if cfg.Portal.BaseURL == "" {
		slog.Warn("portal.base_url is not configured. Email verification links will be relative or broken.")
	}

	missingProviders := []string{}

	// At least one default provider should be configured
	if cfg.DefaultEmailProvider() == "" {
		missingProviders = append(missingProviders, "EMAIL")
	} else if !registry.IsEmailProviderRegistered(cfg.DefaultEmailProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_EMAIL_PROVIDER (unknown provider: %s)",
			cfg.DefaultEmailProvider(),
		)
	}

	if cfg.DefaultSMSProvider() == "" {
		missingProviders = append(missingProviders, "SMS")
	} else if !registry.IsSMSProviderRegistered(cfg.DefaultSMSProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_SMS_PROVIDER (unknown provider: %s)",
			cfg.DefaultSMSProvider(),
		)
	}

	if cfg.DefaultPushProvider() == "" {
		missingProviders = append(missingProviders, "PUSH")
	} else if !registry.IsPushProviderRegistered(cfg.DefaultPushProvider()) {
		return fmt.Errorf(
			"missing or invalid required configuration: MESSAGE_DEFAULT_PUSH_PROVIDER (unknown provider: %s)",
			cfg.DefaultPushProvider(),
		)
	}

	if cfg.DefaultChatProvider() == "" {
		missingProviders = append(missingProviders, "CHAT")
	} else if !registry.IsChatProviderRegistered(cfg.DefaultChatProvider()) {
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

	if len(missingProviders) > 0 {
		slog.Warn("no default provider configured for channels; those channels will use memory fallback",
			"channels", strings.Join(missingProviders, ", "))
	}

	return nil
}
