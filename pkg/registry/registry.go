// Package registry provides the provider factory registry — a shared kernel
// that allows providers to self-register via init() without importing any
// specific layer. Both the embedded SDK (pkg/gateway) and the HTTP server
// (internal/app) consume this registry.
//
// Layer rule: this package may only depend on pkg/contracts interfaces.
// Infrastructure providers register here; Core services look up factories here.
package registry

import (
	"fmt"
	"sync"

	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
)

// CommonConfig contains credential fields shared across all provider channels.
// JSON keys match provider_config_fields / Portal integration payloads (snake_case).
type CommonConfig struct {
	APIKey    string            `json:"api_key"`
	APISecret string            `json:"api_secret"`
	Region    string            `json:"region"`
	BaseURL   string            `json:"base_url"`
	Extra     map[string]string `json:"-"`
}

// EmailConfig holds email-specific provider configuration.
type EmailConfig struct {
	CommonConfig
	Domain    string `json:"domain"`
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
}

// SMSConfig holds SMS-specific provider configuration.
type SMSConfig struct {
	CommonConfig
	FromPhone string `json:"from_phone"`
}

// PushConfig holds push-notification-specific provider configuration.
type PushConfig struct {
	CommonConfig
	AppID string `json:"app_id"`
	Topic string `json:"topic"`
}

// ChatConfig holds chat-specific provider configuration.
type ChatConfig struct {
	CommonConfig
	FromPhone  string `json:"from_phone"`
	WebhookURL string `json:"webhook_url"`
}

// EmailProviderFactory constructs an EmailSender from static config.
type EmailProviderFactory func(cfg EmailConfig) (contracts.EmailSender, error)

// SMSProviderFactory constructs an SMSSender from static config.
type SMSProviderFactory func(cfg SMSConfig) (contracts.SMSSender, error)

// PushProviderFactory constructs a PushSender from static config.
type PushProviderFactory func(cfg PushConfig) (contracts.PushSender, error)

// ChatProviderFactory constructs a ChatSender from static config.
type ChatProviderFactory func(cfg ChatConfig) (contracts.ChatSender, error)

var (
	mu             sync.RWMutex
	emailFactories = make(map[string]EmailProviderFactory)
	smsFactories   = make(map[string]SMSProviderFactory)
	pushFactories  = make(map[string]PushProviderFactory)
	chatFactories  = make(map[string]ChatProviderFactory)
)

// RegisterEmailProvider registers an email provider factory by name.
// Call this inside your provider's init() function.
func RegisterEmailProvider(name string, factory EmailProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	emailFactories[name] = factory
}

// RegisterSMSProvider registers an SMS provider factory by name.
func RegisterSMSProvider(name string, factory SMSProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	smsFactories[name] = factory
}

// RegisterPushProvider registers a push provider factory by name.
func RegisterPushProvider(name string, factory PushProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	pushFactories[name] = factory
}

// RegisterChatProvider registers a chat provider factory by name.
func RegisterChatProvider(name string, factory ChatProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	chatFactories[name] = factory
}

// GetEmailFactory returns the registered factory for the named email provider.
func GetEmailFactory(name string) (EmailProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := emailFactories[name]
	if !ok {
		return nil, fmt.Errorf("registry: unknown email provider %q (not registered)", name)
	}
	return f, nil
}

// GetSMSFactory returns the registered factory for the named SMS provider.
func GetSMSFactory(name string) (SMSProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := smsFactories[name]
	if !ok {
		return nil, fmt.Errorf("registry: unknown SMS provider %q (not registered)", name)
	}
	return f, nil
}

// GetPushFactory returns the registered factory for the named push provider.
func GetPushFactory(name string) (PushProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := pushFactories[name]
	if !ok {
		return nil, fmt.Errorf("registry: unknown push provider %q (not registered)", name)
	}
	return f, nil
}

// GetChatFactory returns the registered factory for the named chat provider.
func GetChatFactory(name string) (ChatProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := chatFactories[name]
	if !ok {
		return nil, fmt.Errorf("registry: unknown chat provider %q (not registered)", name)
	}
	return f, nil
}

// IsEmailProviderRegistered reports whether an email provider is registered.
func IsEmailProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := emailFactories[name]
	return ok
}

// IsSMSProviderRegistered reports whether an SMS provider is registered.
func IsSMSProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := smsFactories[name]
	return ok
}

// IsPushProviderRegistered reports whether a push provider is registered.
func IsPushProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := pushFactories[name]
	return ok
}

// IsChatProviderRegistered reports whether a chat provider is registered.
func IsChatProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := chatFactories[name]
	return ok
}
