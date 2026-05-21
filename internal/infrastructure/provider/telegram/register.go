package telegram

import (
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/registry"
)

func init() {
	registry.RegisterOTPProvider(ProviderName, func(cfg registry.OTPConfig) (port.OTPSender, error) {
		return New(Config{
			APIKey:  cfg.APIKey,
			BaseURL: cfg.BaseURL,
		})
	})
}
