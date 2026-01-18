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

	cfg := f.toRegistryEmailConfig(name)
	mailpit := registry.MailpitConfig{Enabled: f.cfg.Mailpit.Enabled}
	return factory(cfg, mailpit)
}

// CreateSMSProvider creates an SMS provider by name.
func (f *ProviderFactory) CreateSMSProvider(name string) (port.SMSSender, error) {
	factory, err := registry.GetSMSFactory(name)
	if err != nil {
		return nil, err
	}

	cfg := f.toRegistrySMSConfig(name)
	return factory(cfg)
}

// CreatePushProvider creates a push provider by name.
func (f *ProviderFactory) CreatePushProvider(name string) (port.PushSender, error) {
	factory, err := registry.GetPushFactory(name)
	if err != nil {
		return nil, err
	}

	cfg := f.toRegistryPushConfig(name)
	return factory(cfg)
}

// CreateChatProvider creates a chat provider by name.
func (f *ProviderFactory) CreateChatProvider(name string) (port.ChatSender, error) {
	factory, err := registry.GetChatFactory(name)
	if err != nil {
		return nil, err
	}

	cfg := f.toRegistryChatConfig(name)
	return factory(cfg)
}

func (f *ProviderFactory) toRegistryEmailConfig(name string) registry.EmailConfig {
	appCfg := f.cfg.EmailProviders[name]
	return registry.EmailConfig{
		CommonConfig: registry.CommonConfig{
			APIKey:    appCfg.APIKey,
			APISecret: appCfg.APISecret,
			Region:    appCfg.Region,
			BaseURL:   appCfg.BaseURL,
			Extra:     appCfg.Extra,
		},
		Domain:    appCfg.Domain,
		FromEmail: appCfg.FromEmail,
		FromName:  appCfg.FromName,
	}
}

func (f *ProviderFactory) toRegistrySMSConfig(name string) registry.SMSConfig {
	appCfg := f.cfg.SMSProviders[name]
	return registry.SMSConfig{
		CommonConfig: registry.CommonConfig{
			APIKey:    appCfg.APIKey,
			APISecret: appCfg.APISecret,
			Region:    appCfg.Region,
			BaseURL:   appCfg.BaseURL,
			Extra:     appCfg.Extra,
		},
		FromPhone: appCfg.FromPhone,
	}
}

func (f *ProviderFactory) toRegistryPushConfig(name string) registry.PushConfig {
	appCfg := f.cfg.PushProviders[name]
	return registry.PushConfig{
		CommonConfig: registry.CommonConfig{
			APIKey:    appCfg.APIKey,
			APISecret: appCfg.APISecret,
			Region:    appCfg.Region,
			BaseURL:   appCfg.BaseURL,
			Extra:     appCfg.Extra,
		},
		AppID: appCfg.AppID,
		Topic: appCfg.Topic,
	}
}

func (f *ProviderFactory) toRegistryChatConfig(name string) registry.ChatConfig {
	appCfg := f.cfg.ChatProviders[name]
	return registry.ChatConfig{
		CommonConfig: registry.CommonConfig{
			APIKey:    appCfg.APIKey,
			APISecret: appCfg.APISecret,
			Region:    appCfg.Region,
			BaseURL:   appCfg.BaseURL,
			Extra:     appCfg.Extra,
		},
		FromPhone:  appCfg.FromPhone,
		WebhookURL: appCfg.WebhookURL,
	}
}
