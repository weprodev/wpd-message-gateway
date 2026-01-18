package mailgun

import (
	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

func init() {
	registry.RegisterEmailProvider("mailgun", func(cfg registry.EmailConfig, _ registry.MailpitConfig) (port.EmailSender, error) {
		return New(Config{
			APIKey:    cfg.APIKey,
			Domain:    cfg.Domain,
			BaseURL:   cfg.BaseURL,
			FromEmail: cfg.FromEmail,
			FromName:  cfg.FromName,
		})
	})
}
