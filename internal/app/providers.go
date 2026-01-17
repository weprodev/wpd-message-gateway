package app

import (
	"fmt"

	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/mailgun"
	"github.com/weprodev/wpd-message-gateway/internal/infrastructure/provider/memory"
)

// ProviderFactory creates provider instances based on configuration.
type ProviderFactory struct {
	cfg   *Config
	store *memory.Store
}

// NewProviderFactory creates a new provider factory.
func NewProviderFactory(cfg *Config, store *memory.Store) *ProviderFactory {
	return &ProviderFactory{
		cfg:   cfg,
		store: store,
	}
}

// CreateEmailProvider creates an email provider by name.
func (f *ProviderFactory) CreateEmailProvider(name string) (port.EmailSender, error) {
	switch name {
	case "memory":
		mailpitCfg := memory.MailpitConfig{Enabled: f.cfg.Mailpit.Enabled}
		return memory.NewEmailProvider(f.store, mailpitCfg), nil

	case "mailgun":
		providerCfg, exists := f.cfg.EmailProviders[name]
		if !exists {
			return nil, fmt.Errorf("no configuration found for email provider: %s", name)
		}
		return mailgun.New(mailgun.Config{
			APIKey:    providerCfg.APIKey,
			Domain:    providerCfg.Domain,
			BaseURL:   providerCfg.BaseURL,
			FromEmail: providerCfg.FromEmail,
			FromName:  providerCfg.FromName,
		})

	default:
		return nil, fmt.Errorf("unknown email provider: %s", name)
	}
}

// CreateSMSProvider creates an SMS provider by name.
func (f *ProviderFactory) CreateSMSProvider(name string) (port.SMSSender, error) {
	switch name {
	case "memory":
		return memory.NewSMSProvider(f.store), nil

	default:
		return nil, fmt.Errorf("unknown SMS provider: %s", name)
	}
}

// CreatePushProvider creates a push notification provider by name.
func (f *ProviderFactory) CreatePushProvider(name string) (port.PushSender, error) {
	switch name {
	case "memory":
		return memory.NewPushProvider(f.store), nil

	default:
		return nil, fmt.Errorf("unknown push provider: %s", name)
	}
}

// CreateChatProvider creates a chat provider by name.
func (f *ProviderFactory) CreateChatProvider(name string) (port.ChatSender, error) {
	switch name {
	case "memory":
		return memory.NewChatProvider(f.store), nil

	default:
		return nil, fmt.Errorf("unknown chat provider: %s", name)
	}
}
