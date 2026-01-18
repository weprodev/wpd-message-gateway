// Package registry provides provider registration functionality.
// This is a sub-package of app to avoid circular imports when
// providers register themselves via init().
package registry

import (
	"fmt"
	"sync"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

// MailpitConfig holds SMTP forwarding configuration.
type MailpitConfig struct {
	Enabled bool
}

// CommonConfig shared fields across all providers.
type CommonConfig struct {
	APIKey    string
	APISecret string
	Region    string
	BaseURL   string
	Extra     map[string]string
}

// EmailConfig holds email provider configuration.
type EmailConfig struct {
	CommonConfig
	Domain    string
	FromEmail string
	FromName  string
}

// SMSConfig holds SMS provider configuration.
type SMSConfig struct {
	CommonConfig
	FromPhone string
}

// PushConfig holds push notification provider configuration.
type PushConfig struct {
	CommonConfig
	AppID string
	Topic string
}

// ChatConfig holds chat provider configuration.
type ChatConfig struct {
	CommonConfig
	FromPhone  string
	WebhookURL string
}

// EmailProviderFactory creates an email provider from config.
type EmailProviderFactory func(cfg EmailConfig, store port.MessageStore, mailpit MailpitConfig) (port.EmailSender, error)

// SMSProviderFactory creates an SMS provider from config.
type SMSProviderFactory func(cfg SMSConfig, store port.MessageStore) (port.SMSSender, error)

// PushProviderFactory creates a push provider from config.
type PushProviderFactory func(cfg PushConfig, store port.MessageStore) (port.PushSender, error)

// ChatProviderFactory creates a chat provider from config.
type ChatProviderFactory func(cfg ChatConfig, store port.MessageStore) (port.ChatSender, error)

var (
	emailFactories = make(map[string]EmailProviderFactory)
	smsFactories   = make(map[string]SMSProviderFactory)
	pushFactories  = make(map[string]PushProviderFactory)
	chatFactories  = make(map[string]ChatProviderFactory)
	mu             sync.RWMutex
)

// RegisterEmailProvider registers an email provider factory.
// Call this in your provider's init() function.
func RegisterEmailProvider(name string, factory EmailProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	emailFactories[name] = factory
}

// RegisterSMSProvider registers an SMS provider factory.
func RegisterSMSProvider(name string, factory SMSProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	smsFactories[name] = factory
}

// RegisterPushProvider registers a push provider factory.
func RegisterPushProvider(name string, factory PushProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	pushFactories[name] = factory
}

// RegisterChatProvider registers a chat provider factory.
func RegisterChatProvider(name string, factory ChatProviderFactory) {
	mu.Lock()
	defer mu.Unlock()
	chatFactories[name] = factory
}

// GetEmailFactory returns an email provider factory by name.
func GetEmailFactory(name string) (EmailProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()

	factory, exists := emailFactories[name]
	if !exists {
		return nil, fmt.Errorf("unknown email provider: %s (not registered)", name)
	}
	return factory, nil
}

// GetSMSFactory returns an SMS provider factory by name.
func GetSMSFactory(name string) (SMSProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()

	factory, exists := smsFactories[name]
	if !exists {
		return nil, fmt.Errorf("unknown SMS provider: %s (not registered)", name)
	}
	return factory, nil
}

// GetPushFactory returns a push provider factory by name.
func GetPushFactory(name string) (PushProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()

	factory, exists := pushFactories[name]
	if !exists {
		return nil, fmt.Errorf("unknown push provider: %s (not registered)", name)
	}
	return factory, nil
}

// GetChatFactory returns a chat provider factory by name.
func GetChatFactory(name string) (ChatProviderFactory, error) {
	mu.RLock()
	defer mu.RUnlock()

	factory, exists := chatFactories[name]
	if !exists {
		return nil, fmt.Errorf("unknown chat provider: %s (not registered)", name)
	}
	return factory, nil
}

// IsEmailProviderRegistered checks if an email provider is registered.
func IsEmailProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := emailFactories[name]
	return ok
}

// IsSMSProviderRegistered checks if an SMS provider is registered.
func IsSMSProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := smsFactories[name]
	return ok
}

// IsPushProviderRegistered checks if a push provider is registered.
func IsPushProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := pushFactories[name]
	return ok
}

// IsChatProviderRegistered checks if a chat provider is registered.
func IsChatProviderRegistered(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := chatFactories[name]
	return ok
}
