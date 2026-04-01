package mailgun

import (
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/registry"
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
