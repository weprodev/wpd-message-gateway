package app

import (
	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

// ProviderFactory creates provider instances using the registry.
type ProviderFactory struct {
	cfg *Config
}

// NewProviderFactory creates a new provider factory.
func NewProviderFactory(cfg *Config) *ProviderFactory {
	return &ProviderFactory{
		cfg: cfg,
	}
}

// CreateEmailProvider creates an email provider by name.
func (f *ProviderFactory) CreateEmailProvider(name string) (port.EmailSender, error) {
	factory, err := registry.GetEmailFactory(name)
	if err != nil {
		return nil, err
	}

	cfg := f.cfg.EmailProviders[name]
	mailpit := registry.MailpitConfig{Enabled: f.cfg.Mailpit.Enabled}
	return factory(cfg, mailpit)
}

// CreateSMSProvider creates an SMS provider by name.
func (f *ProviderFactory) CreateSMSProvider(name string) (port.SMSSender, error) {
	factory, err := registry.GetSMSFactory(name)
	if err != nil {
		return nil, err
	}

	return factory(f.cfg.SMSProviders[name])
}

// CreatePushProvider creates a push provider by name.
func (f *ProviderFactory) CreatePushProvider(name string) (port.PushSender, error) {
	factory, err := registry.GetPushFactory(name)
	if err != nil {
		return nil, err
	}

	return factory(f.cfg.PushProviders[name])
}

// CreateChatProvider creates a chat provider by name.
func (f *ProviderFactory) CreateChatProvider(name string) (port.ChatSender, error) {
	factory, err := registry.GetChatFactory(name)
	if err != nil {
		return nil, err
	}

	return factory(f.cfg.ChatProviders[name])
}
