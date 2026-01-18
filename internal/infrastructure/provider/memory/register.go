package memory

import (
	"github.com/weprodev/wpd-message-gateway/internal/app/registry"
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
)

func init() {
	registry.RegisterEmailProvider("memory", func(cfg registry.EmailConfig, mailpit registry.MailpitConfig) (port.EmailSender, error) {
		mailpitCfg := MailpitConfig{Enabled: mailpit.Enabled}
		return NewEmailProvider(GetStore(), mailpitCfg), nil
	})

	registry.RegisterSMSProvider("memory", func(cfg registry.SMSConfig) (port.SMSSender, error) {
		return NewSMSProvider(GetStore()), nil
	})

	registry.RegisterPushProvider("memory", func(cfg registry.PushConfig) (port.PushSender, error) {
		return NewPushProvider(GetStore()), nil
	})

	registry.RegisterChatProvider("memory", func(cfg registry.ChatConfig) (port.ChatSender, error) {
		return NewChatProvider(GetStore()), nil
	})
}
