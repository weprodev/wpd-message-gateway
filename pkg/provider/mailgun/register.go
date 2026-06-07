package mailgun

import (
	"github.com/weprodev/wpd-message-gateway/pkg/contracts"
	"github.com/weprodev/wpd-message-gateway/pkg/registry"
)

func init() {
	registry.RegisterEmailProvider(ProviderName, func(cfg registry.EmailConfig) (contracts.EmailSender, error) {
		return New(Config{
			APIKey:    cfg.APIKey,
			Domain:    cfg.Domain,
			BaseURL:   cfg.BaseURL,
			FromEmail: cfg.FromEmail,
			FromName:  cfg.FromName,
		})
	})
}
