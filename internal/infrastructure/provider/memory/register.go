package memory

import (
	"github.com/weprodev/wpd-message-gateway/internal/core/port"
	"github.com/weprodev/wpd-message-gateway/internal/registry"
)

func init() {
	registry.RegisterEmailProvider("memory", func(cfg registry.EmailConfig) (port.EmailSender, error) {
		return NewEmailProvider(GetStore()), nil
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

	registry.RegisterOTPProvider("memory", func(cfg registry.OTPConfig) (port.OTPSender, error) {
		return NewChatProvider(GetStore()), nil
	})
}
