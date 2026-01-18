package gateway

import (
	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
	"github.com/weprodev/wpd-message-gateway/internal/core/service"
)

// Config holds the gateway configuration.
type Config struct {
	DefaultEmailProvider string
	DefaultSMSProvider   string
	DefaultPushProvider  string
	DefaultChatProvider  string

	// Provider-specific configurations keyed by provider name.
	// Uses registry types as single source of truth.
	EmailProviders map[string]registry.EmailConfig
	SMSProviders   map[string]registry.SMSConfig
	PushProviders  map[string]registry.PushConfig
	ChatProviders  map[string]registry.ChatConfig

	// MailpitEnabled enables SMTP forwarding for the memory provider.
	MailpitEnabled bool
}

// Type aliases for SDK users - these reference the canonical registry types.
// Users can use either gateway.EmailConfig or registry.EmailConfig.
type (
	CommonConfig = registry.CommonConfig
	EmailConfig  = registry.EmailConfig
	SMSConfig    = registry.SMSConfig
	PushConfig   = registry.PushConfig
	ChatConfig   = registry.ChatConfig
)

// Gateway is the main entry point for sending messages.
type Gateway struct {
	service *service.GatewayService
	cfg     Config
}

// configAdapter adapts Gateway config to service.GatewayConfig interface.
type configAdapter struct {
	cfg Config
}
