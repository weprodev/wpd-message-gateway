package telegram

import (
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/registry"
)

func init() {
	registry.RegisterPushProvider(ProviderName, func(cfg registry.PushConfig) (port.PushSender, error) {
		return New(Config{
			APIToken: cfg.APIKey,
			BaseURL:  cfg.BaseURL,
		})
	})
}
